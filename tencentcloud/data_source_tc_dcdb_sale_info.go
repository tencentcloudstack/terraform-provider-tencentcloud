package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbSaleInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbSaleInfoRead,
		Schema: map[string]*schema.Schema{
			"region_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of sale region info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region name(en).",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "region id.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region name(zh).",
						},
						"zone_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of az zone.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "zone name(en).",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "zone id.",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "zone name(zh).",
									},
									"on_sale": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "is zone on sale.",
									},
								},
							},
						},
						"available_choice": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "available zone choice.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_zone": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "master zone.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "zone name(en).",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "zone id.",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "zone name(zh).",
												},
												"on_sale": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "is zone on sale.",
												},
											},
										},
									},
									"slave_zones": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "slave zones.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "zone name(en).",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "zone id.",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "zone name(zh).",
												},
												"on_sale": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "is zone on sale.",
												},
											},
										},
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

func dataSourceTencentCloudDcdbSaleInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_sale_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var regionList []*dcdb.RegionInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbSaleInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		regionList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionList))
	tmpList := make([]map[string]interface{}, 0, len(regionList))

	if regionList != nil {
		for _, regionInfo := range regionList {
			regionInfoMap := map[string]interface{}{}

			if regionInfo.Region != nil {
				regionInfoMap["region"] = regionInfo.Region
			}

			if regionInfo.RegionId != nil {
				regionInfoMap["region_id"] = regionInfo.RegionId
			}

			if regionInfo.RegionName != nil {
				regionInfoMap["region_name"] = regionInfo.RegionName
			}

			if regionInfo.ZoneList != nil {
				zoneListList := []interface{}{}
				for _, zoneList := range regionInfo.ZoneList {
					zoneListMap := map[string]interface{}{}

					if zoneList.Zone != nil {
						zoneListMap["zone"] = zoneList.Zone
					}

					if zoneList.ZoneId != nil {
						zoneListMap["zone_id"] = zoneList.ZoneId
					}

					if zoneList.ZoneName != nil {
						zoneListMap["zone_name"] = zoneList.ZoneName
					}

					if zoneList.OnSale != nil {
						zoneListMap["on_sale"] = zoneList.OnSale
					}

					zoneListList = append(zoneListList, zoneListMap)
				}

				regionInfoMap["zone_list"] = zoneListList
			}

			if regionInfo.AvailableChoice != nil {
				availableChoiceList := []interface{}{}
				for _, availableChoice := range regionInfo.AvailableChoice {
					availableChoiceMap := map[string]interface{}{}

					if availableChoice.MasterZone != nil {
						masterZoneMap := map[string]interface{}{}

						if availableChoice.MasterZone.Zone != nil {
							masterZoneMap["zone"] = availableChoice.MasterZone.Zone
						}

						if availableChoice.MasterZone.ZoneId != nil {
							masterZoneMap["zone_id"] = availableChoice.MasterZone.ZoneId
						}

						if availableChoice.MasterZone.ZoneName != nil {
							masterZoneMap["zone_name"] = availableChoice.MasterZone.ZoneName
						}

						if availableChoice.MasterZone.OnSale != nil {
							masterZoneMap["on_sale"] = availableChoice.MasterZone.OnSale
						}

						availableChoiceMap["master_zone"] = []interface{}{masterZoneMap}
					}

					if availableChoice.SlaveZones != nil {
						slaveZonesList := []interface{}{}
						for _, slaveZones := range availableChoice.SlaveZones {
							slaveZonesMap := map[string]interface{}{}

							if slaveZones.Zone != nil {
								slaveZonesMap["zone"] = slaveZones.Zone
							}

							if slaveZones.ZoneId != nil {
								slaveZonesMap["zone_id"] = slaveZones.ZoneId
								ids = append(ids, helper.Int64ToStr(*slaveZones.ZoneId))
							}

							if slaveZones.ZoneName != nil {
								slaveZonesMap["zone_name"] = slaveZones.ZoneName
							}

							if slaveZones.OnSale != nil {
								slaveZonesMap["on_sale"] = slaveZones.OnSale
							}

							slaveZonesList = append(slaveZonesList, slaveZonesMap)
						}

						availableChoiceMap["slave_zones"] = slaveZonesList
					}

					availableChoiceList = append(availableChoiceList, availableChoiceMap)
				}

				regionInfoMap["available_choice"] = availableChoiceList
			}

			tmpList = append(tmpList, regionInfoMap)
		}

		_ = d.Set("region_list", tmpList)
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
