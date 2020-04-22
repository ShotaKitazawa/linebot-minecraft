package rcon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	command            string
	expression         string
	expressionNotFound string
}

func (c Client) command(command Command) ([]string, error) {
	response, err := c.Client.SendCommand(command.command)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(command.expression)
	extracted := re.FindStringSubmatch(response)
	if len(extracted) == 0 {
		re = regexp.MustCompile(command.expressionNotFound)
		extracted = re.FindStringSubmatch(response)
		if len(extracted) == 0 {
			return nil, fmt.Errorf(`"%s" is not match to "%s"`, command.expression, response)
		}
		return nil, nil
	}

	return extracted[1:], nil
}

func (c Client) List() ([]string, error) {
	result, err := c.command(Command{
		command:            `list`,
		expression:         `There are [0-9].* of a max [0-9].* players online: (.*)$`,
		expressionNotFound: `There are 0 of a max [0-9].* players online:`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}
	return strings.Split(result[0], ", "), nil
}

func (c Client) WhitelistAdd(username string) error {
	_, err := c.command(Command{
		command:            fmt.Sprintf(`whitelist add %s`, username),
		expression:         fmt.Sprintf(`Added %s to the whitelist`, username),
		expressionNotFound: `!!!not much!!!`,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c Client) WhitelistRemove(username string) error {
	_, err := c.command(Command{
		command:    fmt.Sprintf(`whitelist remove %s`, username),
		expression: fmt.Sprintf(`Removed %s from the whitelist`, username),

		expressionNotFound: `!!!not much!!!`,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c Client) WhitelistList() ([]string, error) {
	result, err := c.command(Command{
		command:            `whitelist list`,
		expression:         `There are [0-9].* whitelisted players: (.*)`,
		expressionNotFound: `There are no whitelisted players`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}
	return strings.Split(result[0], ", "), nil
}

func (c Client) DataGetEntity(username string) (*User, error) {
	array, err := c.command(Command{
		command:            fmt.Sprintf(`data get entity %s`, username),
		expression:         fmt.Sprintf(`%s has the following entity data: {.*Health: (.*?),.*XpLevel: (.*?),.*Pos: \[(.*?)d, (.*?)d, (.*?)d\].*$`, username),
		expressionNotFound: `!!!not much!!!`,
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

func (c Client) Title(msg string) ([]string, error) {
	result, err := c.command(Command{
		command:            fmt.Sprintf(`title @a title {"text": "%s"}`, msg),
		expression:         `Showing new title for (.*)$`,
		expressionNotFound: `No player was found`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}
	return strings.Split(result[0], ", "), nil
}
