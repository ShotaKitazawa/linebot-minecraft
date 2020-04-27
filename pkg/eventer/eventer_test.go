package eventer

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

var (
	loggerForTest            = logrus.New()
	minecraftHostnameForTest = `test`
)

func TestEventer(t *testing.T) {

	t.Run(`New()`, func(t *testing.T) {
		botSenderValid := &botSenderMockValid{}
		sharedMemValid := &sharedmemMockValid{}
		rconValid := &rconClientMockValid{}
		e, err := New(minecraftHostnameForTest, botSenderValid, sharedMemValid, rconValid, loggerForTest)
		assert.Nil(t, err)
		assert.Equal(t, e, &Eventer{
			BotPluginSender:   botSenderValid,
			MinecraftHostname: minecraftHostnameForTest,
			sharedMem:         sharedMemValid,
			rcon:              rconValid,
			Logger:            loggerForTest,
		})
	})
	t.Run(`job()`, func(t *testing.T) {
		t.Run(`normal`, func(t *testing.T) {
			t.Run(`numOfAllLoginUser:0,numOfCurrentLoginUser:0,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockValid{},
					&sharedmemMockValid{},
					&rconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:0,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockValid{},
					&sharedmemMockValid{data: &domain.Entity{
						AllUsers: []domain.User{{Name: `test`}},
					}},
					&rconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:1,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockValid{},
					&sharedmemMockValid{data: &domain.Entity{
						AllUsers: []domain.User{{Name: `test`}},
					}},
					&rconClientMockValid{LoginUsernames: []string{`test`}},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:0,numOfPreviousLoginUser:1`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockValid{},
					&sharedmemMockValid{data: &domain.Entity{
						AllUsers:   []domain.User{{Name: `test`}},
						LoginUsers: []domain.User{{Name: `test`}},
					}},
					&rconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:1,numOfPreviousLoginUser:1`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockValid{},
					&sharedmemMockValid{data: &domain.Entity{
						AllUsers:   []domain.User{{Name: `test`}},
						LoginUsers: []domain.User{{Name: `test`}},
					}},
					&rconClientMockValid{LoginUsernames: []string{`test`}},
				)
				assert.Nil(t, e.job())
			})
		})
		t.Run(`abnormal(BotPluginSender)`, func(t *testing.T) {
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:0,numOfPreviousLoginUser:1`, func(t *testing.T) {
				e := newForTest(
					&botSenderMockInvalid{},
					&sharedmemMockValid{data: &domain.Entity{
						AllUsers:   []domain.User{{Name: `test`}},
						LoginUsers: []domain.User{{Name: `test`}},
					}},
					&rconClientMockValid{},
				)
				assert.NotNil(t, e.job())
			})
		})
		t.Run(`abnormal(SharedMem)`, func(t *testing.T) {
			e := newForTest(
				&botSenderMockValid{},
				&sharedmemMockInvalid{},
				&rconClientMockValid{},
			)
			assert.NotNil(t, e.job())
		})
		t.Run(`abnormal(RconClient)`, func(t *testing.T) {
			e := newForTest(
				&botSenderMockValid{},
				&sharedmemMockValid{},
				&rconClientMockInvalid{},
			)
			assert.NotNil(t, e.job())
		})
	})
}

func newForTest(sender botplug.BotPluginSender, m sharedmem.SharedMem, rcon rcon.RconClient) *Eventer {
	return &Eventer{
		BotPluginSender:   sender,
		MinecraftHostname: minecraftHostnameForTest,
		sharedMem:         m,
		rcon:              rcon,
		Logger:            loggerForTest,
	}
}

type botSenderMockValid struct {
	msg string
}

func (sender *botSenderMockValid) SendTextMessage(msg string) error {
	sender.msg = msg
	return nil
}

type botSenderMockInvalid struct{}

func (sender *botSenderMockInvalid) SendTextMessage(msg string) error {
	return errors.New(``)
}

type sharedmemMockValid struct {
	data *domain.Entity
}

func (m *sharedmemMockValid) SyncReadEntityFromSharedMem() (*domain.Entity, error) {
	if m.data == nil {
		return nil, errors.New(``)
	}
	return m.data, nil
}

func (m *sharedmemMockValid) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.data = &data
	return nil
}

type sharedmemMockInvalid struct {
}

func (m *sharedmemMockInvalid) SyncReadEntityFromSharedMem() (*domain.Entity, error) {
	return nil, errors.New(``)
}

func (m *sharedmemMockInvalid) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	return errors.New(``)
}

type rconClientMockValid struct {
	LoginUsernames       []string
	WhitelistedUsernames []string
}

func (r *rconClientMockValid) List() ([]string, error) {
	return r.LoginUsernames, nil
}
func (r *rconClientMockValid) WhitelistAdd(username string) error {
	r.WhitelistedUsernames = append(r.WhitelistedUsernames, username)
	return nil
}
func (r *rconClientMockValid) WhitelistRemove(username string) error {
	var matched bool
	for idx, whitelistedUsername := range r.WhitelistedUsernames {
		if username == whitelistedUsername {
			matched = true
			if idx == len(r.WhitelistedUsernames)-1 {
				r.WhitelistedUsernames = r.WhitelistedUsernames[0:idx]
			} else {
				r.WhitelistedUsernames = append(r.WhitelistedUsernames[0:idx], r.WhitelistedUsernames[idx+1:]...)
			}
		}
	}
	if matched {
		return errors.New(``)
	}
	return nil
}
func (r *rconClientMockValid) WhitelistList() ([]string, error) {
	return r.WhitelistedUsernames, nil
}
func (r *rconClientMockValid) DataGetEntity(string) (*rcon.User, error) {
	return &rcon.User{}, nil
}
func (r *rconClientMockValid) Title(string) ([]string, error) {
	return r.LoginUsernames, nil
}

type rconClientMockInvalid struct {
}

func (r *rconClientMockInvalid) List() ([]string, error)                  { return nil, errors.New(``) }
func (r *rconClientMockInvalid) WhitelistAdd(string) error                { return errors.New(``) }
func (r *rconClientMockInvalid) WhitelistRemove(string) error             { return errors.New(``) }
func (r *rconClientMockInvalid) WhitelistList() ([]string, error)         { return nil, errors.New(``) }
func (r *rconClientMockInvalid) DataGetEntity(string) (*rcon.User, error) { return nil, errors.New(``) }
func (r *rconClientMockInvalid) Title(string) ([]string, error)           { return nil, errors.New(``) }
