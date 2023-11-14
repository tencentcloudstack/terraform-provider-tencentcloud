/*
Provides a resource to create a live record_rule

Example Usage

```hcl
resource "tencentcloud_live_record_rule" "record_rule" {
  domain_name = ""
  template_id =
  app_name = ""
  stream_name = ""
}
```

Import

live record_rule can be imported using the id, e.g.

```
terraform import tencentcloud_live_record_rule.record_rule record_rule_id
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

func resourceTencentCloudLiveRecordRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveRecordRuleCreate,
		Read:   resourceTencentCloudLiveRecordRuleRead,
		Delete: resourceTencentCloudLiveRecordRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Push domain name.",
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
				Description: "Push path, which is the same as the AppName in push and playback addresses and is live by default.",
			},

			"stream_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream name.Note: If the parameter is a non-empty string, the rule will be only applicable to the particular stream.",
			},
		},
	}
}

func resourceTencentCloudLiveRecordRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveRecordRuleRequest()
		response   = live.NewCreateLiveRecordRuleResponse()
		domainName string
		appName    string
		streamName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
		request.StreamName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveRecordRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live recordRule failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(strings.Join([]string{domainName, appName, streamName}, FILED_SP))

	return resourceTencentCloudLiveRecordRuleRead(d, meta)
}

func resourceTencentCloudLiveRecordRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_rule.read")()
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

	recordRule, err := service.DescribeLiveRecordRuleById(ctx, domainName, appName, streamName)
	if err != nil {
		return err
	}

	if recordRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveRecordRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudLiveRecordRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_rule.delete")()
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

	if err := service.DeleteLiveRecordRuleById(ctx, domainName, appName, streamName); err != nil {
		return err
	}

	return nil
}
