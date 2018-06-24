package misc

// Tag is an enum representing miscellaneous media attributes
type Tag int

// List is a list of miscellaneous tags
type List []Tag

// Has returns true if tag t is in the list
func (l List) Has(t Tag) bool {
	for _, _t := range l {
		if _t == t {
			return true
		}
	}
	return false
}

const (
	// Video3D is a tag for media which in in stereoscophic 3D format
	Video3D Tag = iota
	// HC is a tag for hard coded media
	HC
	// DTS is a tag for the DTS multichannel audio technology
	DTS
	// DolbyDigital is a tag for the DolbyDigital multichannel audio technology
	DolbyDigital
	// AC3 is a tag for the Dolby Digital AC3 audio codec
	AC3
	// Extended is a tag for extended versions of the theatrical release
	Extended
)
