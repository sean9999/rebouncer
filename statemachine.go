package rebouncer

import (
	"os"

	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/rjeczalik/notify"
)

// StateMachine is used to access all necessary methods and data
type StateMachine interface {
	Subscribe() chan NiceEvent
	Version() string
	Info() machineryInfo
	WatchDir(string)
}

// pass this in to the New() constructor
//
//	should contain any info we have at instantiation-time
type Config struct {
	BufferSize int
}

// The easiest way to create a new StateMachine
func New(config Config) StateMachine {
	m := machinery{
		NiceChannel: make(chan NiceEvent),
		machineId:   getMachineId(),
		bootId:      getBootId(),
		processId:   os.Getpid(),
		sessionId:   uuid.NewString(),
		bufferSize:  config.BufferSize,
	}
	return m
}

func (m machinery) WatchDir(dir string) {

	var fsEvents = make(chan notify.EventInfo, DefaultBufferSize)
	err := notify.Watch(dir+"/...", fsEvents, notify.All)
	if err != nil {
		panic(err)
	}
	go func() {
		for fsEvent := range fsEvents {
			NotifyEventInfoToNiceEvent(fsEvent, dir, m.NiceChannel)
		}
	}()
}

// an ID unique to this computer
func getMachineId() string {
	id, err := machineid.ProtectedID("rebouncer/" + version)
	if err != nil {
		id = "0"
	}
	return id
}

// returns the Linux & SystemD BootID if available, otherwise "0"
func getBootId() string {
	j, err := sdjournal.NewJournal()
	if err != nil {
		return "0"
	}
	defer j.Close()
	bootId, err := j.GetBootID()
	if err != nil {
		return "0"
	}
	return bootId
}

type machineryInfo struct {
	machineId  string
	bootId     string
	processId  int
	sessionId  string
	bufferSize int
}

type machinery struct {
	NiceChannel chan NiceEvent
	machineId   string
	bootId      string
	processId   int
	sessionId   string
	bufferSize  int
}

func (m machinery) Subscribe() chan NiceEvent {
	return m.NiceChannel
}
func (m machinery) Version() string {
	return version
}
func (m machinery) Info() machineryInfo {
	return machineryInfo{
		machineId:  m.machineId,
		bootId:     m.bootId,
		processId:  m.processId,
		sessionId:  m.sessionId,
		bufferSize: m.bufferSize,
	}
}
