/*
Use this data source to query detailed information of dcdb shards

Example Usage

```hcl
data "tencentcloud_dcdb_shards" "shards" {
  instance_id = "your_instance_id"
  shard_instance_ids = ["shard1_id"]
  }
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbShards() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbShardsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"shard_instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "shard instance ids.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "shard list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"shard_serial_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "shard serial id.",
						},
						"shard_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "shard instance id.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet id.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "zone.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "memory, the unit is GB.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "memory, the unit is GB.",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expired time.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "node count.",
						},
						"storage_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "storage usage.",
						},
						"memory_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "memory usage.",
						},
						"proxy_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy version.",
						},
						"paymode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "pay mode.",
						},
						"shard_master_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "shard master zone.",
						},
						"shard_slave_zones": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "shard slave zones.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "cpu cores.",
						},
						"range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the range of shard key.",
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

func dataSourceTencentCloudDcdbShardsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_shards.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("shard_instance_ids"); ok {
		shard_instance_idsSet := v.(*schema.Set).List()
		ids := make([]*string, 0, len(shard_instance_idsSet))
		for _, vv := range shard_instance_idsSet {
			ids = append(ids, helper.String(vv.(string)))
		}
		paramMap["shard_instance_ids"] = ids
	}

	dcdbService := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var shards []*dcdb.DCDBShardInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbShardsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		shards = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb shards failed, reason:%+v", logId, err)
		return err
	}

	// shardList := []interface{}{}
	shardList := make([]map[string]interface{}, 0, len(shards))
	ids := make([]string, 0, len(shards))
	if shards != nil {
		for _, shard := range shards {
			shardMap := map[string]interface{}{}
			if shard.InstanceId != nil {
				shardMap["instance_id"] = shard.InstanceId
			}
			if shard.ShardSerialId != nil {
				shardMap["shard_serial_id"] = shard.ShardSerialId
			}
			if shard.ShardInstanceId != nil {
				shardMap["shard_instance_id"] = shard.ShardInstanceId
			}
			if shard.Status != nil {
				shardMap["status"] = shard.Status
			}
			if shard.StatusDesc != nil {
				shardMap["status_desc"] = shard.StatusDesc
			}
			if shard.CreateTime != nil {
				shardMap["create_time"] = shard.CreateTime
			}
			if shard.VpcId != nil {
				shardMap["vpc_id"] = shard.VpcId
			}
			if shard.SubnetId != nil {
				shardMap["subnet_id"] = shard.SubnetId
			}
			if shard.ProjectId != nil {
				shardMap["project_id"] = shard.ProjectId
			}
			if shard.Region != nil {
				shardMap["region"] = shard.Region
			}
			if shard.Zone != nil {
				shardMap["zone"] = shard.Zone
			}
			if shard.Memory != nil {
				shardMap["memory"] = shard.Memory
			}
			if shard.Storage != nil {
				shardMap["storage"] = shard.Storage
			}
			if shard.PeriodEndTime != nil {
				shardMap["period_end_time"] = shard.PeriodEndTime
			}
			if shard.NodeCount != nil {
				shardMap["node_count"] = shard.NodeCount
			}
			if shard.StorageUsage != nil {
				shardMap["storage_usage"] = shard.StorageUsage
			}
			if shard.MemoryUsage != nil {
				shardMap["memory_usage"] = shard.MemoryUsage
			}
			if shard.ProxyVersion != nil {
				shardMap["proxy_version"] = shard.ProxyVersion
			}
			if shard.Paymode != nil {
				shardMap["paymode"] = shard.Paymode
			}
			if shard.ShardMasterZone != nil {
				shardMap["shard_master_zone"] = shard.ShardMasterZone
			}
			if shard.ShardSlaveZones != nil {
				shardMap["shard_slave_zones"] = shard.ShardSlaveZones
			}
			if shard.Cpu != nil {
				shardMap["cpu"] = shard.Cpu
			}
			if shard.Range != nil {
				shardMap["range"] = shard.Range
			}

			ids = append(ids, *shard.ShardInstanceId)
			shardList = append(shardList, shardMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		err = d.Set("list", shardList)
		if err != nil {
			log.Printf("[CRITAL]%s set Dcdb shards failed, reason:%+v", logId, err)
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), shardList); e != nil {
			return e
		}
	}

	return nil
}
