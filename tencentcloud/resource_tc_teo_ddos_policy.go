/*
Provides a resource to create a teo ddos_policy

Example Usage

```hcl
resource "tencentcloud_teo_ddos_policy" "ddos_policy" {
  policy_id = 1278
  zone_id   = "zone-2983wizgxqvm"

  ddos_rule {
    switch = "on"

    acl {
      switch = "on"
    }

    allow_block {
      switch = "on"
    }

    anti_ply {
      abnormal_connect_num      = 0
      abnormal_syn_num          = 0
      abnormal_syn_ratio        = 0
      connect_timeout           = 0
      destination_connect_limit = 0
      destination_create_limit  = 0
      drop_icmp                 = "off"
      drop_other                = "off"
      drop_tcp                  = "off"
      drop_udp                  = "off"
      empty_connect_protect     = "off"
      source_connect_limit      = 0
      source_create_limit       = 0
      udp_shard                 = "off"
    }

    geo_ip {
      region_ids = []
      switch     = "on"
    }

    packet_filter {
      switch = "on"
    }

    speed_limit {
      flux_limit    = "0 bps"
      package_limit = "0 pps"
    }

    status_info {
      ply_level = "middle"
    }
  }
}


```
Import

teo ddos_policy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_ddos_policy.ddos_policy ddosPolicy_id
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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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
										Description: "Limit the number of packages. Valid range: 1 pps-10000 Gpps, 0 means no limitation, supported units: `pps`,`Kpps`,`Mpps`,`Gpps`.",
									},
									"flux_limit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Limit the number of fluxes. Valid range: 1 bps-10000 Gbps, 0 means no limitation, supported units: `pps`,`Kpps`,`Mpps`,`Gpps`.",
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
		zoneId   string
		policyId int64
		service  = TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("policy_id"); ok {
		if policyId, ok = v.(int64); !ok {
			var tmpPolicyId int
			if tmpPolicyId, ok = v.(int); !ok {
				return fmt.Errorf("create teo ddosPolicy failed, reason: invalid policyId %+v", v)
			}
			policyId = int64(tmpPolicyId)
		}
	}

	var (
		policyIdChecked bool
		ddosPolicy      *teo.DescribeZoneDDoSPolicyResponseParams
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := service.DescribeTeoZoneDDoSPolicyByFilter(ctx, map[string]interface{}{
			"zone_id": zoneId,
		})
		if e != nil {
			return retryError(e)
		}
		ddosPolicy = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo planInfo failed, reason:%+v", logId, err)
		return err
	}
	for _, areas := range ddosPolicy.ShieldAreas {
		if *areas.PolicyId == policyId {
			policyIdChecked = true
		}
	}
	if !policyIdChecked {
		return fmt.Errorf("create teo ddosPolicy failed, reason: invalid policy id %v", policyId)
	}

	d.SetId(zoneId + FILED_SP + strconv.Itoa(int(policyId)))
	err = resourceTencentCloudTeoDdosPolicyUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo ddosPolicy failed, reason:%+v", logId, err)
		return err
	}
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

	policyId64, err := strconv.ParseInt(policyId, 10, 64)
	if err != nil {
		log.Printf("[READ]%s read teo ddosPolicy parseInt[%v] failed, reason:%+v", logId, policyId, err)
		return err
	}
	ddosPolicy, err := service.DescribeTeoDdosPolicy(ctx, zoneId, policyId64)

	if err != nil {
		return err
	}

	if ddosPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `ddosPolicy` %s does not exist", policyId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("policy_id", policyId64)

	if ddosPolicy.DDoSRule != nil {
		dDoSRuleMap := map[string]interface{}{}
		if ddosPolicy.DDoSRule.Switch != nil {
			dDoSRuleMap["switch"] = ddosPolicy.DDoSRule.Switch
		}
		if ddosPolicy.DDoSRule.DDoSStatusInfo != nil {
			statusInfoMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSStatusInfo.PlyLevel != nil {
				statusInfoMap["ply_level"] = ddosPolicy.DDoSRule.DDoSStatusInfo.PlyLevel
			}

			dDoSRuleMap["status_info"] = []interface{}{statusInfoMap}
		}
		if ddosPolicy.DDoSRule.DDoSGeoIp != nil {
			geoIpMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSGeoIp.RegionIds != nil {
				geoIpMap["region_ids"] = ddosPolicy.DDoSRule.DDoSGeoIp.RegionIds
			}
			if ddosPolicy.DDoSRule.DDoSGeoIp.Switch != nil {
				geoIpMap["switch"] = ddosPolicy.DDoSRule.DDoSGeoIp.Switch
			}

			dDoSRuleMap["geo_ip"] = []interface{}{geoIpMap}
		}
		if ddosPolicy.DDoSRule.DDoSAllowBlock != nil {
			allowBlockMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSAllowBlock.Switch != nil {
				allowBlockMap["switch"] = ddosPolicy.DDoSRule.DDoSAllowBlock.Switch
			}
			if ddosPolicy.DDoSRule.DDoSAllowBlock.DDoSAllowBlockRules != nil {
				allowBlockIpsList := []interface{}{}
				for _, allowBlockIps := range ddosPolicy.DDoSRule.DDoSAllowBlock.DDoSAllowBlockRules {
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
				allowBlockMap["allow_block_ips"] = allowBlockIpsList
			}

			dDoSRuleMap["allow_block"] = []interface{}{allowBlockMap}
		}
		if ddosPolicy.DDoSRule.DDoSAntiPly != nil {
			antiPlyMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DropTcp != nil {
				antiPlyMap["drop_tcp"] = ddosPolicy.DDoSRule.DDoSAntiPly.DropTcp
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DropUdp != nil {
				antiPlyMap["drop_udp"] = ddosPolicy.DDoSRule.DDoSAntiPly.DropUdp
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DropIcmp != nil {
				antiPlyMap["drop_icmp"] = ddosPolicy.DDoSRule.DDoSAntiPly.DropIcmp
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DropOther != nil {
				antiPlyMap["drop_other"] = ddosPolicy.DDoSRule.DDoSAntiPly.DropOther
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.SourceCreateLimit != nil {
				antiPlyMap["source_create_limit"] = ddosPolicy.DDoSRule.DDoSAntiPly.SourceCreateLimit
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.SourceConnectLimit != nil {
				antiPlyMap["source_connect_limit"] = ddosPolicy.DDoSRule.DDoSAntiPly.SourceConnectLimit
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DestinationCreateLimit != nil {
				antiPlyMap["destination_create_limit"] = ddosPolicy.DDoSRule.DDoSAntiPly.DestinationCreateLimit
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.DestinationConnectLimit != nil {
				antiPlyMap["destination_connect_limit"] = ddosPolicy.DDoSRule.DDoSAntiPly.DestinationConnectLimit
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalConnectNum != nil {
				antiPlyMap["abnormal_connect_num"] = ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalConnectNum
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalSynRatio != nil {
				antiPlyMap["abnormal_syn_ratio"] = ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalSynRatio
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalSynNum != nil {
				antiPlyMap["abnormal_syn_num"] = ddosPolicy.DDoSRule.DDoSAntiPly.AbnormalSynNum
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.ConnectTimeout != nil {
				antiPlyMap["connect_timeout"] = ddosPolicy.DDoSRule.DDoSAntiPly.ConnectTimeout
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.EmptyConnectProtect != nil {
				antiPlyMap["empty_connect_protect"] = ddosPolicy.DDoSRule.DDoSAntiPly.EmptyConnectProtect
			}
			if ddosPolicy.DDoSRule.DDoSAntiPly.UdpShard != nil {
				antiPlyMap["udp_shard"] = ddosPolicy.DDoSRule.DDoSAntiPly.UdpShard
			}

			dDoSRuleMap["anti_ply"] = []interface{}{antiPlyMap}
		}
		if ddosPolicy.DDoSRule.DDoSPacketFilter != nil {
			packetFilterMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSPacketFilter.Switch != nil {
				packetFilterMap["switch"] = ddosPolicy.DDoSRule.DDoSPacketFilter.Switch
			}
			if ddosPolicy.DDoSRule.DDoSPacketFilter.DDoSFeaturesFilters != nil {
				packetFiltersList := []interface{}{}
				for _, packetFilters := range ddosPolicy.DDoSRule.DDoSPacketFilter.DDoSFeaturesFilters {
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
				packetFilterMap["packet_filters"] = packetFiltersList
			}

			dDoSRuleMap["packet_filter"] = []interface{}{packetFilterMap}
		}
		if ddosPolicy.DDoSRule.DDoSAcl != nil {
			aclMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSAcl.Switch != nil {
				aclMap["switch"] = ddosPolicy.DDoSRule.DDoSAcl.Switch
			}
			if ddosPolicy.DDoSRule.DDoSAcl.DDoSAclRules != nil {
				aclsList := []interface{}{}
				for _, acls := range ddosPolicy.DDoSRule.DDoSAcl.DDoSAclRules {
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
				aclMap["acls"] = aclsList
			}

			dDoSRuleMap["acl"] = []interface{}{aclMap}
		}
		if ddosPolicy.DDoSRule.DDoSSpeedLimit != nil {
			speedLimitMap := map[string]interface{}{}
			if ddosPolicy.DDoSRule.DDoSSpeedLimit.PackageLimit != nil {
				speedLimitMap["package_limit"] = ddosPolicy.DDoSRule.DDoSSpeedLimit.PackageLimit
			}
			if ddosPolicy.DDoSRule.DDoSSpeedLimit.FluxLimit != nil {
				speedLimitMap["flux_limit"] = ddosPolicy.DDoSRule.DDoSSpeedLimit.FluxLimit
			}

			dDoSRuleMap["speed_limit"] = []interface{}{speedLimitMap}
		}

		_ = d.Set("ddos_rule", []interface{}{dDoSRuleMap})
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

	policyId64, err := strconv.ParseInt(policyId, 10, 64)
	if err != nil {
		log.Printf("[UPDATE]%s update teo ddosPolicy parseInt[%v] failed, reason:%+v", logId, policyId, err)
		return err
	}
	request.ZoneId = &zoneId
	request.PolicyId = helper.Int64(policyId64)

	if d.HasChange("zone_id") {
		if old, _ := d.GetChange("zone_id"); old.(string) != "" {
			return fmt.Errorf("`zone_id` do not support change now.")
		}
	}

	if d.HasChange("policy_id") {
		if old, _ := d.GetChange("policy_id"); old.(int) != 0 {
			return fmt.Errorf("`policy_id` do not support change now.")
		}
	}

	if d.HasChange("ddos_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ddos_rule"); ok {
			ddosRule := teo.DDoSRule{}
			if v, ok := dMap["switch"]; ok {
				ddosRule.Switch = helper.String(v.(string))
			}
			if StatusInfoMap, ok := helper.InterfaceToMap(dMap, "status_info"); ok {
				dDoSStatusInfo := teo.DDoSStatusInfo{}
				if v, ok := StatusInfoMap["ply_level"]; ok {
					dDoSStatusInfo.PlyLevel = helper.String(v.(string))
				}
				ddosRule.DDoSStatusInfo = &dDoSStatusInfo
			}
			if GeoIpMap, ok := helper.InterfaceToMap(dMap, "geo_ip"); ok {
				dDoSGeoIp := teo.DDoSGeoIp{}
				if v, ok := GeoIpMap["region_ids"]; ok {
					regionIdsSet := v.(*schema.Set).List()
					for i := range regionIdsSet {
						regionIds := regionIdsSet[i].(int)
						dDoSGeoIp.RegionIds = append(dDoSGeoIp.RegionIds, helper.IntInt64(regionIds))
					}
				}
				if v, ok := GeoIpMap["switch"]; ok {
					dDoSGeoIp.Switch = helper.String(v.(string))
				}
				ddosRule.DDoSGeoIp = &dDoSGeoIp
			}
			if AllowBlockMap, ok := helper.InterfaceToMap(dMap, "allow_block"); ok {
				dDoSAllowBlock := teo.DDoSAllowBlock{}
				if v, ok := AllowBlockMap["switch"]; ok {
					dDoSAllowBlock.Switch = helper.String(v.(string))
				}
				if v, ok := AllowBlockMap["allow_block_ips"]; ok {
					for _, item := range v.([]interface{}) {
						AllowBlockIpsMap := item.(map[string]interface{})
						dDoSUserAllowBlockIP := teo.DDoSAllowBlockRule{}
						if v, ok := AllowBlockIpsMap["ip"]; ok {
							dDoSUserAllowBlockIP.Ip = helper.String(v.(string))
						}
						if v, ok := AllowBlockIpsMap["type"]; ok {
							dDoSUserAllowBlockIP.Type = helper.String(v.(string))
						}
						dDoSAllowBlock.DDoSAllowBlockRules = append(dDoSAllowBlock.DDoSAllowBlockRules, &dDoSUserAllowBlockIP)
					}
				}
				ddosRule.DDoSAllowBlock = &dDoSAllowBlock
			}
			if AntiPlyMap, ok := helper.InterfaceToMap(dMap, "anti_ply"); ok {
				dDoSAntiPly := teo.DDoSAntiPly{}
				if v, ok := AntiPlyMap["drop_tcp"]; ok {
					dDoSAntiPly.DropTcp = helper.String(v.(string))
				}
				if v, ok := AntiPlyMap["drop_udp"]; ok {
					dDoSAntiPly.DropUdp = helper.String(v.(string))
				}
				if v, ok := AntiPlyMap["drop_icmp"]; ok {
					dDoSAntiPly.DropIcmp = helper.String(v.(string))
				}
				if v, ok := AntiPlyMap["drop_other"]; ok {
					dDoSAntiPly.DropOther = helper.String(v.(string))
				}
				if v, ok := AntiPlyMap["source_create_limit"]; ok {
					dDoSAntiPly.SourceCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["source_connect_limit"]; ok {
					dDoSAntiPly.SourceConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["destination_create_limit"]; ok {
					dDoSAntiPly.DestinationCreateLimit = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["destination_connect_limit"]; ok {
					dDoSAntiPly.DestinationConnectLimit = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["abnormal_connect_num"]; ok {
					dDoSAntiPly.AbnormalConnectNum = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["abnormal_syn_ratio"]; ok {
					dDoSAntiPly.AbnormalSynRatio = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["abnormal_syn_num"]; ok {
					dDoSAntiPly.AbnormalSynNum = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["connect_timeout"]; ok {
					dDoSAntiPly.ConnectTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := AntiPlyMap["empty_connect_protect"]; ok {
					dDoSAntiPly.EmptyConnectProtect = helper.String(v.(string))
				}
				if v, ok := AntiPlyMap["udp_shard"]; ok {
					dDoSAntiPly.UdpShard = helper.String(v.(string))
				}
				ddosRule.DDoSAntiPly = &dDoSAntiPly
			}
			if PacketFilterMap, ok := helper.InterfaceToMap(dMap, "packet_filter"); ok {
				dDoSPacketFilter := teo.DDoSPacketFilter{}
				if v, ok := PacketFilterMap["switch"]; ok {
					dDoSPacketFilter.Switch = helper.String(v.(string))
				}
				if v, ok := PacketFilterMap["packet_filters"]; ok {
					for _, item := range v.([]interface{}) {
						PacketFiltersMap := item.(map[string]interface{})
						dDoSFeaturesFilter := teo.DDoSFeaturesFilter{}
						if v, ok := PacketFiltersMap["action"]; ok {
							dDoSFeaturesFilter.Action = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["protocol"]; ok {
							dDoSFeaturesFilter.Protocol = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["dport_start"]; ok {
							dDoSFeaturesFilter.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["dport_end"]; ok {
							dDoSFeaturesFilter.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["packet_min"]; ok {
							dDoSFeaturesFilter.PacketMin = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["packet_max"]; ok {
							dDoSFeaturesFilter.PacketMax = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["sport_start"]; ok {
							dDoSFeaturesFilter.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["sport_end"]; ok {
							dDoSFeaturesFilter.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["match_type"]; ok {
							dDoSFeaturesFilter.MatchType = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["is_not"]; ok {
							dDoSFeaturesFilter.IsNot = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["offset"]; ok {
							dDoSFeaturesFilter.Offset = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["depth"]; ok {
							dDoSFeaturesFilter.Depth = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["match_begin"]; ok {
							dDoSFeaturesFilter.MatchBegin = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["str"]; ok {
							dDoSFeaturesFilter.Str = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["match_type2"]; ok {
							dDoSFeaturesFilter.MatchType2 = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["is_not2"]; ok {
							dDoSFeaturesFilter.IsNot2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["offset2"]; ok {
							dDoSFeaturesFilter.Offset2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["depth2"]; ok {
							dDoSFeaturesFilter.Depth2 = helper.IntInt64(v.(int))
						}
						if v, ok := PacketFiltersMap["match_begin2"]; ok {
							dDoSFeaturesFilter.MatchBegin2 = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["str2"]; ok {
							dDoSFeaturesFilter.Str2 = helper.String(v.(string))
						}
						if v, ok := PacketFiltersMap["match_logic"]; ok {
							dDoSFeaturesFilter.MatchLogic = helper.String(v.(string))
						}
						dDoSPacketFilter.DDoSFeaturesFilters = append(dDoSPacketFilter.DDoSFeaturesFilters, &dDoSFeaturesFilter)
					}
				}
				ddosRule.DDoSPacketFilter = &dDoSPacketFilter
			}
			if AclMap, ok := helper.InterfaceToMap(dMap, "acl"); ok {
				dDoSAcls := teo.DDoSAcl{}
				if v, ok := AclMap["switch"]; ok {
					dDoSAcls.Switch = helper.String(v.(string))
				}
				if v, ok := AclMap["acls"]; ok {
					for _, item := range v.([]interface{}) {
						AclsMap := item.(map[string]interface{})
						dDoSAcl := teo.DDoSAclRule{}
						if v, ok := AclsMap["dport_end"]; ok {
							dDoSAcl.DportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := AclsMap["dport_start"]; ok {
							dDoSAcl.DportStart = helper.IntInt64(v.(int))
						}
						if v, ok := AclsMap["sport_end"]; ok {
							dDoSAcl.SportEnd = helper.IntInt64(v.(int))
						}
						if v, ok := AclsMap["sport_start"]; ok {
							dDoSAcl.SportStart = helper.IntInt64(v.(int))
						}
						if v, ok := AclsMap["protocol"]; ok {
							dDoSAcl.Protocol = helper.String(v.(string))
						}
						if v, ok := AclsMap["action"]; ok {
							dDoSAcl.Action = helper.String(v.(string))
						}
						dDoSAcls.DDoSAclRules = append(dDoSAcls.DDoSAclRules, &dDoSAcl)
					}
				}
				ddosRule.DDoSAcl = &dDoSAcls
			}
			if SpeedLimitMap, ok := helper.InterfaceToMap(dMap, "speed_limit"); ok {
				dDoSSpeedLimit := teo.DDoSSpeedLimit{}
				if v, ok := SpeedLimitMap["package_limit"]; ok {
					dDoSSpeedLimit.PackageLimit = helper.String(v.(string))
				}
				if v, ok := SpeedLimitMap["flux_limit"]; ok {
					dDoSSpeedLimit.FluxLimit = helper.String(v.(string))
				}
				ddosRule.DDoSSpeedLimit = &dDoSSpeedLimit
			}

			request.DDoSRule = &ddosRule
		}

	}

	modifyErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDDoSPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if modifyErr != nil {
		log.Printf("[CRITAL]%s create teo ddosPolicy failed, reason:%+v", logId, modifyErr)
		return modifyErr
	}

	return resourceTencentCloudTeoDdosPolicyRead(d, meta)
}

func resourceTencentCloudTeoDdosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ddos_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
