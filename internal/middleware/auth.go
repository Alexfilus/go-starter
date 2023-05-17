package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-pkg/response"
)

func AuthRequired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := ExtractTokenFromHeader(ctx)
		if err != nil {
			rsp := response.Unauthorized(err.Error(), "")
			return ctx.Status(http.StatusUnauthorized).JSON(rsp)
		}

		userInfo, err := ValidateToken(token)
		if err != nil {
			rsp := response.Unauthorized(err.Error(), "")
			return ctx.Status(http.StatusUnauthorized).JSON(rsp)
		}

		IdentityInContext(ctx, userInfo)
		return ctx.Next()
	}
}

// ExtractTokenFromHeader from a header
func ExtractTokenFromHeader(ctx *fiber.Ctx) (string, error) {
	authorizationHeader := strings.TrimSpace(ctx.Get("Authorization"))

	if authorizationHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Bearer 38ej93lee0e--e9ikdj
	token := strings.TrimPrefix(authorizationHeader, "Bearer ")

	if token == authorizationHeader {
		return "", errors.New("malformed or unsupported authorization header")
	}

	return token, nil
}

// ValidateToken validates a token
func ValidateToken(token string) (bool, error) {
	return false, nil
}

const contextKey = "__CTX_identity_id"

// IdentityInContext: Places the  user identity in a context
func IdentityInContext(ctx *fiber.Ctx, identity any) {
	ctx.Locals(contextKey, identity)
}

// IdentityFromContext gets the user identity from a context
func IdentityFromContext(ctx *fiber.Ctx) any {
	return ctx.Locals(contextKey)
}
