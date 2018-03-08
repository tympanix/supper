package score

import (
	"math"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
	"github.com/tympanix/supper/types"
	"github.com/xrash/smetrics"
)

const (
	qualityWeight = 0.50
	sourceWeight  = 0.75
	codecWeight   = 0.15
	groupWeight   = 0.33
)

const (
	missingMultiplier      = 1.75
	missingCodecMultiplier = 1.25
)

// DefaultEvaluator uses string similarity to rate subtitles against media files
type DefaultEvaluator struct{}

// Evaluate determines how well the subtitle matches
func (e *DefaultEvaluator) Evaluate(f types.Media, s types.Media) float32 {
	if s == nil || s.Meta() == nil {
		return 0.0
	}
	if _m, ok := f.TypeMovie(); ok {
		if _s, ok := s.TypeMovie(); ok {
			return e.EvaluateMovie(_m, _s)
		}
		return 0.0
	} else if _e, ok := f.TypeEpisode(); ok {
		if _s, ok := s.TypeEpisode(); ok {
			return e.EvaluateEpisode(_e, _s)
		}
		return 0.0
	} else {
		return 0.0
	}
}

// EvaluateMovie returns the matching score for a movie
func (e *DefaultEvaluator) EvaluateMovie(media types.Movie, sub types.Movie) float32 {
	if media.Year() > 0 && sub.Year() > 0 && media.Year() != sub.Year() {
		return 0.0
	}

	prob := NewWeighted()
	score := smetrics.JaroWinkler(media.MovieName(), sub.MovieName(), 0.7, 4)
	//tags := math.Min(float64(len(sub.AllTags())/len(media.AllTags())), 1)

	prob.AddScore(score, 0.5)
	//prob.AddScore(tags, 0.25)

	e.evaluateMetadata(prob, media, sub)
	prob.AddEquals(media.Year(), sub.Year(), groupWeight)

	return float32(prob.Score())
}

func (e *DefaultEvaluator) EvaluateEpisode(media types.Episode, sub types.Episode) float32 {
	if media.Season() != sub.Season() || media.Episode() != sub.Episode() {
		return 0
	}
	prob := NewWeighted()
	show := smetrics.JaroWinkler(media.TVShow(), sub.TVShow(), 0.0, 1)

	prob.AddScore(show, 0.5)
	e.evaluateMetadata(prob, media, sub)

	return float32(prob.Score())
}

func (e *DefaultEvaluator) evaluateMetadata(p *Weighted, media types.Metadata, sub types.Metadata) {
	p.AddEquals(media.Group(), sub.Group(), 1)

	if media.Quality() != quality.None {
		if sub.Quality() == quality.None {
			p.AddScore(0.0, missingMultiplier*qualityWeight)
		} else if media.Quality() == sub.Quality() {
			p.AddScore(1.0, qualityWeight)
		} else {
			diff := math.Abs(float64(media.Quality() - sub.Quality()))
			p.AddScore(0.0, diff*qualityWeight)
		}
	}

	if media.Source() != source.None {
		if sub.Source() == source.None {
			p.AddScore(0.0, missingMultiplier*sourceWeight)
		} else if media.Source() == sub.Source() {
			p.AddScore(1.0, sourceWeight)
		} else {
			diff := math.Abs(float64(media.Source() - sub.Source()))
			p.AddScore(0.0, diff*sourceWeight)
		}
	}

	if media.Codec() != codec.None {
		if sub.Codec() == codec.None {
			p.AddScore(0.0, missingCodecMultiplier*codecWeight)
		} else if media.Codec() == sub.Codec() {
			p.AddScore(1.0, codecWeight)
		} else {
			diff := math.Abs(float64(media.Codec() - sub.Codec()))
			p.AddScore(0.0, diff*codecWeight)
		}
	}
}
