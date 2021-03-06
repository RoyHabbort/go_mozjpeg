package main

import (
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

	argsWithoutProg := os.Args[1:]

	var path = "./"
	var quality uint = 75

	if len(argsWithoutProg) >= 1 {
		path = argsWithoutProg[0]
	}

	if len(argsWithoutProg) >= 2 {
		q, err := strconv.ParseUint(argsWithoutProg[1], 10, 0)

		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		} else {
			quality = uint(q)
		}
	}

	processDir(path, quality)
}

func processDir (dirPath string, quality uint) error {
	count := 0

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
			errorImage := optimizeImage(path, quality)

			if errorImage != nil {
				fmt.Printf("Image optimization error \n")
				panic(errorImage)
			} else {
				fmt.Printf("Complete optimize: %v\n",  info.Name())
			}
		}

		return nil
	})

	fmt.Println("Complete for " + strconv.Itoa(count) + " images")

	return nil
}

func optimizeImage(path string, quality uint) error {
	id := guuid.New()
	tmp := "/tmp/" + id.String()
	tmp = "./../tmp/result/" + id.String()

	err := mozjpegbin.NewCJpeg().
		Quality(quality).
		InputFile(path).
		OutputFile(tmp).
		Run()

	if  err != nil {
		return err
	}

	err = copy(tmp, path)
	if  err != nil {
		return err
	}

	err = os.Remove(tmp)
	if err != nil {
		return err
	}

	return nil
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