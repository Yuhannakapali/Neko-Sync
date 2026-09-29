package user

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"nekosync/internal/platform/entity"
)

// fakeDeviceRepo models a table where ReplaceActive is all-or-nothing, as it is
// in the Postgres implementation.
type fakeDeviceRepo struct {
	devices   []*Device
	createErr error
}

func (f *fakeDeviceRepo) Create(_ context.Context, d *Device) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.devices = append(f.devices, d)
	return nil
}

func (f *fakeDeviceRepo) GetByUserID(context.Context, entity.UUID) ([]*Device, error) {
	return f.devices, nil
}
func (f *fakeDeviceRepo) Update(context.Context, *Device) error     { return nil }
func (f *fakeDeviceRepo) Delete(context.Context, entity.UUID) error { return nil }

func (f *fakeDeviceRepo) DeactivateAllForUser(_ context.Context, userID entity.UUID) error {
	for _, d := range f.devices {
		if d.UserID == userID {
			d.IsActive = false
		}
	}
	return nil
}

func (f *fakeDeviceRepo) ReplaceActive(ctx context.Context, d *Device) error {
	if f.createErr != nil {
		return f.createErr
	}
	_ = f.DeactivateAllForUser(ctx, d.UserID)
	f.devices = append(f.devices, d)
	return nil
}

func TestRegisterDeviceKeepsExistingDeviceActiveWhenInsertFails(t *testing.T) {
	existing := &Device{UserID: "u1", DeviceName: "old", IsActive: true}
	repo := &fakeDeviceRepo{devices: []*Device{existing}, createErr: errors.New("insert failed")}
	svc := NewService(nil, repo, nil, nil)

	if _, err := svc.RegisterDevice(context.Background(), "u1", "new", "web"); err == nil {
		t.Fatal("expected error")
	}
	if !existing.IsActive {
		t.Fatal("failed registration deactivated the user's existing device")
	}
}

func TestRegisterDeviceReplacesActiveDevice(t *testing.T) {
	existing := &Device{UserID: "u1", DeviceName: "old", IsActive: true}
	repo := &fakeDeviceRepo{devices: []*Device{existing}}
	svc := NewService(nil, repo, nil, nil)

	d, err := svc.RegisterDevice(context.Background(), "u1", "new", "web")
	if err != nil {
		t.Fatal(err)
	}
	if existing.IsActive || !d.IsActive || len(repo.devices) != 2 {
		t.Fatalf("unexpected state: existing=%v new=%v count=%d", existing.IsActive, d.IsActive, len(repo.devices))
	}
}

type fakeUserRepo struct {
	Repository
	created *User
}

func (f *fakeUserRepo) GetByEmail(context.Context, string) (*User, error)    { return nil, ErrNotFound }
func (f *fakeUserRepo) GetByUsername(context.Context, string) (*User, error) { return nil, ErrNotFound }
func (f *fakeUserRepo) Create(_ context.Context, u *User) error              { f.created = u; return nil }

// Postgres stores ids in a uuid column and returns them in canonical form, so an
// id in any other form changes between registration and every later read.
func TestCreateUserUsesCanonicalUUID(t *testing.T) {
	repo := &fakeUserRepo{}
	u, err := NewService(repo, nil, nil, nil).CreateUser(context.Background(), "neko", "n@example.com", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if !canonicalUUID.MatchString(string(u.ID)) {
		t.Fatalf("id %q is not a canonical UUID", u.ID)
	}
}

var canonicalUUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
