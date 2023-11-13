/*
Provides a resource to create a autoscaling scale_in_instances

Example Usage

```hcl
resource "tencentcloud_autoscaling_scale_in_instances" "scale_in_instances" {
  auto_scaling_group_id = "asg-xxxxxxxx"
  scale_in_number = 1
}
```

Import

autoscaling scale_in_instances can be imported using the id, e.g.

```
terraform import tencentcloud_autoscaling_scale_in_instances.scale_in_instances scale_in_instances_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudAutoscalingScaleInInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAutoscalingScaleInInstancesCreate,
		Read:   resourceTencentCloudAutoscalingScaleInInstancesRead,
		Delete: resourceTencentCloudAutoscalingScaleInInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scaling group ID.",
			},

			"scale_in_number": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of instances to be reduced.",
			},
		},
	}
}

func resourceTencentCloudAutoscalingScaleInInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_scale_in_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = autoscaling.NewScaleInInstancesRequest()
		response           = autoscaling.NewScaleInInstancesResponse()
		autoScalingGroupId string
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		autoScalingGroupId = v.(string)
		request.AutoScalingGroupId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("scale_in_number"); v != nil {
		request.ScaleInNumber = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAutoscalingClient().ScaleInInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate autoscaling scaleInInstances failed, reason:%+v", logId, err)
		return err
	}

	autoScalingGroupId = *response.Response.AutoScalingGroupId
	d.SetId(autoScalingGroupId)

	return resourceTencentCloudAutoscalingScaleInInstancesRead(d, meta)
}

func resourceTencentCloudAutoscalingScaleInInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_scale_in_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAutoscalingScaleInInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_scale_in_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
