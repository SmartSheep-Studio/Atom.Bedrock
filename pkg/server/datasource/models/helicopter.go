package models

import (
	"os"
	"time"
)

type SubAppManifest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Workdir     string   `json:"workdir"`
	Executable  string   `json:"executable"`
	Arguments   []string `json:"arguments"`
	Environment []string `json:"environment"`
	Order       int      `json:"order"`
}

type SubAppExposedOptions struct {
	URL   string              `json:"url"`
	Pages []SubAppExposedPage `json:"pages"`
}

type SubAppExposedPage struct {
	Icon  string         `json:"icon"`
	Name  string         `json:"name"`
	Title string         `json:"title"`
	Path  string         `json:"path"`
	Meta  map[string]any `json:"meta"`
}

type SubApp struct {
	Manifest *SubAppManifest `json:"manifest"`

	Process        *os.Process           `json:"-"`
	ExposedOptions *SubAppExposedOptions `json:"exposed_options"`
	LastHealthyAt  *time.Time            `json:"last_healthy_at"`
}
