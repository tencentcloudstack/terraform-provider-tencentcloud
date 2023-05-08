/*
Provides a resource to create a tsf lane

Example Usage

```hcl
resource "tencentcloud_tsf_lane" "lane" {
  lane_name = "lane-name-1"
  remark = "lane desc1"
  lane_group_list {
		group_id = "group-yn7j5l8a"
		entrance = true
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

func resourceTencentCloudTsfLane() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfLaneCreate,
		Read:   resourceTencentCloudTsfLaneRead,
		Update: resourceTencentCloudTsfLaneUpdate,
		Delete: resourceTencentCloudTsfLaneDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"lane_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lane id.",
			},

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

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "creation time.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "update time.",
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
							Description: "deployment group name.",
						},
						"application_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "application ID.",
						},
						"application_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "application name.",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace ID.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "namespace name.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "creation time.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "update time.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "cluster type.",
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

			"entrance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enter the application.",
			},
			"namespace_id_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
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

	laneId = *response.Response.Result
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

	_ = d.Set("lane_id", laneId)

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

			if laneGroupList.GroupId != nil {
				laneGroupListMap["group_id"] = laneGroupList.GroupId
			}

			if laneGroupList.Entrance != nil {
				laneGroupListMap["entrance"] = laneGroupList.Entrance
			}

			if laneGroupList.LaneGroupId != nil {
				laneGroupListMap["lane_group_id"] = laneGroupList.LaneGroupId
			}

			if laneGroupList.LaneId != nil {
				laneGroupListMap["lane_id"] = laneGroupList.LaneId
			}

			if laneGroupList.GroupName != nil {
				laneGroupListMap["group_name"] = laneGroupList.GroupName
			}

			if laneGroupList.ApplicationId != nil {
				laneGroupListMap["application_id"] = laneGroupList.ApplicationId
			}

			if laneGroupList.ApplicationName != nil {
				laneGroupListMap["application_name"] = laneGroupList.ApplicationName
			}

			if laneGroupList.NamespaceId != nil {
				laneGroupListMap["namespace_id"] = laneGroupList.NamespaceId
			}

			if laneGroupList.NamespaceName != nil {
				laneGroupListMap["namespace_name"] = laneGroupList.NamespaceName
			}

			if laneGroupList.CreateTime != nil {
				laneGroupListMap["create_time"] = laneGroupList.CreateTime
			}

			if laneGroupList.UpdateTime != nil {
				laneGroupListMap["update_time"] = laneGroupList.UpdateTime
			}

			if laneGroupList.ClusterType != nil {
				laneGroupListMap["cluster_type"] = laneGroupList.ClusterType
			}

			laneGroupListList = append(laneGroupListList, laneGroupListMap)
		}

		_ = d.Set("lane_group_list", laneGroupListList)

	}

	// if lane.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", lane.ProgramIdList)
	// }

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

	immutableArgs := []string{"lane_group_list", "program_id_list", "result"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("lane_name"); ok {
		request.LaneName = helper.String(v.(string))
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
