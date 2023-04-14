# Amazon EKS Anywhere Prow Jobs

This repository contains Prowjob configuration for the Amazon EKS Anywhere project, which includes the [eks-anywhere](https://github.com/aws/eks-anywhere) and [eks-anywhere-build-tooling](https://github.com/aws/eks-anywhere-build-tooling) repositories. Jobs can be viewed on the EKS installation of
Prow, which is available at https://prow.eks.amazonaws.com/.

For more info on how to write a Prowjob, read the [test-infra
introduction](https://github.com/kubernetes/test-infra/blob/master/prow/jobs.md)
to prow jobs.

## Contributing

Please read our [CONTRIBUTING](CONTRIBUTING.md) guide before making a pull
request.

## Security

If you discover a potential security issue in this project, or think you may
have discovered a security issue, we ask that you notify AWS Security via our
[vulnerability reporting
page](http://aws.amazon.com/security/vulnerability-reporting/). Please do
**not** create a public GitHub issue.

## License

This project is licensed under the Apache-2.0 License.

## Adding/Removing Kubernetes Versions

Kubernetes/Release branch versions are managed [here](https://github.com/aws/eks-distro-prow-jobs/blob/c23ec5a128f01b11672886da3c9a6da3c2bb846b/templater/jobs/utils/utils.go#L16).

To update the templates, update the eks-distro-prow-jobs module:
```
go get github.com/aws/eks-distro-prow-jobs
make prowjobs -C templater
```
