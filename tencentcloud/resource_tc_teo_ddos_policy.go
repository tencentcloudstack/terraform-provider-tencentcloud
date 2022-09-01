/*
Provides a resource to create a teo ddosPolicy

Example Usage

```hcl
resource "tencentcloud_teo_ddos_policy" "ddosPolicy" {
  zone_id = ""
  policy_id = ""
  ddos_rule {
			switch = ""
			udp_shard_open = ""
		ddos_status_info {
				ply_level = ""
		}
		ddos_geo_ip {
				region_id = ""
				switch = ""
		}
		ddos_allow_block {
				switch = ""
			user_allow_block_ip {
					ip = ""
					mask = ""
					type = ""
					ip2 = ""
					mask2 = ""
			}
		}
		ddos_anti_ply {
				drop_tcp = ""
				drop_udp = ""
				drop_icmp = ""
				drop_other = ""
				source_create_limit = ""
				source_connect_limit = ""
				destination_create_limit = ""
				destination_connect_limit = ""
				abnormal_connect_num = ""
				abnormal_syn_ratio = ""
				abnormal_syn_num = ""
				connect_timeout = ""
				empty_connect_protect = ""
				udp_shard = ""
		}
		ddos_packet_filter {
				switch = ""
			packet_filter {
					action = ""
					protocol = ""
					dport_start = ""
					dport_end = ""
					packet_min = ""
					packet_max = ""
					sport_start = ""
					sport_end = ""
					match_type = ""
					is_not = ""
					offset = ""
					depth = ""
					match_begin = ""
					str = ""
					match_type2 = ""
					is_not2 = ""
					offset2 = ""
					depth2 = ""
					match_begin2 = ""
					str2 = ""
					match_logic = ""
			}
		}
		ddos_acl {
				switch = ""
			acl {
					dport_end = ""
					dport_start = ""
					sport_end = ""
					sport_start = ""
					protocol = ""
					action = ""
					default = ""
			}
		}

  }
}

```
Import

teo ddosPolicy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_ddos_policy.ddosPolicy ddosPolicy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoDdosPolicy() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoDdosPolicyRead,
		Create: resourceTencentCloudTeoDdosPolicyCreate,
		Update: resourceTencentCloudTeoDdosPolicyUpdate,
		Delete: resourceTencentCloudTeoDdosPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy ID.",
			},

			"ddos_rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "DDoS Configuration of the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "DDoS protection switch. Valid values:- on: Enable.- off: Disable.",
						},
						"udp_shard_open": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "UDP shard switch. Valid values:- on: Enable.- off: Disable.",
						},
						"ddos_status_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DDoS protection level.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ply_level": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Policy level. Valid values:- low: loose.- middle: moderate.- high: strict.",
									},
								},
							},
						},
						"ddos_geo_ip": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DDoS Protection by Geo Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
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
										Description: "- on: Enable.- off: Disable.",
									},
								},
							},
						},
						"ddos_allow_block": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DDoS black-white list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- on: Enable. `UserAllowBlockIp` parameter is required.- off: Disable.",
									},
									"user_allow_block_ip": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "DDoS black-white list detail.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Client IP.",
												},
												"mask": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "IP Mask.",
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
												"ip2": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "End of the IP address when setting an IP range.",
												},
												"mask2": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "IP mask of the end IP address.",
												},
											},
										},
									},
								},
							},
						},
						"ddos_anti_ply": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
										Description: "Limitation of new connection to origin website per second. Valid value range: 0-4294967295.",
									},
									"source_connect_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Limitation of connections to origin website. Valid value range: 0-4294967295.",
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
						"ddos_packet_filter": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DDoS feature filtering configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- on: Enable. `ModifyDDoSPolicy` parameter is required.- off: Disable.",
									},
									"packet_filter": {
										Type:        schema.TypeList,
										Optional:    true,
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
													Description: "Match type of feature 1. Valid values:- pcre: regex expression.- sunday: string match.",
												},
												"is_not": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Negate the match condition of feature 1. Valid values:- 0: match.- 1: not match.",
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
													Description: "Packet layer for matching begin of feature 1. Valid values:- begin_l5: matching from packet payload.- begin_l4: matching from TCP/UDP header.- begin_l3: matching from IP header.",
												},
												"str": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Regex expression or string to match.",
												},
												"match_type2": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Match type of feature 2. Valid values:- pcre: regex expression.- sunday: string match.",
												},
												"is_not2": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Negate the match condition of feature 2. Valid values:- 0: match.- 1: not match.",
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
													Description: "Packet layer for matching begin of feature 2. Valid values:- begin_l5: matching from packet payload.- begin_l4: matching from TCP/UDP header.- begin_l3: matching from IP header.",
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
						"ddos_acl": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DDoS ACL rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- on: Enable. `Acl` parameter is require.- off: Disable.",
									},
									"acl": {
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
												"default": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Whether it is default configuration. Valid value:- 0: custom configuration.- 1: default configuration.",
												},
											},
										},
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

	logId := getLogId(contextNil)

	var (
		request  = teo.NewModifyDDoSPolicyRequest()
		response *teo.ModifyDDoSPolicyResponse
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_id"); ok {
		request.PolicyId = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ddos_rule"); ok {
		ddosRule := teo.DdosRule{}
		if v, ok := dMap["switch"]; ok {
			ddosRule.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["udp_shard_open"]; ok {
			ddosRule.UdpShardOpen = helper.String(v.(string))
		}
		if DdosStatusInfoMap, ok := helper.InterfaceToMap(dMap, "ddos_status_info"); ok {
			dDoSStatusInfo := teo.DDoSStatusInfo{}
			if v, ok := DdosStatusInfoMap["ply_level"]; ok {
				dDoSStatusInfo.PlyLevel = helper.String(v.(string))
			}
			ddosRule.DdosStatusInfo = &dDoSStatusInfo
		}
		if DdosGeoIpMap, ok := helper.InterfaceToMap(dMap, "ddos_geo_ip"); ok {
			dDoSGeoIp := teo.DDoSGeoIp{}
			if v, ok := DdosGeoIpMap["region_id"]; ok {
				regionIdSet := v.(*schema.Set).List()
				for i := range regionIdSet {
					regionId := regionIdSet[i].(int)
					dDoSGeoIp.RegionId = append(dDoSGeoIp.RegionId, helper.IntInt64(regionId))
				}
			}
			if v, ok := DdosGeoIpMap["switch"]; ok {
				dDoSGeoIp.Switch = helper.String(v.(string))
			}
			ddosRule.DdosGeoIp = &dDoSGeoIp
		}
		if DdosAllowBlockMap, ok := helper.InterfaceToMap(dMap, "ddos_allow_block"); ok {
			ddosAllowBlock := teo.DdosAllowBlock{}
			if v, ok := DdosAllowBlockMap["switch"]; ok {
				ddosAllowBlock.Switch = helper.String(v.(string))
			}
			if v, ok := DdosAllowBlockMap["user_allow_block_ip"]; ok {
				for _, item := range v.([]interface{}) {
					UserAllowBlockIpMap := item.(map[string]interface{})
					dDoSUserAllowBlockIP := teo.DDoSUserAllowBlockIP{}
					if v, ok := UserAllowBlockIpMap["ip"]; ok {
						dDoSUserAllowBlockIP.Ip = helper.String(v.(string))
					}
					if v, ok := UserAllowBlockIpMap["mask"]; ok {
						dDoSUserAllowBlockIP.Mask = helper.IntInt64(v.(int))
					}
					if v, ok := UserAllowBlockIpMap["type"]; ok {
						dDoSUserAllowBlockIP.Type = helper.String(v.(string))
					}
					if v, ok := UserAllowBlockIpMap["ip2"]; ok {
						dDoSUserAllowBlockIP.Ip2 = helper.String(v.(string))
					}
					if v, ok := UserAllowBlockIpMap["mask2"]; ok {
						dDoSUserAllowBlockIP.Mask2 = helper.IntInt64(v.(int))
					}
					ddosAllowBlock.UserAllowBlockIp = append(ddosAllowBlock.UserAllowBlockIp, &dDoSUserAllowBlockIP)
				}
			}
			ddosRule.DdosAllowBlock = &ddosAllowBlock
		}
		if DdosAntiPlyMap, ok := helper.InterfaceToMap(dMap, "ddos_anti_ply"); ok {
			dDoSAntiPly := teo.DDoSAntiPly{}
			if v, ok := DdosAntiPlyMap["drop_tcp"]; ok {
				dDoSAntiPly.DropTcp = helper.String(v.(string))
			}
			if v, ok := DdosAntiPlyMap["drop_udp"]; ok {
				dDoSAntiPly.DropUdp = helper.String(v.(string))
			}
			if v, ok := DdosAntiPlyMap["drop_icmp"]; ok {
				dDoSAntiPly.DropIcmp = helper.String(v.(string))
			}
			if v, ok := DdosAntiPlyMap["drop_other"]; ok {
				dDoSAntiPly.DropOther = helper.String(v.(string))
			}
			if v, ok := DdosAntiPlyMap["source_create_limit"]; ok {
				dDoSAntiPly.SourceCreateLimit = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["source_connect_limit"]; ok {
				dDoSAntiPly.SourceConnectLimit = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["destination_create_limit"]; ok {
				dDoSAntiPly.DestinationCreateLimit = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["destination_connect_limit"]; ok {
				dDoSAntiPly.DestinationConnectLimit = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["abnormal_connect_num"]; ok {
				dDoSAntiPly.AbnormalConnectNum = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["abnormal_syn_ratio"]; ok {
				dDoSAntiPly.AbnormalSynRatio = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["abnormal_syn_num"]; ok {
				dDoSAntiPly.AbnormalSynNum = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["connect_timeout"]; ok {
				dDoSAntiPly.ConnectTimeout = helper.IntInt64(v.(int))
			}
			if v, ok := DdosAntiPlyMap["empty_connect_protect"]; ok {
				dDoSAntiPly.EmptyConnectProtect = helper.String(v.(string))
			}
			if v, ok := DdosAntiPlyMap["udp_shard"]; ok {
				dDoSAntiPly.UdpShard = helper.String(v.(string))
			}
			ddosRule.DdosAntiPly = &dDoSAntiPly
		}
		if DdosPacketFilterMap, ok := helper.InterfaceToMap(dMap, "ddos_packet_filter"); ok {
			ddosPacketFilter := teo.DdosPacketFilter{}
			if v, ok := DdosPacketFilterMap["switch"]; ok {
				ddosPacketFilter.Switch = helper.String(v.(string))
			}
			if v, ok := DdosPacketFilterMap["packet_filter"]; ok {
				for _, item := range v.([]interface{}) {
					PacketFilterMap := item.(map[string]interface{})
					dDoSFeaturesFilter := teo.DDoSFeaturesFilter{}
					if v, ok := PacketFilterMap["action"]; ok {
						dDoSFeaturesFilter.Action = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["protocol"]; ok {
						dDoSFeaturesFilter.Protocol = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["dport_start"]; ok {
						dDoSFeaturesFilter.DportStart = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["dport_end"]; ok {
						dDoSFeaturesFilter.DportEnd = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["packet_min"]; ok {
						dDoSFeaturesFilter.PacketMin = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["packet_max"]; ok {
						dDoSFeaturesFilter.PacketMax = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["sport_start"]; ok {
						dDoSFeaturesFilter.SportStart = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["sport_end"]; ok {
						dDoSFeaturesFilter.SportEnd = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["match_type"]; ok {
						dDoSFeaturesFilter.MatchType = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["is_not"]; ok {
						dDoSFeaturesFilter.IsNot = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["offset"]; ok {
						dDoSFeaturesFilter.Offset = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["depth"]; ok {
						dDoSFeaturesFilter.Depth = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["match_begin"]; ok {
						dDoSFeaturesFilter.MatchBegin = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["str"]; ok {
						dDoSFeaturesFilter.Str = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["match_type2"]; ok {
						dDoSFeaturesFilter.MatchType2 = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["is_not2"]; ok {
						dDoSFeaturesFilter.IsNot2 = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["offset2"]; ok {
						dDoSFeaturesFilter.Offset2 = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["depth2"]; ok {
						dDoSFeaturesFilter.Depth2 = helper.IntInt64(v.(int))
					}
					if v, ok := PacketFilterMap["match_begin2"]; ok {
						dDoSFeaturesFilter.MatchBegin2 = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["str2"]; ok {
						dDoSFeaturesFilter.Str2 = helper.String(v.(string))
					}
					if v, ok := PacketFilterMap["match_logic"]; ok {
						dDoSFeaturesFilter.MatchLogic = helper.String(v.(string))
					}
					ddosPacketFilter.PacketFilter = append(ddosPacketFilter.PacketFilter, &dDoSFeaturesFilter)
				}
			}
			ddosRule.DdosPacketFilter = &ddosPacketFilter
		}
		if DdosAclMap, ok := helper.InterfaceToMap(dMap, "ddos_acl"); ok {
			ddosAcls := teo.DdosAcls{}
			if v, ok := DdosAclMap["switch"]; ok {
				ddosAcls.Switch = helper.String(v.(string))
			}
			if v, ok := DdosAclMap["acl"]; ok {
				for _, item := range v.([]interface{}) {
					AclMap := item.(map[string]interface{})
					dDoSAcl := teo.DDoSAcl{}
					if v, ok := AclMap["dport_end"]; ok {
						dDoSAcl.DportEnd = helper.IntInt64(v.(int))
					}
					if v, ok := AclMap["dport_start"]; ok {
						dDoSAcl.DportStart = helper.IntInt64(v.(int))
					}
					if v, ok := AclMap["sport_end"]; ok {
						dDoSAcl.SportEnd = helper.IntInt64(v.(int))
					}
					if v, ok := AclMap["sport_start"]; ok {
						dDoSAcl.SportStart = helper.IntInt64(v.(int))
					}
					if v, ok := AclMap["protocol"]; ok {
						dDoSAcl.Protocol = helper.String(v.(string))
					}
					if v, ok := AclMap["action"]; ok {
						dDoSAcl.Action = helper.String(v.(string))
					}
					if v, ok := AclMap["default"]; ok {
						dDoSAcl.Default = helper.IntInt64(v.(int))
					}
					ddosAcls.Acl = append(ddosAcls.Acl, &dDoSAcl)
				}
			}
			ddosRule.DdosAcl = &ddosAcls
		}

		request.DdosRule = &ddosRule
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDDoSPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo ddosPolicy failed, reason:%+v", logId, err)
		return err
	}

	ddosPolicyId := strconv.FormatInt(*response.Response.PolicyId, 10)

	d.SetId(zoneId + FILED_SP + ddosPolicyId)
	return resourceTencentCloudTeoDdosPolicyRead(d, meta)
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

	ddosPolicy, err := service.DescribeTeoDdosPolicy(ctx, zoneId, policyId)

	if err != nil {
		return err
	}

	if ddosPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `ddosPolicy` %s does not exist", policyId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("policy_id", policyId)

	if ddosPolicy.DdosRule != nil {
		ddosRuleMap := map[string]interface{}{}
		if ddosPolicy.DdosRule.Switch != nil {
			ddosRuleMap["switch"] = ddosPolicy.DdosRule.Switch
		}
		if ddosPolicy.DdosRule.UdpShardOpen != nil {
			ddosRuleMap["udp_shard_open"] = ddosPolicy.DdosRule.UdpShardOpen
		}
		if ddosPolicy.DdosRule.DdosStatusInfo != nil {
			ddosStatusInfoMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosStatusInfo.PlyLevel != nil {
				ddosStatusInfoMap["ply_level"] = ddosPolicy.DdosRule.DdosStatusInfo.PlyLevel
			}

			ddosRuleMap["ddos_status_info"] = []interface{}{ddosStatusInfoMap}
		}
		if ddosPolicy.DdosRule.DdosGeoIp != nil {
			ddosGeoIpMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosGeoIp.RegionId != nil {
				ddosGeoIpMap["region_id"] = ddosPolicy.DdosRule.DdosGeoIp.RegionId
			}
			if ddosPolicy.DdosRule.DdosGeoIp.Switch != nil {
				ddosGeoIpMap["switch"] = ddosPolicy.DdosRule.DdosGeoIp.Switch
			}

			ddosRuleMap["ddos_geo_ip"] = []interface{}{ddosGeoIpMap}
		}
		if ddosPolicy.DdosRule.DdosAllowBlock != nil {
			ddosAllowBlockMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosAllowBlock.Switch != nil {
				ddosAllowBlockMap["switch"] = ddosPolicy.DdosRule.DdosAllowBlock.Switch
			}
			if ddosPolicy.DdosRule.DdosAllowBlock.UserAllowBlockIp != nil {
				userAllowBlockIpList := []interface{}{}
				for _, userAllowBlockIp := range ddosPolicy.DdosRule.DdosAllowBlock.UserAllowBlockIp {
					userAllowBlockIpMap := map[string]interface{}{}
					if userAllowBlockIp.Ip != nil {
						userAllowBlockIpMap["ip"] = userAllowBlockIp.Ip
					}
					if userAllowBlockIp.Mask != nil {
						userAllowBlockIpMap["mask"] = userAllowBlockIp.Mask
					}
					if userAllowBlockIp.Type != nil {
						userAllowBlockIpMap["type"] = userAllowBlockIp.Type
					}
					if userAllowBlockIp.UpdateTime != nil {
						userAllowBlockIpMap["update_time"] = userAllowBlockIp.UpdateTime
					}
					if userAllowBlockIp.Ip2 != nil {
						userAllowBlockIpMap["ip2"] = userAllowBlockIp.Ip2
					}
					if userAllowBlockIp.Mask2 != nil {
						userAllowBlockIpMap["mask2"] = userAllowBlockIp.Mask2
					}

					userAllowBlockIpList = append(userAllowBlockIpList, userAllowBlockIpMap)
				}
				ddosAllowBlockMap["user_allow_block_ip"] = userAllowBlockIpList
			}

			ddosRuleMap["ddos_allow_block"] = []interface{}{ddosAllowBlockMap}
		}
		if ddosPolicy.DdosRule.DdosAntiPly != nil {
			ddosAntiPlyMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosAntiPly.DropTcp != nil {
				ddosAntiPlyMap["drop_tcp"] = ddosPolicy.DdosRule.DdosAntiPly.DropTcp
			}
			if ddosPolicy.DdosRule.DdosAntiPly.DropUdp != nil {
				ddosAntiPlyMap["drop_udp"] = ddosPolicy.DdosRule.DdosAntiPly.DropUdp
			}
			if ddosPolicy.DdosRule.DdosAntiPly.DropIcmp != nil {
				ddosAntiPlyMap["drop_icmp"] = ddosPolicy.DdosRule.DdosAntiPly.DropIcmp
			}
			if ddosPolicy.DdosRule.DdosAntiPly.DropOther != nil {
				ddosAntiPlyMap["drop_other"] = ddosPolicy.DdosRule.DdosAntiPly.DropOther
			}
			if ddosPolicy.DdosRule.DdosAntiPly.SourceCreateLimit != nil {
				ddosAntiPlyMap["source_create_limit"] = ddosPolicy.DdosRule.DdosAntiPly.SourceCreateLimit
			}
			if ddosPolicy.DdosRule.DdosAntiPly.SourceConnectLimit != nil {
				ddosAntiPlyMap["source_connect_limit"] = ddosPolicy.DdosRule.DdosAntiPly.SourceConnectLimit
			}
			if ddosPolicy.DdosRule.DdosAntiPly.DestinationCreateLimit != nil {
				ddosAntiPlyMap["destination_create_limit"] = ddosPolicy.DdosRule.DdosAntiPly.DestinationCreateLimit
			}
			if ddosPolicy.DdosRule.DdosAntiPly.DestinationConnectLimit != nil {
				ddosAntiPlyMap["destination_connect_limit"] = ddosPolicy.DdosRule.DdosAntiPly.DestinationConnectLimit
			}
			if ddosPolicy.DdosRule.DdosAntiPly.AbnormalConnectNum != nil {
				ddosAntiPlyMap["abnormal_connect_num"] = ddosPolicy.DdosRule.DdosAntiPly.AbnormalConnectNum
			}
			if ddosPolicy.DdosRule.DdosAntiPly.AbnormalSynRatio != nil {
				ddosAntiPlyMap["abnormal_syn_ratio"] = ddosPolicy.DdosRule.DdosAntiPly.AbnormalSynRatio
			}
			if ddosPolicy.DdosRule.DdosAntiPly.AbnormalSynNum != nil {
				ddosAntiPlyMap["abnormal_syn_num"] = ddosPolicy.DdosRule.DdosAntiPly.AbnormalSynNum
			}
			if ddosPolicy.DdosRule.DdosAntiPly.ConnectTimeout != nil {
				ddosAntiPlyMap["connect_timeout"] = ddosPolicy.DdosRule.DdosAntiPly.ConnectTimeout
			}
			if ddosPolicy.DdosRule.DdosAntiPly.EmptyConnectProtect != nil {
				ddosAntiPlyMap["empty_connect_protect"] = ddosPolicy.DdosRule.DdosAntiPly.EmptyConnectProtect
			}
			if ddosPolicy.DdosRule.DdosAntiPly.UdpShard != nil {
				ddosAntiPlyMap["udp_shard"] = ddosPolicy.DdosRule.DdosAntiPly.UdpShard
			}

			ddosRuleMap["ddos_anti_ply"] = []interface{}{ddosAntiPlyMap}
		}
		if ddosPolicy.DdosRule.DdosPacketFilter != nil {
			ddosPacketFilterMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosPacketFilter.Switch != nil {
				ddosPacketFilterMap["switch"] = ddosPolicy.DdosRule.DdosPacketFilter.Switch
			}
			if ddosPolicy.DdosRule.DdosPacketFilter.PacketFilter != nil {
				packetFilterList := []interface{}{}
				for _, packetFilter := range ddosPolicy.DdosRule.DdosPacketFilter.PacketFilter {
					packetFilterMap := map[string]interface{}{}
					if packetFilter.Action != nil {
						packetFilterMap["action"] = packetFilter.Action
					}
					if packetFilter.Protocol != nil {
						packetFilterMap["protocol"] = packetFilter.Protocol
					}
					if packetFilter.DportStart != nil {
						packetFilterMap["dport_start"] = packetFilter.DportStart
					}
					if packetFilter.DportEnd != nil {
						packetFilterMap["dport_end"] = packetFilter.DportEnd
					}
					if packetFilter.PacketMin != nil {
						packetFilterMap["packet_min"] = packetFilter.PacketMin
					}
					if packetFilter.PacketMax != nil {
						packetFilterMap["packet_max"] = packetFilter.PacketMax
					}
					if packetFilter.SportStart != nil {
						packetFilterMap["sport_start"] = packetFilter.SportStart
					}
					if packetFilter.SportEnd != nil {
						packetFilterMap["sport_end"] = packetFilter.SportEnd
					}
					if packetFilter.MatchType != nil {
						packetFilterMap["match_type"] = packetFilter.MatchType
					}
					if packetFilter.IsNot != nil {
						packetFilterMap["is_not"] = packetFilter.IsNot
					}
					if packetFilter.Offset != nil {
						packetFilterMap["offset"] = packetFilter.Offset
					}
					if packetFilter.Depth != nil {
						packetFilterMap["depth"] = packetFilter.Depth
					}
					if packetFilter.MatchBegin != nil {
						packetFilterMap["match_begin"] = packetFilter.MatchBegin
					}
					if packetFilter.Str != nil {
						packetFilterMap["str"] = packetFilter.Str
					}
					if packetFilter.MatchType2 != nil {
						packetFilterMap["match_type2"] = packetFilter.MatchType2
					}
					if packetFilter.IsNot2 != nil {
						packetFilterMap["is_not2"] = packetFilter.IsNot2
					}
					if packetFilter.Offset2 != nil {
						packetFilterMap["offset2"] = packetFilter.Offset2
					}
					if packetFilter.Depth2 != nil {
						packetFilterMap["depth2"] = packetFilter.Depth2
					}
					if packetFilter.MatchBegin2 != nil {
						packetFilterMap["match_begin2"] = packetFilter.MatchBegin2
					}
					if packetFilter.Str2 != nil {
						packetFilterMap["str2"] = packetFilter.Str2
					}
					if packetFilter.MatchLogic != nil {
						packetFilterMap["match_logic"] = packetFilter.MatchLogic
					}

					packetFilterList = append(packetFilterList, packetFilterMap)
				}
				ddosPacketFilterMap["packet_filter"] = packetFilterList
			}

			ddosRuleMap["ddos_packet_filter"] = []interface{}{ddosPacketFilterMap}
		}
		if ddosPolicy.DdosRule.DdosAcl != nil {
			ddosAclMap := map[string]interface{}{}
			if ddosPolicy.DdosRule.DdosAcl.Switch != nil {
				ddosAclMap["switch"] = ddosPolicy.DdosRule.DdosAcl.Switch
			}
			if ddosPolicy.DdosRule.DdosAcl.Acl != nil {
				aclList := []interface{}{}
				for _, acl := range ddosPolicy.DdosRule.DdosAcl.Acl {
					aclMap := map[string]interface{}{}
					if acl.DportEnd != nil {
						aclMap["dport_end"] = acl.DportEnd
					}
					if acl.DportStart != nil {
						aclMap["dport_start"] = acl.DportStart
					}
					if acl.SportEnd != nil {
						aclMap["sport_end"] = acl.SportEnd
					}
					if acl.SportStart != nil {
						aclMap["sport_start"] = acl.SportStart
					}
					if acl.Protocol != nil {
						aclMap["protocol"] = acl.Protocol
					}
					if acl.Action != nil {
						aclMap["action"] = acl.Action
					}
					if acl.Default != nil {
						aclMap["default"] = acl.Default
					}

					aclList = append(aclList, aclMap)
				}
				ddosAclMap["acl"] = aclList
			}

			ddosRuleMap["ddos_acl"] = []interface{}{ddosAclMap}
		}

		_ = d.Set("ddos_rule", []interface{}{ddosRuleMap})
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

	policyId64, errRet := strconv.ParseInt(policyId, 10, 64)
	if errRet != nil {
		return fmt.Errorf("Type conversion failed, [%s] conversion int64 failed\n", policyId)
	}

	request.ZoneId = &zoneId
	request.PolicyId = &policyId64

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("policy_id") {
		return fmt.Errorf("`policy_id` do not support change now.")
	}

	if d.HasChange("ddos_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ddos_rule"); ok {
			ddosRule := teo.DdosRule{}
			if v, ok := dMap["switch"]; ok {
				ddosRule.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["udp_shard_open"]; ok {
				ddosRule.UdpShardOpen = helper.String(v.(string))
			}
			if DdosStatusInfoMap, ok := helper.InterfaceToMap(dMap, "ddos_status_info"); ok {
				dDoSStatusInfo := teo.DDoSStatusInfo{}
				if v, ok := DdosStatusInfoMap["ply_level"]; ok {
					dDoSStatusInfo.PlyLevel = helper.String(v.(string))
				}
				ddosRule.DdosStatusInfo = &dDoSStatusInfo
			}
			if DdosGeoIpMap, ok := helper.InterfaceToMap(dMap, "ddos_geo_ip"); ok {
				dDoSGeoIp := teo.DDoSGeoIp{}
				if v, ok := DdosGeoIpMap["region_id"]; ok {
					regionIdSet := v.(*schema.Set).List()
					for i := range regionIdSet {
						regionId := regionIdSet[i].(int)
						dDoSGeoIp.RegionId = append(dDoSGeoIp.RegionId, helper.IntInt64(regionId))
					}
				}
				if v, ok := DdosGeoIpMap["switch"]; ok {
					dDoSGeoIp.Switch = helper.String(v.(string))
				}
				ddosRule.DdosGeoIp = &dDoSGeoIp
			}
			if DdosAllowBlockMap, ok := helper.InterfaceToMap(dMap, "ddos_allow_block"); ok {
				ddosAllowBlock := teo.DdosAllowBlock{}
				if v, ok := DdosAllowBlockMap["switch"]; ok {
					ddosAllowBlock.Switch = helper.String(v.(string))
				}
				if v, ok := DdosAllowBlockMap["user_allow_block_ip"]; ok {
					for _, item := range v.([]interface{}) {
						UserAllowBlockIpMap := item.(map[string]interface{})
						dDoSUserAllowBlockIP := teo.DDoSUserAllowBlockIP{}
						if v, ok := UserAllowBlockIpMap["ip"]; ok {
							dDoSUserAllowBlockIP.Ip = helper.String(v.(string))
						}
						if v, ok := UserAllowBlockIpMap["mask"]; ok {
							dDoSUserAllowBlockIP.Mask = helper.IntInt64(v.(int))
						}
						if v, ok := UserAllowBlockIpMap["type"]; ok {
							dDoSUserAllowBlockIP.Type = helper.String(v.(string))
						}
						if v, ok := UserAllowBlockIpMap["ip2"]; ok {
							dDoSUserAllowBlockIP.Ip2 = helper.String(v.(string))
						}
						if v, ok := UserAllowBlockIpMap["mask2"]; ok {
							dDoSUserAllowBlockIP.Mask2 = helper.IntInt64(v.(int))
						}
						ddosAllowBlock.UserAllowBlockIp = append(ddosAllowBlock.UserAllowBlockIp, &dDoSUserAllowBlockIP)
					}
				}
				ddosRule.DdosAllowBlock = &ddosAllowBlock
			}
			if DdosAntiPlyMap, ok := helper.InterfaceToMap(dMap, "ddos_anti_ply"); ok {
				dDoSAntiPly := teo.DDoSAntiPly{}
				if v, ok := DdosAntiPlyMap["drop_tcp"]; ok {
					dDoSAntiPly.DropTcp = helper.String(v.(string))
				}
				if v, ok := DdosAntiPlyMap["drop_udp"]; ok {
					dDoSAntiPly.DropUdp = helper.String(v.(string))
				}
				if v, ok := DdosAntiPlyMap["drop_icmp"]; ok {
					dDoSAntiPly.DropIcmp = helper.String(v.(string))
				}
				if v, ok := DdosAntiPlyMap["drop_other"]; ok {
					dDoSAntiPly.DropOther = helper.String(v.(string))
				}
				if v, ok := DdosAntiPlyMap["source_create_limit"]; ok {
					dDoSAntiPly.SourceCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["source_connect_limit"]; ok {
					dDoSAntiPly.SourceConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["destination_create_limit"]; ok {
					dDoSAntiPly.DestinationCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["destination_connect_limit"]; ok {
					dDoSAntiPly.DestinationConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["abnormal_connect_num"]; ok {
					dDoSAntiPly.AbnormalConnectNum = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["abnormal_syn_ratio"]; ok {
					dDoSAntiPly.AbnormalSynRatio = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["abnormal_syn_num"]; ok {
					dDoSAntiPly.AbnormalSynNum = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["connect_timeout"]; ok {
					dDoSAntiPly.ConnectTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := DdosAntiPlyMap["empty_connect_protect"]; ok {
					dDoSAntiPly.EmptyConnectProtect = helper.String(v.(string))
				}
				if v, ok := DdosAntiPlyMap["udp_shard"]; ok {
					dDoSAntiPly.UdpShard = helper.String(v.(string))
				}
				ddosRule.DdosAntiPly = &dDoSAntiPly
			}
			if DdosPacketFilterMap, ok := helper.InterfaceToMap(dMap, "ddos_packet_filter"); ok {
				ddosPacketFilter := teo.DdosPacketFilter{}
				if v, ok := DdosPacketFilterMap["switch"]; ok {
					ddosPacketFilter.Switch = helper.String(v.(string))
				}
				if v, ok := DdosPacketFilterMap["packet_filter"]; ok {
					for _, item := range v.([]interface{}) {
						PacketFilterMap := item.(map[string]interface{})
						dDoSFeaturesFilter := teo.DDoSFeaturesFilter{}
						if v, ok := PacketFilterMap["action"]; ok {
							dDoSFeaturesFilter.Action = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["protocol"]; ok {
							dDoSFeaturesFilter.Protocol = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["dport_start"]; ok {
							dDoSFeaturesFilter.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["dport_end"]; ok {
							dDoSFeaturesFilter.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["packet_min"]; ok {
							dDoSFeaturesFilter.PacketMin = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["packet_max"]; ok {
							dDoSFeaturesFilter.PacketMax = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["sport_start"]; ok {
							dDoSFeaturesFilter.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["sport_end"]; ok {
							dDoSFeaturesFilter.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["match_type"]; ok {
							dDoSFeaturesFilter.MatchType = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["is_not"]; ok {
							dDoSFeaturesFilter.IsNot = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["offset"]; ok {
							dDoSFeaturesFilter.Offset = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["depth"]; ok {
							dDoSFeaturesFilter.Depth = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["match_begin"]; ok {
							dDoSFeaturesFilter.MatchBegin = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["str"]; ok {
							dDoSFeaturesFilter.Str = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["match_type2"]; ok {
							dDoSFeaturesFilter.MatchType2 = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["is_not2"]; ok {
							dDoSFeaturesFilter.IsNot2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["offset2"]; ok {
							dDoSFeaturesFilter.Offset2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["depth2"]; ok {
							dDoSFeaturesFilter.Depth2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFilterMap["match_begin2"]; ok {
							dDoSFeaturesFilter.MatchBegin2 = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["str2"]; ok {
							dDoSFeaturesFilter.Str2 = helper.String(v.(string))
						}
						if v, ok := PacketFilterMap["match_logic"]; ok {
							dDoSFeaturesFilter.MatchLogic = helper.String(v.(string))
						}
						ddosPacketFilter.PacketFilter = append(ddosPacketFilter.PacketFilter, &dDoSFeaturesFilter)
					}
				}
				ddosRule.DdosPacketFilter = &ddosPacketFilter
			}
			if DdosAclMap, ok := helper.InterfaceToMap(dMap, "ddos_acl"); ok {
				ddosAcls := teo.DdosAcls{}
				if v, ok := DdosAclMap["switch"]; ok {
					ddosAcls.Switch = helper.String(v.(string))
				}
				if v, ok := DdosAclMap["acl"]; ok {
					for _, item := range v.([]interface{}) {
						AclMap := item.(map[string]interface{})
						dDoSAcl := teo.DDoSAcl{}
						if v, ok := AclMap["dport_end"]; ok {
							dDoSAcl.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := AclMap["dport_start"]; ok {
							dDoSAcl.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := AclMap["sport_end"]; ok {
							dDoSAcl.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := AclMap["sport_start"]; ok {
							dDoSAcl.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := AclMap["protocol"]; ok {
							dDoSAcl.Protocol = helper.String(v.(string))
						}
						if v, ok := AclMap["action"]; ok {
							dDoSAcl.Action = helper.String(v.(string))
						}
						if v, ok := AclMap["default"]; ok {
							dDoSAcl.Default = helper.IntInt64(v.(int))
						}
						ddosAcls.Acl = append(ddosAcls.Acl, &dDoSAcl)
					}
				}
				ddosRule.DdosAcl = &ddosAcls
			}

			request.DdosRule = &ddosRule
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDDoSPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTeoDdosPolicyRead(d, meta)
}

func resourceTencentCloudTeoDdosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.delete")()
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

	if err := service.DeleteTeoDdosPolicyById(ctx, zoneId, policyId); err != nil {
		return err
	}

	return nil
}
