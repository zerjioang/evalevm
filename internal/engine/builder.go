package engine

import (
	"context"
	"evalevm/internal/datatype"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Builder struct{}

func (b *Builder) Build(ctx context.Context, toolsPath string, analyzer datatype.Analyzer, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	log.Println("running docker build command...")

	dockerImageName := fmt.Sprintf("local/%s:latest", analyzer.Name())
	dockerfilePath, err := analyzer.DockerfilePath()
	if err != nil {
		return err
	}
	toolRepoRootPath := toolsPath + "/" + analyzer.Name()

	// Check if image exists locally
	imageExists := false
	inspectCmd := exec.CommandContext(ctx, "docker", "image", "inspect", dockerImageName)
	if err := inspectCmd.Run(); err == nil {
		imageExists = true
	}

	if imageExists && !force {
		log.Printf("Docker image %s already exists locally. Skipping build.", dockerImageName)
		return nil
	}

	buildArgs := []string{
		"build",
		"--platform", analyzer.DockerPlatform(),
		"-t", dockerImageName,
		"-f", dockerfilePath, toolRepoRootPath,
	}

	log.Println(dockerfilePath, toolRepoRootPath)
	log.Println(buildArgs)

	return RunDockerCommand(ctx, buildArgs, func(line string, isStderr bool) {
		if isStderr {
			fmt.Printf("err: %s\n", line)
		} else {
			fmt.Printf("out: %s\n", line)
		}
	})
}
