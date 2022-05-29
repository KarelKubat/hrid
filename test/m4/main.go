// file: test/m4/main.go
package main

import (
	"fmt"
	"log"

	"github.com/KarelKubat/hrid/conv"
)

const (
	alphabet = "ABCDEFGH"
)

func main() {
	for checksumLen := uint(0); checksumLen < 8; checksumLen++ {
		converter, err := conv.New("ABCDEFGH", checksumLen) // Base-8 conversion, but using funky digits
		if err != nil {
			log.Fatal(err)
		}
		for nr := uint64(0); nr <= 15; nr++ {
			id := converter.ToString(nr)
			decoded, err := converter.ToNr(id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%2v yields ID %12q (with %v checksum digits) which decodes to %v\n", nr, id, checksumLen, decoded)
			if nr != decoded {
				log.Fatal("decoding failed")
			}
		}
	}
}
