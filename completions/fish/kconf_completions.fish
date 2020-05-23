set -l kconf_commands add help list rm use version view

function __fish_kconf_complete_contexts
   kconf list | tr -d '*[:blank:]'
end

function __fish_kconf_complete_namespaces
   set -l cmd (commandline -opc)
   if [ (count $cmd) -lt 3 ]
      return 1
   end

   kubectl --context=$cmd[3] get namespaces -o custom-columns=:metadata.name
end

complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -a "add" \
  -d "Add in a new kubeconfig file and optional context name"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "help" \
  -d "Help about any command"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "list" \
  -d "Lists available contexts"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "rm" \
  -d "Remove a kubeconfig from main file"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "use" \
  -d "Set the current context"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "version" \
  -d "Print version"
complete -c kconf -n "not __fish_seen_subcommand_from $kconf_commands" \
  -f \
  -a "view" \
  -d "View a specific context's config"

# use, view, and rm take context arguments
complete -c kconf \
  -f \
  -n "__fish_seen_subcommand_from use view rm" \
  -a "(__fish_kconf_complete_contexts)"

# use allows a namespace option
complete -c kconf \
  -x \
  -n "__fish_seen_subcommand_from use" \
  -s "n" \
  -a "(__fish_kconf_complete_namespaces)" \
  -d "set a namespace to use"
