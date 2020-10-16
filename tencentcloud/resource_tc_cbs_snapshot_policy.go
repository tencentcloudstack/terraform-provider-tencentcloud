/*
Provides a snapshot policy resource.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "mysnapshotpolicyname"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}
```

Import

CBS snapshot policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot_policy.snapshot_policy asp-jliex1tn
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCbsSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotPolicyCreate,
		Read:   resourceTencentCloudCbsSnapshotPolicyRead,
		Update: resourceTencentCloudCbsSnapshotPolicyUpdate,
		Delete: resourceTencentCloudCbsSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"snapshot_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Name of snapshot policy. The maximum length can not exceed 60 bytes.",
			},
			"repeat_weekdays": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validateIntegerInRange(0, 6),
				},
				Description: "Periodic snapshot is enabled. Valid values: [0, 1, 2, 3, 4, 5, 6]. 0 means Sunday, 1-6 means Monday to Saturday.",
			},
			"repeat_hours": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validateIntegerInRange(0, 23),
				},
				Description: "Trigger times of periodic snapshot. Valid value ranges: (0~23). The 0 means 00:00, and so on.",
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7,
				Description: "Retention days of the snapshot, and the default value is 7.",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_policy.create")()

	logId := getLogId(contextNil)

	request := cbs.NewCreateAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyName = helper.String(d.Get("snapshot_policy_name").(string))

	request.Policy = make([]*cbs.Policy, 0, 1)
	policy := &cbs.Policy{}
	repeatWeekdays := d.Get("repeat_weekdays").([]interface{})
	policy.DayOfWeek = make([]*uint64, 0, len(repeatWeekdays))
	for _, v := range repeatWeekdays {
		policy.DayOfWeek = append(policy.DayOfWeek, helper.IntUint64(v.(int)))
	}
	repeatHours := d.Get("repeat_hours").([]interface{})
	policy.Hour = make([]*uint64, 0, len(repeatHours))
	for _, v := range repeatHours {
		policy.Hour = append(policy.Hour, helper.IntUint64(v.(int)))
	}
	request.Policy = append(request.Policy, policy)

	if v, ok := d.GetOk("retention_days"); ok {
		request.RetentionDays = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateAutoSnapshotPolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		if response.Response.AutoSnapshotPolicyId == nil {
			return resource.NonRetryableError(fmt.Errorf("snapshot policy id is nil"))
		}
		d.SetId(*response.Response.AutoSnapshotPolicyId)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot policy failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudCbsSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	policyId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var policy *cbs.AutoSnapshotPolicy
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		policy, e = cbsService.DescribeSnapshotPolicyById(ctx, policyId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot policy failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("snapshot_policy_name", policy.AutoSnapshotPolicyName)
	if len(policy.Policy) > 0 {
		_ = d.Set("repeat_weekdays", helper.Uint64sInterfaces(policy.Policy[0].DayOfWeek))
		_ = d.Set("repeat_hours", helper.Uint64sInterfaces(policy.Policy[0].Hour))
	}
	_ = d.Set("retention_days", policy.RetentionDays)

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_policy.update")()

	logId := getLogId(contextNil)

	policyId := d.Id()
	request := cbs.NewModifyAutoSnapshotPolicyAttributeRequest()
	request.AutoSnapshotPolicyId = &policyId
	if d.HasChange("snapshot_policy_name") {
		request.AutoSnapshotPolicyName = helper.String(d.Get("snapshot_policy_name").(string))
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = helper.IntUint64(d.Get("retention_days").(int))
	}
	if d.HasChange("repeat_weekdays") || d.HasChange("repeat_hours") {
		request.Policy = make([]*cbs.Policy, 0, 1)
		policy := &cbs.Policy{}
		repeatWeekdays := d.Get("repeat_weekdays").([]interface{})
		policy.DayOfWeek = make([]*uint64, 0, len(repeatWeekdays))
		for _, v := range repeatWeekdays {
			policy.DayOfWeek = append(policy.DayOfWeek, helper.IntUint64(v.(int)))
		}
		repeatHours := d.Get("repeat_hours").([]interface{})
		policy.Hour = make([]*uint64, 0, len(repeatHours))
		for _, v := range repeatHours {
			policy.Hour = append(policy.Hour, helper.IntUint64(v.(int)))
		}
		request.Policy = append(request.Policy, policy)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().ModifyAutoSnapshotPolicyAttribute(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cbs snapshot policy failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	policyId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteSnapshotPolicy(ctx, policyId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs snapshot policy failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
