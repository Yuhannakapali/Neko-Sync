package user_test

import (
	"context"
	"testing"

	"nekosync/internal/platform/postgres/pgtest"
	"nekosync/internal/user"
)

func TestRegisterDeviceRollsBackOnFailedInsert(t *testing.T) {
	ctx := context.Background()
	db := pgtest.Open(t)
	users := user.NewUserRepository(db)
	devices := user.NewDeviceRepository(db)
	svc := user.NewService(users, devices, user.NewFollowRepository(db), user.NewNotificationRepository(db))

	u, err := svc.CreateUser(ctx, "neko", "neko@example.com", "password1")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := users.GetByID(ctx, u.ID)
	if err != nil || stored.ID != u.ID {
		t.Fatalf("id changed form on read: created %q, read %v (%v)", u.ID, stored, err)
	}

	first, err := svc.RegisterDevice(ctx, u.ID, "laptop", user.PlatformDesktop)
	if err != nil {
		t.Fatal(err)
	}

	// An invalid platform passes Go but fails the platform_type enum, so the
	// insert fails after the deactivate has run inside the transaction.
	if _, err := svc.RegisterDevice(ctx, u.ID, "fridge", user.PlatformType("toaster")); err == nil {
		t.Fatal("expected insert to fail")
	}

	got, err := devices.GetByUserID(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != first.ID || !got[0].IsActive {
		t.Fatalf("failed registration must leave the existing device active: %+v", got)
	}

	second, err := svc.RegisterDevice(ctx, u.ID, "phone", user.PlatformMobile)
	if err != nil {
		t.Fatal(err)
	}
	got, _ = devices.GetByUserID(ctx, u.ID)
	for _, d := range got {
		if d.IsActive != (d.ID == second.ID) {
			t.Fatalf("only the newest device should be active: %+v", d)
		}
	}
}
