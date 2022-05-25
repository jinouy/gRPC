package model

import (
	"context"
	"gRPC_User/client/auth"
)

type Authentication struct {
	User     string
	Password string
}

type Auth struct {
	User string
}

func (a *Auth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {

	user := Auth{User: auth.InputName()}
	return map[string]string{"user": user.User}, nil
}

func (a *Auth) RequireTransportSecurity() bool {

	return false
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
