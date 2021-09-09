package integration_tests

import (
	"fmt"
	"gopkg.in/h2non/baloo.v3"
	"home24-page-analyser/cmd/app"
	"home24-page-analyser/cmd/config"
	"home24-page-analyser/cmd/router"
	"os"
	"testing"
	"time"
)

var client *baloo.Client

func TestMain(m *testing.M) {
	testConfig := &config.Config{Port: "9000"}

	client = baloo.New(fmt.Sprintf("http://localhost:%s", testConfig.Port))

	go app.NewApp(testConfig).StartServer(router.NewRouter(testConfig))

	time.Sleep(1 * time.Second)

	exitCode := m.Run()
	os.Exit(exitCode)
}
