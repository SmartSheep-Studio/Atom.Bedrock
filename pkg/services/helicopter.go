package services

import (
	"code.smartsheep.studio/atom/bedrock/pkg/datasource/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"os/exec"
)

type HeLiCoPtErService struct {
	log  zerolog.Logger
	conf *viper.Viper

	apps []*HeLiCoPtErSubApp
}

type HeLiCoPtErSubApp struct {
	Manifest *models.SubApp `json:"manifest"`

	Process *os.Process
}

func NewHeLiCoPtErSubApp(manifest *models.SubApp) *HeLiCoPtErSubApp {
	return &HeLiCoPtErSubApp{manifest, nil}
}

func NewHeLiCoPtErService(cycle fx.Lifecycle, log zerolog.Logger, conf *viper.Viper) *HeLiCoPtErService {
	inst := &HeLiCoPtErService{log, conf, []*HeLiCoPtErSubApp{}}

	// Load config and parse apps
	manifests, _ := inst.GetSubApps()
	inst.apps = lo.Map(manifests, func(item *models.SubApp, index int) *HeLiCoPtErSubApp {
		return NewHeLiCoPtErSubApp(item)
	})

	// Hook into lifecycles to start subapps automatically
	cycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return inst.StopAll()
		},
	})

	return inst
}

func (v *HeLiCoPtErService) GetSubApps() ([]*models.SubApp, error) {
	var apps []*models.SubApp
	if raw, err := json.Marshal(v.conf.Get("helicopter.subapps")); err != nil {
		return apps, err
	} else {
		if err := json.Unmarshal(raw, &apps); err != nil {
			return apps, err
		}
		return apps, nil
	}
}

func (v *HeLiCoPtErService) SaveSubApps() error {
	v.conf.Set("helicopter.subapps", v.apps)
	return v.conf.SafeWriteConfig()
}

func (v *HeLiCoPtErService) NewSubApp(item models.SubApp) error {
	v.apps = append(v.apps, NewHeLiCoPtErSubApp(lo.ToPtr(item)))
	return v.SaveSubApps()
}

func (v *HeLiCoPtErService) StartOne(app *HeLiCoPtErSubApp) error {
	// Genshin Impact, Boot!
	cmd := exec.Command(app.Manifest.Executable, app.Manifest.Arguments...)
	cmd.Dir = app.Manifest.Workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	v.apps[lo.IndexOf(v.apps, app)].Process = cmd.Process

	go func() {
		err := cmd.Run()
		if err != nil {
			v.log.Err(err).Msgf("HeLiCoPtEr sub app %s quit unexpectedly.", app.Manifest.Name)
		}
	}()

	return nil
}

func (v *HeLiCoPtErService) StartAll() error {
	v.log.Info().Msg("HeLiCoPtEr is starting...")
	for _, app := range v.apps {
		v.log.Info().Msgf("HeLiCoPtEr is starting %s...", app.Manifest.Name)
		if err := v.StartOne(app); err != nil {
			v.log.Info().Msgf("HeLiCoPtEr failed to start %s... %q", app.Manifest.Name, err)
		} else {
			v.log.Info().Msgf("HeLiCoPtEr successfully started %s!", app.Manifest.Name)
		}
	}
	v.log.Info().Msg("HeLiCoPtEr is started!")
	return nil
}

func (v *HeLiCoPtErService) StopOne(app *HeLiCoPtErSubApp) error {
	// Genshin Impact, Stop!
	if app.Process == nil {
		return fmt.Errorf("app %s isn't started yet", app.Manifest.Name)
	}

	if err := app.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to quit app %s: %q", app.Manifest.Name, err)
	} else {
		return nil
	}
}

func (v *HeLiCoPtErService) StopAll() error {
	v.log.Info().Msg("HeLiCoPtEr is stopping...")
	for _, app := range v.apps {
		v.log.Info().Msgf("HeLiCoPtEr now stopping %s...", app.Manifest.Name)
		if err := v.StopOne(app); err != nil {
			v.log.Info().Msgf("HeLiCoPtEr failed to stop %s... %q", app.Manifest.Name, err)
		} else {
			v.log.Info().Msgf("HeLiCoPtEr successfully stopped %s!", app.Manifest.Name)
		}
	}
	v.log.Info().Msg("HeLiCoPtEr is stopped!")
	return nil
}
