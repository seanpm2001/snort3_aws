package lightspd

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/snort3_aws/message"
)

type LspdrTestSuite struct {
	suite.Suite
	lspdr *LightSpdReload
}

func (s *LspdrTestSuite) BeforeTest(suiteName, testName string) {
	lspdReload := message.ReloadLsp{
		LspVersion: "2021-11-09-001",
	}
	s.lspdr = NewLightSpdReload(&lspdReload)
}

func (s *LspdrTestSuite) AfterTest(suiteName, testName string) {
}

func (s *LspdrTestSuite) TestResetState() {
	s.lspdr.resetState()
	assert.Equal(s.T(), stateStart, s.lspdr.state.Download)
	assert.Equal(s.T(), stateStart, s.lspdr.state.Untar)
	assert.Equal(s.T(), stateStart, s.lspdr.state.Swap)
	assert.Equal(s.T(), stateStart, s.lspdr.state.StopSnort)
	assert.Equal(s.T(), stateStart, s.lspdr.state.StartSnort)
}

func (s *LspdrTestSuite) TestStoreReloadState() {
	s.lspdr.state.Download = stateSuccess
	s.lspdr.state.Untar = stateSuccess
	s.lspdr.state.StopSnort = stateSuccess
	s.lspdr.state.Swap = stateSuccess
	s.lspdr.state.StartSnort = stateFail
	err := s.lspdr.storeReloadState()
	assert.Equal(s.T(), nil, err)

	err = s.lspdr.loadReloadState()
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Download)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Untar)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.StopSnort)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Swap)
	assert.Equal(s.T(), stateFail, s.lspdr.state.StartSnort)
}

func (s *LspdrTestSuite) TestStoreReloadData() {
	err := s.lspdr.storeReloadData("/var/tmp/lspd.ver")
	assert.Equal(s.T(), nil, err)
}

func (s *LspdrTestSuite) TestGetRequestedVer() {
	req := s.lspdr.getRequestedVer()
	assert.Equal(s.T(), "2021-11-09-001", req)
}

func (s *LspdrTestSuite) TestGetLoadedSpd() {
	err := s.lspdr.storeReloadData("/var/tmp/lspd.ver")
	lspdConfig, err := s.lspdr.getLoadedLsp("/var/tmp/lspd.ver")
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "2021-11-09-001", lspdConfig.LspVersion)
}

func (s *LspdrTestSuite) TestCheckLspdVersion() {
	lspdReload := message.ReloadLsp{
		LspVersion: "2020/09/14",
	}
	s.lspdr.reloadData = &lspdReload
	err := s.lspdr.storeReloadData("/var/tmp/lspd.ver")
	assert.Equal(s.T(), nil, err)
	lspdReload = message.ReloadLsp{
                LspVersion: "2020/10/14",
        }
	s.lspdr.state.Download = stateSuccess
	s.lspdr.state.Untar = stateSuccess
	s.lspdr.state.StopSnort = stateFail
	s.lspdr.state.Swap = stateSuccess
	s.lspdr.state.StartSnort = stateFail
	s.lspdr.reloadData = &lspdReload
	check := s.lspdr.checkLspVersion("/var/tmp/lspd.ver")
	assert.Equal(s.T(), true, check)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Download)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Untar)
	assert.Equal(s.T(), stateFail, s.lspdr.state.StopSnort)
	assert.Equal(s.T(), stateSuccess, s.lspdr.state.Swap)
	assert.Equal(s.T(), stateFail, s.lspdr.state.StartSnort)
}

func TestLspdrTestSuite(t *testing.T) {
        suite.Run(t, new(LspdrTestSuite))
}
