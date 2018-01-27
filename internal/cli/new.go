package cli

import (
	"archive/zip"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/tcnksm/go-input"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	cmd := Command("new", "Create a new Hemlock project")
	cmdDest := cmd.Arg("dest", "").String()

	cmd.Action(func(context *kingpin.ParseContext) error {
		dest := path.Clean(*cmdDest)
		q := fmt.Sprintf("Create new project in %s? [Y/n]", dest)
		answer, err := UI.Ask(q, &input.Options{
			Default:     "Yes",
			HideDefault: true,
			HideOrder:   true,
		})

		if err != nil {
			return err
		}

		if strings.ToLower(answer)[0] == 'n' {
			fmt.Printf("Cancelled.")
			return nil
		}

		fmt.Printf("[hemlock] Downloading sample project\n")
		resp, err := http.Get("https://github.com/gschier/hemlock-starter/archive/master.zip")
		if err != nil {
			return err
		}

		f, err := ioutil.TempFile("", "hemlock-")
		if err != nil {
			return err
		}

		err = func() error {
			fmt.Printf("[hemlock] Writing package\n")
			if err != nil {
				return err
			}
			defer f.Close()

			_, err := io.Copy(f, resp.Body)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}

		fmt.Printf("[hemlock] Extracting files\n")
		_, err = unzip(f.Name(), dest, true)
		if err != nil {
			return err
		}

		err = os.Remove(f.Name())
		if err != nil {
			return err
		}

		fmt.Printf("[hemlock] Application created\n")

		return nil
	})
}

// unzip will un-compress a zip archive,
// moving all files and folders to an output directory
func unzip(src, dest string, ignoreBaseDir bool) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		filename, err := handleZipFile(f, dest, ignoreBaseDir)
		if err != nil {
			return filenames, err
		}
		if filename == "" {
			continue
		}

		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func handleZipFile(f *zip.File, dest string, ignoreBaseDir bool) (string, error) {
	name := f.Name

	if ignoreBaseDir {
		i := strings.Index(name, "/")

		// Is it the base dir? Bail
		if i == -1 {
			return "", nil
		}

		name = name[i:]
	}

	// Store filename/path for returning and using later on
	filename := filepath.Join(dest, name)

	// Open the file from zip archive
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	// Make folder if it's a folder
	if f.FileInfo().IsDir() {
		err = os.MkdirAll(filename, os.ModePerm)
		return filename, err
	}

	// Ensure base directory exists
	if dir := path.Dir(filename); dir != "" {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Open destination for writing
	fDest, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return "", err
	}
	defer fDest.Close()

	// Write file into destination
	_, err = io.Copy(fDest, rc)
	if err != nil {
		return "", err
	}

	return filename, nil
}
