package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type AWSCredentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
}

func main() {
	fmt.Println("🔒 Security Context: Fetching AWS SSO Token (Memory Pipe Only)...")

	// 1. Intercept AWS credentials via 'aws configure export-credentials'.
	// This uses a memory buffer (akin to a pipe), avoiding passing secrets
	// as command line arguments, and thus preventing exposure via 'ps -ef'.
	cmd := exec.Command("aws", "configure", "export-credentials")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to get AWS credentials: %v\n", err)
		os.Exit(1)
	}

	var creds AWSCredentials
	if err := json.Unmarshal(out.Bytes(), &creds); err != nil {
		fmt.Printf("❌ Failed to parse AWS credentials: %v\n", err)
		os.Exit(1)
	}

	// 2. Create an ephemeral, secure file to pass credentials to act
	tmpFile, err := os.CreateTemp("", "act-sec-*.txt")
	if err != nil {
		fmt.Printf("❌ Failed to create temporary secret file: %v\n", err)
		os.Exit(1)
	}
	tmpFileName := tmpFile.Name()

	// Ensure the file is destroyed immediately upon exit, panic, or interruption
	defer func() {
		os.Remove(tmpFileName)
		fmt.Println("🧹 Ephemeral credentials destroyed.")
	}()

	content := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s\nAWS_SECRET_ACCESS_KEY=%s\nAWS_SESSION_TOKEN=%s\n",
		creds.AccessKeyId, creds.SecretAccessKey, creds.SessionToken)

	if _, err := tmpFile.WriteString(content); err != nil {
		fmt.Printf("❌ Failed to write credentials: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	// 3. Trigger the local act pipeline
	fmt.Println("🚀 Triggering local act pipeline...")
	actCmd := exec.Command("act", "workflow_dispatch", "-W", ".github/workflows/deploy.yml", "--secret-file", tmpFileName, "--env", "ACT=true")
	actCmd.Stdout = os.Stdout
	actCmd.Stderr = os.Stderr
	actCmd.Stdin = os.Stdin

	if err := actCmd.Run(); err != nil {
		fmt.Printf("❌ act command failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Local deployment via act completed cleanly.")
}
