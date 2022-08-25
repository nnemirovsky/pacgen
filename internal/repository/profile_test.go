package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/pkg/logutil"
	"testing"
)

func testPrepareProxyProfileRepository(t *testing.T) (*ProxyProfileRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewProxyProfileRepository(dbx, logutil.DiscardLogger)

	return repo, mock
}

func TestProxyProfileRepository_GetAll_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectQuery(`SELECT id, name, type, address FROM proxy_profiles`).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "name", "type", "address"}).
				AddRow(10, "shadowsocks", model.Https, "127.0.0.1:1080").
				AddRow(20, "simple socks", model.Socks5, "127.0.0.1:9999").
				AddRow(123456789, "tor", model.Http, "localhost:9050").
				AddRow(1111, "some name", model.Socks4, "::1:1080"),
		)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	want := []model.ProxyProfile{
		{10, "shadowsocks", model.Https, "127.0.0.1:1080"},
		{20, "simple socks", model.Socks5, "127.0.0.1:9999"},
		{123456789, "tor", model.Http, "localhost:9050"},
		{1111, "some name", model.Socks4, "::1:1080"},
	}

	assert.Equal(t, got, want)
}

func TestProxyProfileRepository_GetByID_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectQuery(`SELECT id, name, type, address FROM proxy_profiles WHERE id = \?`).
		WithArgs(10).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "name", "type", "address"}).
				AddRow(
					driver.Value(10),
					driver.Value("shadowsocks"),
					driver.Value(model.Https),
					driver.Value("127.0.0.1:1081"),
				),
		)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := repo.GetByID(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}

	want := model.ProxyProfile{ID: 10, Name: "shadowsocks", Type: model.Https, Address: "127.0.0.1:1081"}

	assert.Equal(t, got, want)
}

func TestProxyProfileRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectQuery(`SELECT id, name, type, address FROM proxy_profiles WHERE id = \?`).
		WithArgs(0).
		WillReturnError(sql.ErrNoRows)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := repo.GetByID(ctx, 0)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected error errs.EntityNotFoundError")
	}
}

func TestProxyProfileRepository_Create_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	const insertedID = 15

	mock.
		ExpectExec(`INSERT INTO proxy_profiles \(name, type, address\) VALUES \(\?, \?, \?\)`).
		WithArgs("some name", model.Http, "::1:1080").
		WillReturnResult(sqlmock.NewResult(insertedID, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	profile := model.ProxyProfile{Name: "some name", Type: model.Http, Address: "::1:1080"}
	err := repo.Create(ctx, &profile)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, profile.ID, insertedID)
}

func TestProxyProfileRepository_Create_AlreadyExists(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`INSERT INTO proxy_profiles \(name, type, address\) VALUES \(\?, \?, \?\)`).
		WithArgs("some socks", model.Socks5, "1.1.1.1:1080").
		WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	profile := model.ProxyProfile{Name: "some socks", Type: model.Socks5, Address: "1.1.1.1:1080"}
	err := repo.Create(ctx, &profile)

	if _, ok := err.(*errs.EntityAlreadyExistsError); !ok {
		t.Fatal("expected error errs.EntityAlreadyExistsError")
	}
}

func TestProxyProfileRepository_Update_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`UPDATE proxy_profiles SET name = \?, type = \?, address = \? WHERE id = \?`).
		WithArgs("some name", model.Https, "127.0.0.1:1080", 10).
		WillReturnResult(sqlmock.NewResult(10, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	profile := model.ProxyProfile{ID: 10, Name: "some name", Type: model.Https, Address: "127.0.0.1:1080"}
	err := repo.Update(ctx, profile)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProxyProfileRepository_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`UPDATE proxy_profiles SET name = \?, type = \?, address = \? WHERE id = \?`).
		WithArgs("some name", model.Https, "127.0.0.1:1080", 10).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	profile := model.ProxyProfile{ID: 10, Name: "some name", Type: model.Https, Address: "127.0.0.1:1080"}
	err := repo.Update(ctx, profile)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected error errs.EntityNotFoundError")
	}
}

func TestProxyProfileRepository_Update_AlreadyExists(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`UPDATE proxy_profiles SET name = \?, type = \?, address = \? WHERE id = \?`).
		WithArgs("some socks", model.Socks5, "localhost:1080", 10).
		WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	profile := model.ProxyProfile{ID: 10, Name: "some socks", Type: model.Socks5, Address: "localhost:1080"}
	err := repo.Update(ctx, profile)

	if _, ok := err.(*errs.EntityAlreadyExistsError); !ok {
		t.Fatal("expected error errs.EntityAlreadyExistsError")
	}
}

func TestProxyProfileRepository_Delete_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`DELETE FROM proxy_profiles WHERE id = \?`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(10, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Delete(ctx, 10)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProxyProfileRepository_Delete_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`DELETE FROM proxy_profiles WHERE id = \?`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Delete(ctx, 10)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected error errs.EntityNotFoundError")
	}
}

func TestProxyProfileRepository_Delete_StillReferenced(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareProxyProfileRepository(t)

	mock.
		ExpectExec(`DELETE FROM proxy_profiles WHERE id = \?`).
		WithArgs(10).
		WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Delete(ctx, 10)

	if _, ok := err.(*errs.EntityStillReferencedError); !ok {
		t.Fatal("expected error errs.EntityStillReferencedError")
	}
}
