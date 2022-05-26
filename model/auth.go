package model

import (
	"context"
)

type Authentication struct {
	UserName string
	Password string
}

type Auth struct {
	User string
}

func (a *Auth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {

	return map[string]string{"user": a.User}, nil
}

func (a *Auth) RequireTransportSecurity() bool {

	return false
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"user-name": a.UserName, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
