package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"ai-aws-profiles/internal"
	"github.com/spf13/cobra"
)

var (
	format         string
	output         string
	roleSessionName string
	version        = "1.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "aws-prof",
	Short: "Generate AWS profiles for Extend Switch Roles and ~/.aws/config",
	Long: `A CLI tool that generates AWS profile configurations for:
- AWS Extend Switch Roles extension
- ~/.aws/config file

The tool fetches account information from AWS Organizations automatically.`,
	RunE: runAWSProf,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.Flags().StringVar(&format, "format", "", "Output format (extension/config)")
	rootCmd.Flags().StringVar(&output, "output", "", "Output file path")
	rootCmd.Flags().StringVar(&roleSessionName, "role-session-name", "user_name", "Role session name for AWS config")
	rootCmd.AddCommand(versionCmd)
}

func runAWSProf(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	generator, err := internal.GenerateFromAWSAccounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch accounts from AWS: %w", err)
	}

	generator.SetRoleSessionName(roleSessionName)

	if format == "" {
		return writeDefaultOutput(generator)
	}

	content, err := generateContent(generator, format)
	if err != nil {
		return err
	}

	if output != "" {
		return writeToFile(content, output)
	}

	fmt.Print(content)
	return nil
}


func generateContent(generator *internal.Generator, format string) (string, error) {
	switch format {
	case "extension":
		return generator.GenerateExtensionFormat(), nil
	case "config":
		return generator.GenerateConfigFormat(), nil
	default:
		return "", fmt.Errorf("invalid format: %s (supported: extension, config)", format)
	}
}

func writeDefaultOutput(generator *internal.Generator) error {
	extensionContent := generator.GenerateExtensionFormat()
	configContent := generator.GenerateConfigFormat()

	if err := writeToFile(extensionContent, "extension.ini"); err != nil {
		return fmt.Errorf("failed to write extension.ini: %w", err)
	}

	if err := writeToFile(configContent, "config.ini"); err != nil {
		return fmt.Errorf("failed to write config.ini: %w", err)
	}

	fmt.Println("Generated extension.ini and config.ini")
	return nil
}

func writeToFile(content, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	return err
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}