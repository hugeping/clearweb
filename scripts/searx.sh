#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No query"
	exit 1
fi
q="$1"
page=1
if [ "$#" -gt 1 ]; then
	page="$2"
	exit 1
elif echo "$1" | grep -q -E '\/[0-9]+$'; then
	page=`echo "$1"| sed -e 's|^.*/\([0-9]\+\)$|\1|'`
	q=`echo "$1"| sed -e 's|^\(.*\)/[0-9]\+$|\1|'`
fi
curl -s -X POST -F 'q='"$q" -F 'q=category_general' -F 'pageno='"$page" https://search.fedi.life/search | htmlcut -val '(result.content|external.link|result.header)' -regexp | htmlgmi -l 0 -m -n
