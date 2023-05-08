/*
Provides a resource to attach or detach CVM instances to a specified scaling group.

Example Usage

```hcl
resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = "sg-afasfa"
  instance_ids     = ["ins-01", "ins-02"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsAttachmentCreate,
		Read:   resourceTencentCloudAsAttachmentRead,
		Update: resourceTencentCloudAsAttachmentUpdate,
		Delete: resourceTencentCloudAsAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"instance_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID list of CVM instances to be attached to the scaling group.",
			},
		},
	}
}

func resourceTencentCloudAsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Get("scaling_group_id").(string)
	instanceIds := helper.InterfacesStrings(d.Get("instance_ids").(*schema.Set).List())
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.AttachInstances(ctx, scalingGroupId, instanceIds)
	if err != nil {
		return err
	}
	d.SetId(scalingGroupId)

	return resourceTencentCloudAsAttachmentRead(d, meta)
}

func resourceTencentCloudAsAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instanceIds []string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, errRet := asService.DescribeAutoScalingAttachment(ctx, scalingGroupId, false)
		if errRet != nil {
			return retryError(errRet)
		}
		instanceIds = result
		return nil
	})
	if err != nil {
		return err
	}
	if len(instanceIds) < 1 {
		d.SetId("")
		return nil
	}
	_ = d.Set("instance_ids", instanceIds)
	return nil
}

func resourceTencentCloudAsAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_attachment.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Id()
	if d.HasChange("instance_ids") {
		oldInterface, newInterface := d.GetChange("instance_ids")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)
		remove := helper.InterfacesStrings(oldInstances.Difference(newInstances).List())
		add := helper.InterfacesStrings(newInstances.Difference(oldInstances).List())

		asService := AsService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		if len(add) > 0 {
			err := asService.AttachInstances(ctx, scalingGroupId, add)
			if err != nil {
				return err
			}
		}
		if len(remove) > 0 {
			err := asService.DetachInstances(ctx, scalingGroupId, remove)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceTencentCloudAsAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instanceIds []string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, errRet := asService.DescribeAutoScalingAttachment(ctx, scalingGroupId, false)
		if errRet != nil {
			return retryError(errRet)
		}
		instanceIds = result
		return nil
	})
	if err != nil {
		return err
	}
	if len(instanceIds) < 1 {
		return nil
	}

	err = asService.DetachInstances(ctx, scalingGroupId, instanceIds)
	if err != nil {
		return err
	}
	return nil
}
