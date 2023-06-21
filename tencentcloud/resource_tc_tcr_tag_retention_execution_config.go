/*
Provides a resource to create a tcr tag_retention_execution_config

Example Usage

```hcl
resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id    = tencentcloud_tcr_instance.mytcr_retention.id
  name           = "tf_test_ns_retention"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id    = tencentcloud_tcr_instance.mytcr_retention.id
  namespace_name = tencentcloud_tcr_namespace.my_ns.name
  retention_rule {
    key   = "nDaysSinceLastPush"
    value = 2
  }
  cron_setting = "manual"
  disabled     = true
}

resource "tencentcloud_tcr_tag_retention_execution_config" "tag_retention_execution_config" {
  registry_id  = tencentcloud_tcr_tag_retention_rule.my_rule.registry_id
  retention_id = tencentcloud_tcr_tag_retention_rule.my_rule.retention_id
  dry_run      = false
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrTagRetentionExecutionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrTagRetentionExecutionConfigCreate,
		Read:   resourceTencentCloudTcrTagRetentionExecutionConfigRead,
		Update: resourceTencentCloudTcrTagRetentionExecutionConfigUpdate,
		Delete: resourceTencentCloudTcrTagRetentionExecutionConfigDelete,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"retention_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "retention id.",
			},

			"execution_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "execution id.",
			},

			"dry_run": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to simulate execution, the default value is false, that is, non-simulation execution.",
			},
		},
	}
}

func resourceTencentCloudTcrTagRetentionExecutionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		registryId  string
		retentionId string
	)

	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
	}

	if v, ok := d.GetOk("retention_id"); ok {
		retentionId = helper.IntToStr(v.(int))
	}

	d.SetId(strings.Join([]string{registryId, retentionId}, FILED_SP))

	return resourceTencentCloudTcrTagRetentionExecutionConfigUpdate(d, meta)
}

func resourceTencentCloudTcrTagRetentionExecutionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	retentionId := idSplit[1]

	TagRetentionExecutionConfig, err := service.DescribeTcrTagRetentionExecutionConfigById(ctx, registryId, retentionId)
	if err != nil {
		return err
	}

	if TagRetentionExecutionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrTagRetentionExecutionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("registry_id", registryId)

	if TagRetentionExecutionConfig.RetentionId != nil {
		_ = d.Set("retention_id", TagRetentionExecutionConfig.RetentionId)
	}

	if TagRetentionExecutionConfig.ExecutionId != nil {
		_ = d.Set("execution_id", TagRetentionExecutionConfig.ExecutionId)
	}

	return nil
}

func resourceTencentCloudTcrTagRetentionExecutionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewCreateTagRetentionExecutionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	retentionId := idSplit[1]

	request.RegistryId = &registryId
	request.RetentionId = helper.StrToInt64Point(retentionId)

	if d.HasChange("dry_run") {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request.DryRun = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().CreateTagRetentionExecution(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr TagRetentionExecutionConfig failed, reason:%+v", logId, err)
		return err
	}

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Succeed"}, 3*readRetryTimeout, time.Second, service.TcrTagRetentionExecutionConfigStateRefreshFunc(registryId, retentionId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTcrTagRetentionExecutionConfigRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionExecutionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
