/*
Provides a resource to create a tcr tag_retention_execution

Example Usage

```hcl
resource "tencentcloud_tcr_tag_retention_execution" "tag_retention_execution" {
  registry_id = "tcr-xx"
  retention_id = 1
  dry_run = false
}
```

Import

tcr tag_retention_execution can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_tag_retention_execution.tag_retention_execution tag_retention_execution_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTcrTagRetentionExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrTagRetentionExecutionCreate,
		Read:   resourceTencentCloudTcrTagRetentionExecutionRead,
		Delete: resourceTencentCloudTcrTagRetentionExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"retention_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Retention id.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to simulate execution, the default value is false, that is, non-simulation execution.",
			},
		},
	}
}

func resourceTencentCloudTcrTagRetentionExecutionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tcr.NewCreateTagRetentionExecutionRequest()
		response    = tcr.NewCreateTagRetentionExecutionResponse()
		registryId  string
		retentionId string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("retention_id"); v != nil {
		retentionId = v.(int64)
		request.RetentionId = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateTagRetentionExecution(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr TagRetentionExecution failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, retentionId}, FILED_SP))

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Succeed"}, 600*readRetryTimeout, time.Second, service.TcrTagRetentionExecutionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTcrTagRetentionExecutionRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionExecutionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrTagRetentionExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_tag_retention_execution.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
