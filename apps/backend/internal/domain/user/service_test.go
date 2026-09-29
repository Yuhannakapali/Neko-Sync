package user

import (
	"context"
	"errors"
	"testing"

	"nekosync/internal/domain/shared"
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

func (f *fakeDeviceRepo) GetByUserID(context.Context, shared.UUID) ([]*Device, error) {
	return f.devices, nil
}
func (f *fakeDeviceRepo) Update(context.Context, *Device) error     { return nil }
func (f *fakeDeviceRepo) Delete(context.Context, shared.UUID) error { return nil }

func (f *fakeDeviceRepo) DeactivateAllForUser(_ context.Context, userID shared.UUID) error {
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
