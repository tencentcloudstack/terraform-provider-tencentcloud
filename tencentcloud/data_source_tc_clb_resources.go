/*
Use this data source to query detailed information of clb resources

Example Usage

```hcl
data "tencentcloud_clb_resources" "resources" {
  filters {
    name = "isp"
    values = ["BGP"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbResourcesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter to query the list of AZ resources as detailed below: zone - String - Optional - Filter by AZ, such as ap-guangzhou-1. isp -- String - Optional - Filter by the ISP. Values: BGP, CMCC, CUCC and CTCC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value array.",
						},
					},
				},
			},

			"zone_resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of resources supported by the AZ.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary AZ, such as ap-guangzhou-1.",
						},
						"resource_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of resources. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Specific ISP resource information, Vaules: CMCC, CUCC, CTCC, BGP, and INTERNAL.",
									},
									"isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISP information, such as CMCC, CUCC, CTCC, BGP, and INTERNAL.",
									},
									"availability_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Available resources. Note: This field may return null, indicating that no valid values can be obtaine.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specific ISP resource information. Values: CMCC, CUCC, CTCC, BGP.",
												},
												"availability": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether the resource is available. Values: Available, Unavailable.",
												},
											},
										},
									},
								},
							},
						},
						"slave_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secondary AZ, such as ap-guangzhou-2. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version. Values: IPv4, IPv6, and IPv6_Nat.",
						},
						"zone_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the AZ, such as ap-guangzhou.",
						},
						"local_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the AZ is a LocalZone. Values: true, false.",
						},
						"zone_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of resources in the zone. Values: SHARED, EXCLUSIVE.",
						},
						"edge_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the AZ is an edge zone. Values: true, false.",
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

func dataSourceTencentCloudClbResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_resources.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var zoneResourceSet []*clb.ZoneResource

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbResourcesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		zoneResourceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(zoneResourceSet))
	tmpList := make([]map[string]interface{}, 0, len(zoneResourceSet))

	if zoneResourceSet != nil {
		for _, zoneResource := range zoneResourceSet {
			zoneResourceMap := map[string]interface{}{}

			if zoneResource.MasterZone != nil {
				zoneResourceMap["master_zone"] = zoneResource.MasterZone
			}

			if zoneResource.ResourceSet != nil {
				resourceSetList := []interface{}{}
				for _, resourceSet := range zoneResource.ResourceSet {
					resourceSetMap := map[string]interface{}{}

					if resourceSet.Type != nil {
						resourceSetMap["type"] = resourceSet.Type
					}

					if resourceSet.Isp != nil {
						resourceSetMap["isp"] = resourceSet.Isp
					}

					if resourceSet.AvailabilitySet != nil {
						availabilitySetList := []interface{}{}
						for _, availabilitySet := range resourceSet.AvailabilitySet {
							availabilitySetMap := map[string]interface{}{}

							if availabilitySet.Type != nil {
								availabilitySetMap["type"] = availabilitySet.Type
							}

							if availabilitySet.Availability != nil {
								availabilitySetMap["availability"] = availabilitySet.Availability
							}

							availabilitySetList = append(availabilitySetList, availabilitySetMap)
						}

						resourceSetMap["availability_set"] = availabilitySetList
					}

					resourceSetList = append(resourceSetList, resourceSetMap)
				}

				zoneResourceMap["resource_set"] = resourceSetList
			}

			if zoneResource.SlaveZone != nil {
				zoneResourceMap["slave_zone"] = zoneResource.SlaveZone
			}

			if zoneResource.IPVersion != nil {
				zoneResourceMap["ip_version"] = zoneResource.IPVersion
			}

			if zoneResource.ZoneRegion != nil {
				zoneResourceMap["zone_region"] = zoneResource.ZoneRegion
			}

			if zoneResource.LocalZone != nil {
				zoneResourceMap["local_zone"] = zoneResource.LocalZone
			}

			if zoneResource.ZoneResourceType != nil {
				zoneResourceMap["zone_resource_type"] = zoneResource.ZoneResourceType
			}

			if zoneResource.EdgeZone != nil {
				zoneResourceMap["edge_zone"] = zoneResource.EdgeZone
			}

			ids = append(ids, *zoneResource.MasterZone)
			tmpList = append(tmpList, zoneResourceMap)
		}

		_ = d.Set("zone_resource_set", tmpList)
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
