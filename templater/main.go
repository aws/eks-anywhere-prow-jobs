package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-distro-prow-jobs/templater/jobs/types"
	"github.com/aws/eks-distro-prow-jobs/templater/jobs/utils"

	"github.com/aws/eks-anywhere-prow-jobs/templater/jobs"
)

var (
	jobsFolder    = "jobs"
	orgsSupported = []string{"aws"}
	jobTypes      = []string{"periodic", "postsubmit", "presubmit"}
)

//go:embed templates/presubmits.yaml
var presubmitTemplate string

//go:embed templates/postsubmits.yaml
var postsubmitTemplate string

//go:embed templates/periodics.yaml
var periodicTemplate string

//go:embed templates/warning.txt
var editWarning string

//go:generate cp ../BUILDER_BASE_TAG_FILE ./BUILDER_BASE_TAG_FILE
//go:embed BUILDER_BASE_TAG_FILE
var builderBaseTag string

var buildkitImageTag = "v0.12.5-rootless"

func main() {
	jobsFolderPath, err := getJobsFolderPath()
	if err != nil {
		fmt.Printf("Error getting jobs folder path: %v", err)
		os.Exit(1)
	}

	for _, org := range orgsSupported {
		if err = os.RemoveAll(filepath.Join(jobsFolderPath, org)); err != nil {
			fmt.Printf("Error removing jobs folder path: %v", err)
			os.Exit(1)
		}
	}

	for _, jobType := range jobTypes {
		jobList, err := jobs.GetJobList(jobType)
		if err != nil {
			fmt.Printf("Error getting job list: %v\n", err)
			os.Exit(1)
		}
		template, err := useTemplate(jobType)
		if err != nil {
			fmt.Printf("Error getting job list: %v\n", err)
			os.Exit(1)
		}

		for repoName, jobConfigs := range jobList {
			for fileName, jobConfig := range jobConfigs {
				envVars := jobConfig.EnvVars

				if jobConfig.UseDockerBuildX {
					envVars = append(envVars, &types.EnvVar{Name: "BUILDKITD_IMAGE", Value: "moby/buildkit:" + buildkitImageTag})
					envVars = append(envVars, &types.EnvVar{Name: "USE_BUILDX", Value: "true"})
				}

				if jobConfig.ImageBuild {
					if jobs.IsCuratedPackagesPresubmit(jobConfig.JobName) {
						envVars = append(envVars, &types.EnvVar{Name: "ADDITIONAL_IMAGE_CACHE_REPOS", Value: "067575901363.dkr.ecr.us-west-2.amazonaws.com"})
					} else {
						envVars = append(envVars, &types.EnvVar{Name: "ADDITIONAL_IMAGE_CACHE_REPOS", Value: "857151390494.dkr.ecr.us-west-2.amazonaws.com"})
					}
				}

				cluster, bucket, serviceAccountName := clusterDetails(jobType, jobConfig.Cluster, jobConfig.Bucket, jobConfig.ServiceAccountName)

				data := map[string]interface{}{
					"architecture":                 jobConfig.Architecture,
					"repoName":                     repoName,
					"prowjobName":                  jobConfig.JobName,
					"runIfChanged":                 jobConfig.RunIfChanged,
					"skipIfOnlyChanged":            jobConfig.SkipIfOnlyChanged,
					"branches":                     jobConfig.Branches,
					"cronExpression":               jobConfig.CronExpression,
					"maxConcurrency":               jobConfig.MaxConcurrency,
					"timeout":                      jobConfig.Timeout,
					"extraRefs":                    jobConfig.ExtraRefs,
					"imageBuild":                   jobConfig.ImageBuild,
					"useDockerBuildX":              jobConfig.UseDockerBuildX,
					"prCreation":                   jobConfig.PRCreation,
					"runtimeImage":                 jobConfig.RuntimeImage,
					"localRegistry":                jobConfig.LocalRegistry,
					"serviceAccountName":           serviceAccountName,
					"command":                      strings.Join(jobConfig.Commands, "\n&&\n"),
					"builderBaseTag":               builderBaseTag,
					"buildkitImageTag":             buildkitImageTag,
					"resources":                    jobConfig.Resources,
					"envVars":                      envVars,
					"volumes":                      jobConfig.Volumes,
					"volumeMounts":                 jobConfig.VolumeMounts,
					"editWarning":                  editWarning,
					"automountServiceAccountToken": jobConfig.AutomountServiceAccountToken,
					"cluster":                      cluster,
					"bucket":                       bucket,
					"projectPath":                  jobConfig.ProjectPath,
					"diskUsage":                    true,
					"runAsUser":                    jobConfig.RunAsUser,
					"runAsGroup":                   jobConfig.RunAsGroup,
				}

				err := GenerateProwjob(fileName, template, data)
				if err != nil {
					fmt.Printf("Error generating Prowjob %s: %v\n", fileName, err)
					os.Exit(1)
				}
			}
		}
	}
}

func GenerateProwjob(prowjobFileName, templateContent string, data map[string]interface{}) error {
	bytes, err := utils.ExecuteTemplate(templateContent, data)
	if err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	jobsFolderPath, err := getJobsFolderPath()
	if err != nil {
		return fmt.Errorf("getting jobs folder path: %v", err)
	}

	prowjobPath := filepath.Join(jobsFolderPath, data["repoName"].(string), prowjobFileName)
	if err = os.MkdirAll(filepath.Dir(prowjobPath), 0o755); err != nil {
		return fmt.Errorf("creating Prowjob directory: %v", err)
	}

	if err = os.WriteFile(prowjobPath, bytes, 0o644); err != nil {
		return fmt.Errorf("writing to path %s: %v", prowjobPath, err)
	}

	return nil
}

func getJobsFolderPath() (string, error) {
	gitRootOutput, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("running the git command: %v", err)
	}
	gitRoot := strings.Fields(string(gitRootOutput))[0]

	return filepath.Join(gitRoot, jobsFolder), nil
}

func useTemplate(jobType string) (string, error) {
	switch jobType {
	case "periodic":
		return periodicTemplate, nil
	case "postsubmit":
		return postsubmitTemplate, nil
	case "presubmit":
		return presubmitTemplate, nil
	default:
		return "", fmt.Errorf("unsupported job type: %s", jobType)
	}
}

func clusterDetails(jobType, cluster, bucket, serviceAccountName string) (string, string, string) {
	if jobType == "presubmit" && len(cluster) == 0 {
		cluster = "prow-presubmits-cluster"
		bucket = "s3://prowpresubmitsdataclusterstack-prowbucket7c73355c-vfwwxd2eb4gp"
		serviceAccountName = "presubmits-build-account"
	}

	if (jobType == "postsubmit" || jobType == "periodic") && len(cluster) == 0 {
		cluster = "prow-postsubmits-cluster"
		bucket = "s3://prowdataclusterstack-316434458-prowbucket7c73355c-1n9f9v93wpjcm"
		serviceAccountName = "postsubmits-build-account"
	}

	return cluster, bucket, serviceAccountName
}
