/*
Use this data source to get the available zones in current region.
* Must set product param to fetch the product infomations(e.g. => cvm, vpc)
* By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.
Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "all" {
  product="cvm"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api/v20201106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAvailabilityZonesByProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityZonesByProductRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When specified, only the zone with the exactly name match will be returned.",
			},
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A string variable indicates that the query will use product infomation.",
			},
			"include_unavailable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A bool variable indicates that the query will include `UNAVAILABLE` zones.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values.
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of zones will be exported and its every element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An internal id for the zone, like `200003`, usually not so useful.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the zone, like `ap-guangzhou-3`.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the zone, like `Guangzhou Zone 3`.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the zone, indicate availability using `AVAILABLE` and `UNAVAILABLE` values.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityZonesByProductRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_availability_zones.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	apiService := APIService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var name string
	var product string
	var includeUnavailable = false
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
	}
	if v, ok := d.GetOkExists("include_unavailable"); ok {
		includeUnavailable = v.(bool)
	}

	var zones []*api.ZoneInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		zones, errRet = apiService.DescribeZonesWithProduct(ctx, product)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	zoneList := make([]map[string]interface{}, 0, len(zones))
	ids := make([]string, 0, len(zones))
	for _, zone := range zones {
		if name != "" && name != *zone.Zone {
			continue
		}
		if !includeUnavailable && *zone.ZoneState == ZONE_STATE_UNAVAILABLE {
			continue
		}
		mapping := map[string]interface{}{
			"id":          zone.ZoneId,
			"name":        zone.Zone,
			"description": zone.ZoneName,
			"state":       zone.ZoneState,
		}
		zoneList = append(zoneList, mapping)
		ids = append(ids, *zone.ZoneId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("zones", zoneList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set zones list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), zoneList); err != nil {
			return err
		}
	}

	return nil
}
