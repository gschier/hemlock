package cli

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/howeyc/fsnotify"
	"go/build"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var stopChannel chan bool
var isRunning = false

func init() {
	stopChannel = make(chan bool)

	cmd := Command("serve", "Run a Hemlock project")

	watch := cmd.Flag("watch", "Restart server when files change").Bool()

	cmd.Action(func(context *kingpin.ParseContext) error {
		if *watch {
			go buildAndRunApp()
			watchApp()
		} else {
			// Block until it finishes
			<-buildAndRunApp()
		}

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

func buildApp() {
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
}

func runApp() chan bool {
	doneChannel := make(chan bool)

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

	go func() {
		select {
		case <-stopChannel:
			cmd.Process.Kill()
		}
	}()

	isRunning = true

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)

	go func() {
		cmd.Wait()
		isRunning = false
		doneChannel <- true
	}()

	return doneChannel
}

func watchApp() {
	fmt.Printf("Watching for changes...\n")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		needsToBuild := false
		nextBuild := time.Now()
		for {
			select {
			case <-watcher.Event:
				needsToBuild = true
			case err := <-watcher.Error:
				fmt.Printf("File watch error: %v\n", err)
			default:
				if needsToBuild && time.Now().After(nextBuild) {
					buildAndRunApp()
					needsToBuild = false
					nextBuild = time.Now().Add(time.Second * 5)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	err = watcher.Watch(".")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func buildAndRunApp() chan bool {
	// Build before we kill the existing one
	buildApp()

	if isRunning {
		stopChannel <- true
	}

	return runApp()
}
