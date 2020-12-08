/*
Provides a resource to manage protocol template group.

Example Usage

```hcl
resource "tencentcloud_protocol_template_group" "foo" {
  name                = "group-test"
  protocols = ["ipl-axaf24151","ipl-axaf24152"]
}
```

Import

Protocol template group can be imported using the id, e.g.

```
$ terraform import tencentcloud_protocol_template_group.foo ppmg-0np3u974
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudProtocolTemplateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudProtocolTemplateGroupCreate,
		Read:   resourceTencentCloudProtocolTemplateGroupRead,
		Update: resourceTencentCloudProtocolTemplateGroupUpdate,
		Delete: resourceTencentCloudProtocolTemplateGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the protocol template group.",
			},
			"template_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Service template ID list.",
			},
		},
	}
}

func resourceTencentCloudProtocolTemplateGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template_group.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	name := d.Get("name").(string)
	protocols := d.Get("template_ids").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var outErr, inErr error
	var templateGroupId string

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		templateGroupId, inErr = vpcService.CreateServiceTemplateGroup(ctx, name, protocols)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateGroupId)

	return resourceTencentCloudProtocolTemplateGroupRead(d, meta)
}

func resourceTencentCloudProtocolTemplateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateGroup, has, outErr := vpcService.DescribeServiceTemplateGroupById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			templateGroup, has, inErr = vpcService.DescribeServiceTemplateGroupById(ctx, templateId)
			if inErr != nil {
				return retryError(inErr)
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

	_ = d.Set("name", templateGroup.ServiceTemplateGroupName)
	_ = d.Set("template_ids", templateGroup.ServiceTemplateIdSet)

	return nil
}

func resourceTencentCloudProtocolTemplateGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateGroupId := d.Id()

	if d.HasChange("name") || d.HasChange("template_ids") {
		var outErr, inErr error
		name := d.Get("name").(string)
		templadteIds := d.Get("template_ids").(*schema.Set).List()
		vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyServiceTemplateGroup(ctx, templateGroupId, name, templadteIds)
			if inErr != nil {
				return retryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	return resourceTencentCloudProtocolTemplateGroupRead(d, meta)
}

func resourceTencentCloudProtocolTemplateGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateGroupId := d.Id()
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error

	outErr = vpcService.DeleteServiceTemplateGroup(ctx, templateGroupId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteServiceTemplateGroup(ctx, templateGroupId)
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
		_, has, inErr := vpcService.DescribeServiceTemplateGroupById(ctx, templateGroupId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("protocol template group %s is still exists, retry...", templateGroupId))
		} else {
			return nil
		}
	})

	return outErr
}
