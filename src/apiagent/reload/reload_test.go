package reload

import (
	"os"
	"testing"

	"github.com/snort3_aws/ipspolicy"
	"github.com/snort3_aws/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReloadTestSuite struct {
	suite.Suite
	sr       *SnortReload
	tempFile string
}

func (s *ReloadTestSuite) BeforeTest(suiteName, testName string) {
	policy := message.IpsPolicy{
		PolicyName: ipspolicy.ConnectivityOverSecurity,
	}
	s.sr = NewSnortReload(&policy)
	s.tempFile = "/var/tmp/policy.json"
}

func (s *ReloadTestSuite) AfterTest(suiteName, testName string) {
	os.Remove(s.tempFile)
}

func (s *ReloadTestSuite) TestStorePolicyData() {
	err := s.sr.storePolicyData(s.tempFile)
	assert.Equal(s.T(), nil, err)
	policy, err := LoadPolicyData(s.tempFile)
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), ipspolicy.ConnectivityOverSecurity, policy.PolicyName)

	policy = &message.IpsPolicy{
		PolicyName: "test-policy",
	}
	s.sr = NewSnortReload(policy)
	err = s.sr.storePolicyData(s.tempFile)
	assert.Equal(s.T(), nil, err)
	policy, err = LoadPolicyData(s.tempFile)
	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "test-policy", policy.PolicyName)
}

func TestReloadTestSuite(t *testing.T) {
	suite.Run(t, new(ReloadTestSuite))
}
