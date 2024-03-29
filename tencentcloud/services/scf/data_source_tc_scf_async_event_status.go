package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfAsyncEventStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfAsyncEventStatusRead,
		Schema: map[string]*schema.Schema{
			"invoke_request_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the async execution request.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Async event status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Async event status. Values: `RUNNING` (running); `FINISHED` (invoked successfully); `ABORTED` (invocation ended); `FAILED` (invocation failed).",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request status code.",
						},
						"invoke_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Async execution request ID.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfAsyncEventStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_async_event_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		result          *scf.AsyncEventStatus
		invokeRequestId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invoke_request_id"); ok {
		invokeRequestId = v.(string)
		paramMap["InvokeRequestId"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		res, e := service.DescribeScfAsyncEventStatus(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = res
		return nil
	})
	if err != nil {
		return err
	}

	asyncEventStatusMap := map[string]interface{}{}
	if result != nil {

		if result.Status != nil {
			asyncEventStatusMap["status"] = result.Status
		}

		if result.StatusCode != nil {
			asyncEventStatusMap["status_code"] = result.StatusCode
		}

		if result.InvokeRequestId != nil {
			asyncEventStatusMap["invoke_request_id"] = result.InvokeRequestId
		}

		_ = d.Set("result", asyncEventStatusMap)
	}

	d.SetId(invokeRequestId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), asyncEventStatusMap); e != nil {
			return e
		}
	}
	return nil
}
