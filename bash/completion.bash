# This comes from: https://git.zx2c4.com/password-store/tree/src/completion/pass.bash-completion

_loki_complete_entries () {
	prefix="${PASSWORD_STORE_DIR:-$HOME/.loki/}"
	prefix="${prefix%/}/"
	suffix=".loki"
	autoexpand=${1:-0}

	local IFS=$'\n'
	local items=($(compgen -f $prefix$cur))

	# Remember the value of the first item, to see if it is a directory. If
	# it is a directory, then don't add a space to the completion
	local firstitem=""
	# Use counter, can't use ${#items[@]} as we skip hidden directories
	local i=0

	for item in ${items[@]}; do
		[[ $item =~ /\.[^/]*$ ]] && continue

		# if there is a unique match, and it is a directory with one entry
		# autocomplete the subentry as well (recursively)
		if [[ ${#items[@]} -eq 1 && $autoexpand -eq 1 ]]; then
			while [[ -d $item ]]; do
				local subitems=($(compgen -f "$item/"))
				local filtereditems=( )
				for item2 in "${subitems[@]}"; do
					[[ $item2 =~ /\.[^/]*$ ]] && continue
					filtereditems+=( "$item2" )
				done
				if [[ ${#filtereditems[@]} -eq 1 ]]; then
					item="${filtereditems[0]}"
				else
					break
				fi
			done
		fi

		# append / to directories
		[[ -d $item ]] && item="$item/"

		item="${item%$suffix}"
		COMPREPLY+=("${item#$prefix}")
		if [[ $i -eq 0 ]]; then
			firstitem=$item
		fi
		let i+=1
	done

	# The only time we want to add a space to the end is if there is only
	# one match, and it is not a directory
	if [[ $i -gt 1 || ( $i -eq 1 && -d $firstitem ) ]]; then
		compopt -o nospace
	fi
}

_loki_complete_folders () {
	prefix="${PASSWORD_STORE_DIR:-$HOME/.loki/}"
	prefix="${prefix%/}/"

	local IFS=$'\n'
	local items=($(compgen -d $prefix$cur))
	for item in ${items[@]}; do
		[[ $item == $prefix.* ]] && continue
		COMPREPLY+=("${item#$prefix}/")
	done
}

_loki()
{
	COMPREPLY=()
	local cur="${COMP_WORDS[COMP_CWORD]}"
	local commands="search grep shutdown stop insert add login pw pass help ls list show import init change edit remove rm del copy cp move mv version ver complete"
	if [[ $COMP_CWORD -gt 1 ]]; then
		local lastarg="${COMP_WORDS[$COMP_CWORD-1]}"
		case "${COMP_WORDS[1]}" in
			ls|list|edit)
				_loki_complete_entries
				;;
			show|-*)
				_loki_complete_entries 1
				;;
			insert|add)
				_loki_complete_entries
				;;
			cp|copy|mv|rename)
				_loki_complete_entries
				;;
			rm|remove|delete)
				_loki_complete_entries
				;;
			git)
				COMPREPLY+=($(compgen -W "init push pull config log reflog rebase" -- ${cur}))
				;;
		esac

	else
		COMPREPLY+=($(compgen -W "${commands}" -- ${cur}))
		_loki_complete_entries 1
	fi
}

complete -o filenames -F _loki loki
