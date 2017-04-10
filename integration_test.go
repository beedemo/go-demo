// +build integration

package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	hostIp      string
}

func (s *IntegrationTestSuite) SetupTest() {
}

// Integration

func (s IntegrationTestSuite) Test_Hello_ReturnsStatus200() {
	address := fmt.Sprintf("http://%s/demo/hello", s.hostIp)
	resp, err := http.Get(address)

	s.NoError(err)
	s.Equal(200, resp.StatusCode)
}

func (s IntegrationTestSuite) Test_Person_ReturnsStatus200() {
	address := fmt.Sprintf("http://%s/demo/person", s.hostIp)
	resp, err := http.Get(address)

	s.NoError(err)
	s.Equal(200, resp.StatusCode)
}

// Suite

func TestIntegrationTestSuite(t *testing.T) {
	s := new(IntegrationTestSuite)
	s.hostIp = os.Getenv("HOST_IP")
	suite.Run(t, s)
}