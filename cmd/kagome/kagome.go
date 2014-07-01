package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ikawaha/kagome/tokenizer"
)

var usageMessage = "usage: kagome [-f input_file]\n"

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	fInputFile = flag.String("f", "", "input file")
)

func Main() {
	var fp = os.Stdin
	t := tokenizer.NewTokenizer()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		morphs, err := t.Tokenize(line)
		if err != nil {
			log.Println(err)
		}
		for _, m := range morphs {
			c, _ := m.Content()
			if m.Class == tokenizer.DUMMY {
				fmt.Printf("%s\n", m.Surface)
			} else {
				fmt.Printf("%s\t%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
					m.Surface, c.Pos, c.Pos1, c.Pos2, c.Pos3, c.Katuyougata, c.Katuyoukei, c.Kihonkei, c.Yomi, c.Pronunciation)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	Main()
}
