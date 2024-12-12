if status is-interactive
    echo "use Ctrl + T to search through files"
    echo "use Ctrl + R to search through history"
    # Commands to run in interactive sessions can go here
    oh-my-posh init fish --config $HOME/.config/ohmyposh/zen.toml | source
    fzf --fish | source
	bind \cs '__ethp_commandline_toggle_sudo'
end

# if test -z (pgrep ssh-agent | string collect)
#    eval (ssh-agent -c) > /dev/null
#    set -Ux SSH_AUTH_SOCK $SSH_AUTH_SOCK
#    set -Ux SSH_AGENT_PID $SSH_AGENT_PID
#end

# Wasmer
export WASMER_DIR="/home/k0in/.wasmer"
[ -s "$WASMER_DIR/wasmer.sh" ] && source "$WASMER_DIR/wasmer.sh"


# printf '\eP$f{"hook": "SourcedRcFileForWarp", "value": { "shell": "fish" }}\x9c'
