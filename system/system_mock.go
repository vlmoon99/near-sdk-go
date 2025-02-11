package system

type MockSystem struct {
	Registers map[uint64][]byte
	Storage   map[string][]byte
}

func NewMockSystem() *MockSystem {
	return &MockSystem{
		Registers: make(map[uint64][]byte),
		Storage:   make(map[string][]byte),
	}
}

// No need for ptr, just returning stored data
func (m *MockSystem) ReadRegisterSys(registerId, ptr uint64) {
}

func (m *MockSystem) RegisterLenSys(registerId uint64) uint64 {
	return uint64(len(m.Registers[registerId]))
}

// Simulating storing data (ignoring dataPtr)
func (m *MockSystem) WriteRegisterSys(registerId, dataLen, dataPtr uint64) {
	m.Registers[registerId] = make([]byte, dataLen)
}
