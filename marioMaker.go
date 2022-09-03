package main

type marioMaker struct {
	drawingBoard
}

func (m *marioMaker) init() marioMaker {
	m.defaultInit()

	m.width = 320
	m.height = 180

	m.colorList = [][3]uint8{
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

	return *m
}
