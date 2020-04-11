#!/bin/bash
# This script recursively finds .WAV files beneath the current directory 
# and adds them to a JSON formatted file for use with a github.com/aoeu/audio sampler.
buf=/tmp/sound_files.$$
out=/tmp/waves.json
echo -e '[\n' > $out
find . -name '*.[Ww][Aa][Vv]' -exec readlink -f {} \; > $buf
counter=0
while read f; do 
	echo -e '{\n"NoteNum" : '$counter',\n"FileName" : "'$f'"\n},\n' >> $out
	counter=$(($counter+1))
done <$buf
head -n -2 $out > $out.$$ && mv $out.$$ $out
echo -e '}\n]\n' >> $out
python -m json.tool $out > $out.$$ && mv $out.$$ $out 
rm $buf
mv /tmp/waves.json config/waves.json
