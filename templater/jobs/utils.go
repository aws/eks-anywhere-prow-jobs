package jobs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/eks-distro-prow-jobs/templater/jobs/types"
	"github.com/ghodss/yaml"
)

var releaseBranches = []string{
	"1-27",
	"1-28",
	"1-29",
	"1-30",
	"1-31",
}

func GetJobsByType(repos []string, jobType string) (map[string]map[string]types.JobConfig, error) {
	jobsListByType := map[string]map[string]types.JobConfig{}
	for _, repo := range repos {
		jobDir := filepath.Join("jobs", jobType, repo)

		jobList, err := UnmarshalJobs(jobDir)
		if err != nil {
			return nil, fmt.Errorf("reading job directory %s: %v", jobDir, err)
		}

		jobsListByType[fmt.Sprintf("aws/%s", repo)] = jobList
	}

	return jobsListByType, nil
}

func AppendMap(current, new map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range current {
		newMap[k] = v
	}
	for k, v := range new {
		newMap[k] = v
	}

	return newMap
}

func AddReleaseBranch(fileName string, data map[string]interface{}) map[string]map[string]interface{} {
	jobList := map[string]map[string]interface{}{}
	if !strings.Contains(fileName, "1-X") {
		return jobList
	}

	for i, releaseBranch := range releaseBranches {
		releaseBranchBasedFileName := strings.ReplaceAll(fileName, "1-X", releaseBranch)
		otherReleaseBranches := append(append([]string{}, releaseBranches[:i]...),
			releaseBranches[i+1:]...)
		jobList[releaseBranchBasedFileName] = AppendMap(data, map[string]interface{}{
			"releaseBranch":        releaseBranch,
			"otherReleaseBranches": strings.Join(otherReleaseBranches, "|"),
		})

		// If latest release branch, check if the release branch dir exists before executing cmd
		// This allows us to experiment with adding prow jobs for new branches without failing other runs
		if len(releaseBranches)-1 == i {
			jobList[releaseBranchBasedFileName]["latestReleaseBranch"] = true
		}
	}

	return jobList
}

func RunMappers(jobsToData map[string]map[string]interface{}, mappers []func(string, map[string]interface{}) map[string]map[string]interface{}) {
	if len(mappers) == 0 {
		return
	}

	for fileName, data := range jobsToData {
		newJobList := mappers[0](fileName, data)
		if len(newJobList) == 0 {
			continue
		}

		for k, v := range newJobList {
			jobsToData[k] = v
			if _, ok := data["templateFileName"]; !ok {
				jobsToData[k]["templateFileName"] = fileName
			}
		}
		delete(jobsToData, fileName)
	}

	RunMappers(jobsToData, mappers[1:])
}

func UnmarshalJobs(jobDir string) (map[string]types.JobConfig, error) {
	files, err := os.ReadDir(jobDir)
	if err != nil {
		return nil, fmt.Errorf("reading job directory %s: %v", jobDir, err)
	}

	var mappers []func(string, map[string]interface{}) map[string]map[string]interface{}
	mappers = append(mappers, AddReleaseBranch)

	jobsToData := map[string]map[string]interface{}{}
	for _, file := range files {
		jobsToData[file.Name()] = map[string]interface{}{}
	}

	RunMappers(jobsToData, mappers)

	finalJobList := map[string]types.JobConfig{}
	for fileName, data := range jobsToData {
		templateFileName := fileName
		if name, ok := data["templateFileName"]; ok {
			templateFileName = name.(string)
		}

		jobConfig, err := GenerateJobConfig(data, filepath.Join(jobDir, templateFileName))
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}

		if latest, ok := data["latestReleaseBranch"]; ok && latest.(bool) {
			for j, command := range jobConfig.Commands {
				jobConfig.Commands[j] = "if make check-for-supported-release-branch -C $PROJECT_PATH; then " + command + "; fi"
			}
		}

		finalJobList[fileName] = jobConfig

	}
	return finalJobList, nil
}

func ExecuteTemplate(templateContent string, data interface{}) ([]byte, error) {
	temp := template.New("template")
	funcMap := map[string]interface{}{
		"indent": func(spaces int, v string) string {
			pad := strings.Repeat(" ", spaces)
			return pad + strings.Replace(v, "\n", "\n"+pad, -1)
		},
		"stringsJoin": strings.Join,
		"trim":        strings.TrimSpace,
	}
	temp = temp.Funcs(funcMap)

	temp, err := temp.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = temp.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("substituting values for template: %v", err)
	}
	return buf.Bytes(), nil
}

func GenerateJobConfig(data interface{}, filePath string) (types.JobConfig, error) {
	var jobConfig types.JobConfig
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return jobConfig, fmt.Errorf("reading job YAML %s: %v", filePath, err)
	}
	var templatedContents []byte
	if data != nil {
		templatedContents, err = ExecuteTemplate(string(contents), data)
		if err != nil {
			return jobConfig, fmt.Errorf("executing template: %v", err)
		}
		err = yaml.Unmarshal(templatedContents, &jobConfig)
		if err != nil {
			return jobConfig, fmt.Errorf("unmarshaling contents of file %s: %v", filePath, err)
		}
	} else {
		err = yaml.Unmarshal(contents, &jobConfig)
		if err != nil {
			return jobConfig, fmt.Errorf("unmarshaling contents of file %s: %v", filePath, err)
		}
	}
	return jobConfig, nil
}

func IsCuratedPackagesPresubmit(config string) bool {
	return strings.Contains(config, "autoscaler") ||
		strings.Contains(config, "cloud-provider-aws") ||
		strings.Contains(config, "harbor") ||
		strings.Contains(config, "prometheus") ||
		config == "aws-otel-collector-tooling-presubmit" ||
		config == "distribution-tooling-presubmit" ||
		config == "eks-anywhere-packages-image-tooling-presubmit" ||
		config == "emissary-tooling-presubmit" ||
		config == "hello-eks-anywhere-tooling-presubmit" ||
		config == "metallb-tooling-presubmit" ||
		config == "metrics-server-presubmit" ||
		config == "redis-tooling-presubmit" ||
		config == "rolesanywhere-credential-helper-presubmit" ||
		config == "trivy-tooling-presubmit"
}
