package shared_test

import (
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared"
	"testing"
)

func TestReturnsTrueWhenAllFieldsAreFilled(t *testing.T) {
	src := shared.Source{
		Target: "Target",
		SonarToken: "Token",
		Metrics: "Metrics",
		Component: "Component",
	}
	if !src.Valid(){
		t.Error("Wasn't valid")
	}
}

func TestReturnsFalseWhenTargetIsMissing(t *testing.T) {
	src := shared.Source{
		Target: "",
		SonarToken: "Token",
		Metrics: "Metrics",
		Component: "Component",
	}
	if src.Valid(){
		t.Error("Is still valid")
	}
}

func TestReturnsFalseWhenSonarTokenIsMissing(t *testing.T) {
	src := shared.Source{
		Target: "Target",
		SonarToken: "",
		Metrics: "Metrics",
		Component: "Component",
	}
	if src.Valid(){
		t.Error("Is still valid")
	}
}

func TestReturnsFalseWhenMetricsIsMissing(t *testing.T) {
	src := shared.Source{
		Target: "Target",
		SonarToken: "Token",
		Metrics: "",
		Component: "Component",
	}
	if src.Valid(){
		t.Error("Is still valid")
	}
}

func TestReturnsFalseWhenComponentIsMissing(t *testing.T) {
	src := shared.Source{
		Target: "Target",
		SonarToken: "Token",
		Metrics: "Metrics",
		Component: "",
	}
	if src.Valid(){
		t.Error("Is still valid")
	}
}
