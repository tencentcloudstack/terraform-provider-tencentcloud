/*
Use this data source to query dayu DDoS policies

Example Usage

```hcl
data "tencentcloud_dayu_ddos_policies" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuDdosPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPoliciesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the DDoS policy to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of DDoS policies. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DDoS policy.",
						},
						"drop_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"drop_tcp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to drop TCP protocol or not.",
									},
									"drop_udp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate to drop UDP protocol or not.",
									},
									"drop_icmp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to drop ICMP protocol or not.",
									},
									"drop_other": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to drop other protocols(exclude TCP/UDP/ICMP) or not.",
									},
									"drop_abroad": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Indicate whether to drop abroad traffic or not.",
									},
									"check_sync_conn": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to check null connection or not.",
									},
									"d_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of new connections based on destination IP.",
									},
									"d_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of concurrent connections based on destination IP.", //?
									},
									"s_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of concurrent connections based on source IP.",
									},
									"s_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of new connections based on source IP.",
									},
									"bad_conn_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of new connections based on destination IP that trigger suppression of connections.",
									},
									"null_conn_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate to enable null connection or not.",
									},
									"conn_timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Connection timeout of abnormal connection check.",
									},
									"syn_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The percentage of syn in ack of abnormal connection check.",
									},
									"syn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of syn of abnormal connection check.",
									},
									"tcp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of TCP traffic.",
									},
									"udp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of UDP traffic rate.",
									},
									"icmp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of ICMP traffic rate.",
									},
									"other_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of other protocols(exclude TCP/UDP/ICMP) traffic rate.",
									},
								},
							},
							Description: "Option list of abnormal check of the DDoS policy.",
						},
						"port_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol.",
									},
									"start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start port.",
									},
									"end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End port.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action of port to take.",
									},
									"kind": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The type of forbidden port, and valid values are 0, 1, 2. 0 for destination port, 1 for source port and 2 for both destination and source posts.",
									},
								},
							},
							Description: "Port limits of abnormal check of the DDoS policy.",
						},
						"black_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateIp,
							},
							Optional:    true,
							Description: "Black ip list.",
						},
						"white_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateIp,
							},
							Optional:    true,
							Description: "White ip list.",
						},
						"packet_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol.",
									},
									"d_start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start port of the destination.",
									},
									"d_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End port of the destination.",
									},
									"s_start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start port of the source.",
									},
									"s_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End port of the source.",
									},
									"pkt_length_min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum length of the packet.",
									},
									"pkt_length_max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The max length of the packet.",
									},
									"match_begin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicate whether to check load or not.",
									},
									"match_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Match type.",
									},
									"match_str": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key word or regular expression.",
									},
									"depth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The depth of match.",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The offset of match.",
									},
									"is_include": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to include the key word/regular expression or not.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action of port to take.",
									},
								},
							},
							Description: "Message filter options list.",
						},
						"watermark_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcp_port_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Port range of TCP.",
									},
									"udp_port_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Port range of TCP.",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The offset of watermark.",
									},
									"auto_remove": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to auto-remove the watermark or not.",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to open watermark or not.",
									},
								},
							},
							Description: "Watermark policy options, and only support one watermark policy at most.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the DDoS policy.",
						},
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of policy case that the DDoS policy works for.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of policy.",
						},
						"watermark_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Id of the watermark.",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content of the watermark.",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether to auto-remove the watermark or not.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Create time of the watermark.",
									},
								},
							},
							Description: "Watermark content.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_ddos_policies.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	resourceType := d.Get("resource_type").(string)
	policyId := d.Get("policy_id").(string)

	policies := make([]*dayu.DDosPolicy, 0)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeDdosPolicies(ctx, resourceType, policyId)
		if err != nil {
			return retryError(err)
		}
		policies = result
		return nil
	})

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(policies))
	ids := make([]string, 0, len(policies))
	for _, ddosPolicy := range policies {
		listItem := make(map[string]interface{})
		listItem["drop_options"] = flattenDdosDropOptionList([]*dayu.DDoSPolicyDropOption{ddosPolicy.DropOptions})
		listItem["port_filters"] = flattenDdosPortLimitList(ddosPolicy.PortLimits)
		listItem["packet_filters"] = flattenDdosPacketFilterList(ddosPolicy.PacketFilters)
		listItem["black_ips"], listItem["white_ips"] = flattenIpBlackWhiteList(ddosPolicy.IpBlackWhiteLists)
		listItem["watermark_filters"] = flattenWaterPrintPolicyList(ddosPolicy.WaterPrint)
		listItem["create_time"] = *ddosPolicy.CreateTime
		listItem["name"] = *ddosPolicy.PolicyName
		listItem["policy_id"] = *ddosPolicy.PolicyId
		listItem["scene_id"] = *ddosPolicy.SceneId
		listItem["watermark_key"] = flattenWaterPrintKeyList(ddosPolicy.WaterKey)
		list = append(list, listItem)
		ids = append(ids, *ddosPolicy.PolicyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
