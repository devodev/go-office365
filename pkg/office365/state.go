package office365

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
)

// State is an interface for storinm and retrievinm Watcher state.
type State interface {
	setLastContentCreated(*schema.ContentType, time.Time)
	getLastContentCreated(*schema.ContentType) time.Time
	setLastRequestTime(*schema.ContentType, time.Time)
	getLastRequestTime(*schema.ContentType) time.Time
	Read(io.Reader) error
	Write(io.Writer) error
}

// MemoryState is an in-memory State interface implementation.
type MemoryState struct {
	muCreated          *sync.RWMutex
	lastContentCreated map[schema.ContentType]time.Time
	muRequest          *sync.RWMutex
	lastRequestTime    map[schema.ContentType]time.Time
}

// NewMemoryState returns a new MemoryState.
func NewMemoryState() *MemoryState {
	return &MemoryState{
		muCreated:          &sync.RWMutex{},
		lastContentCreated: make(map[schema.ContentType]time.Time),
		muRequest:          &sync.RWMutex{},
		lastRequestTime:    make(map[schema.ContentType]time.Time),
	}
}

func (m *MemoryState) setLastContentCreated(ct *schema.ContentType, t time.Time) {
	m.muCreated.Lock()
	defer m.muCreated.Unlock()

	last, ok := m.lastContentCreated[*ct]
	if !ok || last.Before(t) {
		m.lastContentCreated[*ct] = t
	}
}

func (m *MemoryState) getLastContentCreated(ct *schema.ContentType) time.Time {
	m.muCreated.RLock()
	defer m.muCreated.RUnlock()

	t, ok := m.lastContentCreated[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

func (m *MemoryState) setLastRequestTime(ct *schema.ContentType, t time.Time) {
	m.muRequest.Lock()
	defer m.muRequest.Unlock()

	last, ok := m.lastRequestTime[*ct]
	if !ok || last.Before(t) {
		m.lastRequestTime[*ct] = t
	}
}

func (m *MemoryState) getLastRequestTime(ct *schema.ContentType) time.Time {
	m.muRequest.RLock()
	defer m.muRequest.RUnlock()

	t, ok := m.lastRequestTime[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

func (m *MemoryState) returnState() *StateData {
	m.muCreated.RLock()
	m.muRequest.RLock()
	defer m.muCreated.RUnlock()
	defer m.muRequest.RUnlock()

	return &StateData{
		LastContentCreated: m.lastContentCreated,
		LastRequestTime:    m.lastRequestTime,
	}
}

func (m *MemoryState) setState(b *StateData) {
	m.muCreated.Lock()
	m.muRequest.Lock()
	defer m.muCreated.Unlock()
	defer m.muRequest.Unlock()

	m.lastContentCreated = b.LastContentCreated
	m.lastRequestTime = b.LastRequestTime
}

// Read will decode json from a reader and populate its state.
func (m *MemoryState) Read(r io.Reader) error {
	decoder := json.NewDecoder(r)

	var blob StateData
	if err := decoder.Decode(&blob); err != nil {
		return err
	}
	m.setState(&blob)
	return nil
}

// Write will encode its state as json to a writer.
func (m *MemoryState) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)

	blob := m.returnState()
	if err := encoder.Encode(&blob); err != nil {
		return err
	}
	return nil
}

// StateData holds the internal state of MemoryState.
type StateData struct {
	LastContentCreated map[schema.ContentType]time.Time
	LastRequestTime    map[schema.ContentType]time.Time
}
