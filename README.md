# hrid: Human Readable IDs

<!-- toc -->
- [But why?](#but-why)
- [Overview](#overview)
- [Package hrid/id](#package-hridid)
  - [Synopsis](#synopsis)
- [Package hrid/conv](#package-hridconv)
  - [Synopsis](#synopsis-1)
  - [Checksumming](#checksumming)
- [Errors](#errors)
<!-- /toc -->

## But why?

![But why?](ButWhy.gif)

There are still some interfaces where computer-generated numbers need to be shown to humans (say on a paper bill), who dutifully need to copy them and eventually enter them in some UI for further processing. 

A typical example is an [IBAN](https://en.wikipedia.org/wiki/International_Bank_Account_Number) (international bank number). You receive your bill for whatever you bought and you need to enter this IBAN in your e-banking app to initiate a payment. This IBAN has a country code and up to 24 digits. The IBAN format already has some safeguards against human errors: there are two checksum digits, so that's good. 

> **But copying (up to) 24 digits is not an easy task for humans: the sheer number of digits go get right introduces the possibility of errors.**

So why only digits? That makes the length unnecessarily long: base-10 yields a way longer representation than base-16. But why stop at base-16? Why not base-20 or even higher?

> Why not "invent" an *alphabet* of digits that's long enough so that large numbers can be respresented using short sequences, whilst the alphabet is designed to avoid typos? How about these digits: `0123456789ABCDEFGHKLMNPQRTUVWXY`. This is a base-30-something notation which attempts to avoid typos. There is no `O` because it looks too much like a zero, there is no `I` because it looks too much like a one, etc.. And there are only uppercase letters, so humans can enter things using whatever casing they like - the computer can compensate.

## Overview

This is a Go package that you can include in your own code to generate IDs in string form from `uint64` numbers, or to reverse that and to decode a string into a number. There is also a stand-alone program `hrid`. Some examples:

```shell
# Convert a number to an ID. The generated output is space-separated into sets to improve readability.
$ hrid 9999999999999999999
CNH M74 XCQ Y4Q H24

# Reverse, the ID is interpreted without regard to casing and spaces.
$ hrid -id 'cn hm 74 xc qy 4q h24'
9999999999999999999
```

Out-of-the-box defaults are applied that are meant to be as sane as possible for humans:

- The "alphabet" for the conversion consists of digits and uppercase letters. This default tries to avoid tokens that are similar to one another: there is no I (looks as a 1), there is no O (looks as a 0), etc. See `id/id.go` for the actual value. (You can always supply a different alphabet for your conversions.)
- Each generated ID is appended with two checksum runes.
- Generated IDs (strings) are padded to a length of 13 runes, which plays well with the alphabet: you don't need more tokens to represent a `uint64`. With the two checksum runes this yields 15 runes (nicely separated into equal-length groups).
- Casing is ignored when converting an ID to a number; an `A` and an `a` are treated the same. This also plays well with the default alphabet (but would have to be turned off if you want to use an alphabet that has upper and lower case tokens).
- Generated IDs are split into groups of five for better readability.

All these settings can be programmatically overruled in the package `hrid/id`, or by the flags that `hrid` accepts (try `hrid -help`). As a silly example, here's a binary converter using the standard 0 and 1, or using smileys:

```shell
$ hrid -alphabet=01 -groupsize=4 12345678
1011 1100 0110 0001 0100 1110 00

$ hrid -alphabet=🥵😀 -groupsize=4 12345678
😀🥵😀😀 😀😀🥵🥵 🥵😀😀🥵 🥵🥵🥵😀 🥵😀🥵🥵 😀😀😀🥵 🥵🥵
```

In this example the default ID length of 13 (plus 2 for the checksum, making it 15) cannot be met. The binary representation needs more digits - here 26. `hrid` won't shorten an ID, but it may pad it up to the desired length.

The alphabet is interpreted as follows:

- The first rune of the alphabet is always interpreted as the value zero, the second rune is always interpreted as the value 1, the third as 2, and so on,
- Therefore, the length of the alphabet is the "base" of the number system,
- Hence, `0123456789ABCDEF` defines a hexadecimal converter using the normally applicable digits 0-F.

## Package hrid/id

Package `hrid/id` exports global functions (that use the default alphabet, grouping etc.). Also there is a `New()` constructor to instantiate a converter using different values. This package calls `hrid/conv` for the actual conversions, and then applies padding and formatting and takes care of case-insensitivy (if so requested).

### Synopsis for `hrid/id`

```go
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
		Alphabet:    "0123456789ABCDEF", // Hex converter
		StringLen:   8,                  // Pad IDs to 8 tokens if needed
		IgnoreCase:  true,               // treat an `a` as `A`
		GroupSize:   4,                  // group by 4 tokens as in "DEAD BEEF"
		ChecksumLen: 0,                  // Don't add a checksum runes when generating IDs
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
```

## Package hrid/conv

This package is responsible for the actual conversions (with checksums, if so requested). It can be directly called from your program if you don't care about padding, grouping or case-insensitivity in the string representations.

### Synopsis for `hrid/conv`

```go
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
```

### Checksumming

The checksum over an ID is computed and postifixed to the ID as follows:

- The checksum starts at zero.
- For each rune in the ID, its position in the alphabet (i.e., its numeric value) is added to the checksum, and then the checksum is resized to "fit" the base of the alphabet using a modulo operation.
- The resulting checksum rune is added to the ID.
- When the checksum length indicates that more than 1 checksum runes should be added, then the process repeats. I.e., the second checksum rune that is added represents the ID *and* the first checksum rune, the third checksum rune represents the ID plus the two checksum runes, etc..

Assuming that the alphabet is `ABCDEFGH`, then an `A` is the value zero, a `B` is the value one etc. (This is in fact a base-8 conversion, but with funky digits `A-H` instead of `0-7`.) When converting the number 14 to an ID, with 2 checksum runes, the following applies:

- 14 is represented as `BG` (check your octal converter, 14 decimal is 016 octal).
- The first checksum rune is `H`, because `B`=1 plus `G`=6 is 7, which still fits the octal system.
- The second checksum rune is `G`, because `B`=1 plus `G`=6 plus `H`=7 is 14, and 17%8=6, or `G`.
- The overall ID is then `BGHG`, with the last 2 runes representing the checksum.

An example is `test/m4/main.go`:

```go
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
	for checksumLen := 0; checksumLen < 8; checksumLen++ {
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
```

Abbreviated output:

```
14 yields ID         "BG" (with 0 checksum digits) which decodes to 14
14 yields ID        "BGH" (with 1 checksum digits) which decodes to 14
14 yields ID       "BGHG" (with 2 checksum digits) which decodes to 14 
14 yields ID      "BGHGE" (with 3 checksum digits) which decodes to 14
14 yields ID     "BGHGEA" (with 4 checksum digits) which decodes to 14
14 yields ID    "BGHGEAA" (with 5 checksum digits) which decodes to 14
```

## Errors

The following errors may be raised:

**Programming errors** (the converter can't work):

- *Alphabet too short*: The converter needs at least two runes to work with, which is a base-2 number system.
- *Token repeats*: Tokens in the conversion alphabet may not repeat. Note that this also depends on whether case insensitivity is requested: the alphabet `abcABC` is perfectly valid when case matters.

**User input errors** (the converter works, but can't decode this):

- *ID too short*: An ID must contain at least one rune that leads to a number, plus checksum runes (if checksumming applies). E.g., the ID `a` is only valid without checksumming. ID `ab` succeeds when no checksumming is requested, or when the checksum length is 1.
- *Checksum error*: The last runes of an ID, when taken as the checksum, don't match.
- *No such token*: An ID contains a token that's not in the conversion alphabet. E.g., given the alphabet `ABCD`, the ID `ZZZ` isn't valid.

`hrid` implements error handling where besides a description, a code is present that can be inspected. The codes are in `er/er.go`. For each returned error your code may inspect the `.Code` field to see whether this is a system error, or a user error. For example:

```go
/// file: test/m5/main.go
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
		Alphabet:    "0123456789ABCDEF",
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

	// Let's cause a broken converter that just can't work.
	converter, err = id.New(&id.Opts{
		Alphabet:    "0123456789ABCDE0F", // 0 repeats
		IgnoreCase:  true,
		ChecksumLen: 2,
	})
	checkError(err)
	// Output will be similar to:
	//   System error, the conversion will never ever work.
	//   Detail: TokenRepeatsError: 0 repeats in alphabet "0123456789ABCDE0F"
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
```
