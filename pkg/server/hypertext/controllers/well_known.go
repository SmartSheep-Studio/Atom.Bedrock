package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type WellKnownController struct{}

func NewWellKnownController() *WellKnownController {
	return &WellKnownController{}
}

func (ctrl *WellKnownController) Map(router *fiber.App) {
	router.Get(
		"/.well-known/oauth-authorization-hypertext",
		ctrl.oauth,
	)
	router.Get(
		"/.well-known/openid-configuration",
		ctrl.openid,
	)
}

func (ctrl *WellKnownController) oauth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"issuer":                   viper.GetString("general.base_url"),
		"authorization_endpoint":   fmt.Sprintf("%s/users/openid/connect", viper.GetString("base_url")),
		"token_endpoint":           fmt.Sprintf("%s/api/users/openid/exchange", viper.GetString("base_url")),
		"end_session_endpoint":     fmt.Sprintf("%s/api/users/sign-out", viper.GetString("base_url")),
		"response_types_supported": []string{"code", "token"},
		"grant_types_supported":    []string{"authorization_code", "implicit", "refresh_token", "password"},
		"subject_types_supported":  []string{"public"},
		"scopes_supported": []string{
			"principal",
			"sessions.read",
			"sessions.create",
			"oauth.read",
			"oauth.create",
			"oauth.update",
			"oauth.delete",
			"oauth.openid",
			"oauth.openid.approve",
			"assets.create",
			"users.tokens.read",
			"users.tokens.create",
			"users.update",
			"users.update.avatar",
			"users.update.banner",
		},
	})
}

func (ctrl *WellKnownController) openid(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"issuer":                   viper.GetString("general.base_url"),
		"authorization_endpoint":   fmt.Sprintf("%s/users/openid/connect", viper.GetString("base_url")),
		"token_endpoint":           fmt.Sprintf("%s/api/users/openid/exchange", viper.GetString("base_url")),
		"userinfo_endpoint":        fmt.Sprintf("%s/api/users", viper.GetString("base_url")),
		"end_session_endpoint":     fmt.Sprintf("%s/api/users/sign-out", viper.GetString("base_url")),
		"response_types_supported": []string{"code", "token"},
		"grant_types_supported":    []string{"authorization_code", "implicit", "refresh_token", "password"},
		"subject_types_supported":  []string{"public"},
		"scopes_supported": []string{
			"principal",
			"sessions.read",
			"sessions.create",
			"oauth.read",
			"oauth.create",
			"oauth.update",
			"oauth.delete",
			"oauth.openid",
			"oauth.openid.approve",
			"assets.create",
			"users.tokens.read",
			"users.tokens.create",
			"users.update",
			"users.update.avatar",
			"users.update.banner",
		},
	})
}
