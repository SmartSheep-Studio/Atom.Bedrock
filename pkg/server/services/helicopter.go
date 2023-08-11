package services

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"os/exec"
	"strings"
)

type HeLiCoPtErService struct {
	log  zerolog.Logger
	conf *viper.Viper

	Apps []*models.SubApp
}

func NewSubApp(manifest *models.SubAppManifest) *models.SubApp {
	return &models.SubApp{manifest, nil, nil, nil}
}

func NewHeLiCoPtErService(cycle fx.Lifecycle, log zerolog.Logger, conf *viper.Viper) *HeLiCoPtErService {
	inst := &HeLiCoPtErService{log, conf, []*models.SubApp{}}

	// Load config and parse apps
	manifests, _ := inst.GetSubApps()
	inst.Apps = lo.Map(manifests, func(item *models.SubAppManifest, index int) *models.SubApp {
		return NewSubApp(item)
	})

	// Hook into lifecycles to start subapps automatically
	cycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return inst.StopAll()
		},
	})

	return inst
}

func (v *HeLiCoPtErService) GetSubApps() ([]*models.SubAppManifest, error) {
	var apps []*models.SubAppManifest
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
	v.conf.Set("helicopter.subapps", v.Apps)
	return v.conf.SafeWriteConfig()
}

func (v *HeLiCoPtErService) NewSubApp(item models.SubAppManifest) error {
	v.Apps = append(v.Apps, NewSubApp(lo.ToPtr(item)))
	return v.SaveSubApps()
}

func (v *HeLiCoPtErService) StartOne(app *models.SubApp) error {
	// Genshin Impact, Boot!
	cmd := exec.Command(app.Manifest.Executable, app.Manifest.Arguments...)
	cmd.Dir = app.Manifest.Workdir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(
		app.Manifest.Environment,
		append(
			os.Environ(),
			fmt.Sprintf(
				"BEDROCK_ENDPOINT_URL=http://127.0.0.1:%s",
				strings.SplitN(v.conf.GetString("hypertext.bind_addr"), ":", 2)[1],
			),
		)...,
	)

	v.Apps[lo.IndexOf(v.Apps, app)].Process = cmd.Process

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
	for _, app := range v.Apps {
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

func (v *HeLiCoPtErService) StopOne(app *models.SubApp) error {
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
	for _, app := range v.Apps {
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
