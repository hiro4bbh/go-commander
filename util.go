package gocommander

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hiro4bbh/go-log"
)

// BoxString represents the states BoxString(str) or BoxString(none).
type BoxString struct {
	str string
	ok  bool
}

// NewBoxString returns a new BoxString.
func NewBoxString(o interface{}) BoxString {
	str, ok := o.(string)
	return BoxString{str, ok}
}

// String returns the string representation.
func (box BoxString) String() string {
	if box.ok {
		return fmt.Sprintf("BoxString(%q)", box.str)
	}
	return "BoxString(none)"
}

// Unwrap returns the string if ok, otherwise returns an error.
func (box BoxString) Unwrap() (string, error) {
	if box.ok {
		return box.str, nil
	}
	return "", fmt.Errorf("tried to unwrap %s", box)
}

// UnwrapOr returns the string if ok, otherwise returns defval.
func (box BoxString) UnwrapOr(defval string) string {
	if box.ok {
		return box.str
	}
	return defval
}

// UnwrapFilePath returns the FilePath string if ok, otherwise returns an error.
func (box BoxString) UnwrapFilePath() (FilePath, error) {
	if box.ok {
		return FilePath(box.str), nil
	}
	return FilePath(""), fmt.Errorf("tried to unwrap %s", box)
}

// UnwrapFilePathOr returns the FilePath string if ok, otherwise returns defval.
func (box BoxString) UnwrapFilePathOr(defval FilePath) FilePath {
	if box.ok {
		return FilePath(box.str)
	}
	return defval
}

// FilePath is the string with filepath methods.
type FilePath string

// Base is filepath.Base(p).
func (p FilePath) Base() FilePath {
	return FilePath(filepath.Base(string(p)))
}

// Dir is filepath.Dir(p).
func (p FilePath) Dir() FilePath {
	return FilePath(filepath.Dir(string(p)))
}

// Ext is filepath.Ext(p).
func (p FilePath) Ext() string {
	return filepath.Ext(string(p))
}

// Join returns the joined FilePath.
func (p FilePath) Join(q FilePath) FilePath {
	return FilePath(filepath.Join(string(p), string(q)))
}

// CreateFile creates the parent directories if needed, creates a new file, and returns the *os.File.
//
// This function returns an error in file operations.
func CreateFile(name FilePath) (*os.File, error) {
	if err := os.MkdirAll(string(name.Dir()), 0750); err != nil {
		return nil, err
	}
	return os.Create(string(name))
}

// DownloadFile returns a downloaded file.
// If the file does not exist or cache is false, then this function downloads the file.
//
// This function returns an error in downloading.
func DownloadFile(dirpath FilePath, rawurl string, cache bool, logger *golog.Logger) (*os.File, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	logger.Debugf("downloading %s to %s ...", rawurl, dirpath)
	filename := dirpath.Join(FilePath(u.Path).Base())
	if cache {
		if file, err := os.Open(string(filename)); err == nil {
			logger.Debugf("used the downloaded file at %s for %s ...", filename, rawurl)
			return file, nil
		}
	}
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	file, err := CreateFile(filename)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	if _, err := io.Copy(writer, resp.Body); err != nil {
		file.Close()
		return nil, err
	}
	writer.Flush()
	file.Seek(0, 0)
	return file, nil
}

// DecompressFile returns a decompressed file.
// This function automatically closes the given file if returning the decompressed file.
// If the file does not exist or cache is false, then this function decompress the file.
// Currently, this function supports gzip (extension ".gz").
//
// This function returns an error in decompression.
func DecompressFile(file *os.File, cache bool, logger *golog.Logger) (*os.File, error) {
	filename := FilePath(file.Name())
	switch filename.Ext() {
	case ".gz":
		decompFilename := filename[:len(filename)-len(".gz")]
		if decompFile, err := os.Open(string(decompFilename)); err == nil {
			file.Close()
			logger.Debugf("used the decompressed file at %s", filename)
			return decompFile, err
		}
		logger.Debugf("decompressing %s with gzip ...", filename)
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, err
		}
		decompFile, err := CreateFile(decompFilename)
		if err != nil {
			return nil, err
		}
		writer := bufio.NewWriter(decompFile)
		if _, err := io.Copy(writer, reader); err != nil {
			return nil, err
		}
		if err := reader.Close(); err != nil {
			return nil, err
		}
		writer.Flush()
		file.Close()
		decompFile.Seek(0, 0)
		return decompFile, nil
	}
	return file, nil
}

// DownloadAndDecompressFile is DownloadFile followed by DecompressFile.
//
// This function returns an error by DownloadFile or DecompressFile.
func DownloadAndDecompressFile(dirpath FilePath, rawurl string, cache bool, logger *golog.Logger) (*os.File, error) {
	file, err := DownloadFile(dirpath, rawurl, cache, logger)
	if err != nil {
		return nil, err
	}
	decompFile, err := DecompressFile(file, cache, logger)
	if err != nil {
		file.Close()
		return nil, err
	}
	return decompFile, err
}

// Env returns the environment variable with the given name.
func Env(name string) BoxString {
	if val := os.Getenv(name); val != "" {
		return NewBoxString(val)
	}
	return NewBoxString(nil)
}

// HomeDir returns the user's home directory path.
//
// This function calls panic in getting the path.
func HomeDir() FilePath {
	u, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("gocommander.HomeDir: user.Current: %s", err))
	}
	return FilePath(u.HomeDir)
}
