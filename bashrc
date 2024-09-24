# ~/.bashrc: executed by bash for non-login shells

# Enable color support for ls and other commands
if command -v dircolors >/dev/null; then
    test -r ~/.dircolors && eval "$(dircolors ~/.dircolors)"
fi

# Alias definitions
alias ll='ls -la'
alias gs='git status'
alias gd='git diff'
alias gc='git commit'
alias gl='git pull'
alias gp='git push'

# Add a custom prompt
PS1='\[\e[1;32m\]\u@\h:\[\e[1;34m\]\w\[\e[0m\]\$ '

# Enable programmable completion features (you might need to install bash-completion)
if [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
fi

# Custom PATH
export PATH="$HOME/bin:$PATH"

# Set default editor
export EDITOR='nano'

# Enable vi mode for command line
set -o vi

# Export additional environment variables if needed
export MY_VAR='some_value'
