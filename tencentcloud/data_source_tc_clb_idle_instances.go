package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbIdleInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbIdleInstancesRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CLB instance region.",
			},

			"idle_load_balancers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of idle CLBs. Note: This field may return null, indicating that no valid values can be obtained.",
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
						"idle_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason why the load balancer is considered idle. NO_RULES: No rules configured. NO_RS: The rules are not associated with servers.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB instance status, including:0: Creating; 1: Running.",
						},
						"forward": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB type. Value range: 1 (CLB); 0 (classic CLB).",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancing hostname.Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudClbIdleInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_idle_loadbalancers.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_region"); ok {
		paramMap["LoadBalancerRegion"] = helper.String(v.(string))
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var idleLoadBalancers []*clb.IdleLoadBalancer

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbIdleInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		idleLoadBalancers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(idleLoadBalancers))
	tmpList := make([]map[string]interface{}, 0, len(idleLoadBalancers))

	if idleLoadBalancers != nil {
		for _, idleLoadBalancer := range idleLoadBalancers {
			idleLoadBalancerMap := map[string]interface{}{}

			if idleLoadBalancer.LoadBalancerId != nil {
				idleLoadBalancerMap["load_balancer_id"] = idleLoadBalancer.LoadBalancerId
			}

			if idleLoadBalancer.LoadBalancerName != nil {
				idleLoadBalancerMap["load_balancer_name"] = idleLoadBalancer.LoadBalancerName
			}

			if idleLoadBalancer.Region != nil {
				idleLoadBalancerMap["region"] = idleLoadBalancer.Region
			}

			if idleLoadBalancer.Vip != nil {
				idleLoadBalancerMap["vip"] = idleLoadBalancer.Vip
			}

			if idleLoadBalancer.IdleReason != nil {
				idleLoadBalancerMap["idle_reason"] = idleLoadBalancer.IdleReason
			}

			if idleLoadBalancer.Status != nil {
				idleLoadBalancerMap["status"] = idleLoadBalancer.Status
			}

			if idleLoadBalancer.Forward != nil {
				idleLoadBalancerMap["forward"] = idleLoadBalancer.Forward
			}

			if idleLoadBalancer.Domain != nil {
				idleLoadBalancerMap["domain"] = idleLoadBalancer.Domain
			}

			ids = append(ids, *idleLoadBalancer.LoadBalancerId)
			tmpList = append(tmpList, idleLoadBalancerMap)
		}

		_ = d.Set("idle_load_balancers", tmpList)
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
