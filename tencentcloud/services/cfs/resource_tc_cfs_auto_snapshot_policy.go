package cfs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAutoSnapshotPolicyCreate,
		Read:   resourceTencentCloudCfsAutoSnapshotPolicyRead,
		Update: resourceTencentCloudCfsAutoSnapshotPolicyUpdate,
		Delete: resourceTencentCloudCfsAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"hour": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The time point when to repeat the snapshot operation.",
			},

			"policy_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Policy name.",
			},

			"day_of_week": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The day of the week on which to repeat the snapshot operation.",
			},

			"alive_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Snapshot retention period.",
			},

			"day_of_month": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The specific day (day 1 to day 31) of the month on which to create a snapshot.",
			},

			"interval_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The snapshot interval, in days.",
			},
		},
	}
}

func resourceTencentCloudCfsAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request              = cfs.NewCreateAutoSnapshotPolicyRequest()
		response             = cfs.NewCreateAutoSnapshotPolicyResponse()
		autoSnapshotPolicyId string
	)
	if v, ok := d.GetOk("day_of_week"); ok {
		request.DayOfWeek = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hour"); ok {
		request.Hour = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("alive_days"); ok {
		request.AliveDays = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("day_of_month"); ok {
		request.DayOfMonth = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("interval_days"); ok {
		request.IntervalDays = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().CreateAutoSnapshotPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs autoSnapshotPolicy failed, reason:%+v", logId, err)
		return err
	}

	autoSnapshotPolicyId = *response.Response.AutoSnapshotPolicyId
	d.SetId(autoSnapshotPolicyId)

	return resourceTencentCloudCfsAutoSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCfsAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	autoSnapshotPolicyId := d.Id()

	autoSnapshotPolicy, err := service.DescribeCfsAutoSnapshotPolicyById(ctx, autoSnapshotPolicyId)
	if err != nil {
		return err
	}

	if autoSnapshotPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_cfs_auto_snapshot_policy` %s does not exist", d.Id())
	}

	if autoSnapshotPolicy.DayOfWeek != nil {
		_ = d.Set("day_of_week", autoSnapshotPolicy.DayOfWeek)
	}

	if autoSnapshotPolicy.Hour != nil {
		_ = d.Set("hour", autoSnapshotPolicy.Hour)
	}

	if autoSnapshotPolicy.PolicyName != nil {
		_ = d.Set("policy_name", autoSnapshotPolicy.PolicyName)
	}

	if autoSnapshotPolicy.AliveDays != nil {
		_ = d.Set("alive_days", autoSnapshotPolicy.AliveDays)
	}

	if autoSnapshotPolicy.DayOfMonth != nil {
		_ = d.Set("day_of_month", autoSnapshotPolicy.DayOfMonth)
	}

	if autoSnapshotPolicy.IntervalDays != nil {
		_ = d.Set("interval_days", autoSnapshotPolicy.IntervalDays)
	}

	return nil
}

func resourceTencentCloudCfsAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cfs.NewUpdateAutoSnapshotPolicyRequest()

	autoSnapshotPolicyId := d.Id()

	request.AutoSnapshotPolicyId = &autoSnapshotPolicyId
	if d.HasChange("day_of_week") {
		if v, ok := d.GetOk("day_of_week"); ok {
			request.DayOfWeek = helper.String(v.(string))
		}
	}

	if d.HasChange("hour") {
		if v, ok := d.GetOk("hour"); ok {
			request.Hour = helper.String(v.(string))
		}
	}

	if d.HasChange("policy_name") {
		if v, ok := d.GetOk("policy_name"); ok {
			request.PolicyName = helper.String(v.(string))
		}
	}

	if d.HasChange("alive_days") {
		if v, ok := d.GetOkExists("alive_days"); ok {
			request.AliveDays = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("day_of_month") {
		if v, ok := d.GetOk("day_of_month"); ok {
			request.DayOfMonth = helper.String(v.(string))
		}
	}

	if d.HasChange("interval_days") {
		if v, ok := d.GetOkExists("interval_days"); ok {
			request.IntervalDays = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().UpdateAutoSnapshotPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs autoSnapshotPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfsAutoSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudCfsAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	autoSnapshotPolicyId := d.Id()

	if err := service.DeleteCfsAutoSnapshotPolicyById(ctx, autoSnapshotPolicyId); err != nil {
		return err
	}

	return nil
}
