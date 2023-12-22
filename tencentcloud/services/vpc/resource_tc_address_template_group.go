package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudAddressTemplateGroup() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_address_template_group.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	name := d.Get("name").(string)
	addresses := d.Get("template_ids").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var outErr, inErr error
	var templateGroupId string

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		templateGroupId, inErr = vpcService.CreateAddressTemplateGroup(ctx, name, addresses)
		if inErr != nil {
			return tccommon.RetryError(inErr)
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
	defer tccommon.LogElapsed("resource.tencentcloud_address_template_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	templateGroup, has, outErr := vpcService.DescribeAddressTemplateGroupById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			templateGroup, has, inErr = vpcService.DescribeAddressTemplateGroupById(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
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
	defer tccommon.LogElapsed("resource.tencentcloud_address_template_group.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateGroupId := d.Id()

	if d.HasChange("name") || d.HasChange("template_ids") {
		var outErr, inErr error
		name := d.Get("name").(string)
		templadteIds := d.Get("template_ids").(*schema.Set).List()
		vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyAddressTemplateGroup(ctx, templateGroupId, name, templadteIds)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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
	defer tccommon.LogElapsed("resource.tencentcloud_address_template_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateGroupId := d.Id()
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error

	outErr = vpcService.DeleteAddressTemplateGroup(ctx, templateGroupId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteAddressTemplateGroup(ctx, templateGroupId)
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
		_, has, inErr := vpcService.DescribeAddressTemplateGroupById(ctx, templateGroupId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("address template group %s is still exists, retry...", templateGroupId))
		} else {
			return nil
		}
	})

	return outErr
}
