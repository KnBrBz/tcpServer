package setup

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestSetupHost(t *testing.T) {
	// default value
	stp := New()
	if host := stp.Host(); host != envHostDefault {
		t.Fatalf("host expected `%s`, got `%s`", envHostDefault, host)
	}
	// custom file
	if err := godotenv.Load(".envExm"); err != nil {
		t.Fatal(err)
	}
	stp = New()
	if host := stp.Host(); host == envHostDefault || len(host) == 0 {
		t.Fatalf("host expected to be not default and not empty, got `%s`", host)
	}
}
