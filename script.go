package main

import (
	"flag"
	"fmt"
	mimetype "github.com/gabriel-vasile/mimetype"
	guuid "github.com/google/uuid"
	mozjpegbin "github.com/nickalie/go-mozjpegbin"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	var path string
	var quality int

	flag.StringVar(&path, "path", "./", "a string var")
	flag.IntVar(&quality, "quality", 75, "a int var")

	flag.Parse()

	q := uint(quality)

	processDir(path, q)
}

func processDir (dirPath string, quality uint) error {
	count := 0

	var pathChanel chan string = make(chan string)
	go optimizeImageWithChanel(pathChanel, quality)

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		mime, err := mimetype.DetectFile(path)
		if err != nil {
			panic(err)
		}

		if mime.String() == "image/jpeg" {
			count++
			fmt.Println("Start optimize: " + path)
			pathChanel <- path
		}

		return nil
	})

	fmt.Println("Complete for " + strconv.Itoa(count) + " images")

	return nil
}

func optimizeImageWithChanel(pathChanel chan string, quality uint) {
	id := guuid.New()
	tmp := "/tmp/" + id.String()

	for {
		path := <- pathChanel

		err := mozjpegbin.NewCJpeg().
			Quality(quality).
			InputFile(path).
			OutputFile(tmp).
			Run()

		if  err != nil {
			panic(err)
		} else {
			err = copy(tmp, path)
			if  err != nil {
				panic(err)
			} else {
				err = os.Remove(tmp)
				if err != nil {
					panic(err)
				} else {
					fmt.Printf("Complete optimize: %v\n",  path)
				}
			}
		}
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}