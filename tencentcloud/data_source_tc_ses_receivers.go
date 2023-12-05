package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSesReceivers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesReceiversRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Group status (`1`: to be uploaded; `2`: uploading; `3`: uploaded). To query groups in all states, do not pass in this parameter.",
			},

			"key_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group name keyword for fuzzy query.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"receiver_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recipient group ID.",
						},
						"receivers_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recipient group name.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of recipient email addresses.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recipient group descriptionNote: This field may return `null`, indicating that no valid value can be found.",
						},
						"receivers_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Group status (`1`: to be uploaded; `2` uploading; `3` uploaded)Note: This field may return `null`, indicating that no valid value can be found.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, such as 2021-09-28 16:40:35.",
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

func dataSourceTencentCloudSesReceiversRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ses_receivers.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("status"); v != nil {
		paramMap["Status"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("key_word"); ok {
		paramMap["KeyWord"] = helper.String(v.(string))
	}

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*ses.ReceiverData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesReceiversByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, receiverData := range data {
			receiverDataMap := map[string]interface{}{}

			if receiverData.ReceiverId != nil {
				receiverDataMap["receiver_id"] = receiverData.ReceiverId
			}

			if receiverData.ReceiversName != nil {
				receiverDataMap["receivers_name"] = receiverData.ReceiversName
			}

			if receiverData.Count != nil {
				receiverDataMap["count"] = receiverData.Count
			}

			if receiverData.Desc != nil {
				receiverDataMap["desc"] = receiverData.Desc
			}

			if receiverData.ReceiversStatus != nil {
				receiverDataMap["receivers_status"] = receiverData.ReceiversStatus
			}

			if receiverData.CreateTime != nil {
				receiverDataMap["create_time"] = receiverData.CreateTime
			}

			ids = append(ids, strconv.Itoa(int(*receiverData.ReceiverId)))
			tmpList = append(tmpList, receiverDataMap)
		}

		_ = d.Set("data", tmpList)
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
