package wedata

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataRunSqlScriptOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataRunSqlScriptOperationCreate,
		Read:   resourceTencentCloudWedataRunSqlScriptOperationRead,
		Delete: resourceTencentCloudWedataRunSqlScriptOperationDelete,
		Schema: map[string]*schema.Schema{
			"script_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Script id.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"script_content": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Script content. executed by default if not transmitted.",
			},

			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced running parameter.",
			},

			// computed
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job ID of the SQL script operation.",
			},
		},
	}
}

func resourceTencentCloudWedataRunSqlScriptOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_run_sql_script_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wedatav20250806.NewRunSQLScriptRequest()
		response = wedatav20250806.NewRunSQLScriptResponse()
		jobId    string
	)

	if v, ok := d.GetOk("script_id"); ok {
		request.ScriptId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("script_content"); ok {
		request.ScriptContent = helper.String(tccommon.StringToBase64(v.(string)))
	}

	if v, ok := d.GetOk("params"); ok {
		request.Params = helper.String(tccommon.StringToBase64(v.(string)))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RunSQLScriptWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata run sql script operation failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata run sql script operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.JobId == nil {
		return fmt.Errorf("JobId is nil")
	}

	jobId = *response.Response.Data.JobId
	_ = d.Set("job_id", jobId)

	d.SetId(jobId)
	return resourceTencentCloudWedataRunSqlScriptOperationRead(d, meta)
}

func resourceTencentCloudWedataRunSqlScriptOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_run_sql_script_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWedataRunSqlScriptOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_run_sql_script_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
