package model

import "errors"

type ProxyType int

const (
	Http ProxyType = iota + 1
	Https
	Socks4
	Socks5
)

func (t ProxyType) String() string {
	switch t {
	case Http:
		return "HTTP"
	case Https:
		return "HTTPS"
	case Socks4:
		return "SOCKS4"
	case Socks5:
		return "SOCKS5"
	default:
		return "UNKNOWN"
	}
}

func ParseType(s string) (ProxyType, error) {
	switch s {
	case "HTTP", "http":
		return Http, nil
	case "HTTPS", "https":
		return Https, nil
	case "SOCKS4", "socks4":
		return Socks4, nil
	case "SOCKS5", "socks5":
		return Socks5, nil
	default:
		return 0, errors.New("unknown type, possible values: HTTP, HTTPS, SOCKS4, SOCKS5")
	}
}

type ProxyProfile struct {
	ID      int       `db:"id"`
	Name    string    `db:"name"`
	Type    ProxyType `db:"type"`
	Address string    `db:"address"`
}

type Rule struct {
	ID           int           `db:"id"`
	Regex        string        `db:"regex"`
	ProxyProfile *ProxyProfile `db:"proxy_profile"`
}
