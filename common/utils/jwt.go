package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

type (
	VerifyTokenCode int
)

const (
	VerifyTokenCodeOK VerifyTokenCode = iota
	VerifyTokenCodeError
	VerifyTokenCodeExpire
	VerifyTokenCodeBaned
)

func VerifyToken(ctx context.Context, rc redis.UniversalClient, token string) (VerifyTokenCode, string) {
	// todo add your logic here and delete this line
	return VerifyTokenCodeOK, strings.TrimPrefix(token, "token.uid.")
}
