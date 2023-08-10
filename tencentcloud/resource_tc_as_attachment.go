/*
Provides a resource to attach or detach CVM instances to a specified scaling group.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones_by_product.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_instance" "example" {
  instance_name              = "tf_example_instance"
  availability_zone          = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  image_id                   = data.tencentcloud_images.image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = tencentcloud_as_scaling_group.example.id
  instance_ids     = [tencentcloud_instance.example.id]
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
