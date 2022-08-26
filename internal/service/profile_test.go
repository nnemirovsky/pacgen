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

func testPrepareProxyProfileService(t *testing.T) (
	*ProxyProfileService,
	*mock.ProxyProfileRepository,
	*mock.PacService,
) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewProxyProfileRepository(ctrl)
	pacSrvcMock := mock.NewPacService(ctrl)

	pacSrvcMock.EXPECT().GeneratePACFile(gomock.Any()).Return(nil).AnyTimes()

	return NewProxyProfileService(repoMock, pacSrvcMock, logutil.DiscardLogger), repoMock, pacSrvcMock
}

func TestProxyProfileService_GetAll_OK(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	want := []model.ProxyProfile{
		{
			ID:      1,
			Name:    "shadowsocks",
			Type:    model.Socks5,
			Address: "192.168.1.1:1080",
		},
		{
			ID:      2,
			Name:    "some http proxy",
			Type:    model.Http,
			Address: "::1:8080",
		},
	}

	repoMock.EXPECT().GetAll(gomock.Any()).Return(want, nil)

	got, err := proxyProfileSrvc.GetAll(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, got, want)
}

func TestProxyProfileService_GetByID_OK(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	want := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "1.1.1.1:1080",
	}

	repoMock.EXPECT().GetByID(gomock.Any(), want.ID).Return(want, nil)

	got, err := proxyProfileSrvc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, got, want)
}

func TestProxyProfileService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	repoMock.EXPECT().GetByID(gomock.Any(), 100).Return(model.ProxyProfile{}, &errs.EntityNotFoundError{})

	_, err := proxyProfileSrvc.GetByID(context.Background(), 100)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}

func TestProxyProfileService_Create_OK(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	profile := model.ProxyProfile{
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "10.1.1.1:9999",
	}

	const insertedID = 15

	repoMock.EXPECT().Create(gomock.Any(), &profile).DoAndReturn(
		func(ctx context.Context, profile *model.ProxyProfile) error {
			profile.ID = insertedID
			return nil
		},
	)

	err := proxyProfileSrvc.Create(context.Background(), &profile)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	assert.Equal(t, profile.ID, insertedID)
}

func TestProxyProfileService_Create_AlreadyExists(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	profile := model.ProxyProfile{
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	repoMock.EXPECT().Create(gomock.Any(), &profile).Return(&errs.EntityAlreadyExistsError{})

	err := proxyProfileSrvc.Create(context.Background(), &profile)

	if _, ok := err.(*errs.EntityAlreadyExistsError); !ok {
		t.Errorf("expected error is errs.EntityAlreadyExistsError, but got %#v", err)
	}
}

func TestProxyProfileService_Update_OK(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "127.0.0.1:1080",
	}

	repoMock.EXPECT().Update(gomock.Any(), profile).Return(nil)

	err := proxyProfileSrvc.Update(context.Background(), profile)

	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestProxyProfileService_Update_NotFound(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "127.0.0.1:1080",
	}

	repoMock.EXPECT().Update(gomock.Any(), profile).Return(&errs.EntityNotFoundError{})

	err := proxyProfileSrvc.Update(context.Background(), profile)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}

func TestProxyProfileService_Update_AlreadyExists(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	repoMock.EXPECT().Update(gomock.Any(), profile).Return(&errs.EntityAlreadyExistsError{})

	err := proxyProfileSrvc.Update(context.Background(), profile)

	if _, ok := err.(*errs.EntityAlreadyExistsError); !ok {
		t.Errorf("expected error is errs.EntityAlreadyExistsError, but got %#v", err)
	}
}

func TestProxyProfileService_Delete_OK(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	repoMock.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	err := proxyProfileSrvc.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestProxyProfileService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	repoMock.EXPECT().Delete(gomock.Any(), 1).Return(&errs.EntityNotFoundError{})

	err := proxyProfileSrvc.Delete(context.Background(), 1)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Errorf("expected error is errs.EntityNotFoundError, but got %#v", err)
	}
}

func TestProxyProfileService_Delete_StillReferenced(t *testing.T) {
	t.Parallel()

	proxyProfileSrvc, repoMock, _ := testPrepareProxyProfileService(t)

	repoMock.EXPECT().Delete(gomock.Any(), 1).Return(&errs.EntityStillReferencedError{})

	err := proxyProfileSrvc.Delete(context.Background(), 1)

	if _, ok := err.(*errs.EntityStillReferencedError); !ok {
		t.Errorf("expected error is errs.EntityStillReferencedError, but got %#v", err)
	}
}
