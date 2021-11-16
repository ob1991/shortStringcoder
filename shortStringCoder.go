package main

type ShortStringCoder struct {
	encodeMap map[byte]byte
	decodeMap map[byte]byte
	size      int
	mask      map[int]byte
	b         int
	c         int
}

func MaxComDivisor(x int, y int) int {
	for {
		if x > y {
			x -= y
		} else if x < y {
			y -= x
		} else {
			return x
		}
	}
}
func NewShortStringCoder(seed []byte) *ShortStringCoder {
	shortStringCoder := ShortStringCoder{
		encodeMap: make(map[byte]byte),
		decodeMap: make(map[byte]byte),
		mask:      map[int]byte{7: 1, 6: 3, 5: 7, 4: 15, 3: 31, 2: 63, 1: 127},
	}
	for k := len(seed); k != 0; {
		shortStringCoder.size++
		k >>= 1
	}
	for k, v := range seed {
		shortStringCoder.encodeMap[v] = byte(k + 1)
		shortStringCoder.decodeMap[byte(k+1)] = v
	}
	a := MaxComDivisor(8, shortStringCoder.size)
	shortStringCoder.b = 8 / a
	shortStringCoder.c = shortStringCoder.size / a
	return &shortStringCoder
}
func (shortStringCoder *ShortStringCoder) Encode(content []byte) []byte {
	length := len(content)
	if length == 0 {
		return content
	}
	resultLen := length/shortStringCoder.b*shortStringCoder.c + length%shortStringCoder.b
	resultPtr := 0
	bytePtr := 0
	for k, v := range content {
		if k == 0 {
			content[0] = 0
		}
		if bytePtr == 0 {
			content[resultPtr] |= byte(shortStringCoder.encodeMap[v] << (8 - shortStringCoder.size))
			bytePtr = (bytePtr + 8 - shortStringCoder.size) % 8
		} else {
			if shortStringCoder.size > bytePtr {
				content[resultPtr] |= byte(shortStringCoder.encodeMap[v] >> (shortStringCoder.size - bytePtr))
				resultPtr++
				content[resultPtr] = 0
				if resultPtr >= resultLen {
					break
				}
				content[resultPtr] |= byte(shortStringCoder.encodeMap[v] << (8 - shortStringCoder.size + bytePtr))
				bytePtr = (bytePtr + 8 - shortStringCoder.size) % 8
			} else {
				content[resultPtr] |= byte(shortStringCoder.encodeMap[v] << (bytePtr - shortStringCoder.size))
				bytePtr = bytePtr - shortStringCoder.size
				if bytePtr == 0 {
					resultPtr++
					content[resultPtr] = 0
				}
			}
		}
	}
	return content[:resultLen]
}
func (shortStringCoder *ShortStringCoder) Decode(content []byte) []byte {
	length := len(content)
	if length == 0 {
		return content
	}
	var resultLength int
	if content[length-1]&shortStringCoder.mask[8-shortStringCoder.size] == 0 && length%shortStringCoder.c == 0 {
		resultLength = length/shortStringCoder.c*shortStringCoder.b + length%shortStringCoder.c - shortStringCoder.b + shortStringCoder.c
	} else {
		resultLength = length/shortStringCoder.c*shortStringCoder.b + length%shortStringCoder.c
	}
	result := make([]byte, resultLength)
	resultPtr := 0
	bytePtr := 0
	for i := 0; i < length && resultPtr < resultLength; {
		if bytePtr == 0 {
			result[resultPtr] |= shortStringCoder.decodeMap[content[i]>>(8-shortStringCoder.size)]
			bytePtr = (bytePtr + shortStringCoder.size) % 8
		} else {
			if shortStringCoder.size+bytePtr-8 > 0 {
				cur := (shortStringCoder.mask[bytePtr] & content[i]) << (shortStringCoder.size + bytePtr - 8)
				i++
				if i >= length {
					if resultPtr >= resultLength {
						break
					}
					result[resultPtr] = shortStringCoder.decodeMap[cur]
					break
				}
				cur |= content[i] >> (16 - shortStringCoder.size - bytePtr)
				bytePtr = (bytePtr + shortStringCoder.size) % 8
				result[resultPtr] = shortStringCoder.decodeMap[cur]
			} else {
				cur := (shortStringCoder.mask[bytePtr] & content[i]) >> (8 - shortStringCoder.size - bytePtr)
				bytePtr = bytePtr + shortStringCoder.size
				result[resultPtr] = shortStringCoder.decodeMap[cur]
			}
		}
		resultPtr++
	}
	return result
}
