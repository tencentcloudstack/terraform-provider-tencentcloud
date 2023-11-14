/*
Use this data source to query detailed information of mariadb flow

Example Usage

```hcl
data "tencentcloud_mariadb_flow" "flow" {
  flow_id =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbFlow() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_mariadb_flow.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("flow_id"); v != nil {
		paramMap["FlowId"] = helper.IntInt64(v.(int))
	}

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbFlowByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		status = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(status))
	if status != nil {
		_ = d.Set("status", status)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
