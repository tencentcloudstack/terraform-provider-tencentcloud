/*
Provides a resource to create a live callback_rule

Example Usage

```hcl
resource "tencentcloud_live_callback_rule" "callback_rule" {
  domain_name = "5000.livepush.myqcloud.com"
  app_name = "live"
  template_id = 1000
}
```

Import

live callback_rule can be imported using the id, e.g.

```
terraform import tencentcloud_live_callback_rule.callback_rule callback_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudLiveCallbackRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveCallbackRuleCreate,
		Read:   resourceTencentCloudLiveCallbackRuleRead,
		Update: resourceTencentCloudLiveCallbackRuleUpdate,
		Delete: resourceTencentCloudLiveCallbackRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Streaming domain name.",
			},

			"app_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The streaming path is consistent with the AppName in the streaming and playback addresses. The default is live.",
			},

			"template_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Template ID.",
			},
		},
	}
}

func resourceTencentCloudLiveCallbackRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveCallbackRuleRequest()
		response   = live.NewCreateLiveCallbackRuleResponse()
		domainName string
		appName    string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveCallbackRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live callbackRule failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(strings.Join([]string{domainName, appName}, FILED_SP))

	return resourceTencentCloudLiveCallbackRuleRead(d, meta)
}

func resourceTencentCloudLiveCallbackRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]

	callbackRule, err := service.DescribeLiveCallbackRuleById(ctx, domainName, appName)
	if err != nil {
		return err
	}

	if callbackRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveCallbackRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if callbackRule.DomainName != nil {
		_ = d.Set("domain_name", callbackRule.DomainName)
	}

	if callbackRule.AppName != nil {
		_ = d.Set("app_name", callbackRule.AppName)
	}

	if callbackRule.TemplateId != nil {
		_ = d.Set("template_id", callbackRule.TemplateId)
	}

	return nil
}

func resourceTencentCloudLiveCallbackRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"domain_name", "app_name", "template_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudLiveCallbackRuleRead(d, meta)
}

func resourceTencentCloudLiveCallbackRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]

	if err := service.DeleteLiveCallbackRuleById(ctx, domainName, appName); err != nil {
		return err
	}

	return nil
}
