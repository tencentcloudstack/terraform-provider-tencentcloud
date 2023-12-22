package mariadb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbFlow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbFlowRead,
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Flow ID returned by async request API.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Flow status. 0: succeeded, 1: failed, 2: running.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbFlowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_flow.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		status  *mariadb.DescribeFlowResponseParams
		flowId  int
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("flow_id"); v != nil {
		paramMap["FlowId"] = helper.IntInt64(v.(int))
		flowId = v.(int)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbFlowByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		status = result
		return nil
	})

	if err != nil {
		return err
	}

	if status.Status != nil {
		_ = d.Set("status", status.Status)
	}

	d.SetId(strconv.Itoa(flowId))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
