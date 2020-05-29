package os

import (
	sys_os "os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	sys_os.Clearenv()
	sys_os.Setenv("TESTING", "value")
	if GetEnv("TESTING", "") != "value" {
		t.Fatalf("unexpected environment value")
	}
}

func TestGetEnv_default(t *testing.T) {
	sys_os.Clearenv()
	if GetEnv("TESTING", "default") != "default" {
		t.Fatalf("unexpected environment value")
	}
}

func TestGetIntEnv(t *testing.T) {
	sys_os.Clearenv()
	sys_os.Setenv("TESTING", "200")
	if GetIntEnv("TESTING", 0) != 200 {
		t.Fatalf("unexpected environment value")
	}
}

func TestGetIntEnv_default(t *testing.T) {
	sys_os.Clearenv()
	if GetIntEnv("TESTING", 100) != 100 {
		t.Fatalf("unexpected environment value")
	}
}
