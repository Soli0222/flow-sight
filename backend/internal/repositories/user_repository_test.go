package repositories

import "testing"

func TestUserRepository_Skipped(t *testing.T) {
	t.Skip("User repository removed in single-user mode")
}
