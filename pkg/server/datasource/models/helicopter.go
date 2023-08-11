package models

type SubApp struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	Workdir      string   `json:"workdir"`
	Executable   string   `json:"executable"`
	Arguments    []string `json:"arguments"`
	Environments []string `json:"environments"`
	Order        int      `json:"order"`
}
