#! /bin/bash

FHASH=`md5sum $1`
while true; do
    NHASH=`md5sum $1`
    if [ "$NHASH" != "$FHASH" ]; then
        ./mdp -file $1
        FHASH=$NHASH
    fi
    sleep 5
done

# This script receives the name of the file you want to preview as an argument
# It calculates the checksum of this file every 5 seconds. If the result is different,
# from the previous one, the content of the file was changed, triggering the execution 
# of the markdown-preview tool to preview it.