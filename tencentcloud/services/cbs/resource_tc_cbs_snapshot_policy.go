package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsSnapshotPolicy() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "Name of snapshot policy. The maximum length can not exceed 60 bytes.",
			},

			"repeat_hours": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: tccommon.ValidateIntegerInRange(0, 23),
				},
				Description: "Trigger times of periodic snapshot. Valid value ranges: (0~23). The 0 means 00:00, and so on.",
			},

			"repeat_weekdays": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: tccommon.ValidateIntegerInRange(0, 6),
				},
				Description: "Periodic snapshot is enabled. Valid values: [0, 1, 2, 3, 4, 5, 6]. 0 means Sunday, 1-6 means Monday to Saturday.",
			},

			"repeat_months": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Specifies the days of the month on which scheduled snapshots are triggered. Valid value ranges: [1, 31]. For example, 5 means triggering snapshots on the 5th day of each month. Note: if a date such as 29, 30 or 31 does not exist in some months, those months will be skipped.",
			},

			"interval_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the interval in days for creating scheduled snapshots. Valid value ranges: [1, 365]. For example, setting it to 5 means a scheduled snapshot is created every 5 days. Note: when backing up by days, theoretically the first backup is on the day the policy is created; if the policy creation time is later than the configured backup time on that day, the first backup will be in the next backup cycle.",
			},

			"is_activated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the snapshot policy is activated. Default is true.",
			},

			"is_permanent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether snapshots created by this policy are retained permanently. Default is false. Note: if this parameter is set to true, retention_days must not be specified, otherwise a conflict error will be returned by the API.",
			},

			"retention_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Retention days of the snapshot, and the default value is 7.",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cbs.NewCreateAutoSnapshotPolicyRequest()
		response = cbs.NewCreateAutoSnapshotPolicyResponse()
		policy   = &cbs.Policy{}
		policyId string
	)

	if v, ok := d.GetOk("snapshot_policy_name"); ok {
		request.AutoSnapshotPolicyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repeat_weekdays"); ok {
		for _, item := range v.([]interface{}) {
			policy.DayOfWeek = append(policy.DayOfWeek, helper.IntUint64(item.(int)))
		}
	}

	if v, ok := d.GetOk("repeat_hours"); ok {
		for _, item := range v.([]interface{}) {
			policy.Hour = append(policy.Hour, helper.IntUint64(item.(int)))
		}
	}

	if v, ok := d.GetOk("repeat_months"); ok {
		for _, item := range v.([]interface{}) {
			policy.DayOfMonth = append(policy.DayOfMonth, helper.IntUint64(item.(int)))
		}
	}

	if v, ok := d.GetOkExists("interval_days"); ok {
		policy.IntervalDays = helper.IntUint64(v.(int))
	}

	request.Policy = []*cbs.Policy{policy}

	if v, ok := d.GetOkExists("is_activated"); ok {
		request.IsActivated = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_permanent"); ok {
		request.IsPermanent = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("retention_days"); ok {
		request.RetentionDays = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CreateAutoSnapshotPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create cbs snapshot policy failed, Response is nil."))
		}

		if result.Response.AutoSnapshotPolicyId == nil {
			return resource.NonRetryableError(fmt.Errorf("snapshot policy id is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot policy failed, reason:%+v", logId, err)
		return err
	}

	policyId = *response.Response.AutoSnapshotPolicyId
	d.SetId(policyId)

	return resourceTencentCloudCbsSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		policyId = d.Id()
	)

	var policy *cbs.AutoSnapshotPolicy
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		policy, e = service.DescribeSnapshotPolicyById(ctx, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot policy failed, reason:%+v", logId, err)
		return err
	}

	if policy == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cbs_snapshot_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if policy.AutoSnapshotPolicyName != nil {
		_ = d.Set("snapshot_policy_name", policy.AutoSnapshotPolicyName)
	}

	if len(policy.Policy) > 0 && policy.Policy[0] != nil {
		if policy.Policy[0].Hour != nil {
			_ = d.Set("repeat_hours", helper.Uint64sInterfaces(policy.Policy[0].Hour))
		}

		if policy.Policy[0].DayOfWeek != nil {
			_ = d.Set("repeat_weekdays", helper.Uint64sInterfaces(policy.Policy[0].DayOfWeek))
		}

		if policy.Policy[0].DayOfMonth != nil {
			_ = d.Set("repeat_months", helper.Uint64sInterfaces(policy.Policy[0].DayOfMonth))
		}

		if policy.Policy[0].IntervalDays != nil {
			_ = d.Set("interval_days", policy.Policy[0].IntervalDays)
		}
	}

	if policy.IsActivated != nil {
		_ = d.Set("is_activated", policy.IsActivated)
	}

	if policy.IsPermanent != nil {
		_ = d.Set("is_permanent", policy.IsPermanent)
	}

	if policy.RetentionDays != nil {
		_ = d.Set("retention_days", policy.RetentionDays)
	}

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cbs.NewModifyAutoSnapshotPolicyAttributeRequest()
		policyId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"snapshot_policy_name", "repeat_hours", "repeat_weekdays", "repeat_months", "interval_days", "is_activated", "is_permanent", "retention_days"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request.AutoSnapshotPolicyId = helper.String(policyId)

		if d.HasChange("snapshot_policy_name") {
			if v, ok := d.GetOk("snapshot_policy_name"); ok {
				request.AutoSnapshotPolicyName = helper.String(v.(string))
			}
		}

		if d.HasChange("repeat_weekdays") || d.HasChange("repeat_hours") || d.HasChange("repeat_months") || d.HasChange("interval_days") {
			policy := &cbs.Policy{}
			if v, ok := d.GetOk("repeat_weekdays"); ok {
				for _, item := range v.([]interface{}) {
					policy.DayOfWeek = append(policy.DayOfWeek, helper.IntUint64(item.(int)))
				}
			}

			if v, ok := d.GetOk("repeat_hours"); ok {
				for _, item := range v.([]interface{}) {
					policy.Hour = append(policy.Hour, helper.IntUint64(item.(int)))
				}
			}

			if v, ok := d.GetOk("repeat_months"); ok {
				for _, item := range v.([]interface{}) {
					policy.DayOfMonth = append(policy.DayOfMonth, helper.IntUint64(item.(int)))
				}
			}

			if v, ok := d.GetOk("interval_days"); ok {
				policy.IntervalDays = helper.IntUint64(v.(int))
			}

			request.Policy = []*cbs.Policy{policy}
		}

		if d.HasChange("is_activated") {
			if v, ok := d.GetOkExists("is_activated"); ok {
				request.IsActivated = helper.Bool(v.(bool))
			}
		}

		if d.HasChange("is_permanent") {
			if v, ok := d.GetOkExists("is_permanent"); ok {
				request.IsPermanent = helper.Bool(v.(bool))
			}
		}

		if d.HasChange("retention_days") {
			if v, ok := d.GetOkExists("retention_days"); ok {
				request.RetentionDays = helper.IntUint64(v.(int))
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().ModifyAutoSnapshotPolicyAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("modify cbs snapshot policy failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cbs snapshot policy failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCbsSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		policyId = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := service.DeleteSnapshotPolicy(ctx, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs snapshot policy failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
