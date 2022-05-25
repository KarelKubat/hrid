# hrid

`hrid` stands for *Human Readable ID*. 

This is a Go package that you can include in your own code to generate IDs in string form from `uint64` numbers, or to reverse that and to decode a string into a number. There is also a stand-alone program `hrid`. Some examples:

```shell
# Convert a number to an ID. The output is space-separated into sets to improve readability.
$ hrid 9999999999999999999
08NH S30K 4XFY YY

# Reverse, the ID is interpreted without regard to casing and spaces.
hrid -id '08nhs 30k 4xf YYY'
9999999999999999999
```

Out-of-the-box defaults are applied that are meant to be as sane as possible for humans:

- The "alphabet" for the conversion is `0123456789ABCDEFGHKLMNPQRSTUVWXY`: digits and uppercase letters. This default tries to avoid tokens that are similar to one another: there is no I (looks as a 1), there is no O (looks as a 0), etc.
- Generated IDs (strings) are padded to a length of 14 runes, which plays well with the alphabet: you don't need more tokens to represent a `uint64`.
- Casing is ignored when converting an ID to a number; an `A` and an `a` are treated the same. This also plays well with the default alphabet (but would have to be turned off if you want to use an alphabet that has upper and lower case tokens).
- Generated IDs are split into groups of four runes for better readability.

All these settings can be programmatically overruled in the package `hrid/id`, or by the flags that `hrid` accepts (try `hrid -help`). As a silly example, here's a binary converter using the standard 0 and 1, or using smileys:

```shell
$ hrid -tokens=01  12345678
1011 1100 0110 0001 0100 1110

$ hrid -tokens=🥵😀  12345678
😀🥵😀😀 😀😀🥵🥵 🥵😀😀🥵 🥵🥵🥵😀 🥵😀🥵🥵 😀😀😀🥵
```

The alphabet is interpreted as follows:

- The first rune of the alphabet is always interpreted as the value zero, the second rune is always interpreted as the value 1, the third as 2, and so on,
- Therefore, the length of the alphabet is the "base" of the number system,
- Hence, `0123456789ABCDEF` defines a hexadecimal converter using the normally applicable digits 0-F.


## Package hrid/id

Package `hrid/id` exports global functions (that use the default alphabet, grouping etc.). Also there is a `New()` constructor to instantiate a converter using different values. This package calls `hrid/conv` for the actual conversions, and then applies padding and formatting and takes care of case-insensitivy (if so requested).

### Synopsis

```go
// file: test/m1.go
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
```

```go
// file: test/m2.go
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
```

## Package hrid/conv

This package is respoonsible for the actual conversion. It can be directly called from your program if you don't care about padding or grouping in the string representations.

### Synopsis

```go
/ file: test/m3.go
package main

import (
	"fmt"
	"log"

	"github.com/KarelKubat/hrid/conv"
)

func main() {
	converter, err := conv.New("abcdefgh") // octal converter but with digits a-h
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
```
