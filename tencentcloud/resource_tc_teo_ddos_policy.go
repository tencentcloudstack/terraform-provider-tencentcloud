/*
Provides a resource to create a teo ddos_policy

Example Usage

```hcl
resource "tencentcloud_teo_ddos_policy" "ddos_policy" {
  zone_id = &lt;nil&gt;
  policy_id = &lt;nil&gt;
  d_do_s_rule {
		switch = &lt;nil&gt;
		status_info {
			ply_level = &lt;nil&gt;
		}
		geo_ip {
			region_ids = &lt;nil&gt;
			switch = &lt;nil&gt;
		}
		allow_block {
			switch = &lt;nil&gt;
			allow_block_ips {
				ip = &lt;nil&gt;
				type = &lt;nil&gt;
			}
		}
		anti_ply {
			drop_tcp = &lt;nil&gt;
			drop_udp = &lt;nil&gt;
			drop_icmp = &lt;nil&gt;
			drop_other = &lt;nil&gt;
			source_create_limit = &lt;nil&gt;
			source_connect_limit = &lt;nil&gt;
			destination_create_limit = &lt;nil&gt;
			destination_connect_limit = &lt;nil&gt;
			abnormal_connect_num = &lt;nil&gt;
			abnormal_syn_ratio = &lt;nil&gt;
			abnormal_syn_num = &lt;nil&gt;
			connect_timeout = &lt;nil&gt;
			empty_connect_protect = &lt;nil&gt;
			udp_shard = &lt;nil&gt;
		}
		packet_filter {
			switch = &lt;nil&gt;
			packet_filters {
				action = &lt;nil&gt;
				protocol = &lt;nil&gt;
				dport_start = &lt;nil&gt;
				dport_end = &lt;nil&gt;
				packet_min = &lt;nil&gt;
				packet_max = &lt;nil&gt;
				sport_start = &lt;nil&gt;
				sport_end = &lt;nil&gt;
				match_type = &lt;nil&gt;
				is_not = &lt;nil&gt;
				offset = &lt;nil&gt;
				depth = &lt;nil&gt;
				match_begin = &lt;nil&gt;
				str = &lt;nil&gt;
				match_type2 = &lt;nil&gt;
				is_not2 = &lt;nil&gt;
				offset2 = &lt;nil&gt;
				depth2 = &lt;nil&gt;
				match_begin2 = &lt;nil&gt;
				str2 = &lt;nil&gt;
				match_logic = &lt;nil&gt;
			}
		}
		acl {
			switch = &lt;nil&gt;
			acls {
				dport_end = &lt;nil&gt;
				dport_start = &lt;nil&gt;
				sport_end = &lt;nil&gt;
				sport_start = &lt;nil&gt;
				protocol = &lt;nil&gt;
				action = &lt;nil&gt;
			}
		}
		speed_limit {
			package_limit = &lt;nil&gt;
			flux_limit = &lt;nil&gt;
		}

  }
}
```

Import

teo ddos_policy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_ddos_policy.ddos_policy ddos_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTeoDdosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDdosPolicyCreate,
		Read:   resourceTencentCloudTeoDdosPolicyRead,
		Update: resourceTencentCloudTeoDdosPolicyUpdate,
		Delete: resourceTencentCloudTeoDdosPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"policy_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},

			"d_do_s_rule": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "DDoS Configuration of the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "DDoS protection switch. Valid values:- `on`: Enable.- `off`: Disable.",
						},
						"status_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS protection level.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ply_level": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Policy level. Valid values:- `low`: loose.- `middle`: moderate.- `high`: strict.",
									},
								},
							},
						},
						"geo_ip": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS Protection by Geo Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Optional:    true,
										Computed:    true,
										Description: "Region ID. See details in data source `security_policy_regions`.",
									},
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
								},
							},
						},
						"allow_block": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS black-white list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable. `AllowBlockIps` parameter is required.- `off`: Disable.",
									},
									"allow_block_ips": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "DDoS black-white list detail.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Valid value format:- ip, for example 1.1.1.1- ip range, for example 1.1.1.2-1.1.1.3- network segment, for example 1.2.1.0/24- network segment range, for example 1.2.1.0/24-1.2.2.0/24.",
												},
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Valid values: `block`, `allow`.",
												},
												"update_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Last modification date.",
												},
											},
										},
									},
								},
							},
						},
						"anti_ply": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS protocol and connection protection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"drop_tcp": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Block TCP protocol. Valid values: `on`, `off`.",
									},
									"drop_udp": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Block UDP protocol. Valid values: `on`, `off`.",
									},
									"drop_icmp": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Block ICMP protocol. Valid values: `on`, `off`.",
									},
									"drop_other": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Block other protocols. Valid values: `on`, `off`.",
									},
									"source_create_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Limitation of new connection to origin site per second. Valid value range: 0-4294967295.",
									},
									"source_connect_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Limitation of connections to origin site. Valid value range: 0-4294967295.",
									},
									"destination_create_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Limitation of new connection to dest port per second. Valid value range: 0-4294967295.",
									},
									"destination_connect_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Limitation of connections to dest port. Valid value range: 0-4294967295.",
									},
									"abnormal_connect_num": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Abnormal connections threshold. Valid value range: 0-4294967295.",
									},
									"abnormal_syn_ratio": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Abnormal syn packet ratio threshold. Valid value range: 0-100.",
									},
									"abnormal_syn_num": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Abnormal syn packet number threshold. Valid value range: 0-65535.",
									},
									"connect_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Connection timeout detection per second. Valid value range: 0-65535.",
									},
									"empty_connect_protect": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Empty connection protection switch. Valid values: `on`, `off`.",
									},
									"udp_shard": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "UDP shard protection switch. Valid values: `on`, `off`.",
									},
								},
							},
						},
						"packet_filter": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS feature filtering configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable. `PacketFilters` parameter is required.- `off`: Disable.",
									},
									"packet_filters": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "DDoS feature filtering configuration detail.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Action to take. Valid values: `drop`, `transmit`, `drop_block`, `forward`.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Valid value: `tcp`, `udp`, `icmp`, `all`.",
												},
												"dport_start": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start of the dest port range. Valid value range: 0-65535.",
												},
												"dport_end": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "End of the dest port range. Valid value range: 0-65535.",
												},
												"packet_min": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Min packet size. Valid value range: 0-1500.",
												},
												"packet_max": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Max packet size. Valid value range: 0-1500.",
												},
												"sport_start": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start of the source port range. Valid value range: 0-65535.",
												},
												"sport_end": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "End of the source port range. Valid value range: 0-65535.",
												},
												"match_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Match type of feature 1. Valid values:- `pcre`: regex expression.- `sunday`: string match.",
												},
												"is_not": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Negate the match condition of feature 1. Valid values:- `0`: match.- `1`: not match.",
												},
												"offset": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Offset of feature 1. Valid value range: 1-1500.",
												},
												"depth": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Packet character depth to check of feature 1. Valid value range: 1-1500.",
												},
												"match_begin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Packet layer for matching begin of feature 1. Valid values:- `begin_l5`: matching from packet payload.- `begin_l4`: matching from TCP/UDP header.- `begin_l3`: matching from IP header.",
												},
												"str": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Regex expression or string to match.",
												},
												"match_type2": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Match type of feature 2. Valid values:- `pcre`: regex expression.- `sunday`: string match.",
												},
												"is_not2": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Negate the match condition of feature 2. Valid values:- `0`: match.- `1`: not match.",
												},
												"offset2": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Offset of feature 2. Valid value range: 1-1500.",
												},
												"depth2": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Packet character depth to check of feature 2. Valid value range: 1-1500.",
												},
												"match_begin2": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Packet layer for matching begin of feature 2. Valid values:- `begin_l5`: matching from packet payload.- `begin_l4`: matching from TCP/UDP header.- `begin_l3`: matching from IP header.",
												},
												"str2": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Regex expression or string to match.",
												},
												"match_logic": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Relation between multi features. Valid values: `and`, `or`, `none` (only feature 1 is used).",
												},
											},
										},
									},
								},
							},
						},
						"acl": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS ACL rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable. `Acl` parameter is require.- `off`: Disable.",
									},
									"acls": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "DDoS ACL rule configuration detail.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dport_end": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "End of the dest port range. Valid value range: 0-65535.",
												},
												"dport_start": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start of the dest port range. Valid value range: 0-65535.",
												},
												"sport_end": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "End of the source port range. Valid value range: 0-65535.",
												},
												"sport_start": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start of the source port range. Valid value range: 0-65535.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Valid values: `tcp`, `udp`, `all`.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Action to take. Valid values: `drop`, `transmit`, `forward`.",
												},
											},
										},
									},
								},
							},
						},
						"speed_limit": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "DDoS access origin site speed limit configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package_limit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Limit the number of packages. Valid range: 1 pps-10000 Gpps, 0 means no limitation, supported units: pps、Kpps、Mpps、Gpps.",
									},
									"flux_limit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Limit the number of fluxes. Valid range: 1 bps-10000 Gbps, 0 means no limitation, supported units: bps、Kbps、Mbps、Gbps.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDdosPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.create")()
	defer inconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	var policyId int64
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int64)
	}

	d.SetId(strings.Join([]string{zoneId, helper.Int64ToStr(policyId)}, FILED_SP))

	return resourceTencentCloudTeoDdosPolicyUpdate(d, meta)
}

func resourceTencentCloudTeoDdosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	policyId := idSplit[1]

	ddosPolicy, err := service.DescribeTeoDdosPolicyById(ctx, zoneId, policyId)
	if err != nil {
		return err
	}

	if ddosPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoDdosPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ddosPolicy.ZoneId != nil {
		_ = d.Set("zone_id", ddosPolicy.ZoneId)
	}

	if ddosPolicy.PolicyId != nil {
		_ = d.Set("policy_id", ddosPolicy.PolicyId)
	}

	if ddosPolicy.DDoSRule != nil {
		dDoSRuleMap := map[string]interface{}{}

		if ddosPolicy.DDoSRule.Switch != nil {
			dDoSRuleMap["switch"] = ddosPolicy.DDoSRule.Switch
		}

		if ddosPolicy.DDoSRule.StatusInfo != nil {
			statusInfoMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.StatusInfo.PlyLevel != nil {
				statusInfoMap["ply_level"] = ddosPolicy.DDoSRule.StatusInfo.PlyLevel
			}

			dDoSRuleMap["status_info"] = []interface{}{statusInfoMap}
		}

		if ddosPolicy.DDoSRule.GeoIp != nil {
			geoIpMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.GeoIp.RegionIds != nil {
				geoIpMap["region_ids"] = ddosPolicy.DDoSRule.GeoIp.RegionIds
			}

			if ddosPolicy.DDoSRule.GeoIp.Switch != nil {
				geoIpMap["switch"] = ddosPolicy.DDoSRule.GeoIp.Switch
			}

			dDoSRuleMap["geo_ip"] = []interface{}{geoIpMap}
		}

		if ddosPolicy.DDoSRule.AllowBlock != nil {
			allowBlockMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.AllowBlock.Switch != nil {
				allowBlockMap["switch"] = ddosPolicy.DDoSRule.AllowBlock.Switch
			}

			if ddosPolicy.DDoSRule.AllowBlock.AllowBlockIps != nil {
				allowBlockIpsList := []interface{}{}
				for _, allowBlockIps := range ddosPolicy.DDoSRule.AllowBlock.AllowBlockIps {
					allowBlockIpsMap := map[string]interface{}{}

					if allowBlockIps.Ip != nil {
						allowBlockIpsMap["ip"] = allowBlockIps.Ip
					}

					if allowBlockIps.Type != nil {
						allowBlockIpsMap["type"] = allowBlockIps.Type
					}

					if allowBlockIps.UpdateTime != nil {
						allowBlockIpsMap["update_time"] = allowBlockIps.UpdateTime
					}

					allowBlockIpsList = append(allowBlockIpsList, allowBlockIpsMap)
				}

				allowBlockMap["allow_block_ips"] = []interface{}{allowBlockIpsList}
			}

			dDoSRuleMap["allow_block"] = []interface{}{allowBlockMap}
		}

		if ddosPolicy.DDoSRule.AntiPly != nil {
			antiPlyMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.AntiPly.DropTcp != nil {
				antiPlyMap["drop_tcp"] = ddosPolicy.DDoSRule.AntiPly.DropTcp
			}

			if ddosPolicy.DDoSRule.AntiPly.DropUdp != nil {
				antiPlyMap["drop_udp"] = ddosPolicy.DDoSRule.AntiPly.DropUdp
			}

			if ddosPolicy.DDoSRule.AntiPly.DropIcmp != nil {
				antiPlyMap["drop_icmp"] = ddosPolicy.DDoSRule.AntiPly.DropIcmp
			}

			if ddosPolicy.DDoSRule.AntiPly.DropOther != nil {
				antiPlyMap["drop_other"] = ddosPolicy.DDoSRule.AntiPly.DropOther
			}

			if ddosPolicy.DDoSRule.AntiPly.SourceCreateLimit != nil {
				antiPlyMap["source_create_limit"] = ddosPolicy.DDoSRule.AntiPly.SourceCreateLimit
			}

			if ddosPolicy.DDoSRule.AntiPly.SourceConnectLimit != nil {
				antiPlyMap["source_connect_limit"] = ddosPolicy.DDoSRule.AntiPly.SourceConnectLimit
			}

			if ddosPolicy.DDoSRule.AntiPly.DestinationCreateLimit != nil {
				antiPlyMap["destination_create_limit"] = ddosPolicy.DDoSRule.AntiPly.DestinationCreateLimit
			}

			if ddosPolicy.DDoSRule.AntiPly.DestinationConnectLimit != nil {
				antiPlyMap["destination_connect_limit"] = ddosPolicy.DDoSRule.AntiPly.DestinationConnectLimit
			}

			if ddosPolicy.DDoSRule.AntiPly.AbnormalConnectNum != nil {
				antiPlyMap["abnormal_connect_num"] = ddosPolicy.DDoSRule.AntiPly.AbnormalConnectNum
			}

			if ddosPolicy.DDoSRule.AntiPly.AbnormalSynRatio != nil {
				antiPlyMap["abnormal_syn_ratio"] = ddosPolicy.DDoSRule.AntiPly.AbnormalSynRatio
			}

			if ddosPolicy.DDoSRule.AntiPly.AbnormalSynNum != nil {
				antiPlyMap["abnormal_syn_num"] = ddosPolicy.DDoSRule.AntiPly.AbnormalSynNum
			}

			if ddosPolicy.DDoSRule.AntiPly.ConnectTimeout != nil {
				antiPlyMap["connect_timeout"] = ddosPolicy.DDoSRule.AntiPly.ConnectTimeout
			}

			if ddosPolicy.DDoSRule.AntiPly.EmptyConnectProtect != nil {
				antiPlyMap["empty_connect_protect"] = ddosPolicy.DDoSRule.AntiPly.EmptyConnectProtect
			}

			if ddosPolicy.DDoSRule.AntiPly.UdpShard != nil {
				antiPlyMap["udp_shard"] = ddosPolicy.DDoSRule.AntiPly.UdpShard
			}

			dDoSRuleMap["anti_ply"] = []interface{}{antiPlyMap}
		}

		if ddosPolicy.DDoSRule.PacketFilter != nil {
			packetFilterMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.PacketFilter.Switch != nil {
				packetFilterMap["switch"] = ddosPolicy.DDoSRule.PacketFilter.Switch
			}

			if ddosPolicy.DDoSRule.PacketFilter.PacketFilters != nil {
				packetFiltersList := []interface{}{}
				for _, packetFilters := range ddosPolicy.DDoSRule.PacketFilter.PacketFilters {
					packetFiltersMap := map[string]interface{}{}

					if packetFilters.Action != nil {
						packetFiltersMap["action"] = packetFilters.Action
					}

					if packetFilters.Protocol != nil {
						packetFiltersMap["protocol"] = packetFilters.Protocol
					}

					if packetFilters.DportStart != nil {
						packetFiltersMap["dport_start"] = packetFilters.DportStart
					}

					if packetFilters.DportEnd != nil {
						packetFiltersMap["dport_end"] = packetFilters.DportEnd
					}

					if packetFilters.PacketMin != nil {
						packetFiltersMap["packet_min"] = packetFilters.PacketMin
					}

					if packetFilters.PacketMax != nil {
						packetFiltersMap["packet_max"] = packetFilters.PacketMax
					}

					if packetFilters.SportStart != nil {
						packetFiltersMap["sport_start"] = packetFilters.SportStart
					}

					if packetFilters.SportEnd != nil {
						packetFiltersMap["sport_end"] = packetFilters.SportEnd
					}

					if packetFilters.MatchType != nil {
						packetFiltersMap["match_type"] = packetFilters.MatchType
					}

					if packetFilters.IsNot != nil {
						packetFiltersMap["is_not"] = packetFilters.IsNot
					}

					if packetFilters.Offset != nil {
						packetFiltersMap["offset"] = packetFilters.Offset
					}

					if packetFilters.Depth != nil {
						packetFiltersMap["depth"] = packetFilters.Depth
					}

					if packetFilters.MatchBegin != nil {
						packetFiltersMap["match_begin"] = packetFilters.MatchBegin
					}

					if packetFilters.Str != nil {
						packetFiltersMap["str"] = packetFilters.Str
					}

					if packetFilters.MatchType2 != nil {
						packetFiltersMap["match_type2"] = packetFilters.MatchType2
					}

					if packetFilters.IsNot2 != nil {
						packetFiltersMap["is_not2"] = packetFilters.IsNot2
					}

					if packetFilters.Offset2 != nil {
						packetFiltersMap["offset2"] = packetFilters.Offset2
					}

					if packetFilters.Depth2 != nil {
						packetFiltersMap["depth2"] = packetFilters.Depth2
					}

					if packetFilters.MatchBegin2 != nil {
						packetFiltersMap["match_begin2"] = packetFilters.MatchBegin2
					}

					if packetFilters.Str2 != nil {
						packetFiltersMap["str2"] = packetFilters.Str2
					}

					if packetFilters.MatchLogic != nil {
						packetFiltersMap["match_logic"] = packetFilters.MatchLogic
					}

					packetFiltersList = append(packetFiltersList, packetFiltersMap)
				}

				packetFilterMap["packet_filters"] = []interface{}{packetFiltersList}
			}

			dDoSRuleMap["packet_filter"] = []interface{}{packetFilterMap}
		}

		if ddosPolicy.DDoSRule.Acl != nil {
			aclMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.Acl.Switch != nil {
				aclMap["switch"] = ddosPolicy.DDoSRule.Acl.Switch
			}

			if ddosPolicy.DDoSRule.Acl.Acls != nil {
				aclsList := []interface{}{}
				for _, acls := range ddosPolicy.DDoSRule.Acl.Acls {
					aclsMap := map[string]interface{}{}

					if acls.DportEnd != nil {
						aclsMap["dport_end"] = acls.DportEnd
					}

					if acls.DportStart != nil {
						aclsMap["dport_start"] = acls.DportStart
					}

					if acls.SportEnd != nil {
						aclsMap["sport_end"] = acls.SportEnd
					}

					if acls.SportStart != nil {
						aclsMap["sport_start"] = acls.SportStart
					}

					if acls.Protocol != nil {
						aclsMap["protocol"] = acls.Protocol
					}

					if acls.Action != nil {
						aclsMap["action"] = acls.Action
					}

					aclsList = append(aclsList, aclsMap)
				}

				aclMap["acls"] = []interface{}{aclsList}
			}

			dDoSRuleMap["acl"] = []interface{}{aclMap}
		}

		if ddosPolicy.DDoSRule.SpeedLimit != nil {
			speedLimitMap := map[string]interface{}{}

			if ddosPolicy.DDoSRule.SpeedLimit.PackageLimit != nil {
				speedLimitMap["package_limit"] = ddosPolicy.DDoSRule.SpeedLimit.PackageLimit
			}

			if ddosPolicy.DDoSRule.SpeedLimit.FluxLimit != nil {
				speedLimitMap["flux_limit"] = ddosPolicy.DDoSRule.SpeedLimit.FluxLimit
			}

			dDoSRuleMap["speed_limit"] = []interface{}{speedLimitMap}
		}

		_ = d.Set("d_do_s_rule", []interface{}{dDoSRuleMap})
	}

	return nil
}

func resourceTencentCloudTeoDdosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyDDoSPolicyRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	policyId := idSplit[1]

	request.ZoneId = &zoneId
	request.PolicyId = &policyId

	immutableArgs := []string{"zone_id", "policy_id", "d_do_s_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("d_do_s_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "d_do_s_rule"); ok {
			ddosRule := teo.DdosRule{}
			if v, ok := dMap["switch"]; ok {
				ddosRule.Switch = helper.String(v.(string))
			}
			if statusInfoMap, ok := helper.InterfaceToMap(dMap, "status_info"); ok {
				dDoSStatusInfo := teo.DDoSStatusInfo{}
				if v, ok := statusInfoMap["ply_level"]; ok {
					dDoSStatusInfo.PlyLevel = helper.String(v.(string))
				}
				ddosRule.StatusInfo = &dDoSStatusInfo
			}
			if geoIpMap, ok := helper.InterfaceToMap(dMap, "geo_ip"); ok {
				dDoSGeoIp := teo.DDoSGeoIp{}
				if v, ok := geoIpMap["region_ids"]; ok {
					regionIdsSet := v.(*schema.Set).List()
					for i := range regionIdsSet {
						regionIds := regionIdsSet[i].(int)
						dDoSGeoIp.RegionIds = append(dDoSGeoIp.RegionIds, helper.IntInt64(regionIds))
					}
				}
				if v, ok := geoIpMap["switch"]; ok {
					dDoSGeoIp.Switch = helper.String(v.(string))
				}
				ddosRule.GeoIp = &dDoSGeoIp
			}
			if allowBlockMap, ok := helper.InterfaceToMap(dMap, "allow_block"); ok {
				dDoSAllowBlock := teo.DDoSAllowBlock{}
				if v, ok := allowBlockMap["switch"]; ok {
					dDoSAllowBlock.Switch = helper.String(v.(string))
				}
				if v, ok := allowBlockMap["allow_block_ips"]; ok {
					for _, item := range v.([]interface{}) {
						allowBlockIpsMap := item.(map[string]interface{})
						dDoSUserAllowBlockIP := teo.DDoSUserAllowBlockIP{}
						if v, ok := allowBlockIpsMap["ip"]; ok {
							dDoSUserAllowBlockIP.Ip = helper.String(v.(string))
						}
						if v, ok := allowBlockIpsMap["type"]; ok {
							dDoSUserAllowBlockIP.Type = helper.String(v.(string))
						}
						dDoSAllowBlock.AllowBlockIps = append(dDoSAllowBlock.AllowBlockIps, &dDoSUserAllowBlockIP)
					}
				}
				ddosRule.AllowBlock = &dDoSAllowBlock
			}
			if antiPlyMap, ok := helper.InterfaceToMap(dMap, "anti_ply"); ok {
				dDoSAntiPly := teo.DDoSAntiPly{}
				if v, ok := antiPlyMap["drop_tcp"]; ok {
					dDoSAntiPly.DropTcp = helper.String(v.(string))
				}
				if v, ok := antiPlyMap["drop_udp"]; ok {
					dDoSAntiPly.DropUdp = helper.String(v.(string))
				}
				if v, ok := antiPlyMap["drop_icmp"]; ok {
					dDoSAntiPly.DropIcmp = helper.String(v.(string))
				}
				if v, ok := antiPlyMap["drop_other"]; ok {
					dDoSAntiPly.DropOther = helper.String(v.(string))
				}
				if v, ok := antiPlyMap["source_create_limit"]; ok {
					dDoSAntiPly.SourceCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["source_connect_limit"]; ok {
					dDoSAntiPly.SourceConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["destination_create_limit"]; ok {
					dDoSAntiPly.DestinationCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["destination_connect_limit"]; ok {
					dDoSAntiPly.DestinationConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["abnormal_connect_num"]; ok {
					dDoSAntiPly.AbnormalConnectNum = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["abnormal_syn_ratio"]; ok {
					dDoSAntiPly.AbnormalSynRatio = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["abnormal_syn_num"]; ok {
					dDoSAntiPly.AbnormalSynNum = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["connect_timeout"]; ok {
					dDoSAntiPly.ConnectTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := antiPlyMap["empty_connect_protect"]; ok {
					dDoSAntiPly.EmptyConnectProtect = helper.String(v.(string))
				}
				if v, ok := antiPlyMap["udp_shard"]; ok {
					dDoSAntiPly.UdpShard = helper.String(v.(string))
				}
				ddosRule.AntiPly = &dDoSAntiPly
			}
			if packetFilterMap, ok := helper.InterfaceToMap(dMap, "packet_filter"); ok {
				dDoSPacketFilter := teo.DDoSPacketFilter{}
				if v, ok := packetFilterMap["switch"]; ok {
					dDoSPacketFilter.Switch = helper.String(v.(string))
				}
				if v, ok := packetFilterMap["packet_filters"]; ok {
					for _, item := range v.([]interface{}) {
						packetFiltersMap := item.(map[string]interface{})
						dDoSFeaturesFilter := teo.DDoSFeaturesFilter{}
						if v, ok := packetFiltersMap["action"]; ok {
							dDoSFeaturesFilter.Action = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["protocol"]; ok {
							dDoSFeaturesFilter.Protocol = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["dport_start"]; ok {
							dDoSFeaturesFilter.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["dport_end"]; ok {
							dDoSFeaturesFilter.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["packet_min"]; ok {
							dDoSFeaturesFilter.PacketMin = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["packet_max"]; ok {
							dDoSFeaturesFilter.PacketMax = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["sport_start"]; ok {
							dDoSFeaturesFilter.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["sport_end"]; ok {
							dDoSFeaturesFilter.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["match_type"]; ok {
							dDoSFeaturesFilter.MatchType = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["is_not"]; ok {
							dDoSFeaturesFilter.IsNot = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["offset"]; ok {
							dDoSFeaturesFilter.Offset = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["depth"]; ok {
							dDoSFeaturesFilter.Depth = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["match_begin"]; ok {
							dDoSFeaturesFilter.MatchBegin = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["str"]; ok {
							dDoSFeaturesFilter.Str = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["match_type2"]; ok {
							dDoSFeaturesFilter.MatchType2 = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["is_not2"]; ok {
							dDoSFeaturesFilter.IsNot2 = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["offset2"]; ok {
							dDoSFeaturesFilter.Offset2 = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["depth2"]; ok {
							dDoSFeaturesFilter.Depth2 = helper.IntInt64(v.(int))
						}
						if v, ok := packetFiltersMap["match_begin2"]; ok {
							dDoSFeaturesFilter.MatchBegin2 = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["str2"]; ok {
							dDoSFeaturesFilter.Str2 = helper.String(v.(string))
						}
						if v, ok := packetFiltersMap["match_logic"]; ok {
							dDoSFeaturesFilter.MatchLogic = helper.String(v.(string))
						}
						dDoSPacketFilter.PacketFilters = append(dDoSPacketFilter.PacketFilters, &dDoSFeaturesFilter)
					}
				}
				ddosRule.PacketFilter = &dDoSPacketFilter
			}
			if aclMap, ok := helper.InterfaceToMap(dMap, "acl"); ok {
				dDoSAcls := teo.DDoSAcls{}
				if v, ok := aclMap["switch"]; ok {
					dDoSAcls.Switch = helper.String(v.(string))
				}
				if v, ok := aclMap["acls"]; ok {
					for _, item := range v.([]interface{}) {
						aclsMap := item.(map[string]interface{})
						dDoSAcl := teo.DDoSAcl{}
						if v, ok := aclsMap["dport_end"]; ok {
							dDoSAcl.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := aclsMap["dport_start"]; ok {
							dDoSAcl.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := aclsMap["sport_end"]; ok {
							dDoSAcl.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := aclsMap["sport_start"]; ok {
							dDoSAcl.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := aclsMap["protocol"]; ok {
							dDoSAcl.Protocol = helper.String(v.(string))
						}
						if v, ok := aclsMap["action"]; ok {
							dDoSAcl.Action = helper.String(v.(string))
						}
						dDoSAcls.Acls = append(dDoSAcls.Acls, &dDoSAcl)
					}
				}
				ddosRule.Acl = &dDoSAcls
			}
			if speedLimitMap, ok := helper.InterfaceToMap(dMap, "speed_limit"); ok {
				dDoSSpeedLimit := teo.DDoSSpeedLimit{}
				if v, ok := speedLimitMap["package_limit"]; ok {
					dDoSSpeedLimit.PackageLimit = helper.String(v.(string))
				}
				if v, ok := speedLimitMap["flux_limit"]; ok {
					dDoSSpeedLimit.FluxLimit = helper.String(v.(string))
				}
				ddosRule.SpeedLimit = &dDoSSpeedLimit
			}
			request.DDoSRule = &ddosRule
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDDoSPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo ddosPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDdosPolicyRead(d, meta)
}

func resourceTencentCloudTeoDdosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
