/*
Use this data source to query detailed information of dcdb upgrade_price

Example Usage

```hcl
data "tencentcloud_dcdb_upgrade_price" "add_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "ADD"
  add_shard_config {
		shard_count = 2
		shard_memory = 2
		shard_storage = 100
  }
  amount_unit = "pent"
}

data "tencentcloud_dcdb_upgrade_price" "expand_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "EXPAND"

  expand_shard_config {
		shard_instance_ids = ["shard-1b5r04az"]
		shard_memory = 2
		shard_storage = 40
		shard_node_count = 2
  }
  amount_unit = "pent"
}

data "tencentcloud_dcdb_upgrade_price" "split_upgrade_price" {
  instance_id = local.dcdb_id
  upgrade_type = "SPLIT"

  split_shard_config {
		shard_instance_ids = ["shard-1b5r04az"]
		split_rate = 50
		shard_memory = 2
		shard_storage = 100
  }
  amount_unit = "pent"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbUpgradePrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbUpgradePriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"upgrade_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Upgrade type, ADD: add new shard, EXPAND: upgrade the existing shard, SPLIT: split existing shard.",
			},

			"add_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Config for adding new shard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of new shards.",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size in GB.",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity in GB.",
						},
					},
				},
			},

			"expand_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Config for expanding existing shard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of shard ID.",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size in GB.",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity in GB.",
						},
						"shard_node_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Shard node count.",
						},
					},
				},
			},

			"split_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Config for splitting existing shard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of shard ID.",
						},
						"split_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data split ratio, fixed at 50%.",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size in GB.",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity in GB.",
						},
					},
				},
			},

			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. Valid values: `pent` (cent), `microPent` (microcent).",
			},

			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).",
			},

			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).",
			},

			"formula": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Price calculation formula.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcdbUpgradePriceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_upgrade_price.read")()
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

	if v, ok := d.GetOk("upgrade_type"); ok {
		paramMap["UpgradeType"] = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "add_shard_config"); ok {
		addShardConfig := dcdb.AddShardConfig{}
		if v, ok := dMap["shard_count"]; ok {
			addShardConfig.ShardCount = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_memory"]; ok {
			addShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			addShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		paramMap["AddShardConfig"] = &addShardConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "expand_shard_config"); ok {
		expandShardConfig := dcdb.ExpandShardConfig{}
		if v, ok := dMap["shard_instance_ids"]; ok {
			shardInstanceIdsSet := v.(*schema.Set).List()
			expandShardConfig.ShardInstanceIds = helper.InterfacesStringsPoint(shardInstanceIdsSet)
		}
		if v, ok := dMap["shard_memory"]; ok {
			expandShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			expandShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_node_count"]; ok {
			expandShardConfig.ShardNodeCount = helper.IntInt64(v.(int))
		}
		paramMap["ExpandShardConfig"] = &expandShardConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "split_shard_config"); ok {
		splitShardConfig := dcdb.SplitShardConfig{}
		if v, ok := dMap["shard_instance_ids"]; ok {
			shardInstanceIdsSet := v.(*schema.Set).List()
			splitShardConfig.ShardInstanceIds = helper.InterfacesStringsPoint(shardInstanceIdsSet)
		}
		if v, ok := dMap["split_rate"]; ok {
			splitShardConfig.SplitRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_memory"]; ok {
			splitShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			splitShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		paramMap["SplitShardConfig"] = &splitShardConfig
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *dcdb.DescribeDCDBUpgradePriceResponseParams
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeDcdbUpgradePriceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.OriginalPrice != nil {
			_ = d.Set("original_price", result.OriginalPrice)
		}

		if result.Price != nil {
			_ = d.Set("price", result.Price)
		}

		if result.Formula != nil {
			_ = d.Set("formula", result.Formula)
		}
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
