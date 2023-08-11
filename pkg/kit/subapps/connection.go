package subapps

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"os"
)

type HeLiCoPtErConnection struct {
	Endpoint      string          `json:"endpoint"`
	OtherApps     []models.SubApp `json:"others"`
	Configuration map[string]any  `json:"configuration"`
}

func PublishApp(url string, name string) (*HeLiCoPtErConnection, error) {
	endpoint := os.Getenv("BEDROCK_ENDPOINT_URL")
	if len(endpoint) == 0 {
		return nil, fmt.Errorf("couldn't get endpoint from environment variables")
	}

	client := resty.New()
	res, err := client.R().
		SetBody(fiber.Map{
			"url": url,
		}).
		SetResult(&HeLiCoPtErConnection{}).
		Post(fmt.Sprintf("%s/cgi/subapps/%s", endpoint, name))

	if err != nil {
		return nil, err
	} else if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to publish app into endpoint: %s", string(res.Body()))
	} else {
		res := res.Result().(*HeLiCoPtErConnection)
		res.Endpoint = endpoint

		return res, nil
	}
}

func (conn *HeLiCoPtErConnection) GetEndpointPath(path string) string {
	return fmt.Sprintf("%s%s", conn.Endpoint, path)
}

func (conn *HeLiCoPtErConnection) GetOtherSubApp(name string) (models.SubApp, bool) {
	return lo.Find(conn.OtherApps, func(item models.SubApp) bool {
		return item.Manifest.Name == name
	})
}

func (conn *HeLiCoPtErConnection) GetOtherSubAppEndpoint(name string) (string, bool) {
	if app, ok := conn.GetOtherSubApp(name); ok {
		return app.ExposedOptions.URL, ok
	} else {
		return "404", ok
	}
}

func (conn *HeLiCoPtErConnection) GetOtherSubAppEndpointPath(name string, path string) (string, bool) {
	if address, ok := conn.GetOtherSubAppEndpoint(name); !ok {
		return "404", ok
	} else {
		return fmt.Sprintf("%s%s", address, path), ok
	}
}
