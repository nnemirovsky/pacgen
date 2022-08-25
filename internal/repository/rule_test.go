package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/pkg/logutil"
	"testing"
)

func testPrepareRuleRepository(t *testing.T) (*RuleRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewRuleRepository(dbx, logutil.DiscardLogger)

	return repo, mock
}

func TestRuleRepository_GetAll_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectQuery(`SELECT id, regex, proxy_profile_id AS "proxy_profile.id" FROM rules`).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "regex", "proxy_profile.id"}).
				AddRow(10, `^google\.com$`, 1).
				AddRow(20, `(?:^|\.)aws\.com$`, 2).
				AddRow(123456789, `^facebook\.com$`, 3),
		)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	want := []model.Rule{
		{10, `^google\.com$`, &model.ProxyProfile{ID: 1}},
		{20, `(?:^|\.)aws\.com$`, &model.ProxyProfile{ID: 2}},
		{123456789, `^facebook\.com$`, &model.ProxyProfile{ID: 3}},
	}

	assert.Equal(t, want, got)
}

func TestRuleRepository_GetAllWithProfiles_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectQuery(
			`SELECT r.id,
					r.regex,
					p.id AS "proxy_profile.id",
					p.name AS "proxy_profile.name",
					p.type AS "proxy_profile.type",
					p.address AS "proxy_profile.address"
			 FROM rules r
			 JOIN proxy_profiles p ON r.proxy_profile_id = p.id`,
		).
		WillReturnRows(
			sqlmock.
				NewRows(
					[]string{
						"id",
						"regex",
						"proxy_profile.id",
						"proxy_profile.name",
						"proxy_profile.type",
						"proxy_profile.address",
					},
				).
				AddRow(10, `^google\.com$`, 1, "shadowsocks", model.Socks5, "localhost:1080").
				AddRow(20, `(?:^|\.)aws\.com$`, 1, "shadowsocks", model.Socks5, "localhost:1080").
				AddRow(123456789, `^facebook\.com$`, 2, "tor", model.Socks5, "10.100.100.50:9050"),
		)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := repo.GetAllWithProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	want := []model.Rule{
		{
			ID:    10,
			Regex: `^google\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      1,
				Name:    "shadowsocks",
				Type:    model.Socks5,
				Address: "localhost:1080",
			},
		},
		{
			ID:    20,
			Regex: `(?:^|\.)aws\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      1,
				Name:    "shadowsocks",
				Type:    model.Socks5,
				Address: "localhost:1080",
			},
		},
		{
			ID:    123456789,
			Regex: `^facebook\.com$`,
			ProxyProfile: &model.ProxyProfile{
				ID:      2,
				Name:    "tor",
				Type:    model.Socks5,
				Address: "10.100.100.50:9050",
			},
		},
	}

	assert.Equal(t, want, got)
}

func TestRuleRepository_GetByID_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectQuery(
			`SELECT r.id,
					r.regex,
					p.id AS "proxy_profile.id",
					p.name AS "proxy_profile.name",
					p.type AS "proxy_profile.type",
					p.address AS "proxy_profile.address"
			 FROM rules r
			 JOIN proxy_profiles p ON r.proxy_profile_id = p.id
			 WHERE r.id = \?`,
		).
		WithArgs(10).
		WillReturnRows(
			sqlmock.
				NewRows(
					[]string{
						"id",
						"regex",
						"proxy_profile.id",
						"proxy_profile.name",
						"proxy_profile.type",
						"proxy_profile.address",
					},
				).
				AddRow(10, `^google\.com$`, 1, "shadowsocks", model.Socks5, "localhost:1080"),
		)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := repo.GetByID(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}

	want := model.Rule{
		ID:    10,
		Regex: `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{
			ID:      1,
			Name:    "shadowsocks",
			Type:    model.Socks5,
			Address: "localhost:1080",
		},
	}

	assert.Equal(t, want, got)
}

func TestRuleRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectQuery(
			`SELECT r.id,
					r.regex,
					p.id AS "proxy_profile.id",
					p.name AS "proxy_profile.name",
					p.type AS "proxy_profile.type",
					p.address AS "proxy_profile.address"
			 FROM rules r
			 JOIN proxy_profiles p ON r.proxy_profile_id = p.id
			 WHERE r.id = \?`,
		).
		WithArgs(10).
		WillReturnError(sql.ErrNoRows)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := repo.GetByID(ctx, 10)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected errs.EntityNotFoundError")
	}
}

func TestRuleRepository_Create_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	const insertedID = 15

	mock.
		ExpectExec(`INSERT INTO rules \(regex, proxy_profile_id\) VALUES \(\?, \?\) RETURNING id`).
		WithArgs(`^google\.com$`, 1).
		WillReturnResult(sqlmock.NewResult(insertedID, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rule := model.Rule{Regex: `^google\.com$`, ProxyProfile: &model.ProxyProfile{ID: 1}}
	err := repo.Create(ctx, &rule)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, insertedID, rule.ID)
}

func TestRuleRepository_Create_InvalidReference(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`INSERT INTO rules \(regex, proxy_profile_id\) VALUES \(\?, \?\) RETURNING id`).
		WithArgs(`^google\.com$`, 1).
		WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rule := model.Rule{Regex: `^google\.com$`, ProxyProfile: &model.ProxyProfile{ID: 1}}
	err := repo.Create(ctx, &rule)

	if err != errs.InvalidReferenceError {
		t.Fatal("expected errs.InvalidReferenceError")
	}
}

func TestRuleRepository_Update_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`UPDATE rules SET regex = \?, proxy_profile_id = \? WHERE id = \?`).
		WithArgs(`^google\.com$`, 1, 10).
		WillReturnResult(sqlmock.NewResult(10, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rule := model.Rule{ID: 10, Regex: `^google\.com$`, ProxyProfile: &model.ProxyProfile{ID: 1}}
	err := repo.Update(ctx, rule)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRuleRepository_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`UPDATE rules SET regex = \?, proxy_profile_id = \? WHERE id = \?`).
		WithArgs(`^google\.com$`, 1, 10).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rule := model.Rule{ID: 10, Regex: `^google\.com$`, ProxyProfile: &model.ProxyProfile{ID: 1}}
	err := repo.Update(ctx, rule)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected errs.EntityNotFoundError")
	}
}

func TestRuleRepository_Update_InvalidReference(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`UPDATE rules SET regex = \?, proxy_profile_id = \? WHERE id = \?`).
		WithArgs(`^google\.com$`, 1, 10).
		WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rule := model.Rule{ID: 10, Regex: `^google\.com$`, ProxyProfile: &model.ProxyProfile{ID: 1}}
	err := repo.Update(ctx, rule)

	if err != errs.InvalidReferenceError {
		t.Fatal("expected errs.InvalidReferenceError")
	}
}

func TestRuleRepository_Delete_OK(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`DELETE FROM rules WHERE id = \?`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(10, 1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Delete(ctx, 10)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRuleRepository_Delete_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := testPrepareRuleRepository(t)

	mock.
		ExpectExec(`DELETE FROM rules WHERE id = \?`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Delete(ctx, 10)

	if _, ok := err.(*errs.EntityNotFoundError); !ok {
		t.Fatal("expected errs.EntityNotFoundError")
	}
}
