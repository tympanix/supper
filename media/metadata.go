package media

import (
	"encoding/json"
	"strings"

	"github.com/tympanix/supper/parse"
)

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

func (m Metadata) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Group   string `json:"group"`
		Codec   string `json:"codec"`
		Quality string `json:"quality"`
		Source  string `json:"source"`
	}{
		m.group,
		m.codec,
		m.quality,
		m.source,
	})
}

// String return a description of the metadata
func (m Metadata) String() string {
	return strings.Join([]string{
		m.Group(),
		m.Codec(),
		m.Quality(),
		m.Source(),
	}, ",")
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
