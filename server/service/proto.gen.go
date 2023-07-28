// api v1.0.0 6946c2ef0345f2f58437b0e43106d2858b4b4380
// --
// Code generated by webrpc-gen@v0.12.0 with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=./server/api.ridl -target=golang -pkg=service -server -out=./server/service/proto.gen.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v1.0.0"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "6946c2ef0345f2f58437b0e43106d2858b4b4380"
}

//
// Types
//

type UserRegister struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type AuthService interface {
	Register(ctx context.Context, user *UserRegister) error
	Login(ctx context.Context, usernameOrEmail string, password string) (string, error)
}

type UserService interface {
	Me(ctx context.Context) (*User, error)
}

var WebRPCServices = map[string][]string{
	"AuthService": {
		"Register",
		"Login",
	},
	"UserService": {
		"Me",
	},
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type authServiceServer struct {
	AuthService
}

func NewAuthServiceServer(svc AuthService) WebRPCServer {
	return &authServiceServer{
		AuthService: svc,
	}
}

func (s *authServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "AuthService")

	if r.Method != "POST" {
		err := ErrorWithCause(ErrWebrpcBadMethod, fmt.Errorf("unsupported method %q (only POST is allowed)", r.Method))
		RespondWithError(w, err)
		return
	}

	switch r.URL.Path {
	case "/rpc/AuthService/Register":
		s.serveRegister(ctx, w, r)
		return
	case "/rpc/AuthService/Login":
		s.serveLogin(ctx, w, r)
		return
	default:
		err := ErrorWithCause(ErrWebrpcBadRoute, fmt.Errorf("no handler for path %q", r.URL.Path))
		RespondWithError(w, err)
		return
	}
}

func (s *authServiceServer) serveRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveRegisterJSON(ctx, w, r)
	default:
		err := ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		RespondWithError(w, err)
	}
}

func (s *authServiceServer) serveRegisterJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Register")
	reqContent := struct {
		Arg0 *UserRegister `json:"user"`
	}{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("failed to read request data: %w", err))
		RespondWithError(w, err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, &reqContent)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("failed to unmarshal request data: %w", err))
		RespondWithError(w, err)
		return
	}

	// Call service method
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorWithCause(ErrWebrpcServerPanic, fmt.Errorf("%v", rr)))
				panic(rr)
			}
		}()
		err = s.AuthService.Register(ctx, reqContent.Arg0)
	}()

	if err != nil {
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *authServiceServer) serveLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveLoginJSON(ctx, w, r)
	default:
		err := ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		RespondWithError(w, err)
	}
}

func (s *authServiceServer) serveLoginJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Login")
	reqContent := struct {
		Arg0 string `json:"usernameOrEmail"`
		Arg1 string `json:"password"`
	}{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("failed to read request data: %w", err))
		RespondWithError(w, err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, &reqContent)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("failed to unmarshal request data: %w", err))
		RespondWithError(w, err)
		return
	}

	// Call service method
	var ret0 string
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorWithCause(ErrWebrpcServerPanic, fmt.Errorf("%v", rr)))
				panic(rr)
			}
		}()
		ret0, err = s.AuthService.Login(ctx, reqContent.Arg0, reqContent.Arg1)
	}()
	respContent := struct {
		Ret0 string `json:"token"`
	}{ret0}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to marshal json response: %w", err))
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

type userServiceServer struct {
	UserService
}

func NewUserServiceServer(svc UserService) WebRPCServer {
	return &userServiceServer{
		UserService: svc,
	}
}

func (s *userServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "UserService")

	if r.Method != "POST" {
		err := ErrorWithCause(ErrWebrpcBadMethod, fmt.Errorf("unsupported method %q (only POST is allowed)", r.Method))
		RespondWithError(w, err)
		return
	}

	switch r.URL.Path {
	case "/rpc/UserService/Me":
		s.serveMe(ctx, w, r)
		return
	default:
		err := ErrorWithCause(ErrWebrpcBadRoute, fmt.Errorf("no handler for path %q", r.URL.Path))
		RespondWithError(w, err)
		return
	}
}

func (s *userServiceServer) serveMe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveMeJSON(ctx, w, r)
	default:
		err := ErrorWithCause(ErrWebrpcBadRequest, fmt.Errorf("unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		RespondWithError(w, err)
	}
}

func (s *userServiceServer) serveMeJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Me")

	// Call service method
	var ret0 *User
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorWithCause(ErrWebrpcServerPanic, fmt.Errorf("%v", rr)))
				panic(rr)
			}
		}()
		ret0, err = s.UserService.Me(ctx)
	}()
	respContent := struct {
		Ret0 *User `json:"user"`
	}{ret0}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		err = ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to marshal json response: %w", err))
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(WebRPCError)
	if !ok {
		rpcErr = ErrorWithCause(ErrWebrpcEndpoint, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

//
// Helpers
//

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	// For Client
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}

	// For Server
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	err := rpcErr
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Webrpc errors
var (
	ErrWebrpcEndpoint      = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 0}
	ErrWebrpcBadRoute      = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod     = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest    = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse   = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic   = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
)

// Schema errors
var (
	ErrInvalidToken = WebRPCError{Code: 100, Name: "InvalidToken", Message: "invalid token", HTTPStatus: 401}
)
