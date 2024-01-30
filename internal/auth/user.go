package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(v repo.User) User {
	return User{
		ID:           v.ID,
		Email:        v.Email,
		Username:     v.Username,
		PasswordHash: v.Password,
	}
}

type User struct {
	ID           int64
	Email        string `validate:"required,lte=128,email,excludes= "`
	Username     string `validate:"gte=3,lte=64,excludes=@,excludes= "`
	Password     string `validate:"gte=8"`
	PasswordHash string
}

func (u *User) normalize() {
	u.Email = strings.ToLower(u.Email)
	u.Username = strings.ToLower(u.Username)
}

func (u *User) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

func CreateUser(ctx context.Context, db repo.DB, arg User) (int64, error) {
	arg.normalize()

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	if err := arg.hashPassword(); err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return db.AuthCreateUser(ctx, repo.AuthCreateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  arg.PasswordHash,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func UpdateUser(ctx context.Context, db repo.DB, arg User) error {
	arg.normalize()

	if arg.Password == "" {
		if err := validate.Validate.StructExcept(arg, "Password"); err != nil {
			return err
		}
	} else {
		if err := validate.Validate.Struct(arg); err != nil {
			return err
		}

		if err := arg.hashPassword(); err != nil {
			return err
		}
	}

	_, err := db.AuthUpdateUser(ctx, repo.AuthUpdateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  arg.PasswordHash,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        arg.ID,
	})
	return err
}

func UpdateUserDisable(ctx context.Context, db repo.DB, userID int64, disable bool) error {
	if disable {
		_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         userID,
		})
		return err
	}
	_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         userID,
	})
	return err
}

func UpdateUserAdmin(ctx context.Context, db repo.DB, userID int64, admin bool) error {
	if admin {
		_, err := db.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    userID,
			CreatedAt: types.NewTime(time.Now()),
		})
		return err
	}
	return db.AuthDeleteAdmin(ctx, userID)
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
