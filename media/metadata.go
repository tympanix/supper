package media

import "github.com/Tympanix/supper/parse"

// Metadata provides release information for media
type Metadata struct {
	group   string
	codec   string
	quality string
	source  string
	tags    []string
}

// ParseMetadata generates meta data from a string
func ParseMetadata(tags string) Metadata {
	return Metadata{
		group:   parse.Group(tags),
		codec:   parse.Codec(tags),
		quality: parse.Quality(tags),
		source:  parse.Source(tags),
		tags:    parse.Tags(tags),
	}
}

// Group returns the release group
func (m Metadata) Group() string {
	return m.group
}

// Codec returns the codec
func (m Metadata) Codec() string {
	return m.codec
}

// Quality returns the quality of the media
func (m Metadata) Quality() string {
	return m.quality
}

// Source returns the source of the media
func (m Metadata) Source() string {
	return m.source
}

// AllTags returns all metadata as a list of tags
func (m Metadata) AllTags() []string {
	return m.tags
}
