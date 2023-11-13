/*
Provides a resource to create a clb instance_mix_ip_target_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_mix_ip_target_config" "instance_mix_ip_target_config" {
  load_balancer_ids =
  mix_ip_target =
}
```

Import

clb instance_mix_ip_target_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config instance_mix_ip_target_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"log"
	"time"
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
			"load_balancer_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of IDs of CLB instances to be queried.",
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

	d.SetId()

	return resourceTencentCloudClbInstanceMixIpTargetConfigUpdate(d, meta)
}

func resourceTencentCloudClbInstanceMixIpTargetConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceMixIpTargetConfigId := d.Id()

	instanceMixIpTargetConfig, err := service.DescribeClbInstanceMixIpTargetConfigById(ctx, loadBalancerIds)
	if err != nil {
		return err
	}

	if instanceMixIpTargetConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbInstanceMixIpTargetConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceMixIpTargetConfig.LoadBalancerIds != nil {
		_ = d.Set("load_balancer_ids", instanceMixIpTargetConfig.LoadBalancerIds)
	}

	if instanceMixIpTargetConfig.MixIpTarget != nil {
		_ = d.Set("mix_ip_target", instanceMixIpTargetConfig.MixIpTarget)
	}

	return nil
}

func resourceTencentCloudClbInstanceMixIpTargetConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := clb.NewModifyLoadBalancerMixIpTargetRequest()

	instanceMixIpTargetConfigId := d.Id()

	request.LoadBalancerIds = &loadBalancerIds

	immutableArgs := []string{"load_balancer_ids", "mix_ip_target"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerMixIpTarget(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb instanceMixIpTargetConfig failed, reason:%+v", logId, err)
		return err
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 60*readRetryTimeout, time.Second, service.ClbInstanceMixIpTargetConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudClbInstanceMixIpTargetConfigRead(d, meta)
}

func resourceTencentCloudClbInstanceMixIpTargetConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_mix_ip_target_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
