package teo

import (
	"context"
	"fmt"
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
				Description: "Site ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. The upper limit of Filters.Values is 20. Only `OriginACLFamily` is supported, used to filter by origin ACL control domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name. Valid value: `OriginACLFamily`.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Filter value.",
						},
					},
				},
			},

			"origin_acl_family_info_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Origin ACL family info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin protection version number.",
						},
						"active_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version active time, Beijing time UTC+8, following ISO 8601 standard.",
						},
						"entire_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin IP CIDR details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IPv4 CIDR list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"ipv6": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IPv6 CIDR list.",
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
							Description: "Origin ACL control domain.",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of available origin ACL families.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
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
			filter := teov20220901.Filter{}
			filterMap := item.(map[string]interface{})
			if name, ok := filterMap["name"].(string); ok && name != "" {
				filter.Name = helper.String(name)
			}
			if values, ok := filterMap["values"]; ok {
				valuesSet := values.(*schema.Set).List()
				if len(valuesSet) > 20 {
					return fmt.Errorf("the number of values in filter `%s` cannot exceed 20", *filter.Name)
				}
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	var (
		totalCount          *int64
		originACLFamilyList []*teov20220901.OriginACLFamilyInfo
	)

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		count, result, e := service.DescribeTeoAvailableOriginAclFamilyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if count == nil || result == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("DescribeAvailableOriginACLFamily response is empty"))
		}
		totalCount = count
		originACLFamilyList = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	if totalCount != nil {
		_ = d.Set("total_count", *totalCount)
	}

	originACLFamilyInfoSet := make([]map[string]interface{}, 0, len(originACLFamilyList))
	if originACLFamilyList != nil {
		for _, info := range originACLFamilyList {
			infoMap := map[string]interface{}{}
			if info.Version != nil {
				infoMap["version"] = info.Version
			}
			if info.ActiveTime != nil {
				infoMap["active_time"] = info.ActiveTime
			}
			if info.EntireAddresses != nil {
				entireAddressesMap := map[string]interface{}{}
				if info.EntireAddresses.IPv4 != nil {
					entireAddressesMap["ipv4"] = info.EntireAddresses.IPv4
				}
				if info.EntireAddresses.IPv6 != nil {
					entireAddressesMap["ipv6"] = info.EntireAddresses.IPv6
				}
				infoMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}
			if info.OriginACLFamily != nil {
				infoMap["origin_acl_family"] = info.OriginACLFamily
			}
			originACLFamilyInfoSet = append(originACLFamilyInfoSet, infoMap)
		}
		_ = d.Set("origin_acl_family_info_set", originACLFamilyInfoSet)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originACLFamilyInfoSet); e != nil {
			return e
		}
	}

	return nil
}
