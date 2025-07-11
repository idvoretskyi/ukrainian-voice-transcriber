package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/auth"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Cloud using OAuth",
	Long: `Simplified authentication setup guidance.

For the easiest setup, use gcloud CLI:
• Install gcloud: https://cloud.google.com/sdk/docs/install
• Run: gcloud auth login
• Run: gcloud auth application-default login

Alternative methods:
• Place service-account.json in current directory
• Set up custom OAuth client ID

Examples:
  ukrainian-voice-transcriber auth           # Show setup instructions
  ukrainian-voice-transcriber auth --status  # Check authentication status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oauth := auth.NewOAuthManager()
		
		status, _ := cmd.Flags().GetBool("status")
		
		if status {
			return showAuthStatus(oauth)
		}
		
		// Show authentication setup instructions
		ctx := context.Background()
		err := oauth.StartAuthFlow(ctx)
		if err != nil {
			return err
		}
		
		return nil
	},
}

// showAuthStatus shows current authentication status
func showAuthStatus(oauth *auth.OAuthManager) error {
	fmt.Printf("🔍 Checking authentication status...\n\n")
	
	// Check if gcloud is available and authenticated
	fmt.Printf("1. Checking gcloud CLI:\n")
	
	// Check if gcloud is installed
	if _, err := exec.LookPath("gcloud"); err != nil {
		fmt.Printf("   ❌ gcloud CLI not found\n")
		fmt.Printf("   Install: https://cloud.google.com/sdk/docs/install\n")
	} else {
		fmt.Printf("   ✅ gcloud CLI found\n")
		
		// Check if user is logged in
		cmd := exec.Command("gcloud", "auth", "list", "--filter=status:ACTIVE", "--format=value(account)")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			fmt.Printf("   ❌ No active gcloud authentication\n")
			fmt.Printf("   Run: gcloud auth login\n")
		} else {
			fmt.Printf("   ✅ gcloud user authenticated\n")
			
			// Check application default credentials
			cmd = exec.Command("gcloud", "auth", "application-default", "print-access-token")
			_, err = cmd.Output()
			if err != nil {
				fmt.Printf("   ❌ Application default credentials not set\n")
				fmt.Printf("   Run: gcloud auth application-default login\n")
			} else {
				fmt.Printf("   ✅ Application default credentials configured\n")
			}
		}
	}
	
	fmt.Printf("\n2. Checking service account:\n")
	if _, err := os.Stat("service-account.json"); err == nil {
		fmt.Printf("   ✅ service-account.json found\n")
	} else {
		fmt.Printf("   ❌ service-account.json not found\n")
	}
	
	fmt.Printf("\n💡 For setup instructions, run: ukrainian-voice-transcriber auth\n")
	
	return nil
}

func init() {
	authCmd.Flags().Bool("status", false, "Show current authentication status")
}