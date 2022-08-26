package service

import (
	"context"
	"github.com/nnemirovsky/pacgen/internal/model"
)

type RuleRepository interface {
	GetAll(ctx context.Context) ([]model.Rule, error)
	GetAllWithProfiles(ctx context.Context) ([]model.Rule, error)
	GetByID(ctx context.Context, id int) (model.Rule, error)
	Create(ctx context.Context, rule *model.Rule) error
	Update(ctx context.Context, rule model.Rule) error
	Delete(ctx context.Context, id int) error
}

type ProxyProfileRepository interface {
	GetAll(ctx context.Context) ([]model.ProxyProfile, error)
	GetByID(ctx context.Context, id int) (model.ProxyProfile, error)
	Create(ctx context.Context, profile *model.ProxyProfile) error
	Update(ctx context.Context, profile model.ProxyProfile) error
	Delete(ctx context.Context, id int) error
}

type pacService interface {
	GeneratePACFile(ctx context.Context) error
}
