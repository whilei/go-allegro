package allegro

/*
#cgo pkg-config: allegro-5.0
#include <allegro5/allegro.h>
*/
import "C"
import "container/list"

var displayList *list.List

func init() {
	displayList = list.New()
}

type Display struct {
	Width, Height int
	ptr *C.ALLEGRO_DISPLAY
}

type DisplayFlags int
const (
	Windowed DisplayFlags = C.ALLEGRO_WINDOWED
	Fullscreen DisplayFlags = C.ALLEGRO_FULLSCREEN
	FullscreenWindow DisplayFlags = C.ALLEGRO_FULLSCREEN_WINDOW
	Resizable DisplayFlags = C.ALLEGRO_RESIZABLE
	// TODO: add the rest of these flags
)

func CreateDisplay(w, h int) *Display {
	d := C.al_create_display(cint(w), cint(h))
	if d == nil {
		return nil
	}
	display := &Display{Width:w, Height:h, ptr:d}
	displayList.PushBack(display)
	return display
}

func (display *Display) Destroy() {
	for e := displayList.Front(); e != nil; e = e.Next() {
		if e.Value.(*Display) == display {
			displayList.Remove(e)
			break
		}
	}
	C.al_destroy_display(display.ptr)
}

// utility function to update struct values after the display changes,
// e.g. changing Width and Height after a resize
func (display *Display) Update() {
	display.Width = int(C.al_get_display_width(display.ptr))
	display.Height = int(C.al_get_display_height(display.ptr))
}

func FlipDisplay() {
	C.al_flip_display()
}

func GetNewDisplayFlags() DisplayFlags {
	return (DisplayFlags)(C.al_get_new_display_flags())
}

func SetNewDisplayFlags(flags DisplayFlags) {
	C.al_set_new_display_flags(cint(int(flags)))
}

func ResetDisplayFlags() {
	C.al_set_new_display_flags(cint(0))
}

func (display *Display) GetEventSource() *EventSource {
	source := C.al_get_display_event_source(display.ptr)
	if source == nil {
		return nil
	}
	return &EventSource{ptr:source}
}

func (display *Display) AcknowledgeResize() bool {
	success := gobool(C.al_acknowledge_resize(display.ptr))
	if success {
		display.Update()
	}
	return success
}

func (display *Display) SetWindowTitle(title string) {
	title_ := C.CString(title) ; defer FreeString(title_)
	C.al_set_window_title(display.ptr, title_)
}
