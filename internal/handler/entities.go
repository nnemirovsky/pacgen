package handler

import (
	"errors"
	"fmt"
	"pacgen/internal/model"
	"pacgen/pkg/regexp"
)

type RuleR struct {
	ID             int    `json:"id"`
	Regexp         string `json:"regexp"`
	ProxyProfileID int    `json:"proxy_profile_id"`
}

func (r *RuleR) FromModel(rule model.Rule) {
	r.ID = rule.ID
	r.Regexp = rule.Regex
	r.ProxyProfileID = rule.ProxyProfile.ID
}

type RuleCU struct {
	Domain         string `json:"domain" validate:"required"`
	Mode           string `json:"mode" validate:"required,oneof=domain domain_and_subdomains"`
	ProxyProfileID int    `json:"proxy_profile_id" validate:"required"`
}

func (r *RuleCU) ToModel() (model.Rule, error) {
	var regex string
	switch r.Mode {
	case "domain":
		regex = regexp.Domain(r.Domain)
	case "domain_and_subdomains":
		regex = regexp.DomainAndSubdomains(r.Domain)
	default:
		return model.Rule{}, errors.New("invalid mode")
	}

	return model.Rule{Regex: regex, ProxyProfile: &model.ProxyProfile{ID: r.ProxyProfileID}}, nil
}

type ProxyProfileR struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func (p *ProxyProfileR) FromModel(profile model.ProxyProfile) {
	p.Address = profile.Address
	p.ID = profile.ID
	p.Name = profile.Name
	p.Type = profile.Type.String()
}

type ProxyProfileCU struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func (p *ProxyProfileCU) ToModel() (model.ProxyProfile, error) {
	t, err := model.ParseType(p.Type)
	if err != nil {
		return model.ProxyProfile{}, fmt.Errorf("invalid profile type: %w", err)
	}

	return model.ProxyProfile{
		Name:    p.Name,
		Type:    t,
		Address: p.Address,
	}, nil
}
