package internal

import (
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()
	if generator == nil {
		t.Fatal("NewGenerator() returned nil")
	}
	if len(generator.Profiles) != 0 {
		t.Errorf("Expected empty profiles, got %d", len(generator.Profiles))
	}
}

func TestAddProfile(t *testing.T) {
	generator := NewGenerator()
	generator.AddProfile("test-admin", "arn:aws:iam::123456789012:role/AdminSwitchRole")
	
	if len(generator.Profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(generator.Profiles))
	}
	
	profile := generator.Profiles[0]
	if profile.Name != "test-admin" {
		t.Errorf("Expected name 'test-admin', got '%s'", profile.Name)
	}
	if profile.RoleArn != "arn:aws:iam::123456789012:role/AdminSwitchRole" {
		t.Errorf("Expected role ARN 'arn:aws:iam::123456789012:role/AdminSwitchRole', got '%s'", profile.RoleArn)
	}
	if profile.Color != "00aa00" {
		t.Errorf("Expected color '00aa00', got '%s'", profile.Color)
	}
}

func TestGenerateProfiles(t *testing.T) {
	generator := NewGenerator()
	generator.GenerateProfiles("123456789012", "test")
	
	expectedProfiles := 2
	if len(generator.Profiles) != expectedProfiles {
		t.Errorf("Expected %d profiles, got %d", expectedProfiles, len(generator.Profiles))
	}
	
	expectedNames := []string{"test-admin", "test-readonly"}
	for i, expectedName := range expectedNames {
		if generator.Profiles[i].Name != expectedName {
			t.Errorf("Expected profile name '%s', got '%s'", expectedName, generator.Profiles[i].Name)
		}
	}
}

func TestGenerateExtensionFormat(t *testing.T) {
	generator := NewGenerator()
	generator.AddProfile("test-admin", "arn:aws:iam::123456789012:role/AdminSwitchRole")
	
	output := generator.GenerateExtensionFormat()
	expected := `[profile test-admin]
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole
region = ap-northeast-1
color = 00aa00`
	
	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestGenerateConfigFormat(t *testing.T) {
	generator := NewGenerator()
	generator.AddProfile("test-admin", "arn:aws:iam::123456789012:role/AdminSwitchRole")
	
	output := generator.GenerateConfigFormat()
	
	expectedParts := []string{
		"[default]",
		"region = ap-northeast-1",
		"output = json",
		"role_session_name = user_name",
		"[profile test-admin]",
		"source_profile = default",
		"role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole",
	}
	
	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected output to contain '%s', but it didn't. Output:\n%s", part, output)
		}
	}
}

