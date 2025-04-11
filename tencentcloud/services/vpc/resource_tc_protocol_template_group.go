package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudProtocolTemplateGroup() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template_group.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	name := d.Get("name").(string)
	protocols := d.Get("template_ids").(*schema.Set).List()

	vpcService := VpcService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var outErr, inErr error
	var templateGroupId string

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		templateGroupId, inErr = vpcService.CreateServiceTemplateGroup(ctx, name, protocols)
		if inErr != nil {
			return tccommon.RetryError(inErr)
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
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	templateGroup, has, outErr := vpcService.DescribeServiceTemplateGroupById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			templateGroup, has, inErr = vpcService.DescribeServiceTemplateGroupById(ctx, templateId)
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

	_ = d.Set("name", templateGroup.ServiceTemplateGroupName)
	_ = d.Set("template_ids", templateGroup.ServiceTemplateIdSet)

	return nil
}

func resourceTencentCloudProtocolTemplateGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template_group.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateGroupId := d.Id()

	if d.HasChange("name") || d.HasChange("template_ids") {
		var outErr, inErr error
		name := d.Get("name").(string)
		templadteIds := d.Get("template_ids").(*schema.Set).List()
		vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyServiceTemplateGroup(ctx, templateGroupId, name, templadteIds)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
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
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateGroupId := d.Id()
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error

	outErr = vpcService.DeleteServiceTemplateGroup(ctx, templateGroupId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteServiceTemplateGroup(ctx, templateGroupId)
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
		_, has, inErr := vpcService.DescribeServiceTemplateGroupById(ctx, templateGroupId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("protocol template group %s is still exists, retry...", templateGroupId))
		} else {
			return nil
		}
	})

	return outErr
}
