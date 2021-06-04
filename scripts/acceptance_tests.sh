#!/usr/bin/sh

files=$(ls data/ -I "*exp")
red='\e[1;31m'
grn='\e[1;32m'
end='\e[0m'

for f in $files; do
  ./authorizer < "data/$f" | diff - "data/$f.exp" 1>/dev/null && printf "%s %bOK%b\n" "$f" "$grn" "$end" || printf "%s %bFAIL%b\n" "$f" "$red" "$end"
done