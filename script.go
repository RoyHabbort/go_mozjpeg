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
	processDir(".")
}

func processDir (dirPath string) error {
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
			errorImage := optimizeImage(path)

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

func optimizeImage(path string) error {
	id := guuid.New()
	tmp := "/tmp/" + id.String()
	tmp = "./../tmp/result/" + id.String()

	// Красиво? а вот хуй работает. jpeg файл видите ли не того формата. и не намёка на то, что не нравится
	err := mozjpegbin.NewCJpeg().
		Quality(75).
		InputFile(path).
		OutputFile(tmp).
		Run()

	if  err != nil {
		return err
	} else {
		fmt.Println(tmp)
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