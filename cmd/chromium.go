package cmd

import (
	"fmt"
	"github.com/aerokube/images/build"
	"github.com/spf13/cobra"
)

var (
	chromiumDebian       bool
	chromiumArchitecture string

	chromiumCmd = &cobra.Command{
		Use:   "chromium",
		Short: "build Chromium image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource: build.BrowserSource(browserSource),
				NoCache:       noCache,
				TestsDir:      testsDir,
				RunTests:      test,
				IgnoreTests:   ignoreTests,
				Tags:          tags,
				PushImage:     push,
			}
			if chromiumDebian {
				chromium := &build.DebianChromium{Requirements: req, Architecture: chromiumArchitecture}
				return chromium.Build()
			}
			if chromiumArchitecture != "" {
				return fmt.Errorf("--architecture can only be used with --debian")
			}
			chromium := &build.Chromium{req}
			return chromium.Build()
		},
	}
)

func init() {
	chromiumCmd.Flags().BoolVar(&chromiumDebian, "debian", false, "build using Debian Chromium packages")
	chromiumCmd.Flags().StringVar(&chromiumArchitecture, "architecture", "", "target architecture for Debian builds (amd64 or arm64)")
}
