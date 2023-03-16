package rebouncer

import (
	"os"

	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
)

const version = "v0.0.1"
const DefaultBufferSize = 1024

// an ID unique to this machine
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

type UUID string

type Config struct {
	BufferSize int
}

type machineryInfo struct {
	machineId  string
	bootId     string
	processId  int
	sessionId  string
	bufferSize int
}

type StateMachine interface {
	Subscribe() chan NiceEvent
	Version() string
	Info() machineryInfo
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
