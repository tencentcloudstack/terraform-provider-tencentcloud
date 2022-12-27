/*
Provides a resource to create a cfs auto_snapshot_policy

Example Usage

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  day_of_week = "1,2"
  hour = "2,3"
  policy_name = "policy_name"
  alive_days = 7
}
```

Import

cfs auto_snapshot_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy.auto_snapshot_policy auto_snapshot_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCfsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAutoSnapshotPolicyCreate,
		Read:   resourceTencentCloudCfsAutoSnapshotPolicyRead,
		Update: resourceTencentCloudCfsAutoSnapshotPolicyUpdate,
		Delete: resourceTencentCloudCfsAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"day_of_week": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The day of the week on which to repeat the snapshot operation.",
			},

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

			"alive_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Snapshot retention period.",
			},
		},
	}
}

func resourceTencentCloudCfsAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	if v, _ := d.GetOk("alive_days"); v != nil {
		request.AliveDays = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().CreateAutoSnapshotPolicy(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	autoSnapshotPolicyId := d.Id()

	autoSnapshotPolicy, err := service.DescribeCfsAutoSnapshotPolicyById(ctx, autoSnapshotPolicyId)
	if err != nil {
		return err
	}

	if autoSnapshotPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
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

	return nil
}

func resourceTencentCloudCfsAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
		if v, _ := d.GetOk("alive_days"); v != nil {
			request.AliveDays = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().UpdateAutoSnapshotPolicy(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	autoSnapshotPolicyId := d.Id()

	if err := service.DeleteCfsAutoSnapshotPolicyById(ctx, autoSnapshotPolicyId); err != nil {
		return err
	}

	return nil
}
