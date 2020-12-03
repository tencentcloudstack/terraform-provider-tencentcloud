/*
Provides a resource to manage service template group.

Example Usage

```hcl
resource "tencentcloud_service_template_group" "foo" {
  name                = "group-test"
  services = ["ipl-axaf24151","ipl-axaf24152"]
}
```

Import

CAM user can be imported using the service template, e.g.

```
$ terraform import tencentcloud_service_template.foo ppmg-0np3u974
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudServiceTemplateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudServiceTemplateGroupCreate,
		Read:   resourceTencentCloudServiceTemplateGroupRead,
		Update: resourceTencentCloudServiceTemplateGroupUpdate,
		Delete: resourceTencentCloudServiceTemplateGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the service template group.",
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

func resourceTencentCloudServiceTemplateGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	name := d.Get("name").(string)
	servicees := d.Get("template_ids").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var outErr, inErr error
	var templateGroupId string

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		templateGroupId, inErr = vpcService.CreateServiceTemplateGroup(ctx, name, servicees)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateGroupId)

	return resourceTencentCloudServiceTemplateGroupRead(d, meta)
}

func resourceTencentCloudServiceTemplateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.read")()
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

func resourceTencentCloudServiceTemplateGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.update")()

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

		d.SetPartial("name")
		d.SetPartial("templadte_ids")
	}

	return resourceTencentCloudServiceTemplateGroupRead(d, meta)
}

func resourceTencentCloudServiceTemplateGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.delete")()

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
			return resource.RetryableError(fmt.Errorf("service template group %s is still exists, retry...", templateGroupId))
		} else {
			return nil
		}
	})

	return outErr
}
