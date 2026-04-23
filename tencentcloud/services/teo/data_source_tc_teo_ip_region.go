package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoIPRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoIPRegionRead,
		Schema: map[string]*schema.Schema{
			"ips": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    100,
				Description: "List of IP addresses (IPv4/IPv6) to query, up to 100 entries.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_region_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "IP region information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address, IPv4 or IPv6.",
						},
						"is_edge_one_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the IP belongs to an EdgeOne node. Values: `yes` (belongs to EdgeOne node), `no` (does not belong to EdgeOne node).",
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

func dataSourceTencentCloudTeoIPRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_ip_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ips"); ok {
		ipsList := v.([]interface{})
		ips := make([]*string, 0, len(ipsList))
		for _, item := range ipsList {
			ips = append(ips, helper.String(item.(string)))
		}
		paramMap["IPs"] = ips
	}

	ipRegionInfo, err := service.DescribeTeoIPRegionByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ipRegionInfo))
	tmpList := make([]map[string]interface{}, 0, len(ipRegionInfo))
	if ipRegionInfo != nil {
		for _, info := range ipRegionInfo {
			infoMap := map[string]interface{}{}
			if info.IP != nil {
				infoMap["ip"] = info.IP
				ids = append(ids, *info.IP)
			}
			if info.IsEdgeOneIP != nil {
				infoMap["is_edge_one_ip"] = info.IsEdgeOneIP
			}
			tmpList = append(tmpList, infoMap)
		}
		_ = d.Set("ip_region_info", tmpList)
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
