package main

import (
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/mzyy94/nscon"
	"image"
	"image/color"
	"log"
	"math"
	"time"
)

// img is a 2D array of color indexes
// data[x][y] is index of the colorList
type img struct {
	data      [][]int
	imageType ImageType
	height    int
	width     int
}

type ImageType uint8

const (
	BlackAndWhite ImageType = iota
	Colored
)

type drawing interface {
	ink(i img)
}

type MoveState uint8

const (
	MoveStateStart MoveState = iota
	MoveStateMoveToRight
	MoveStateMoveToLeft
	MoveStateMoveToDown
	MoveStateEnd
)

type drawingBoard struct {
	width     int
	height    int
	colorList [][3]uint8

	clickPerMove    int
	moveState       MoveState
	currentPosition [2]int
	colorState      int
}

func (d *drawingBoard) defaultInit() {
	d.moveState = MoveStateStart
	d.currentPosition = [2]int{0, 0}
	d.colorState = 0
	d.clickPerMove = 1
}

func (d *drawingBoard) checkImgSize(i img) bool {
	imgData := i.data
	if d.height != len(imgData) || d.width != len(imgData[0]) {
		return false
	}
	return true
}

func setInputWithTimes(input *uint8, times int) {
	for i := 0; i < times; i++ {
		setInput(input)
		time.Sleep(50 * time.Millisecond)
	}
}

func (d *drawingBoard) cursorInit(con *nscon.Controller) {
	con.Input.Stick.Left.X = -1
	con.Input.Stick.Left.Y = -1
	time.AfterFunc(2*time.Second, func() {
		con.Input.Stick.Left.X = 0
		con.Input.Stick.Left.Y = 0
	})
}

func (d *drawingBoard) ink(i img, con *nscon.Controller) {
	if !d.checkImgSize(i) {
		log.Println("Image size is not correct")
	}
	log.Println("Init cursor")
	d.cursorInit(con)
	for {
		x, y := &d.currentPosition[0], &d.currentPosition[1]
		height, width := d.height, d.width

		// position move
		switch d.moveState {
		case MoveStateEnd:
			log.Println("Drawing finished")
			break
		case MoveStateStart:
			log.Println("Drawing start")
			d.currentPosition = [2]int{0, 0}
			d.moveState = MoveStateMoveToRight
		case MoveStateMoveToRight:
			*x += d.clickPerMove

			setInputWithTimes(&con.Input.Dpad.Right, d.clickPerMove)

			if *x >= width-1 {
				d.moveState = MoveStateMoveToDown
			}
		case MoveStateMoveToLeft:
			*x -= d.clickPerMove

			setInputWithTimes(&con.Input.Dpad.Left, d.clickPerMove)

			if *x == 0 {
				d.moveState = MoveStateMoveToDown
			}
		case MoveStateMoveToDown:
			*y += d.clickPerMove

			setInputWithTimes(&con.Input.Dpad.Down, d.clickPerMove)

			if *y >= height-1 {
				d.moveState = MoveStateEnd
			} else if *x <= 0 {
				d.moveState = MoveStateMoveToRight
			} else if *x >= width-1 {
				d.moveState = MoveStateMoveToLeft
			}
		}

		// Choose color
		if i.imageType == Colored {
			colorIndex := i.data[*x][*y]
			if colorIndex != d.colorState {
				colorChooseMoveTimes := colorIndex - d.colorState
				if colorChooseMoveTimes < 0 {
					setInputWithTimes(&con.Input.Button.ZL, -colorChooseMoveTimes)
				} else {
					setInputWithTimes(&con.Input.Button.ZR, colorChooseMoveTimes)
				}
			}
		}

		// Click

		// ignore white color
		//if (i.imageType == BlackAndWhite && i.data[x][y] == 0) ||
		//	(i.imageType == Colored && d.colorList[i.data[x][y]] == [3]uint8{255, 255, 255}) {
		//	continue
		//}

		setInput(&con.Input.Button.A)

		log.Println("Drawing", *x, *y, d.colorState)
	}
}

// Euclidean distance
func distance(a, b [3]uint8) float64 {
	var sum float64
	for i := 0; i < 3; i++ {
		sum += math.Pow(float64(a[i]-b[i]), 2)
	}
	return math.Sqrt(sum)
}

func (d *drawingBoard) getNearestColorIndex(c [3]uint8) int {
	var minDistance float64
	var minIndex int
	for i, color1 := range d.colorList {
		d := distance(c, color1)
		if i == 0 || d < minDistance {
			minDistance = d
			minIndex = i
		}
	}
	return minIndex
}

func (d *drawingBoard) getColorIndex(c [3]uint8) int {
	for i, color1 := range d.colorList {
		if c == color1 {
			return i
		}
	}
	return -1
}

func (d *drawingBoard) convertToImg(i image.Image, it ImageType) img {
	bounds := i.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	colors := make([]color.Color, 0)

	switch it {
	case BlackAndWhite:
		colors = append(colors, color.White, color.Black)
	case Colored:
		for colorIndex := range d.colorList {
			r, g, b := d.colorList[colorIndex][0], d.colorList[colorIndex][1], d.colorList[colorIndex][2]
			rgba := color.RGBA{R: r, G: g, B: b, A: 255}
			colors = append(colors, rgba)
		}
	}

	di := dither.NewDitherer(colors)
	di.Matrix = dither.FloydSteinberg
	i = di.Dither(i)

	imgData := make([][]int, height)
	for y := 0; y < height; y++ {
		imgData[y] = make([]int, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := i.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			index := d.getNearestColorIndex([3]uint8{uint8(r), uint8(g), uint8(b)})
			imgData[y][x] = index
		}
	}

	return img{imgData, it, height, width}
}
