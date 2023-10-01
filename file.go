package transcribe

import (
	"os"
	"path/filepath"
)

const DIRECTORY_WRITE_PERMISSIONS = 0700
const FILE_WRITE_PERMISSIONS = 0600

func OpenFileForWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), DIRECTORY_WRITE_PERMISSIONS); err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FILE_WRITE_PERMISSIONS)
	} else {
		return nil, err
	}
}
