package test

import "github.com/stretchr/testify/suite"

// Suite Suite
type Suite struct {
	suite.Suite
}

// Option Option
type Option func(*Suite)

// New new
func New(opts ...Option) *Suite {
	testSuite := &Suite{}
	for _, opt := range opts {
		opt(testSuite)
	}
	return testSuite
}

// SetupTest SetupTest
func (suite *Suite) SetupTest() {
}
