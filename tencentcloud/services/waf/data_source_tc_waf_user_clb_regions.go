package waf

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
)

func DataSourceTencentCloudWafUserClbRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafUserClbRegionsRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Region list(ap-xxx format).",
			},
			"rich_datas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detail info for region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese description for region.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English description for region.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region code.",
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

func dataSourceTencentCloudWafUserClbRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_user_clb_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userClbRegions *waf.DescribeUserClbWafRegionsResponseParams
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafUserClbRegionsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		userClbRegions = result
		return nil
	})

	if err != nil {
		return err
	}

	if userClbRegions.Data != nil {
		_ = d.Set("data", userClbRegions.Data)
	}

	if userClbRegions.RichDatas != nil {
		tmpList := make([]map[string]interface{}, 0, len(userClbRegions.RichDatas))
		for _, clbWafRegionItem := range userClbRegions.RichDatas {
			clbWafRegionItemMap := map[string]interface{}{}

			if clbWafRegionItem.Id != nil {
				clbWafRegionItemMap["id"] = clbWafRegionItem.Id
			}

			if clbWafRegionItem.Text != nil {
				clbWafRegionItemMap["text"] = clbWafRegionItem.Text
			}

			if clbWafRegionItem.Value != nil {
				clbWafRegionItemMap["value"] = clbWafRegionItem.Value
			}

			if clbWafRegionItem.Code != nil {
				clbWafRegionItemMap["code"] = clbWafRegionItem.Code
			}

			tmpList = append(tmpList, clbWafRegionItemMap)
		}

		_ = d.Set("rich_datas", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
