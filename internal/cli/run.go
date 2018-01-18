package cli

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func init() {
	cmd := Command("run", "")
	cmd.Action(func(context *kingpin.ParseContext) error {
		func() {
			fmt.Printf("Building %v...\n", filepath.Base(buildPath()))
			cmd := exec.Command("go", "build", "-o", buildPath(), ".")
			stderr, err := cmd.StderrPipe()
			if err != nil {
				panic(err)
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}

			err = cmd.Start()
			if err != nil {
				panic(err)
			}

			io.Copy(os.Stdout, stdout)
			io.Copy(os.Stderr, stderr)

			err = cmd.Wait()
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}

			if err != nil {
				panic(err)
			}
		}()

		func() {
			fmt.Printf("Running...\n")
			cmd := exec.Command(buildPath())

			stderr, err := cmd.StderrPipe()
			if err != nil {
				panic(err)
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}

			err = cmd.Start()
			if err != nil {
				panic(err)
			}

			go io.Copy(os.Stderr, stderr)
			go io.Copy(os.Stdout, stdout)

			err = cmd.Wait()
			if err != nil {
				panic(err)
			}
		}()

		return nil
	})
}

func buildPath() string {
	p := filepath.Join(tmpPath(), buildName())
	if runtime.GOOS == "windows" && filepath.Ext(p) != ".exe" {
		p += ".exe"
	}
	return p
}

func tmpPath() string {
	return os.TempDir()
}

func buildName() string {
	dir, _ := os.Getwd()
	return filepath.Base(dir)
}

func importPathFromCurrentDir() string {
	pwd, _ := os.Getwd()
	importPath, _ := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), pwd)
	return filepath.ToSlash(importPath)
}
