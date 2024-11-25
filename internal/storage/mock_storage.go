package storage

type MockStorage struct{}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) SaveOriginal(hash string, data []byte) error {
	return nil
}

func (m *MockStorage) CheckAndRetrieveResized(hash string, width, height int) (string, bool) {
	return "", false
}

func (m *MockStorage) SaveResized(hash string, width, height int, data []byte) error {
	return nil
}
