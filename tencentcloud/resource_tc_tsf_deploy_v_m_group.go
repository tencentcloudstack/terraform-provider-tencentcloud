/*
Provides a resource to create a tsf deploy_v_m_group

Example Usage

```hcl
resource "tencentcloud_tsf_deploy_v_m_group" "deploy_v_m_group" {
  group_id = ""
  pkg_id = ""
  startup_parameters = ""
  deploy_desc = ""
  force_start =
  enable_health_check =
  health_check_settings {
		liveness_probe {
			action_type = ""
			initial_delay_seconds =
			timeout_seconds =
			period_seconds =
			success_threshold =
			failure_threshold =
			scheme = ""
			port =
			path = ""
			command =
			type = ""
		}
		readiness_probe {
			action_type = ""
			initial_delay_seconds =
			timeout_seconds =
			period_seconds =
			success_threshold =
			failure_threshold =
			scheme = ""
			port =
			path = ""
			command =
			type = ""
		}

  }
  update_type =
  deploy_beta_enable =
  deploy_batch =
  deploy_exe_mode = ""
  deploy_wait_time =
  start_script = ""
  stop_script = ""
  incremental_deployment =
  jdk_name = ""
  jdk_version = ""
  agent_profile_list {
		agent_type = ""
		agent_version = ""

  }
  warmup_setting {
		enabled =
		warmup_time =
		curvature =
		enabled_protection =

  }
}
```

Import

tsf deploy_v_m_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_deploy_v_m_group.deploy_v_m_group deploy_v_m_group_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfDeployVMGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfDeployVMGroupCreate,
		Read:   resourceTencentCloudTsfDeployVMGroupRead,
		Delete: resourceTencentCloudTsfDeployVMGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "GroupId .",
			},

			"pkg_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Program package ID.",
			},

			"startup_parameters": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Start args of group.",
			},

			"deploy_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group description.",
			},

			"force_start": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow forced start.",
			},

			"enable_health_check": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable health check.",
			},

			"health_check_settings": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "When enabling health check, configure the health check settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"liveness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Survival health check. Note: This field may return null, indicating that no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Health check method. HTTP: check through HTTP interface; CMD: check through executing command; TCP: check through establishing TCP connection. Note: This field may return null, indicating that no valid value was found.",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The time delay for the container to start the health check. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum timeout period for each health check response. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The time interval for performing health checks. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of consecutive successful health checks required for the backend container to transition from failure to success. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of consecutive successful health checks required for the backend container to transition from success to failure. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The protocol used for HTTP health checks. HTTP and HTTPS are supported. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The port used for health checks, ranging from 1 to 65535. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The request path for HTTP health checks. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "The command to be executed for command health checks. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of readiness probe. TSF_DEFAULT represents the default readiness probe of TSF, while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified, the native readiness probe of Kubernetes is used by default. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"readiness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Readiness health check. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The health check method. HTTP indicates checking through an HTTP interface, CMD indicates checking through executing a command, and TCP indicates checking through establishing a TCP connection. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The time to delay the start of the container health check. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum timeout period for each health check response. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The time interval for performing health checks. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of consecutive successful health checks required for the backend container to transition from failure to success. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of consecutive successful health checks required for the backend container to transition from success to failure. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The protocol used for HTTP health checks. HTTP and HTTPS are supported. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The port used for health checks, ranging from 1 to 65535. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The request path for HTTP health checks. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "The command to be executed for command check. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of readiness probe. TSF_DEFAULT represents the default readiness probe of TSF, while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified, the native readiness probe of Kubernetes is used by default. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"update_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Update method: 0 for fast update, 1 for rolling update.",
			},

			"deploy_beta_enable": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable beta batch.",
			},

			"deploy_batch": {
				Optional:    true,
				ForceNew:    true,
				Description: "The ratio of instances participating in each batch during rolling release.",
			},

			"deploy_exe_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The execution method of rolling release.",
			},

			"deploy_wait_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The time interval for each batch during rolling release.",
			},

			"start_script": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The base64-encoded startup script.",
			},

			"stop_script": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The base64-encoded stop script.",
			},

			"incremental_deployment": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to perform incremental deployment. The default value is false, which means full update.",
			},

			"jdk_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "JDK name: konaJDK or openJDK.",
			},

			"jdk_version": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "JDK version : 8 or 11 (openJDK only support 8).",
			},

			"agent_profile_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Javaagent info: SERVICE_AGENT/OT_AGENT.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Agent type.",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Agent version.",
						},
					},
				},
			},

			"warmup_setting": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Warmup setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable preheating.",
						},
						"warmup_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Warmup time.",
						},
						"curvature": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Preheating curvature, with a value between 1 and 5.",
						},
						"enabled_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable preheating protection. If protection is enabled and more than 50% of nodes are in preheating state, preheating will be aborted.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfDeployVMGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_v_m_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewDeployGroupRequest()
		response = tsf.NewDeployGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pkg_id"); ok {
		request.PkgId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("startup_parameters"); ok {
		request.StartupParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deploy_desc"); ok {
		request.DeployDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("force_start"); ok {
		request.ForceStart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_health_check"); ok {
		request.EnableHealthCheck = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "health_check_settings"); ok {
		healthCheckSettings := tsf.HealthCheckSettings{}
		if livenessProbeMap, ok := helper.InterfaceToMap(dMap, "liveness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := livenessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["initial_delay_seconds"]; ok {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["timeout_seconds"]; ok {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["period_seconds"]; ok {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["success_threshold"]; ok {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["failure_threshold"]; ok {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["port"]; ok {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := livenessProbeMap["type"]; ok {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.LivenessProbe = &healthCheckSetting
		}
		if readinessProbeMap, ok := helper.InterfaceToMap(dMap, "readiness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := readinessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["initial_delay_seconds"]; ok {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["timeout_seconds"]; ok {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["period_seconds"]; ok {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["success_threshold"]; ok {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["failure_threshold"]; ok {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["port"]; ok {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := readinessProbeMap["type"]; ok {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.ReadinessProbe = &healthCheckSetting
		}
		request.HealthCheckSettings = &healthCheckSettings
	}

	if v, ok := d.GetOkExists("update_type"); ok {
		request.UpdateType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("deploy_beta_enable"); ok {
		request.DeployBetaEnable = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("deploy_batch"); v != nil {
	}

	if v, ok := d.GetOk("deploy_exe_mode"); ok {
		request.DeployExeMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deploy_wait_time"); ok {
		request.DeployWaitTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("start_script"); ok {
		request.StartScript = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stop_script"); ok {
		request.StopScript = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("incremental_deployment"); ok {
		request.IncrementalDeployment = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("jdk_name"); ok {
		request.JdkName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("jdk_version"); ok {
		request.JdkVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_profile_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			agentProfile := tsf.AgentProfile{}
			if v, ok := dMap["agent_type"]; ok {
				agentProfile.AgentType = helper.String(v.(string))
			}
			if v, ok := dMap["agent_version"]; ok {
				agentProfile.AgentVersion = helper.String(v.(string))
			}
			request.AgentProfileList = append(request.AgentProfileList, &agentProfile)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "warmup_setting"); ok {
		warmupSetting := tsf.WarmupSetting{}
		if v, ok := dMap["enabled"]; ok {
			warmupSetting.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["warmup_time"]; ok {
			warmupSetting.WarmupTime = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["curvature"]; ok {
			warmupSetting.Curvature = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["enabled_protection"]; ok {
			warmupSetting.EnabledProtection = helper.Bool(v.(bool))
		}
		request.WarmupSetting = &warmupSetting
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().DeployGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf deployVMGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(groupId)

	return resourceTencentCloudTsfDeployVMGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployVMGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_v_m_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	deployVMGroupId := d.Id()

	deployVMGroup, err := service.DescribeTsfDeployVMGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if deployVMGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfDeployVMGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if deployVMGroup.GroupId != nil {
		_ = d.Set("group_id", deployVMGroup.GroupId)
	}

	if deployVMGroup.PkgId != nil {
		_ = d.Set("pkg_id", deployVMGroup.PkgId)
	}

	if deployVMGroup.StartupParameters != nil {
		_ = d.Set("startup_parameters", deployVMGroup.StartupParameters)
	}

	if deployVMGroup.DeployDesc != nil {
		_ = d.Set("deploy_desc", deployVMGroup.DeployDesc)
	}

	if deployVMGroup.ForceStart != nil {
		_ = d.Set("force_start", deployVMGroup.ForceStart)
	}

	if deployVMGroup.EnableHealthCheck != nil {
		_ = d.Set("enable_health_check", deployVMGroup.EnableHealthCheck)
	}

	if deployVMGroup.HealthCheckSettings != nil {
		healthCheckSettingsMap := map[string]interface{}{}

		if deployVMGroup.HealthCheckSettings.LivenessProbe != nil {
			livenessProbeMap := map[string]interface{}{}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.ActionType != nil {
				livenessProbeMap["action_type"] = deployVMGroup.HealthCheckSettings.LivenessProbe.ActionType
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds != nil {
				livenessProbeMap["initial_delay_seconds"] = deployVMGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds != nil {
				livenessProbeMap["timeout_seconds"] = deployVMGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds != nil {
				livenessProbeMap["period_seconds"] = deployVMGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold != nil {
				livenessProbeMap["success_threshold"] = deployVMGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.FailureThreshold != nil {
				livenessProbeMap["failure_threshold"] = deployVMGroup.HealthCheckSettings.LivenessProbe.FailureThreshold
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.Scheme != nil {
				livenessProbeMap["scheme"] = deployVMGroup.HealthCheckSettings.LivenessProbe.Scheme
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.Port != nil {
				livenessProbeMap["port"] = deployVMGroup.HealthCheckSettings.LivenessProbe.Port
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.Path != nil {
				livenessProbeMap["path"] = deployVMGroup.HealthCheckSettings.LivenessProbe.Path
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.Command != nil {
				livenessProbeMap["command"] = deployVMGroup.HealthCheckSettings.LivenessProbe.Command
			}

			if deployVMGroup.HealthCheckSettings.LivenessProbe.Type != nil {
				livenessProbeMap["type"] = deployVMGroup.HealthCheckSettings.LivenessProbe.Type
			}

			healthCheckSettingsMap["liveness_probe"] = []interface{}{livenessProbeMap}
		}

		if deployVMGroup.HealthCheckSettings.ReadinessProbe != nil {
			readinessProbeMap := map[string]interface{}{}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.ActionType != nil {
				readinessProbeMap["action_type"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.ActionType
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds != nil {
				readinessProbeMap["initial_delay_seconds"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds != nil {
				readinessProbeMap["timeout_seconds"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds != nil {
				readinessProbeMap["period_seconds"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold != nil {
				readinessProbeMap["success_threshold"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold != nil {
				readinessProbeMap["failure_threshold"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.Scheme != nil {
				readinessProbeMap["scheme"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.Scheme
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.Port != nil {
				readinessProbeMap["port"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.Port
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.Path != nil {
				readinessProbeMap["path"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.Path
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.Command != nil {
				readinessProbeMap["command"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.Command
			}

			if deployVMGroup.HealthCheckSettings.ReadinessProbe.Type != nil {
				readinessProbeMap["type"] = deployVMGroup.HealthCheckSettings.ReadinessProbe.Type
			}

			healthCheckSettingsMap["readiness_probe"] = []interface{}{readinessProbeMap}
		}

		_ = d.Set("health_check_settings", []interface{}{healthCheckSettingsMap})
	}

	if deployVMGroup.UpdateType != nil {
		_ = d.Set("update_type", deployVMGroup.UpdateType)
	}

	if deployVMGroup.DeployBetaEnable != nil {
		_ = d.Set("deploy_beta_enable", deployVMGroup.DeployBetaEnable)
	}

	if deployVMGroup.DeployBatch != nil {
		_ = d.Set("deploy_batch", deployVMGroup.DeployBatch)
	}

	if deployVMGroup.DeployExeMode != nil {
		_ = d.Set("deploy_exe_mode", deployVMGroup.DeployExeMode)
	}

	if deployVMGroup.DeployWaitTime != nil {
		_ = d.Set("deploy_wait_time", deployVMGroup.DeployWaitTime)
	}

	if deployVMGroup.StartScript != nil {
		_ = d.Set("start_script", deployVMGroup.StartScript)
	}

	if deployVMGroup.StopScript != nil {
		_ = d.Set("stop_script", deployVMGroup.StopScript)
	}

	if deployVMGroup.IncrementalDeployment != nil {
		_ = d.Set("incremental_deployment", deployVMGroup.IncrementalDeployment)
	}

	if deployVMGroup.JdkName != nil {
		_ = d.Set("jdk_name", deployVMGroup.JdkName)
	}

	if deployVMGroup.JdkVersion != nil {
		_ = d.Set("jdk_version", deployVMGroup.JdkVersion)
	}

	if deployVMGroup.AgentProfileList != nil {
		agentProfileListList := []interface{}{}
		for _, agentProfileList := range deployVMGroup.AgentProfileList {
			agentProfileListMap := map[string]interface{}{}

			if deployVMGroup.AgentProfileList.AgentType != nil {
				agentProfileListMap["agent_type"] = deployVMGroup.AgentProfileList.AgentType
			}

			if deployVMGroup.AgentProfileList.AgentVersion != nil {
				agentProfileListMap["agent_version"] = deployVMGroup.AgentProfileList.AgentVersion
			}

			agentProfileListList = append(agentProfileListList, agentProfileListMap)
		}

		_ = d.Set("agent_profile_list", agentProfileListList)

	}

	if deployVMGroup.WarmupSetting != nil {
		warmupSettingMap := map[string]interface{}{}

		if deployVMGroup.WarmupSetting.Enabled != nil {
			warmupSettingMap["enabled"] = deployVMGroup.WarmupSetting.Enabled
		}

		if deployVMGroup.WarmupSetting.WarmupTime != nil {
			warmupSettingMap["warmup_time"] = deployVMGroup.WarmupSetting.WarmupTime
		}

		if deployVMGroup.WarmupSetting.Curvature != nil {
			warmupSettingMap["curvature"] = deployVMGroup.WarmupSetting.Curvature
		}

		if deployVMGroup.WarmupSetting.EnabledProtection != nil {
			warmupSettingMap["enabled_protection"] = deployVMGroup.WarmupSetting.EnabledProtection
		}

		_ = d.Set("warmup_setting", []interface{}{warmupSettingMap})
	}

	return nil
}

func resourceTencentCloudTsfDeployVMGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_v_m_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	deployVMGroupId := d.Id()

	if err := service.DeleteTsfDeployVMGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
