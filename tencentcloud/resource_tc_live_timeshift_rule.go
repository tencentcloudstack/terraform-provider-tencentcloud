/*
Provides a resource to create a live timeshift_rule

Example Usage

```hcl
resource "tencentcloud_live_timeshift_rule" "timeshift_rule" {
  domain_name = ""
  app_name = ""
  stream_name = ""
  template_id =
}
```

Import

live timeshift_rule can be imported using the id, e.g.

```
terraform import tencentcloud_live_timeshift_rule.timeshift_rule timeshift_rule_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudLiveTimeshiftRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveTimeshiftRuleCreate,
		Read:   resourceTencentCloudLiveTimeshiftRuleRead,
		Delete: resourceTencentCloudLiveTimeshiftRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The push domain.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The push path, which should be the same as `AppName` in the push and playback URLs. The default value is `live`.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The stream name.Note: If you pass in a non-empty string, the rule will only be applied to the specified stream.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The template ID.",
			},
		},
	}
}

func resourceTencentCloudLiveTimeshiftRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveTimeShiftRuleRequest()
		response   = live.NewCreateLiveTimeShiftRuleResponse()
		domainName string
		appName    string
		streamName string
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
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveTimeShiftRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live timeshiftRule failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(strings.Join([]string{domainName, appName, streamName}, FILED_SP))

	return resourceTencentCloudLiveTimeshiftRuleRead(d, meta)
}

func resourceTencentCloudLiveTimeshiftRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	timeshiftRule, err := service.DescribeLiveTimeshiftRuleById(ctx, domainName, appName, streamName)
	if err != nil {
		return err
	}

	if timeshiftRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveTimeshiftRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if timeshiftRule.DomainName != nil {
		_ = d.Set("domain_name", timeshiftRule.DomainName)
	}

	if timeshiftRule.AppName != nil {
		_ = d.Set("app_name", timeshiftRule.AppName)
	}

	if timeshiftRule.StreamName != nil {
		_ = d.Set("stream_name", timeshiftRule.StreamName)
	}

	if timeshiftRule.TemplateId != nil {
		_ = d.Set("template_id", timeshiftRule.TemplateId)
	}

	return nil
}

func resourceTencentCloudLiveTimeshiftRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	if err := service.DeleteLiveTimeshiftRuleById(ctx, domainName, appName, streamName); err != nil {
		return err
	}

	return nil
}
