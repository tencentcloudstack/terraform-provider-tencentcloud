package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssXp2pDetailInfoList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssXp2pDetailInfoListRead,
		Schema: map[string]*schema.Schema{
			"query_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The UTC minute granularity query time for querying usage data for a specific minute is in the format: yyyy-mm-ddTHH:MM:00Z. Please refer to the link https://cloud.tencent.com/document/product/266/11732#I.For example, if the local time is 2019-01-08 10:00:00 in Beijing, the corresponding UTC time would be 2019-01-08T10:00:00+08:00.This query supports data from the past six months.",
			},

			"type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The type array can be used to specify the type of media content to be queried. The two available options are live for live streaming and vod for video on demand. If no type is specified, the query will include both live and VOD content by default.",
			},

			"stream_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The stream array can be used to specify the streams to be queried. If no stream is specified, the query will include all streams by default.",
			},

			"dimension": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The dimension parameter can be used to specify the dimension for the query. If this parameter is not passed, the query will default to stream-level data. If you pass this parameter, it will only retrieve data for the specified dimension. The available dimension currently supported is AppId dimension, which allows you to query data based on the application ID. Please note that the returned fields will be related to the specified dimension.",
			},

			"data_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "P2P streaming statistical information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cdn_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CDN traffic.",
						},
						"p2p_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "P2P traffic.",
						},
						"stuck_people": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "People count.",
						},
						"stuck_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Count.",
						},
						"online_people": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Online numbers.",
						},
						"request": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request numbers.",
						},
						"request_success": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Success numbers.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The requested format for time in UTC with one-minute granularity is yyyy-mm-ddTHH:MM:SSZ. This format follows the ISO 8601 standard and is commonly used for representing timestamps in UTC. For more information and examples, you can refer to the link provided: https://cloud.tencent.com/document/product/266/11732#I.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type, divided into two categories: live and vod.Note: This field may return null, indicating that no valid value is available.",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stream ID.Note: This field may return null, indicating that no valid value is available.",
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AppId. Note: This field may return null, indicating that no valid value is available.",
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

func dataSourceTencentCloudCssXp2pDetailInfoListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_xp2p_detail_info_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("query_time"); ok {
		paramMap["QueryTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		typeSet := v.(*schema.Set).List()
		paramMap["Type"] = helper.InterfacesStringsPoint(typeSet)
	}

	if v, ok := d.GetOk("stream_names"); ok {
		streamNamesSet := v.(*schema.Set).List()
		paramMap["StreamNames"] = helper.InterfacesStringsPoint(streamNamesSet)
	}

	if v, ok := d.GetOk("dimension"); ok {
		dimensionSet := v.(*schema.Set).List()
		paramMap["Dimension"] = helper.InterfacesStringsPoint(dimensionSet)
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dataInfoList []*css.XP2PDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssXp2pDetailInfoListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataInfoList))
	tmpList := make([]map[string]interface{}, 0, len(dataInfoList))

	if dataInfoList != nil {
		for _, xP2PDetailInfo := range dataInfoList {
			xP2PDetailInfoMap := map[string]interface{}{}

			if xP2PDetailInfo.CdnBytes != nil {
				xP2PDetailInfoMap["cdn_bytes"] = xP2PDetailInfo.CdnBytes
			}

			if xP2PDetailInfo.P2pBytes != nil {
				xP2PDetailInfoMap["p2p_bytes"] = xP2PDetailInfo.P2pBytes
			}

			if xP2PDetailInfo.StuckPeople != nil {
				xP2PDetailInfoMap["stuck_people"] = xP2PDetailInfo.StuckPeople
			}

			if xP2PDetailInfo.StuckTimes != nil {
				xP2PDetailInfoMap["stuck_times"] = xP2PDetailInfo.StuckTimes
			}

			if xP2PDetailInfo.OnlinePeople != nil {
				xP2PDetailInfoMap["online_people"] = xP2PDetailInfo.OnlinePeople
			}

			if xP2PDetailInfo.Request != nil {
				xP2PDetailInfoMap["request"] = xP2PDetailInfo.Request
			}

			if xP2PDetailInfo.RequestSuccess != nil {
				xP2PDetailInfoMap["request_success"] = xP2PDetailInfo.RequestSuccess
			}

			if xP2PDetailInfo.Time != nil {
				xP2PDetailInfoMap["time"] = xP2PDetailInfo.Time
			}

			if xP2PDetailInfo.Type != nil {
				xP2PDetailInfoMap["type"] = xP2PDetailInfo.Type
			}

			if xP2PDetailInfo.StreamName != nil {
				xP2PDetailInfoMap["stream_name"] = xP2PDetailInfo.StreamName
			}

			if xP2PDetailInfo.AppId != nil {
				xP2PDetailInfoMap["app_id"] = xP2PDetailInfo.AppId
			}

			ids = append(ids, *xP2PDetailInfo.StreamName)
			tmpList = append(tmpList, xP2PDetailInfoMap)
		}

		_ = d.Set("data_info_list", tmpList)
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
