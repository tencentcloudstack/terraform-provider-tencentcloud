package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbCmdTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbCmdTemplateCreate,
		Read:   resourceTencentCloudDasbCmdTemplateRead,
		Update: resourceTencentCloudDasbCmdTemplateUpdate,
		Delete: resourceTencentCloudDasbCmdTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name, maximum length 32 characters, cannot contain blank characters.",
			},
			"cmd_list": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Command list, n separated, maximum length 32768 bytes.",
			},
		},
	}
}

func resourceTencentCloudDasbCmdTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_cmd_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = dasb.NewCreateCmdTemplateRequest()
		response   = dasb.NewCreateCmdTemplateResponse()
		templateId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cmd_list"); ok {
		request.CmdList = helper.String(v.(string))
	}

	request.Encoding = helper.IntUint64(0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateCmdTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.Id == nil {
			e = fmt.Errorf("dasb CmdTemplate not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb CmdTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateIdInt := *response.Response.Id
	templateId = strconv.FormatUint(templateIdInt, 10)
	d.SetId(templateId)

	return resourceTencentCloudDasbCmdTemplateRead(d, meta)
}

func resourceTencentCloudDasbCmdTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_cmd_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		templateId = d.Id()
	)

	CmdTemplate, err := service.DescribeDasbCmdTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if CmdTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbCmdTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if CmdTemplate.Name != nil {
		_ = d.Set("name", CmdTemplate.Name)
	}

	if CmdTemplate.CmdList != nil {
		_ = d.Set("cmd_list", CmdTemplate.CmdList)
	}

	return nil
}

func resourceTencentCloudDasbCmdTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_cmd_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = dasb.NewModifyCmdTemplateRequest()
		templateId = d.Id()
	)

	request.Id = helper.StrToUint64Point(templateId)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cmd_list"); ok {
		request.CmdList = helper.String(v.(string))
	}

	request.Encoding = helper.IntUint64(0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().ModifyCmdTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dasb CmdTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbCmdTemplateRead(d, meta)
}

func resourceTencentCloudDasbCmdTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_cmd_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		templateId = d.Id()
	)

	if err := service.DeleteDasbCmdTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
