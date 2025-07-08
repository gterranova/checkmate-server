package loader

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ResourceLoader struct {
	ResourceName string
	BasePath     string
	Data         map[string][]byte
}

func (r *ResourceLoader) Name() string {
	return r.ResourceName
}

func (r *ResourceLoader) Load(filename string) ([]byte, bool) {
	if _, err := os.Stat(filename); err == nil {
		if data, err := os.ReadFile(filename); err == nil {
			return data, true
		}
	}
	return nil, false
}

func (r *ResourceLoader) Get(filename string) ([]byte, bool) {
	path := path.Join(r.BasePath, r.ResourceName, filename)
	if content, ok := r.Data[path]; ok {
		return content, ok
	}

	if content, ok := r.Load(path); ok {
		r.Data[path] = content
		return content, true
	}
	return nil, false
}

func addFile(zipWriter *zip.Writer, filename string, src io.Reader) error {
	w1, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w1, src); err != nil {
		return err
	}
	return nil
}

func (r *ResourceLoader) Set(filename string, data []byte) error {
	path := path.Join(r.BasePath, r.ResourceName, filename)
	r.Data[filename] = data
	os.WriteFile(path, data, 0644)
	return nil
}

func (r *ResourceLoader) Save() error {

	outputFile := r.Name() + CHECKLIST_EXT
	return r.SaveAs(outputFile)
}

func (r *ResourceLoader) SaveAs(outputFile string) error {

	path := path.Join(r.BasePath, r.ResourceName) + "/"
	if name, found := strings.CutSuffix(outputFile, CHECKLIST_EXT); found {
		r.ResourceName = name
	} else {
		r.ResourceName = outputFile
		outputFile += CHECKLIST_EXT
	}

	archive, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Iterate through the files in the zip archive
	for k, v := range r.Data {
		zippedFilename, _ := strings.CutPrefix(k, path)
		//fmt.Println(k, zippedFilename, path)
		if err := addFile(zipWriter, zippedFilename, bytes.NewReader(v)); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceLoader) MustGet(filename string) []byte {
	if content, ok := r.Get(filename); ok {
		return content
	}

	panic(fmt.Errorf("file %s not found", filename))
}

func (r *ResourceLoader) Unpack(b *[]byte) error {

	s := bytes.NewReader(*b)

	read, err := zip.NewReader(s, s.Size())

	if err != nil {
		return err
	}

	// Iterate through the files in the zip archive
	for _, f := range read.File {
		// Open the current file
		v, err := f.Open()
		if err != nil {
			return err
		}
		defer v.Close()

		// Read the contents of the file
		b, err := io.ReadAll(v)
		if err != nil {
			return err
		}
		path := path.Join(r.BasePath, r.ResourceName, f.Name)
		r.Data[path] = b
		// Print the file name and contents
		//fmt.Printf("File Name: %s\n", path)
		//fmt.Printf("%s\n", string(b))
	}
	return nil
}

func (r *ResourceLoader) FromZip(filename string) (*ResourceLoader, error) {
	if b, ok := r.Load(filename); ok {
		return r.FromBuffer(b)
	}
	return r, fmt.Errorf("cannot open %s", filename)
}

func (r *ResourceLoader) FromBuffer(buf []byte) (*ResourceLoader, error) {
	if err := r.Unpack(&buf); err != nil {
		return r, err
	}
	return r, nil
}

func NewEmptyLoader(path string) *ResourceLoader {
	basePath := filepath.Dir(path)
	name := filepath.Base(path)
	r := ResourceLoader{ResourceName: name, BasePath: basePath, Data: make(map[string][]byte)}
	return &r
}

func NewLoader(pkg string) (*ResourceLoader, error) {
	var err error
	resLoader := NewEmptyLoader(pkg)

	if resLoader, err = resLoader.FromZip(pkg + CHECKLIST_EXT); err != nil {
		if _, ok := resLoader.Get("config.json"); !ok {
			return nil, fmt.Errorf("cant find config.json in %s", pkg)
		}
	}
	return resLoader, nil
}
