package internal

import (
	"os"
	"testing"
)

func TestColorManager_GetColorForProfile(t *testing.T) {
	cm := NewColorManager()
	cm.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}

	tests := []struct {
		profileName string
		expected    string
	}{
		{"test-admin", "6644FF"},
		{"test-readonly", "22CCAA"},
		{"test-admin-readonly", "22CCAA"}, // 後に記載されたルールが優先
		{"test-readonly-admin", "22CCAA"}, // 後に記載されたルールが優先
		{"production", "00aa00"},          // デフォルトカラー
		{"staging", "00aa00"},             // デフォルトカラー
	}

	for _, test := range tests {
		result := cm.GetColorForProfile(test.profileName)
		if result != test.expected {
			t.Errorf("GetColorForProfile(%q) = %q, expected %q", test.profileName, result, test.expected)
		}
	}
}

func TestColorManager_LoadColorSettings_Default(t *testing.T) {
	// color-setting.iniが存在しない場合のテスト
	testDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(testDir)

	cm := NewColorManager()
	err := cm.LoadColorSettings()
	if err != nil {
		t.Fatalf("LoadColorSettings() failed: %v", err)
	}

	// デフォルトルールが設定されているか確認
	if len(cm.Rules) != 2 {
		t.Errorf("Expected 2 default rules, got %d", len(cm.Rules))
	}

	expectedRules := []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}

	for i, expected := range expectedRules {
		if i >= len(cm.Rules) {
			t.Errorf("Rule %d missing", i)
			continue
		}
		if cm.Rules[i].Pattern != expected.Pattern || cm.Rules[i].Color != expected.Color {
			t.Errorf("Rule %d: got {%q, %q}, expected {%q, %q}",
				i, cm.Rules[i].Pattern, cm.Rules[i].Color, expected.Pattern, expected.Color)
		}
	}

	// ファイルが作成されているか確認
	if _, err := os.Stat("color-setting.ini"); os.IsNotExist(err) {
		t.Error("color-setting.ini was not created")
	}
}

func TestColorManager_LoadColorSettings_CustomFile(t *testing.T) {
	// カスタムcolor-setting.iniのテスト
	testDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(testDir)

	customContent := `prod FF0000
dev 00FF00
staging 0000FF
admin 6644FF
readonly 22CCAA
`
	err := os.WriteFile("color-setting.ini", []byte(customContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cm := NewColorManager()
	err = cm.LoadColorSettings()
	if err != nil {
		t.Fatalf("LoadColorSettings() failed: %v", err)
	}

	if len(cm.Rules) != 5 {
		t.Errorf("Expected 5 rules, got %d", len(cm.Rules))
	}

	// 優先順位のテスト（後に記載されたルールが優先）
	testCases := []struct {
		profileName string
		expected    string
	}{
		{"prod-admin", "6644FF"},     // prodとadminにマッチするが、adminが後なのでadmin(6644FF)が優先
		{"dev-readonly", "22CCAA"},   // devとreadonlyにマッチするが、readonlyが後なのでreadonly(22CCAA)が優先
		{"staging-only", "0000FF"},   // stagingのみにマッチ
		{"unknown", "00aa00"},        // どのルールにもマッチしない場合はデフォルト
	}

	for _, test := range testCases {
		result := cm.GetColorForProfile(test.profileName)
		if result != test.expected {
			t.Errorf("GetColorForProfile(%q) = %q, expected %q", test.profileName, result, test.expected)
		}
	}
}