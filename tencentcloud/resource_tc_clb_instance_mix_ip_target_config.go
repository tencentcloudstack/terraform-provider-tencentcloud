/*
Provides a resource to create a clb instance_mix_ip_target_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_mix_ip_target_config" "instance_mix_ip_target_config" {
  load_balancer_id = "lb-5dnrkgry"
  mix_ip_target = false
}
```

Import

clb instance_mix_ip_target_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config instance_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbInstanceMixIpTargetConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceMixIpTargetConfigCreate,
		Read:   resourceTencentCloudClbInstanceMixIpTargetConfigRead,
		Update: resourceTencentCloudClbInstanceMixIpTargetConfigUpdate,
		Delete: resourceTencentCloudClbInstanceMixIpTargetConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of CLB instances to be queried.",
			},

			"mix_ip_target": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "False: closed True:open.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceMixIpTargetConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.create")()
	defer inconsistentCheck(d, meta)()

	lbId := d.Get("load_balancer_id").(string)

	d.SetId(lbId)

	return resourceTencentCloudClbInstanceMixIpTargetConfigUpdate(d, meta)
}

func resourceTencentCloudClbInstanceMixIpTargetConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	lbId := d.Id()

	instance, err := service.DescribeLoadBalancerById(ctx, lbId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.LoadBalancerId != nil {
		_ = d.Set("load_balancer_id", instance.LoadBalancerId)
	}

	if instance.MixIpTarget != nil {
		_ = d.Set("mix_ip_target", instance.MixIpTarget)
	}

	return nil
}

func resourceTencentCloudClbInstanceMixIpTargetConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := clb.NewModifyLoadBalancerMixIpTargetRequest()

	lbId := d.Id()

	request.LoadBalancerIds = []*string{&lbId}

	if v, ok := d.GetOkExists("mix_ip_target"); ok {
		request.MixIpTarget = helper.Bool(v.(bool))
	}

	var taskId string
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerMixIpTarget(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		taskId = *result.Response.RequestId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb instanceMixIpTargetConfig failed, reason:%+v", logId, err)
		return err
	}

	retryErr := waitForTaskFinish(taskId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
	if retryErr != nil {
		return retryErr
	}

	return resourceTencentCloudClbInstanceMixIpTargetConfigRead(d, meta)
}

func resourceTencentCloudClbInstanceMixIpTargetConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
