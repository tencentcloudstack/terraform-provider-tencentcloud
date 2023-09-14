/*
Use this data source to query detailed information of waf waf_infos

Example Usage

```hcl
data "tencentcloud_waf_waf_infos" "example" {
  params {
    load_balancer_id = "lb-A8VF445"
  }
}
```

Or

```hcl
data "tencentcloud_waf_waf_infos" "example" {
  params {
    load_balancer_id = "lb-A8VF445"
    listener_id      = "lbl-nonkgvc2"
    domain_id        = "waf-MPtWPK5Q"
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafWafInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafWafInfosRead,
		Schema: map[string]*schema.Schema{
			"params": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Parameters of interfaces for clb.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Loadbalancer unique ID.If this parameter is not passed, it will operate all listeners of this appid. If this parameter is not empty, it will operate listeners of the LoadBalancer only.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Listener ID of LoadBalancer.",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain unique ID.",
						},
					},
				},
			},
			"host_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Host info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "LoadBalancer info bound by waf.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer ID.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer name.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique ID of listener in LB.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener name.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer IP.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "LoadBalancer port.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer region.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol of listener，http or https.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer zone.",
									},
									"numerical_vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "VPCID for load balancer, public network is -1, and internal network is filled in according to actual conditionsNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network type for load balancerNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain unique ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Waf switch,0 off 1 on.",
						},
						"flow_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "WAF traffic mode, 1 cleaning mode, 0 mirroring mode.",
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

func dataSourceTencentCloudWafWafInfosRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_waf_infos.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		hostList []*waf.ClbHostResult
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("params"); ok {
		paramsSet := v.([]interface{})
		tmpSet := make([]*waf.ClbHostsParams, 0, len(paramsSet))

		for _, item := range paramsSet {
			clbHostsParams := waf.ClbHostsParams{}
			clbHostsParamsMap := item.(map[string]interface{})

			if v, ok := clbHostsParamsMap["load_balancer_id"]; ok {
				clbHostsParams.LoadBalancerId = helper.String(v.(string))
			}

			if v, ok := clbHostsParamsMap["listener_id"]; ok {
				clbHostsParams.ListenerId = helper.String(v.(string))
			}

			if v, ok := clbHostsParamsMap["domain_id"]; ok {
				clbHostsParams.DomainId = helper.String(v.(string))
			}

			tmpSet = append(tmpSet, &clbHostsParams)
		}

		paramMap["Params"] = tmpSet
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafWafInfosByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		hostList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(hostList))
	tmpList := make([]map[string]interface{}, 0, len(hostList))

	if hostList != nil {
		for _, clbHostResult := range hostList {
			clbHostResultMap := map[string]interface{}{}

			if clbHostResult.LoadBalancer != nil {
				loadBalancerMap := map[string]interface{}{}

				if clbHostResult.LoadBalancer.LoadBalancerId != nil {
					loadBalancerMap["load_balancer_id"] = clbHostResult.LoadBalancer.LoadBalancerId
				}

				if clbHostResult.LoadBalancer.LoadBalancerName != nil {
					loadBalancerMap["load_balancer_name"] = clbHostResult.LoadBalancer.LoadBalancerName
				}

				if clbHostResult.LoadBalancer.ListenerId != nil {
					loadBalancerMap["listener_id"] = clbHostResult.LoadBalancer.ListenerId
				}

				if clbHostResult.LoadBalancer.ListenerName != nil {
					loadBalancerMap["listener_name"] = clbHostResult.LoadBalancer.ListenerName
				}

				if clbHostResult.LoadBalancer.Vip != nil {
					loadBalancerMap["vip"] = clbHostResult.LoadBalancer.Vip
				}

				if clbHostResult.LoadBalancer.Vport != nil {
					loadBalancerMap["vport"] = clbHostResult.LoadBalancer.Vport
				}

				if clbHostResult.LoadBalancer.Region != nil {
					loadBalancerMap["region"] = clbHostResult.LoadBalancer.Region
				}

				if clbHostResult.LoadBalancer.Protocol != nil {
					loadBalancerMap["protocol"] = clbHostResult.LoadBalancer.Protocol
				}

				if clbHostResult.LoadBalancer.Zone != nil {
					loadBalancerMap["zone"] = clbHostResult.LoadBalancer.Zone
				}

				if clbHostResult.LoadBalancer.NumericalVpcId != nil {
					loadBalancerMap["numerical_vpc_id"] = clbHostResult.LoadBalancer.NumericalVpcId
				}

				if clbHostResult.LoadBalancer.LoadBalancerType != nil {
					loadBalancerMap["load_balancer_type"] = clbHostResult.LoadBalancer.LoadBalancerType
				}

				clbHostResultMap["load_balancer"] = []interface{}{loadBalancerMap}
			}

			if clbHostResult.Domain != nil {
				clbHostResultMap["domain"] = clbHostResult.Domain
			}

			if clbHostResult.DomainId != nil {
				clbHostResultMap["domain_id"] = clbHostResult.DomainId
			}

			if clbHostResult.Status != nil {
				clbHostResultMap["status"] = clbHostResult.Status
			}

			if clbHostResult.FlowMode != nil {
				clbHostResultMap["flow_mode"] = clbHostResult.FlowMode
			}

			ids = append(ids, *clbHostResult.DomainId)
			tmpList = append(tmpList, clbHostResultMap)
		}

		_ = d.Set("host_list", tmpList)
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
