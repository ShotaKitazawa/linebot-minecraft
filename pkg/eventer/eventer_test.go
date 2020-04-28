package eventer

/*
import (
	"testing"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
)

func TestEventer(t *testing.T) {
	sharedMemValid := &sharedmemMockValid{}
	New("test", domain.LineConfig{}, sharedMemValid)
}

type sharedmemMockValid struct {
	data domain.Entity
}

func (m *sharedmemMockValid) SyncReadEntityFromSharedMem() (*domain.Entity, error) {
	return &domain.Entity{
		AllUsers:           []domain.User{},
		LoginUsers:         []domain.User{},
		LogoutUsers:        []domain.User{},
		WhitelistUsernames: []string{},
	}, nil
}

func (m *sharedmemMockValid) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.data = data
	return nil
}

*/
