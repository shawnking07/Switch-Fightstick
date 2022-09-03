package main

type splatoon3 struct {
	drawingBoard
}

func (s *splatoon3) init() splatoon3 {
	s.width = 320
	s.height = 120

	return *s
}
