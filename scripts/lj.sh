#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No url"
	exit 1
fi
curl -s "$1" | htmlcut -type article -key class -val entry-content -contains | htmlgmi -m -n
