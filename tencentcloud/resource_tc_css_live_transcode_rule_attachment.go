package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssLiveTranscodeRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCssLiveTranscodeRuleAttachmentRead,
		Create: resourceTencentCloudCssLiveTranscodeRuleAttachmentCreate,
		Delete: resourceTencentCloudCssLiveTranscodeRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "domain name hich you want to bind the transcode template.",
			},

			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "app name which you want to bind, can be empty string if not binding specific app name.",
			},

			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "stream name which you want to bind, can be empty string if not binding specific stream.",
			},

			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "template created by css_live_transcode_template.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "update time.",
			},
		},
	}
}

func resourceTencentCloudCssLiveTranscodeRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveTranscodeRuleRequest()
		domainName string
		appName    string
		streamName string
		templateId int64
	)

	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
		request.StreamName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_id"); ok {
		templateId = (int64)(v.(int))
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveTranscodeRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create css liveTranscodeRuleAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domainName + FILED_SP + appName + FILED_SP + streamName + FILED_SP + helper.Int64ToStr(templateId))
	return resourceTencentCloudCssLiveTranscodeRuleAttachmentRead(d, meta)
}

func resourceTencentCloudCssLiveTranscodeRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]
	templateId := idSplit[3]

	rules, err := service.DescribeCssLiveTranscodeRuleAttachment(ctx, helper.String(domainName), helper.String(templateId))
	if err != nil {
		return err
	}

	if rules == nil {
		d.SetId("")
		return fmt.Errorf("resource `liveTranscodeRuleAttachment` %s does not exist", d.Id())
	}

	for _, v := range rules {
		if *v.DomainName != domainName || *v.AppName != appName || *v.StreamName != streamName || *v.TemplateId != helper.StrToInt64(templateId) {
			log.Printf("[DEBUG]%s api[%s] this rule does not match with:[%s]\n", logId, "query attachment", d.Id())
			continue
		}

		liveTranscodeRuleAttachment := v
		if liveTranscodeRuleAttachment.DomainName != nil {
			_ = d.Set("domain_name", liveTranscodeRuleAttachment.DomainName)
		}

		if liveTranscodeRuleAttachment.AppName != nil {
			_ = d.Set("app_name", liveTranscodeRuleAttachment.AppName)
		}

		if liveTranscodeRuleAttachment.StreamName != nil {
			_ = d.Set("stream_name", liveTranscodeRuleAttachment.StreamName)
		}

		if liveTranscodeRuleAttachment.TemplateId != nil {
			_ = d.Set("template_id", liveTranscodeRuleAttachment.TemplateId)
		}

		if liveTranscodeRuleAttachment.CreateTime != nil {
			_ = d.Set("create_time", liveTranscodeRuleAttachment.CreateTime)
		}

		if liveTranscodeRuleAttachment.UpdateTime != nil {
			_ = d.Set("update_time", liveTranscodeRuleAttachment.UpdateTime)
		}
	}

	return nil
}

func resourceTencentCloudCssLiveTranscodeRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]
	templateId := idSplit[3]

	if err := service.DeleteCssLiveTranscodeRuleAttachmentById(ctx, helper.String(domainName), helper.String(appName), helper.String(streamName), helper.String(templateId)); err != nil {
		return err
	}

	return nil
}
