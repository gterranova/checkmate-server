package converter

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	TempDirPrefix string
}

func (p *File) WriteToTempFile(filename string, data []byte) (fname string, err error) {
	tempFileName, err := p.TempFile(filename)
	if err != nil {
		return
	}

	err = os.WriteFile(tempFileName, data, 0644)

	if err != nil {
		err = fmt.Errorf("write file %s failure", tempFileName)
		return
	}

	fname = tempFileName

	return
}

func (p *File) TempFile(f string) (filename string, err error) {
	tmpDir := os.TempDir()

	dir := filepath.Join(tmpDir, p.TempDirPrefix)

	err = os.MkdirAll(dir, 0755)

	if err != nil {
		return "", fmt.Errorf("make temp dir failure: %s, error: %s", dir, err)
	}

	return filepath.Join(dir, f), nil
}
