package lightspd

import (
	"log"
	"os"
	"io"
	"io/ioutil"
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/otiai10/copy"
	"github.com/snort3_aws/apiagent/download"
	"github.com/snort3_aws/apiagent/reload"
	"github.com/snort3_aws/message"
	"github.com/snort3_aws/ipspolicy"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	lspdPath = "/var/snort/lightspd"
	tmpDir = "/var/tmp"
	tmpLspdPath = tmpDir + "/lightspd"
	lspdFileName = "Talos_LightSPD.tar.gz"
	lspdReloadStateFile = "lspd_reload.json"
	stateSuccess = "success"
	stateFail = "fail"
	stateStart = "start"
	lspdJsonFile = "/var/snort/lspd.json"
	lspVerLen = 14
	PolicyVer = "3.1.0.0-0"
	ModuleVer = "3.1.44.0"
)

type LightSpdReload struct {
	state ReloadStatus
	reloadData *message.ReloadLsp
}

type ReloadStatus struct {
	Download string `json:"download"`
	Untar string `json:"untar"`
	StopSnort string `json:"stopSnort"`
	Swap string `json:"swap"`
	StartSnort string `json:"startSnort"`
}

func NewLightSpdReload(reloadSpd *message.ReloadLsp) *LightSpdReload {
	lspdr := LightSpdReload{}
	lspdr.state = ReloadStatus{}
	lspdr.reloadData = reloadSpd
	return &lspdr
}

func (lspdr *LightSpdReload) createSymLinks(policyJsonFile string) error {
	target := lspdPath + "/policies/" + PolicyVer
	current := lspdPath + "/policies/current"
	if err := os.Symlink(target, current); err != nil {
		return errors.Wrap(err, "failed to create policies symlink")
	}
	target = lspdPath + "/modules/" + ModuleVer
	current = lspdPath + "/modules/current"
	if err := os.Symlink(target, current); err != nil {
		return errors.Wrap(err, "failed to create modules symlink")
	}
	policy, err := reload.LoadPolicyData(policyJsonFile)
	if err != nil {
		policy = &message.IpsPolicy {
			PolicyName: ipspolicy.BalancedSecurityAndConnectivity,
		}
	}
	target = lspdPath + "/policies/current/" + policy.PolicyName + ".lua"
	current = lspdPath + "/policies/current/snort.lua"
	if err := os.Symlink(target, current); err != nil {
		return errors.Wrap(err, "failed to create snort.lua symlink")
	}
	return nil
}

func (lspdr *LightSpdReload) resetState() {
	lspdr.state.Download = stateStart
	lspdr.state.Untar = stateStart
	lspdr.state.StopSnort = stateStart
	lspdr.state.Swap = stateStart
	lspdr.state.StartSnort = stateStart
}

func (lspdr *LightSpdReload) storeReloadState() error {
	dat, err := json.Marshal(&lspdr.state)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(tmpDir + "/" + lspdReloadStateFile, dat, 0644)
	return err
}

func (lspdr *LightSpdReload) storeReloadData(path string) error {
	out, err := protojson.Marshal(lspdr.reloadData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, out, 0644)
	return err
}

func (lspdr *LightSpdReload) loadReloadState() error {
	lspdr.resetState()
	dat, err := ioutil.ReadFile(tmpDir + "/" + lspdReloadStateFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &lspdr.state)
	if err != nil {
		lspdr.resetState()
	}
	return err
}

func (lspdr *LightSpdReload) downloadLsp() error {
	dh := download.NewDownloadHandle()
	err := dh.Download("/" + lspdFileName, tmpDir + "/" + lspdFileName)
	if err != nil {
		return errors.Wrap(err, "failed to download spd package")
	}
	return nil
}

func (lspdr *LightSpdReload) untarLsp() error {
	if err := os.RemoveAll(tmpLspdPath); err != nil {
		return errors.Wrap(err, "failed to remove old spd")
	}

	gzStream, err := os.Open(tmpDir + "/" + lspdFileName)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	uncompressedStream, err := gzip.NewReader(gzStream)
	if err != nil {
		return errors.Wrap(err, "failed to read gzip stream")
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "ExtractTarGz: Next() failed")
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(tmpDir + "/" + header.Name, 0755); err != nil {
				return errors.Wrap(err, "ExtractTarGz: Mkdir() failed")
			}
		case tar.TypeReg:
			outFile, err := os.Create(tmpDir + "/" + header.Name)
			if err != nil {
				return errors.Wrap(err, "ExtractTarGz: Create() failed")
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.Wrap(err, "ExtractTarGz: Copy() failed")
			}
			outFile.Close()
		default:
			log.Fatalf("ExtractTarGz: uknown type: %v in %s", header.Typeflag, header.Name)
		}
	}

	log.Printf("successfully uncompressed spd tar.gz file")
	return nil
}

func (lspdr *LightSpdReload) cleanup() {
	if err := os.Remove("/var/tmp/" + lspdFileName); err != nil {
		log.Printf("failed to remove lightspd tar gz file")
	}
	if err := os.RemoveAll(tmpLspdPath); err != nil {
		log.Printf("failed to remove lightspd tar gz file")
	}
}

