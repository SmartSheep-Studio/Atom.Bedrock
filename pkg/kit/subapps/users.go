package subapps

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type HeLiCoPtErAccountResp struct {
	User    models.User        `json:"user"`
	Session models.UserSession `json:"session"`
	Claims  models.UserClaims  `json:"claims"`
}

func (conn *HeLiCoPtErConnection) GetAccount(token string) (HeLiCoPtErAccountResp, error) {
	client := resty.New()
	res, err := client.R().
		SetAuthToken(token).
		SetResult(&HeLiCoPtErAccountResp{}).
		Get(conn.GetEndpointPath("/api/users/self"))

	if err != nil {
		return HeLiCoPtErAccountResp{}, err
	} else if res.StatusCode() != 200 {
		return HeLiCoPtErAccountResp{}, fmt.Errorf("failed to get principal: %s", string(res.Body()))
	} else {
		ply := res.Result().(*HeLiCoPtErAccountResp)
		return *ply, nil
	}
}

func (conn *HeLiCoPtErConnection) GetAccountWithID(id uint) (HeLiCoPtErAccountResp, error) {
	client := resty.New()
	res, err := client.R().
		SetResult(&HeLiCoPtErAccountResp{}).
		Get(conn.GetEndpointPath("/cgi/users/" + strconv.Itoa(int(id))))

	if err != nil {
		return HeLiCoPtErAccountResp{}, err
	} else if res.StatusCode() != 200 {
		return HeLiCoPtErAccountResp{}, fmt.Errorf("failed to get principal: %s", string(res.Body()))
	} else {
		ply := res.Result().(*HeLiCoPtErAccountResp)
		return *ply, nil
	}
}
