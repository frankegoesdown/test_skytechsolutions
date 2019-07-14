package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	cipher = "aes-128-cbc"
    notFind = "Couldn't find password in that file"
)

type results struct {
	password string
	fileName string
}


var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func argParse() (string, string) {
	wordlistPath := flag.String("wordlist", "keys/users.txt", "Wordlist to use.")
	encFile := flag.String("file", "keys/enc.pem", "File to decrypt. (Required)")

	flag.Parse()
	if *encFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	return *wordlistPath, *encFile
}

func printResults(info results) {
	if info.password != "" {
		log.Println(strings.Repeat("-", 50))
		log.Printf("Found password [ %s ] using [ %s ] algorithm!!\n", info.password, cipher)
		log.Println(strings.Repeat("-", 50))
	} else {
		log.Println(strings.Repeat("-", 50))
		log.Println(notFind)
		log.Println(strings.Repeat("-", 50))
		info.fileName = RandStringBytesMaskImprSrcSB(32)
		info.password = notFind
	}

	f, err := os.OpenFile(info.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(info.password); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Results you can find in: %s", info.fileName)
}

func crack(encFile string, wordlistPath string, found chan<- results) {
	cmdFormat := "openssl ec -in %s -passin pass:%s -outform PEM -%s -noout"

	// read big file
	inFile, err := os.Open(wordlistPath)
	defer inFile.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		select {
		default:
			word := scanner.Text()
			// Sprintf is not a very good choice for bruteforce
			// https://dev.to/pmalhaire/concatenate-strings-in-golang-a-quick-benchmark-4ahh
			cmd := fmt.Sprintf(cmdFormat, encFile, word, cipher)
			command := strings.Split(cmd, " ")
			_, err := exec.Command(command[0], command[1:]...).Output()

			if err == nil {
				fileName := RandStringBytesMaskImprSrcSB(32)
				found <- results{word, fileName}
				return
			}
		}
	}
}


func main() {
	wordlist, encryptedFile := argParse()
	println("Bruteforcing Started")
	var info results
	alreadyFound := false
	found := make(chan results)

	go crack(encryptedFile, wordlist, found)
Loop:
	for {
		select {
		case info = <-found:
			if !alreadyFound {
				alreadyFound = true
				break Loop
			}
		}
	}
	printResults(info)
}
