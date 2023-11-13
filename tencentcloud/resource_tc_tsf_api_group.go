/*
Provides a resource to create a tsf api_group

Example Usage

```hcl
resource "tencentcloud_tsf_api_group" "api_group" {
  group_name = ""
  group_context = ""
  auth_type = ""
  description = ""
  group_type = ""
  gateway_instance_id = ""
  namespace_name_key = ""
  service_name_key = ""
  namespace_name_key_position = ""
  service_name_key_position = ""
  }
```

Import

tsf api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_group.api_group api_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfApiGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApiGroupCreate,
		Read:   resourceTencentCloudTsfApiGroupRead,
		Update: resourceTencentCloudTsfApiGroupUpdate,
		Delete: resourceTencentCloudTsfApiGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group name, cannot contain Chinese.",
			},

			"group_context": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grouping context.",
			},

			"auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Authentication type. secret: key authentication; none: no authentication.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},

			"group_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Grouping type, default ms. ms: microservice grouping; external: external Api grouping.",
			},

			"gateway_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gateway entity ID.",
			},

			"namespace_name_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace parameter key value.",
			},

			"service_name_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Microservice name parameter key value.",
			},

			"namespace_name_key_position": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace parameter position, path, header or query, the default is path.",
			},

			"service_name_key_position": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Microservice name parameter position, path, header or query, the default is path.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API group informationNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api Group Id.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api Group name.",
						},
						"group_context": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grouping context.",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authentication type. secret: key authentication; none: no authentication.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Release status, drafted: Not published. released: released.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group creation time such as: 2019-06-20 15:51:28.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group update time such as: 2019-06-20 15:51:28.",
						},
						"binded_gateway_deploy_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api group bound gateway deployment group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deploy_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway deployment group ID.",
									},
									"deploy_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway deployment group name.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application Name.",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application classification: V: virtual machine application, C: container application.",
									},
									"group_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment group application status, values: Running, Waiting, Paused, Updating, RollingBack, Abnormal, Unknown.",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster type, C: container, V: virtual machine.",
									},
								},
							},
						},
						"api_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of APIs.",
						},
						"acl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access group ACL type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe.",
						},
						"group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grouping type. ms: microservice grouping; external: external Api grouping.",
						},
						"gateway_instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of gateway instance.",
						},
						"gateway_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway instance ID.",
						},
						"namespace_name_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace parameter key value.",
						},
						"service_name_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice name parameter key value.",
						},
						"namespace_name_key_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace parameter position, path, header or query, the default is path.",
						},
						"service_name_key_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice name parameter position, path, header or query, the default is path.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfApiGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateApiGroupRequest()
		response = tsf.NewCreateApiGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_context"); ok {
		request.GroupContext = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_type"); ok {
		request.AuthType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_type"); ok {
		request.GroupType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gateway_instance_id"); ok {
		request.GatewayInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name_key"); ok {
		request.NamespaceNameKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name_key"); ok {
		request.ServiceNameKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name_key_position"); ok {
		request.NamespaceNameKeyPosition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name_key_position"); ok {
		request.ServiceNameKeyPosition = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateApiGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf apiGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(groupId)

	return resourceTencentCloudTsfApiGroupRead(d, meta)
}

func resourceTencentCloudTsfApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	apiGroupId := d.Id()

	apiGroup, err := service.DescribeTsfApiGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if apiGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApiGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if apiGroup.GroupName != nil {
		_ = d.Set("group_name", apiGroup.GroupName)
	}

	if apiGroup.GroupContext != nil {
		_ = d.Set("group_context", apiGroup.GroupContext)
	}

	if apiGroup.AuthType != nil {
		_ = d.Set("auth_type", apiGroup.AuthType)
	}

	if apiGroup.Description != nil {
		_ = d.Set("description", apiGroup.Description)
	}

	if apiGroup.GroupType != nil {
		_ = d.Set("group_type", apiGroup.GroupType)
	}

	if apiGroup.GatewayInstanceId != nil {
		_ = d.Set("gateway_instance_id", apiGroup.GatewayInstanceId)
	}

	if apiGroup.NamespaceNameKey != nil {
		_ = d.Set("namespace_name_key", apiGroup.NamespaceNameKey)
	}

	if apiGroup.ServiceNameKey != nil {
		_ = d.Set("service_name_key", apiGroup.ServiceNameKey)
	}

	if apiGroup.NamespaceNameKeyPosition != nil {
		_ = d.Set("namespace_name_key_position", apiGroup.NamespaceNameKeyPosition)
	}

	if apiGroup.ServiceNameKeyPosition != nil {
		_ = d.Set("service_name_key_position", apiGroup.ServiceNameKeyPosition)
	}

	if apiGroup.Result != nil {
		resultMap := map[string]interface{}{}

		if apiGroup.Result.GroupId != nil {
			resultMap["group_id"] = apiGroup.Result.GroupId
		}

		if apiGroup.Result.GroupName != nil {
			resultMap["group_name"] = apiGroup.Result.GroupName
		}

		if apiGroup.Result.GroupContext != nil {
			resultMap["group_context"] = apiGroup.Result.GroupContext
		}

		if apiGroup.Result.AuthType != nil {
			resultMap["auth_type"] = apiGroup.Result.AuthType
		}

		if apiGroup.Result.Status != nil {
			resultMap["status"] = apiGroup.Result.Status
		}

		if apiGroup.Result.CreatedTime != nil {
			resultMap["created_time"] = apiGroup.Result.CreatedTime
		}

		if apiGroup.Result.UpdatedTime != nil {
			resultMap["updated_time"] = apiGroup.Result.UpdatedTime
		}

		if apiGroup.Result.BindedGatewayDeployGroups != nil {
			bindedGatewayDeployGroupsList := []interface{}{}
			for _, bindedGatewayDeployGroups := range apiGroup.Result.BindedGatewayDeployGroups {
				bindedGatewayDeployGroupsMap := map[string]interface{}{}

				if bindedGatewayDeployGroups.DeployGroupId != nil {
					bindedGatewayDeployGroupsMap["deploy_group_id"] = bindedGatewayDeployGroups.DeployGroupId
				}

				if bindedGatewayDeployGroups.DeployGroupName != nil {
					bindedGatewayDeployGroupsMap["deploy_group_name"] = bindedGatewayDeployGroups.DeployGroupName
				}

				if bindedGatewayDeployGroups.ApplicationId != nil {
					bindedGatewayDeployGroupsMap["application_id"] = bindedGatewayDeployGroups.ApplicationId
				}

				if bindedGatewayDeployGroups.ApplicationName != nil {
					bindedGatewayDeployGroupsMap["application_name"] = bindedGatewayDeployGroups.ApplicationName
				}

				if bindedGatewayDeployGroups.ApplicationType != nil {
					bindedGatewayDeployGroupsMap["application_type"] = bindedGatewayDeployGroups.ApplicationType
				}

				if bindedGatewayDeployGroups.GroupStatus != nil {
					bindedGatewayDeployGroupsMap["group_status"] = bindedGatewayDeployGroups.GroupStatus
				}

				if bindedGatewayDeployGroups.ClusterType != nil {
					bindedGatewayDeployGroupsMap["cluster_type"] = bindedGatewayDeployGroups.ClusterType
				}

				bindedGatewayDeployGroupsList = append(bindedGatewayDeployGroupsList, bindedGatewayDeployGroupsMap)
			}

			resultMap["binded_gateway_deploy_groups"] = []interface{}{bindedGatewayDeployGroupsList}
		}

		if apiGroup.Result.ApiCount != nil {
			resultMap["api_count"] = apiGroup.Result.ApiCount
		}

		if apiGroup.Result.AclMode != nil {
			resultMap["acl_mode"] = apiGroup.Result.AclMode
		}

		if apiGroup.Result.Description != nil {
			resultMap["description"] = apiGroup.Result.Description
		}

		if apiGroup.Result.GroupType != nil {
			resultMap["group_type"] = apiGroup.Result.GroupType
		}

		if apiGroup.Result.GatewayInstanceType != nil {
			resultMap["gateway_instance_type"] = apiGroup.Result.GatewayInstanceType
		}

		if apiGroup.Result.GatewayInstanceId != nil {
			resultMap["gateway_instance_id"] = apiGroup.Result.GatewayInstanceId
		}

		if apiGroup.Result.NamespaceNameKey != nil {
			resultMap["namespace_name_key"] = apiGroup.Result.NamespaceNameKey
		}

		if apiGroup.Result.ServiceNameKey != nil {
			resultMap["service_name_key"] = apiGroup.Result.ServiceNameKey
		}

		if apiGroup.Result.NamespaceNameKeyPosition != nil {
			resultMap["namespace_name_key_position"] = apiGroup.Result.NamespaceNameKeyPosition
		}

		if apiGroup.Result.ServiceNameKeyPosition != nil {
			resultMap["service_name_key_position"] = apiGroup.Result.ServiceNameKeyPosition
		}

		_ = d.Set("result", []interface{}{resultMap})
	}

	return nil
}

func resourceTencentCloudTsfApiGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewUpdateApiGroupRequest()

	apiGroupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_name", "group_context", "auth_type", "description", "group_type", "gateway_instance_id", "namespace_name_key", "service_name_key", "namespace_name_key_position", "service_name_key_position", "result"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("group_name") {
		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("group_context") {
		if v, ok := d.GetOk("group_context"); ok {
			request.GroupContext = helper.String(v.(string))
		}
	}

	if d.HasChange("auth_type") {
		if v, ok := d.GetOk("auth_type"); ok {
			request.AuthType = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_name_key") {
		if v, ok := d.GetOk("namespace_name_key"); ok {
			request.NamespaceNameKey = helper.String(v.(string))
		}
	}

	if d.HasChange("service_name_key") {
		if v, ok := d.GetOk("service_name_key"); ok {
			request.ServiceNameKey = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_name_key_position") {
		if v, ok := d.GetOk("namespace_name_key_position"); ok {
			request.NamespaceNameKeyPosition = helper.String(v.(string))
		}
	}

	if d.HasChange("service_name_key_position") {
		if v, ok := d.GetOk("service_name_key_position"); ok {
			request.ServiceNameKeyPosition = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().UpdateApiGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf apiGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfApiGroupRead(d, meta)
}

func resourceTencentCloudTsfApiGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	apiGroupId := d.Id()

	if err := service.DeleteTsfApiGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
