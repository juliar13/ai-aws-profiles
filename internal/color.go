package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ColorRule struct {
	Pattern string
	Color   string
}

type ColorManager struct {
	Rules []ColorRule
}

func NewColorManager() *ColorManager {
	return &ColorManager{}
}

func (cm *ColorManager) LoadColorSettings() error {
	if _, err := os.Stat("color-setting.ini"); os.IsNotExist(err) {
		return cm.createDefaultColorSettings()
	}

	file, err := os.Open("color-setting.ini")
	if err != nil {
		return fmt.Errorf("failed to open color-setting.ini: %w", err)
	}
	defer file.Close()

	cm.Rules = []ColorRule{}
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			pattern := parts[0]
			color := parts[1]
			cm.Rules = append(cm.Rules, ColorRule{
				Pattern: pattern,
				Color:   color,
			})
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read color-setting.ini: %w", err)
	}
	
	return nil
}

func (cm *ColorManager) createDefaultColorSettings() error {
	defaultContent := `admin 6644FF
readonly 22CCAA
`
	
	err := os.WriteFile("color-setting.ini", []byte(defaultContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create default color-setting.ini: %w", err)
	}
	
	cm.Rules = []ColorRule{
		{Pattern: "admin", Color: "6644FF"},
		{Pattern: "readonly", Color: "22CCAA"},
	}
	
	return nil
}

func (cm *ColorManager) GetColorForProfile(profileName string) string {
	defaultColor := "00aa00"
	matchedColor := defaultColor
	
	for _, rule := range cm.Rules {
		if strings.Contains(profileName, rule.Pattern) {
			matchedColor = rule.Color
		}
	}
	
	return matchedColor
}