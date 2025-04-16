package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type RSAKeyGenerator struct {
	outputPath string
}

func (k *RSAKeyGenerator) cleanup(keepCount int) error {
	// Read all files in the output directory
	files, err := os.ReadDir(k.outputPath)
	if err != nil {
		return fmt.Errorf("failed to read output directory: %v", err)
	}

	// Filter and group key files by timestamp
	type keyPair struct {
		timestamp string
		files     []string
	}
	keyPairs := make(map[string]*keyPair)

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".private.pem") && !strings.HasSuffix(name, ".public.pem") {
			continue
		}

		// Extract timestamp from filename (first 14 characters)
		if len(name) < 14 {
			continue
		}
		timestamp := name[:14]

		if _, exists := keyPairs[timestamp]; !exists {
			keyPairs[timestamp] = &keyPair{
				timestamp: timestamp,
				files:     make([]string, 0, 2),
			}
		}
		keyPairs[timestamp].files = append(keyPairs[timestamp].files, name)
	}

	// Convert map to slice for sorting
	pairs := make([]*keyPair, 0, len(keyPairs))
	for _, pair := range keyPairs {
		pairs = append(pairs, pair)
	}

	// Sort by timestamp in descending order (newest first)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].timestamp > pairs[j].timestamp
	})

	// Keep only the specified number of most recent pairs, delete the rest
	for i := keepCount; i < len(pairs); i++ {
		for _, filename := range pairs[i].files {
			fullPath := filepath.Join(k.outputPath, filename)
			if err := os.Remove(fullPath); err != nil {
				return fmt.Errorf("failed to remove old key file %s: %v", filename, err)
			}
		}
	}

	return nil
}

func (k *RSAKeyGenerator) generateKeyPair() error {
	timestamp := time.Now().UTC().Format("20060102150405")

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(k.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate private key
	privateKeyPath := filepath.Join(k.outputPath, fmt.Sprintf("%s.private.pem", timestamp))
	privateKeyCmd := exec.Command(
		"openssl", "genpkey",
		"-algorithm", "RSA",
		"-pkeyopt", "rsa_keygen_bits:2048",
		"-out", privateKeyPath,
	)

	if err := privateKeyCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Extract public key from private key
	publicKeyPath := filepath.Join(k.outputPath, fmt.Sprintf("%s.public.pem", timestamp))
	publicKeyCmd := exec.Command(
		"openssl", "rsa",
		"-pubout",
		"-in", privateKeyPath,
		"-out", publicKeyPath,
	)

	if err := publicKeyCmd.Run(); err != nil {
		return fmt.Errorf("failed to extract public key: %v", err)
	}

	fmt.Printf("Generated key pair:\n")
	fmt.Printf("Private key: %s\n", privateKeyPath)
	fmt.Printf("Public key: %s\n", publicKeyPath)

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "RSAKeyGenerator",
	Short: "A tool for generating and managing RSA key pairs",
	Long: `RSAKeyGenerator is a CLI tool that generates RSA key pairs and manages their lifecycle.
It can generate new key pairs and clean up old ones while maintaining a specified number of recent pairs.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new RSA key pair",
	Long:  `Generate a new RSA key pair and store it in the specified output directory. Automatically cleans up old pairs, keeping the 2 most recent.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath, _ := cmd.Flags().GetString("output")
		if outputPath == "" {
			return fmt.Errorf("output directory path is required")
		}

		keygen := &RSAKeyGenerator{
			outputPath: outputPath,
		}

		return keygen.generateKeyPair()
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up old key pairs",
	Long:  `Clean up old key pairs while keeping a specified number of most recent pairs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath, _ := cmd.Flags().GetString("output")
		if outputPath == "" {
			return fmt.Errorf("output directory path is required")
		}

		keepCount, _ := cmd.Flags().GetInt("keep")
		if keepCount < 1 {
			return fmt.Errorf("keep count must be at least 1")
		}

		keygen := &RSAKeyGenerator{
			outputPath: outputPath,
		}

		if err := keygen.cleanup(keepCount); err != nil {
			return err
		}

		fmt.Printf("Successfully cleaned up keys, keeping %d most recent pairs\n", keepCount)
		return nil
	},
}

func init() {
	// Add generate command
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().String("output", "", "Output directory path for the key files")
	generateCmd.MarkFlagRequired("output")

	// Add cleanup command
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().String("output", "", "Output directory path containing the key files")
	cleanupCmd.Flags().Int("keep", 2, "Number of most recent key pairs to keep")
	cleanupCmd.MarkFlagRequired("output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
