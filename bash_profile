# ~/.bash_profile: executed by bash for login shells

# Load .bashrc if it exists
if [ -f ~/.bashrc ]; then
    . ~/.bashrc
fi

# Set the default language to UTF-8
export LANG=en_US.UTF-8

# Set history options
export HISTCONTROL=ignoredups:erasedups  # No duplicate entries
export HISTSIZE=10000                     # Increase history size
export HISTFILESIZE=20000                 # Increase history file size
shopt -s histappend                        # Append to history, don't overwrite

# Set the terminal type
export TERM=xterm-256color
