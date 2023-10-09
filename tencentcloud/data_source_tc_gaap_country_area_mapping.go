/*
Use this data source to query detailed information of gaap country area mapping

Example Usage

```hcl
data "tencentcloud_gaap_country_area_mapping" "country_area_mapping" {
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapCountryAreaMapping() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCountryAreaMappingRead,
		Schema: map[string]*schema.Schema{
			"country_area_mapping_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Country/region code mapping table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nation_country_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country name.",
						},
						"nation_country_inner_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country code.",
						},
						"geographical_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"geographical_zone_inner_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region code.",
						},
						"continent_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the continent.",
						},
						"continent_inner_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Continental Code.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Annotation InformationNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudGaapCountryAreaMappingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_country_area_mapping.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}
	var countryAreaMappingList []*gaap.CountryAreaMap

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapCountryAreaMapping(ctx)
		if e != nil {
			return retryError(e)
		}
		countryAreaMappingList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(countryAreaMappingList))
	tmpList := make([]map[string]interface{}, 0, len(countryAreaMappingList))

	if countryAreaMappingList != nil {
		for _, countryAreaMap := range countryAreaMappingList {
			countryAreaMapMap := map[string]interface{}{}

			if countryAreaMap.NationCountryName != nil {
				countryAreaMapMap["nation_country_name"] = countryAreaMap.NationCountryName
			}

			if countryAreaMap.NationCountryInnerCode != nil {
				countryAreaMapMap["nation_country_inner_code"] = countryAreaMap.NationCountryInnerCode
			}

			if countryAreaMap.GeographicalZoneName != nil {
				countryAreaMapMap["geographical_zone_name"] = countryAreaMap.GeographicalZoneName
			}

			if countryAreaMap.GeographicalZoneInnerCode != nil {
				countryAreaMapMap["geographical_zone_inner_code"] = countryAreaMap.GeographicalZoneInnerCode
			}

			if countryAreaMap.ContinentName != nil {
				countryAreaMapMap["continent_name"] = countryAreaMap.ContinentName
			}

			if countryAreaMap.ContinentInnerCode != nil {
				countryAreaMapMap["continent_inner_code"] = countryAreaMap.ContinentInnerCode
			}

			if countryAreaMap.Remark != nil {
				countryAreaMapMap["remark"] = countryAreaMap.Remark
			}

			ids = append(ids, *countryAreaMap.NationCountryInnerCode)
			tmpList = append(tmpList, countryAreaMapMap)
		}

		_ = d.Set("country_area_mapping_list", tmpList)
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
