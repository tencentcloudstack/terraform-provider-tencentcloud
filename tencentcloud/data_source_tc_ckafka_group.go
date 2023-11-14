/*
Use this data source to query detailed information of ckafka group

Example Usage

```hcl
data "tencentcloud_ckafka_group" "group" {
  instance_id = "InstanceId"
  search_word = "SearchWord"
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

func dataSourceTencentCloudCkafkaGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaGroupRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search for the keyword.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total Count of the Result.",
						},
						"group_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "GroupList.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "GroupId.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The protocol used by this group.",
									},
								},
							},
						},
						"group_count_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "GroupCountQuota.",
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

func dataSourceTencentCloudCkafkaGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.GroupResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaGroupByFilter(ctx, paramMap)
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
		groupResponseMap := map[string]interface{}{}

		if result.TotalCount != nil {
			groupResponseMap["total_count"] = result.TotalCount
		}

		if result.GroupList != nil {
			groupListList := []interface{}{}
			for _, groupList := range result.GroupList {
				groupListMap := map[string]interface{}{}

				if groupList.Group != nil {
					groupListMap["group"] = groupList.Group
				}

				if groupList.Protocol != nil {
					groupListMap["protocol"] = groupList.Protocol
				}

				groupListList = append(groupListList, groupListMap)
			}

			groupResponseMap["group_list"] = []interface{}{groupListList}
		}

		if result.GroupCountQuota != nil {
			groupResponseMap["group_count_quota"] = result.GroupCountQuota
		}

		ids = append(ids, *result.InstanceId)
		_ = d.Set("result", groupResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupResponseMap); e != nil {
			return e
		}
	}
	return nil
}
