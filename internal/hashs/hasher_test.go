package hashs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		notWant string
		wantErr bool
	}{
		{
			name: "Hash a password",
			args: args{
				password: "123456",
			},
			notWant: "123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewHasher()
			got, err := hasher.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.notWant {
				t.Errorf("HashPassword() = %v, not want %v", got, tt.notWant)
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	t.Run("Should not return an error when comparing a password with its valid hash", func(t *testing.T) {
		// Arrange
		password := "123456"

		hasher := NewHasher()

		hashedPassword, err := hasher.HashPassword(password)
		assert.NoError(t, err)

		// Act
		err = hasher.CheckPassword(password, hashedPassword)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return a error when comparing a password with its invalid hash", func(t *testing.T) {
		// Arrange
		password := "123456"

		hasher := NewHasher()

		// Act
		err := hasher.CheckPassword(password, "123")

		// Assert
		assert.Error(t, err)
	})
}
