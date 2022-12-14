package main

import (
	"flag"
	"github.com/mzyy94/nscon"
	"image/png"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func setInput(input *uint8) {
	*input++
	time.AfterFunc(100*time.Millisecond, func() {
		*input--
	})
}

func main() {
	target := *flag.String("target", "/dev/hidg0", "hid device")
	imgPath := *flag.String("img", "print.png", "image path")
	flag.Parse()

	con := nscon.NewController(target)
	con.LogLevel = 0
	defer con.Close()
	err := con.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, 1)

	m := marioMaker{}
	m.init()

	sp := splatoon3{}
	sp.init()

	f, err := os.Open(imgPath)
	if err != nil {
		log.Println(err)
	}
	// load png
	i, err := png.Decode(f)
	if err != nil {
		log.Println(err)
	}

	// Set tty break for read keyboard input directly
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	go func() {
		for {
			os.Stdin.Read(buf)
			switch buf[0] {
			case 'a':
				setInput(&con.Input.Dpad.Left)
			case 'd':
				setInput(&con.Input.Dpad.Right)
			case 'w':
				setInput(&con.Input.Dpad.Up)
			case 's':
				setInput(&con.Input.Dpad.Down)
			case ' ':
				setInput(&con.Input.Button.B)
			case 0x0a: // Enter
				setInput(&con.Input.Button.A)
			case '.':
				setInput(&con.Input.Button.X)
			case '/':
				setInput(&con.Input.Button.Y)
			case 0x1b: // Escape
				setInput(&con.Input.Button.Home)
			case '`':
				setInput(&con.Input.Button.Capture)
			case '-':
				setInput(&con.Input.Button.ZL)
			case 'q':
				setInput(&con.Input.Button.L)
			case ']':
				setInput(&con.Input.Button.R)
			case '=':
				setInput(&con.Input.Button.ZR)
			case 'g':
				setInput(&con.Input.Button.Plus)
			case 'f':
				setInput(&con.Input.Button.Minus)
			case 'n': // Golden finger
				log.Printf("Golden finger: Print Super Mario Maker")
				im := m.convertToImg(i, Colored)
				err = m.ink(im, con)
				if err != nil {
					log.Println(err)
				}

			case 'm':
				log.Printf("Golden finger: Print Splatoon 3")
				im := sp.convertToImg(i, BlackAndWhite)
				err = sp.ink(im, con)
				if err != nil {
					log.Println(err)
				}

			default:
				log.Printf("unknown: %c = 0x%02x\n", buf[0], buf[0])
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		return
	}
}
