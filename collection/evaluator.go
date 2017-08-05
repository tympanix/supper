package collection

import (
	"github.com/Tympanix/supper/types"
	"github.com/xrash/smetrics"
)

// DefaultEvaluator uses string similarity to rate subtitles against media files
type DefaultEvaluator struct{}

// Evaluate determines how well the subtitle matches
func (e *DefaultEvaluator) Evaluate(f types.LocalMedia, s types.Subtitle) float32 {
	if s == nil || s.Meta() == nil {
		return 0
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

func (e *DefaultEvaluator) EvaluateMovie(media types.Movie, sub types.Movie) float32 {
	score := smetrics.WagnerFischer(media.MovieName(), sub.MovieName(), 1, 1, 2)
	return float32(score)
}

func (e *DefaultEvaluator) EvaluateEpisode(media types.Episode, sub types.Episode) float32 {
	score := smetrics.WagnerFischer(media.TVShow(), sub.TVShow(), 1, 1, 2)
	return float32(score)
}
