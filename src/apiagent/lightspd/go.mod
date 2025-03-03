module github.com/snort3_aws/apiagent/lightspd

go 1.20

replace github.com/snort3_aws/apiagent/download => ../download

replace github.com/snort3_aws/apiagent/reload => ../reload

replace github.com/snort3_aws/message => ../../message

require (
	github.com/otiai10/copy v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/snort3_aws/apiagent/download v0.0.0-00010101000000-000000000000
	github.com/snort3_aws/apiagent/reload v0.0.0-00010101000000-000000000000
	github.com/snort3_aws/ipspolicy v0.0.0-00010101000000-000000000000
	github.com/snort3_aws/message v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/abrander/go-supervisord v0.0.0-20210517172913-a5469a4c50e2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/grpc v1.57.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20220521103104-8f96da9f5d5e // indirect
	k8s.io/api v0.22.3 // indirect
	k8s.io/apimachinery v0.22.3 // indirect
	k8s.io/client-go v0.22.3 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/snort3_aws/ipspolicy => ../../ipspolicy
