module github.com/IBM/ibmcloud-volume-interface

go 1.15

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v0.6.7
	github.com/IBM/secret-common-lib v1.0.5
	github.com/IBM/secret-utils-lib v1.0.4
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
)

replace (
	k8s.io/api => k8s.io/api v0.25.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.0
	k8s.io/client-go => k8s.io/client-go v0.25.0
)
