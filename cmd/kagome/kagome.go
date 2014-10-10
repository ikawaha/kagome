package main

import (
	kagome "github.com/ikawaha/kagome"

	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var usageMessage = "usage: kagome [-f input_file] [-u userdic_file]"

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(0)
}

var (
	fInputFile   = flag.String("f", "", "input file")
	fUserDicFile = flag.String("u", "", "user dic")
)

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

	t := kagome.NewTokenizer()
	if *fUserDicFile != "" {
		if udic, err := kagome.NewUserDic(*fUserDicFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			t.SetUserDic(udic)
		}
	}

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := t.Tokenize(line)
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			c := tok.Features()
			if tok.Class == kagome.DUMMY {
				fmt.Printf("%s\n", tok.Surface)
			} else {
				fmt.Printf("%s\t%v\n", tok.Surface, strings.Join(c, ","))
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
