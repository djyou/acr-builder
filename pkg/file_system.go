package build

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// FileSystem contains a few interfaces that we are using in acr-builder
type FileSystem interface {
	Getwd() (string, error)
	Chdir(path string) error
	DoesDirExist(path string) (bool, error)
	DoesFileExist(path string) (bool, error)
	IsDirEmpty(path string) (bool, error)
}

// ContextAware are any objects that are aware of build contexts
type ContextAware interface {
	SetContext(context *BuilderContext)
	GetContext() *BuilderContext
}

// ContextAwareFileSystem is a collections of file sysetm operations
// which will resolves the input with respect to BuilderContext
type ContextAwareFileSystem struct {
	context *BuilderContext
}

// NewContextAwareFileSystem creates a new file system with no context
func NewContextAwareFileSystem(context *BuilderContext) *ContextAwareFileSystem {
	return &ContextAwareFileSystem{context: context}
}

// SetContext changes the context that the file system uses
func (r *ContextAwareFileSystem) SetContext(context *BuilderContext) {
	r.context = context
}

// Getwd gets the current working directory
func (r *ContextAwareFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

// Chdir changes current working directory for the runner
func (r *ContextAwareFileSystem) Chdir(path string) error {
	path = r.context.Expand(path)
	logrus.Debugf("Chdir to %s", path)
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("Error chdir to %s", path)
	}
	return nil
}

// DoesDirExist checks if a given directory exists
func (r *ContextAwareFileSystem) DoesDirExist(path string) (bool, error) {
	return r.lookupPath(path, true)
}

// DoesFileExist checks if a given file exists
func (r *ContextAwareFileSystem) DoesFileExist(path string) (bool, error) {
	return r.lookupPath(path, false)
}

// IsDirEmpty checks if given directory is empty
func (r *ContextAwareFileSystem) IsDirEmpty(path string) (bool, error) {
	path = r.context.Expand(path)
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(info) == 0, nil
}

func (r *ContextAwareFileSystem) lookupPath(path string, isDir bool) (bool, error) {
	path = r.context.Expand(path)
	fileInfo, err := os.Stat(path)
	if err == nil {
		if fileInfo.IsDir() == isDir {
			return true, nil
		}
		err = fmt.Errorf("Path is expected to be IsDir: %t", isDir)
	} else if os.IsNotExist(err) {
		err = nil
	} else {
		logrus.Warnf("Unexpected error while getting path: %s", path)
	}
	return false, err
}
