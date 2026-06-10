package cfw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwNatFwClusterRegionStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwNatFwClusterRegionStatusRead,
		Schema: map[string]*schema.Schema{
			"nat_cluster_region_status_query_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of query conditions for NAT firewall cluster region status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ccn_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CCN ID.",
						},
						"nat_ins_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "NAT gateway ID.",
						},
						"asset_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Asset type. Valid values: `nat_ccn` (CCN+NAT scenario), `nat` (standalone NAT scenario).",
						},
						"routing_mode": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Traffic steering routing method. 0: multi-route table mode, 1: policy routing mode.",
						},
					},
				},
			},

			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of regions returned.",
			},

			"region_fw_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of regional firewall cluster status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ins_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT gateway ID.",
						},
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CCN ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region, e.g. ap-guangzhou.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region cluster status. Valid values: `NotDeployed` (cluster not deployed), `Deployed` (cluster deployed but traffic steering network not created), `DeployedCustomOnly` (cluster deployed but internal segment covered, need custom traffic steering segment), `Auto` (traffic steering network created with auto-assigned CIDR), `Custom` (traffic steering network created with custom CIDR).",
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic steering network CIDR. Only has value when Status is Auto or Custom.",
						},
						"routing_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Traffic steering routing method. 0: multi-route table mode, 1: policy routing mode.",
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

func dataSourceTencentCloudCfwNatFwClusterRegionStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_nat_fw_cluster_region_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cfwv20190904.NewDescribeNatFwClusterRegionStatusRequest()

	if v, ok := d.GetOk("nat_cluster_region_status_query_list"); ok {
		queryList := v.([]interface{})
		tmpList := make([]*cfwv20190904.NatClusterRegionStatusQuery, 0, len(queryList))
		for _, item := range queryList {
			queryMap := item.(map[string]interface{})
			query := &cfwv20190904.NatClusterRegionStatusQuery{}
			if v, ok := queryMap["ccn_id"].(string); ok && v != "" {
				query.CcnId = helper.String(v)
			}
			if v, ok := queryMap["nat_ins_id"].(string); ok && v != "" {
				query.NatInsId = helper.String(v)
			}
			if v, ok := queryMap["asset_type"].(string); ok && v != "" {
				query.AssetType = helper.String(v)
			}
			if v, ok := queryMap["routing_mode"].(int); ok {
				query.RoutingMode = helper.Int64(int64(v))
			}
			tmpList = append(tmpList, query)
		}
		request.NatClusterRegionStatusQueryList = tmpList
	}

	var response *cfwv20190904.DescribeNatFwClusterRegionStatusResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient()
		result, e := client.DescribeNatFwClusterRegionStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("nat_fw_cluster_region_status DescribeNatFwClusterRegionStatus response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return err
	}

	if response.Response.Total != nil {
		_ = d.Set("total", response.Response.Total)
	}

	regionFwStatusList := make([]map[string]interface{}, 0)
	if response.Response.RegionFwStatus != nil {
		for _, item := range response.Response.RegionFwStatus {
			itemMap := map[string]interface{}{}
			if item.NatInsId != nil {
				itemMap["nat_ins_id"] = item.NatInsId
			}
			if item.CcnId != nil {
				itemMap["ccn_id"] = item.CcnId
			}
			if item.Region != nil {
				itemMap["region"] = item.Region
			}
			if item.Status != nil {
				itemMap["status"] = item.Status
			}
			if item.Cidr != nil {
				itemMap["cidr"] = item.Cidr
			}
			if item.RoutingMode != nil {
				itemMap["routing_mode"] = item.RoutingMode
			}
			regionFwStatusList = append(regionFwStatusList, itemMap)
		}
	}

	_ = d.Set("region_fw_status", regionFwStatusList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
