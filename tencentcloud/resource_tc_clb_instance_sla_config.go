/*
Provides a resource to create a clb instance_sla_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_sla_config" "instance_sla_config" {
  load_balancer_id = "lb-5dnrkgry"
  sla_type         = "SLA"
}
```

Import

clb instance_sla_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_sla_config.instance_sla_config instance_id
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
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the CLB instance.",
			},
			"sla_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "This parameter is required to create LCU-supported instances. Values:" +
					"`SLA`: Super Large 4. When you have activated Super Large models, `SLA` refers to Super Large 4; " +
					"`clb.c2.medium`: Standard; " +
					"`clb.c3.small`: Advanced 1; " +
					"`clb.c3.medium`: Advanced 1; " +
					"`clb.c4.small`: Super Large 1; " +
					"`clb.c4.medium`: Super Large 2; " +
					"`clb.c4.large`: Super Large 3; " +
					"`clb.c4.xlarge`: Super Large 4. " +
					"For more details, see [Instance Specifications](https://intl.cloud.tencent.com/document/product/214/84689?from_cn_redirect=1).",
			},
		},
	}
}

func resourceTencentCloudClbInstanceSlaConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.create")()
	defer inconsistentCheck(d, meta)()

	lbId := d.Get("load_balancer_id").(string)
	d.SetId(lbId)

	return resourceTencentCloudClbInstanceSlaConfigUpdate(d, meta)
}

func resourceTencentCloudClbInstanceSlaConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.read")()
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
		log.Printf("[WARN]%s resource `ClbInstanceSlaConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.LoadBalancerId != nil {
		_ = d.Set("load_balancer_id", instance.LoadBalancerId)
	}

	if instance.SlaType != nil {
		_ = d.Set("sla_type", instance.SlaType)
	}

	return nil
}

func resourceTencentCloudClbInstanceSlaConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := clb.NewModifyLoadBalancerSlaRequest()

	lbId := d.Id()

	param := clb.SlaUpdateParam{}
	param.LoadBalancerId = &lbId
	param.SlaType = helper.String(d.Get("sla_type").(string))

	request.LoadBalancerSla = []*clb.SlaUpdateParam{&param}

	var taskId string
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerSla(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		taskId = *result.Response.RequestId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb instanceSlaConfig failed, reason:%+v", logId, err)
		return err
	}

	retryErr := waitForTaskFinish(taskId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
	if retryErr != nil {
		return retryErr
	}

	return resourceTencentCloudClbInstanceSlaConfigRead(d, meta)
}

func resourceTencentCloudClbInstanceSlaConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance_sla_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
