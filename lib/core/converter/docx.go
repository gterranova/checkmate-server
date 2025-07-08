package converter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func addFile(zipWriter *zip.Writer, filename string, src io.Reader) {
	w1, err := zipWriter.Create(filename)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w1, src); err != nil {
		panic(err)
	}
}

func FixDocxTableStyle(inputFile string, outputFile string) error {
	if _, err := os.Stat(inputFile); err == nil {
		if b, err := os.ReadFile(inputFile); err == nil {
			s := bytes.NewReader(b)

			read, err := zip.NewReader(s, s.Size())

			if err != nil {
				return err
			}

			archive, err := os.Create(outputFile)
			if err != nil {
				return err
			}
			defer archive.Close()

			zipWriter := zip.NewWriter(archive)
			defer zipWriter.Close()

			// Iterate through the files in the zip archive
			for _, f := range read.File {
				// Open the current file
				v, err := f.Open()
				if err != nil {
					return err
				}
				defer v.Close()

				if f.Name == "word/document.xml" {
					// Read the contents of the file
					b, err := io.ReadAll(v)
					if err != nil {
						return err
					}

					addFile(zipWriter, f.Name, bytes.NewReader([]byte(strings.ReplaceAll(string(b),
						"<w:tblStyle w:val=\"Table\" />",
						"<w:tblStyle w:val=\"StileTable\" />"))))
				} else {
					addFile(zipWriter, f.Name, v)
				}
			}
			return nil
		}
	}

	return fmt.Errorf("cannot open %s", inputFile)
}
