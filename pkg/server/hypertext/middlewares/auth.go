package middlewares

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	services "code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type AuthHandler func(force bool, scope []string, perms []string) fiber.Handler

type AuthMiddleware struct {
	auth  *services.AuthService
	users *services.UserService

	Fn AuthHandler
}

type AuthConfig struct {
	Next        func(c *fiber.Ctx) bool
	LookupToken string
}

func NewAuth(auth *services.AuthService, users *services.UserService) *AuthMiddleware {
	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: authorization",
	}

	inst := &AuthMiddleware{auth: auth, users: users}
	inst.Fn = func(force bool, scope []string, perms []string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			if cfg.Next != nil && cfg.Next(c) {
				return c.Next()
			}

			err := inst.LookupAuthToken(c, strings.Split(cfg.LookupToken, ","))
			if err != nil && force {
				return fiber.NewError(fiber.StatusUnauthorized, err.Error())
			} else {
				if err == nil {
					for _, lock := range c.Locals("principal").(models.User).Locks {
						if lock.ExpiredAt == nil || lock.ExpiredAt.Unix() < time.Now().Unix() {
							return fiber.NewError(
								fiber.StatusForbidden,
								fmt.Sprintf("your account has been locked, reason: %s", lock.Reason),
							)
						}
					}

					if err := c.Locals("principal-session").(models.UserSession).HasScope(scope...); err != nil {
						return fiber.NewError(fiber.StatusForbidden, err.Error())
					} else if err := c.Locals("principal").(models.User).HasPermissions(perms...); err != nil {
						return fiber.NewError(fiber.StatusForbidden, err.Error())
					}
				}

				c.Locals("principal-ok", err == nil)
			}

			return c.Next()
		}
	}

	return inst
}

func (v *AuthMiddleware) LookupAuthToken(c *fiber.Ctx, args []string) error {
	var str string
	for _, arg := range args {
		parts := strings.Split(strings.TrimSpace(arg), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "header":
			if len(c.GetReqHeaders()[v]) > 0 {
				str = strings.TrimSpace(strings.ReplaceAll(c.GetReqHeaders()[v], "Bearer", ""))
			}
		case "query":
			if len(c.Query(v)) > 0 {
				str = c.Query(v)
			}
		case "cookie":
			if len(c.Cookies(v)) > 0 {
				str = c.Cookies(v)
			}
		}
	}

	if len(str) == 0 {
		return fmt.Errorf("missing token in request")
	}

	claims, err := v.auth.ReadJwt(str)
	if err != nil {
		return fmt.Errorf("failed to parse token: %q", err)
	}
	session, user, err := v.auth.ReadClaims(*claims)
	if err != nil {
		return fmt.Errorf("failed to read details: %q", err)
	}

	c.Locals("principal", user)
	c.Locals("principal-claims", *claims)
	c.Locals("principal-session", session)

	return nil
}
