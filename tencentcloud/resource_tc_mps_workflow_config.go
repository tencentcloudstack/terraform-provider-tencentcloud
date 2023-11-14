/*
Provides a resource to create a mps workflow_config

Example Usage

```hcl
resource "tencentcloud_mps_workflow_config" "workflow_config" {
  workflow_id =
}
```

Import

mps workflow_config can be imported using the id, e.g.

```
terraform import tencentcloud_mps_workflow_config.workflow_config workflow_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsWorkflowConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWorkflowConfigCreate,
		Read:   resourceTencentCloudMpsWorkflowConfigRead,
		Update: resourceTencentCloudMpsWorkflowConfigUpdate,
		Delete: resourceTencentCloudMpsWorkflowConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Workflow ID.",
			},
		},
	}
}

func resourceTencentCloudMpsWorkflowConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow_config.create")()
	defer inconsistentCheck(d, meta)()

	var workflowId int64
	if v, ok := d.GetOkExists("workflow_id"); ok {
		workflowId = v.(int64)
	}

	d.SetId(helper.Int64ToStr(workflowId))

	return resourceTencentCloudMpsWorkflowConfigUpdate(d, meta)
}

func resourceTencentCloudMpsWorkflowConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	workflowConfigId := d.Id()

	workflowConfig, err := service.DescribeMpsWorkflowConfigById(ctx, workflowId)
	if err != nil {
		return err
	}

	if workflowConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWorkflowConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if workflowConfig.WorkflowId != nil {
		_ = d.Set("workflow_id", workflowConfig.WorkflowId)
	}

	return nil
}

func resourceTencentCloudMpsWorkflowConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewEnableWorkflowRequest()

	workflowConfigId := d.Id()

	request.WorkflowId = &workflowId

	immutableArgs := []string{"workflow_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().EnableWorkflow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps workflowConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsWorkflowConfigRead(d, meta)
}

func resourceTencentCloudMpsWorkflowConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
