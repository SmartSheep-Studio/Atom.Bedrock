package subapps

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
)

func (conn *HeLiCoPtErConnection) GetConnectURL(id string, callback string) string {
	return fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=all",
		conn.GetEndpointPath("/auth/openid/connect"),
		id,
		callback,
	)
}

func (conn *HeLiCoPtErConnection) GetAccessToken(code string, id string, secret string, callback string) (string, models.UserSession, error) {
	type reply struct {
		AccessToken  string             `json:"access_token"`
		RefreshToken string             `json:"refresh_token"`
		RedirectURI  string             `json:"redirect_uri"`
		Session      models.UserSession `json:"session"`
	}

	client := resty.New()
	res, err := client.R().SetBody(fiber.Map{
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  callback,
		"client_id":     id,
		"client_secret": secret,
		"scope":         "all",
	}).SetResult(&reply{}).Post(conn.GetEndpointPath("/api/auth/openid/exchange"))

	if err != nil {
		return "", models.UserSession{}, err
	} else if res.StatusCode() != 200 {
		return "", models.UserSession{}, fmt.Errorf("failed to exchange access token: %s", string(res.Body()))
	} else {
		ply := res.Result().(*reply)
		return ply.AccessToken, ply.Session, nil
	}
}
