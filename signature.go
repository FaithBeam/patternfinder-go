package patternfinder

type signature struct {
	Name        string
	Pattern     []patternByte
	FoundOffset int
}

func (s signature) String() string {
	return s.Name
}
