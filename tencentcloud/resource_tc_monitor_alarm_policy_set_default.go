package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorAlarmPolicySetDefault() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorAlarmPolicySetDefaultCreate,
		Read:   resourceTencentCloudMonitorAlarmPolicySetDefaultRead,
		Delete: resourceTencentCloudMonitorAlarmPolicySetDefaultDelete,

		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Fixed value, as `monitor`.",
			},

			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Policy id.",
			},
		},
	}
}

func resourceTencentCloudMonitorAlarmPolicySetDefaultCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy_set_default.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewSetDefaultAlarmPolicyRequest()
		policyId string
	)
	if v, ok := d.GetOk("module"); ok {
		request.Module = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_id"); ok {
		policyId = v.(string)
		request.PolicyId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().SetDefaultAlarmPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate monitor policySetDefault failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(policyId)

	return resourceTencentCloudMonitorAlarmPolicySetDefaultRead(d, meta)
}

func resourceTencentCloudMonitorAlarmPolicySetDefaultRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy_set_default.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMonitorAlarmPolicySetDefaultDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy_set_default.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
