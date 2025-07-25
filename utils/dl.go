package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\r\t[...] Downloading - %s complete", humanize.Bytes(wc.Total))
}

func DownloadTarget(url, filePath string, forMacro bool) (bool, string, error) {
	fmt.Println("- Download Started @", url)

	err := DownloadFile(url, filePath, forMacro)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("- Download Finished:", filePath)

	return true, filePath, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(url, filePath string, forMacro bool) error {
	fileDir := "custom/"

	if !forMacro {
		fileDir = "custom/files/"
	}

	if _, err := os.Stat(fileDir); os.IsNotExist(err) { // Folder doesn't exist
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			log.Fatalln("custom folder couldn't be created:", err)
		}
	}

	// Folder exists
	if _, err := os.Stat(fileDir); !os.IsNotExist(err) {
		// Create the file, but give it a tmp file extension, this means we won't overwrite a
		// file until it's downloaded, but we'll remove the tmp extension once downloaded.
		out, err := os.Create(fileDir + filePath + ".tmp")
		if err != nil {
			return err
		}

		// Get the data
		resp, err := http.Get(url)
		if err != nil {
			out.Close()
			return err
		}
		defer resp.Body.Close()

		// Create our progress reporter and pass it to be used alongside our writer
		counter := &WriteCounter{}
		if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
			out.Close()
			return err
		}

		// The progress use the same line so print a new line once it's finished downloading
		fmt.Print("\n")

		// Close the file without defer so it can happen before Rename()
		out.Close()

		if err = os.Rename(fileDir+filePath+".tmp", fileDir+filePath); err != nil {
			return err
		}
	}

	return nil
}
