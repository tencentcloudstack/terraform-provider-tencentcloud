/*
Provides a resource to create a css watermark_rule

Example Usage

```hcl
resource "tencentcloud_css_watermark_rule" "watermark_rule" {
  domain_name = &lt;nil&gt;
  app_name = &lt;nil&gt;
  stream_name = &lt;nil&gt;
  watermark_id = &lt;nil&gt;
    }
```

Import

css watermark_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_watermark_rule.watermark_rule watermark_rule_id
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

func resourceTencentCloudCssWatermarkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssWatermarkRuleCreate,
		Read:   resourceTencentCloudCssWatermarkRuleRead,
		Delete: resourceTencentCloudCssWatermarkRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Rule domain name.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Rule app name.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Rule stream name.",
			},

			"watermark_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Watermark id created by AddLiveWatermark.",
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

func resourceTencentCloudCssWatermarkRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = css.NewCreateLiveWatermarkRuleRequest()
		response    = css.NewCreateLiveWatermarkRuleResponse()
		domainName  string
		appName     string
		streamName  string
		watermarkId string
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

	if v, ok := d.GetOkExists("watermark_id"); ok {
		watermarkId = v.(int64)
		request.WatermarkId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveWatermarkRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css watermarkRule failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(strings.Join([]string{domainName, appName, streamName, watermarkId}, FILED_SP))

	return resourceTencentCloudCssWatermarkRuleRead(d, meta)
}

func resourceTencentCloudCssWatermarkRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule.read")()
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
	watermarkId := idSplit[3]

	watermarkRule, err := service.DescribeCssWatermarkRuleById(ctx, domainName, appName, streamName, watermarkId)
	if err != nil {
		return err
	}

	if watermarkRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssWatermarkRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if watermarkRule.DomainName != nil {
		_ = d.Set("domain_name", watermarkRule.DomainName)
	}

	if watermarkRule.AppName != nil {
		_ = d.Set("app_name", watermarkRule.AppName)
	}

	if watermarkRule.StreamName != nil {
		_ = d.Set("stream_name", watermarkRule.StreamName)
	}

	if watermarkRule.WatermarkId != nil {
		_ = d.Set("watermark_id", watermarkRule.WatermarkId)
	}

	if watermarkRule.CreateTime != nil {
		_ = d.Set("create_time", watermarkRule.CreateTime)
	}

	if watermarkRule.UpdateTime != nil {
		_ = d.Set("update_time", watermarkRule.UpdateTime)
	}

	return nil
}

func resourceTencentCloudCssWatermarkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule.delete")()
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
	watermarkId := idSplit[3]

	if err := service.DeleteCssWatermarkRuleById(ctx, domainName, appName, streamName, watermarkId); err != nil {
		return err
	}

	return nil
}
