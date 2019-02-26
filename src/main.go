package main

/**
 * dZip
 * Dependency-less Zip Tool
 *
 * @author Noah Halstead <nhalstead00@gmail.com>
 * @link https://golangcode.com/unzip-files-in-go/
 * @link https://stackoverflow.com/a/52394699/5779200
 * @link http://www.golangprograms.com/go-program-to-compress-list-of-files-into-zip.html
 */

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vjeantet/jodaTime"
)

const NoChange = 0
const ZipGood = 0
const FileInuse = 1
const ZipFailed = 9

var ZipFile string
var PackMode bool
var FileTails FileList

// FileList Used to Add Files to the Zip from a JSON List
type FileList []string

func main() {
	flag.StringVar(&ZipFile, "file", "", "Target Zip File to act on")
	flag.BoolVar(&PackMode, "zip", false, "Set to Unzip or Zip Mode, Pass the -zip flag to zip files on the input")
	flag.Parse()
	FileTails = flag.Args()

	if PackMode == true {
		// Write Mode, Pack Files into the Zip File

		flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		file, err := os.OpenFile(ZipFile, flags, 0644)
		if err != nil {
			log("Failed to open zip for writing.\n" + err.Error())
			os.Exit(FileInuse)
		}
		defer file.Close()

		zipw := zip.NewWriter(file)
		defer zipw.Close()

		for _, filename := range FileTails {
			if err := writeFiles(filename, zipw); err != nil {
				log("Failed to add file " + filename + " to zip: " + err.Error())
			}
		}

	} else {

		ZipFile, err := filepath.Abs(ZipFile)
		if err != nil {
			log("Failed to get the Realpath of: " + err.Error())
			os.Exit(FileInuse)
		}

		// Read Mode, Extract the Data from the Zip
		dir, file := filepath.Split(ZipFile)
		name := strings.TrimSuffix(file, filepath.Ext(file))
		finalTarget := dir + name
		files, err := unzip(ZipFile, finalTarget)

		if err != nil {
			log(fmt.Sprintf("Failed to unzip file from zip: %s \nError -> %s\nSource -> %s\nTarget -> %s", ZipFile, err, file, finalTarget))
			os.Exit(ZipFailed)
		} else {
			fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
			os.Exit(ZipGood)
		}

	}

}

/**
 * The Zip Lib for Golang does not Support appending to the Current Files within a ZipFile
 */
func writeFiles(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

/**
 * log
 *
 * Console Logs all of the Args provided to the Function.
 * The Prepends a timestamp, Good for Data Logging and
 *  when they happened.
 */
func log(in ...string) {
	dateTime := jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now())
	fmt.Print("[")
	fmt.Print(dateTime)
	fmt.Print("] ")
	fmt.Println(strings.Join(in, ""))
}
