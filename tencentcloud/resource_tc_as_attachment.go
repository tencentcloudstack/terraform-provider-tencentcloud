package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudAsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsAttachmentCreate,
		Read:   resourceTencentCloudAsAttachmentRead,
		Update: resourceTencentCloudAsAttachmentUpdate,
		Delete: resourceTencentCloudAsAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTencentCloudAsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Get("scaling_group_id").(string)
	instanceIds := expandStringList(d.Get("instance_ids").(*schema.Set).List())
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceIds, err := asService.DescribeAutoScalingAttachment(ctx, scalingGroupId)
	if err != nil {
		return err
	}
	d.Set("instance_ids", instanceIds)

	return nil
}
func resourceTencentCloudAsAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Id()
	if d.HasChange("instance_ids") {
		old, new := d.GetChange("instance_ids")
		oldInstances := old.(*schema.Set)
		newInstances := new.(*schema.Set)
		remove := expandStringList(oldInstances.Difference(newInstances).List())
		add := expandStringList(newInstances.Difference(oldInstances).List())

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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceIds, err := asService.DescribeAutoScalingAttachment(ctx, scalingGroupId)
	if err != nil {
		return err
	}

	err = asService.DetachInstances(ctx, scalingGroupId, instanceIds)
	if err != nil {
		return err
	}
	return nil
}
