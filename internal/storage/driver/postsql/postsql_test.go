package postsql_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
)

func brokerStorageDatabaseTestConfig() storage.Config {
	return storage.Config{
		Host:            "localhost",
		User:            "test",
		Password:        "test",
		Port:            "5430",
		Name:            "testbroker",
		SSLMode:         "disable",
		SecretKey:       "################################",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Minute,
	}
}

func TestMain(m *testing.M) {
	exitVal := 0
	defer func() {
		os.Exit(exitVal)
	}()

	config := brokerStorageDatabaseTestConfig()

	docker, err := internal.NewDockerHandler()
	if err != nil {
		log.Fatal(err)
	}
	defer func(docker *internal.DockerHelper) {
		err := docker.CloseDockerClient()
		if err != nil {
			log.Fatal(err)
		}
	}(docker)

	cleanupContainer, err := docker.CreateDBContainer(internal.ContainerCreateRequest{
		Port:          config.Port,
		User:          config.User,
		Password:      config.Password,
		Name:          config.Name,
		Host:          config.Host,
		ContainerName: "keb-storage-tests",
		Image:         internal.PostgresImage,
	})
	defer func() {
		if cleanupContainer != nil {
			err := cleanupContainer()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	exitVal = m.Run()
}

func GetStorageForDatabaseTests() (func() error, storage.BrokerStorage, error) {
	return storage.GetStorageForTest(brokerStorageDatabaseTestConfig())
}
