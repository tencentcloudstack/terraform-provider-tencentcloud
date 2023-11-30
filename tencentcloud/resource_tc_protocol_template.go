package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudProtocolTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudProtocolTemplateCreate,
		Read:   resourceTencentCloudProtocolTemplateRead,
		Update: resourceTencentCloudProtocolTemplateUpdate,
		Delete: resourceTencentCloudProtocolTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the protocol template.",
			},
			"protocols": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateLowCase,
				},
				Required:    true,
				Description: "Protocol list. Valid protocols are  `tcp`, `udp`, `icmp`, `gre`. Single port(tcp:80), multi-port(tcp:80,443), port range(tcp:3306-20000), all(tcp:all) format are support. Protocol `icmp` and `gre` cannot specify port.",
			},
		},
	}
}

func resourceTencentCloudProtocolTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	name := d.Get("name").(string)
	protocols := d.Get("protocols").(*schema.Set).List()

	vpcProtocol := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var outErr, inErr error
	var templateId string

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		templateId, inErr = vpcProtocol.CreateServiceTemplate(ctx, name, protocols)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateId)

	return resourceTencentCloudProtocolTemplateRead(d, meta)
}

func resourceTencentCloudProtocolTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcProtocol := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	template, has, outErr := vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			template, has, inErr = vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
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
	_ = d.Set("protocols", template.ServiceSet)

	return nil
}

func resourceTencentCloudProtocolTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("protocols") {
		var outErr, inErr error
		name := d.Get("name").(string)
		protocols := d.Get("protocols").(*schema.Set).List()
		vpcProtocol := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcProtocol.ModifyServiceTemplate(ctx, templateId, name, protocols)
			if inErr != nil {
				return retryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudProtocolTemplateRead(d, meta)
}

func resourceTencentCloudProtocolTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_protocol_template.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	templateId := d.Id()
	vpcProtocol := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error

	outErr = vpcProtocol.DeleteServiceTemplate(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = vpcProtocol.DeleteServiceTemplate(ctx, templateId)
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
		_, has, inErr := vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("protocol template %s is still exists, retry...", templateId))
		} else {
			return nil
		}
	})

	return outErr
}
