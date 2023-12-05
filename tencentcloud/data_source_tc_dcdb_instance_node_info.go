package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbInstanceNodeInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbInstanceNodeInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, such as tdsqlshard-6ltok4u9.",
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
						"shard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance shard ID.",
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

func dataSourceTencentCloudDcdbInstanceNodeInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_instance_node_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var nodesInfo []*dcdb.BriefNodeInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbInstanceNodeInfoByFilter(ctx, paramMap)
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
		for _, briefNodeInfo := range nodesInfo {
			briefNodeInfoMap := map[string]interface{}{}

			if briefNodeInfo.NodeId != nil {
				briefNodeInfoMap["node_id"] = briefNodeInfo.NodeId
			}

			if briefNodeInfo.Role != nil {
				briefNodeInfoMap["role"] = briefNodeInfo.Role
			}

			if briefNodeInfo.ShardId != nil {
				briefNodeInfoMap["shard_id"] = briefNodeInfo.ShardId
			}

			tmpList = append(tmpList, briefNodeInfoMap)
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
