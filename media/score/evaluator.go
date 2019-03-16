package score

import (
	"math"

	"github.com/tympanix/supper/media/meta/codec"
	"github.com/tympanix/supper/media/meta/misc"
	"github.com/tympanix/supper/media/meta/quality"
	"github.com/tympanix/supper/media/meta/source"
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
	missingMultiplier     = 2.25
	unavailableMultiplier = 0.18
	diffMultiplier        = 0.66
)

func diffVal(f float64) float64 {
	return math.Sqrt(f)
}

// DefaultEvaluator uses string similarity to rate subtitles against media files
type DefaultEvaluator struct{}

// Evaluate determines how well the subtitle matches
func (e *DefaultEvaluator) Evaluate(f types.Media, s types.Media) float32 {
	if s == nil || s.Meta() == nil {
		return 0.0
	}
	if !f.Similar(s) {
		return 0.0
	}
	if !f.Meta().Misc().Has(misc.Video3D) && s.Meta().Misc().Has(misc.Video3D) {
		return 0.0
	}
	if _m, ok := f.TypeMovie(); ok {
		if _s, ok := s.TypeMovie(); ok {
			return e.evaluateMovie(_m, _s)
		}
		return 0.0
	} else if _e, ok := f.TypeEpisode(); ok {
		if _s, ok := s.TypeEpisode(); ok {
			return e.evaluateEpisode(_e, _s)
		}
		return 0.0
	} else {
		return 0.0
	}
}

// EvaluateMovie returns the matching score for a movie
func (e *DefaultEvaluator) evaluateMovie(media types.Movie, sub types.Movie) float32 {
	prob := NewWeighted()
	score := smetrics.JaroWinkler(media.MovieName(), sub.MovieName(), 0.7, 4)

	prob.AddScore(score, 0.5)
	e.evaluateMetadata(prob, media, sub)

	return float32(prob.Score())
}

func (e *DefaultEvaluator) evaluateEpisode(media types.Episode, sub types.Episode) float32 {
	prob := NewWeighted()
	show := smetrics.JaroWinkler(media.TVShow(), sub.TVShow(), 0.7, 4)

	prob.AddScore(show, 0.5)
	e.evaluateMetadata(prob, media, sub)

	return float32(prob.Score())
}

func (e *DefaultEvaluator) evaluateMetadata(p *Weighted, media types.Metadata, sub types.Metadata) {
	if media.Group() != "" {
		p.AddEquals(media.Group(), sub.Group(), groupWeight)
	} else {
		if sub.Group() == "" {
			p.AddScore(0.0, unavailableMultiplier*groupWeight)
		}
	}

	if sub.Quality() == quality.None {
		p.AddScore(0.0, missingMultiplier*qualityWeight)
	}
	if media.Quality() != quality.None {
		if media.Quality() == sub.Quality() {
			p.AddScore(1.0, qualityWeight)
		} else {
			diff := math.Abs(float64(media.Quality() - sub.Quality()))
			p.AddScore(0.0, diffVal(diff)*diffMultiplier*qualityWeight)
		}
	} else {
		if sub.Quality() == quality.None {
			// unavailable quality, apply penalty
			p.AddScore(0.0, unavailableMultiplier*qualityWeight)
		} else {
			// not comparable, favour 720p
			diff := math.Abs(float64(sub.Quality() - quality.HD720p))
			p.AddScore(0.0, diffVal(diff)*unavailableMultiplier*diffMultiplier*qualityWeight)
		}
	}

	if sub.Source() == source.None {
		p.AddScore(0.0, missingMultiplier*sourceWeight)
	}
	if media.Source() != source.None {
		if media.Source() == sub.Source() {
			p.AddScore(1.0, sourceWeight)
		} else {
			diff := math.Abs(float64(media.Source() - sub.Source()))
			p.AddScore(0.0, diffVal(diff)*diffMultiplier*sourceWeight)
		}
	} else {
		if sub.Source() == source.None {
			p.AddScore(0.0, unavailableMultiplier*sourceWeight)
		} else {
			diff := math.Abs(float64(sub.Source() - source.BluRay))
			p.AddScore(0.0, diffVal(diff)*unavailableMultiplier*diffMultiplier*sourceWeight)
		}
	}

	if sub.Codec() == codec.None {
		p.AddScore(0.0, missingMultiplier*codecWeight)
	}
	if media.Codec() != codec.None {
		if media.Codec() == sub.Codec() {
			p.AddScore(1.0, codecWeight)
		} else {
			diff := math.Abs(float64(media.Codec() - sub.Codec()))
			p.AddScore(0.0, diffVal(diff)*diffMultiplier*codecWeight)
		}
	} else {
		if sub.Codec() == codec.None {
			p.AddScore(0.0, unavailableMultiplier*codecWeight)
		}
	}
}
