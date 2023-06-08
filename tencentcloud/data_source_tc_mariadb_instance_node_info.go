/*
Use this data source to query detailed information of mariadb instance_node_info

Example Usage

```hcl
data "tencentcloud_mariadb_instance_node_info" "instance_node_info" {
  instance_id = "tdsql-9vqvls95"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbInstanceNodeInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbInstanceNodeInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, such as tdsql-6ltok4u9.",
			},
			"nodes_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node role. Valid values: `master`, `slave`.",
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

func dataSourceTencentCloudMariadbInstanceNodeInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_instance_node_info.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		nodesInfo  []*mariadb.NodeInfo
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbInstanceNodeInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		nodesInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(nodesInfo))

	if nodesInfo != nil {
		for _, nodeInfo := range nodesInfo {
			nodeInfoMap := map[string]interface{}{}

			if nodeInfo.NodeId != nil {
				nodeInfoMap["node_id"] = nodeInfo.NodeId
			}

			if nodeInfo.Role != nil {
				nodeInfoMap["role"] = nodeInfo.Role
			}

			tmpList = append(tmpList, nodeInfoMap)
		}

		_ = d.Set("nodes_info", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
