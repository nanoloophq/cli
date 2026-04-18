package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nanoloop/cli/internal/api"
	"github.com/nanoloop/cli/internal/sourcemap"
	"github.com/spf13/cobra"
)

var (
	token     string
	appID     string
	dist      string
	release   string
	urlPrefix string
	dryRun    bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload source maps to nanoloop",
	RunE:  runUpload,
}

func init() {
	uploadCmd.Flags().StringVar(&token, "token", "", "API token (or set NANOLOOP_TOKEN)")
	uploadCmd.Flags().StringVar(&appID, "app", "", "App ID (or set NANOLOOP_APP_ID)")
	uploadCmd.Flags().StringVar(&dist, "dist", "./dist", "Directory containing source maps")
	uploadCmd.Flags().StringVar(&release, "release", "", "Release version (defaults to git commit)")
	uploadCmd.Flags().StringVar(&urlPrefix, "url-prefix", "", "URL prefix for source files")
	uploadCmd.Flags().BoolVar(&dryRun, "dry-run", false, "List files without uploading")
}

func runUpload(cmd *cobra.Command, args []string) error {
	if token == "" {
		token = os.Getenv("NANOLOOP_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("--token or NANOLOOP_TOKEN required")
	}

	if appID == "" {
		appID = os.Getenv("NANOLOOP_APP_ID")
	}
	if appID == "" {
		return fmt.Errorf("--app or NANOLOOP_APP_ID required")
	}

	if release == "" {
		release = detectRelease()
	}
	if release == "" {
		return fmt.Errorf("--release required (or run from a git repository)")
	}

	absPath, err := filepath.Abs(dist)
	if err != nil {
		return fmt.Errorf("invalid dist path: %w", err)
	}

	maps, err := sourcemap.Discover(absPath)
	if err != nil {
		return fmt.Errorf("failed to discover source maps: %w", err)
	}

	if len(maps) == 0 {
		fmt.Println("No source maps found in", absPath)
		return nil
	}

	fmt.Printf("Found %d source map(s) for release %s\n", len(maps), release)

	if dryRun {
		for _, m := range maps {
			fmt.Printf("  %s\n", m.Filename)
		}
		return nil
	}

	client := api.NewClient(token)
	result, err := client.UploadSourceMaps(appID, release, maps, urlPrefix)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("Uploaded %d source map(s)\n", len(result.Uploaded))
	for _, u := range result.Uploaded {
		fmt.Printf("  %s\n", u.Filename)
	}

	return nil
}

func detectRelease() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
