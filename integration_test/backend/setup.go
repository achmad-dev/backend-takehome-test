package backend

import (
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Load docker compose for integration test in root path
func SetupDockerAppTest() (tc testcontainers.DockerCompose, err error) {
	root := filepath.Join(filepath.Dir(b), "../..")
	dockerComposeFile, err := filepath.Abs(root + "/docker-compose-it.yml")
	composeFilePaths := []string{dockerComposeFile}
	composeNameProject := "backend-takehome-test"

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, composeNameProject)
	execError := compose.WithCommand([]string{"up", "-d", "--build"}).
		WaitForService("backend-takehome-test", wait.ForHTTP("/health_check/db").WithPort("8080")).
		Invoke()
	err = execError.Error
	if err != nil {
		compose.Down()
		fmt.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
		return nil, err
	}

	return compose, nil
}
