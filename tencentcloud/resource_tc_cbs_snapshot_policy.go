/*
Provides a snapshot policy resource.

Example Usage

```hcl
resource "tencentcloud-cbs_snapshot_policy" "snapshot_policy" {
	snapshot_policy_name  = "mysnapshotpolicyname"
	repeat_weekdays = [1, 4]
	repeat_hours = [1]
	retention_days = 7
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

	"github.com/hashicorp/terraform/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
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
				Description: "Periodic snapshot is enabled, the available values are [0, 1, 2, 3, 4, 5, 6]. 0 means Sunday, 1-6 means Monday to Saturday.",
			},
			"repeat_hours": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validateIntegerInRange(0, 23),
				},
				Description: "Trigger times of periodic snapshot, the available values are 0 to 23. The 0 means 00:00, and so on.",
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
	logId := GetLogId(nil)

	request := cbs.NewCreateAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyName = stringToPointer(d.Get("snapshot_policy_name").(string))

	request.Policy = make([]*cbs.Policy, 0, 1)
	policy := &cbs.Policy{}
	repeatWeekdays := d.Get("repeat_weekdays").([]interface{})
	policy.DayOfWeek = make([]*uint64, 0, len(repeatWeekdays))
	for _, v := range repeatWeekdays {
		policy.DayOfWeek = append(policy.DayOfWeek, intToPointer(v.(int)))
	}
	repeatHours := d.Get("repeat_hours").([]interface{})
	policy.Hour = make([]*uint64, 0, len(repeatHours))
	for _, v := range repeatHours {
		policy.Hour = append(policy.Hour, intToPointer(v.(int)))
	}
	request.Policy = append(request.Policy, policy)

	if v, ok := d.GetOk("retention_days"); ok {
		request.RetentionDays = intToPointer(v.(int))
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateAutoSnapshotPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	if response.Response.AutoSnapshotPolicyId == nil {
		return fmt.Errorf("snapshot policy id is nil")
	}
	d.SetId(*response.Response.AutoSnapshotPolicyId)
	return resourceTencentCloudCbsSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	policyId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	policy, err := cbsService.DescribeSnapshotPolicyById(ctx, policyId)
	if err != nil {
		return err
	}

	d.Set("snapshot_policy_name", policy.AutoSnapshotPolicyName)
	if len(policy.Policy) > 0 {
		d.Set("repeat_weekdays", flattenIntList(policy.Policy[0].DayOfWeek))
		d.Set("repeat_hours", flattenIntList(policy.Policy[0].Hour))
	}
	d.Set("retention_days", policy.RetentionDays)
	return nil
}

func resourceTencentCloudCbsSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	policyId := d.Id()
	request := cbs.NewModifyAutoSnapshotPolicyAttributeRequest()
	request.AutoSnapshotPolicyId = &policyId
	if d.HasChange("snapshot_policy_name") {
		request.AutoSnapshotPolicyName = stringToPointer(d.Get("snapshot_policy_name").(string))
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = intToPointer(d.Get("retention_days").(int))
	}
	if d.HasChange("repeat_weekdays") || d.HasChange("repeat_hours") {
		request.Policy = make([]*cbs.Policy, 0, 1)
		policy := &cbs.Policy{}
		repeatWeekdays := d.Get("repeat_weekdays").([]interface{})
		policy.DayOfWeek = make([]*uint64, 0, len(repeatWeekdays))
		for _, v := range repeatWeekdays {
			policy.DayOfWeek = append(policy.DayOfWeek, intToPointer(v.(int)))
		}
		repeatHours := d.Get("repeat_hours").([]interface{})
		policy.Hour = make([]*uint64, 0, len(repeatHours))
		for _, v := range repeatHours {
			policy.Hour = append(policy.Hour, intToPointer(v.(int)))
		}
		request.Policy = append(request.Policy, policy)
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().ModifyAutoSnapshotPolicyAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	policyId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cbsService.DeleteSnapshotPolicy(ctx, policyId)
	if err != nil {
		return err
	}
	return nil
}
