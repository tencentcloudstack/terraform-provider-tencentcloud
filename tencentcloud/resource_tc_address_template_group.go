/*
Provides a resource to manage address template group.

Example Usage

```hcl
resource "tencentcloud_address_template_group" "foo" {
  name                = "group-test"
  template_ids = ["ipl-axaf24151","ipl-axaf24152"]
}
```

Import

Address template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template_group.foo ipmg-0np3u974
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudAddressTemplateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAddressTemplateGroupCreate,
		Read:   resourceTencentCloudAddressTemplateGroupRead,
		Update: resourceTencentCloudAddressTemplateGroupUpdate,
		Delete: resourceTencentCloudAddressTemplateGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the address template group.",
			},
			"template_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Template ID list.",
			},
		},
	}
}

func resourceTencentCloudAddressTemplateGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_address_template_group.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	name := d.Get("name").(string)
	addresses := d.Get("template_ids").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var outErr, inErr error
	var templateGroupId string

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		templateGroupId, inErr = vpcService.CreateAddressTemplateGroup(ctx, name, addresses)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateGroupId)

	return resourceTencentCloudAddressTemplateGroupRead(d, meta)
}

func resourceTencentCloudAddressTemplateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_address_template_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateGroup, has, outErr := vpcService.DescribeAddressTemplateGroupById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			templateGroup, has, inErr = vpcService.DescribeAddressTemplateGroupById(ctx, templateId)
			if inErr != nil {
				return retryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", templateGroup.AddressTemplateGroupName)
	_ = d.Set("template_ids", templateGroup.AddressTemplateIdSet)

	return nil
}

func resourceTencentCloudAddressTemplateGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_address_template_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateGroupId := d.Id()

	if d.HasChange("name") || d.HasChange("template_ids") {
		var outErr, inErr error
		name := d.Get("name").(string)
		templadteIds := d.Get("template_ids").(*schema.Set).List()
		vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyAddressTemplateGroup(ctx, templateGroupId, name, templadteIds)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudAddressTemplateGroupRead(d, meta)
}

func resourceTencentCloudAddressTemplateGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_address_template_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateGroupId := d.Id()
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error

	outErr = vpcService.DeleteAddressTemplateGroup(ctx, templateGroupId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteAddressTemplateGroup(ctx, templateGroupId)
			if inErr != nil {
				return retryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	//check not exist
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := vpcService.DescribeAddressTemplateGroupById(ctx, templateGroupId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("address template group %s is still exists, retry...", templateGroupId))
		} else {
			return nil
		}
	})

	return outErr
}
