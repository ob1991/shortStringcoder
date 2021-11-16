package main

import (
	"fmt"
)

func main() {
	var shortStringCoder = NewShortStringCoder([]byte("1234567890abcdef"))
	str := "123456781234"
	encodeStr := shortStringCoder.Encode([]byte(str))
	decodeStr := shortStringCoder.Decode(encodeStr)
	fmt.Printf("%v,orin %v", string(decodeStr), str)
}
