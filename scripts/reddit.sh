#!/bin/sh
if [ "$#" -eq 0 ]; then
	echo "No reddit"
	exit 1
fi
curl -s -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36" https://www.reddit.com/r/"$1"/.rss | atom2rss | rss2gmi -h -r
