package reload

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/abrander/go-supervisord"
	"github.com/pkg/errors"
	"github.com/snort3_aws/message"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	PolicyDir       = "/var/snort/lightspd/policies"
	PolicyVer       = "current"
	SupervisordSock = "/var/run/supervisor.sock"
	PolicyJsonFile  = "/var/snort/policy.json"
)

type SnortReload struct {
	ipsPolicy *message.IpsPolicy
}

func NewSnortReload(policy *message.IpsPolicy) *SnortReload {
	sr := SnortReload{ipsPolicy: policy}
	return &sr
}

func (sr *SnortReload) storePolicyData(filePath string) error {
	out, err := protojson.Marshal(sr.ipsPolicy)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	return err
}

func LoadPolicyData(filePath string) (*message.IpsPolicy, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("failed to read file %s\n", filePath)
		return nil, err
	}
	policy := &message.IpsPolicy{}
	err = protojson.Unmarshal(data, policy)
	if err != nil {
		log.Printf("failed to unmarshall policy data %v\n", err)
		return nil, err
	}
	return policy, nil
}

func (sr *SnortReload) UpdateSnortConfig() error {
	snort_lua_path := filepath.Join(PolicyDir+"/"+PolicyVer, "snort.lua")
	if err := os.Remove(snort_lua_path); err != nil {
		return errors.Wrap(err, "failed to remove snort.lua")
	}
	target := filepath.Join(PolicyDir+"/"+PolicyVer, sr.ipsPolicy.PolicyName+".lua")
	if err := os.Symlink(target, snort_lua_path); err != nil {
		return errors.Wrap(err, "failed to create snort.lua symlink")
	}
	return nil
}

func HupSnort() error {
	client, err := supervisord.NewUnixSocketClient(SupervisordSock)
	if err != nil {
		return errors.Wrap(err, "failed to create supervisord client")
	}
	err = client.SignalProcess("snort", syscall.SIGHUP)
	if err != nil {
		return errors.Wrap(err, "failed to send SIGHUP to snort")
	}

	return nil
}

func StopSnort() error {
	client, err := supervisord.NewUnixSocketClient(SupervisordSock)
	if err != nil {
		return errors.Wrap(err, "failed to create supervisord client")
	}
	err = client.StopProcess("snort", true)
	if err != nil {
		return errors.Wrap(err, "failed to stop snort")
	}
	return nil
}

func StartSnort() error {
	client, err := supervisord.NewUnixSocketClient(SupervisordSock)
	if err != nil {
		return errors.Wrap(err, "failed to create supervisord client")
	}
	err = client.StartProcess("snort", true)
	if err != nil {
		return errors.Wrap(err, "failed to start snort")
	}
	return nil
}

func RestartSnort() error {
	client, err := supervisord.NewUnixSocketClient(SupervisordSock)
	if err != nil {
		return errors.Wrap(err, "failed to create supervisord client")
	}
	err = client.StopProcess("snort", true)
	if err != nil {
		return errors.Wrap(err, "failed to stop snort")
	}
	err = client.StartProcess("snort", true)
	if err != nil {
		return errors.Wrap(err, "failed to restart snort")
	}
	return nil
}

func (sr *SnortReload) Reload() error {
	if err := sr.UpdateSnortConfig(); err != nil {
		return errors.Wrap(err, "failed to update snort.lua")
	}
	if err := HupSnort(); err != nil {
		return errors.Wrap(err, "failed to hup snort for reload")
	}
	if err := sr.storePolicyData(PolicyJsonFile); err != nil {
		log.Printf("failed to store policy data in %s\n", PolicyJsonFile)
	}
	return nil
}
