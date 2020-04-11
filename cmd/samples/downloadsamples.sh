#!/bin/sh
set -e
which sox || echo 'Install sox (http://sox.sourceforge.net/soxformat.html), it isn't in your path.'
which wget || echo 'Install wget, it isn't in your path.'
tmpdir=/tmp/samples.$$
outfolder=808_samples
mkdir $tmpdir
prevdir=$PWD
cd $tmpdir
wget http://surachai.org/thedeepelement/rolandtr808/808normalized.zip
unzip $tmpdir/808normalized.zip -d .
rm -rf '__MACOSX'
mkdir $outfolder
find . -name '*.aif' > samplelist.txt
while read f; do 
	s=`basename "$f" | sed 's/ //g'`
	echo sox "$f" $outfolder/${s%.*}.wav
	sox -G -v 0.99 "$f" -b 16 $outfolder/${s%.*}.wav channels 2 rate 44100
done <samplelist.txt
cd $prevdir
cp -r $tmpdir/$outfolder .
rm -rf $tmpdir
