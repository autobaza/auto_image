package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Pic struct {
	Name    string
	SaveDir string
}

func (p *Pic) getPath() string {
	return fmt.Sprintf("https://autobaza.kg/uploads/%s/%s/%s/%s", p.Name[0:2], p.Name[2:4], "1024x768", p.Name)
}

func (p *Pic) Save() error {
	resp, err := http.Get(p.getPath())
	defer resp.Body.Close()
	if err != nil {
		return errors.New(err.Error())
	}

	dirPath := fmt.Sprintf("%s/%s/%s", p.SaveDir, p.Name[0:2], p.Name[2:4])
	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return errors.New(err.Error())
	}

	out, err := os.Create(fmt.Sprintf("%s/%s", dirPath, p.Name))
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
