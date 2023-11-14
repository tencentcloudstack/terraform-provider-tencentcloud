/*
Provides a resource to create a clb instance_sla_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_sla_config" "instance_sla_config" {
  load_balancer_sla {
		load_balancer_id = "lb-xxxxxxxx"
		sla_type = ""

  }
}
```

Import

clb instance_sla_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_sla_config.instance_sla_config instance_sla_config_id
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

func resourceTencentCloudClbInstanceSlaConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceSlaConfigCreate,
		Read:   resourceTencentCloudClbInstanceSlaConfigRead,
		Update: resourceTencentCloudClbInstanceSlaConfigUpdate,
		Delete: resourceTencentCloudClbInstanceSlaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_sla": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "CLB instance information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the CLB instance.",
						},
						"sla_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "To upgrade to LCU-supported CLB instances. It must be SLA.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbInstanceSlaConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId()

	return resourceTencentCloudClbInstanceSlaConfigUpdate(d, meta)
}

func resourceTencentCloudClbInstanceSlaConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceSlaConfigId := d.Id()

	instanceSlaConfig, err := service.DescribeClbInstanceSlaConfigById(ctx, slaUpdateParam)
	if err != nil {
		return err
	}

	if instanceSlaConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbInstanceSlaConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceSlaConfig.LoadBalancerSla != nil {
		loadBalancerSlaList := []interface{}{}
		for _, loadBalancerSla := range instanceSlaConfig.LoadBalancerSla {
			loadBalancerSlaMap := map[string]interface{}{}

			if instanceSlaConfig.LoadBalancerSla.LoadBalancerId != nil {
				loadBalancerSlaMap["load_balancer_id"] = instanceSlaConfig.LoadBalancerSla.LoadBalancerId
			}

			if instanceSlaConfig.LoadBalancerSla.SlaType != nil {
				loadBalancerSlaMap["sla_type"] = instanceSlaConfig.LoadBalancerSla.SlaType
			}

			loadBalancerSlaList = append(loadBalancerSlaList, loadBalancerSlaMap)
		}

		_ = d.Set("load_balancer_sla", loadBalancerSlaList)

	}

	return nil
}

func resourceTencentCloudClbInstanceSlaConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := clb.NewModifyLoadBalancerSlaRequest()

	instanceSlaConfigId := d.Id()

	request.SlaUpdateParam = &slaUpdateParam

	immutableArgs := []string{"load_balancer_sla"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerSla(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb instanceSlaConfig failed, reason:%+v", logId, err)
		return err
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 60*readRetryTimeout, time.Second, service.ClbInstanceSlaConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudClbInstanceSlaConfigRead(d, meta)
}

func resourceTencentCloudClbInstanceSlaConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
