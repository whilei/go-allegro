package main

import (
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/image"
	"fmt"
	"os"
)

/*
// Without this, cgo complains about atexit not being defined.
// Linking explicitly against /usr/lib64/libc_nonshared.a might fix it
void atexit() {}
*/
import "C"

func main() {
	var (
		display *allegro.Display
		eventQueue *allegro.EventQueue
		running bool = true
	)

	if err := allegro.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	allegro.SetNewDisplayFlags(allegro.Windowed)
	if display = allegro.CreateDisplay(640, 480); display != nil {
		defer display.Destroy()
		display.SetWindowTitle("Hello, Go!")
	} else {
		fmt.Fprintf(os.Stderr, "failed to create display\n")
		return
	}

	if eventQueue = allegro.CreateEventQueue(); eventQueue != nil {
		defer eventQueue.Destroy()
	} else {
		fmt.Fprintf(os.Stderr, "failed to create event queue\n")
		return
	}

	image.Init()
	gordon := allegro.LoadBitmap("img/gordon-the-gopher.png") ; defer gordon.Destroy()
	eventQueue.RegisterEventSource(display.GetEventSource())

	allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
	allegro.FlipDisplay()

	for {
		event, found := eventQueue.WaitForEventUntil(allegro.NewTimeout(0.06))
		if found {
			switch event.Type {
				case allegro.DisplayCloseEvent:
					running = false
					break
			}
		}

		if !running {
			break
		}

		allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
		gordon.Draw(float32((display.Width-gordon.Width)/2), float32((display.Height-gordon.Height)/2), 0)
		allegro.FlipDisplay()
	}
}
