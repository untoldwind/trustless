#!/bin/sh

target=../secrets/generate/default_words.go
cat >$target <<END
package generate

var defaultWords = []string{
END

for i in $(cat wordlist2)
do
  echo "\"$i\"," >>$target
done

echo "}" >> $target
