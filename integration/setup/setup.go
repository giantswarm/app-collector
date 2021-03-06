// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/apptest"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/app-exporter/integration/env"
	"github.com/giantswarm/app-exporter/integration/key"
	"github.com/giantswarm/app-exporter/integration/templates"
	"github.com/giantswarm/app-exporter/pkg/project"
)

func Setup(m *testing.M, config Config) {
	var err error

	ctx := context.Background()

	err = installResources(ctx, config)
	if err != nil {
		config.Logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("failed to install %#q", project.Name()), "stack", fmt.Sprintf("%#v", err))
		os.Exit(-1)
	}

	os.Exit(m.Run())
}

func installResources(ctx context.Context, config Config) error {
	var err error

	{
		apps := []apptest.App{
			{
				CatalogName:   key.ControlPlaneTestCatalogName(),
				Name:          project.Name(),
				Namespace:     key.Namespace(),
				SHA:           env.CircleSHA(),
				ValuesYAML:    templates.AppExporterValues,
				WaitForDeploy: true,
			},
		}
		err = config.AppTest.InstallApps(ctx, apps)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
