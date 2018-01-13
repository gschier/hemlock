package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Unzip will un-compress a zip archive,
// moving all files and folders to an output directory
func Unzip(src, dest string, ignoreBaseDir bool) ([]string, error) {
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
