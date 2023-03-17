package rebouncer

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rjeczalik/notify"
	"golang.org/x/sys/unix"
)

// only the inotify events we care about
const WatchMask = notify.InModify |
	notify.InCloseWrite |
	notify.InMovedFrom |
	notify.InMovedTo |
	notify.InCreate |
	notify.InDelete |
	notify.InDeleteSelf |
	notify.InMoveSelf

// the string ends with a "~" character
func endsInTilde(s string) bool {
	pattern := `~$`
	r, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	return r
}

// the string is just numbers
func containsOnlyNumbers(s string) bool {
	pattern := `^\d+$`
	r, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	return r
}

func isTempFile(path string) bool {
	return endsInTilde(path) && containsOnlyNumbers(path)
}

func NotifyToNiceEvent(ei notify.EventInfo, path string) NiceEvent {

	//	strip out unnecessary parts of the path, because the root (watchDir) is known
	abs, _ := filepath.Abs(path)
	normalFile := strings.TrimPrefix(ei.Path(), abs+"/")

	//	original low-level event
	data := ei.Sys().(*unix.InotifyEvent)

	e := NewNiceEvent("rebouncer/incoming/inotify")
	e.File = normalFile
	e.Operation = ei.Event().String()
	e.TransactionId = data.Cookie

	return e

}

// launch an inotify process that converts to NiceEvents and appends to machinery.batch
func (m *machinery) WatchDir(dir string) {
	var fsEvents = make(chan notify.EventInfo, DefaultBufferSize)
	err := notify.Watch(dir+"/...", fsEvents, notify.All)
	if err != nil {
		panic(err)
	}
	go func() {
		for fsEvent := range fsEvents {
			//m.batch = append(m.batch, NotifyToNiceEvent(fsEvent, dir))
			m.Injest(NotifyToNiceEvent(fsEvent, dir))
		}
	}()
}
