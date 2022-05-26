// file: test/m2/main.go
// Example of the implicit converter that uses defaults.
package main

import (
	"fmt"
	"log"

	"github.com/KarelKubat/hrid/id"
)

func main() {
	str := id.ToString(3735928559)
	fmt.Println("3735928559 as string is", str)

	nr, err := id.ToNr(str)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str, "as number is", nr)
}
