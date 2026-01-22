package privatedns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPrivateDnsEndPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPrivateDnsEndPointsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter parameters. Valid values: EndPointName, EndPointId, EndPointServiceId, and EndPointVip.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Array of parameter values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"end_point_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Endpoint list.\nNote: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_point_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint ID.",
						},
						"end_point_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint name.",
						},
						"end_point_service_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint service ID.",
						},
						"end_point_vip_set": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "VIP list of the endpoint.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ap-guangzhou\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag key-value pair collection.\nNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
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

func dataSourceTencentCloudPrivateDnsEndPointsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_private_dns_end_points.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*privatednsIntlv20201028.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := privatednsIntlv20201028.Filter{}
			if v, ok := filtersMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*privatednsIntlv20201028.EndPointInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePrivateDnsEndPointsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		}

		respData = result
		return nil
	})

	if err != nil {
		return err
	}

	var ids []string
	endPointSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, endPointSet := range respData {
			endPointSetMap := map[string]interface{}{}
			if endPointSet.EndPointId != nil {
				endPointSetMap["end_point_id"] = endPointSet.EndPointId
				ids = append(ids, *endPointSet.EndPointId)
			}

			if endPointSet.EndPointName != nil {
				endPointSetMap["end_point_name"] = endPointSet.EndPointName
			}

			if endPointSet.EndPointServiceId != nil {
				endPointSetMap["end_point_service_id"] = endPointSet.EndPointServiceId
			}

			if endPointSet.EndPointVipSet != nil {
				endPointSetMap["end_point_vip_set"] = endPointSet.EndPointVipSet
			}

			if endPointSet.RegionCode != nil {
				endPointSetMap["region_code"] = endPointSet.RegionCode
			}

			tagsList := make([]map[string]interface{}, 0, len(endPointSet.Tags))
			if endPointSet.Tags != nil {
				for _, tags := range endPointSet.Tags {
					tagsMap := map[string]interface{}{}
					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				endPointSetMap["tags"] = tagsList
			}

			endPointSetList = append(endPointSetList, endPointSetMap)
		}

		_ = d.Set("end_point_set", endPointSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), endPointSetList); e != nil {
			return e
		}
	}

	return nil
}
