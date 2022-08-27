package regexp

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestDomain(t *testing.T) {
	t.Parallel()

	data := []struct{ name, input, want string }{
		{name: "google", input: "google.com", want: `^google\.com$`},
		{name: "aws", input: "aws.com", want: `^aws\.com$`},
		{name: "fb", input: "www.facebook.com", want: `^www\.facebook\.com$`},
		{name: "abc", input: "abc.def.ghi.jkl", want: `^abc\.def\.ghi\.jkl$`},
		{name: "localhost", input: "localhost:80", want: `^localhost:80$`},
		{name: "ip", input: "182.234.77.22", want: `^182\.234\.77\.22$`},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := Domain(d.input)
			assert.Equal(t, got, d.want)
		})
	}
}

func TestDomainAndSubdomains(t *testing.T) {
	t.Parallel()

	data := []struct{ name, input, want string }{
		{name: "google", input: "google.com", want: `(?:^|\.)google\.com$`},
		{name: "aws", input: "aws.com", want: `(?:^|\.)aws\.com$`},
		{name: "fb", input: "www.facebook.com", want: `(?:^|\.)www\.facebook\.com$`},
		{name: "abc", input: "abc.def.ghi.jkl", want: `(?:^|\.)abc\.def\.ghi\.jkl$`},
		{name: "localhost", input: "localhost:80", want: `(?:^|\.)localhost:80$`},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := DomainAndSubdomains(d.input)
			assert.Equal(t, got, d.want)
		})
	}
}
