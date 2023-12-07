package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsEnableWorkflowConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsEnableWorkflowConfigCreate,
		Read:   resourceTencentCloudMpsEnableWorkflowConfigRead,
		Update: resourceTencentCloudMpsEnableWorkflowConfigUpdate,
		Delete: resourceTencentCloudMpsEnableWorkflowConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Workflow ID.",
			},

			"enabled": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "true: enable; false: disable.",
			},
		},
	}
}

func resourceTencentCloudMpsEnableWorkflowConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_workflow_config.create")()
	defer inconsistentCheck(d, meta)()

	var workflowId int
	if v, ok := d.GetOkExists("workflow_id"); ok {
		workflowId = v.(int)
	}
	d.SetId(helper.IntToStr(workflowId))

	return resourceTencentCloudMpsEnableWorkflowConfigUpdate(d, meta)
}

func resourceTencentCloudMpsEnableWorkflowConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_workflow_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	workflowId := d.Id()

	enableWorkflowConfig, err := service.DescribeMpsWorkflowById(ctx, workflowId)
	if err != nil {
		return err
	}

	if enableWorkflowConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsEnableWorkflowConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if enableWorkflowConfig.WorkflowId != nil {
		_ = d.Set("workflow_id", enableWorkflowConfig.WorkflowId)
	}

	status := enableWorkflowConfig.Status
	if status != nil {
		if *status == "Enabled" {
			_ = d.Set("enabled", true)
		} else {
			_ = d.Set("enabled", false)
		}
	}

	return nil
}

func resourceTencentCloudMpsEnableWorkflowConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_workflow_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enableRequest  = mps.NewEnableWorkflowRequest()
		disableRequest = mps.NewDisableWorkflowRequest()
		workflowId     *int64
		enabled        bool
	)

	workflowId = helper.StrToInt64Point(d.Id())

	if v, ok := d.GetOkExists("enabled"); ok && v != nil {
		enabled = v.(bool)

		if enabled {
			enableRequest.WorkflowId = workflowId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().EnableWorkflow(enableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate mps enableWorkflowConfig failed, reason:%+v", logId, err)
				return err
			}
		} else {
			disableRequest.WorkflowId = workflowId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().DisableWorkflow(disableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate mps disableWorkflowConfig failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudMpsEnableWorkflowConfigRead(d, meta)
}

func resourceTencentCloudMpsEnableWorkflowConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_workflow_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
