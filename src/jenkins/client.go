/**
 * @Time : 2019-06-25 17:52
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package jenkins

import (
	"fmt"
	"github.com/icowan/config"
	"strings"
)

type Params struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	GitAddr    string `json:"git_addr"`
	GitType    string `json:"git_type"`
	GitVersion string `json:"git_version"`
	Command    string `json:"command"`
}

func NewJenkins(cf *config.Config) (Jenkins, error) {

	fmt.Printf("姜猛修改标志--jenkins连接信息---" + cf.GetString("jenkins", "host"))
	fmt.Printf("姜猛修改标记,测试查看读取的配置文件信息  " + cf.GetEnv() + "   打印结束")

	return NewService(&Auth{
		Username: cf.GetString("jenkins", "user"),
		ApiToken: cf.GetString("jenkins", "token"),
	}, cf.GetString("jenkins", "host"), "v3")
}

func (c *service) CreateJobParams(params Params) error {
	var gitVersion = params.GitVersion
	if index := strings.Index(gitVersion, ":"); index >= 0 {
		gitVersion = gitVersion[index+1:]
	}

	var mavenJobItem = MavenJobItem{}

	mavenJobItem.Builders.HudsonTasksShell.Command = params.Command
	mavenJobItem.Scm.Class = "hudson.plugins.git.GitSCM"
	mavenJobItem.Scm.Plugin = "git@3.10.1"
	mavenJobItem.Scm.ConfigVersion = "2"
	mavenJobItem.Scm.SubmoduleCfg.Class = "list"
	mavenJobItem.Scm.UserRemoteConfigs.GitUserRemoteConfig.Url = params.GitAddr
	mavenJobItem.Scm.UserRemoteConfigs.GitUserRemoteConfig.CredentialsId = "jenkins"

	if params.GitType == "branch" {
		mavenJobItem.Scm.Branches.BranchSpec.Name = fmt.Sprintf("*/%s", gitVersion)
	} else if params.GitType == "tag" {
		mavenJobItem.Scm.Branches.BranchSpec.Name = "tags/$TAGNAME"
	}

	var parameterDefinitions []HudsonModelTextParameterDefinition
	parameterDefinitions = append(parameterDefinitions, HudsonModelTextParameterDefinition{Name: "TAGNAME", Trim: "false"})
	mavenJobItem.Properties.HudsonModelParametersDefinitionProperty.ParameterDefinitions.HudsonModelTextParameterDefinition = parameterDefinitions

	fmt.Printf("姜猛修改标志--------------%%%%%%%%%%%%------------------%%%%%%%%-------------------------------------------")
	fmt.Printf(mavenJobItem.Builders.HudsonTasksShell.Command)

	return c.CreateJob(mavenJobItem, params.Name+"."+params.Namespace)
}

func (c *service) CreateJavaJobParams(params Params) error {

	fmt.Printf("姜猛修改标志-----测试jenkins的shell指令：" + params.Command)

	var gitVersion = params.GitVersion
	if index := strings.Index(gitVersion, ":"); index >= 0 {
		gitVersion = gitVersion[index+1:]
	}

	var parameterDefinitions []HudsonModelTextParameterDefinition
	parameterDefinitions = append(parameterDefinitions, HudsonModelTextParameterDefinition{
		Name:        "TAGNAME",
		Trim:        "false",
		Description: "Tags名称或Branch名称",
	}, HudsonModelTextParameterDefinition{
		Name:        "POMFILE",
		Description: "pom.xml路径",
	}, HudsonModelTextParameterDefinition{
		Name:        "BUILD_ENV",
		Description: "build环境 dev,pre,online",
	}, HudsonModelTextParameterDefinition{
		Name:        "BRANCH",
		Description: "分支名",
	})
	var mavenJobItem MavenJobItem
	mavenJobItem = MavenJobItem{
		Plugin:           "maven-plugin@3.1.2",
		KeepDependencies: "false",
		Properties: Properties{
			HudsonModelParametersDefinitionProperty: HudsonModelParametersDefinitionProperty{
				ParameterDefinitions: ParameterDefinitions{
					HudsonModelTextParameterDefinition: parameterDefinitions,
				},
			},
		},
		Scm: Scm{
			Class:         "hudson.plugins.git.GitSCM",
			Plugin:        "git@3.10.1",
			ConfigVersion: "2",
			SubmoduleCfg:  SubmoduleCfg{Class: "list"},
			UserRemoteConfigs: UserRemoteConfigs{
				GitUserRemoteConfig: GitUserRemoteConfig{
					Url:           params.GitAddr,
					CredentialsId: "jenkins",
				},
			},
			Branches: Branches{
				BranchSpec: HudsonPluginsGitBranchSpec{Name: "${BRANCH}"},
			},
		},
		CanRoam:                          "true",
		Disabled:                         "false",
		BlockBuildWhenUpstreamBuilding:   "false",
		BlockBuildWhenDownstreamBuilding: "false",
		ConcurrentBuild:                  "false",
		RootPOM:                          "$POMFILE",
		Goals:                            "-Dmaven.test.skip=true ${BUILD_ENV} -U clean package",
		AggregatorStyleBuild:             "true",
		IncrementalBuild:                 "false",
		IgnoreUpstremChanges:             "true",
		ArchivingDisabled:                "false",
		SiteArchivingDisabled:            "false",
		FingerprintingDisabled:           "false",
		ResolveDependencies:              "false",
		ProcessPlugins:                   "false",
		MavenValidationLevel:             "-1",
		RunHeadless:                      "false",
		DisableTriggerDownstreamProjects: "false",
		Settings: Settings{
			Class: "jenkins.mvn.DefaultSettingsProvider",
		},
		GlobalSettings: Settings{
			Class: "jenkins.mvn.DefaultGlobalSettingsProvider",
		},
		Postbuilders: Postbuilders{
			HudsonTasksShell: HudsonTasksShell{
				Command: params.Command,
			},
		},
		//姜猛修改标记，新增了Builders标签的内容
		Builders: Postbuilders{
			HudsonTasksShell: HudsonTasksShell{
				Command: params.Command,
			},
		},
		RunPostStepsIfResult: RunPostStepsIfResult{
			Name:          "Success",
			Ordinal:       "0",
			Color:         "BLUE",
			CompleteBuild: "true",
		},
	}
	command := mavenJobItem.Postbuilders.HudsonTasksShell.Command
	fmt.Printf("姜猛修改标志-----测试打印shell的xml：" + command)
	return c.CreateJob(mavenJobItem, params.Name+"."+params.Namespace)
}
