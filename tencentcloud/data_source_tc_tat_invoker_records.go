package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTatInvokerRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatInvokerRecordsRead,
		Schema: map[string]*schema.Schema{
			"invoker_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of invoker IDs. Up to 100 IDs are allowed.",
			},

			"invoker_record_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Execution history of an invoker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invoker_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invoker ID.",
						},
						"invoke_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution time.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution reason.",
						},
						"invocation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command execution ID.",
						},
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger result.",
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

func dataSourceTencentCloudTatInvokerRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tat_invoker_records.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invoker_ids"); ok {
		invokerIdsSet := v.(*schema.Set).List()
		paramMap["InvokerIds"] = helper.InterfacesStringsPoint(invokerIdsSet)
	}

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var invokerRecordSet []*tat.InvokerRecord

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTatInvokerRecordsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		invokerRecordSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(invokerRecordSet))
	tmpList := make([]map[string]interface{}, 0, len(invokerRecordSet))

	if invokerRecordSet != nil {
		for _, invokerRecord := range invokerRecordSet {
			invokerRecordMap := map[string]interface{}{}

			if invokerRecord.InvokerId != nil {
				invokerRecordMap["invoker_id"] = invokerRecord.InvokerId
			}

			if invokerRecord.InvokeTime != nil {
				invokerRecordMap["invoke_time"] = invokerRecord.InvokeTime
			}

			if invokerRecord.Reason != nil {
				invokerRecordMap["reason"] = invokerRecord.Reason
			}

			if invokerRecord.InvocationId != nil {
				invokerRecordMap["invocation_id"] = invokerRecord.InvocationId
			}

			if invokerRecord.Result != nil {
				invokerRecordMap["result"] = invokerRecord.Result
			}

			ids = append(ids, *invokerRecord.InvokerId)
			tmpList = append(tmpList, invokerRecordMap)
		}

		_ = d.Set("invoker_record_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
