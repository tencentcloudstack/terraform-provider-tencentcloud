package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClickhouseInstanceShards() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseInstanceShardsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster instance ID.",
			},

			"instance_shards_list": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance shard information.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClickhouseInstanceShardsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clickhouse_instance_shards.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceShards *cdwch.DescribeInstanceShardsResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClickhouseInstanceShardsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceShards = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceShards.InstanceShardsList != nil {
		_ = d.Set("instance_shards_list", instanceShards.InstanceShardsList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
