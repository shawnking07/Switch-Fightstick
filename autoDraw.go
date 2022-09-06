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
	ink(im img, con *nscon.Controller) error
}

type drawingBoard struct {
	width     int
	height    int
	colorList []color.Color

	clickPerMove    int
	currentPosition [2]int
	colorState      int
}

func (d *drawingBoard) defaultInit() {
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
		*input++
		time.Sleep(50 * time.Millisecond)
		*input--
		time.Sleep(50 * time.Millisecond)
	}
}

func (d *drawingBoard) cursorInit(con *nscon.Controller) {
	con.Input.Stick.Left.X = -1
	con.Input.Stick.Left.Y = 1
	time.Sleep(2 * time.Second)
	con.Input.Stick.Left.X = 0
	con.Input.Stick.Left.Y = 0
	time.Sleep(50 * time.Millisecond)
}

func (d *drawingBoard) ink(im img, con *nscon.Controller) error {
	if !d.checkImgSize(im) {
		log.Println("Image size is not correct")
		//return errors.New("image size is not correct")
	}

	x, y := &d.currentPosition[0], &d.currentPosition[1]
	height, width := d.height, d.width

	log.Println("Init cursor")
	d.cursorInit(con)
	d.cursorInit(con)

	for i := 0; i < height; i += d.clickPerMove {
		for j := 0; j < width; j += d.clickPerMove {
			if i%2 == 0 {
				*x, *y = j, i
			} else {
				*x, *y = width-j-1, i
			}

			// Choose color
			colorIndexValue := im.data[*y][*x]
			if im.imageType == Colored {
				if colorIndexValue != d.colorState {
					colorChooseMoveTimes := colorIndexValue - d.colorState
					if colorChooseMoveTimes < 0 {
						setInputWithTimes(&con.Input.Button.ZL, -colorChooseMoveTimes)
					} else {
						setInputWithTimes(&con.Input.Button.ZR, colorChooseMoveTimes)
					}
					d.colorState = colorIndexValue
				}
			}

			// ignore white color
			if im.imageType == Colored || (im.imageType == BlackAndWhite && colorIndexValue != 0) {
				setInputWithTimes(&con.Input.Button.A, 1)
				log.Println("Drawing", *x, *y, colorIndexValue)
			}

			// move cursor
			// will click one more time when move to the edge, but it's ok
			if i%2 == 0 {
				setInputWithTimes(&con.Input.Dpad.Right, d.clickPerMove)
			} else {
				setInputWithTimes(&con.Input.Dpad.Left, d.clickPerMove)
			}
		}
		setInputWithTimes(&con.Input.Dpad.Down, d.clickPerMove)
	}

	return nil
}

// Color distance
func distance(a, b color.Color) float64 {
	r1, g1, b1, _ := a.RGBA()
	r2, g2, b2, _ := b.RGBA()
	return math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)))
}

func (d *drawingBoard) getNearestColorIndex(c color.Color) int {
	min := math.MaxFloat64
	minIndex := -1
	for i, color1 := range d.colorList {
		dist := distance(c, color1)
		if dist < min {
			min = dist
			minIndex = i
		}
	}
	return minIndex
}

func (d *drawingBoard) getColorIndex(c color.Color) int {
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
		colors = append(colors, d.colorList...)
	}

	di := dither.NewDitherer(colors)
	di.Matrix = dither.FloydSteinberg
	i = di.Dither(i)

	imgData := make([][]int, height)
	for y := 0; y < height; y++ {
		imgData[y] = make([]int, width)
		for x := 0; x < width; x++ {
			index := d.getColorIndex(i.At(x, y))
			imgData[y][x] = index
		}
	}

	return img{imgData, it, height, width}
}
