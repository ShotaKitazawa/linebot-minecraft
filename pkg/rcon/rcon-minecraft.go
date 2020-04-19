package rcon

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/seeruk/minecraft-rcon/rcon"
)

type Client struct {
	*rcon.Client
}

type User struct {
	Health  string
	XpLevel string
	Position
}

type Position struct {
	X float32
	Y float32
	Z float32
}

func New(host string, port int, password string) (*Client, error) {
	client, err := rcon.NewClient(host, port, password)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

type Command struct {
	command    string
	expression string
}

func (c Client) command(command Command) ([]string, error) {
	response, err := c.Client.SendCommand(command.command)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(command.expression)
	extracted := re.FindStringSubmatch(response)
	if len(extracted) == 0 {
		return nil, fmt.Errorf(`"%s" is not match to "%s"`, command.expression, response)
	}
	return extracted[1:], nil
}

func (c Client) List() ([]string, error) {
	result, err := c.command(Command{
		command:    `list`,
		expression: `There are [0-9].* of a max [0-9].* players online: ?(.*)$`,
	})
	if err != nil {
		return nil, err
	} else if result[0] == "" {
		return nil, nil
	}
	return result, nil
}

func (c Client) WhitelistList() ([]string, error) {
	result, err := c.command(Command{
		command:    `whitelist list`,
		expression: `There are ([0-9].*) whitelisted players: ?(.*)`,
	})
	if err != nil {
		return nil, err
	} else if result[0] == "" {
		return nil, nil
	}
	return result, nil
}

func (c Client) DataGetEntity(username string) (*User, error) {
	array, err := c.command(Command{
		command:    fmt.Sprintf(`data get entity %s`, username),
		expression: fmt.Sprintf(`%s has the following entity data: {.*Health: (.*?),.*XpLevel: (.*?),.*Pos: \[(.*?)d, (.*?)d, (.*?)d\].*$`, username),
	})
	if err != nil {
		return nil, err
	}
	posX, err := strconv.ParseFloat(array[2], 32)
	if err != nil {
		return nil, err
	}
	posY, err := strconv.ParseFloat(array[3], 32)
	if err != nil {
		return nil, err
	}
	posZ, err := strconv.ParseFloat(array[4], 32)
	if err != nil {
		return nil, err
	}
	user := &User{
		Health:  array[0],
		XpLevel: array[1],
		Position: Position{
			X: float32(posX),
			Y: float32(posY),
			Z: float32(posZ),
		},
	}
	return user, nil
}
