#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No query"
	exit 1
fi

curl -s -X POST -F 'q='"$1" -F 'q=category_general' -F 'pageno=1' https://search.fedi.life/search | htmlcut -val '(result.content|external.link|result.header)' -regexp | htmlgmi -l 0 -m -n
