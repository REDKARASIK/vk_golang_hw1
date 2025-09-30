// go build gen/codegen.go && ./codegen pack/unpack.go  pack/marshaller.go
package main

import "fmt"

// lets generate code for this struct
// cgen: binpack
type User struct {
	ID       int
	Login    string
	RealName string `cgen:"-"`
	Flags    int
}

type Avatar struct {
	ID  int
	Url string
}

var test = 42

func main() {
	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "a.sulaev", 16);
			print map { ord.", "  } split("", $b); '
	*/
	data := []byte{
		128, 36, 17, 0,

		8, 0, 0, 0,
		97, 46, 115, 117, 108, 97, 101, 118,

		16, 0, 0, 0,
	}

	u := User{}
	u.Unpack(data)
	fmt.Printf("Unpacked user %#v", u)
}
