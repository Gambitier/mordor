package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RSAKeyGenerator struct {
	outputPath string
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

func main() {
	var outputPath string
	flag.StringVar(&outputPath, "output", "", "Output directory path for the key files")
	flag.Parse()

	if outputPath == "" {
		fmt.Println("Error: output directory path is required")
		flag.Usage()
		os.Exit(1)
	}

	keygen := &RSAKeyGenerator{
		outputPath: outputPath,
	}

	if err := keygen.generateKeyPair(); err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		os.Exit(1)
	}
}
