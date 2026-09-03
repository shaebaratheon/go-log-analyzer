package storage

import (
	"encoding/binary"
	"os"
	"sync"
)

type JournalRecord struct {
	Timestamp int64
	Length    uint32
	Payload   []byte
}

type JournalWriter struct {
	mu   sync.Mutex
	file *os.File
}

func NewJournalWriter(path string) (*JournalWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &JournalWriter{file: f}, nil
}

func (j *JournalWriter) Append(ts int64, payload []byte) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	header := make([]byte, 12)
	binary.BigEndian.PutUint64(header[0:8], uint64(ts))
	binary.BigEndian.PutUint32(header[8:12], uint32(len(payload)))

	if _, err := j.file.Write(header); err != nil {
		return err
	}
	_, err := j.file.Write(payload)
	return err
}

func (j *JournalWriter) Close() error {
	return j.file.Close()
}
