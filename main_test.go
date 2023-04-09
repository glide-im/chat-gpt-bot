package main

import (
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"testing"
)

func TestGPT2(t *testing.T) {
}

func TestGenerateSecret(t *testing.T) {
	jwtAuthorize := jwt_auth.NewAuthorizeImpl("")
	token, _ := jwtAuthorize.GetToken(&jwt_auth.JwtAuthInfo{
		UID:         "",
		Device:      "1",
		ExpiredHour: 23,
	})
	t.Log(token.Token)
}
