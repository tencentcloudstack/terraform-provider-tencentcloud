/*
Use this data source to get the available zones in the current region. By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.

Example Usage

```hcl
data "tencentcloud_availability_zones" "my_favourite_zone" {
  name = "ap-guangzhou-3"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When specified, only the zone with the exactly name match will return.",
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
							Description: "An internal id for the zone, like `200003`, usually not so useful for end user.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The english name for the zone, like `ap-guangzhou-3`.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description for the zone, unfortunately only Chinese characters at this stage.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state for the zone, indicate availability using `AVAILABLE` and `UNAVAILABLE` values.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_availability_zones.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var name string
	var includeUnavailable = false
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOkExists("include_unavailable"); ok {
		includeUnavailable = v.(bool)
	}

	var zones []*cvm.ZoneInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		zones, errRet = cvmService.DescribeZones(ctx)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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

	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("zones", zoneList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set zone list fail, reason:%s\n ", logId, err.Error())
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
