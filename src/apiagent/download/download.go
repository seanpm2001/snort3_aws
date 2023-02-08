package download

import (
	"context"
	"log"
	"strings"
	"net/http"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type DownloadHandle struct {
	storageServerIp string
}

func NewDownloadHandle() *DownloadHandle {
	return &DownloadHandle{storageServerIp: ""}
}

func (dh *DownloadHandle) getStorageServerIp() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	namespace := os.Getenv("POD_NAMESPACE")
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for i := 0; i < len(pods.Items); i++ {
		if strings.HasPrefix(pods.Items[i].Name, "storage-") {
			dh.storageServerIp = pods.Items[i].Status.PodIP
			log.Printf("Found storage pod with name: %s IP: %s\n", pods.Items[i].Name, pods.Items[i].Status.PodIP)
			break
		}
	}
	if dh.storageServerIp == "" {
		return errors.New("storage pod not found")
	}
	return nil
}

func (dh *DownloadHandle) constructUrl(path string) string {
	return "http://" + dh.storageServerIp + path
}

func (dh *DownloadHandle) Download(path string, dest string) error {
	if dh.storageServerIp == "" {
		err := dh.getStorageServerIp()
		if err != nil {
			return errors.Wrap(err, "failed to get storage server ip")
		}
	}
	resp, err := http.Get(dh.constructUrl(path))
	if err != nil {
		return errors.Wrap(err, "failed to run http get")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("%s http respose code: %v\n", path, resp.StatusCode)
		return errors.New("failed to download " + path + ", http response code is " + strconv.Itoa(resp.StatusCode))
	}
	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy file from resp body")
	}
	log.Printf("successfully downloaded file %s\n", path)
	return nil
}
