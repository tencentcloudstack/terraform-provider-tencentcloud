package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMonitorPolicyBindingObject() *schema.Resource {
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
					if vmap["region"] != nil {
						hashMap["region"] = vmap["region"]
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
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Region.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentMonitorPolicyBindingObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_object.create")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
		var region string

		if v, ok := m["region"]; ok && v.(string) != "" {
			region = v.(string)
		} else {
			region = monitorService.client.Region
		}
		if v, ok := MonitorRegionMap[region]; ok {
			dimension.Region = helper.String(v)
		} else {
			return fmt.Errorf("monitor not support region `%s` bind", region)
		}
		dimension.Dimensions = &dimensionsJson
		request.Dimensions = append(request.Dimensions, &dimension)
	}

	request.Module = helper.String("monitor")
	if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().BindingPolicyObject(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_policy_binding_object.read")()
	defer tccommon.InconsistentCheck(d, meta)()
	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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

	if info.OriginId == nil {
		return fmt.Errorf("OriginId is nil")
	}
	originId, err := strconv.Atoi(*info.OriginId)
	if err != nil {
		return err
	}
	regionList, err := monitorService.DescribePolicyObjectCount(ctx, originId)
	if err != nil {
		return err
	}

	newDimensions := make([]interface{}, 0, 10)
	for _, regionInfo := range regionList {
		if regionInfo.Count != nil && *regionInfo.Count == 0 {
			continue
		}
		region := MonitorRegionMapName[*regionInfo.Region]
		objects, err := monitorService.DescribeBindingAlarmPolicyObjectList(ctx, policyId, region)
		if err != nil {
			return err
		}

		for _, item := range objects {
			dimensionsJson := item.Dimensions
			uniqueId := item.UniqueId
			newDimension := map[string]interface{}{
				"dimensions_json": dimensionsJson,
				"unique_id":       uniqueId,
				"region":          region,
			}
			newDimensions = append(newDimensions, newDimension)

		}
	}

	return d.Set("dimensions", newDimensions)
}

func resourceTencentMonitorPolicyBindingObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_object.delete")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		policyId       = d.Id()
	)

	info, err := monitorService.DescribeAlarmPolicyById(ctx, policyId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("alarm policy %s not exist", policyId)
	}

	if info.OriginId == nil {
		return fmt.Errorf("OriginId is nil")
	}
	originId, err := strconv.Atoi(*info.OriginId)
	if err != nil {
		return err
	}
	regionList, err := monitorService.DescribePolicyObjectCount(ctx, originId)
	if err != nil {
		return err
	}

	for _, regionInfo := range regionList {
		if regionInfo.Count != nil && *regionInfo.Count == 0 {
			continue
		}

		request := monitor.NewUnBindingAllPolicyObjectRequest()
		request.Module = helper.String("monitor")
		request.GroupId = helper.Int64(0)
		request.PolicyId = &policyId

		region := MonitorRegionMapName[*regionInfo.Region]
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, e := monitorService.client.UseMonitorClientRegion(region).UnBindingAllPolicyObject(request); e != nil {
				return tccommon.RetryError(e, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
