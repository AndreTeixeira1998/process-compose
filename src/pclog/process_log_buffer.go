package pclog

import (
	"sync"
)

const (
	slack = 100
)

type ProcessLogBuffer struct {
	buffer   []string
	size     int
	observer PcLogObserver
	mx       sync.Mutex
}

func NewLogBuffer(size int) *ProcessLogBuffer {
	return &ProcessLogBuffer{
		size:     size,
		buffer:   make([]string, 0, size+slack),
		observer: nil,
	}
}

func (b *ProcessLogBuffer) Write(message string) {
	b.mx.Lock()
	defer b.mx.Unlock()
	b.buffer = append(b.buffer, message)
	if len(b.buffer) > b.size+slack {
		b.buffer = b.buffer[slack:]
	}
	if b.observer != nil {
		b.observer.AddLine(message)
	}

}

func (b *ProcessLogBuffer) GetLogRange(offsetFromEnd, limit int) []string {
	if len(b.buffer) == 0 {
		return []string{}
	}
	if offsetFromEnd < 0 {
		offsetFromEnd = 0
	}
	if offsetFromEnd > len(b.buffer) {
		offsetFromEnd = len(b.buffer)
	}

	if limit < 1 {
		limit = 0
	}
	if limit > len(b.buffer) {
		limit = len(b.buffer)
	}
	if offsetFromEnd+limit > len(b.buffer) {
		limit = len(b.buffer) - offsetFromEnd
	}
	if limit == 0 {
		return b.buffer[len(b.buffer)-offsetFromEnd:]
	}
	return b.buffer[len(b.buffer)-offsetFromEnd : offsetFromEnd+limit]
}

func (b *ProcessLogBuffer) GetLogLine(lineIndex int) string {
	if len(b.buffer) == 0 {
		return ""
	}

	if lineIndex >= len(b.buffer) {
		lineIndex = len(b.buffer) - 1
	}

	if lineIndex < 0 {
		lineIndex = 0
	}

	return b.buffer[lineIndex]
}

func (b *ProcessLogBuffer) GetLogLength() int {
	return len(b.buffer)
}

func (b *ProcessLogBuffer) GetLogsAndSubscribe(observer PcLogObserver) {
	b.mx.Lock()
	defer b.mx.Unlock()
	observer.SetLines(b.buffer)
	b.observer = observer
}

func (b *ProcessLogBuffer) UnSubscribe() {
	b.mx.Lock()
	defer b.mx.Unlock()
	b.observer = nil
}
