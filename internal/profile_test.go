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
	// デフォルトのカラールールを設定
	generator.ColorManager.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}
	
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
	if profile.Color != "6644FF" {
		t.Errorf("Expected color '6644FF', got '%s'", profile.Color)
	}
}

func TestGenerateProfiles(t *testing.T) {
	generator := NewGenerator()
	// デフォルトのカラールールを設定
	generator.ColorManager.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}
	
	generator.GenerateProfiles("123456789012", "test")
	
	expectedProfiles := 2
	if len(generator.Profiles) != expectedProfiles {
		t.Errorf("Expected %d profiles, got %d", expectedProfiles, len(generator.Profiles))
	}
	
	expectedData := []struct {
		name  string
		color string
	}{
		{"test-admin", "6644FF"},
		{"test-readonly", "22CCAA"},
	}
	
	for i, expected := range expectedData {
		if generator.Profiles[i].Name != expected.name {
			t.Errorf("Expected profile name '%s', got '%s'", expected.name, generator.Profiles[i].Name)
		}
		if generator.Profiles[i].Color != expected.color {
			t.Errorf("Expected profile color '%s', got '%s'", expected.color, generator.Profiles[i].Color)
		}
	}
}

func TestGenerateExtensionFormat(t *testing.T) {
	generator := NewGenerator()
	// デフォルトのカラールールを設定
	generator.ColorManager.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}
	
	generator.AddProfile("test-admin", "arn:aws:iam::123456789012:role/AdminSwitchRole")
	
	output := generator.GenerateExtensionFormat()
	expected := `[profile test-admin]
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole
region = ap-northeast-1
color = 6644FF`
	
	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestGenerateConfigFormat(t *testing.T) {
	generator := NewGenerator()
	// デフォルトのカラールールを設定（configフォーマットでは色は使用されないが、一貫性のため）
	generator.ColorManager.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}
	
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

func TestSetRoleSessionName(t *testing.T) {
	generator := NewGenerator()
	
	if generator.RoleSessionName != "user_name" {
		t.Errorf("Expected default RoleSessionName 'user_name', got '%s'", generator.RoleSessionName)
	}
	
	generator.SetRoleSessionName("claude")
	
	if generator.RoleSessionName != "claude" {
		t.Errorf("Expected RoleSessionName 'claude', got '%s'", generator.RoleSessionName)
	}
}

func TestGenerateConfigFormatWithCustomRoleSessionName(t *testing.T) {
	generator := NewGenerator()
	generator.ColorManager.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
	}
	
	generator.SetRoleSessionName("claude")
	generator.AddProfile("test-admin", "arn:aws:iam::123456789012:role/AdminSwitchRole")
	
	output := generator.GenerateConfigFormat()
	
	if !strings.Contains(output, "role_session_name = claude") {
		t.Errorf("Expected output to contain 'role_session_name = claude', but it didn't. Output:\n%s", output)
	}
	
	if strings.Contains(output, "role_session_name = user_name") {
		t.Errorf("Expected output NOT to contain 'role_session_name = user_name', but it did. Output:\n%s", output)
	}
}

