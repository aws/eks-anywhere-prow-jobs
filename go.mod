module github.com/aws/eks-anywhere-prow-jobs

go 1.13

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/client-go => k8s.io/client-go v0.20.1
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3-0.20210301154926-12660d4f2255
)

require (
	k8s.io/test-infra v0.0.0-20210608224924-94f3f2343d63
	sigs.k8s.io/yaml v1.2.0
)
