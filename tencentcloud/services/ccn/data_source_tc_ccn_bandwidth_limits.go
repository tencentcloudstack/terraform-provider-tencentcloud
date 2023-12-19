package ccn

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudCcnBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the CCN to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values
			"limits": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The bandwidth limits of regions:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Limitation of region.",
						},
						"bandwidth_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limitation of bandwidth.",
						},
						"dst_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination area restriction.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCcnBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_bandwidth_limit.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		ccnId = d.Get("ccn_id").(string)
	)

	var infos, err = service.GetCcnRegionBandwidthLimits(ctx, ccnId)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["region"] = item.Region
		infoMap["bandwidth_limit"] = item.BandwidthLimit
		infoMap["dst_region"] = item.DstRegion
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("limits", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  ccn  bandwidth limits fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(ccnId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
