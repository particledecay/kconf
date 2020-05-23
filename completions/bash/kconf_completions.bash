__bash_kconf_commands () {
   echo add help list rm use version view
}

__contains_word () {
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
}

__bash_kconf_complete_contexts () {
   kconf list |tr -d '*[:blank:]' | tr -s '[:space:]'
}

__bash_kconf_complete () {
   local cmd=${COMP_WORDS[COMP_CWORD]}
   local i verb comps

   for ((i=0; i < COMP_CWORD; i++)); do
      if __contains_word "${COMP_WORDS[i]}" $(__bash_kconf_commands); then
         verb=${COMP_WORDS[i]}
         break
      fi
   done

   if [[ -z $verb ]]; then
      comps=$(__bash_kconf_commands)
   else
      case $verb in
         add)
            COMPREPLY=( $(compgen -o plusdirs -f -- "$cmd") )
            return 0
            ;;
         help)
            ;;
         list)
            ;;
         rm)
            comps=$(__bash_kconf_complete_contexts)
            ;;
         use)
            comps=$(__bash_kconf_complete_contexts)
            ;;
         version)
            ;;
         view)
            comps=$(__bash_kconf_complete_contexts)
            ;;
      esac
   fi

   COMPREPLY=( $(compgen -W "$comps" -- "$cmd") )
   return 0

}

complete -F __bash_kconf_complete kconf
