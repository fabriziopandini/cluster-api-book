module github.com/fabriziopandini/cluster-api-website/hack/tools

go 1.19

replace github.com/fabriziopandini/cluster-api-website => ../../

require (
	github.com/onsi/gomega v1.24.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	k8s.io/utils v0.0.0-20221101230645-61b03e2f6476
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
