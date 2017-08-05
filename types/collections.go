package types

// SubtitleCollection is a collection of subtitles
type SubtitleCollection interface {
	RemoveNotHI()
	RemoveHI()
	Add(Subtitle)
}
