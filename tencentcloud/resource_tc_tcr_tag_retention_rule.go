/*
Provides a resource to create a tcr tag retention rule.

Example Usage

Create a tcr tag retention rule instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id 	 = tencentcloud_tcr_instance.example.id
  name			 = "tf_example_ns_retention"
  is_public		 = true
  is_auto_scan	 = true
  is_prevent_vul = true
  severity		 = "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  retention_rule {
		key = "nDaysSinceLastPush"
		value = 2
  }
  cron_setting = "daily"
  disabled = true
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The Name of the namespace.",
			},

			"retention_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The ID of the retention task.",
			},

			"retention_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Retention Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The supported policies are latestPushedK (retain the latest `k` pushed versions) and nDaysSinceLastPush (retain pushed versions within the last `n` days).",
						},
						"value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "corresponding values for rule settings.",
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

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		request       = tcr.NewCreateTagRetentionRuleRequest()
		registryId    string
		namespaceName string
		tcrService    = TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("namespace_id"); ok {
		request.NamespaceId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		namespace, has, err := tcrService.DescribeTCRNameSpaceById(ctx, registryId, namespaceName)
		if !has || namespace == nil {
			return fmt.Errorf("TCR namespace not found.")
		}
		if err != nil {
			return err
		}
		request.NamespaceId = namespace.NamespaceId
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().CreateTagRetentionRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr TagRetentionRule failed, reason:%+v", logId, err)
		return err
	}

	TagRetentionRule, err := tcrService.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, nil)
	if err != nil {
		return fmt.Errorf("Query retention rule by id failed, reason:[%s]", err.Error())
	}

	if TagRetentionRule != nil {
		retentionId := helper.Int64ToStr(*TagRetentionRule.RetentionId)
		d.SetId(strings.Join([]string{registryId, namespaceName, retentionId}, FILED_SP))
	} else {
		log.Printf("[CRITAL]%s TagRetentionRule is nil! Set unique id as empty.", logId)
		d.SetId("")
	}

	return resourceTencentCloudTcrTagRetentionRuleRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	TagRetentionRule, err := service.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, &retentionId)
	if err != nil {
		return err
	}

	if TagRetentionRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrTagRetentionRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("registry_id", registryId)

	if TagRetentionRule.RetentionId != nil {
		_ = d.Set("retention_id", TagRetentionRule.RetentionId)
	}

	if TagRetentionRule.NamespaceName != nil {
		_ = d.Set("namespace_name", TagRetentionRule.NamespaceName)
	}

	if len(TagRetentionRule.RetentionRuleList) > 0 {
		retentionRuleMap := map[string]interface{}{}
		retentionRule := TagRetentionRule.RetentionRuleList[0]

		if retentionRule.Key != nil {
			retentionRuleMap["key"] = retentionRule.Key
		}

		if retentionRule.Value != nil {
			retentionRuleMap["value"] = retentionRule.Value
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
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	namespace, has, err := tcrService.DescribeTCRNameSpaceById(ctx, registryId, namespaceName)
	if !has || namespace == nil {
		return fmt.Errorf("TCR namespace not found.")
	}
	if err != nil {
		return err
	}

	request.RegistryId = &registryId
	request.NamespaceId = namespace.NamespaceId
	request.RetentionId = helper.StrToInt64Point(retentionId)
	if v, ok := d.GetOkExists("cron_setting"); ok {
		request.CronSetting = helper.String(v.(string))
	}

	immutableArgs := []string{"registry_id", "namespace_name", "cron_setting"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
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

	if d.HasChange("disabled") {
		if v, ok := d.GetOkExists("disabled"); ok {
			request.Disabled = helper.Bool(v.(bool))
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().ModifyTagRetentionRule(request)
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

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	retentionId := idSplit[2]

	if err := service.DeleteTcrTagRetentionRuleById(ctx, registryId, retentionId); err != nil {
		return err
	}

	return nil
}
