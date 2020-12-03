/*
Provides a resource to manage service template.

Example Usage

```hcl
resource "tencentcloud_service_template" "foo" {
  name                = "service-template-test"
  services = ["tcp:80","udp:all","icmp:10-30"]
}
```

Import

CAM user can be imported using the service template, e.g.

```
$ terraform import tencentcloud_service_template.foo ppm-nwrggd14
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudServiceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudServiceTemplateCreate,
		Read:   resourceTencentCloudServiceTemplateRead,
		Update: resourceTencentCloudServiceTemplateUpdate,
		Delete: resourceTencentCloudServiceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the service template.",
			},
			"services": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateLowCase,
				},
				Required:    true,
				Description: "Service list. Valid protocols are  `tcp`, `udp`, `icmp`, `gre`. Single port(tcp:80), multi-port(tcp:80,443), port range(tcp:3306-20000), all(tcp:all) format are support. Protocol `icmp` and `gre` cannot specify port.",
			},
		},
	}
}

func resourceTencentCloudServiceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	name := d.Get("name").(string)
	services := d.Get("services").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var outErr, inErr error
	var templateId string

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		templateId, inErr = vpcService.CreateServiceTemplate(ctx, name, services)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateId)

	return resourceTencentCloudServiceTemplateRead(d, meta)
}

func resourceTencentCloudServiceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	template, has, outErr := vpcService.DescribeServiceTemplateById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			template, has, inErr = vpcService.DescribeServiceTemplateById(ctx, templateId)
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

	_ = d.Set("name", template.ServiceTemplateName)
	_ = d.Set("services", template.ServiceSet)

	return nil
}

func resourceTencentCloudServiceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("services") {
		var outErr, inErr error
		name := d.Get("name").(string)
		services := d.Get("services").(*schema.Set).List()
		vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyServiceTemplate(ctx, templateId, name, services)
			if inErr != nil {
				return retryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("name")
		d.SetPartial("services")
	}

	return resourceTencentCloudServiceTemplateRead(d, meta)
}

func resourceTencentCloudServiceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_service_template.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateId := d.Id()
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error

	outErr = vpcService.DeleteServiceTemplate(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteServiceTemplate(ctx, templateId)
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
		_, has, inErr := vpcService.DescribeServiceTemplateById(ctx, templateId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("service template %s is still exists, retry...", templateId))
		} else {
			return nil
		}
	})

	return outErr
}
