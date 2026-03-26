package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env file
	_ = godotenv.Load()

	expectedID := os.Getenv("EXPECTED_AWS_SSO_ACCOUNT_ID")
	if expectedID == "" {
		log.Fatal("❌ [CONFIG ERROR] EXPECTED_AWS_SSO_ACCOUNT_ID is not set in .env or environment variables.")
	}

	expectedRegion := os.Getenv("EXPECTED_AWS_SSO_REGION")

	// 2. Load AWS Configuration
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("❌ [AWS ERROR] Failed to load AWS config: %v\n(Check if your AWS CLI is configured correctly.)", err)
	}

	// 3. Verify Identity via STS
	svc := sts.NewFromConfig(cfg)
	identity, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("❌ [AUTH ERROR] Failed to verify AWS session: %v\n(Hint: Run 'aws sso login' or check your AWS_PROFILE.)", err)
	}

	currentID := *identity.Account
	fmt.Printf("🔍 Checking AWS Context... [Account: %s, Region: %s]\n", currentID, cfg.Region)

	// 4. Validate Account ID
	if currentID != expectedID {
		fmt.Printf("\n❌ [CRITICAL SECURITY ERROR] AWS Account Mismatch!\n")
		fmt.Printf("   - EXPECTED: %s\n", expectedID)
		fmt.Printf("   - CURRENT : %s\n", currentID)
		fmt.Println("\nActions required:")
		fmt.Println("1. Check your AWS_PROFILE or AWS_ACCESS_KEY_ID.")
		fmt.Println("2. Make sure you are logged into the correct SSO account.")
		fmt.Println("3. If this is intentional, update EXPECTED_AWS_SSO_ACCOUNT_ID in your .env file.")
		os.Exit(1)
	}

	// 5. Validate Region (Strict)
	if expectedRegion != "" && cfg.Region != expectedRegion {
		fmt.Printf("\n❌ [CRITICAL SECURITY ERROR] AWS Region Mismatch!\n")
		fmt.Printf("   - EXPECTED: %s\n", expectedRegion)
		fmt.Printf("   - CURRENT : %s\n", cfg.Region)
		fmt.Println("\nActions required:")
		fmt.Printf("1. Set environment variable: export AWS_REGION=%s\n", expectedRegion)
		fmt.Println("2. Or update your AWS_PROFILE default region.")
		os.Exit(1)
	}

	fmt.Println("✅ [SUCCESS] AWS Context verified (Account & Region). Proceeding safely.")
}
