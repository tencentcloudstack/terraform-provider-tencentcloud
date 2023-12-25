package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfApiGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApiGroupCreate,
		Read:   resourceTencentCloudTsfApiGroupRead,
		Update: resourceTencentCloudTsfApiGroupUpdate,
		Delete: resourceTencentCloudTsfApiGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api Group Id.",
			},
			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group name, cannot contain Chinese.",
			},

			"group_context": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "grouping context.",
			},

			"auth_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "authentication type. secret: key authentication; none: no authentication.",
			},

			"description": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "remarks.",
			},

			"group_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "grouping type, default ms. ms: microservice grouping; external: external Api grouping.",
			},

			"gateway_instance_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "gateway entity ID.",
			},

			"namespace_name_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace parameter key value.",
			},

			"service_name_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "microservice name parameter key value.",
			},

			"namespace_name_key_position": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace parameter position, path, header or query, the default is path.",
			},

			"service_name_key_position": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "microservice name parameter position, path, header or query, the default is path.",
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
				Description: "api group bound gateway deployment group.",
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
							Description: "application ID.",
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
				Description: "number of APIs.",
			},
			"acl_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access group ACL type.",
			},
			"gateway_instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of gateway instance.",
			},
		},
	}
}

func resourceTencentCloudTsfApiGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_api_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().CreateApiGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	groupId = *response.Response.Result
	d.SetId(groupId)

	return resourceTencentCloudTsfApiGroupRead(d, meta)
}

func resourceTencentCloudTsfApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_api_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()

	apiGroup, err := service.DescribeTsfApiGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if apiGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApiGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("group_id", groupId)

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

	if apiGroup.Status != nil {
		_ = d.Set("status", apiGroup.Status)
	}

	if apiGroup.CreatedTime != nil {
		_ = d.Set("created_time", apiGroup.CreatedTime)
	}

	if apiGroup.UpdatedTime != nil {
		_ = d.Set("updated_time", apiGroup.UpdatedTime)
	}

	if apiGroup.BindedGatewayDeployGroups != nil {
		bindedGatewayDeployGroupsList := []interface{}{}
		for _, bindedGatewayDeployGroups := range apiGroup.BindedGatewayDeployGroups {
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

		_ = d.Set("binded_gateway_deploy_groups", bindedGatewayDeployGroupsList)
	}

	if apiGroup.ApiCount != nil {
		_ = d.Set("api_count", apiGroup.ApiCount)
	}

	if apiGroup.AclMode != nil {
		_ = d.Set("acl_mode", apiGroup.AclMode)
	}

	if apiGroup.GatewayInstanceType != nil {
		_ = d.Set("gateway_instance_type", apiGroup.GatewayInstanceType)
	}

	return nil
}

func resourceTencentCloudTsfApiGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_api_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tsf.NewUpdateApiGroupRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_name", "group_type", "gateway_instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().UpdateApiGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_api_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	groupId := d.Id()

	if err := service.DeleteTsfApiGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
