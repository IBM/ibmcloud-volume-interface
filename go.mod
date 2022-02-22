module github.com/IBM/ibmcloud-volume-interface

go 1.15

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v0.6.7
	github.com/IBM/secret-common-lib v0.0.0-20220222031547-939ad5dfc3a9
	github.com/IBM/secret-utils-lib v0.0.0-20220222031021-e3e6d5002fff
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
)

replace (
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
)
