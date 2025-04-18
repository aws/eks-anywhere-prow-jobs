package jobs

import (
	"fmt"

	"github.com/aws/eks-distro-prow-jobs/templater/jobs/types"
)

func GetJobList(jobType string) (map[string]map[string]types.JobConfig, error) {
	switch jobType {
	case "periodic":
		repos := []string{"eks-anywhere", "eks-anywhere-build-tooling"}
		periodicsList, err := GetJobsByType(repos, "periodic")
		if err != nil {
			return nil, fmt.Errorf("getting periodic list:%v", err)
		}
		return periodicsList, nil
	case "postsubmit":
		repos := []string{"eks-anywhere"}
		postsubmitsList, err := GetJobsByType(repos, "postsubmit")
		if err != nil {
			return nil, fmt.Errorf("getting postsubmits list:%v", err)
		}
		return postsubmitsList, nil
	case "presubmit":
		repos := []string{"eks-anywhere", "eks-anywhere-build-tooling", "eks-anywhere-packages", "eks-anywhere-prow-jobs"}
		presubmitsList, err := GetJobsByType(repos, "presubmit")
		if err != nil {
			return nil, fmt.Errorf("getting presubmits list:%v", err)
		}
		return presubmitsList, nil
	default:
		return nil, fmt.Errorf("unsupported job type: %s", jobType)
	}
}
