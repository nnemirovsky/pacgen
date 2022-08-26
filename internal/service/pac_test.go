package service

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/nnemirovsky/pacgen/internal/model"
	"testing"
)

func TestGeneratePAC_OK(t *testing.T) {
	t.Parallel()

	buff := bytes.NewBuffer([]byte{})

	rules := []model.Rule{
		{
			ID:    1,
			Regex: `^www\.google\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      1,
				Name:    "tor",
				Type:    model.Socks5,
				Address: "localhost:9050",
			},
		},
		{
			ID:    2,
			Regex: `(?:^|\.)facebook\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      2,
				Name:    "shadowsocks",
				Type:    model.Socks5,
				Address: "localhost:1080",
			},
		},
	}

	err := generatePAC(buff, rules)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	want := `function FindProxyForURL(url, host) {
	if (/^www\.google\.com$/.test(host)) return 'SOCKS5 localhost:9050';
	if (/(?:^|\.)facebook\.com$/.test(host)) return 'SOCKS5 localhost:1080';
	return 'DIRECT';
}`
	got := buff.String()

	assert.Equal(t, got, want)
}
