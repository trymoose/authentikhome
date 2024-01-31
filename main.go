package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/trymoose/authentikhome/pkg/homeassistant"
	"github.com/trymoose/authentikhome/pkg/ldap"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func main() {
	defer RecoverAndExit()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	args := ReadArgs()
	cfg := ReadConfig[struct {
		BaseDN   string  `yaml:"base"`
		BindUser string  `yaml:"username"`
		BindPass string  `yaml:"password"`
		Group    *string `yaml:"group"`

		Host   string `yaml:"host"`
		Port   uint16 `yaml:"port"`
		Secure bool   `yaml:"secure"`
	}](args.ConfigFile)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn := Must((&ldap.Dialer{Secure: cfg.Secure}).Dial(
		ctx,
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	))
	defer Defer(conn.Close)

	user := Must(conn.Login(
		cfg.BaseDN,
		cfg.BindUser,
		cfg.BindPass,
		args.Username,
		args.Password,
		cfg.Group,
	))
	if !user.Active {
		return
	}

	meta := homeassistant.Meta{Name: user.Name, Group: homeassistant.GroupUsers}
	if user.Admin {
		meta.Group = homeassistant.GroupAdmin
	}
	fmt.Print(string(Must(meta.MarshalText())))
	Login()
}

func ReadArgs() (args struct {
	Username, Password string
	ConfigFile         string
}) {
	args.Username = Must(homeassistant.EnvKeyUsername.Value())
	args.Password = Must(homeassistant.EnvKeyPassword.Value())
	args.ConfigFile = "./config.yml"

	flag.StringVar(&args.ConfigFile, "config", args.ConfigFile, "ldap config file to read")
	flag.BoolVar(&printStack, "debug", printStack, "print to terminal")
	flag.Parse()
	return
}

func ReadConfig[T any](filename string) (cfg T) {
	f := Must(os.Open(filename))
	defer Defer(f.Close)
	Check(yaml.NewDecoder(f).Decode(&cfg))
	return
}
