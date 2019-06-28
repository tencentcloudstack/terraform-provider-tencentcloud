package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCcnBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			// Computed values
			"limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCcnBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "data_source.tencentcloud_ccn_bandwidth_limit.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId string = d.Get("ccn_id").(string)
	)

	var infos, err = service.DescribeCcnRegionBandwidthLimits(ctx, ccnId)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["region"] = item.region
		infoMap["bandwidth_limit"] = item.limit
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("limits", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  ccn  bandwidth limits fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(ccnId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
