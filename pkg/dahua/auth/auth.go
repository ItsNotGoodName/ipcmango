package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/pkg/dahua"
	"github.com/ItsNotGoodName/pkg/dahua/modules/global"
)

const (
	WatchNet = "WatchNet"
	TimeOut  = 60 * time.Second
)

func Logout(conn *dahua.Conn) {
	global.Logout(conn)
	conn.Set(dahua.StateLogout)
}

func KeepAlive(conn *dahua.Conn) (bool, error) {
	if time.Now().Sub(conn.LastLogin) > TimeOut {
		_, err := global.KeepAlive(conn)
		if err != nil {
			if !errors.Is(err, dahua.ErrRequestFailed) {
				conn.Set(dahua.StateLogout)
				return false, nil
			}

			return true, err
		}

		conn.Set(dahua.StateLogin)
	}

	return true, nil
}

func Login(conn *dahua.Conn, username, password string) error {
	if conn.State == dahua.StateLogin {
		Logout(conn)
	} else if conn.State == dahua.StateError {
		return conn.Error
	}

	err := login(conn, username, password)
	if err != nil {
		var e *LoginError
		if errors.As(err, &e) {
			conn.SetError(err)
		} else {
			conn.Set(dahua.StateLogout)
		}

		return err
	}

	conn.Set(dahua.StateLogin)

	return nil
}

func login(conn *dahua.Conn, username, password string) error {
	// Do a first login
	firstLogin, err := global.FirstLogin(conn, username)
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
	if err := conn.UpdateSession(firstLogin.Session.Value); err != nil {
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
	err = global.SecondLogin(conn, username, passwordHash, loginType, firstLogin.Params.Encryption)
	if err != nil {
		var responseErr *dahua.ResponseError
		if errors.As(err, &responseErr) {
			if loginErr := intoLoginError(responseErr); loginErr != nil {
				return errors.Join(loginErr, err)
			}
		}

		return err
	}

	return nil
}

func intoLoginError(r *dahua.ResponseError) *LoginError {
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

func newErrLogin(message string) LoginError {
	return LoginError{
		Message: message,
	}
}

func (e *LoginError) Error() string {
	return e.Message
}

var (
	ErrLoginClosed                 = newErrLogin("Client is closed")
	ErrLoginUserOrPasswordNotValid = newErrLogin("User or password not valid")
	ErrLoginUserNotValid           = newErrLogin("User not valid")
	ErrLoginPasswordNotValid       = newErrLogin("Password not valid")
	ErrLoginInBlackList            = newErrLogin("User in blackList")
	ErrLoginHasBeedUsed            = newErrLogin("User has be used")
	ErrLoginHasBeenLocked          = newErrLogin("User locked")
)
