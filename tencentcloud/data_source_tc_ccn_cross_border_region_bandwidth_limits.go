/*
Use this data source to query detailed information of ccn_cross_border_region_bandwidth_limits

-> **NOTE:** This resource is dedicated to Unicom.

Example Usage

```hcl
data "tencentcloud_ccn_cross_border_region_bandwidth_limits" "ccn_region_bandwidth_limits" {
  filters {
    name   = "source-region"
    values = ["ap-guangzhou"]
  }

  filters {
    name   = "destination-region"
    values = ["ap-shanghai"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCcnCrossBorderRegionBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnCrossBorderRegionBandwidthLimitsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter condition. Currently, only one value is supported. The supported fields, 1)source-region, the value is like ap-guangzhou; 2)destination-region, the value is like ap-shanghai; 3)ccn-ids,cloud network ID array, the value is like ccn-12345678; 4)user-account-id,user account ID, the value is like 12345678.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "attribute name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Value of the field.",
						},
					},
				},
			},

			"ccn_bandwidth_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Info of cross region ccn instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ccn id.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expired time.",
						},
						"region_flow_control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of RegionFlowControl.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "renew flag.",
						},
						"ccn_region_bandwidth_limit": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "bandwidth limit of cross region.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "source region, such as &#39;ap-shanghai&#39;.",
									},
									"destination_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "destination region, such as.",
									},
									"bandwidth_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "bandwidth list(Mbps).",
									},
								},
							},
						},
						"market_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "market id.",
						},
						"user_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user account id.",
						},
						"is_cross_border": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "if cross region.",
						},
						"is_security_lock": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "`true` means locked.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`POSTPAID` or `PREPAID`.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
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

func dataSourceTencentCloudCcnCrossBorderRegionBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ccn_cross_border_region_bandwidth_limits.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var ccnBandwidthSet []*vpc.CcnBandwidth

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCcnRegionBandwidthLimitsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ccnBandwidthSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ccnBandwidthSet))
	tmpList := make([]map[string]interface{}, 0, len(ccnBandwidthSet))

	if ccnBandwidthSet != nil {
		for _, ccnBandwidth := range ccnBandwidthSet {
			ccnBandwidthMap := map[string]interface{}{}

			if ccnBandwidth.CcnId != nil {
				ccnBandwidthMap["ccn_id"] = ccnBandwidth.CcnId
			}

			if ccnBandwidth.CreatedTime != nil {
				ccnBandwidthMap["created_time"] = ccnBandwidth.CreatedTime
			}

			if ccnBandwidth.ExpiredTime != nil {
				ccnBandwidthMap["expired_time"] = ccnBandwidth.ExpiredTime
			}

			if ccnBandwidth.RegionFlowControlId != nil {
				ccnBandwidthMap["region_flow_control_id"] = ccnBandwidth.RegionFlowControlId
			}

			if ccnBandwidth.RenewFlag != nil {
				ccnBandwidthMap["renew_flag"] = ccnBandwidth.RenewFlag
			}

			if ccnBandwidth.CcnRegionBandwidthLimit != nil {
				ccnRegionBandwidthLimitMap := map[string]interface{}{}

				if ccnBandwidth.CcnRegionBandwidthLimit.SourceRegion != nil {
					ccnRegionBandwidthLimitMap["source_region"] = ccnBandwidth.CcnRegionBandwidthLimit.SourceRegion
				}

				if ccnBandwidth.CcnRegionBandwidthLimit.DestinationRegion != nil {
					ccnRegionBandwidthLimitMap["destination_region"] = ccnBandwidth.CcnRegionBandwidthLimit.DestinationRegion
				}

				if ccnBandwidth.CcnRegionBandwidthLimit.BandwidthLimit != nil {
					ccnRegionBandwidthLimitMap["bandwidth_limit"] = ccnBandwidth.CcnRegionBandwidthLimit.BandwidthLimit
				}

				ccnBandwidthMap["ccn_region_bandwidth_limit"] = []interface{}{ccnRegionBandwidthLimitMap}
			}

			if ccnBandwidth.MarketId != nil {
				ccnBandwidthMap["market_id"] = ccnBandwidth.MarketId
			}

			if ccnBandwidth.UserAccountID != nil {
				ccnBandwidthMap["user_account_id"] = ccnBandwidth.UserAccountID
			}

			if ccnBandwidth.IsCrossBorder != nil {
				ccnBandwidthMap["is_cross_border"] = ccnBandwidth.IsCrossBorder
			}

			if ccnBandwidth.IsSecurityLock != nil {
				ccnBandwidthMap["is_security_lock"] = ccnBandwidth.IsSecurityLock
			}

			if ccnBandwidth.InstanceChargeType != nil {
				ccnBandwidthMap["instance_charge_type"] = ccnBandwidth.InstanceChargeType
			}

			if ccnBandwidth.UpdateTime != nil {
				ccnBandwidthMap["update_time"] = ccnBandwidth.UpdateTime
			}

			ids = append(ids, *ccnBandwidth.CcnId)
			tmpList = append(tmpList, ccnBandwidthMap)
		}

		_ = d.Set("ccn_bandwidth_set", tmpList)
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
