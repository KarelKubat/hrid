// file: test/m1/main.go
// Example of an explicit converter.
package main

import (
	"fmt"
	"log"

	"github.com/KarelKubat/hrid/id"
)

func main() {
	converter, err := id.New(&id.Opts{
		Tokens:     "0123456789ABCDEF", // Hex converter
		StringLen:  8,                  // Pad IDs to 8 tokens if needed
		IgnoreCase: true,               // treat an `a` as `A`
		GroupSize:  4,                  // group by 4 tokens as in "DEAD BEEF"
	})
	if err != nil {
		log.Fatal(err)
	}

	str := converter.ToString(3735928559)
	fmt.Println("3735928559 as string is", str)

	nr, err := converter.ToNr(str)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str, "as number is", nr)
}
