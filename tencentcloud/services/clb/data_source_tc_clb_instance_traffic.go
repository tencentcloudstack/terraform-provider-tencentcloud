package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstanceTraffic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstanceTrafficRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CLB instance region. If this parameter is not passed in, CLB instances in all regions will be returned.",
			},

			"load_balancer_traffic": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of CLB instances sorted by outbound bandwidth from highest to lowest. Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance ID.",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance region.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance VIP.",
						},
						"out_bandwidth": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Maximum outbound bandwidth in Mbps.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB domain name. Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudClbInstanceTrafficRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instance_traffic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_region"); ok {
		paramMap["LoadBalancerRegion"] = helper.String(v.(string))
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var loadBalancerTraffic []*clb.LoadBalancerTraffic

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbInstanceTraffic(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		loadBalancerTraffic = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(loadBalancerTraffic))
	tmpList := make([]map[string]interface{}, 0, len(loadBalancerTraffic))

	if loadBalancerTraffic != nil {
		for _, loadBalancerTraffic := range loadBalancerTraffic {
			loadBalancerTrafficMap := map[string]interface{}{}

			if loadBalancerTraffic.LoadBalancerId != nil {
				loadBalancerTrafficMap["load_balancer_id"] = loadBalancerTraffic.LoadBalancerId
			}

			if loadBalancerTraffic.LoadBalancerName != nil {
				loadBalancerTrafficMap["load_balancer_name"] = loadBalancerTraffic.LoadBalancerName
			}

			if loadBalancerTraffic.Region != nil {
				loadBalancerTrafficMap["region"] = loadBalancerTraffic.Region
			}

			if loadBalancerTraffic.Vip != nil {
				loadBalancerTrafficMap["vip"] = loadBalancerTraffic.Vip
			}

			if loadBalancerTraffic.OutBandwidth != nil {
				loadBalancerTrafficMap["out_bandwidth"] = loadBalancerTraffic.OutBandwidth
			}

			if loadBalancerTraffic.Domain != nil {
				loadBalancerTrafficMap["domain"] = loadBalancerTraffic.Domain
			}

			ids = append(ids, *loadBalancerTraffic.LoadBalancerId)
			tmpList = append(tmpList, loadBalancerTrafficMap)
		}

		_ = d.Set("load_balancer_traffic", tmpList)
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