func (lspdr *LightSpdReload) getLoadedLsp(path string) (*message.ReloadLsp, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file %s\n", path)
		return nil, err
	}
	spd := &message.ReloadLsp{}
	err = protojson.Unmarshal(data, spd)
	if err != nil {
		log.Printf("failed to unmarshall spd data %v\n", err)
		return nil, err
	}
	return spd, nil
}

func (lspdr *LightSpdReload) getRequestedVer() string {
	return lspdr.reloadData.LspVersion
}

func (lspdr *LightSpdReload) checkLspVersion(path string) bool {
	req := lspdr.getRequestedVer()
	// force a fresh reload without using previous state
	if strings.Contains(req, "force") {
		lspdr.resetState()
		if err := lspdr.storeReloadState(); err != nil {
			log.Printf("failed to reset lsp reload state: %v\n", err)
		}
	}

	// 2021-11-09-001 version format
	if len(req) < lspVerLen  {
		log.Printf("requested lsp version %s\n", req)
		return true
	}
	lspConfig, err := lspdr.getLoadedLsp(path)
	if err != nil {
		log.Printf("failed to load lsp config: %v\n", err)
		return true
	}
	if len(lspConfig.LspVersion) < lspVerLen {
		log.Printf("loaded lsp version %s\n", lspConfig.LspVersion)
		return true
	}
	if req[:lspVerLen] != lspConfig.LspVersion[:lspVerLen] {
		return true
	}

	log.Printf("loaded version %s and request %s are the same\n", lspConfig.LspVersion, req)
	return false
}

func (lspdr *LightSpdReload) download() error {
	if lspdr.state.Download != stateSuccess {
		err := lspdr.downloadLsp()
		if err != nil {
			lspdr.state.Download = stateFail
			lspdr.storeReloadState()
			return errors.Wrap(err, "failed to download spd package")
		}
	}
	lspdr.state.Download = stateSuccess
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store download state: %v\n", err)
	}
	return nil
}

func (lspdr *LightSpdReload) untar() error {
	if lspdr.state.Untar != stateSuccess {
		err := lspdr.untarLsp()
		if err != nil {
			lspdr.state.Untar = stateFail
			lspdr.storeReloadState()
			return errors.Wrap(err, "failed to untar lsp package")
		}
	}
	lspdr.state.Untar = stateSuccess
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store untar state: %v\n", err)
	}
	return nil
}

func (lspdr *LightSpdReload) swap() error {
	// Copy will keep existing symlinks
	if lspdr.state.Swap != stateSuccess {
		err := os.RemoveAll(lspdPath)
		if err != nil {
			log.Printf("failed to delete dir %s\n", lspdPath)
		}
		copyOptions := copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			},
		}
		err = copy.Copy(tmpLspdPath, lspdPath, copyOptions)
		if err != nil {
			lspdr.state.Swap = stateFail
			lspdr.storeReloadState()
			return errors.Wrap(err, "failed to copy lightspd directory")
		}
		err = lspdr.createSymLinks(reload.PolicyJsonFile)
		if err != nil {
			lspdr.state.Swap = stateFail
			lspdr.storeReloadState()
			return errors.Wrap(err, "failed to create sym links")
		}
	}
	lspdr.state.Swap = stateSuccess
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store swap state: %v\n", err)
	}
	return nil
}

func (lspdr *LightSpdReload) stopSnort() error {
	if lspdr.state.StopSnort != stateSuccess {
		err := reload.StopSnort()
		if err != nil {
			lspdr.state.StopSnort = stateFail
			return errors.Wrap(err, "failed to stop snort")
		}
	}
	lspdr.state.StopSnort = stateSuccess
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store stop snort state: %v\n", err)
	}
	log.Printf("successfully stopped snort")
	return nil
}

func (lspdr *LightSpdReload) startSnort() error {
	if lspdr.state.StartSnort != stateSuccess {
		err := reload.StartSnort()
		if err != nil {
			lspdr.state.StartSnort = stateFail
			lspdr.storeReloadState()
			return errors.Wrap(err, "failed to start snort")
		}
	}
	lspdr.state.StartSnort = stateSuccess
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store sstart snort state: %v\n", err)
	}
	log.Printf("successfully started snort")
	return nil
}

func (lspdr *LightSpdReload) Reload() error {
	log.Printf("reload lightspd started...\n")
	if !lspdr.checkLspVersion(lspdJsonFile) {
		return nil
	}

	if err := lspdr.loadReloadState(); err != nil {
		log.Printf("failed to load reload state: %v\n", err)
	}
	if err := lspdr.download(); err != nil {
		return err
	}
	if err := lspdr.untar(); err != nil {
		return err
	}
	if err := lspdr.stopSnort(); err != nil {
		// ignore snort snort error and continue
		log.Printf("failed to stop snort %v\n", err)
	}
	if err := lspdr.swap(); err != nil {
		return err
	}
	if err := lspdr.startSnort(); err != nil {
		return err
	}

	lspdr.resetState()
	if err := lspdr.storeReloadState(); err != nil {
		log.Printf("failed to store reload state: %v\n", err)
	}
	if err := lspdr.storeReloadData(lspdJsonFile); err != nil {
		log.Printf("failed to store reload data: %v\n", err)
	}
	log.Printf("successfully reloaded talos light spd")
	lspdr.cleanup()
	return nil
}
