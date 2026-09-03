package main


import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	root := newRootCommand()

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "infer-blog",
		Short:         "Build and manage the Infer Origins blog",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(
		newBuildCommand(),
		newRebuildCommand(),
		newServeCommand(),
		newCleanCommand(),
		newCleanCacheCommand(),
		newCleanAllCommand(),
		newCheckCommand(),
		newVersionCommand(),
	)

	return cmd
}

func newBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build blog content and deployment artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliStage("Building blog")

			if err := buildBlog(); err != nil {
				return cliError("build failed", err)
			}

			cliSuccess("generated", filepath.Join(outputDir, "index.json"))
			cliSuccess("generated", filepath.Join(publicRootDir, "index.json"))

			if strings.TrimSpace(blogCNAME) != "" {
				cliSuccess("generated", filepath.Join(publicRootDir, "CNAME"))
			}

			fmt.Println()
			cliDone("Blog build complete")
			fmt.Printf("  Content  %s\n", sourceDir)
			fmt.Printf("  Output   %s\n", publicRootDir)
			fmt.Printf("  Blog     %s\n", publicBaseURL)
			fmt.Printf("  Website  %s\n", siteBaseURL)

			return nil
		},
	}
}

func newRebuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rebuild",
		Short: "Clean generated output and perform a fresh build",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliStage("Rebuilding blog")

			if err := os.RemoveAll(publicRootDir); err != nil {
				return cliError("clean generated output", err)
			}
			cliSuccess("cleaned", publicRootDir)

			if err := buildBlog(); err != nil {
				return cliError("build failed", err)
			}

			cliSuccess("generated", filepath.Join(outputDir, "index.json"))
			cliSuccess("generated", filepath.Join(publicRootDir, "index.json"))

			if strings.TrimSpace(blogCNAME) != "" {
				cliSuccess("generated", filepath.Join(publicRootDir, "CNAME"))
			}

			fmt.Println()
			cliDone("Blog rebuild complete")
			fmt.Printf("  Output   %s\n", publicRootDir)
			fmt.Printf("  Blog     %s\n", publicBaseURL)

			return nil
		},
	}
}

func newServeCommand() *cobra.Command {
	var port string
	var host string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Build and serve the generated blog locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := buildBlog(); err != nil {
				return cliError("build failed", err)
			}

			cliStage("Serving blog")
			fmt.Printf("  Local    http://%s:%s\n", host, port)
			fmt.Printf("  Index    http://%s:%s/index.json\n", host, port)
			fmt.Printf("  Blog API http://%s:%s/blog/index.json\n\n", host, port)

			server := exec.Command(
				"python3",
				"-m",
				"http.server",
				port,
				"--bind",
				host,
				"--directory",
				publicRootDir,
			)

			server.Stdout = os.Stdout
			server.Stderr = os.Stderr
			server.Stdin = os.Stdin

			if err := server.Run(); err != nil {
				return cliError("local server failed", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(
		&host,
		"host",
		envOrDefault("LOCAL_HOST", "localhost"),
		"local server host",
	)
	cmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		envOrDefault("LOCAL_PORT", "8080"),
		"local server port",
	)

	return cmd
}

func newCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove generated deployment output",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.RemoveAll(publicRootDir); err != nil {
				return cliError("clean failed", err)
			}

			cliDone("Removed " + publicRootDir)
			return nil
		},
	}
}

func newCleanCacheCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean-cache",
		Short: "Remove the image build cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			cacheRoot := envOrDefault("CACHE_DIR", ".cache")

			if err := os.RemoveAll(cacheRoot); err != nil {
				return cliError("clean cache failed", err)
			}

			cliDone("Removed " + cacheRoot)
			return nil
		},
	}
}

func newCleanAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean-all",
		Short: "Remove generated output and image cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.RemoveAll(publicRootDir); err != nil {
				return cliError("clean output failed", err)
			}

			cacheRoot := envOrDefault("CACHE_DIR", ".cache")
			if err := os.RemoveAll(cacheRoot); err != nil {
				return cliError("clean cache failed", err)
			}

			cliDone("Removed generated output and cache")
			return nil
		},
	}
}

func newCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check required development dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliStage("Checking dependencies")

			deps := []struct {
				name string
				bin  string
				args []string
			}{
				{"Go", "go", []string{"version"}},
				{"FFmpeg", "ffmpeg", []string{"-version"}},
				{"Git", "git", []string{"--version"}},
				{"Python", "python3", []string{"--version"}},
			}

			failed := false

			for _, dep := range deps {
				path, err := exec.LookPath(dep.bin)
				if err != nil {
					cliFailure(dep.name, "not found")
					failed = true
					continue
				}

				out, err := exec.Command(path, dep.args...).CombinedOutput()
				if err != nil {
					cliFailure(dep.name, err.Error())
					failed = true
					continue
				}

				line := strings.SplitN(
					strings.TrimSpace(string(out)),
					"\n",
					2,
				)[0]

				cliSuccess(dep.name, line)
			}

			if failed {
				return fmt.Errorf("one or more required dependencies are missing")
			}

			fmt.Println()
			cliDone("All dependencies available")
			return nil
		},
	}
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the infer-blog version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("infer-blog %s\n", version)
		},
	}
}

func envOrDefault(
	key string,
	fallback string,
) string {
	if value := strings.TrimSpace(
		os.Getenv(key),
	); value != "" {
		return value
	}

	return fallback
}

func cliStage(message string) {
	fmt.Printf("\n==> %s\n", message)
}

func cliSuccess(label string, value string) {
	fmt.Printf("  ✓ %-10s %s\n", label, value)
}

func cliDetail(label string, value string) {
	fmt.Printf("    %-10s %s\n", label, value)
}

func cliFailure(label string, value string) {
	fmt.Fprintf(
		os.Stderr,
		"  ✗ %-10s %s\n",
		label,
		value,
	)
}

func cliDone(message string) {
	fmt.Printf("✓ %s\n", message)
}

func cliError(
	message string,
	err error,
) error {
	cliFailure(message, err.Error())
	return err
}
