package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssRecordRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssRecordRuleAttachmentCreate,
		Read:   resourceTencentCloudCssRecordRuleAttachmentRead,
		Delete: resourceTencentCloudCssRecordRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Streaming domain name.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Template ID.",
			},

			"app_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The streaming path is consistent with the AppName in the streaming and playback addresses. The default is live.",
			},

			"stream_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream name. Note: If this parameter is set to a non empty string, the rule will only work on this streaming.",
			},
		},
	}
}

func resourceTencentCloudCssRecordRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_record_rule_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveRecordRuleRequest()
		templateId int
		domainName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		request.StreamName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveRecordRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css recordRule failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{helper.IntToStr(templateId), domainName}, FILED_SP))

	return resourceTencentCloudCssRecordRuleAttachmentRead(d, meta)
}

func resourceTencentCloudCssRecordRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_record_rule_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}
	domainName := idSplit[1]

	recordRule, err := service.DescribeCssRecordRuleById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	if recordRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssRecordRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if recordRule.DomainName != nil {
		_ = d.Set("domain_name", recordRule.DomainName)
	}

	if recordRule.TemplateId != nil {
		_ = d.Set("template_id", recordRule.TemplateId)
	}

	if recordRule.AppName != nil {
		_ = d.Set("app_name", recordRule.AppName)
	}

	if recordRule.StreamName != nil {
		_ = d.Set("stream_name", recordRule.StreamName)
	}

	return nil
}

func resourceTencentCloudCssRecordRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_record_rule_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)
	domainName := idSplit[1]

	recordRule, err := service.DescribeCssRecordRuleById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	appName := ""
	if recordRule.AppName != nil {
		appName = *recordRule.AppName
	}

	streamName := ""
	if recordRule.StreamName != nil {
		streamName = *recordRule.StreamName
	}

	if err := service.DeleteCssRecordRuleById(ctx, domainName, appName, streamName); err != nil {
		return err
	}

	return nil
}
