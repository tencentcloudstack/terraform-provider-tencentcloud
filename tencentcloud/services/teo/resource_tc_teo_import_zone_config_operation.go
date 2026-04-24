package teo

import (
	"fmt"
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoImportZoneConfigOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoImportZoneConfigOperationCreate,
		Read:   resourceTencentCloudTeoImportZoneConfigOperationRead,
		Delete: resourceTencentCloudTeoImportZoneConfigOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The configuration content to import. It must be in JSON format and encoded in UTF-8. You can obtain the configuration content via the tencentcloud_teo_export_zone_config data source.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task ID of the import configuration operation.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The import task status. Valid values: success, failure, doing.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status message of the import task. When the configuration import fails, you can view the failure reason through this field.",
			},
			"import_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the import task.",
			},
			"finish_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The end time of the import task.",
			},
		},
	}
}

func resourceTencentCloudTeoImportZoneConfigOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	content := d.Get("content").(string)

	request := teov20220901.NewImportZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)
	request.Content = helper.String(content)

	var taskId string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ImportZoneConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result.Response.TaskId != nil {
			taskId = *result.Response.TaskId
		}
		return nil
	})
	if err != nil {
		return err
	}

	if taskId == "" {
		return fmt.Errorf("ImportZoneConfig API returned empty TaskId")
	}

	_ = d.Set("task_id", taskId)

	// Poll DescribeZoneConfigImportResult until the task is complete
	describeRequest := teov20220901.NewDescribeZoneConfigImportResultRequest()
	describeRequest.ZoneId = helper.String(zoneId)
	describeRequest.TaskId = helper.String(taskId)

	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeZoneConfigImportResult(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeZoneConfigImportResult API returned empty response"))
		}

		status := ""
		if result.Response.Status != nil {
			status = *result.Response.Status
		}

		if status == "doing" {
			return resource.RetryableError(fmt.Errorf("import zone config task is still processing, status: doing"))
		}

		if status == "failure" {
			msg := ""
			if result.Response.Message != nil {
				msg = *result.Response.Message
			}
			return resource.NonRetryableError(fmt.Errorf("import zone config task failed, message: %s", msg))
		}

		// status == "success"
		if result.Response.Status != nil {
			_ = d.Set("status", *result.Response.Status)
		}
		if result.Response.Message != nil {
			_ = d.Set("message", *result.Response.Message)
		}
		if result.Response.ImportTime != nil {
			_ = d.Set("import_time", *result.Response.ImportTime)
		}
		if result.Response.FinishTime != nil {
			_ = d.Set("finish_time", *result.Response.FinishTime)
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + tccommon.FILED_SP + taskId)

	return resourceTencentCloudTeoImportZoneConfigOperationRead(d, meta)
}

func resourceTencentCloudTeoImportZoneConfigOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoImportZoneConfigOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
