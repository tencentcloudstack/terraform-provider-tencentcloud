package ses

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSesTemplateRead,
		Create: resourceTencentCloudSesTemplateCreate,
		Update: resourceTencentCloudSesTemplateUpdate,
		Delete: resourceTencentCloudSesTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "smsTemplateName, which must be required.",
			},

			"template_content": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Sms Template Content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Html code after base64.",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text content after base64.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSesTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ses.NewCreateEmailTemplateRequest()
		response   *ses.CreateEmailTemplateResponse
		templateId uint64
	)

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
		prometheusTemp := ses.TemplateContent{}
		if v, ok := dMap["html"]; ok {
			html := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Html = helper.String(html)
		}
		if v, ok := dMap["text"]; ok {
			text := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Text = helper.String(text)
		}
		request.TemplateContent = &prometheusTemp
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateEmailTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses template failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateID

	d.SetId(strconv.FormatUint(templateId, 10))
	return resourceTencentCloudSesTemplateRead(d, meta)
}

func resourceTencentCloudSesTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	templateID := d.Id()
	templateId, ee := strconv.Atoi(templateID)
	if ee != nil {
		return ee
	}

	templateResponse, err := service.DescribeSesTemplate(ctx, uint64(templateId))
	if err != nil {
		return err
	}

	if templateResponse == nil {
		d.SetId("")
		return fmt.Errorf("resource `template` %v does not exist", templateId)
	}

	template := templateResponse.TemplateContent
	if template != nil {
		var templateContents []map[string]interface{}
		templateContent := map[string]interface{}{}
		if template.Html != nil {
			html, err := base64.StdEncoding.DecodeString(*template.Html)
			if err != nil {
				return err
			}
			contentHtml := string(html)
			templateContent["html"] = &contentHtml
		}
		if template.Text != nil {
			text, err := base64.StdEncoding.DecodeString(*template.Text)
			if err != nil {
				return err
			}
			contentText := string(text)
			templateContent["text"] = &contentText
		}
		templateContents = append(templateContents, templateContent)
		err = d.Set("template_content", templateContents)
		if ee != nil {
			return fmt.Errorf("set template_content err: %v", err)
		}
	}

	_ = d.Set("template_name", templateResponse.TemplateName)

	return nil
}

func resourceTencentCloudSesTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := ses.NewUpdateEmailTemplateRequest()

	templateID, ee := strconv.Atoi(d.Id())
	if ee != nil {
		return ee
	}
	templateId := uint64(templateID)
	request.TemplateID = &templateId

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
		prometheusTemp := ses.TemplateContent{}
		if v, ok := dMap["html"]; ok {
			html := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Html = helper.String(html)
		}
		if v, ok := dMap["text"]; ok {
			text := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Text = helper.String(text)
		}
		request.TemplateContent = &prometheusTemp
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().UpdateEmailTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses template failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSesTemplateRead(d, meta)
}

func resourceTencentCloudSesTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	templateID, ee := strconv.Atoi(d.Id())
	if ee != nil {
		return ee
	}
	templateId := uint64(templateID)

	if err := service.DeleteSesTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
