module guthub.com/snort3_aws/apiagent/download

go 1.17

require (
	github.com/abrander/go-supervisord v0.0.0-20210517172913-a5469a4c50e2
	github.com/pkg/errors v0.9.1
)

require (
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/snort3_aws/ipspolicy v0.0.0-00010101000000-000000000000
	github.com/snort3_aws/message v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/snort3_aws/message => ../../message

replace github.com/snort3_aws/ipspolicy => ../../ipspolicy
