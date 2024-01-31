package homeassistant

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

//go:generate go run github.com/dmarkham/enumer@latest -type Group -linecomment

type Group int

const (
	GroupUsers Group = iota // system-users
	GroupAdmin              // system-admin
)

type ErrInvalidGroup Group

func (e ErrInvalidGroup) Error() string { return fmt.Sprintf("'%d' is not a valid group", int(e)) }

type Meta struct {
	Name      string
	Group     Group
	LocalOnly bool
}

func (m *Meta) MarshalText() ([]byte, error) {
	if !m.Group.IsAGroup() {
		return nil, ErrInvalidGroup(m.Group)
	}

	var buf bytes.Buffer
	for _, meta := range [][2]string{
		{"name", m.Name},
		{"group", m.Group.String()},
		{"local_only", strconv.FormatBool(m.LocalOnly)},
	} {
		_, _ = fmt.Fprintf(&buf, "%s = %s\n", meta[0], meta[1])
	}
	return buf.Bytes(), nil
}

type EnvKey string

const (
	EnvKeyUsername EnvKey = "username"
	EnvKeyPassword EnvKey = "password"
)

type ErrEnvKeyNotExist string

func (e ErrEnvKeyNotExist) Error() string {
	return fmt.Sprintf("env key %q not set", string(e))
}

func (e EnvKey) Value() (string, error) {
	v, ok := os.LookupEnv(string(e))
	if !ok {
		panic(ErrEnvKeyNotExist(e))
	}
	return v, nil
}
