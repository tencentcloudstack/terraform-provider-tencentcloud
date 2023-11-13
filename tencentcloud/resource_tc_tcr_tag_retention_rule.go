/*
Provides a resource to create a tcr tag_retention_rule

Example Usage

```hcl
resource "tencentcloud_tcr_tag_retention_rule" "tag_retention_rule" {
  registry_id = "tcr-12345"
  namespace_id = 1
  retention_rule {
		key = "latestPushedK"
		value = 1

  }
  cron_setting = "manual"
  disabled = false
}
```

Import

tcr tag_retention_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_tag_retention_rule.tag_retention_rule tag_retention_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTcrTagRetentionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrTagRetentionRuleCreate,
		Read:   resourceTencentCloudTcrTagRetentionRuleRead,
		Update: resourceTencentCloudTcrTagRetentionRuleUpdate,
		Delete: resourceTencentCloudTcrTagRetentionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The main instance ID.",
			},

			"namespace_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The id of the namespace.",
			},

			"retention_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Retention policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The supported policies are latestPushedK (retain the latest “k” pushed versions) and nDaysSinceLastPush (retain pushed versions within the last “n” days).",
						},
						"value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Corresponding values for rule settings.",
						},
					},
				},
			},

			"cron_setting": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Execution cycle, currently only available selections are: manual; daily; weekly; monthly.",
			},

			"disabled": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to disable the rule, with the default value of false.",
			},
		},
	}
}

func resourceTencentCloudTcrTagRetentionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tcr.NewCreateTagRetentionRuleRequest()
		response      = tcr.NewCreateTagRetentionRuleResponse()
		registryId    string
		namespaceName string
		retentionId   int
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("namespace_id"); ok {
		request.NamespaceId = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "retention_rule"); ok {
		retentionRule := tcr.RetentionRule{}
		if v, ok := dMap["key"]; ok {
			retentionRule.Key = helper.String(v.(string))
		}
		if v, ok := dMap["value"]; ok {
			retentionRule.Value = helper.IntInt64(v.(int))
		}
		request.RetentionRule = &retentionRule
	}

	if v, ok := d.GetOk("cron_setting"); ok {
		request.CronSetting = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disabled"); ok {
		request.Disabled = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateTagRetentionRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr TagRetentionRule failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, namespaceName, helper.Int64ToStr(retentionId)}, FILED_SP))

	return resourceTencentCloudTcrTagRetentionRuleRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	TagRetentionRule, err := service.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, retentionId)
	if err != nil {
		return err
	}

	if TagRetentionRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrTagRetentionRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if TagRetentionRule.RegistryId != nil {
		_ = d.Set("registry_id", TagRetentionRule.RegistryId)
	}

	if TagRetentionRule.NamespaceId != nil {
		_ = d.Set("namespace_id", TagRetentionRule.NamespaceId)
	}

	if TagRetentionRule.RetentionRule != nil {
		retentionRuleMap := map[string]interface{}{}

		if TagRetentionRule.RetentionRule.Key != nil {
			retentionRuleMap["key"] = TagRetentionRule.RetentionRule.Key
		}

		if TagRetentionRule.RetentionRule.Value != nil {
			retentionRuleMap["value"] = TagRetentionRule.RetentionRule.Value
		}

		_ = d.Set("retention_rule", []interface{}{retentionRuleMap})
	}

	if TagRetentionRule.CronSetting != nil {
		_ = d.Set("cron_setting", TagRetentionRule.CronSetting)
	}

	if TagRetentionRule.Disabled != nil {
		_ = d.Set("disabled", TagRetentionRule.Disabled)
	}

	return nil
}

func resourceTencentCloudTcrTagRetentionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewModifyTagRetentionRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	request.RegistryId = &registryId
	request.NamespaceName = &namespaceName
	request.RetentionId = &retentionId

	immutableArgs := []string{"registry_id", "namespace_id", "retention_rule", "cron_setting", "disabled"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("registry_id") {
		if v, ok := d.GetOk("registry_id"); ok {
			request.RegistryId = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_id") {
		if v, ok := d.GetOkExists("namespace_id"); ok {
			request.NamespaceId = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("retention_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "retention_rule"); ok {
			retentionRule := tcr.RetentionRule{}
			if v, ok := dMap["key"]; ok {
				retentionRule.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				retentionRule.Value = helper.IntInt64(v.(int))
			}
			request.RetentionRule = &retentionRule
		}
	}

	if d.HasChange("cron_setting") {
		if v, ok := d.GetOk("cron_setting"); ok {
			request.CronSetting = helper.String(v.(string))
		}
	}

	if d.HasChange("disabled") {
		if v, ok := d.GetOkExists("disabled"); ok {
			request.Disabled = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().ModifyTagRetentionRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr TagRetentionRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcrTagRetentionRuleRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	if err := service.DeleteTcrTagRetentionRuleById(ctx, registryId, namespaceName, retentionId); err != nil {
		return err
	}

	return nil
}
