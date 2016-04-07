package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	imageEigenvalue = map[string][]byte{}
)

func init() {
	imageEigenvalue = map[string][]byte{
		".jpeg": []byte{0xFF, 0xD8},
		".png":  []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		".gif":  []byte{0x47, 0x49, 0x46},
		".bmp":  []byte{0x42, 0x4D},
	}

}

func main() {
	app := cli.NewApp()

	app.Action = action

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "src",
			Value: "",
			Usage: "images source dir",
		},
		cli.StringFlag{
			Name:  "dst",
			Value: "",
			Usage: "base64-images js generate to",
		},
		cli.StringFlag{
			Name:  "tpl",
			Value: "module.exports = { uri:'{{.data}}' }",
			Usage: "base64-image file template",
		},
	}

	app.Run(os.Args)
}

func action(c *cli.Context) {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("gen-base64-images: %s\n", err.Error())
			return
		}
	}()

	src := c.String("src")

	if src == "" {
		err = fmt.Errorf("%s", "arg of src could not be empty")
		return
	}

	if !filepath.IsAbs(src) {
		if src, err = filepath.Abs(src); err != nil {
			return
		}
	}

	dst := c.String("dst")
	if dst == "" {
		err = fmt.Errorf("%s", "arg of dst could not be empty")
		return
	}

	if !filepath.IsAbs(dst) {
		if dst, err = filepath.Abs(dst); err != nil {
			return
		}
	}

	tpl := c.String("tpl")

	if err = genreate(src, dst, tpl); err != nil {
		return
	}

}

func genreate(src, dst, tpl string) (err error) {

	createMap := map[string]string{}

	walkFn := func(path string, info os.FileInfo, inErr error) (err error) {

		if inErr != nil {
			return
		}

		if path == "" {
			return
		}

		if info.IsDir() {
			return
		}

		if strings.HasPrefix(info.Name(), ".") {
			return
		}

		var relPath string
		if relPath, err = filepath.Rel(src, path); err != nil {
			return
		}

		ext := filepath.Ext(path)
		var eigenvalue []byte
		var exist bool
		if eigenvalue, exist = imageEigenvalue[ext]; !exist {
			return
		}

		var fileData []byte
		if fileData, err = ioutil.ReadFile(path); err != nil {
			return
		}

		if len(fileData) < len(eigenvalue) {
			err = fmt.Errorf("bad file:", path)
			return
		}

		if !bytes.Equal(fileData[0:len(eigenvalue)], eigenvalue) {
			err = fmt.Errorf("bad eigenvalue:", path)
			return
		}

		b64 := base64.StdEncoding.EncodeToString(fileData)

		data := fmt.Sprintf("data:image/%s;base64,%s", strings.TrimLeft(ext, "."), b64)

		createMap[relPath] = data

		return
	}

	if err = filepath.Walk(src, walkFn); err != nil {
		return
	}

	var tmpl *template.Template
	if tmpl, err = template.New("js-images").Parse(tpl); err != nil {
		return
	}

	for relPath, data := range createMap {
		jsFilePath := filepath.Join(dst, relPath)

		if err = os.MkdirAll(filepath.Dir(jsFilePath), os.FileMode(0755)); err != nil {
			return
		}

		buf := bytes.NewBuffer(nil)

		if err = tmpl.Execute(buf, map[string]string{"data": data}); err != nil {
			return
		}

		if err = ioutil.WriteFile(jsFilePath+".js", buf.Bytes(), 0644); err != nil {
			return
		}
	}

	return
}
