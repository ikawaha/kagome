package main

import (
	"github.com/ikawaha/kagome/tokenizer"

	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	var fp *os.File
	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		var err error
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		morphs, err := tokenizer.Tokenize(line)
		if err != nil {
			log.Println(err)
		}
		for i, m := range morphs {
			content, _ := m.Content()
			fmt.Printf("%3d, %v(%v, %v)\t%v\n", i, m.Surface, m.Start, m.End, content)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
