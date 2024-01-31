package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"net"
	"strconv"
)

const (
	AttrName   = "name"
	AttrActive = "ak-active"
	AttrAdmin  = "ak-superuser"
	AttrDN     = "dn"
)

type (
	Dialer struct {
		Secure bool
	}
	Conn ldap.Conn
)

func (d Dialer) Dial(ctx context.Context, network, addr string) (*Conn, error) {
	netConn, err := new(net.Dialer).DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	if d.Secure {
		conn := tls.Client(netConn, &tls.Config{InsecureSkipVerify: true})
		if err := conn.HandshakeContext(ctx); err != nil {
			return nil, fmt.Errorf("ssl failed: %w", err)
		}
		netConn = conn
	}

	conn := ldap.NewConn(netConn, d.Secure)
	conn.Start()

	return (*Conn)(conn), nil
}

func (conn *Conn) Close() error { return (*ldap.Conn)(conn).Close() }

type User struct {
	DN     string
	Name   string
	Active bool
	Admin  bool
}

var (
	ErrBindToSearchUserFailed = errors.New("failed to bind to search user")
	ErrSearchFailed           = errors.New("search for user failed")
	ErrUserNotExist           = errors.New("user does not exist")
	ErrMultipleUserExist      = errors.New("search for user returned multiple users")
	ErrResultParseFailed      = errors.New("failed to parse result")
	ErrLoginFailed            = errors.New("failed to login")
)

func (conn *Conn) Login(baseDN, searchUsername, searchPassword, username, password string, group *string) (*User, error) {
	searchUsername = fmt.Sprintf("cn=%s,%s", searchUsername, baseDN)
	cnn := (*ldap.Conn)(conn)

	if err := cnn.Bind(searchUsername, searchPassword); err != nil {
		return nil, errors.Join(ErrBindToSearchUserFailed, err)
	}

	res, err := cnn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		getQuery(baseDN, username, group),
		[]string{AttrDN, AttrName, AttrActive, AttrAdmin},
		nil,
	))
	if err != nil {
		return nil, errors.Join(ErrSearchFailed, err)
	}

	user, err := getUserFromResult(res)
	if err != nil {
		return nil, err
	}

	if err := cnn.Bind(user.DN, password); err != nil {
		return nil, errors.Join(ErrLoginFailed, err)
	}
	return user, nil
}

func getUserFromResult(res *ldap.SearchResult) (_ *User, err error) {
	var user User
	switch len(res.Entries) {
	case 0:
		return nil, ErrUserNotExist
	case 1:
		e := res.Entries[0]
		user.DN = e.DN
		user.Name = e.GetAttributeValue(AttrName)
		if user.Active, err = strconv.ParseBool(e.GetAttributeValue(AttrActive)); err != nil {
			return nil, errors.Join(ErrResultParseFailed, err)
		}
		if user.Admin, err = strconv.ParseBool(e.GetAttributeValue(AttrAdmin)); err != nil {
			return nil, errors.Join(ErrResultParseFailed, err)
		}
	default:
		return nil, ErrMultipleUserExist
	}
	return &user, nil
}

func getQuery(baseDN, username string, group *string) (query string) {
	query = fmt.Sprintf("(&(objectClass=user)(cn=%s,%s))", username, baseDN)
	if group != nil {
		query = fmt.Sprintf("(&(objectClass=user)(memberOf=cn=%s,ou=groups,%s)(cn=%s))", *group, baseDN, username)
	}
	return
}

func Group(s string) *string { return &s }
