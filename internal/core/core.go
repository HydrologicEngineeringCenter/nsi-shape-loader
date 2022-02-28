package core

import (
	"errors"
	"fmt"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/store"
	"github.com/urfave/cli/v2"
)

func Upload(c *cli.Context) error {
	cfg, err := config.NewConfig(c)
	if err != nil {
		return errors.New("core upload: invalid config")
	}
	_, err = store.NewStore(cfg)
	if err != nil {
		return errors.New("core upload: invalid store")
	}
	fmt.Println(cfg)
	return nil
}
