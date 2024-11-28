kubectl exec -it kafka-0 -- kafka-acls --bootstrap-server kafka-service:9092 --list | \
awk '
/principal=/ {
    user = gensub(/.*principal=User:([^,]+).*/, "\\1", "g");
    operation = gensub(/.*operation=([^,]+).*/, "\\1", "g");
    if (operation == "READ") read = "y";
    if (operation == "WRITE") write = "y";
}
/ResourcePattern/ && NR > 1 {
    if (user) print user "," name "," (read ? read : "") "," (write ? write : "");
    read = write = ""; name = gensub(/.*name=([^,]+).*/, "\\1", "g");
}
END {
    if (user) print user "," name "," (read ? read : "") "," (write ? write : "");
}'

######################################

kubectl exec -it kafka-0 -- kafka-acls --bootstrap-server kafka-service:9092 --list | \
awk '
/principal=/ {
    user = gensub(/.*principal=User:([^,]+).*/, "\\1", "g");
    operation = gensub(/.*operation=([^,]+).*/, "\\1", "g");
    if (operation == "READ") read = "read";
    if (operation == "WRITE") write = "write";
}
/ResourcePattern/ && NR > 1 {
    if (user) print user "," name "," (read ? read : "") "," (write ? write : "");
    read = write = ""; name = gensub(/.*name=([^,]+).*/, "\\1", "g");
}
END {
    if (user) print user "," name "," (read ? read : "") "," (write ? write : "");
}'
