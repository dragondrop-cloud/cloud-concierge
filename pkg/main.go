package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Entrypoint on go binary")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	err := RemoveSubDirectories()
	if err != nil {
		log.Errorf("Error removing sub directories: %s", err.Error())
		os.Exit(1)
	}

	env := os.Getenv("DRAGONDROP_EXECUTION_ENVIRONMENT")
	job, err := InitializeJobDependencies(ctx, env)
	if err != nil {
		log.Errorf("Error creating job: %s", err.Error())
		os.Exit(1)
	}

	err = job.Authorize(ctx)
	if err != nil {
		log.Errorf("Error authorizing job: %s", err.Error())
		os.Exit(1)
	}

	err = job.Run(ctx)
	if err != nil {
		log.Errorf("Error running job: %s", err.Error())
		os.Exit(1)
	}

	log.Info("Done executing go binary")
}

// RemoveSubDirectories removes all subdirectories within the container's volume prior to container startup
func RemoveSubDirectories() error {
	if _, err := os.Stat("/driftmitigation/"); err == nil {
		d, err := os.Open("/driftmitigation/")
		if err != nil {
			return fmt.Errorf("[os.Open('/driftmitigation/)]%v", err)
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			return fmt.Errorf("[d.Readdirnames(-1)]%v", err)
		}
		fmt.Printf("All sub directories identified:\n%v\n", names)

		for _, name := range names {
			err = os.RemoveAll(filepath.Join("/driftmitigation/", name))
			if err != nil {
				return fmt.Errorf("[os.RemoveAll(/driftmitigation/%v)]%v", name, err)
			}
		}
	}
	return nil
}
