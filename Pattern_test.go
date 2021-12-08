package patternfinder

import (
	"fmt"
	"os"
	"testing"
)

var pattern = ""
var file = ""

func TestFind(t *testing.T) {
	patternByte := Transform(pattern)
	b, _ := os.ReadFile(file)
	result, offset := Find(b, patternByte)
	if result {
		fmt.Println(offset)
	}
}

func TestScan(t *testing.T) {
	b, _ := os.ReadFile(file)
	signatures := make([]signature, 0)
	sig := signature{
		Name:        "yes",
		Pattern:     Transform(pattern),
		FoundOffset: 0,
	}
	signatures = append(signatures, sig)
	foundSigs := Scan(b, signatures)
	fmt.Println(foundSigs)
}

func TestTransform(t *testing.T) {
	pattern := Transform(pattern)
	fmt.Println(pattern)
}
