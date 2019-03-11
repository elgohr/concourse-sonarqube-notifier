package shared_test

import (
	"encoding/base64"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared"
	"os"
	"testing"
)

var sonarqube shared.Sonarqube

func TestMain(m *testing.M) {
	sonarqube = shared.Sonarqube{}
	os.Exit(m.Run())
}

func getBasicHeader(authToken string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken+":"))
}