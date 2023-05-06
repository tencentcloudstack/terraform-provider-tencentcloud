/*
Use this data source to query detailed information of teo securityPolicyRegions

Example Usage

```hcl
data "tencentcloud_teo_security_policy_regions" "securityPolicyRegions" {
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoSecurityPolicyRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoSecurityPolicyRegionsRead,
		Schema: map[string]*schema.Schema{
			"geo_ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Region info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the country.",
						},
						"continent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the continent.",
						},
						"province": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Province of the region. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTeoSecurityPolicyRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_security_policy_regions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var geoIps []*teo.GeoIp
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoSecurityPolicyRegionsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		geoIps = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo rules failed, reason:%+v", logId, err)
		return err
	}

	ruleList := []interface{}{}
	if geoIps != nil {
		for _, geoIp := range geoIps {
			ruleMap := map[string]interface{}{}
			if geoIp.RegionId != nil {
				ruleMap["region_id"] = geoIp.RegionId
			}
			if geoIp.Country != nil {
				ruleMap["country"] = geoIp.Country
			}
			if geoIp.Continent != nil {
				ruleMap["continent"] = geoIp.Continent
			}
			if geoIp.Province != nil {
				ruleMap["province"] = geoIp.Province
			}

			ruleList = append(ruleList, ruleMap)
		}
		_ = d.Set("geo_ip", ruleList)
	}

	d.SetId("security_policy_regions")

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), ruleList); e != nil {
			return e
		}
	}
	return nil
}
