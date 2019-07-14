#!/bin/bash
FILE="users.txt"
KEY="enc.pem"
COUNTER=0
while read LINE; do
    COUNTER=$[COUNTER + 1]
    echo -ne "\\033[KPassword count [$COUNTER] Trying password [$LINE]\\r"
	openssl ec -in "enc.pem" -passin pass:200300 -outform DER | xxd -ps
done < $FILE
