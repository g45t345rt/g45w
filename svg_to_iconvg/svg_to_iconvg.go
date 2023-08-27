package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"github.com/gio-eui/ivgconv"
)

var pkgName = "app_icons"

func main() {
	fileNameFlag := flag.String("fileName", "", fmt.Sprintf("fileName in %s folder", pkgName))
	nameFlag := flag.String("name", "", fmt.Sprintf("name in %s pkg", pkgName))
	flag.Parse()

	fileName := *fileNameFlag
	name := *nameFlag

	if fileName == "" {
		log.Fatal("missing -fileName arg")
	}

	if name == "" {
		log.Fatal("missing -name arg")
	}

	svgPath := filepath.Join(pkgName, fmt.Sprintf("%s.svg", fileName))
	data, err := ivgconv.FromFile(svgPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = store(fileName, name, data)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func store(fileName string, name string, data []byte) error {
	out := new(bytes.Buffer)

	fmt.Fprintf(out, "package %s\n\n", pkgName)
	fmt.Fprintf(out, "var %s = []byte{", name)
	for i, b := range data {
		if i%16 == 0 {
			fmt.Fprintf(out, "\n\t")
		}
		fmt.Fprintf(out, "0x%02x, ", b)
	}
	fmt.Fprintf(out, "}\n")

	src, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}

	goPath := filepath.Join(pkgName, fmt.Sprintf("%s.go", fileName))
	if err := os.WriteFile(goPath, src, os.ModePerm); err != nil {
		return err
	}

	return nil
}
