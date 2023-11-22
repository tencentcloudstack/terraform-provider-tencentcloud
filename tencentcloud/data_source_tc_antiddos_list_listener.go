/*
Use this data source to query detailed information of antiddos list_listener

# Example Usage

```hcl

data "tencentcloud_antiddos_list_listener" "list_listener" {
}

```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosListListener() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosListListenerRead,
		Schema: map[string]*schema.Schema{
			"layer4_listeners": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "L4 listener list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Origin port, value 1~65535.",
						},
						"frontend_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Forwarding port, value 1~65535.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCP or UDP.",
						},
						"real_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Source server list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_server": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source server addr, ip or domain.",
									},
									"rs_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1: domain, 2: ip.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The return weight of the source station, ranging from 1 to 100.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0~65535.",
									},
								},
							},
						},
						"instance_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Instance ip.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "InstanceId.",
									},
								},
							},
						},
						"instance_detail_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource instance to which the rule belongs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Resource instance ip.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "InstanceId.",
									},
									"cname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance cname.",
									},
								},
							},
						},
					},
				},
			},

			"layer7_listeners": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Layer 7 forwarding listener list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain.",
						},
						"proxy_type_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of forwarding types.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proxy_ports": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "Forwarding listening port list, port value is 1~65535.",
									},
									"proxy_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Http, https.",
									},
								},
							},
						},
						"real_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Source server list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_server": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source server list.",
									},
									"rs_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1: domain, 2: ip.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight: 1-100.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0-65535.",
									},
								},
							},
						},
						"instance_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "InstanceDetails.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Instance ip list.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance id.",
									},
								},
							},
						},
						"instance_detail_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource instance to which the rule belongs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Instance ip list.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance id.",
									},
									"cname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cname.",
									},
								},
							},
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port.",
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

func dataSourceTencentCloudAntiddosListListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_list_listener.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var listListener *antiddos.DescribeListListenerResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosListListenerByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		listListener = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)

	if listListener.Layer4Listeners != nil {
		for _, layer4Rule := range listListener.Layer4Listeners {
			layer4RuleMap := map[string]interface{}{}

			if layer4Rule.BackendPort != nil {
				layer4RuleMap["backend_port"] = layer4Rule.BackendPort
			}

			if layer4Rule.FrontendPort != nil {
				layer4RuleMap["frontend_port"] = layer4Rule.FrontendPort
			}

			if layer4Rule.Protocol != nil {
				layer4RuleMap["protocol"] = layer4Rule.Protocol
			}

			if layer4Rule.RealServers != nil {
				var realServersList []interface{}
				for _, realServers := range layer4Rule.RealServers {
					realServersMap := map[string]interface{}{}

					if realServers.RealServer != nil {
						realServersMap["real_server"] = realServers.RealServer
					}

					if realServers.RsType != nil {
						realServersMap["rs_type"] = realServers.RsType
					}

					if realServers.Weight != nil {
						realServersMap["weight"] = realServers.Weight
					}

					if realServers.Port != nil {
						realServersMap["port"] = realServers.Port
					}

					realServersList = append(realServersList, realServersMap)
				}

				layer4RuleMap["real_servers"] = realServersList
			}

			if layer4Rule.InstanceDetails != nil {
				var instanceDetailsList []interface{}
				for _, instanceDetails := range layer4Rule.InstanceDetails {
					instanceDetailsMap := map[string]interface{}{}

					if instanceDetails.EipList != nil {
						instanceDetailsMap["eip_list"] = instanceDetails.EipList
					}

					if instanceDetails.InstanceId != nil {
						instanceDetailsMap["instance_id"] = instanceDetails.InstanceId
					}

					instanceDetailsList = append(instanceDetailsList, instanceDetailsMap)
				}

				layer4RuleMap["instance_details"] = instanceDetailsList
			}

			if layer4Rule.InstanceDetailRule != nil {
				var instanceDetailRuleList []interface{}
				for _, instanceDetailRule := range layer4Rule.InstanceDetailRule {
					instanceDetailRuleMap := map[string]interface{}{}

					if instanceDetailRule.EipList != nil {
						instanceDetailRuleMap["eip_list"] = instanceDetailRule.EipList
					}

					if instanceDetailRule.InstanceId != nil {
						instanceDetailRuleMap["instance_id"] = instanceDetailRule.InstanceId
					}

					if instanceDetailRule.Cname != nil {
						instanceDetailRuleMap["cname"] = instanceDetailRule.Cname
					}

					instanceDetailRuleList = append(instanceDetailRuleList, instanceDetailRuleMap)
				}

				layer4RuleMap["instance_detail_rule"] = instanceDetailRuleList
			}

			tmpList = append(tmpList, layer4RuleMap)
		}

		_ = d.Set("layer4_listeners", tmpList)
	}

	if listListener.Layer7Listeners != nil {
		for _, layer7Rule := range listListener.Layer7Listeners {
			layer7RuleMap := map[string]interface{}{}

			if layer7Rule.Domain != nil {
				layer7RuleMap["domain"] = layer7Rule.Domain
			}

			if layer7Rule.ProxyTypeList != nil {
				var proxyTypeListList []interface{}
				for _, proxyTypeList := range layer7Rule.ProxyTypeList {
					proxyTypeListMap := map[string]interface{}{}

					if proxyTypeList.ProxyPorts != nil {
						proxyTypeListMap["proxy_ports"] = proxyTypeList.ProxyPorts
					}

					if proxyTypeList.ProxyType != nil {
						proxyTypeListMap["proxy_type"] = proxyTypeList.ProxyType
					}

					proxyTypeListList = append(proxyTypeListList, proxyTypeListMap)
				}

				layer7RuleMap["proxy_type_list"] = proxyTypeListList
			}

			if layer7Rule.RealServers != nil {
				var realServersList []interface{}
				for _, realServers := range layer7Rule.RealServers {
					realServersMap := map[string]interface{}{}

					if realServers.RealServer != nil {
						realServersMap["real_server"] = realServers.RealServer
					}

					if realServers.RsType != nil {
						realServersMap["rs_type"] = realServers.RsType
					}

					if realServers.Weight != nil {
						realServersMap["weight"] = realServers.Weight
					}

					if realServers.Port != nil {
						realServersMap["port"] = realServers.Port
					}

					realServersList = append(realServersList, realServersMap)
				}

				layer7RuleMap["real_servers"] = realServersList
			}

			if layer7Rule.InstanceDetails != nil {
				var instanceDetailsList []interface{}
				for _, instanceDetails := range layer7Rule.InstanceDetails {
					instanceDetailsMap := map[string]interface{}{}

					if instanceDetails.EipList != nil {
						instanceDetailsMap["eip_list"] = instanceDetails.EipList
					}

					if instanceDetails.InstanceId != nil {
						instanceDetailsMap["instance_id"] = instanceDetails.InstanceId
					}

					instanceDetailsList = append(instanceDetailsList, instanceDetailsMap)
				}

				layer7RuleMap["instance_details"] = instanceDetailsList
			}

			if layer7Rule.InstanceDetailRule != nil {
				var instanceDetailRuleList []interface{}
				for _, instanceDetailRule := range layer7Rule.InstanceDetailRule {
					instanceDetailRuleMap := map[string]interface{}{}

					if instanceDetailRule.EipList != nil {
						instanceDetailRuleMap["eip_list"] = instanceDetailRule.EipList
					}

					if instanceDetailRule.InstanceId != nil {
						instanceDetailRuleMap["instance_id"] = instanceDetailRule.InstanceId
					}

					if instanceDetailRule.Cname != nil {
						instanceDetailRuleMap["cname"] = instanceDetailRule.Cname
					}

					instanceDetailRuleList = append(instanceDetailRuleList, instanceDetailRuleMap)
				}

				layer7RuleMap["instance_detail_rule"] = instanceDetailRuleList
			}

			if layer7Rule.Protocol != nil {
				layer7RuleMap["protocol"] = layer7Rule.Protocol
			}

			if layer7Rule.Vport != nil {
				layer7RuleMap["vport"] = layer7Rule.Vport
			}

			tmpList = append(tmpList, layer7RuleMap)
		}

		_ = d.Set("layer7_listeners", tmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
