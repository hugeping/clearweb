#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No url"
	exit 1
fi
curl -s "$1" | htmlcut -val '(post\_full|comment\_\_message|user\-info\_inline)' -regexp | htmlgmi -m -n 
