package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/global"
)

const (
	WatchNet = "WatchNet"
	TimeOut  = 60 * time.Second
)

func Close(ctx context.Context, conn *dahuarpc.Conn) error {
	_, err := global.Logout(ctx, conn)
	conn.UpdateState(dahuarpc.StateClosed)
	return err
}

func KeepAlive(ctx context.Context, conn *dahuarpc.Conn) error {
	if time.Now().Sub(conn.Data().LastLogin) > TimeOut {
		_, err := global.KeepAlive(ctx, conn)
		if err != nil {
			if !errors.Is(err, dahuarpc.ErrRequestFailed) {
				conn.UpdateState(dahuarpc.StateLogout)
			}

			return err
		}

		conn.UpdateState(dahuarpc.StateLogin)
	}

	return nil
}

func Login(ctx context.Context, conn *dahuarpc.Conn, username, password string) error {
	if err := login(ctx, conn, username, password); err != nil {
		var e *LoginError
		if errors.As(err, &e) {
			conn.UpdateState(dahuarpc.StateError, err)
		} else {
			conn.UpdateState(dahuarpc.StateLogout)
		}

		return err
	}

	conn.UpdateState(dahuarpc.StateLogin)

	return nil
}

func login(ctx context.Context, conn *dahuarpc.Conn, username, password string) error {
	firstLogin, err := global.FirstLogin(ctx, conn, username)
	if err != nil {
		return err
	}
	if firstLogin.Error == nil {
		return fmt.Errorf("FirstLogin did not return an error")
	}
	if !(firstLogin.Error.Code == 268632079 || firstLogin.Error.Code == 401) {
		return fmt.Errorf("FirstLogin has invalid error code: %d", firstLogin.Error.Code)
	}

	// Update session
	if err := conn.UpdateSession(firstLogin.Session.String()); err != nil {
		return err
	}

	// Magic
	loginType := func() string {
		if firstLogin.Params.Encryption == WatchNet {
			return WatchNet
		}
		return "Direct"
	}()

	// Encrypt password based on the first login and then do a second login
	passwordHash := firstLogin.Params.HashPassword(username, password)
	err = global.SecondLogin(ctx, conn, username, passwordHash, loginType, firstLogin.Params.Encryption)
	if err != nil {
		var responseErr *dahuarpc.ResponseError
		if errors.As(err, &responseErr) {
			if loginErr := intoLoginError(responseErr); loginErr != nil {
				return errors.Join(loginErr, err)
			}
		}

		return err
	}

	return nil
}

func intoLoginError(r *dahuarpc.ResponseError) *LoginError {
	switch r.Code {
	case 268632085:
		return &ErrLoginUserOrPasswordNotValid
	case 268632081:
		return &ErrLoginHasBeenLocked
	}

	switch r.Message {
	case "UserNotValidt":
		return &ErrLoginUserNotValid
	case "PasswordNotValid":
		return &ErrLoginPasswordNotValid
	case "InBlackList":
		return &ErrLoginInBlackList
	case "HasBeedUsed":
		return &ErrLoginHasBeedUsed
	case "HasBeenLocked":
		return &ErrLoginHasBeenLocked
	}

	return nil
}

type LoginError struct {
	Message string
}

func newLoginError(message string) LoginError {
	return LoginError{
		Message: message,
	}
}

func (e *LoginError) Error() string {
	return e.Message
}

var (
	ErrLoginUserOrPasswordNotValid = newLoginError("User or password not valid")
	ErrLoginUserNotValid           = newLoginError("User not valid")
	ErrLoginPasswordNotValid       = newLoginError("Password not valid")
	ErrLoginInBlackList            = newLoginError("User in blackList")
	ErrLoginHasBeedUsed            = newLoginError("User has be used")
	ErrLoginHasBeenLocked          = newLoginError("User locked")
)
