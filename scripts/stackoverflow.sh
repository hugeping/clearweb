#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No url"
	exit 1
fi
curl -s "$1" | htmlcut -val '(js\-post\-body|answers\-subheader)' -regexp | htmlcut -val '(js\-filter)' -regexp -not | sed -e 's/^\(.*itemprop="text">\)$/\1<div><h2>====<\/h2><\/div>/g' | htmlgmi -l 0 -m -n
