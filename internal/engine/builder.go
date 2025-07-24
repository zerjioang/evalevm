package engine

import (
	"context"
	"evalevm/internal/datatype"
	"fmt"
	"log"
	"time"
)

type Builder struct{}

func (b *Builder) Build(ctx context.Context, toolsPath string, analyzer datatype.Analyzer) error {
	// docker build -t conkas:latest .
	// docker run --rm -v "$(pwd)":/work -w /work conkas:latest \
	// python3 /opt/conkas/conkas.py contract.bin

	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()

	log.Println("Running Docker command...")

	dockerImageName := fmt.Sprintf("local/%s:latest", analyzer.Name())
	dockerfilePath, err := analyzer.DockerfilePath()
	if err != nil {
		return err
	}
	toolRepoRootPath := toolsPath + "/" + analyzer.Name()
	args := []string{
		"build",
		"--platform", analyzer.DockerPlatform(),
		"-t", dockerImageName,
		"-f", dockerfilePath, toolRepoRootPath,
	}

	log.Println(dockerfilePath, toolRepoRootPath)
	log.Println(args)

	return RunDockerCommand(ctx, args, func(line string, isStderr bool) {
		if isStderr {
			fmt.Printf("ERR: %s\n", line)
		} else {
			fmt.Printf("OUT: %s\n", line)
		}
	})
}
