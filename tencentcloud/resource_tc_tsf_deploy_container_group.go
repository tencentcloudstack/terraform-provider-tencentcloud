/*
Provides a resource to create a tsf deploy_container_group

Example Usage

```hcl
resource "tencentcloud_tsf_deploy_container_group" "deploy_container_group" {
	group_id          = "group-yqml6w3a"
	cpu_request       = "0.25"
	mem_request       = "640"
	server            = "ccr.ccs.tencentyun.com"
	reponame          = "tsf_100011913960/terraform"
	tag_name          = "terraform-only-1"
	do_not_start      = false
	instance_num      = 1
	update_type       = 1
	update_ivl        = 10
	mem_limit         = "1280"
	cpu_limit         = "0.5"
	agent_cpu_request = "0.1"
	agent_cpu_limit   = "0.2"
	agent_mem_request = "125"
	agent_mem_limit   = "400"
	max_surge         = "25%"
	max_unavailable   = "0"
	service_setting {
		access_type = 1
		protocol_ports {
			protocol    = "TCP"
			port        = 18081
			target_port = 18081
			node_port   = 30001
		}
		subnet_id						 = ""
		disable_service                  = false
		headless_service                 = false
		allow_delete_service             = true
		open_session_affinity            = false
		session_affinity_timeout_seconds = 10800

	}
	health_check_settings {
		readiness_probe {
			action_type           = "TCP"
			initial_delay_seconds = 0
			timeout_seconds       = 3
			period_seconds        = 30
			success_threshold     = 1
			failure_threshold     = 3
			scheme                = "HTTP"
			port                  = 80
			path                  = "/"
			type                  = "TSF_DEFAULT"
		}
	}
	scheduling_strategy {
		type = "NONE"
	}
	deploy_agent = true
	repo_type = "personal"
	volume_clean = false
	jvm_opts          = "-Xms128m -Xmx512m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=512m"
	warmup_setting {
		enabled = false
	}
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfDeployContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfDeployContainerGroupCreate,
		Read:   resourceTencentCloudTsfDeployContainerGroupRead,
		Update: resourceTencentCloudTsfDeployContainerGroupUpdate,
		Delete: resourceTencentCloudTsfDeployContainerGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "group Id.",
			},

			"tag_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "image version name, v1.",
			},

			"instance_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "instance number.",
			},

			"server": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "image server.",
			},

			"reponame": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "old image name, eg: /tsf/server.",
			},

			"cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum number of CPU cores for the business container, corresponding to the limit in K8S. If not specified, it defaults to twice the request.",
			},

			"mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum memory size in MiB for the business container, corresponding to the limit in K8S. If not specified, it defaults to twice the request.",
			},

			"jvm_opts": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "jvm options.",
			},

			"cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The number of CPU cores allocated to the business container, corresponding to the request in K8S. The default value is 0.25.",
			},

			"mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The amount of memory in MiB allocated to the business container, corresponding to the request in K8S. The default value is 640 MiB.",
			},

			"do_not_start": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Not start right away.",
			},

			"repo_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "(Priority use) New image name, such as /tsf/nginx.",
			},

			"update_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Update method: 0 for fast update, 1 for rolling update.",
			},

			"update_ivl": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "update Interval, is required when rolling update.",
			},

			"agent_cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The number of CPU cores allocated to the agent container corresponds to the request field in Kubernetes.",
			},

			"agent_cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum number of CPU cores allocated to the agent container corresponds to the limit field in Kubernetes.",
			},

			"agent_mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.",
			},

			"agent_mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum amount of memory in MiB allocated to the agent container corresponds to the &amp;#39;limit&amp;#39; field in Kubernetes.",
			},

			"istio_cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The number of CPU cores allocated to the istio proxy container corresponds to the &amp;#39;request&amp;#39; field in Kubernetes.",
			},

			"istio_cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum amount of CPU cores allocated to the istio proxy container corresponds to the &amp;#39;limit&amp;#39; field in Kubernetes.",
			},

			"istio_mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.",
			},

			"istio_mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The maximum amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.",
			},

			"max_surge": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MaxSurge parameter in Kubernetes rolling update strategy.",
			},

			"max_unavailable": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MaxUnavailable parameter in Kubernetes rolling update strategy.",
			},

			"health_check_settings": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The configuration information for health checks. If this parameter is not specified, the health check is not set by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"liveness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Liveness probe. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The health check method. HTTP: checks through an HTTP interface; CMD: checks by executing a command; TCP: checks by establishing a TCP connection. Note: This field may return null, indicating that no valid values can be obtained.",
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

			"envs": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The environment variables that the application runs in the deployment group. If this parameter is not specified, no additional environment variables are set by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "env param name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "value of env.",
						},
						"value_from": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Kubernetes ValueFrom configuration. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_ref": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "The FieldRef configuration of Kubernetes env. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The FieldPath configuration of Kubernetes. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"resource_field_ref": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "The ResourceFieldRef configuration of Kubernetes env. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The Resource configuration of Kubernetes. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"service_setting": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Network settings for container deployment groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "0: Public network, 1: Access within the cluster, 2: NodePort, 3: Access within VPC. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"protocol_ports": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Container port mapping. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "TCP or UDP.",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "port.",
									},
									"target_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "container port.",
									},
									"node_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "node port.",
									},
								},
							},
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "subnet Id.",
						},
						"disable_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to create a Kubernetes service. The default value is false. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"headless_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the service is of headless type. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"allow_delete_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When set to true and DisableService is also true, the previously created service will be deleted. Please use with caution. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"open_session_affinity": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable session affinity. true means enabled, false means disabled. The default value is false. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"session_affinity_timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Session affinity session time. The default value is 10800. Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"deploy_agent": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to deploy the agent container. If this parameter is not specified, the agent container will not be deployed by default.",
			},

			"scheduling_strategy": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Node scheduling strategy. If this parameter is not specified, the node scheduling strategy will not be used by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NONE: Do not use scheduling strategy; CROSS_AZ: Deploy across availability zones. Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"incremental_deployment": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to perform incremental deployment. The default value is false, which means full update.",
			},

			"repo_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "repo type, tcr or leave it blank.",
			},

			"volume_info_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Volume information, as a list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume type.",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume name.",
						},
						"volume_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "volume config.",
						},
					},
				},
			},

			"volume_mount_info_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Volume mount point information, list type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_mount_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mount volume name.",
						},
						"volume_mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mount path.",
						},
						"volume_mount_sub_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "mount subPath.",
						},
						"read_or_write": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Read and write access mode. 1: Read-only. 2: Read-write.",
						},
					},
				},
			},

			"volume_clean": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to clear the volume information. Default is false.",
			},

			"agent_profile_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "javaagent info: SERVICE_AGENT/OT_AGENT.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent type.",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent version.",
						},
					},
				},
			},

			"warmup_setting": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "warmup setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to enable preheating.",
						},
						"warmup_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "warmup time.",
						},
						"curvature": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Preheating curvature, with a value between 1 and 5.",
						},
						"enabled_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to enable preheating protection. If protection is enabled and more than 50% of nodes are in preheating state, preheating will be aborted.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfDeployContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_container_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request  = tsf.NewDeployContainerGroupRequest()
		response = tsf.NewDeployContainerGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_name"); ok {
		request.TagName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("instance_num"); v != nil {
		request.InstanceNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("server"); ok {
		request.Server = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reponame"); ok {
		request.Reponame = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu_limit"); ok {
		request.CpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_limit"); ok {
		request.MemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("jvm_opts"); ok {
		request.JvmOpts = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu_request"); ok {
		request.CpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_request"); ok {
		request.MemRequest = helper.String(v.(string))
	}

	if v, _ := d.GetOk("do_not_start"); v != nil {
		request.DoNotStart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("repo_name"); ok {
		request.RepoName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("update_type"); v != nil {
		request.UpdateType = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("update_ivl"); v != nil {
		request.UpdateIvl = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("agent_cpu_request"); ok {
		request.AgentCpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_cpu_limit"); ok {
		request.AgentCpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_mem_request"); ok {
		request.AgentMemRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_mem_limit"); ok {
		request.AgentMemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_cpu_request"); ok {
		request.IstioCpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_cpu_limit"); ok {
		request.IstioCpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_mem_request"); ok {
		request.IstioMemRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_mem_limit"); ok {
		request.IstioMemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_surge"); ok {
		request.MaxSurge = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_unavailable"); ok {
		request.MaxUnavailable = helper.String(v.(string))
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

	if v, ok := d.GetOk("envs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			env := tsf.Env{}
			if v, ok := dMap["name"]; ok {
				env.Name = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				env.Value = helper.String(v.(string))
			}
			if valueFromMap, ok := helper.InterfaceToMap(dMap, "value_from"); ok {
				valueFrom := tsf.ValueFrom{}
				if fieldRefMap, ok := helper.InterfaceToMap(valueFromMap, "field_ref"); ok {
					fieldRef := tsf.FieldRef{}
					if v, ok := fieldRefMap["field_path"]; ok {
						fieldRef.FieldPath = helper.String(v.(string))
					}
					valueFrom.FieldRef = &fieldRef
				}
				if resourceFieldRefMap, ok := helper.InterfaceToMap(valueFromMap, "resource_field_ref"); ok {
					resourceFieldRef := tsf.ResourceFieldRef{}
					if v, ok := resourceFieldRefMap["resource"]; ok {
						resourceFieldRef.Resource = helper.String(v.(string))
					}
					valueFrom.ResourceFieldRef = &resourceFieldRef
				}
				env.ValueFrom = &valueFrom
			}
			request.Envs = append(request.Envs, &env)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_setting"); ok {
		serviceSetting := tsf.ServiceSetting{}
		if v, ok := dMap["access_type"]; ok {
			serviceSetting.AccessType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["protocol_ports"]; ok {
			for _, item := range v.([]interface{}) {
				protocolPortsMap := item.(map[string]interface{})
				protocolPort := tsf.ProtocolPort{}
				if v, ok := protocolPortsMap["protocol"]; ok {
					protocolPort.Protocol = helper.String(v.(string))
				}
				if v, ok := protocolPortsMap["port"]; ok {
					protocolPort.Port = helper.IntInt64(v.(int))
				}
				if v, ok := protocolPortsMap["target_port"]; ok {
					protocolPort.TargetPort = helper.IntInt64(v.(int))
				}
				if v, ok := protocolPortsMap["node_port"]; ok {
					protocolPort.NodePort = helper.IntInt64(v.(int))
				}
				serviceSetting.ProtocolPorts = append(serviceSetting.ProtocolPorts, &protocolPort)
			}
		}
		if v, ok := dMap["subnet_id"]; ok {
			serviceSetting.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["disable_service"]; ok {
			serviceSetting.DisableService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["headless_service"]; ok {
			serviceSetting.HeadlessService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["allow_delete_service"]; ok {
			serviceSetting.AllowDeleteService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["open_session_affinity"]; ok {
			serviceSetting.OpenSessionAffinity = helper.Bool(v.(bool))
		}
		if v, ok := dMap["session_affinity_timeout_seconds"]; ok {
			serviceSetting.SessionAffinityTimeoutSeconds = helper.IntInt64(v.(int))
		}
		request.ServiceSetting = &serviceSetting
	}

	if v, _ := d.GetOk("deploy_agent"); v != nil {
		request.DeployAgent = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "scheduling_strategy"); ok {
		schedulingStrategy := tsf.SchedulingStrategy{}
		if v, ok := dMap["type"]; ok {
			schedulingStrategy.Type = helper.String(v.(string))
		}
		request.SchedulingStrategy = &schedulingStrategy
	}

	if v, _ := d.GetOk("incremental_deployment"); v != nil {
		request.IncrementalDeployment = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("repo_type"); ok {
		request.RepoType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("volume_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			volumeInfo := tsf.VolumeInfo{}
			if v, ok := dMap["volume_type"]; ok {
				volumeInfo.VolumeType = helper.String(v.(string))
			}
			if v, ok := dMap["volume_name"]; ok {
				volumeInfo.VolumeName = helper.String(v.(string))
			}
			if v, ok := dMap["volume_config"]; ok {
				volumeInfo.VolumeConfig = helper.String(v.(string))
			}
			request.VolumeInfoList = append(request.VolumeInfoList, &volumeInfo)
		}
	}

	if v, ok := d.GetOk("volume_mount_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			volumeMountInfo := tsf.VolumeMountInfo{}
			if v, ok := dMap["volume_mount_name"]; ok {
				volumeMountInfo.VolumeMountName = helper.String(v.(string))
			}
			if v, ok := dMap["volume_mount_path"]; ok {
				volumeMountInfo.VolumeMountPath = helper.String(v.(string))
			}
			if v, ok := dMap["volume_mount_sub_path"]; ok {
				volumeMountInfo.VolumeMountSubPath = helper.String(v.(string))
			}
			if v, ok := dMap["read_or_write"]; ok {
				volumeMountInfo.ReadOrWrite = helper.String(v.(string))
			}
			request.VolumeMountInfoList = append(request.VolumeMountInfoList, &volumeMountInfo)
		}
	}

	if v, _ := d.GetOk("volume_clean"); v != nil {
		request.VolumeClean = helper.Bool(v.(bool))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().DeployContainerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s deploy tsf deployContainerGroup failed, reason:%+v", logId, err)
		return err
	}

	if !*response.Response.Result {
		return fmt.Errorf("[CRITAL]%s deploy tsf deployContainerGroup failed", logId)
	}
	d.SetId(groupId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		groupInfo, err := service.DescribeTsfStartContainerGroupById(ctx, groupId)
		if err != nil {
			return retryError(err)
		}
		if groupInfo == nil {
			err = fmt.Errorf("group %s not exists", groupId)
			return resource.NonRetryableError(err)
		}
		if *groupInfo.Status == "Running" {
			return nil
		}
		if *groupInfo.Status == "Waiting" || *groupInfo.Status == "Updating" {
			return resource.RetryableError(fmt.Errorf("deploy container group status is %s", *groupInfo.Status))
		}
		err = fmt.Errorf("deploy container group status is %v, we won't wait for it finish", *groupInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s deploy container group, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudTsfDeployContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_container_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()

	deployContainerGroup, err := service.DescribeTsfDeployContainerGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if deployContainerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfDeployContainerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if deployContainerGroup.GroupId != nil {
		_ = d.Set("group_id", deployContainerGroup.GroupId)
	}

	if deployContainerGroup.TagName != nil {
		_ = d.Set("tag_name", deployContainerGroup.TagName)
	}

	if deployContainerGroup.InstanceNum != nil {
		_ = d.Set("instance_num", deployContainerGroup.InstanceNum)
	}

	if deployContainerGroup.Server != nil {
		_ = d.Set("server", deployContainerGroup.Server)
	}

	if deployContainerGroup.Reponame != nil {
		_ = d.Set("reponame", deployContainerGroup.Reponame)
	}

	if deployContainerGroup.CpuLimit != nil {
		_ = d.Set("cpu_limit", deployContainerGroup.CpuLimit)
	}

	if deployContainerGroup.MemLimit != nil {
		_ = d.Set("mem_limit", deployContainerGroup.MemLimit)
	}

	if deployContainerGroup.JvmOpts != nil {
		_ = d.Set("jvm_opts", deployContainerGroup.JvmOpts)
	}

	if deployContainerGroup.CpuRequest != nil {
		_ = d.Set("cpu_request", deployContainerGroup.CpuRequest)
	}

	if deployContainerGroup.MemRequest != nil {
		_ = d.Set("mem_request", deployContainerGroup.MemRequest)
	}

	// if deployContainerGroup.DoNotStart != nil {
	// 	_ = d.Set("do_not_start", deployContainerGroup.DoNotStart)
	// }

	// if deployContainerGroup.RepoName != nil {
	// 	_ = d.Set("repo_name", deployContainerGroup.RepoName)
	// }

	if deployContainerGroup.UpdateType != nil {
		_ = d.Set("update_type", deployContainerGroup.UpdateType)
	}

	if deployContainerGroup.UpdateIvl != nil {
		_ = d.Set("update_ivl", deployContainerGroup.UpdateIvl)
	}

	if deployContainerGroup.AgentCpuRequest != nil {
		_ = d.Set("agent_cpu_request", deployContainerGroup.AgentCpuRequest)
	}

	if deployContainerGroup.AgentCpuLimit != nil {
		_ = d.Set("agent_cpu_limit", deployContainerGroup.AgentCpuLimit)
	}

	if deployContainerGroup.AgentMemRequest != nil {
		_ = d.Set("agent_mem_request", deployContainerGroup.AgentMemRequest)
	}

	if deployContainerGroup.AgentMemLimit != nil {
		_ = d.Set("agent_mem_limit", deployContainerGroup.AgentMemLimit)
	}

	if deployContainerGroup.IstioCpuRequest != nil {
		_ = d.Set("istio_cpu_request", deployContainerGroup.IstioCpuRequest)
	}

	if deployContainerGroup.IstioCpuLimit != nil {
		_ = d.Set("istio_cpu_limit", deployContainerGroup.IstioCpuLimit)
	}

	if deployContainerGroup.IstioMemRequest != nil {
		_ = d.Set("istio_mem_request", deployContainerGroup.IstioMemRequest)
	}

	if deployContainerGroup.IstioMemLimit != nil {
		_ = d.Set("istio_mem_limit", deployContainerGroup.IstioMemLimit)
	}

	// if deployContainerGroup.HealthCheckSettings != nil {
	// 	healthCheckSettingsMap := map[string]interface{}{}

	// 	if deployContainerGroup.HealthCheckSettings.LivenessProbe != nil {
	// 		livenessProbeMap := map[string]interface{}{}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.ActionType != nil {
	// 			livenessProbeMap["action_type"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.ActionType
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds != nil {
	// 			livenessProbeMap["initial_delay_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds != nil {
	// 			livenessProbeMap["timeout_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds != nil {
	// 			livenessProbeMap["period_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold != nil {
	// 			livenessProbeMap["success_threshold"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.FailureThreshold != nil {
	// 			livenessProbeMap["failure_threshold"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.FailureThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Scheme != nil {
	// 			livenessProbeMap["scheme"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Scheme
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Port != nil {
	// 			livenessProbeMap["port"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Port
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Path != nil {
	// 			livenessProbeMap["path"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Path
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Command != nil {
	// 			livenessProbeMap["command"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Command
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Type != nil {
	// 			livenessProbeMap["type"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Type
	// 		}

	// 		healthCheckSettingsMap["liveness_probe"] = []interface{}{livenessProbeMap}
	// 	}

	// 	if deployContainerGroup.HealthCheckSettings.ReadinessProbe != nil {
	// 		readinessProbeMap := map[string]interface{}{}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.ActionType != nil {
	// 			readinessProbeMap["action_type"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.ActionType
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds != nil {
	// 			readinessProbeMap["initial_delay_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds != nil {
	// 			readinessProbeMap["timeout_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds != nil {
	// 			readinessProbeMap["period_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold != nil {
	// 			readinessProbeMap["success_threshold"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold != nil {
	// 			readinessProbeMap["failure_threshold"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Scheme != nil {
	// 			readinessProbeMap["scheme"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Scheme
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Port != nil {
	// 			readinessProbeMap["port"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Port
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Path != nil {
	// 			readinessProbeMap["path"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Path
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Command != nil {
	// 			readinessProbeMap["command"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Command
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Type != nil {
	// 			readinessProbeMap["type"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Type
	// 		}

	// 		healthCheckSettingsMap["readiness_probe"] = []interface{}{readinessProbeMap}
	// 	}

	// 	_ = d.Set("health_check_settings", []interface{}{healthCheckSettingsMap})
	// }

	if deployContainerGroup.Envs != nil {
		envsList := []interface{}{}
		for _, envs := range deployContainerGroup.Envs {
			envsMap := map[string]interface{}{}

			if envs.Name != nil {
				envsMap["name"] = envs.Name
			}

			if envs.Value != nil {
				envsMap["value"] = envs.Value
			}

			if envs.ValueFrom != nil {
				valueFromMap := map[string]interface{}{}

				if envs.ValueFrom.FieldRef != nil {
					fieldRefMap := map[string]interface{}{}

					if envs.ValueFrom.FieldRef.FieldPath != nil {
						fieldRefMap["field_path"] = envs.ValueFrom.FieldRef.FieldPath
					}

					valueFromMap["field_ref"] = []interface{}{fieldRefMap}
				}

				if envs.ValueFrom.ResourceFieldRef != nil {
					resourceFieldRefMap := map[string]interface{}{}

					if envs.ValueFrom.ResourceFieldRef.Resource != nil {
						resourceFieldRefMap["resource"] = envs.ValueFrom.ResourceFieldRef.Resource
					}

					valueFromMap["resource_field_ref"] = []interface{}{resourceFieldRefMap}
				}

				envsMap["value_from"] = []interface{}{valueFromMap}
			}

			envsList = append(envsList, envsMap)
		}

		_ = d.Set("envs", envsList)

	}

	if deployContainerGroup.DeployAgent != nil {
		_ = d.Set("deploy_agent", deployContainerGroup.DeployAgent)
	}

	if deployContainerGroup.RepoType != nil {
		_ = d.Set("repo_type", deployContainerGroup.RepoType)
	}

	if deployContainerGroup.VolumeInfos != nil {
		volumeInfoListList := []interface{}{}
		for _, volumeInfoList := range deployContainerGroup.VolumeInfos {
			volumeInfoListMap := map[string]interface{}{}

			if volumeInfoList.VolumeType != nil {
				volumeInfoListMap["volume_type"] = volumeInfoList.VolumeType
			}

			if volumeInfoList.VolumeName != nil {
				volumeInfoListMap["volume_name"] = volumeInfoList.VolumeName
			}

			if volumeInfoList.VolumeConfig != nil {
				volumeInfoListMap["volume_config"] = volumeInfoList.VolumeConfig
			}

			volumeInfoListList = append(volumeInfoListList, volumeInfoListMap)
		}

		_ = d.Set("volume_info_list", volumeInfoListList)

	}

	if deployContainerGroup.VolumeMountInfos != nil {
		volumeMountInfoListList := []interface{}{}
		for _, volumeMountInfoList := range deployContainerGroup.VolumeMountInfos {
			volumeMountInfoListMap := map[string]interface{}{}

			if volumeMountInfoList.VolumeMountName != nil {
				volumeMountInfoListMap["volume_mount_name"] = volumeMountInfoList.VolumeMountName
			}

			if volumeMountInfoList.VolumeMountPath != nil {
				volumeMountInfoListMap["volume_mount_path"] = volumeMountInfoList.VolumeMountPath
			}

			if volumeMountInfoList.VolumeMountSubPath != nil {
				volumeMountInfoListMap["volume_mount_sub_path"] = volumeMountInfoList.VolumeMountSubPath
			}

			if volumeMountInfoList.ReadOrWrite != nil {
				volumeMountInfoListMap["read_or_write"] = volumeMountInfoList.ReadOrWrite
			}

			volumeMountInfoListList = append(volumeMountInfoListList, volumeMountInfoListMap)
		}

		_ = d.Set("volume_mount_info_list", volumeMountInfoListList)

	}

	// if deployContainerGroup.VolumeClean != nil {
	// 	_ = d.Set("volume_clean", deployContainerGroup.VolumeClean)
	// }

	// if deployContainerGroup.AgentProfileList != nil {
	// 	agentProfileListList := []interface{}{}
	// 	for _, agentProfileList := range deployContainerGroup.AgentProfileList {
	// 		agentProfileListMap := map[string]interface{}{}

	// 		if deployContainerGroup.AgentProfileList.AgentType != nil {
	// 			agentProfileListMap["agent_type"] = deployContainerGroup.AgentProfileList.AgentType
	// 		}

	// 		if deployContainerGroup.AgentProfileList.AgentVersion != nil {
	// 			agentProfileListMap["agent_version"] = deployContainerGroup.AgentProfileList.AgentVersion
	// 		}

	// 		agentProfileListList = append(agentProfileListList, agentProfileListMap)
	// 	}

	// 	_ = d.Set("agent_profile_list", agentProfileListList)

	// }

	if deployContainerGroup.WarmupSetting != nil {
		warmupSettingMap := map[string]interface{}{}

		if deployContainerGroup.WarmupSetting.Enabled != nil {
			warmupSettingMap["enabled"] = deployContainerGroup.WarmupSetting.Enabled
		}

		if deployContainerGroup.WarmupSetting.WarmupTime != nil {
			warmupSettingMap["warmup_time"] = deployContainerGroup.WarmupSetting.WarmupTime
		}

		if deployContainerGroup.WarmupSetting.Curvature != nil {
			warmupSettingMap["curvature"] = deployContainerGroup.WarmupSetting.Curvature
		}

		if deployContainerGroup.WarmupSetting.EnabledProtection != nil {
			warmupSettingMap["enabled_protection"] = deployContainerGroup.WarmupSetting.EnabledProtection
		}

		_ = d.Set("warmup_setting", []interface{}{warmupSettingMap})
	}

	return nil
}

func resourceTencentCloudTsfDeployContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_container_group.update")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)

	request := tsf.NewModifyContainerReplicasRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_id", "tag_name", "server", "reponame", "cpu_limit", "mem_limit", "jvm_opts", "cpu_request", "mem_request", "do_not_start", "repo_name", "update_type", "update_ivl", "agent_cpu_request", "agent_cpu_limit", "agent_mem_request", "agent_mem_limit", "istio_cpu_request", "istio_cpu_limit", "istio_mem_request", "istio_mem_limit", "max_surge", "max_unavailable", "health_check_settings", "envs", "service_setting", "deploy_agent", "scheduling_strategy", "incremental_deployment", "repo_type", "volume_infos", "volume_mount_infos", "volume_info_list", "volume_mount_info_list", "volume_clean", "agent_profile_list", "warmup_setting"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_num") {
		if v, ok := d.GetOk("instance_num"); ok {
			request.InstanceNum = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyContainerReplicas(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf unitRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfDeployContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_deploy_container_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
