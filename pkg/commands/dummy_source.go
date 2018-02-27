package commands

import (
	build "github.com/Azure/acr-builder/pkg"
)

type dummySource struct {
}

// NewDummySource create a dummy source
func NewDummySource() build.Source {
	return &dummySource{}
}

func (s *dummySource) Return(runner build.Runner) error {
	return nil
}

func (s *dummySource) Obtain(runner build.Runner) error {
	return nil
}

func (s *dummySource) Export() []build.EnvVar {
	return []build.EnvVar{}
}
