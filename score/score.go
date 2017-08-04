package score

type Metric interface {
	Compute(interface{}, interface{})
}

type ScoreBoard struct {
	Metric
	list []Entry
}

func (s *ScoreBoard) Add(entry Entry) {
	s.list = append(s.list, entry)
}

type Entry interface {
	Score() float32
}
