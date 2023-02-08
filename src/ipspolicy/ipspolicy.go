package ipspolicy

import (
	"github.com/pkg/errors"
)

const (
	BalancedSecurityAndConnectivity = "balanced-security-and-connectivity"
	ConnectivityOverSecurity        = "connectivity-over-security"
	MaximumDetection                = "maximum-detection"
	NoRulesActive                   = "no-rules-active"
	SecurityOverConnectivity        = "security-over-connectivity"
)

func ValidatePolicyName(policy string) error {
	switch policy {
	case BalancedSecurityAndConnectivity:
		return nil
	case ConnectivityOverSecurity:
		return nil
	case MaximumDetection:
		return nil
	case NoRulesActive:
		return nil
	case SecurityOverConnectivity:
		return nil
	}
	return errors.New("invalid policy " + policy)
}
