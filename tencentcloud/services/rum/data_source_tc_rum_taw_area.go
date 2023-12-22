package rum

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRumTawArea() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumTawAreaRead,
		Schema: map[string]*schema.Schema{
			"area_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Area id.",
			},

			"area_keys": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Area key.",
			},

			"area_statuses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Area status `1`:valid; `2`:invalid.",
			},

			"area_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Area list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Area id.",
						},
						"area_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Area status `1`:&amp;#39;valid&amp;#39;; `2`:&amp;#39;invalid&amp;#39;.",
						},
						"area_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Area name.",
						},
						"area_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Area key.",
						},
						"area_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Area code id.",
						},
						"area_region_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Area code.",
						},
						"area_abbr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region abbreviation.",
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

func dataSourceTencentCloudRumTawAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_rum_taw_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("area_ids"); ok {
		areaIdsSet := v.(*schema.Set).List()
		areaList := make([]*int64, 0)
		for i := range areaIdsSet {
			areaIds := areaIdsSet[i].(int)
			areaList = append(areaList, helper.IntInt64(areaIds))
		}
		paramMap["AreaIds"] = areaList
	}

	if v, ok := d.GetOk("area_keys"); ok {
		areaKeysSet := v.(*schema.Set).List()
		paramMap["AreaKeys"] = helper.InterfacesStringsPoint(areaKeysSet)
	}

	if v, ok := d.GetOk("area_statuses"); ok {
		areaStatusesSet := v.(*schema.Set).List()
		areaList := make([]*int64, 0)
		for i := range areaStatusesSet {
			areaStatuses := areaStatusesSet[i].(int)
			areaList = append(areaList, helper.IntInt64(areaStatuses))
		}
		paramMap["AreaStatuses"] = areaList
	}

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var areaSet []*rum.RumAreaInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumTawAreaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		areaSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(areaSet))
	tmpList := make([]map[string]interface{}, 0, len(areaSet))

	if areaSet != nil {
		for _, rumAreaInfo := range areaSet {
			rumAreaInfoMap := map[string]interface{}{}

			if rumAreaInfo.AreaId != nil {
				rumAreaInfoMap["area_id"] = rumAreaInfo.AreaId
			}

			if rumAreaInfo.AreaStatus != nil {
				rumAreaInfoMap["area_status"] = rumAreaInfo.AreaStatus
			}

			if rumAreaInfo.AreaName != nil {
				rumAreaInfoMap["area_name"] = rumAreaInfo.AreaName
			}

			if rumAreaInfo.AreaKey != nil {
				rumAreaInfoMap["area_key"] = rumAreaInfo.AreaKey
			}

			if rumAreaInfo.AreaRegionID != nil {
				rumAreaInfoMap["area_region_id"] = rumAreaInfo.AreaRegionID
			}

			if rumAreaInfo.AreaRegionCode != nil {
				rumAreaInfoMap["area_region_code"] = rumAreaInfo.AreaRegionCode
			}

			if rumAreaInfo.AreaAbbr != nil {
				rumAreaInfoMap["area_abbr"] = rumAreaInfo.AreaAbbr
			}

			ids = append(ids, strconv.FormatInt(*rumAreaInfo.AreaId, 10))
			tmpList = append(tmpList, rumAreaInfoMap)
		}

		_ = d.Set("area_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
