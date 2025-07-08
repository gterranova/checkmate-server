package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"terra9.it/checkmate/loader"
)

type FeatureDef struct {
	Lang      string   `json:"lang"`
	Filenames []string `json:"filenames"`
}

type TemplateDef struct {
	Name         string   `json:"name"`
	Filenames    []string `json:"filenames"`
	Format       string   `json:"format"`
	ReferenceDoc string   `json:"reference_doc"`
}

type Project struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	License string `json:"license"`

	StatusDefs   *TemplateDef   `json:"status"`
	TemplateDefs []*TemplateDef `json:"templates"`
	FeatureDefs  []*FeatureDef  `json:"features_defs"`
}

func loadFile(filename string) ([]byte, bool) {
	if _, err := os.Stat(filename); err == nil {
		if data, err := os.ReadFile(filename); err == nil {
			return data, true
		}
	}
	return nil, false
}

func dataDir() string {

	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../../data")
}

func mkPkg(name string, dataPath string) {
	fmt.Println(name)

	arcFile := name + loader.CHECKLIST_EXT
	arcPath := path.Join(dataPath, arcFile)

	pkgPath := path.Join(dataPath, name)
	projFile := path.Join(pkgPath, "config.json")
	logoFile := path.Join(pkgPath, "logo.png")

	fmt.Println("creating zip archive " + arcPath + "...")
	archive, err := os.Create(arcPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	p := new(Project)

	zipWriter := zip.NewWriter(archive)

	if content, ok := loadFile(projFile); ok {
		if err := json.Unmarshal(content, p); err != nil {
			panic(err)
		}
	}

	addFile(zipWriter, projFile, "config.json")
	addFile(zipWriter, logoFile, "logo.png")

	if p.StatusDefs != nil {
		for _, tmplFile := range p.StatusDefs.Filenames {
			addFile(zipWriter, path.Join(pkgPath, tmplFile), tmplFile)
		}
	}

	storedFiles := make(map[string]uint)
	for _, t := range p.TemplateDefs {
		for _, tmplFile := range t.Filenames {
			addFile(zipWriter, path.Join(pkgPath, tmplFile), tmplFile)
		}
		if len(t.ReferenceDoc) > 0 {
			if _, ok := storedFiles[t.ReferenceDoc]; !ok {
				addFile(zipWriter, path.Join(pkgPath, t.ReferenceDoc), t.ReferenceDoc)
				storedFiles[t.ReferenceDoc] = 1
			}
		}
	}

	for _, t := range p.FeatureDefs {
		for _, tmplFile := range t.Filenames {
			if _, ok := storedFiles[tmplFile]; !ok {
				addFile(zipWriter, path.Join(pkgPath, tmplFile), tmplFile)
				storedFiles[tmplFile] = 1
			}
		}
	}

	fmt.Println("closing zip archive...")
	zipWriter.Close()

}

func addFile(zipWriter *zip.Writer, projFile string, name string) {
	f1, err := os.Open(projFile)
	if err != nil {
		panic(err)
	}

	defer f1.Close()

	fmt.Println("Adding " + name + "...")
	w1, err := zipWriter.Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w1, f1); err != nil {
		panic(err)
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func main() {

	dataPath := dataDir()
	settingsJsonFile := path.Join(dataPath, "settings.json")
	settingsBundleFile := path.Join(dataPath, "settings"+loader.CHECKLIST_EXT)

	pkgs := make([]string, 0)

	entries, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			mkPkg(e.Name(), dataPath)
			pkgs = append(pkgs, e.Name())
		}
	}

	settings := `{"theme":"light","lang":"en","installed_projects":["`

	for i, pkg := range pkgs {

		if i != 0 {
			settings += `","`
		}
		settings += pkg
	}

	settings += `"]}`

	os.WriteFile(settingsJsonFile, []byte(settings), 0)

	archive, err := os.Create(settingsBundleFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	addFile(zipWriter, settingsJsonFile, "settings.json")
	addFile(zipWriter, path.Join(dataPath, "icon.png"), "icon.png")
	for _, pkg := range pkgs {
		addFile(zipWriter, path.Join(dataPath, pkg+loader.CHECKLIST_EXT), pkg+loader.CHECKLIST_EXT)
	}
	fmt.Println("closing zip archive...")
	zipWriter.Close()

	copy(settingsBundleFile, path.Join(dataPath, "../settings"+loader.CHECKLIST_EXT))

	fmt.Printf("deleting %s...\n", settingsJsonFile)
	os.Remove(settingsJsonFile)

	fmt.Printf("deleting %s...\n", settingsBundleFile)
	os.Remove(settingsBundleFile)
	for _, pkg := range pkgs {
		fmt.Printf("deleting %s...\n", path.Join(dataPath, pkg+loader.CHECKLIST_EXT))
		os.Remove(path.Join(dataPath, pkg+loader.CHECKLIST_EXT))
	}
}
