// file: test/m5/main.go
// Example of an error checking.
package main

import (
	"fmt"
	"os"

	"github.com/KarelKubat/hrid/er"
	"github.com/KarelKubat/hrid/id"
)

func main() {
	converter, err := id.New(&id.Opts{
		Tokens:      "0123456789ABCDEF",
		IgnoreCase:  true,
		ChecksumLen: 2,
	})
	checkError(err)

	str := converter.ToString(3735928559)
	fmt.Println("3735928559 as string is", str)

	nr, err := converter.ToNr(str)
	checkError(err)
	fmt.Println(str, "as number is", nr)

	// Let's cause a user input error.
	_, err = converter.ToNr("ZAB5A")
	checkError(err)
	// Output will be similar to:
	//   Check your user input and retry.
	//   Detail: NoSuchTokenError: token Z not in alphabet "0123456789ABCDEF"
}

func checkError(err *er.Err) {
	if err == nil {
		return
	}
	// Find out what's wrong and issue a friendly message.
	var cause string
	if err.Code == er.IDTooShortError || err.Code == er.ChecksumError || err.Code == er.NoSuchTokenError {
		cause = "Check your user input and retry."
	} else {
		cause = "System error, the conversion will never ever work."
	}
	fmt.Fprintf(os.Stderr, "%s\nDetail: %s\n", cause, err)

	// At this point your program might abort, or ask to retry, or whatever.
}
