package score

import (
	"math"

	"github.com/tympanix/supper/types"
	"github.com/xrash/smetrics"
)

// DefaultEvaluator uses string similarity to rate subtitles against media files
type DefaultEvaluator struct{}

// Evaluate determines how well the subtitle matches
func (e *DefaultEvaluator) Evaluate(f types.LocalMedia, s types.Subtitle) float32 {
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
	prob := NewWeighted()
	score := smetrics.JaroWinkler(media.MovieName(), sub.MovieName(), 0.7, 4)
	tags := math.Min(float64(len(sub.AllTags())/len(media.AllTags())), 1)

	prob.AddScore(score, 3)
	prob.AddScore(tags, 2)
	prob.AddEquals(media.Group(), sub.Group(), 1)
	prob.AddEquals(media.Year(), sub.Year(), 1)
	prob.AddEquals(media.Quality(), sub.Quality(), 0.5)
	prob.AddEquals(media.Source(), sub.Source(), 0.75)
	prob.AddEquals(media.Codec(), sub.Codec(), 0.25)

	return float32(prob.Score())
}

func (e *DefaultEvaluator) EvaluateEpisode(media types.Episode, sub types.Episode) float32 {
	if media.Season() != sub.Season() || media.Episode() != sub.Episode() {
		return 0
	}
	prob := NewWeighted()
	show := smetrics.JaroWinkler(media.TVShow(), sub.TVShow(), 0.0, 1)

	prob.AddScore(show, 1)
	prob.AddEquals(media.Group(), sub.Group(), 1)
	prob.AddEquals(media.Quality(), sub.Quality(), 0.5)
	prob.AddEquals(media.Source(), sub.Source(), 0.75)
	prob.AddEquals(media.Codec(), sub.Codec(), 0.25)

	return float32(prob.Score())
}
