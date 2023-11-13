/*
Provides a resource to create a tsf lane

Example Usage

```hcl
resource "tencentcloud_tsf_lane" "lane" {
  lane_name = ""
  remark = ""
  lane_group_list {
		group_id = ""
		entrance =
		lane_group_id = ""
		lane_id = ""
		group_name = ""
		application_id = ""
		application_name = ""
		namespace_id = ""
		namespace_name = ""
		create_time =
		update_time =
		cluster_type = ""

  }
  program_id_list =
        }
```

Import

tsf lane can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_lane.lane lane_id
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

func resourceTencentCloudTsfLane() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfLaneCreate,
		Read:   resourceTencentCloudTsfLaneRead,
		Update: resourceTencentCloudTsfLaneUpdate,
		Delete: resourceTencentCloudTsfLaneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"lane_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lane name.",
			},

			"remark": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lane Remarks.",
			},

			"lane_group_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Swimlane Deployment Group Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Deployment group ID.",
						},
						"entrance": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to enter the application.",
						},
						"lane_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Swimlane deployment group ID.",
						},
						"lane_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Lane ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Deployment group name.",
						},
						"application_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application ID.",
						},
						"application_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application name.",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace ID.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace name.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Update time.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster type.",
						},
					},
				},
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Update time.",
			},

			"entrance": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enter the application.",
			},

			"namespace_id_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of namespaces to which the swimlane has associated deployment groups.",
			},
		},
	}
}

func resourceTencentCloudTsfLaneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_lane.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateLaneRequest()
		response = tsf.NewCreateLaneResponse()
		laneId   string
	)
	if v, ok := d.GetOk("lane_name"); ok {
		request.LaneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lane_group_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			laneGroup := tsf.LaneGroup{}
			if v, ok := dMap["group_id"]; ok {
				laneGroup.GroupId = helper.String(v.(string))
			}
			if v, ok := dMap["entrance"]; ok {
				laneGroup.Entrance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["lane_group_id"]; ok {
				laneGroup.LaneGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["lane_id"]; ok {
				laneGroup.LaneId = helper.String(v.(string))
			}
			if v, ok := dMap["group_name"]; ok {
				laneGroup.GroupName = helper.String(v.(string))
			}
			if v, ok := dMap["application_id"]; ok {
				laneGroup.ApplicationId = helper.String(v.(string))
			}
			if v, ok := dMap["application_name"]; ok {
				laneGroup.ApplicationName = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_id"]; ok {
				laneGroup.NamespaceId = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_name"]; ok {
				laneGroup.NamespaceName = helper.String(v.(string))
			}
			if v, ok := dMap["create_time"]; ok {
				laneGroup.CreateTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["update_time"]; ok {
				laneGroup.UpdateTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cluster_type"]; ok {
				laneGroup.ClusterType = helper.String(v.(string))
			}
			request.LaneGroupList = append(request.LaneGroupList, &laneGroup)
		}
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateLane(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf lane failed, reason:%+v", logId, err)
		return err
	}

	laneId = *response.Response.LaneId
	d.SetId(laneId)

	return resourceTencentCloudTsfLaneRead(d, meta)
}

func resourceTencentCloudTsfLaneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_lane.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	laneId := d.Id()

	lane, err := service.DescribeTsfLaneById(ctx, laneId)
	if err != nil {
		return err
	}

	if lane == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfLane` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if lane.LaneName != nil {
		_ = d.Set("lane_name", lane.LaneName)
	}

	if lane.Remark != nil {
		_ = d.Set("remark", lane.Remark)
	}

	if lane.LaneGroupList != nil {
		laneGroupListList := []interface{}{}
		for _, laneGroupList := range lane.LaneGroupList {
			laneGroupListMap := map[string]interface{}{}

			if lane.LaneGroupList.GroupId != nil {
				laneGroupListMap["group_id"] = lane.LaneGroupList.GroupId
			}

			if lane.LaneGroupList.Entrance != nil {
				laneGroupListMap["entrance"] = lane.LaneGroupList.Entrance
			}

			if lane.LaneGroupList.LaneGroupId != nil {
				laneGroupListMap["lane_group_id"] = lane.LaneGroupList.LaneGroupId
			}

			if lane.LaneGroupList.LaneId != nil {
				laneGroupListMap["lane_id"] = lane.LaneGroupList.LaneId
			}

			if lane.LaneGroupList.GroupName != nil {
				laneGroupListMap["group_name"] = lane.LaneGroupList.GroupName
			}

			if lane.LaneGroupList.ApplicationId != nil {
				laneGroupListMap["application_id"] = lane.LaneGroupList.ApplicationId
			}

			if lane.LaneGroupList.ApplicationName != nil {
				laneGroupListMap["application_name"] = lane.LaneGroupList.ApplicationName
			}

			if lane.LaneGroupList.NamespaceId != nil {
				laneGroupListMap["namespace_id"] = lane.LaneGroupList.NamespaceId
			}

			if lane.LaneGroupList.NamespaceName != nil {
				laneGroupListMap["namespace_name"] = lane.LaneGroupList.NamespaceName
			}

			if lane.LaneGroupList.CreateTime != nil {
				laneGroupListMap["create_time"] = lane.LaneGroupList.CreateTime
			}

			if lane.LaneGroupList.UpdateTime != nil {
				laneGroupListMap["update_time"] = lane.LaneGroupList.UpdateTime
			}

			if lane.LaneGroupList.ClusterType != nil {
				laneGroupListMap["cluster_type"] = lane.LaneGroupList.ClusterType
			}

			laneGroupListList = append(laneGroupListList, laneGroupListMap)
		}

		_ = d.Set("lane_group_list", laneGroupListList)

	}

	if lane.ProgramIdList != nil {
		_ = d.Set("program_id_list", lane.ProgramIdList)
	}

	if lane.CreateTime != nil {
		_ = d.Set("create_time", lane.CreateTime)
	}

	if lane.UpdateTime != nil {
		_ = d.Set("update_time", lane.UpdateTime)
	}

	if lane.Entrance != nil {
		_ = d.Set("entrance", lane.Entrance)
	}

	if lane.NamespaceIdList != nil {
		_ = d.Set("namespace_id_list", lane.NamespaceIdList)
	}

	return nil
}

func resourceTencentCloudTsfLaneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_lane.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyLaneRequest()

	laneId := d.Id()

	request.LaneId = &laneId

	immutableArgs := []string{"lane_name", "remark", "lane_group_list", "program_id_list", "create_time", "update_time", "entrance", "namespace_id_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("lane_name") {
		if v, ok := d.GetOk("lane_name"); ok {
			request.LaneName = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyLane(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf lane failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfLaneRead(d, meta)
}

func resourceTencentCloudTsfLaneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_lane.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	laneId := d.Id()

	if err := service.DeleteTsfLaneById(ctx, laneId); err != nil {
		return err
	}

	return nil
}
