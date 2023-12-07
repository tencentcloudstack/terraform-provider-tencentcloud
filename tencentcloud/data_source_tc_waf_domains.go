package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafDomainsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Unique ID of Instance.",
			},
			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Domain name.",
			},
			"domains": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance unique ID.",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cname address, used for dns access.",
						},
						"edition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"cls_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable access logs, 1 enable, 0 disable.",
						},
						"flow_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLBWAF traffic mode, 1 cleaning mode, 0 mirroring mode.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Waf switch,0 off 1 on.",
						},
						"mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule defense mode, 0 observation mode, 1 interception mode.",
						},
						"engine": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule and AI Defense Mode, 10 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Shutdown Mode 11 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Observation Mode 12 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Interception Mode 20 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Shutdown Mode 21 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Observation Mode 22 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Interception Mode.",
						},
						"cc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Waf sandbox export addresses, should be added to the whitelist by the upstreams.",
						},
						"rs_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Waf engine export addresses, should be added to the whitelist by the upstreams.",
						},
						"ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Listening ports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nginx_server_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Nginx server ID.",
									},
									"port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listening port.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The listening protocol of listening port.",
									},
									"upstream_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The upstream port for listening port.",
									},
									"upstream_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The upstream protocol for listening port.",
									},
								},
							},
						},
						"load_balancer_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of bound LB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener unique IDNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener nameNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer IDNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer nameNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener protocolNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RegionNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LoadBalancer ipNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Listener portNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Loadbalancer zoneNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"numerical_vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "VPCID for load balancer, public network is -1, and internal network is filled in according to actual conditionsNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Loadbalancer typeNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User appid.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Clbwaf domain name listener status, 0 operation successful, 4 binding LB, 6 unbinding LB, 7 unbinding LB failed, 8 binding LB failed, 10 internal error.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"ipv6_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Ipv6 switch status, 0 off, 1 on.",
						},
						"bot_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BOT switch status, 0 off, 1 on.",
						},
						"level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance level.",
						},
						"post_cls_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable the delivery CLS function, 0 off, 1 on.",
						},
						"post_ckafka_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable the delivery of CKafka function, 0 off, 1 on.",
						},
						"cdc_clusters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cdc clustersNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"api_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "API security switch status, 0 off, 1 onNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"alb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic Source: clb represents Tencent Cloud clb, apisix represents apisix gateway, tsegw represents Tencent Cloud API gateway, default clbNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"sg_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Security group status, 0 does not display, 1 non Tencent cloud source site, 2 security group binding failed, 3 security group changedNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"sg_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detailed explanation of security group statusNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudWafDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_domains.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		domains    []*waf.DomainInfo
		instanceID string
		domain     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceID = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafDomainsByFilter(ctx, instanceID, domain)
		if e != nil {
			return retryError(e)
		}

		domains = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domains))
	tmpList := make([]map[string]interface{}, 0, len(domains))

	if domains != nil {
		for _, domainInfo := range domains {
			domainInfoMap := map[string]interface{}{}

			if domainInfo.Domain != nil {
				domainInfoMap["domain"] = domainInfo.Domain
			}

			if domainInfo.DomainId != nil {
				domainInfoMap["domain_id"] = domainInfo.DomainId
			}

			if domainInfo.InstanceId != nil {
				domainInfoMap["instance_id"] = domainInfo.InstanceId
			}

			if domainInfo.Cname != nil {
				domainInfoMap["cname"] = domainInfo.Cname
			}

			if domainInfo.Edition != nil {
				domainInfoMap["edition"] = domainInfo.Edition
			}

			if domainInfo.Region != nil {
				domainInfoMap["region"] = domainInfo.Region
			}

			if domainInfo.InstanceName != nil {
				domainInfoMap["instance_name"] = domainInfo.InstanceName
			}

			if domainInfo.ClsStatus != nil {
				domainInfoMap["cls_status"] = domainInfo.ClsStatus
			}

			if domainInfo.FlowMode != nil {
				domainInfoMap["flow_mode"] = domainInfo.FlowMode
			}

			if domainInfo.Status != nil {
				domainInfoMap["status"] = domainInfo.Status
			}

			if domainInfo.Mode != nil {
				domainInfoMap["mode"] = domainInfo.Mode
			}

			if domainInfo.Engine != nil {
				domainInfoMap["engine"] = domainInfo.Engine
			}

			if domainInfo.CCList != nil {
				domainInfoMap["cc_list"] = domainInfo.CCList
			}

			if domainInfo.RsList != nil {
				domainInfoMap["rs_list"] = domainInfo.RsList
			}

			if domainInfo.Ports != nil {
				portsList := []interface{}{}
				for _, ports := range domainInfo.Ports {
					portsMap := map[string]interface{}{}

					if ports.NginxServerId != nil {
						portsMap["nginx_server_id"] = ports.NginxServerId
					}

					if ports.Port != nil {
						portsMap["port"] = ports.Port
					}

					if ports.Protocol != nil {
						portsMap["protocol"] = ports.Protocol
					}

					if ports.UpstreamPort != nil {
						portsMap["upstream_port"] = ports.UpstreamPort
					}

					if ports.UpstreamProtocol != nil {
						portsMap["upstream_protocol"] = ports.UpstreamProtocol
					}

					portsList = append(portsList, portsMap)
				}

				domainInfoMap["ports"] = portsList
			}

			if domainInfo.LoadBalancerSet != nil {
				loadBalancerSetList := []interface{}{}
				for _, loadBalancerSet := range domainInfo.LoadBalancerSet {
					loadBalancerSetMap := map[string]interface{}{}

					if loadBalancerSet.ListenerId != nil {
						loadBalancerSetMap["listener_id"] = loadBalancerSet.ListenerId
					}

					if loadBalancerSet.ListenerName != nil {
						loadBalancerSetMap["listener_name"] = loadBalancerSet.ListenerName
					}

					if loadBalancerSet.LoadBalancerId != nil {
						loadBalancerSetMap["load_balancer_id"] = loadBalancerSet.LoadBalancerId
					}

					if loadBalancerSet.LoadBalancerName != nil {
						loadBalancerSetMap["load_balancer_name"] = loadBalancerSet.LoadBalancerName
					}

					if loadBalancerSet.Protocol != nil {
						loadBalancerSetMap["protocol"] = loadBalancerSet.Protocol
					}

					if loadBalancerSet.Region != nil {
						loadBalancerSetMap["region"] = loadBalancerSet.Region
					}

					if loadBalancerSet.Vip != nil {
						loadBalancerSetMap["vip"] = loadBalancerSet.Vip
					}

					if loadBalancerSet.Vport != nil {
						loadBalancerSetMap["vport"] = loadBalancerSet.Vport
					}

					if loadBalancerSet.Zone != nil {
						loadBalancerSetMap["zone"] = loadBalancerSet.Zone
					}

					if loadBalancerSet.NumericalVpcId != nil {
						loadBalancerSetMap["numerical_vpc_id"] = loadBalancerSet.NumericalVpcId
					}

					if loadBalancerSet.LoadBalancerType != nil {
						loadBalancerSetMap["load_balancer_type"] = loadBalancerSet.LoadBalancerType
					}

					loadBalancerSetList = append(loadBalancerSetList, loadBalancerSetMap)
				}

				domainInfoMap["load_balancer_set"] = loadBalancerSetList
			}

			if domainInfo.AppId != nil {
				domainInfoMap["app_id"] = domainInfo.AppId
			}

			if domainInfo.State != nil {
				domainInfoMap["state"] = domainInfo.State
			}

			if domainInfo.CreateTime != nil {
				domainInfoMap["create_time"] = domainInfo.CreateTime
			}

			if domainInfo.Ipv6Status != nil {
				domainInfoMap["ipv6_status"] = domainInfo.Ipv6Status
			}

			if domainInfo.BotStatus != nil {
				domainInfoMap["bot_status"] = domainInfo.BotStatus
			}

			if domainInfo.Level != nil {
				domainInfoMap["level"] = domainInfo.Level
			}

			if domainInfo.PostCLSStatus != nil {
				domainInfoMap["post_cls_status"] = domainInfo.PostCLSStatus
			}

			if domainInfo.PostCKafkaStatus != nil {
				domainInfoMap["post_ckafka_status"] = domainInfo.PostCKafkaStatus
			}

			if domainInfo.CdcClusters != nil {
				domainInfoMap["cdc_clusters"] = domainInfo.CdcClusters
			}

			if domainInfo.ApiStatus != nil {
				domainInfoMap["api_status"] = domainInfo.ApiStatus
			}

			if domainInfo.AlbType != nil {
				domainInfoMap["alb_type"] = domainInfo.AlbType
			}

			if domainInfo.SgState != nil {
				domainInfoMap["sg_state"] = domainInfo.SgState
			}

			if domainInfo.SgDetail != nil {
				domainInfoMap["sg_detail"] = domainInfo.SgDetail
			}

			ids = append(ids, *domainInfo.DomainId)
			tmpList = append(tmpList, domainInfoMap)
		}

		_ = d.Set("domains", tmpList)
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
