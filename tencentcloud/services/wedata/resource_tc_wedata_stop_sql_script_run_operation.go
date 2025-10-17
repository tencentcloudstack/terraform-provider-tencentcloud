package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataStopSqlScriptRunOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataStopSqlScriptRunOperationCreate,
		Read:   resourceTencentCloudWedataStopSqlScriptRunOperationRead,
		Delete: resourceTencentCloudWedataStopSqlScriptRunOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the query id.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},
		},
	}
}

func resourceTencentCloudWedataStopSqlScriptRunOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_stop_sql_script_run_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewStopSQLScriptRunRequest()
		projectId string
		jobId     string
	)

	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().StopSQLScriptRunWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata stop sql script run operation failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Create wedata stop sql script run operation failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata stop sql script run operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{projectId, jobId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataStopSqlScriptRunOperationRead(d, meta)
}

func resourceTencentCloudWedataStopSqlScriptRunOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_stop_sql_script_run_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWedataStopSqlScriptRunOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_stop_sql_script_run_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
