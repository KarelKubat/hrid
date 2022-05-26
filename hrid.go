package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/KarelKubat/hrid/id"
)

const (
	usage = `
This is hrid, the Human Readable ID converter.
Usage:
  hrid [FLAGS] NUMBER - generates a human readable ID and prints it on stdout
  hrid [FLAGS] -id ID - re-interprets the ID as a number and prints it on stdout
Try hrid -help for an overview of all flags. 
`
)

var (
	tokensFlag     = flag.String("tokens", id.Tokens, "conversion alphabet: first rune represents 0, second 1, etc.")
	lenFlag        = flag.Int("length", id.StringLen, "minimum length of generated IDs, set to 0 for no padding")
	ignoreCaseFlag = flag.Bool("ignorecase", id.IgnoreCase, "when true, casing is ignored when converting IDs to numbers")
	groupsizeFlag  = flag.Int("groupsize", id.GroupSize, "size of space-delimited groups in generated IDs, for better readability")

	idFlag      = flag.Bool("id", false, "when true, arguments are taken as IDs, default: numbers")
	verboseFlag = flag.Bool("verbose", false, "show options with which the converter is instantiated")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	opts := &id.Opts{
		Tokens:     *tokensFlag,
		StringLen:  *lenFlag,
		IgnoreCase: *ignoreCaseFlag,
		GroupSize:  *groupsizeFlag,
	}
	idConverter, err := id.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	if *verboseFlag {
		log.Printf("Converter options: %+v", *opts)
	}
	for _, a := range flag.Args() {
		if *idFlag {
			n, err := idConverter.ToNr(a)
			if err != nil {
				log.Printf("%v: not a valid ID: %v", a, err)
			} else {
				fmt.Println(n)
			}
		} else {
			u, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				log.Printf("%v: not a valid number: %v", a, err)
			} else {
				fmt.Println(idConverter.ToString(u))
			}
		}
	}
}
