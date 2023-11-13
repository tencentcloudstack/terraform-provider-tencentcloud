/*
Provides a resource to create a autoscaling auto_scaling_group

Example Usage

```hcl
resource "tencentcloud_autoscaling_auto_scaling_group" "auto_scaling_group" {
  auto_scaling_group_id = "asg-xxxxxxxx"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

autoscaling auto_scaling_group can be imported using the id, e.g.

```
terraform import tencentcloud_autoscaling_auto_scaling_group.auto_scaling_group auto_scaling_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"log"
)

func resourceTencentCloudAutoscalingAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAutoscalingAutoScalingGroupCreate,
		Read:   resourceTencentCloudAutoscalingAutoScalingGroupRead,
		Update: resourceTencentCloudAutoscalingAutoScalingGroupUpdate,
		Delete: resourceTencentCloudAutoscalingAutoScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scaling group ID.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudAutoscalingAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_auto_scaling_group.create")()
	defer inconsistentCheck(d, meta)()

	var autoScalingGroupId string
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		autoScalingGroupId = v.(string)
	}

	d.SetId(autoScalingGroupId)

	return resourceTencentCloudAutoscalingAutoScalingGroupUpdate(d, meta)
}

func resourceTencentCloudAutoscalingAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_auto_scaling_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AutoscalingService{client: meta.(*TencentCloudClient).apiV3Conn}

	autoScalingGroupId := d.Id()

	autoScalingGroup, err := service.DescribeAutoscalingAutoScalingGroupById(ctx, autoScalingGroupId)
	if err != nil {
		return err
	}

	if autoScalingGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AutoscalingAutoScalingGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if autoScalingGroup.AutoScalingGroupId != nil {
		_ = d.Set("auto_scaling_group_id", autoScalingGroup.AutoScalingGroupId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "autoscaling", "autoScalingGroupId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAutoscalingAutoScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_auto_scaling_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := autoscaling.NewEnableAutoScalingGroupRequest()

	autoScalingGroupId := d.Id()

	request.AutoScalingGroupId = &autoScalingGroupId

	immutableArgs := []string{"auto_scaling_group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAutoscalingClient().EnableAutoScalingGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update autoscaling autoScalingGroup failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("autoscaling", "autoScalingGroupId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudAutoscalingAutoScalingGroupRead(d, meta)
}

func resourceTencentCloudAutoscalingAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_auto_scaling_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
