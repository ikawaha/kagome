package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ikawaha/kagome/dic"
	"github.com/ikawaha/kagome/tokenizer"
)

var usageMessage = "usage: kagome [-f input_file] [-u userdic_file]"

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	fInputFile   = flag.String("f", "", "input file")
	fUserDicFile = flag.String("u", "", "input file")
)

type Item string

func (this Item) String() string {
	if this == "" {
		return "*"
	}
	return string(this)
}

func Main() {
	var inputFile = os.Stdin
	if *fInputFile != "" {
		var err error
		inputFile, err = os.Open(*fInputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer inputFile.Close()
	}

	t := tokenizer.NewTokenizer()
	if *fUserDicFile != "" {
		userDicFile, err := os.Open(*fUserDicFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer userDicFile.Close()
		if udic, err := dic.NewUserDic(userDicFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			t.SetUserDic(udic)
		}
	}

	scanner := bufio.NewScanner(inputFile)
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
					m.Surface, Item(c.Pos), Item(c.Pos1), Item(c.Pos2), Item(c.Pos3),
					Item(c.Katuyougata), Item(c.Katuyoukei), Item(c.Kihonkei), Item(c.Yomi), Item(c.Pronunciation))
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
