/*
Provides a resource to create a tsf group

Example Usage

```hcl
resource "tencentcloud_tsf_group" "group" {
  application_id = "application-xxxxxx"
  namespace_id = "namespace-xxxxxx"
  group_name = "terraform-test"
  cluster_id = "cluster-xxxxxx"
  group_desc = "terraform group desc"
  group_resource_type = "DEF"
  alias = "terraform desc"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_group.group group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfGroupCreate,
		Read:   resourceTencentCloudTsfGroupRead,
		Update: resourceTencentCloudTsfGroupUpdate,
		Delete: resourceTencentCloudTsfGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Group id.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application id.",
			},

			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace id.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group name.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},

			"group_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group desc.",
			},

			"group_resource_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "DEF",
				Description: "Deployment Group Resource Type, default value `DEF` means DEFAULT, optional `GW` means gateway, `SVL` means SERVERLESS.",
			},

			"alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Deployment Group Notes.",
			},

			"group_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group status.",
			},

			"package_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Package ID.",
			},

			"package_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "package name.",
			},

			"package_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Package version number.",
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

			"application_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application Name.",
			},

			"instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of machines in the deployment group.",
			},

			"run_instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of running machines in the deployment group.",
			},

			"startup_parameters": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group startup parameter information.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment group update time.",
			},

			"off_instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of machines stopped by the deployment group.",
			},

			"microservice_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Microservice Type.",
			},

			"application_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application types.",
			},

			"updated_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Deployment group update timestamp.",
			},

			"deploy_desc": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deploy application description information.",
			},

			"update_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Update method for rolling release.",
			},

			"deploy_beta_enable": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Post whether to enable beta batches.",
			},

			"deploy_batch": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of batch ratios released on a rolling basis.",
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
			},

			"deploy_exe_mode": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Batch Execution Mode for Rolling Releases.",
			},

			"deploy_wait_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The waiting time for each batch of a rolling release.",
			},

			"enable_health_check": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether health check is enabled.",
			},

			"package_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "package type.",
			},

			"start_script": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Start script base64 encoding.",
			},

			"stop_script": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Stop script base64 encoding.",
			},

			"agent_profile_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "javaagent information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent type.",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent version number.",
						},
					},
				},
			},

			"warmup_setting": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Preheat property configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable preheating.",
						},
						"warmup_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "warm-up time.",
						},
						"curvature": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "preheating curvature, value 1~5.",
						},
						"enabled_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable preheating protection, if protection is enabled, if more than 50% of the nodes are in preheating, the preheating will be terminated.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateGroupRequest()
		response = tsf.NewCreateGroupResponse()
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

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_desc"); ok {
		request.GroupDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_resource_type"); ok {
		request.GroupResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf group failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.Result
	d.SetId(groupId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:group/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfGroupRead(d, meta)
}

func resourceTencentCloudTsfGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()

	group, err := service.DescribeTsfGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if group.GroupId != nil {
		_ = d.Set("group_id", group.GroupId)
	}

	if group.ApplicationId != nil {
		_ = d.Set("application_id", group.ApplicationId)
	}

	if group.NamespaceId != nil {
		_ = d.Set("namespace_id", group.NamespaceId)
	}

	if group.GroupName != nil {
		_ = d.Set("group_name", group.GroupName)
	}

	if group.ClusterId != nil {
		_ = d.Set("cluster_id", group.ClusterId)
	}

	if group.GroupDesc != nil {
		_ = d.Set("group_desc", group.GroupDesc)
	}

	if group.GroupResourceType != nil {
		_ = d.Set("group_resource_type", group.GroupResourceType)
	}

	if group.Alias != nil {
		_ = d.Set("alias", group.Alias)
	}

	if group.GroupStatus != nil {
		_ = d.Set("group_status", group.GroupStatus)
	}

	if group.PackageId != nil {
		_ = d.Set("package_id", group.PackageId)
	}

	if group.PackageName != nil {
		_ = d.Set("package_name", group.PackageName)
	}

	if group.PackageVersion != nil {
		_ = d.Set("package_version", group.PackageVersion)
	}

	if group.ClusterName != nil {
		_ = d.Set("cluster_name", group.ClusterName)
	}

	if group.NamespaceName != nil {
		_ = d.Set("namespace_name", group.NamespaceName)
	}

	if group.ApplicationName != nil {
		_ = d.Set("application_name", group.ApplicationName)
	}

	if group.InstanceCount != nil {
		_ = d.Set("instance_count", group.InstanceCount)
	}

	if group.RunInstanceCount != nil {
		_ = d.Set("run_instance_count", group.RunInstanceCount)
	}

	if group.StartupParameters != nil {
		_ = d.Set("startup_parameters", group.StartupParameters)
	}

	if group.CreateTime != nil {
		_ = d.Set("create_time", group.CreateTime)
	}

	if group.UpdateTime != nil {
		_ = d.Set("update_time", group.UpdateTime)
	}

	if group.OffInstanceCount != nil {
		_ = d.Set("off_instance_count", group.OffInstanceCount)
	}

	if group.MicroserviceType != nil {
		_ = d.Set("microservice_type", group.MicroserviceType)
	}

	if group.ApplicationType != nil {
		_ = d.Set("application_type", group.ApplicationType)
	}

	if group.UpdatedTime != nil {
		_ = d.Set("updated_time", group.UpdatedTime)
	}

	if group.DeployDesc != nil {
		_ = d.Set("deploy_desc", group.DeployDesc)
	}

	if group.UpdateType != nil {
		_ = d.Set("update_type", group.UpdateType)
	}

	if group.DeployBetaEnable != nil {
		_ = d.Set("deploy_beta_enable", group.DeployBetaEnable)
	}

	if group.DeployBatch != nil {
		_ = d.Set("deploy_batch", group.DeployBatch)
	}

	if group.DeployExeMode != nil {
		_ = d.Set("deploy_exe_mode", group.DeployExeMode)
	}

	if group.DeployWaitTime != nil {
		_ = d.Set("deploy_wait_time", group.DeployWaitTime)
	}

	if group.EnableHealthCheck != nil {
		_ = d.Set("enable_health_check", group.EnableHealthCheck)
	}

	if group.PackageType != nil {
		_ = d.Set("package_type", group.PackageType)
	}

	if group.StartScript != nil {
		_ = d.Set("start_script", group.StartScript)
	}

	if group.StopScript != nil {
		_ = d.Set("stop_script", group.StopScript)
	}

	if group.AgentProfileList != nil {
		agentProfileListList := []interface{}{}
		for _, agentProfileList := range group.AgentProfileList {
			agentProfileListMap := map[string]interface{}{}

			if agentProfileList.AgentType != nil {
				agentProfileListMap["agent_type"] = agentProfileList.AgentType
			}

			if agentProfileList.AgentVersion != nil {
				agentProfileListMap["agent_version"] = agentProfileList.AgentVersion
			}

			agentProfileListList = append(agentProfileListList, agentProfileListMap)
		}

		_ = d.Set("agent_profile_list", agentProfileListList)

	}

	if group.WarmupSetting != nil {
		warmupSettingMap := map[string]interface{}{}

		if group.WarmupSetting.Enabled != nil {
			warmupSettingMap["enabled"] = group.WarmupSetting.Enabled
		}

		if group.WarmupSetting.WarmupTime != nil {
			warmupSettingMap["warmup_time"] = group.WarmupSetting.WarmupTime
		}

		if group.WarmupSetting.Curvature != nil {
			warmupSettingMap["curvature"] = group.WarmupSetting.Curvature
		}

		if group.WarmupSetting.EnabledProtection != nil {
			warmupSettingMap["enabled_protection"] = group.WarmupSetting.EnabledProtection
		}

		_ = d.Set("warmup_setting", []interface{}{warmupSettingMap})
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "group", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyGroupRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_id", "application_id", "group_name", "namespace_id", "cluster_id", "group_resource_type", "group_status", "package_id", "package_name", "package_version", "cluster_name", "namespace_name", "application_name", "instance_count", "run_instance_count", "startup_parameters", "create_time", "update_time", "off_instance_count", "microservice_type", "application_type", "updated_time", "deploy_desc", "update_type", "deploy_beta_enable", "deploy_batch", "deploy_exe_mode", "deploy_wait_time", "enable_health_check", "package_type", "start_script", "stop_script", "agent_profile_list", "warmup_setting"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("group_desc") {
		if v, ok := d.GetOk("group_desc"); ok {
			request.GroupDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("alias") {
		if v, ok := d.GetOk("alias"); ok {
			request.Alias = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf group failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "group", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfGroupRead(d, meta)
}

func resourceTencentCloudTsfGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	groupId := d.Id()

	if err := service.DeleteTsfGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
