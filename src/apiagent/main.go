package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/snort3_aws/apiagent/lightspd"
	"github.com/snort3_aws/apiagent/reload"
	"github.com/snort3_aws/ipspolicy"
	"github.com/snort3_aws/message"
	"google.golang.org/grpc"
)

var (
	serverAddr      = []string{"0.0.0.0:60011"}
	reloadMutex     sync.Mutex
)

type apiServer struct {
	message.UnimplementedMessageServer
	addr string
}

func initAgent() {
	reloadMutex = sync.Mutex{}
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen ", err)
	}

	s := grpc.NewServer()
	message.RegisterMessageServer(s, &apiServer{addr: addr})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve ", err)
	}
}

func (s *apiServer) ReloadIpsPolicy(ctx context.Context, policy *message.IpsPolicy) (*message.Response, error) {
	log.Printf("Received ips reload message: %v\n", policy)
	if err := ipspolicy.ValidatePolicyName(policy.PolicyName); err != nil {
		return &message.Response{Status: "invalid policy"}, err
	}
	reloadMutex.Lock()
	sr := reload.NewSnortReload(policy)
	reloadMutex.Unlock()
	if err := sr.Reload(); err != nil {
		log.Printf("failed to reload snort %s\n", err)
		return &message.Response{Status: "failed"}, err
	}
	log.Printf("successfully reloaded snort")
	return &message.Response{Status: "ok"}, nil
}

func (s *apiServer) ReloadTalosLsp(ctx context.Context, spdVersion *message.ReloadLsp) (*message.Response, error) {
	log.Printf("Received spd reload message: %v\n", spdVersion)
	reloadMutex.Lock()
	lsp := lightspd.NewLightSpdReload(spdVersion)
	reloadMutex.Unlock()
	if err := lsp.Reload(); err != nil {
		log.Printf("failed to reload talos spd %s\n", err)
		return &message.Response{Status: "failed"}, err
	}
	return &message.Response{Status: "ok"}, nil
}

func main() {
	log.Println("starting api server")
	initAgent()

	var wg sync.WaitGroup
	for _, addr := range serverAddr {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}

	wg.Wait()
}
