package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
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

	flowId := d.Get("flow_id").(int)

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *ckafka.TaskStatusResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, e := service.DescribeCkafkaTaskStatusByFilter(ctx, flowId)
		if e != nil {
			return retryError(e)
		}
		result = taskStatus
		return nil
	})
	if err != nil {
		return err
	}
	taskStatusResponseMapList := make([]interface{}, 0)
	if result != nil {
		taskStatusResponseMap := map[string]interface{}{}
		if result.Status != nil {
			taskStatusResponseMap["status"] = result.Status
		}

		if result.Output != nil {
			taskStatusResponseMap["output"] = result.Output
		}
		taskStatusResponseMapList = append(taskStatusResponseMapList, taskStatusResponseMap)
		_ = d.Set("result", taskStatusResponseMapList)
	}

	d.SetId(strconv.Itoa(flowId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), taskStatusResponseMapList); e != nil {
			return e
		}
	}
	return nil
}
