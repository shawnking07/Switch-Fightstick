package main

import (
	"github.com/mzyy94/nscon"
	"image/color"
)

type marioMaker struct {
	drawingBoard
}

func (m *marioMaker) init() marioMaker {
	m.defaultInit()

	m.width = 320
	m.height = 180

	colorList := [][3]uint8{
		{0, 0, 0},
		{255, 255, 255},
		{255, 0, 0},
		{180, 0, 0},
		{255, 244, 202},
		{163, 117, 62},
		{255, 255, 0},
		{255, 186, 0},
		{0, 255, 0},
		{0, 180, 0},
		{0, 255, 255},
		{0, 0, 255},
		{180, 88, 255},
		{125, 0, 180},
		{255, 185, 254},
		{180, 0, 127},
		{180, 180, 180},
	}

	// convert uint8 color to Color
	for _, c := range colorList {
		m.colorList = append(m.colorList, color.RGBA{R: c[0], G: c[1], B: c[2], A: 255})
	}

	return *m
}

func (m *marioMaker) ink(im img, con *nscon.Controller) error {
	return m.commonInk(im, con)
}
