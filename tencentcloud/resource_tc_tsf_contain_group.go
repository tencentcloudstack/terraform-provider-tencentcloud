/*
Provides a resource to create a tsf contain_group

Example Usage

```hcl
resource "tencentcloud_tsf_contain_group" "contain_group" {
    access_type           = 0
    application_id        = "application-y5r4nejv"
    cluster_id            = "cls-2yu5kxr8"
    cpu_limit             = "0.5"
    cpu_request           = "0.25"
    group_name            = "terraform-test"
    group_resource_type   = "DEF"
    instance_num          = 1
    mem_limit             = "1280"
    mem_request           = "640"
    namespace_id          = "namespace-ydlezgxa"
    update_ivl            = 10
    update_type           = 1

    protocol_ports {
        node_port   = 0
        port        = 333
        protocol    = "TCP"
        target_port = 333
    }
}
```

Import

tsf contain_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_contain_group.contain_group contain_group_id
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

func resourceTencentCloudTsfContainGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfContainGroupCreate,
		Read:   resourceTencentCloudTsfContainGroupRead,
		Update: resourceTencentCloudTsfContainGroupUpdate,
		Delete: resourceTencentCloudTsfContainGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The application ID to which the group belongs.",
			},

			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the namespace to which the group belongs.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group name field, length 1~60, beginning with a letter or underscore, can contain alphanumeric underscore.",
			},

			"instance_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "number of instances.",
			},

			"access_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "0: public network 1: access within the cluster 2: NodePort.",
			},

			"protocol_ports": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Protocol Ports array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TCP UDP.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "service port.",
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
							Description: "host port.",
						},
					},
				},
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cpu_limit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The maximum number of allocated CPU cores, corresponding to the K8S limit.",
			},

			"mem_limit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maximum allocated memory MiB, corresponding to K8S limit.",
			},

			"group_comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group remarks field, the length should not exceed 200 characters.",
			},

			"update_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Update method: 0: fast update 1: rolling update.",
			},

			"update_ivl": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Rolling update is required, update interval.",
			},

			"cpu_request": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Initially allocated CPU cores, corresponding to K8S request.",
			},

			"mem_request": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Initially allocated memory MiB, corresponding to K8S request.",
			},

			"group_resource_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment Group Resource Type.",
			},

			"subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "subnet ID.",
			},

			"agent_cpu_request": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The number of CPU cores allocated by the agent container, corresponding to the K8S request.",
			},

			"agent_cpu_limit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The maximum number of CPU cores for the agent container, corresponding to the limit of K8S.",
			},

			"agent_mem_request": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The number of memory MiB allocated by the agent container, corresponding to the K8S request.",
			},

			"agent_mem_limit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The maximum memory MiB of the agent container, corresponding to the limit of K8S.",
			},

			"istio_cpu_request": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The number of CPU cores allocated by the istioproxy container, corresponding to the K8S request.",
			},

			"istio_cpu_limit": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The maximum number of CPU cores for the istioproxy container corresponds to the limit of K8S.",
			},

			"istio_mem_request": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The number of memory MiB allocated by the istioproxy container, corresponding to the K8S request.",
			},

			"istio_mem_limit": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The maximum memory MiB of the istioproxy container corresponds to the limit of K8S.",
			},

			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group ID.",
			},

			"current_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of instances launched.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "creation time.",
			},

			"server": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "mirror server.",
			},

			"reponame": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Mirror name, such as /tsf/nginx.",
			},

			"tag_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Image version name.",
			},

			"cluster_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster name.",
			},

			"namespace_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"lb_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "load balancing ip.",
			},

			"application_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "App types.",
			},

			"cluster_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service ip.",
			},

			"envs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "environment variable array object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "environment variable name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "environment variable value.",
						},
						"value_from": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "k8s ValueFrom.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_ref": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "FieldRef for k8s env.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "FieldPath of k8s.",
												},
											},
										},
									},
									"resource_field_ref": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ResourceFieldRef of k8s env.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource of k8s.",
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

			"application_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application Name.",
			},

			"message": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "pod error message description.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group status.",
			},

			"microservice_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service type.",
			},

			"instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of deployment group instances.",
			},

			"updated_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Deployment group update timestamp.",
			},

			"max_surge": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The MaxSurge parameter of the kubernetes rolling update policy.",
			},

			"max_unavailable": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The MaxUnavailable parameter of the kubernetes rolling update policy.",
			},

			"health_check_settings": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Deployment group health check settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"liveness_probe": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "live health check.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "health check method. HTTP: check by HTTP interface; CMD: check by executing command; TCP: check by establishing TCP connection.",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The time for the container to delay starting the health check.",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum timeout for each health check response.",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Interval between health checks.",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the number of consecutive health check successes for the backend container from failure to success.",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the number of consecutive health check successes of the backend container from success to failure.",
									},
									"scheme": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The inspection protocol used by the HTTP health check method. HTTP and HTTPS are supported.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health check port, range 1~65535.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The request path of the HTTP health check interface.",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Execute command check mode, the command to execute.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TSF_DEFAULT: tsf default readiness probe. K8S_NATIVE: k8s native probe. If not filled, it defaults to k8s native probe.",
									},
								},
							},
						},
						"readiness_probe": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "readiness health check.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "health check method. HTTP: check by HTTP interface; CMD: check by executing command; TCP: check by establishing TCP connection.",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The time for the container to delay starting the health check.",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum timeout for each health check response.",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The interval at which health checks are performed.",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the number of consecutive health check successes for the backend container from failure to success.",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the number of consecutive health check successes for the backend container from success to failure.",
									},
									"scheme": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The inspection protocol used by the HTTP health check method. HTTP and HTTPS are supported.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health check port, range 1~65535.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The request path of the HTTP health check interface.",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Execute command check mode, the command to execute.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TSF_DEFAULT: tsf default readiness probe. K8S_NATIVE: k8s native probe. If not filled, it defaults to k8s native probe.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfContainGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_contain_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateContainGroupRequest()
		response = tsf.NewCreateContainGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("application_id"); ok {
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("instance_num"); v != nil {
		request.InstanceNum = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("access_type"); v != nil {
		request.AccessType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("protocol_ports"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			protocolPort := tsf.ProtocolPort{}
			if v, ok := dMap["protocol"]; ok {
				protocolPort.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				protocolPort.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["target_port"]; ok {
				protocolPort.TargetPort = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["node_port"]; ok && v.(int) > 0 {
				protocolPort.NodePort = helper.IntInt64(v.(int))
			}
			request.ProtocolPorts = append(request.ProtocolPorts, &protocolPort)
		}
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu_limit"); ok {
		request.CpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_limit"); ok {
		request.MemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_comment"); ok {
		request.GroupComment = helper.String(v.(string))
	}

	if v, _ := d.GetOk("update_type"); v != nil {
		request.UpdateType = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("update_ivl"); v != nil {
		request.UpdateIvl = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu_request"); ok {
		request.CpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_request"); ok {
		request.MemRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_resource_type"); ok {
		request.GroupResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateContainGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf containGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.Result
	d.SetId(groupId)

	return resourceTencentCloudTsfContainGroupRead(d, meta)
}

func resourceTencentCloudTsfContainGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_contain_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()

	containGroup, err := service.DescribeTsfContainGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if containGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfContainGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if containGroup.ApplicationId != nil {
		_ = d.Set("application_id", containGroup.ApplicationId)
	}

	if containGroup.NamespaceId != nil {
		_ = d.Set("namespace_id", containGroup.NamespaceId)
	}

	if containGroup.GroupName != nil {
		_ = d.Set("group_name", containGroup.GroupName)
	}

	if containGroup.InstanceNum != nil {
		_ = d.Set("instance_num", containGroup.InstanceNum)
	}

	if containGroup.AccessType != nil {
		_ = d.Set("access_type", containGroup.AccessType)
	}

	if containGroup.ProtocolPorts != nil {
		protocolPortsList := []interface{}{}
		for _, protocolPorts := range containGroup.ProtocolPorts {
			protocolPortsMap := map[string]interface{}{}

			if protocolPorts.Protocol != nil {
				protocolPortsMap["protocol"] = protocolPorts.Protocol
			}

			if protocolPorts.Port != nil {
				protocolPortsMap["port"] = protocolPorts.Port
			}

			if protocolPorts.TargetPort != nil {
				protocolPortsMap["target_port"] = protocolPorts.TargetPort
			}

			if protocolPorts.NodePort != nil {
				protocolPortsMap["node_port"] = protocolPorts.NodePort
			}

			protocolPortsList = append(protocolPortsList, protocolPortsMap)
		}

		_ = d.Set("protocol_ports", protocolPortsList)

	}

	if containGroup.ClusterId != nil {
		_ = d.Set("cluster_id", containGroup.ClusterId)
	}

	if containGroup.CpuLimit != nil {
		_ = d.Set("cpu_limit", containGroup.CpuLimit)
	}

	if containGroup.MemLimit != nil {
		_ = d.Set("mem_limit", containGroup.MemLimit)
	}

	// if containGroup.GroupComment != nil {
	// 	_ = d.Set("group_comment", containGroup.GroupComment)
	// }

	if containGroup.UpdateType != nil {
		_ = d.Set("update_type", containGroup.UpdateType)
	}

	if containGroup.UpdateIvl != nil {
		_ = d.Set("update_ivl", containGroup.UpdateIvl)
	}

	if containGroup.CpuRequest != nil {
		_ = d.Set("cpu_request", containGroup.CpuRequest)
	}

	if containGroup.MemRequest != nil {
		_ = d.Set("mem_request", containGroup.MemRequest)
	}

	if containGroup.GroupResourceType != nil {
		_ = d.Set("group_resource_type", containGroup.GroupResourceType)
	}

	if containGroup.SubnetId != nil {
		_ = d.Set("subnet_id", containGroup.SubnetId)
	}

	// if containGroup.AgentCpuRequest != nil {
	// 	_ = d.Set("agent_cpu_request", containGroup.AgentCpuRequest)
	// }

	// if containGroup.AgentCpuLimit != nil {
	// 	_ = d.Set("agent_cpu_limit", containGroup.AgentCpuLimit)
	// }

	// if containGroup.AgentMemRequest != nil {
	// 	_ = d.Set("agent_mem_request", containGroup.AgentMemRequest)
	// }

	// if containGroup.AgentMemLimit != nil {
	// 	_ = d.Set("agent_mem_limit", containGroup.AgentMemLimit)
	// }

	// if containGroup.IstioCpuRequest != nil {
	// 	_ = d.Set("istio_cpu_request", containGroup.IstioCpuRequest)
	// }

	// if containGroup.IstioCpuLimit != nil {
	// 	_ = d.Set("istio_cpu_limit", containGroup.IstioCpuLimit)
	// }

	// if containGroup.IstioMemRequest != nil {
	// 	_ = d.Set("istio_mem_request", containGroup.IstioMemRequest)
	// }

	// if containGroup.IstioMemLimit != nil {
	// 	_ = d.Set("istio_mem_limit", containGroup.IstioMemLimit)
	// }

	if containGroup.GroupId != nil {
		_ = d.Set("group_id", containGroup.GroupId)
	}

	if containGroup.CurrentNum != nil {
		_ = d.Set("current_num", containGroup.CurrentNum)
	}

	if containGroup.CreateTime != nil {
		_ = d.Set("create_time", containGroup.CreateTime)
	}

	if containGroup.Server != nil {
		_ = d.Set("server", containGroup.Server)
	}

	if containGroup.Reponame != nil {
		_ = d.Set("reponame", containGroup.Reponame)
	}

	if containGroup.TagName != nil {
		_ = d.Set("tag_name", containGroup.TagName)
	}

	if containGroup.ClusterName != nil {
		_ = d.Set("cluster_name", containGroup.ClusterName)
	}

	if containGroup.NamespaceName != nil {
		_ = d.Set("namespace_name", containGroup.NamespaceName)
	}

	if containGroup.LbIp != nil {
		_ = d.Set("lb_ip", containGroup.LbIp)
	}

	if containGroup.ApplicationType != nil {
		_ = d.Set("application_type", containGroup.ApplicationType)
	}

	if containGroup.ClusterIp != nil {
		_ = d.Set("cluster_ip", containGroup.ClusterIp)
	}

	if containGroup.Envs != nil {
		envsList := []interface{}{}
		for _, envs := range containGroup.Envs {
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

	if containGroup.ApplicationName != nil {
		_ = d.Set("application_name", containGroup.ApplicationName)
	}

	if containGroup.Message != nil {
		_ = d.Set("message", containGroup.Message)
	}

	if containGroup.Status != nil {
		_ = d.Set("status", containGroup.Status)
	}

	if containGroup.MicroserviceType != nil {
		_ = d.Set("microservice_type", containGroup.MicroserviceType)
	}

	if containGroup.InstanceCount != nil {
		_ = d.Set("instance_count", containGroup.InstanceCount)
	}

	if containGroup.UpdatedTime != nil {
		_ = d.Set("updated_time", containGroup.UpdatedTime)
	}

	if containGroup.MaxSurge != nil {
		_ = d.Set("max_surge", containGroup.MaxSurge)
	}

	if containGroup.MaxUnavailable != nil {
		_ = d.Set("max_unavailable", containGroup.MaxUnavailable)
	}

	healthCheckSettingsMap := map[string]interface{}{}
	if containGroup.HealthCheckSettings != nil {

		if containGroup.HealthCheckSettings.LivenessProbe != nil {
			livenessProbeMap := map[string]interface{}{}

			if containGroup.HealthCheckSettings.LivenessProbe.ActionType != nil {
				livenessProbeMap["action_type"] = containGroup.HealthCheckSettings.LivenessProbe.ActionType
			}

			if containGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds != nil {
				livenessProbeMap["initial_delay_seconds"] = containGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds
			}

			if containGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds != nil {
				livenessProbeMap["timeout_seconds"] = containGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds
			}

			if containGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds != nil {
				livenessProbeMap["period_seconds"] = containGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds
			}

			if containGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold != nil {
				livenessProbeMap["success_threshold"] = containGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold
			}

			if containGroup.HealthCheckSettings.LivenessProbe.FailureThreshold != nil {
				livenessProbeMap["failure_threshold"] = containGroup.HealthCheckSettings.LivenessProbe.FailureThreshold
			}

			if containGroup.HealthCheckSettings.LivenessProbe.Scheme != nil {
				livenessProbeMap["scheme"] = containGroup.HealthCheckSettings.LivenessProbe.Scheme
			}

			if containGroup.HealthCheckSettings.LivenessProbe.Port != nil {
				livenessProbeMap["port"] = containGroup.HealthCheckSettings.LivenessProbe.Port
			}

			if containGroup.HealthCheckSettings.LivenessProbe.Path != nil {
				livenessProbeMap["path"] = containGroup.HealthCheckSettings.LivenessProbe.Path
			}

			if containGroup.HealthCheckSettings.LivenessProbe.Command != nil {
				livenessProbeMap["command"] = containGroup.HealthCheckSettings.LivenessProbe.Command
			}

			if containGroup.HealthCheckSettings.LivenessProbe.Type != nil {
				livenessProbeMap["type"] = containGroup.HealthCheckSettings.LivenessProbe.Type
			}

			healthCheckSettingsMap["liveness_probe"] = []interface{}{livenessProbeMap}
		}

		if containGroup.HealthCheckSettings.ReadinessProbe != nil {
			readinessProbeMap := map[string]interface{}{}

			if containGroup.HealthCheckSettings.ReadinessProbe.ActionType != nil {
				readinessProbeMap["action_type"] = containGroup.HealthCheckSettings.ReadinessProbe.ActionType
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds != nil {
				readinessProbeMap["initial_delay_seconds"] = containGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds != nil {
				readinessProbeMap["timeout_seconds"] = containGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds != nil {
				readinessProbeMap["period_seconds"] = containGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold != nil {
				readinessProbeMap["success_threshold"] = containGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold != nil {
				readinessProbeMap["failure_threshold"] = containGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.Scheme != nil {
				readinessProbeMap["scheme"] = containGroup.HealthCheckSettings.ReadinessProbe.Scheme
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.Port != nil {
				readinessProbeMap["port"] = containGroup.HealthCheckSettings.ReadinessProbe.Port
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.Path != nil {
				readinessProbeMap["path"] = containGroup.HealthCheckSettings.ReadinessProbe.Path
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.Command != nil {
				readinessProbeMap["command"] = containGroup.HealthCheckSettings.ReadinessProbe.Command
			}

			if containGroup.HealthCheckSettings.ReadinessProbe.Type != nil {
				readinessProbeMap["type"] = containGroup.HealthCheckSettings.ReadinessProbe.Type
			}

			healthCheckSettingsMap["readiness_probe"] = []interface{}{readinessProbeMap}
		}

	}
	_ = d.Set("health_check_settings", []interface{}{healthCheckSettingsMap})

	return nil
}

func resourceTencentCloudTsfContainGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_contain_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyContainerGroupRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"application_id", "namespace_id", "group_name", "instance_num", "cluster_id", "cpu_limit", "mem_limit", "group_comment", "cpu_request", "mem_request", "group_resource_type", "agent_cpu_request", "agent_cpu_limit", "agent_mem_request", "agent_mem_limit", "istio_cpu_request", "istio_cpu_limit", "istio_mem_request", "istio_mem_limit", "group_id", "current_num", "create_time", "server", "reponame", "tag_name", "cluster_name", "namespace_name", "lb_ip", "application_type", "cluster_ip", "envs", "application_name", "message", "status", "microservice_type", "instance_count", "updated_time", "max_surge", "max_unavailable", "health_check_settings"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("access_type") {
		if v, _ := d.GetOk("access_type"); v != nil {
			request.AccessType = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("protocol_ports") {
		if v, ok := d.GetOk("protocol_ports"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				protocolPort := tsf.ProtocolPort{}
				if v, ok := dMap["protocol"]; ok {
					protocolPort.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["port"]; ok {
					protocolPort.Port = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["target_port"]; ok {
					protocolPort.TargetPort = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["node_port"]; ok {
					protocolPort.NodePort = helper.IntInt64(v.(int))
				}
				request.ProtocolPorts = append(request.ProtocolPorts, &protocolPort)
			}
		}
	}

	if d.HasChange("update_type") {
		if v, _ := d.GetOk("update_type"); v != nil {
			request.UpdateType = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("update_ivl") {
		if v, _ := d.GetOk("update_ivl"); v != nil {
			request.UpdateIvl = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("subnet_id") {
		if v, ok := d.GetOk("subnet_id"); ok {
			request.SubnetId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyContainerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf containGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfContainGroupRead(d, meta)
}

func resourceTencentCloudTsfContainGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_contain_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	groupId := d.Id()

	if err := service.DeleteTsfContainGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
