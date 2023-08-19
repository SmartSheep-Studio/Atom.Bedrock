package subapps

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
)

func (conn *HeLiCoPtErConnection) UploadAssets(name string, reader io.Reader) (models.StorageFile, error) {
	client := resty.New()
	res, err := client.R().
		SetResult(&models.StorageFile{}).
		SetFileReader("file", name, reader).
		Post(conn.GetEndpointPath("/cgi/assets"))

	if err != nil {
		return models.StorageFile{}, err
	} else if res.StatusCode() != 200 {
		return models.StorageFile{}, fmt.Errorf("failed to get principal: %s", string(res.Body()))
	} else {
		ply := res.Result().(*models.StorageFile)
		return *ply, nil
	}
}

func (conn *HeLiCoPtErConnection) UploadAssets2User(tk string, name string, reader io.Reader) (models.StorageFile, error) {
	client := resty.New()
	res, err := client.R().
		SetResult(&models.StorageFile{}).
		SetAuthToken(tk).
		SetFileReader("file", name, reader).
		Post(conn.GetEndpointPath("/cgi/assets"))

	if err != nil {
		return models.StorageFile{}, err
	} else if res.StatusCode() != 200 {
		return models.StorageFile{}, fmt.Errorf("failed to get principal: %s", string(res.Body()))
	} else {
		ply := res.Result().(*models.StorageFile)
		return *ply, nil
	}
}
