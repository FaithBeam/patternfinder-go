package patternfinder

import (
	"strings"
	"sync"
)

type nibble struct {
	wildcard bool
	data     byte
}

type patternByte struct {
	n1 nibble
	n2 nibble
}

func Format(pattern string) string {
	length := len(pattern)
	var result strings.Builder
	runePattern := []rune(pattern)
	for i := 0; i < length; i++ {
		ch := runePattern[i]
		if ch >= '0' && ch <= '9' || ch >= 'A' && ch <= 'F' || ch >= 'a' && ch <= 'f' || ch == '?' {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func hexRuneToInt(ch rune) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	if ch >= 'A' && ch <= 'F' {
		return int(ch - 'A' + 10)
	}
	if ch >= 'a' && ch <= 'f' {
		return int(ch - 'a' + 10)
	}
	return -1
}

func Transform(pattern string) []patternByte {
	pattern = Format(pattern)
	length := len(pattern)
	if length == 0 {
		return nil
	}
	result := make([]patternByte, 0)
	if length%2 != 0 {
		pattern += "?"
		length++
	}
	var newByte patternByte
	j := 0
	runePattern := []rune(pattern)
	for i := 0; i < length; i++ {
		ch := runePattern[i]
		if ch == '?' { // wildcard
			if j == 0 {
				newByte.n1.wildcard = true
			} else {
				newByte.n2.wildcard = true
			}
		} else { // hex
			if j == 0 {
				newByte.n1.wildcard = false
				newByte.n1.data = byte(hexRuneToInt(ch) & 0xF)
			} else {
				newByte.n2.wildcard = false
				newByte.n2.data = byte(hexRuneToInt(ch) & 0xF)
			}
		}

		j++
		if j == 2 {
			j = 0
			result = append(result, newByte)
		}
	}
	return result
}

func Find(data []byte, pattern []patternByte) (bool, int) {
	var temp int
	return FindWithOffset(data, pattern, temp)
}

func matchByte(b byte, p *patternByte) bool {
	if !p.n1.wildcard { // if not a wildcard we need to compare the data.
		n1 := b >> 4
		if n1 != p.n1.data { // if the data is not equal b doesn't match p.
			return false
		}
	}
	if !p.n2.wildcard { // if not a wildcard we need to compare the data.
		n2 := b & 0xF
		if n2 != p.n2.data { // if the data is not equal b doesn't match p.
			return false
		}
	}
	return true
}

func FindWithOffset(data []byte, pattern []patternByte, offsetFound int) (bool, int) {
	offset := 0
	offsetFound = -1
	if data == nil || pattern == nil {
		return false, 0
	}
	patternSize := len(pattern)
	if len(data) == 0 || patternSize == 0 {
		return false, 0
	}

	pos := 0
	for i := offset; i < len(data); i++ {
		if matchByte(data[i], &pattern[pos]) { // check if the current data byte matches the current pattern byte
			pos++
			if pos == patternSize { // everything matched
				offsetFound = i - patternSize + 1
				return true, offsetFound
			}
		} else { // fix by Computer_Angel
			i -= pos
			pos = 0 // reset current pattern position
		}
	}

	return false, 0
}

func FindAll(data []byte, pattern []patternByte) (bool, []int) {
	offsetsFound := make([]int, 0)
	size := len(data)
	pos := 0
	for size > pos {
		result, offsetFound := FindWithOffset(data, pattern, pos)
		if result {
			offsetsFound = append(offsetsFound, offsetFound)
			pos = offsetFound + len(pattern)
		} else {
			break
		}
	}
	if len(offsetsFound) > 0 {
		return true, offsetsFound
	} else {
		return false, offsetsFound
	}
}

func Scan(data []byte, signatures []signature) []signature {
	var wg sync.WaitGroup
	var mu sync.Mutex
	found := make([]signature, 0)
	for _, s := range signatures {
		wg.Add(1)
		go func(data []byte, pattern []patternByte) {
			result, _ := Find(data, pattern)
			if result {
				defer wg.Done()
				mu.Lock()
				found = append(found, s)
				mu.Unlock()
			}
		}(data, s.Pattern)
	}
	wg.Wait()
	return found
}
