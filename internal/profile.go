package internal

import (
	"fmt"
	"strings"
)

type Profile struct {
	Name    string
	RoleArn string
	Color   string
}

type Generator struct {
	Profiles     []Profile
	ColorManager *ColorManager
}

func NewGenerator() *Generator {
	cm := NewColorManager()
	return &Generator{
		ColorManager: cm,
	}
}

func (g *Generator) AddProfile(name, roleArn string) {
	color := g.ColorManager.GetColorForProfile(name)
	profile := Profile{
		Name:    name,
		RoleArn: roleArn,
		Color:   color,
	}
	g.Profiles = append(g.Profiles, profile)
}

func (g *Generator) GenerateProfiles(accountID, envType string) {
	roles := []string{"AdminSwitchRole", "ReadOnlySwitchRole"}
	
	for _, role := range roles {
		roleSuffix := strings.ToLower(strings.Replace(role, "SwitchRole", "", 1))
		profileName := fmt.Sprintf("%s-%s", envType, roleSuffix)
		roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, role)
		
		g.AddProfile(profileName, roleArn)
	}
}

func (g *Generator) GenerateExtensionFormat() string {
	var output strings.Builder
	
	for _, profile := range g.Profiles {
		output.WriteString(fmt.Sprintf("[profile %s]\n", profile.Name))
		output.WriteString(fmt.Sprintf("role_arn = %s\n", profile.RoleArn))
		output.WriteString("region = ap-northeast-1\n")
		output.WriteString(fmt.Sprintf("color = %s\n\n", profile.Color))
	}
	
	return strings.TrimSpace(output.String())
}

func (g *Generator) GenerateConfigFormat() string {
	var output strings.Builder
	
	output.WriteString("[default]\n")
	output.WriteString("region = ap-northeast-1\n")
	output.WriteString("output = json\n")
	output.WriteString("role_session_name = user_name\n\n")
	
	for _, profile := range g.Profiles {
		output.WriteString(fmt.Sprintf("[profile %s]\n", profile.Name))
		output.WriteString("source_profile = default\n")
		output.WriteString(fmt.Sprintf("role_arn = %s\n\n", profile.RoleArn))
	}
	
	return strings.TrimSpace(output.String())
}

