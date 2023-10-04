module guthub.com/snort3_aws/apiagent/download

go 1.20

require (
	github.com/abrander/go-supervisord v0.0.0-20210517172913-a5469a4c50e2
	github.com/pkg/errors v0.9.1
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20220521103104-8f96da9f5d5e // indirect
)

require (
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/snort3_aws/ipspolicy v0.0.0-00010101000000-000000000000
	github.com/snort3_aws/message v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.27.1
)

replace github.com/snort3_aws/message => ../../message

replace github.com/snort3_aws/ipspolicy => ../../ipspolicy
