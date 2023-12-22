package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudAddressTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAddressTemplateCreate,
		Read:   resourceTencentCloudAddressTemplateRead,
		Update: resourceTencentCloudAddressTemplateUpdate,
		Delete: resourceTencentCloudAddressTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the address template.",
			},
			"addresses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Address list. IP(`10.0.0.1`), CIDR(`10.0.1.0/24`), IP range(`10.0.0.1-10.0.0.100`) format are supported.",
			},
		},
	}
}

func resourceTencentCloudAddressTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_template.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	name := d.Get("name").(string)
	addresses := d.Get("addresses").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var outErr, inErr error
	var templateId string

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		templateId, inErr = vpcService.CreateAddressTemplate(ctx, name, addresses)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateId)

	return resourceTencentCloudAddressTemplateRead(d, meta)
}

func resourceTencentCloudAddressTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	template, has, outErr := vpcService.DescribeAddressTemplateById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			template, has, inErr = vpcService.DescribeAddressTemplateById(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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

	_ = d.Set("name", template.AddressTemplateName)
	_ = d.Set("addresses", template.AddressSet)

	return nil
}

func resourceTencentCloudAddressTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_template.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("addresses") {
		var outErr, inErr error
		name := d.Get("name").(string)
		addresses := d.Get("addresses").(*schema.Set).List()
		vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyAddressTemplate(ctx, templateId, name, addresses)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	return resourceTencentCloudAddressTemplateRead(d, meta)
}

func resourceTencentCloudAddressTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error

	outErr = vpcService.DeleteAddressTemplate(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteAddressTemplate(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	//check not exist
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := vpcService.DescribeAddressTemplateById(ctx, templateId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("address template %s is still exists, retry...", templateId))
		} else {
			return nil
		}
	})

	return outErr
}
