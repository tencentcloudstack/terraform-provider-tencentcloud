/*
Provides a resource to create a css live_transcode_rule

Example Usage

```hcl
resource "tencentcloud_css_live_transcode_rule" "live_transcode_rule" {
  domain_name = &lt;nil&gt;
  app_name = ""
  stream_name = ""
  template_id = &lt;nil&gt;
    }
```

Import

css live_transcode_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_live_transcode_rule.live_transcode_rule live_transcode_rule_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCssLiveTranscodeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssLiveTranscodeRuleCreate,
		Read:   resourceTencentCloudCssLiveTranscodeRuleRead,
		Delete: resourceTencentCloudCssLiveTranscodeRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name hich you want to bind the transcode template.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "App name which you want to bind, can be empty string if not binding specific app name.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream name which you want to bind, can be empty string if not binding specific stream.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Template created by css_live_transcode_template.",
			},

			"create_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Create time.",
			},

			"update_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudCssLiveTranscodeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveTranscodeRuleRequest()
		response   = css.NewCreateLiveTranscodeRuleResponse()
		domainName string
		appName    string
		streamName string
		templateId int
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

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int64)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveTranscodeRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css liveTranscodeRule failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(strings.Join([]string{domainName, appName, streamName, helper.Int64ToStr(templateId)}, FILED_SP))

	return resourceTencentCloudCssLiveTranscodeRuleRead(d, meta)
}

func resourceTencentCloudCssLiveTranscodeRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule.read")()
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

	liveTranscodeRule, err := service.DescribeCssLiveTranscodeRuleById(ctx, domainName, appName, streamName, templateId)
	if err != nil {
		return err
	}

	if liveTranscodeRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssLiveTranscodeRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if liveTranscodeRule.DomainName != nil {
		_ = d.Set("domain_name", liveTranscodeRule.DomainName)
	}

	if liveTranscodeRule.AppName != nil {
		_ = d.Set("app_name", liveTranscodeRule.AppName)
	}

	if liveTranscodeRule.StreamName != nil {
		_ = d.Set("stream_name", liveTranscodeRule.StreamName)
	}

	if liveTranscodeRule.TemplateId != nil {
		_ = d.Set("template_id", liveTranscodeRule.TemplateId)
	}

	if liveTranscodeRule.CreateTime != nil {
		_ = d.Set("create_time", liveTranscodeRule.CreateTime)
	}

	if liveTranscodeRule.UpdateTime != nil {
		_ = d.Set("update_time", liveTranscodeRule.UpdateTime)
	}

	return nil
}

func resourceTencentCloudCssLiveTranscodeRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_rule.delete")()
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

	if err := service.DeleteCssLiveTranscodeRuleById(ctx, domainName, appName, streamName, templateId); err != nil {
		return err
	}

	return nil
}
