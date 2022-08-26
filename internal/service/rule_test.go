package service

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/internal/service/mock"
	"github.com/nnemirovsky/pacgen/pkg/logutil"
	"testing"
)

func testPrepareRuleService(t *testing.T) (*RuleService, *mock.RuleRepository, *mock.PacService) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewRuleRepository(ctrl)
	pacSrvcMock := mock.NewPacService(ctrl)

	pacSrvcMock.EXPECT().GeneratePACFile(gomock.Any()).Return(nil).AnyTimes()

	return NewRuleService(repoMock, pacSrvcMock, logutil.DiscardLogger), repoMock, pacSrvcMock
}

func TestRuleService_GetAll_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	want := []model.Rule{
		{
			ID:           1,
			Regex:        `^www\.google\.com$`,
			ProxyProfile: &model.ProxyProfile{ID: 1},
		},
		{
			ID:           2,
			Regex:        `(?:^|\.)facebook\.com$`,
			ProxyProfile: &model.ProxyProfile{ID: 2},
		},
	}

	repoMock.EXPECT().GetAll(gomock.Any()).Return(want, nil)

	got, err := ruleSrvc.GetAll(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, got, want)
}

func TestRuleService_GetAllWithProfiles_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	want := []model.Rule{
		{
			ID:    1,
			Regex: `^www\.google\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      1,
				Name:    "shadowsocks",
				Type:    model.Socks5,
				Address: "localhost:1080",
			},
		},
		{
			ID:    2,
			Regex: `(?:^|\.)facebook\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      2,
				Name:    "some http proxy",
				Type:    model.Http,
				Address: "localhost:8080",
			},
		},
	}

	repoMock.EXPECT().GetAllWithProfiles(gomock.Any()).Return(want, nil)

	got, err := ruleSrvc.GetAllWithProfiles(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, got, want)
}

func TestRuleService_GetByID_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	want := model.Rule{
		ID:           1,
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	repoMock.EXPECT().GetByID(gomock.Any(), want.ID).Return(want, nil)

	got, err := ruleSrvc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, got, want)
}

func TestRuleService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	repoMock.EXPECT().GetByID(gomock.Any(), 1).Return(model.Rule{}, &errs.EntityNotFoundError{})

	_, err := ruleSrvc.GetByID(context.Background(), 1)
	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}

func TestRuleService_Create_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	const insertedID = 15

	rule := model.Rule{
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	repoMock.EXPECT().Create(gomock.Any(), &rule).DoAndReturn(
		func(ctx context.Context, rule *model.Rule) error {
			rule.ID = insertedID
			return nil
		},
	)

	err := ruleSrvc.Create(context.Background(), &rule)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, rule.ID, insertedID)
}

func TestRuleService_Create_InvalidReference(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	rule := model.Rule{
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 15},
	}

	repoMock.EXPECT().Create(gomock.Any(), &rule).Return(errs.InvalidReferenceError)

	err := ruleSrvc.Create(context.Background(), &rule)
	if err != errs.InvalidReferenceError {
		t.Errorf("expected error is errs.InvalidReferenceError, but got %#v", err)
	}
}

func TestRuleService_Update_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	rule := model.Rule{
		ID:           1,
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	repoMock.EXPECT().Update(gomock.Any(), rule).Return(nil)

	err := ruleSrvc.Update(context.Background(), rule)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestRuleService_Update_InvalidReference(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	rule := model.Rule{
		ID:           1,
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 15},
	}

	repoMock.EXPECT().Update(gomock.Any(), rule).Return(errs.InvalidReferenceError)

	err := ruleSrvc.Update(context.Background(), rule)
	if err != errs.InvalidReferenceError {
		t.Errorf("expected error is errs.InvalidReferenceError, but got %#v", err)
	}
}

func TestRuleService_Update_NotFound(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	rule := model.Rule{
		ID:           1,
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	repoMock.EXPECT().Update(gomock.Any(), rule).Return(&errs.EntityNotFoundError{})

	err := ruleSrvc.Update(context.Background(), rule)
	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}

func TestRuleService_Delete_OK(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	repoMock.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	err := ruleSrvc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestRuleService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	ruleSrvc, repoMock, _ := testPrepareRuleService(t)

	repoMock.EXPECT().Delete(gomock.Any(), 1).Return(&errs.EntityNotFoundError{})

	err := ruleSrvc.Delete(context.Background(), 1)
	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}
