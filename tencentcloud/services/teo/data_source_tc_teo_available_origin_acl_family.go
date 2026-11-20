package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoAvailableOriginAclFamily() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoAvailableOriginAclFamilyRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the site ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. the maximum value of Filters.Values is 20. if this parameter is left empty, all available origin acl family information under the current site will be returned. detailed filter criteria are as follows: OriginACLFamily: filter by origin acl control domain, such as gaz/mlc/emc/plat-gaz/plat-mlc/plat-emc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Value of the filtered field.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

			"origin_acl_family_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of origin acl family infos.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin protection version number. Format description: standard version: gaz-xxxxx/mlc-xxxxx/emc-xxxxx; streamlined version: plat-gaz-xxxxxx/plat-mlc-xxxxxx/plat-emc-xxxxxx.",
						},
						"active_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version effective time in UTC+8, following the date and time format of the ISO 8601 standard.",
						},
						"entire_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin IP range details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"i_pv4": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IPv4 subnet list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"i_pv6": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IPv6 subnet list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"origin_acl_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin acl control domain. Values: gaz/mlc/emc/plat-gaz/plat-mlc/plat-emc.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTeoAvailableOriginAclFamilyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_available_origin_acl_family.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teov20220901.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := teov20220901.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
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

	var respData []*teov20220901.OriginACLFamilyInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeTeoAvailableOriginACLFamilyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return reqErr
	}

	originACLFamilyInfosList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, info := range respData {
			infoMap := map[string]interface{}{}

			if info.Version != nil {
				infoMap["version"] = info.Version
			}

			if info.ActiveTime != nil {
				infoMap["active_time"] = info.ActiveTime
			}

			entireAddressesMap := map[string]interface{}{}
			if info.EntireAddresses != nil {
				if info.EntireAddresses.IPv4 != nil {
					entireAddressesMap["i_pv4"] = info.EntireAddresses.IPv4
				}

				if info.EntireAddresses.IPv6 != nil {
					entireAddressesMap["i_pv6"] = info.EntireAddresses.IPv6
				}

				infoMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			if info.OriginACLFamily != nil {
				infoMap["origin_acl_family"] = info.OriginACLFamily
			}

			originACLFamilyInfosList = append(originACLFamilyInfosList, infoMap)
		}

		_ = d.Set("origin_acl_family_infos", originACLFamilyInfosList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originACLFamilyInfosList); e != nil {
			return e
		}
	}

	return nil
}
