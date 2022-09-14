package main

import (
	"github.com/mzyy94/nscon"
	"image/color"
)

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

func (s *splatoon3) ink(im img, con *nscon.Controller) error {
	err := s.commonInk(im, con)
	setInputWithTimes(&con.Input.Button.Minus, 1)
	return err
}
