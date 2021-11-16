/*
Provides a resource for bind objects to a policy group resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_monitor_policy_binding_object.

Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "terraform_test"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}

#for cvm
resource "tencentcloud_monitor_binding_object" "binding" {
  group_id = tencentcloud_monitor_policy_group.group.id
  dimensions {
    dimensions_json = "{\"unInstanceId\":\"${data.tencentcloud_instances.instances.instance_list[0].instance_id}\"}"
  }
}
```

*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentMonitorBindingObject() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.60.5. Please use 'tencentcloud_monitor_policy_binding_object' instead.",
		Create:             resourceTencentMonitorBindingObjectCreate,
		Read:               resourceTencentMonitorBindingObjectRead,
		Delete:             resourceTencentMonitorBindingObjectDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Policy group ID for binding objects.",
			},
			"dimensions": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "A list objects. Each element contains the following attributes:",
				Set: func(v interface{}) int {
					vmap := v.(map[string]interface{})
					hashMap := map[string]interface{}{}
					if vmap["dimensions_json"] != nil {
						hashMap["dimensions_json"] = vmap["dimensions_json"]
					}
					b, _ := json.Marshal(hashMap)
					return hashcode.String(string(b))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimensions_json": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `Represents a collection of dimensions of an object instance, json format.eg:'{"unInstanceId":"ins-ot3cq4bi"}'.`,
						},
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object unique ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentMonitorBindingObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_object.create")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewBindingPolicyObjectRequest()
		idSeeds        []string
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("policy group %d not exist", groupId)
	}
	request.GroupId = &groupId
	dimensions := d.Get("dimensions").(*schema.Set).List()

	request.Dimensions = make([]*monitor.BindingPolicyObjectDimension, 0, len(dimensions))

	idSeeds = append(idSeeds, fmt.Sprintf("%d", groupId))

	for _, v := range dimensions {
		m := v.(map[string]interface{})
		var dimension monitor.BindingPolicyObjectDimension
		var dimensionsJson = m["dimensions_json"].(string)
		var region = MonitorRegionMap[monitorService.client.Region]

		if region == "" {
			return fmt.Errorf("monitor not support region `%s` bind", monitorService.client.Region)
		}
		idSeeds = append(idSeeds, dimensionsJson, region)
		dimension.Dimensions = &dimensionsJson
		dimension.Region = &region
		request.Dimensions = append(request.Dimensions, &dimension)

	}

	request.Module = helper.String("monitor")
	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().BindingPolicyObject(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	d.SetId(helper.DataResourceIdsHash(idSeeds))
	time.Sleep(3 * time.Second)

	objects, err := monitorService.DescribeBindingPolicyObjectList(ctx, groupId)

	if err != nil {
		return err
	}

	successDimensionsJsonMap := make(map[string]bool)
	bindingFails := make([]string, 0, len(request.Dimensions))
	for _, v := range objects {
		successDimensionsJsonMap[*v.Dimensions] = true
	}
	for _, v := range request.Dimensions {
		if !successDimensionsJsonMap[*v.Dimensions] {
			bindingFails = append(bindingFails, *v.Dimensions)
		}
	}

	if len(bindingFails) > 0 {
		return fmt.Errorf("bind objects to policy has partial failure,Please check if it is an instance of this region `%s`,[%s]",
			monitorService.client.Region, helper.SliceFieldSerialize(bindingFails))
	}

	return resourceTencentMonitorBindingObjectRead(d, meta)
}

func resourceTencentMonitorBindingObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_object.read")()
	defer inconsistentCheck(d, meta)()
	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("policy group %d not exist", groupId)
	}

	objects, err := monitorService.DescribeBindingPolicyObjectList(ctx, groupId)

	if err != nil {
		return err
	}

	getUniqueId := func(dimensionsJson string) (has bool, uniqueId string) {
		for _, item := range objects {
			if *item.Dimensions == dimensionsJson {
				uniqueId = *item.UniqueId
				has = true
				return
			}
		}
		return
	}

	dimensions := d.Get("dimensions").(*schema.Set).List()
	newDimensions := make([]interface{}, 0, len(dimensions))

	for _, v := range dimensions {
		m := v.(map[string]interface{})
		var dimensionsJson = m["dimensions_json"].(string)
		var has, uniqueId = getUniqueId(dimensionsJson)
		if has {
			newDimensions = append(newDimensions, map[string]interface{}{
				"dimensions_json": dimensionsJson,
				"unique_id":       uniqueId,
			})
		}
	}

	if len(newDimensions) == 0 {
		d.SetId("")
		return nil
	}

	return d.Set("dimensions", newDimensions)
}
func resourceTencentMonitorBindingObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_object.delete")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("policy group %d not exist", groupId)
	}

	objects, err := monitorService.DescribeBindingPolicyObjectList(ctx, groupId)

	if err != nil {
		return err
	}
	getUniqueId := func(dimensionsJson string) (has bool, uniqueId string) {
		for _, item := range objects {
			if *item.Dimensions == dimensionsJson {
				uniqueId = *item.UniqueId
				has = true
				return
			}
		}
		return
	}

	dimensions := d.Get("dimensions").(*schema.Set).List()
	uniqueIds := make([]*string, 0, len(dimensions))
	for _, v := range dimensions {
		m := v.(map[string]interface{})
		var dimensionsJson = m["dimensions_json"].(string)
		var has, uniqueId = getUniqueId(dimensionsJson)
		if has {
			uniqueIds = append(uniqueIds, &uniqueId)
		}
	}

	if len(uniqueIds) == 0 {
		d.SetId("")
		return nil
	}

	var (
		request = monitor.NewUnBindingPolicyObjectRequest()
	)

	request.Module = helper.String("monitor")
	request.GroupId = &groupId
	request.UniqueId = uniqueIds

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().UnBindingPolicyObject(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
