/*
Use this data source to query detailed information of tdmq rabbitmq_virtual_host_list

Example Usage

```hcl
data "tencentcloud_tdmq_rabbitmq_virtual_host_list" "rabbitmq_virtual_host_list" {
  instance_id = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqRabbitmqVirtualHostList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRabbitmqVirtualHostListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Not applicable, default parameters.",
			},

			"virtual_host_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster listNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual host nameNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the virtual hostNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTdmqRabbitmqVirtualHostListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_rabbitmq_virtual_host_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var virtualHostList []*tdmq.RabbitMQPrivateVirtualHost

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVirtualHostListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		virtualHostList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(virtualHostList))
	tmpList := make([]map[string]interface{}, 0, len(virtualHostList))

	if virtualHostList != nil {
		for _, rabbitMQPrivateVirtualHost := range virtualHostList {
			rabbitMQPrivateVirtualHostMap := map[string]interface{}{}

			if rabbitMQPrivateVirtualHost.VirtualHostName != nil {
				rabbitMQPrivateVirtualHostMap["virtual_host_name"] = rabbitMQPrivateVirtualHost.VirtualHostName
			}

			if rabbitMQPrivateVirtualHost.Description != nil {
				rabbitMQPrivateVirtualHostMap["description"] = rabbitMQPrivateVirtualHost.Description
			}

			ids = append(ids, *rabbitMQPrivateVirtualHost.InstanceId)
			tmpList = append(tmpList, rabbitMQPrivateVirtualHostMap)
		}

		_ = d.Set("virtual_host_list", tmpList)
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
