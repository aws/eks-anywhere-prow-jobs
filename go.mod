module github.com/aws/eks-anywhere-prow-jobs

go 1.16

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.8
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
	github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.2.0
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc95
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.8
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.6.1
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.1
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3-0.20210301154926-12660d4f2255
)

require (
	github.com/aws/eks-distro-prow-jobs v0.0.0-20230103185424-2eeb8942f6a8
	k8s.io/api v0.21.1
	k8s.io/test-infra v0.0.0-20210608224924-94f3f2343d63
	sigs.k8s.io/yaml v1.2.0
)
