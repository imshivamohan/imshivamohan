package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/exec"
)

var (
	clientset    *kubernetes.Clientset
	podNamespace string
	podName      string
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to the kubeconfig file")
	flag.StringVar(&podNamespace, "namespace", "default", "Namespace of the Kafka pod")
	flag.StringVar(&podName, "pod", "kafka-dev-0", "Name of the Kafka pod")
	flag.Parse()

	if *kubeconfig == "" {
		fmt.Println("Error: kubeconfig path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Load Kubernetes config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Set up REST API routes
	http.HandleFunc("/topics", handleTopics)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleTopics handles requests to the /topics endpoint
func handleTopics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// List topics
		topics, err := listTopicsInPod()
		if err != nil {
			http.Error(w, "Failed to list topics: "+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(topics)

	case "POST":
		// Create a new topic (expecting JSON payload with "topicName" field)
		var reqBody struct {
			TopicName string `json:"topicName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.TopicName == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := createTopicInPod(reqBody.TopicName)
		if err != nil {
			http.Error(w, "Failed to create topic: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Topic %s created", reqBody.TopicName)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listTopicsInPod executes the command in the pod to list Kafka topics
func listTopicsInPod() ([]string, error) {
	cmd := []string{
		"/bin/sh", "-c", "kafka-topics.sh --list --bootstrap-server $(cat /mnt/secrets/tls.sh)",
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNamespace).
		SubResource("exec").
		Param("container", "kafka").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	for _, arg := range cmd {
		req.Param("command", arg)
	}

	exec, err := exec.NewSPDYExecutor(clientset.RESTConfig(), "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	// Capture output
	output := &bytes.Buffer{}
	err = exec.Stream(exec.StreamOptions{
		IOStreams: exec.IOStreams{
			Out: output,
			ErrOut: os.Stderr,
		},
		Tty: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command in pod: %w", err)
	}

	// Parse topics from output
	topics := strings.Split(strings.TrimSpace(output.String()), "\n")
	return topics, nil
}

// createTopicInPod executes the command in the pod to create a new Kafka topic
func createTopicInPod(topicName string) error {
	cmd := []string{
		"/bin/sh", "-c", fmt.Sprintf("kafka-topics.sh --create --topic %s --bootstrap-server $(cat /mnt/secrets/tls.sh)", topicName),
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNamespace).
		SubResource("exec").
		Param("container", "kafka").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	for _, arg := range cmd {
		req.Param("command", arg)
	}

	exec, err := exec.NewSPDYExecutor(clientset.RESTConfig(), "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	err = exec.Stream(exec.StreamOptions{
		IOStreams: exec.IOStreams{
			Out: os.Stdout,
			ErrOut: os.Stderr,
		},
		Tty: false,
	})
	if err != nil {
		return fmt.Errorf("failed to execute command in pod: %w", err)
	}

	return nil
}



go run main.go -kubeconfig=/path/to/kubeconfig -namespace=kafka-namespace -pod=kafka-dev-0

