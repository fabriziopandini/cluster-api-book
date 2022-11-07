module github.com/fabriziopandini/cluster-api-website/hack/tools

go 1.19

replace github.com/fabriziopandini/cluster-api-website => ../../

require (
	github.com/spf13/pflag v1.0.5
	k8s.io/utils v0.0.0-20221101230645-61b03e2f6476
)
