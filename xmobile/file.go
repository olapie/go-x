package xmobile

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AppDirectories struct {
	Document  string
	Cache     string
	Temporary string
}

func (d *AppDirectories) MustMake() {
	d.Normalize()
	if d.Document != "" {
		MustMkdir(d.Document)
	} else {
		log.Println("Document directory is not specified")
	}

	if d.Cache != "" {
		MustMkdir(d.Cache)
	} else {
		log.Println("Cache directory is not specified")
	}

	if d.Temporary != "" {
		MustMkdir(d.Temporary)
	} else {
		log.Println("Temporary directory is not specified")
	}
}

func (d *AppDirectories) Normalize() {
	filePrefix := "file:"
	a := []*string{&d.Document, &d.Cache, &d.Temporary}
	for _, p := range a {
		if strings.HasPrefix(*p, filePrefix) {
			*p = (*p)[len(filePrefix):]
			*p = strings.Replace(*p, "///", "/", -1)
			*p = strings.Replace(*p, "//", "/", -1)
		}
	}
}

func NewAppDirectories() *AppDirectories {
	return new(AppDirectories)
}

func NewTestAppDirectories() *AppDirectories {
	return &AppDirectories{
		Document:  "testdata/document",
		Cache:     "testdata/cache",
		Temporary: "testdata/temporary",
	}
}

func GetDiskSize(path string) int64 {
	var sum int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			sum += info.Size()
		}
		return err
	})
	if err != nil {
		log.Println(err)
		return 0
	}
	return sum
}

func MustMkdir(dir string) {
	if f, err := os.Open(dir); err == nil {
		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}

		if fi.IsDir() {
			return
		}
		panic(dir + " is not directory")
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Make dir: %s, %v\n", dir, err)
	}
}
