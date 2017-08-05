package score

import "reflect"

// NewWeighted returns a new weighted probability object
func NewWeighted() *Weighted {
	return &Weighted{}
}

type score struct {
	p float64
	w float64
}

// Weighted is a weighted probability
type Weighted struct {
	scores []score
}

// AddScore adds a new metric to the weighted probability
func (ws *Weighted) AddScore(p float64, w float64) {
	if p < 0 || p > 1 {
		panic("Probability must in percent")
	}
	if w < 0 {
		panic("Weights can't be negative")
	}
	ws.scores = append(ws.scores, score{p, w})
}

// AddEquals adds a new metric to the probability based on equality
func (ws *Weighted) AddEquals(a interface{}, b interface{}, w float64) {
	if reflect.DeepEqual(a, b) {
		ws.AddScore(1, w)
	} else {
		ws.AddScore(0, w)
	}
}

// Score calculates the score from the added metrics
func (ws *Weighted) Score() float64 {
	var totweight float64
	for _, s := range ws.scores {
		totweight += s.w
	}

	var totscore float64
	for _, s := range ws.scores {
		totscore += s.p * s.w
	}

	return totscore / totweight
}
