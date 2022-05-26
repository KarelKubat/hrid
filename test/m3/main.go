// file: test/m3/main.go
package main

import (
	"fmt"
	"log"

	"github.com/KarelKubat/hrid/conv"
)

func main() {
	converter, err := conv.New("abcdefgh", 0) // octal converter but with digits a-h, zero checksum runes
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("1234 as runes:", converter.ToRunes(1234))
	fmt.Println("1234 as string:", converter.ToString(1234))
	nr, err := converter.ToNr("cdcc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cdcc as number:", nr)
}
