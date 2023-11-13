/*
Use this data source to query detailed information of ckafka task_status

Example Usage

```hcl
data "tencentcloud_ckafka_task_status" "task_status" {
  flow_id = flowId
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaTaskStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTaskStatusRead,
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "FlowId.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status.",
						},
						"output": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OutPut Info.",
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

func dataSourceTencentCloudCkafkaTaskStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_task_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("flow_id"); v != nil {
		paramMap["FlowId"] = helper.IntInt64(v.(int))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.TaskStatusResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaTaskStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		taskStatusResponseMap := map[string]interface{}{}

		if result.Status != nil {
			taskStatusResponseMap["status"] = result.Status
		}

		if result.Output != nil {
			taskStatusResponseMap["output"] = result.Output
		}

		ids = append(ids, *result.FlowId)
		_ = d.Set("result", taskStatusResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), taskStatusResponseMap); e != nil {
			return e
		}
	}
	return nil
}
