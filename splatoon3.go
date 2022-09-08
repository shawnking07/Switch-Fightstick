package main

import "image/color"

type splatoon3 struct {
	drawingBoard
}

func (s *splatoon3) init() splatoon3 {
	s.defaultInit()

	s.width = 320
	s.height = 120

	s.colorList = []color.Color{
		color.White,
		color.Black,
	}

	return *s
}
