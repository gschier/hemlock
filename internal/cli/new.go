package cli

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/gschier/hemlock/internal/zip"
	"github.com/tcnksm/go-input"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func init() {
	cmd := Command("new", "")
	cmdDest := cmd.Arg("dest", "").String()

	cmd.Action(func(context *kingpin.ParseContext) error {
		dest := path.Clean(*cmdDest)
		q := fmt.Sprintf("Create new project in %s? [Y/n]", dest)
		answer, err := UI.Ask(q, &input.Options{
			Default: "Yes",
			HideDefault: true,
			HideOrder: true,
		})

		if err != nil {
			return err
		}

		if strings.ToLower(answer)[0] == 'n' {
			fmt.Printf("Cancelled.")
			return nil
		}

		fmt.Printf("Downloading package\n")
		resp, err := http.Get("https://github.com/gschier/hemlock-starter/archive/master.zip")
		if err != nil {
			return err
		}

		f, err := ioutil.TempFile("", "hemlock-")
		if err != nil {
			return err
		}

		err = func() error {
			fmt.Printf("Writing zip to %v\n", f.Name())
			if err != nil {
				return err
			}
			defer f.Close()

			n, err := io.Copy(f, resp.Body)
			if err != nil {
				return err
			}
			fmt.Printf("Wrote %d bytes\n", n)
			return nil
		}()
		if err != nil {
			return err
		}

		fmt.Printf("Extracting files\n")
		_, err = zip.Unzip(f.Name(), dest, true)
		if err != nil {
			return err
		}

		fmt.Printf("Deleting zip\n")
		err = os.Remove(f.Name())
		if err != nil {
			return err
		}

		return nil
	})
}
