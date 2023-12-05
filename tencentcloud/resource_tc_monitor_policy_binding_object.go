package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudMonitorPolicyBindingObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorPolicyBindingObjectCreate,
		Read:   resourceTencentMonitorPolicyBindingObjectRead,
		Delete: resourceTencentMonitorPolicyBindingObjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Alarm policy ID for binding objects.",
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
					return helper.HashString(string(b))
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

func resourceTencentMonitorPolicyBindingObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_object.create")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewBindingPolicyObjectRequest()
		policyId       = d.Get("policy_id").(string)
	)

	info, err := monitorService.DescribeAlarmPolicyById(ctx, policyId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("alarm policy %s not exist", policyId)
	}
	request.GroupId = helper.Int64(0)
	request.PolicyId = &policyId
	dimensions := d.Get("dimensions").(*schema.Set).List()

	request.Dimensions = make([]*monitor.BindingPolicyObjectDimension, 0, len(dimensions))

	for _, v := range dimensions {
		m := v.(map[string]interface{})
		var dimension monitor.BindingPolicyObjectDimension
		var dimensionsJson = m["dimensions_json"].(string)
		var region = MonitorRegionMap[monitorService.client.Region]

		if region == "" {
			return fmt.Errorf("monitor not support region `%s` bind", monitorService.client.Region)
		}
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

	d.SetId(policyId)
	time.Sleep(3 * time.Second)

	return resourceTencentMonitorPolicyBindingObjectRead(d, meta)
}

func resourceTencentMonitorPolicyBindingObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_policy_binding_object.read")()
	defer inconsistentCheck(d, meta)()
	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		policyId       = d.Id()
	)

	_ = d.Set("policy_id", policyId)

	info, err := monitorService.DescribeAlarmPolicyById(ctx, policyId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("alarm policy %s not exist", policyId)
	}

	objects, err := monitorService.DescribeBindingAlarmPolicyObjectList(ctx, policyId)

	if err != nil {
		return err
	}

	newDimensions := make([]interface{}, 0, 10)

	for _, item := range objects {
		dimensionsJson := item.Dimensions
		uniqueId := item.UniqueId
		newDimensions = append(newDimensions, map[string]interface{}{
			"dimensions_json": dimensionsJson,
			"unique_id":       uniqueId,
		})
	}

	return d.Set("dimensions", newDimensions)
}

func resourceTencentMonitorPolicyBindingObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_object.delete")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		policyId       = d.Id()
	)

	info, err := monitorService.DescribeAlarmPolicyById(ctx, policyId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("alarm policy %s not exist", policyId)
	}

	objects, err := monitorService.DescribeBindingAlarmPolicyObjectList(ctx, policyId)

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

	var (
		request = monitor.NewUnBindingPolicyObjectRequest()
	)

	request.Module = helper.String("monitor")
	request.GroupId = helper.Int64(0)
	request.PolicyId = &policyId
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
