package office365

import (
	"encoding/gob"
	"io"
	"sync"
	"time"
)

// State is an interface for storing and retrieving Watcher state.
type State interface {
	setLastContentCreated(*ContentType, time.Time)
	getLastContentCreated(*ContentType) time.Time
	setLastRequestTime(*ContentType, time.Time)
	getLastRequestTime(*ContentType) time.Time
}

// MemoryState is an in-memory State interface implementation.
type MemoryState struct {
	muCreated          *sync.RWMutex
	lastContentCreated map[ContentType]time.Time
	muRequest          *sync.RWMutex
	lastRequestTime    map[ContentType]time.Time
}

// NewMemoryState returns a new MemoryState.
func NewMemoryState() *MemoryState {
	return &MemoryState{
		muCreated:          &sync.RWMutex{},
		lastContentCreated: make(map[ContentType]time.Time),
		muRequest:          &sync.RWMutex{},
		lastRequestTime:    make(map[ContentType]time.Time),
	}
}

func (m *MemoryState) setLastContentCreated(ct *ContentType, t time.Time) {
	m.muCreated.Lock()
	defer m.muCreated.Unlock()

	last, ok := m.lastContentCreated[*ct]
	if !ok || last.Before(t) {
		m.lastContentCreated[*ct] = t
	}
}

func (m *MemoryState) getLastContentCreated(ct *ContentType) time.Time {
	m.muCreated.RLock()
	defer m.muCreated.RUnlock()

	t, ok := m.lastContentCreated[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

func (m *MemoryState) setLastRequestTime(ct *ContentType, t time.Time) {
	m.muRequest.Lock()
	defer m.muRequest.Unlock()

	last, ok := m.lastRequestTime[*ct]
	if !ok || last.Before(t) {
		m.lastRequestTime[*ct] = t
	}
}

func (m *MemoryState) getLastRequestTime(ct *ContentType) time.Time {
	m.muRequest.RLock()
	defer m.muRequest.RUnlock()

	t, ok := m.lastRequestTime[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

// GOBState is an in-memory State interface implementation, but
// also provides Read and Write methods for serializing/deserializing
// on io.Reader/io.Writer.
// It uses the encoding/gob package.
type GOBState struct {
	*MemoryState
}

// NewGOBState returns a new GOBState.
func NewGOBState() *GOBState {
	return &GOBState{NewMemoryState()}
}

func (g *GOBState) createBlob() *GOBStateBlob {
	g.muCreated.RLock()
	g.muRequest.RLock()
	defer g.muCreated.RUnlock()
	defer g.muRequest.RUnlock()

	return &GOBStateBlob{
		LastContentCreated: g.lastContentCreated,
		LastRequestTime:    g.lastRequestTime,
	}
}

func (g *GOBState) setFromBlob(b *GOBStateBlob) {
	g.muCreated.Lock()
	g.muRequest.Lock()
	defer g.muCreated.Unlock()
	defer g.muRequest.Unlock()

	g.lastContentCreated = b.LastContentCreated
	g.lastRequestTime = b.LastRequestTime
}

// Read will deserialize from a reader and populate its internal state.
func (g *GOBState) Read(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	var blob GOBStateBlob
	if err := decoder.Decode(&blob); err != nil {
		return err
	}
	g.setFromBlob(&blob)
	return nil
}

// Write will serialize its internal state and write to a writer.
func (g *GOBState) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	blob := g.createBlob()
	if err := encoder.Encode(&blob); err != nil {
		return err
	}
	return nil
}

// GOBStateBlob is used to serialize/deserialize MemoryState
// internal state.
type GOBStateBlob struct {
	LastContentCreated map[ContentType]time.Time
	LastRequestTime    map[ContentType]time.Time
}
