package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

type pwUsers struct {
	user    *dbsqlc.User
	updated string
}

func (f *pwUsers) GetByID(context.Context, uuid.UUID) (*dbsqlc.User, error) { return f.user, nil }
func (f *pwUsers) GetByEmail(context.Context, string) (*dbsqlc.User, error) {
	return nil, sql.ErrNoRows
}
func (f *pwUsers) Create(context.Context, dbsqlc.CreateUserParams) (*dbsqlc.User, error) {
	return nil, nil
}
func (f *pwUsers) UpdatePasswordHash(_ context.Context, _ uuid.UUID, h string) error {
	f.updated = h
	return nil
}
func (f *pwUsers) UpdateName(context.Context, uuid.UUID, string) (*dbsqlc.User, error) {
	return nil, nil
}
func (f *pwUsers) UpdateEmail(context.Context, uuid.UUID, string) (*dbsqlc.User, error) {
	return nil, nil
}

func TestChangePassword(t *testing.T) {
	hash, err := hashPassword("OldSecure@1")
	if err != nil {
		t.Fatal(err)
	}
	users := &pwUsers{user: &dbsqlc.User{PasswordHash: sql.NullString{String: hash, Valid: true}}}
	svc := NewAccountService(users, nil, nil, nil)
	ctx := context.Background()

	if err := svc.ChangePassword(ctx, uuid.New(), "OldSecure@1", "NewSecure@2", "different"); err != ErrPasswordMismatch {
		t.Errorf("mismatch: got %v", err)
	}
	if err := svc.ChangePassword(ctx, uuid.New(), "wrong", "NewSecure@2", "NewSecure@2"); err != ErrWrongPassword {
		t.Errorf("wrong current: got %v", err)
	}
	if err := svc.ChangePassword(ctx, uuid.New(), "OldSecure@1", "NewSecure@2", "NewSecure@2"); err != nil {
		t.Fatalf("success: %v", err)
	}
	if !verifyPassword("NewSecure@2", users.updated) {
		t.Error("stored hash does not match the new password")
	}
}

func TestSpecStatus(t *testing.T) {
	cases := map[dbsqlc.OrderStatus]string{
		dbsqlc.OrderStatusPendingPayment: "open", dbsqlc.OrderStatusPaid: "open",
		dbsqlc.OrderStatusShipped: "open", dbsqlc.OrderStatusDelivered: "fulfilled",
		dbsqlc.OrderStatusCancelled: "cancelled", dbsqlc.OrderStatusRefunded: "cancelled",
	}
	for in, want := range cases {
		if got := SpecStatus(in); got != want {
			t.Errorf("SpecStatus(%s)=%s want %s", in, got, want)
		}
	}
}
