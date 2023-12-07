package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbZoneRead,
		Schema: map[string]*schema.Schema{
			"include_virtual_zones": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Is virtual zone included.",
			},

			"show_permission": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to display all available zones under the region and display the permissions of each available zone of the user.",
			},

			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region in English.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
						},
						"region_zh": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name in Chinese.",
						},
						"zone_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of available zones for sale.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone name in English.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ZoneId.",
									},
									"zone_zh": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone name in Chinesee.",
									},
									"is_support_serverless": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Does it support serverless clusters, 0:Not supported 1:Support.",
									},
									"is_support_normal": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Does it support normal clusters, 0:Not supported 1:Support.",
									},
									"physical_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Physical zone.",
									},
									"has_permission": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the user have zone permissionsNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"is_whole_rdma_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Is zone Rdma.",
									},
								},
							},
						},
						"db_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database type.",
						},
						"modules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Regional module support.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_disable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Is zone on sale, optional values: yes, no.",
									},
									"module_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Module name.",
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

func dataSourceTencentCloudCynosdbZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("include_virtual_zones"); v != nil {
		paramMap["IncludeVirtualZones"] = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("show_permission"); v != nil {
		paramMap["ShowPermission"] = helper.Bool(v.(bool))
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var regionSet []*cynosdb.SaleRegion
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbZoneByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		regionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionSet))
	tmpList := make([]map[string]interface{}, 0, len(regionSet))
	if regionSet != nil {
		for _, saleRegion := range regionSet {
			saleRegionMap := map[string]interface{}{}

			if saleRegion.Region != nil {
				saleRegionMap["region"] = saleRegion.Region
			}

			if saleRegion.RegionId != nil {
				saleRegionMap["region_id"] = saleRegion.RegionId
			}

			if saleRegion.RegionZh != nil {
				saleRegionMap["region_zh"] = saleRegion.RegionZh
			}

			if saleRegion.ZoneSet != nil {
				zoneSetList := []interface{}{}
				for _, zoneSet := range saleRegion.ZoneSet {
					zoneSetMap := map[string]interface{}{}

					if zoneSet.Zone != nil {
						zoneSetMap["zone"] = zoneSet.Zone
					}

					if zoneSet.ZoneId != nil {
						zoneSetMap["zone_id"] = zoneSet.ZoneId
					}

					if zoneSet.ZoneZh != nil {
						zoneSetMap["zone_zh"] = zoneSet.ZoneZh
					}

					if zoneSet.IsSupportServerless != nil {
						zoneSetMap["is_support_serverless"] = zoneSet.IsSupportServerless
					}

					if zoneSet.IsSupportNormal != nil {
						zoneSetMap["is_support_normal"] = zoneSet.IsSupportNormal
					}

					if zoneSet.PhysicalZone != nil {
						zoneSetMap["physical_zone"] = zoneSet.PhysicalZone
					}

					if zoneSet.HasPermission != nil {
						zoneSetMap["has_permission"] = zoneSet.HasPermission
					}

					if zoneSet.IsWholeRdmaZone != nil {
						zoneSetMap["is_whole_rdma_zone"] = zoneSet.IsWholeRdmaZone
					}

					zoneSetList = append(zoneSetList, zoneSetMap)
				}

				saleRegionMap["zone_set"] = zoneSetList
			}

			if saleRegion.DbType != nil {
				saleRegionMap["db_type"] = saleRegion.DbType
			}

			if saleRegion.Modules != nil {
				modulesList := []interface{}{}
				for _, modules := range saleRegion.Modules {
					modulesMap := map[string]interface{}{}

					if modules.IsDisable != nil {
						modulesMap["is_disable"] = modules.IsDisable
					}

					if modules.ModuleName != nil {
						modulesMap["module_name"] = modules.ModuleName
					}

					modulesList = append(modulesList, modulesMap)
				}

				saleRegionMap["modules"] = modulesList
			}

			ids = append(ids, *saleRegion.Region)
			tmpList = append(tmpList, saleRegionMap)
		}

		_ = d.Set("region_set", tmpList)
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
