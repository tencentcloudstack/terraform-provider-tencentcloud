package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigListAggregators() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigListAggregatorsRead,
		Schema: map[string]*schema.Schema{
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of aggregators.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Aggregator list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aggregator name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aggregator description.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Owner UIN.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"account_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of accounts in the aggregator.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aggregator type.",
						},
						"account_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aggregator ID.",
						},
						"aggregator_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Aggregator status.",
						},
						"member_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member name.",
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

func dataSourceTencentCloudConfigListAggregatorsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_list_aggregators.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	var (
		respItems []map[string]interface{}
		respTotal int
	)

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		items, total, e := service.DescribeConfigListAggregatorsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(items) == 0 && total == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId, config_list_aggregators paramMap=%v", paramMap)
			return resource.NonRetryableError(fmt.Errorf("DescribeConfigListAggregatorsByFilter return empty"))
		}

		respItems = flattenConfigListAggregatorsList(items)
		respTotal = int(total)
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	_ = d.Set("items", respItems)
	_ = d.Set("total", respTotal)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func flattenConfigListAggregatorsList(items []*configv20220802.Aggregator) []map[string]interface{} {
	aggregatorList := make([]map[string]interface{}, 0, len(items))
	for _, aggregator := range items {
		aggregatorMap := map[string]interface{}{}

		if aggregator.Name != nil {
			aggregatorMap["name"] = aggregator.Name
		}

		if aggregator.Description != nil {
			aggregatorMap["description"] = aggregator.Description
		}

		if aggregator.OwnerUin != nil {
			aggregatorMap["owner_uin"] = int(*aggregator.OwnerUin)
		}

		if aggregator.CreateTime != nil {
			aggregatorMap["create_time"] = aggregator.CreateTime
		}

		if aggregator.AccountCount != nil {
			aggregatorMap["account_count"] = int(*aggregator.AccountCount)
		}

		if aggregator.Type != nil {
			aggregatorMap["type"] = aggregator.Type
		}

		if aggregator.AccountGroupId != nil {
			aggregatorMap["account_group_id"] = aggregator.AccountGroupId
		}

		if aggregator.AggregatorStatus != nil {
			aggregatorMap["aggregator_status"] = int(*aggregator.AggregatorStatus)
		}

		if aggregator.MemberName != nil {
			aggregatorMap["member_name"] = aggregator.MemberName
		}

		aggregatorList = append(aggregatorList, aggregatorMap)
	}

	return aggregatorList
}
