package dayuv2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuDdosPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyV2Create,
		Read:   resourceTencentCloudDayuDdosPolicyV2Read,
		Update: resourceTencentCloudDayuDdosPolicyV2Update,
		Delete: resourceTencentCloudDayuDdosPolicyV2Delete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the resource instance.",
			},
			"business": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Business of resource instance. bgpip indicates anti-anti-ip ip; bgp means exclusive package; bgp-multip means shared packet; net indicates anti-anti-ip pro version.",
			},
			"ddos_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "DDoS cleaning threshold, value[0, 60, 80, 100, 150, 200, 250, 300, 400, 500, 700, 1000]; When the value is set to 0, it means that the default value is adopted.",
			},
			"ddos_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protection class, value [`low`, `middle`, `high`].",
			},
			"black_white_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ip of resource instance.",
						},
						"ip_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP type, value [`black`(blacklist IP), `white` (whitelist IP)].",
						},
					},
				},
				Description: "DDoS-protected IP blacklist and whitelist.",
			},
			"acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action, optional values: drop, transmit, forward.",
						},
						"d_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The destination port starts, and the value range is 0~65535.",
						},
						"d_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The destination port ends, and the value range is 0~65535.",
						},
						"s_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The source port starts, and the value range is 0~65535.",
						},
						"s_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The source port ends, and the acceptable value ranges from 0 to 65535.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Policy priority, the lower the number, the higher the level, the higher the rule matches, taking a value of 1-1000.Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"forward_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol type, desirable values tcp, udp, all.",
						},
					},
				},
				Description: "Port ACL policy for DDoS protection.",
			},
			"protocol_block_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drop_icmp": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ICMP block, value [0 (block off), 1 (block on)].",
						},
						"drop_tcp": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "TCP block, value [0 (block off), 1 (block on)].",
						},
						"drop_udp": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "UDP block, value [0 (block off), 1 (block on)].",
						},
						"drop_other": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Other block, value [0 (block off), 1 (block on)].",
						},
					},
				},
				Description: "Protocol block configuration for DDoS protection.",
			},
			"water_print_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"offset": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Watermark offset, value range: [0-100].",
						},
						"open_status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether it is enabled, value [0 (manual open), 1 (immediate operation)].",
						},
						"listeners": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frontend_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Lower limit of forwarding listening port. Values: [1-65535].",
									},
									"forward_protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Forwarding protocol, value [TCP, UDP].",
									},
									"frontend_port_end": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Upper limit of forwarding listening port. Values: [1-65535].",
									},
								},
							},
							Description: "List of forwarding listeners to which the watermark belongs.",
						},
						"verify": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark check mode, value [`checkall`(normal mode), `shortfpcheckall`(simplified mode)].",
						},
					},
				},
				Description: "Water print config.",
			},
			"ddos_connect_limit": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sd_new_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The limit on the number of news per second based on source IP + destination IP.",
						},
						"sd_conn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Concurrent connection control based on source IP + destination IP.",
						},
						"dst_new_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Limit on the number of news per second based on the destination IP.",
						},
						"dst_conn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Concurrent connection control based on destination IP+ destination port.",
						},
						"bad_conn_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Based on connection suppression trigger threshold, value range [0,4294967295].",
						},
						"syn_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Anomalous connection detection condition, percentage of syn ack, value range [0,100].",
						},
						"syn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Anomaly connection detection condition, syn threshold, value range [0,100].",
						},
						"conn_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Abnormal connection detection condition, connection timeout, value range [0,65535].",
						},
						"null_conn_enable": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Abnormal connection detection conditions, empty connection guard switch, value range[0,1].",
						},
					},
				},
				Description: "DDoS connection suppression options.",
			},
			"ddos_ai": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AI protection switch, take the value [`on`, `off`].",
			},
			"ddos_geo_ip_block_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Block action, take the value [`drop`, `trans`].",
						},
						"area_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Required:    true,
							Description: "When the RegionType is customized, the AreaList must be filled in, and a maximum of 128 must be filled in.",
						},
						"region_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone type, value [oversea (overseas),china (domestic),customized (custom region)].",
						},
					},
				},
				Description: "DDoS-protected area block configuration.",
			},
			"ddos_speed_limit_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol_list": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP protocol numbers, take the value[ ALL (all protocols),TCP (tcp protocol),UDP (udp protocol),SMP (smp protocol),1; 2-100 (custom protocol number range, up to 8)].",
						},
						"dst_port_list": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "List of port ranges, up to 8, multiple; Separated, the range is represented with -; this port range must be filled in; fill in the style 1:0-65535, style 2:80; 443; 1000-2000.",
						},
						"mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Speed limit mode, take the value [1 (speed limit based on source IP),2 (speed limit based on destination port)].",
						},
						"packet_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Packet rate pps.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bandwidth bps.",
						},
					},
				},
				Description: "Access speed limit configuration for DDoS protection.",
			},
			"packet_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action, take the value [drop,transmit,drop_black (discard and black out),drop_rst (Interception),drop_black_rst (intercept and block),forward].",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol, value [tcp udp icmp all].",
						},
						"s_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Start the source port, take the value 0~65535.",
						},
						"s_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "End source port, take the value 1~65535, must be greater than or equal to the starting source port.",
						},
						"d_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "From the destination port, take the value 0~65535.",
						},
						"d_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The end destination port, take the value 1~65535, which must be greater than or equal to the starting destination port.",
						},
						"pktlen_min": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minimum message length, 1-1500.",
						},
						"pktlen_max": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The maximum message length, taken from 1 to 1500, must be greater than or equal to the minimum message length.",
						},
						"str": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Detect values, key strings or regular expressions, take the value [When the detection type is sunday, please fill in the string or hexadecimal bytecode, for example \x313233 corresponds to the hexadecimal bytecode of the string `123`;When the detection type is pcre, please fill in the regular expression string;].",
						},
						"str2": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The second detection value, the key string or regular expression, takes the value [When the detection type is sunday, please fill in the string or hexadecimal bytecode, for example \x313233 corresponds to the hexadecimal bytecode of the string `123`;When the detection type is pcre, please fill in the regular expression string;].",
						},
						"match_logic": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "When there is a second detection condition, the and/or relationship with the first detection condition, takes the value [And (and relationship),none (fill in this value when there is no second detection condition)].",
						},
						"match_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Detection type, value [sunday (keyword),pcre (regular expression)].",
						},
						"match_type2": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The second type of detection, takes the value [sunday (keyword),pcre (regular expression)].",
						},
						"match_begin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Detect position, take the value [begin_l3 (IP header),begin_l4 (TCP/UDP header),begin_l5 (T load), no_match (mismatch)].",
						},
						"match_begin2": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The second detection position. take the value [begin_l3 (IP header),begin_l4 (TCP/UDP header),begin_l5 (T load), no_match (mismatch)].",
						},
						"depth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Detection depth from the detection position, value [0,1500].",
						},
						"depth2": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Second detection depth starting from the second detection position, value [0,1500].",
						},
						"offset": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Offset from detection position, value range [0, Depth].",
						},
						"offset2": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Offset from the second detection position, value range [0,Depth2].",
						},
						"is_not": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether to include the detected value, take the value [0 (included),1 (not included)].",
						},
						"is_not2": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether the second detection contains the detected value, the value [0 (included),1 (not included)].",
						},
					},
				},
				Description: "Feature filtering rules for DDoS protection.",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_v2.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	resourceId := d.Get("resource_id").(string)
	business := d.Get("business").(string)

	d.SetId(resourceId + tccommon.FILED_SP + business)
	if v, ok := d.GetOk("black_white_ips"); ok {
		blackWhiteIps := v.([]interface{})
		blacks := make([]string, 0)
		whites := make([]string, 0)
		for _, blackWhiteIpItem := range blackWhiteIps {
			blackWhiteIp := blackWhiteIpItem.(map[string]interface{})
			if blackWhiteIp["ip_type"].(string) == DDOS_BLACK_WHITE_IP_TYPE_WHITE {
				whites = append(whites, blackWhiteIp["ip"].(string))
			}
			if blackWhiteIp["ip_type"].(string) == DDOS_BLACK_WHITE_IP_TYPE_BLACK {
				blacks = append(blacks, blackWhiteIp["ip"].(string))
			}

		}
		if len(blacks) > 0 {
			err := antiddosService.CreateDDoSBlackWhiteIpList(ctx, resourceId, blacks, DDOS_BLACK_WHITE_IP_TYPE_BLACK)
			if err != nil {
				return err
			}

		}
		if len(whites) > 0 {
			err := antiddosService.CreateDDoSBlackWhiteIpList(ctx, resourceId, whites, DDOS_BLACK_WHITE_IP_TYPE_WHITE)
			if err != nil {
				return err
			}
		}

	}

	if v, ok := d.GetOk("ddos_level"); ok {
		ddosLevel := v.(string)
		err := antiddosService.ModifyDDoSLevel(ctx, business, resourceId, ddosLevel)
		if err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("ddos_threshold"); ok {
		threshold := v.(int)
		err := antiddosService.ModifyDDoSThreshold(ctx, business, resourceId, threshold)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("acls"); ok {
		acls := v.([]interface{})
		for _, aclItem := range acls {
			acl := aclItem.(map[string]interface{})
			action := acl["action"].(string)
			dPortStart := acl["d_port_start"].(int)
			dPortEnd := acl["d_port_end"].(int)
			sPortStart := acl["s_port_start"].(int)
			sPortEnd := acl["s_port_end"].(int)
			priority := acl["priority"].(int)
			forwardProtocol := acl["forward_protocol"].(string)
			tmpAclConfig := antiddos.AclConfig{
				Action:          &action,
				DPortStart:      helper.IntUint64(dPortStart),
				DPortEnd:        helper.IntUint64(dPortEnd),
				SPortStart:      helper.IntUint64(sPortStart),
				SPortEnd:        helper.IntUint64(sPortEnd),
				Priority:        helper.IntUint64(priority),
				ForwardProtocol: &forwardProtocol,
			}
			err := antiddosService.CreatePortAclConfig(ctx, resourceId, tmpAclConfig)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("protocol_block_config"); ok {
		protocolBlockConfigs := v.([]interface{})
		for _, protocolBlockConfigItem := range protocolBlockConfigs {
			protocolBlockConfig := protocolBlockConfigItem.(map[string]interface{})
			dropIcmp := protocolBlockConfig["drop_icmp"].(int)
			dropTcp := protocolBlockConfig["drop_tcp"].(int)
			dropUdp := protocolBlockConfig["drop_udp"].(int)
			dropOther := protocolBlockConfig["drop_other"].(int)
			tmpProtocolBlockConfig := antiddos.ProtocolBlockConfig{
				DropIcmp:  helper.IntInt64(dropIcmp),
				DropTcp:   helper.IntInt64(dropTcp),
				DropUdp:   helper.IntInt64(dropUdp),
				DropOther: helper.IntInt64(dropOther),
			}
			err := antiddosService.CreateProtocolBlockConfig(ctx, resourceId, tmpProtocolBlockConfig)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("water_print_config"); ok {
		waterPrintConfigs := v.([]interface{})
		for _, waterPrintConfigItem := range waterPrintConfigs {
			waterPrintConfigs := waterPrintConfigItem.(map[string]interface{})
			offset := waterPrintConfigs["offset"].(int)
			openStatus := waterPrintConfigs["open_status"].(int)
			verify := waterPrintConfigs["verify"].(string)
			listeners := waterPrintConfigs["listeners"].([]interface{})
			listenerList := make([]*antiddos.ForwardListener, 0)
			for _, listenerItem := range listeners {
				listener := listenerItem.(map[string]interface{})
				frontendPort := listener["frontend_port"].(int)
				forwardProtocol := listener["forward_protocol"].(string)
				frontendPortEnd := listener["frontend_port_end"].(int)
				listenerList = append(listenerList, &antiddos.ForwardListener{
					FrontendPort:    helper.IntInt64(frontendPort),
					ForwardProtocol: helper.String(forwardProtocol),
					FrontendPortEnd: helper.IntInt64(frontendPortEnd),
				})
			}
			tmpWaterPrintConfig := antiddos.WaterPrintConfig{
				Offset:     helper.IntInt64(offset),
				OpenStatus: helper.IntInt64(openStatus),
				Verify:     helper.String(verify),
				Listeners:  listenerList,
			}
			err := antiddosService.CreateWaterPrintConfig(ctx, resourceId, tmpWaterPrintConfig)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("ddos_connect_limit"); ok {
		ddosConnectLimits := v.([]interface{})
		for _, ddosConnectLimitItem := range ddosConnectLimits {
			ddosConnectLimit := ddosConnectLimitItem.(map[string]interface{})
			sdNewLimit := ddosConnectLimit["sd_new_limit"].(int)
			sdConnLimit := ddosConnectLimit["sd_conn_limit"].(int)
			dstNewLimit := ddosConnectLimit["dst_new_limit"].(int)
			dstConnLimit := ddosConnectLimit["dst_conn_limit"].(int)
			badConnThreshold := ddosConnectLimit["bad_conn_threshold"].(int)
			synRate := ddosConnectLimit["syn_rate"].(int)
			synLimit := ddosConnectLimit["syn_limit"].(int)
			connTimeout := ddosConnectLimit["conn_timeout"].(int)
			nullConnEnable := ddosConnectLimit["null_conn_enable"].(int)
			tmpConnectLimitConfig := antiddos.ConnectLimitConfig{
				SdNewLimit:       helper.IntUint64(sdNewLimit),
				SdConnLimit:      helper.IntUint64(sdConnLimit),
				DstNewLimit:      helper.IntUint64(dstNewLimit),
				DstConnLimit:     helper.IntUint64(dstConnLimit),
				BadConnThreshold: helper.IntUint64(badConnThreshold),
				SynRate:          helper.IntUint64(synRate),
				SynLimit:         helper.IntUint64(synLimit),
				ConnTimeout:      helper.IntUint64(connTimeout),
				NullConnEnable:   helper.IntUint64(nullConnEnable),
			}
			err := antiddosService.CreateDDoSConnectLimit(ctx, resourceId, tmpConnectLimitConfig)
			if err != nil {
				return err
			}
		}
	}
	if v, ok := d.GetOk("ddos_ai"); ok {
		err := antiddosService.CreateDDoSAI(ctx, resourceId, v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("ddos_geo_ip_block_config"); ok {
		ddosGeoIpBlockConfigs := v.([]interface{})
		for _, ddosGeoIpBlockConfigItem := range ddosGeoIpBlockConfigs {
			ddosGeoIpBlockConfig := ddosGeoIpBlockConfigItem.(map[string]interface{})
			action := ddosGeoIpBlockConfig["action"].(string)
			areaList := ddosGeoIpBlockConfig["area_list"].([]interface{})
			regionType := ddosGeoIpBlockConfig["region_type"].(string)
			areaInt64List := make([]*int64, 0)
			for _, area := range areaList {
				areaInt64List = append(areaInt64List, helper.IntInt64(area.(int)))
			}
			tmpDDoSGeoIPBlockConfig := antiddos.DDoSGeoIPBlockConfig{
				Action:     &action,
				Id:         &resourceId,
				RegionType: &regionType,
				AreaList:   areaInt64List,
			}
			err := antiddosService.CreateDDoSGeoIPBlockConfig(ctx, resourceId, tmpDDoSGeoIPBlockConfig)
			if err != nil {
				return err
			}
		}
	}
	if v, ok := d.GetOk("ddos_speed_limit_config"); ok {
		ddosConnectLimitList := v.([]interface{})
		for _, ddosConnectLimitItem := range ddosConnectLimitList {
			ddosConnectLimit := ddosConnectLimitItem.(map[string]interface{})
			protocolList := ddosConnectLimit["protocol_list"].(string)
			dstPortList := ddosConnectLimit["dst_port_list"].(string)
			mode := ddosConnectLimit["mode"].(int)
			packetRate := ddosConnectLimit["packet_rate"].(int)
			bandwidth := ddosConnectLimit["bandwidth"].(int)

			tmpDDosConnectLimit := antiddos.DDoSSpeedLimitConfig{
				ProtocolList: &protocolList,
				DstPortList:  &dstPortList,
				Mode:         helper.IntUint64(mode),
				SpeedValues: []*antiddos.SpeedValue{
					{
						Type:  helper.IntUint64(1),
						Value: helper.IntUint64(packetRate),
					},
					{
						Type:  helper.IntUint64(2),
						Value: helper.IntUint64(bandwidth),
					},
				},
			}
			err := antiddosService.CreateDDoSSpeedLimitConfig(ctx, resourceId, tmpDDosConnectLimit)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("packet_filters"); ok {
		packetFilters := v.([]interface{})
		for _, packetFilterItem := range packetFilters {
			packetFilter := packetFilterItem.(map[string]interface{})
			action := packetFilter["action"].(string)
			protocol := packetFilter["protocol"].(string)
			sPortStart := packetFilter["s_port_start"].(int)
			sPortEnd := packetFilter["s_port_end"].(int)
			dPortStart := packetFilter["d_port_start"].(int)
			dPortEnd := packetFilter["d_port_end"].(int)
			pktlenMin := packetFilter["pktlen_min"].(int)
			pktlenMax := packetFilter["pktlen_max"].(int)
			str := packetFilter["str"].(string)
			str2 := packetFilter["str2"].(string)
			matchLogic := packetFilter["match_logic"].(string)
			matchType := packetFilter["match_type"].(string)
			matchType2 := packetFilter["match_type2"].(string)
			matchBegin := packetFilter["match_begin"].(string)
			matchBegin2 := packetFilter["match_begin2"].(string)
			depth := packetFilter["depth"].(int)
			depth2 := packetFilter["depth2"].(int)
			offset := packetFilter["offset"].(int)
			offset2 := packetFilter["offset2"].(int)
			isNot := packetFilter["is_not"].(int)
			isNot2 := packetFilter["is_not2"].(int)

			tmpPacketFilter := antiddos.PacketFilterConfig{
				Id:          &resourceId,
				Action:      &action,
				Protocol:    &protocol,
				SportStart:  helper.IntInt64(sPortStart),
				SportEnd:    helper.IntInt64(sPortEnd),
				DportStart:  helper.IntInt64(dPortStart),
				DportEnd:    helper.IntInt64(dPortEnd),
				PktlenMin:   helper.IntInt64(pktlenMin),
				PktlenMax:   helper.IntInt64(pktlenMax),
				Str:         &str,
				Str2:        &str2,
				MatchLogic:  &matchLogic,
				MatchType:   &matchType,
				MatchType2:  &matchType2,
				MatchBegin:  &matchBegin,
				MatchBegin2: &matchBegin2,
				Depth:       helper.IntInt64(depth),
				Offset:      helper.IntInt64(offset),
				IsNot:       helper.IntInt64(isNot),
				Depth2:      helper.IntInt64(depth2),
				Offset2:     helper.IntInt64(offset2),
				IsNot2:      helper.IntInt64(isNot2),
			}
			err := antiddosService.CreatePacketFilterConfig(ctx, resourceId, tmpPacketFilter)
			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudDayuDdosPolicyV2Read(d, meta)
}

func resourceTencentCloudDayuDdosPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	instanceId := items[0]
	protectThresholdRelation, err := antiddosService.DescribeListProtectThresholdConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	ddosLevel := protectThresholdRelation.DDoSLevel
	ddosThreshold := protectThresholdRelation.DDoSThreshold
	_ = d.Set("ddos_level", ddosLevel)
	_ = d.Set("ddos_threshold", ddosThreshold)

	blackWhiteIpRelationList, err := antiddosService.DescribeListBlackWhiteIpList(ctx, instanceId)
	if err != nil {
		return err
	}
	blackWhiteIpInfos := make([]map[string]string, 0)
	for _, blackWhiteIpRelation := range blackWhiteIpRelationList {
		ip := blackWhiteIpRelation.Ip
		ipType := blackWhiteIpRelation.Type
		blackWhiteIpInfos = append(blackWhiteIpInfos, map[string]string{"ip": *ip, "ip_type": *ipType})
	}
	_ = d.Set("black_white_ips", blackWhiteIpInfos)

	aclConfigRelationList, err := antiddosService.DescribeListPortAclList(ctx, instanceId)
	if err != nil {
		return err
	}
	acls := make([]map[string]interface{}, 0)
	for _, aclConfigRelation := range aclConfigRelationList {
		action := aclConfigRelation.AclConfig.Action
		dPortStart := aclConfigRelation.AclConfig.DPortStart
		dPortEnd := aclConfigRelation.AclConfig.DPortEnd
		sPortStart := aclConfigRelation.AclConfig.SPortStart
		sPortEnd := aclConfigRelation.AclConfig.SPortEnd
		priority := aclConfigRelation.AclConfig.Priority
		forwardProtocol := aclConfigRelation.AclConfig.ForwardProtocol
		acl := make(map[string]interface{})
		acl["action"] = action
		acl["d_port_start"] = dPortStart
		acl["d_port_end"] = dPortEnd
		acl["s_port_start"] = sPortStart
		acl["s_port_end"] = sPortEnd
		acl["priority"] = priority
		acl["forward_protocol"] = forwardProtocol
		acls = append(acls, acl)
	}
	_ = d.Set("acls", acls)

	protocolBlockRelation, err := antiddosService.DescribeListProtocolBlockConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	if protocolBlockRelation.ProtocolBlockConfig != nil {
		dropIcmp := protocolBlockRelation.ProtocolBlockConfig.DropIcmp
		dropTcp := protocolBlockRelation.ProtocolBlockConfig.DropTcp
		dropUdp := protocolBlockRelation.ProtocolBlockConfig.DropUdp
		dropOther := protocolBlockRelation.ProtocolBlockConfig.DropOther
		protocolBlockConfig := make(map[string]interface{})
		protocolBlockConfig["drop_icmp"] = dropIcmp
		protocolBlockConfig["drop_tcp"] = dropTcp
		protocolBlockConfig["drop_udp"] = dropUdp
		protocolBlockConfig["drop_other"] = dropOther
		_ = d.Set("protocol_block_config", []map[string]interface{}{protocolBlockConfig})
	}

	waterPrintConfigs, err := antiddosService.DescribeListWaterPrintConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	waterPrintConfigList := make([]map[string]interface{}, 0)

	for _, waterPrintConfig := range waterPrintConfigs {
		waterPrintConfigMap := make(map[string]interface{})
		waterPrintConfigMap["offset"] = waterPrintConfig.WaterPrintConfig.Offset
		waterPrintConfigMap["open_status"] = waterPrintConfig.WaterPrintConfig.OpenStatus
		waterPrintConfigMap["verify"] = waterPrintConfig.WaterPrintConfig.Verify
		listenerList := make([]interface{}, 0)
		for _, listener := range waterPrintConfig.WaterPrintConfig.Listeners {
			listenerMap := make(map[string]interface{})
			listenerMap["frontend_port"] = listener.FrontendPort
			listenerMap["forward_protocol"] = listener.ForwardProtocol
			listenerMap["frontend_port_end"] = listener.FrontendPortEnd
			listenerList = append(listenerList, listenerMap)
		}
		waterPrintConfigMap["listeners"] = listenerList
		waterPrintConfigList = append(waterPrintConfigList, waterPrintConfigMap)
	}
	_ = d.Set("water_print_config", waterPrintConfigList)

	connectLimitRelation, err := antiddosService.DescribeDDoSConnectLimitList(ctx, instanceId)
	if err != nil {
		return err
	}
	ddosConnectLimit := make(map[string]interface{})
	ddosConnectLimit["sd_new_limit"] = connectLimitRelation.SdNewLimit
	ddosConnectLimit["sd_conn_limit"] = connectLimitRelation.SdConnLimit
	ddosConnectLimit["dst_new_limit"] = connectLimitRelation.DstNewLimit
	ddosConnectLimit["dst_conn_limit"] = connectLimitRelation.DstConnLimit
	ddosConnectLimit["bad_conn_threshold"] = connectLimitRelation.BadConnThreshold
	ddosConnectLimit["syn_rate"] = connectLimitRelation.SynRate
	ddosConnectLimit["syn_limit"] = connectLimitRelation.SynLimit
	ddosConnectLimit["conn_timeout"] = connectLimitRelation.ConnTimeout
	ddosConnectLimit["null_conn_enable"] = connectLimitRelation.NullConnEnable
	_ = d.Set("ddos_connect_limit", []map[string]interface{}{ddosConnectLimit})

	ddoSAIRelation, err := antiddosService.DescribeListDDoSAI(ctx, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("ddos_ai", ddoSAIRelation.DDoSAI)

	ddosGeoIPBlockConfigRelations, err := antiddosService.DescribeListDDoSGeoIPBlockConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	ddosGeoIPBlockConfigList := make([]map[string]interface{}, 0)

	for _, ddosGeoIPBlockConfigRelation := range ddosGeoIPBlockConfigRelations {
		ddosGeoIPBlockConfig := make(map[string]interface{})
		ddosGeoIPBlockConfig["action"] = ddosGeoIPBlockConfigRelation.GeoIPBlockConfig.Action
		ddosGeoIPBlockConfig["area_list"] = ddosGeoIPBlockConfigRelation.GeoIPBlockConfig.AreaList
		ddosGeoIPBlockConfig["region_type"] = ddosGeoIPBlockConfigRelation.GeoIPBlockConfig.RegionType
		ddosGeoIPBlockConfigList = append(ddosGeoIPBlockConfigList, ddosGeoIPBlockConfig)
	}

	_ = d.Set("ddos_geo_ip_block_config", ddosGeoIPBlockConfigList)

	ddosSpeedLimitConfigRelations, err := antiddosService.DescribeListDDoSSpeedLimitConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	ddosSpeedLimitConfigs := make([]map[string]interface{}, 0)
	for _, ddosSpeedLimitConfigRelation := range ddosSpeedLimitConfigRelations {
		ddosSpeedLimitConfig := make(map[string]interface{})
		ddosSpeedLimitConfig["protocol_list"] = ddosSpeedLimitConfigRelation.SpeedLimitConfig.ProtocolList
		ddosSpeedLimitConfig["dst_port_list"] = ddosSpeedLimitConfigRelation.SpeedLimitConfig.DstPortList
		ddosSpeedLimitConfig["mode"] = ddosSpeedLimitConfigRelation.SpeedLimitConfig.Mode
		for _, speedValue := range ddosSpeedLimitConfigRelation.SpeedLimitConfig.SpeedValues {
			if *speedValue.Type == uint64(1) {
				ddosSpeedLimitConfig["packet_rate"] = *speedValue.Value
			}
			if *speedValue.Type == uint64(2) {
				ddosSpeedLimitConfig["bandwidth"] = *speedValue.Value
			}

		}
		ddosSpeedLimitConfigs = append(ddosSpeedLimitConfigs, ddosSpeedLimitConfig)
	}

	_ = d.Set("ddos_speed_limit_config", ddosSpeedLimitConfigs)

	packetFilterRelationList, err := antiddosService.DescribeListPacketFilterConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	packetFilters := make([]map[string]interface{}, 0)
	for _, packetFilterRelation := range packetFilterRelationList {
		tmpPacketFilter := make(map[string]interface{})
		tmpPacketFilter["s_port_start"] = packetFilterRelation.PacketFilterConfig.SportStart
		tmpPacketFilter["s_port_end"] = packetFilterRelation.PacketFilterConfig.SportEnd
		tmpPacketFilter["d_port_start"] = packetFilterRelation.PacketFilterConfig.DportStart
		tmpPacketFilter["d_port_end"] = packetFilterRelation.PacketFilterConfig.DportEnd
		tmpPacketFilter["pktlen_min"] = packetFilterRelation.PacketFilterConfig.PktlenMin
		tmpPacketFilter["pktlen_max"] = packetFilterRelation.PacketFilterConfig.PktlenMax
		tmpPacketFilter["str"] = packetFilterRelation.PacketFilterConfig.Str
		tmpPacketFilter["str2"] = packetFilterRelation.PacketFilterConfig.Str2
		tmpPacketFilter["match_logic"] = packetFilterRelation.PacketFilterConfig.MatchLogic
		tmpPacketFilter["match_type"] = packetFilterRelation.PacketFilterConfig.MatchType
		tmpPacketFilter["match_type2"] = packetFilterRelation.PacketFilterConfig.MatchType2
		tmpPacketFilter["match_begin"] = packetFilterRelation.PacketFilterConfig.MatchBegin
		tmpPacketFilter["match_begin2"] = packetFilterRelation.PacketFilterConfig.MatchBegin2
		tmpPacketFilter["action"] = packetFilterRelation.PacketFilterConfig.Action
		tmpPacketFilter["depth"] = packetFilterRelation.PacketFilterConfig.Depth
		tmpPacketFilter["depth2"] = packetFilterRelation.PacketFilterConfig.Depth2
		tmpPacketFilter["offset"] = packetFilterRelation.PacketFilterConfig.Offset
		tmpPacketFilter["offset2"] = packetFilterRelation.PacketFilterConfig.Offset2
		tmpPacketFilter["protocol"] = packetFilterRelation.PacketFilterConfig.Protocol
		packetFilters = append(packetFilters, tmpPacketFilter)
	}
	_ = d.Set("packet_filters", packetFilters)

	return nil
}

func resourceTencentCloudDayuDdosPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceId := items[0]
	business := items[1]
	d.Partial(true)

	if d.HasChange("ddos_level") {
		ddosLevel := d.Get("ddos_level").(string)
		err := antiddosService.ModifyDDoSLevel(ctx, business, resourceId, ddosLevel)
		if err != nil {
			return err
		}
	}
	if d.HasChange("ddos_threshold") {
		threshold := d.Get("ddos_threshold").(int)
		err := antiddosService.ModifyDDoSThreshold(ctx, business, resourceId, threshold)
		if err != nil {
			return err
		}
	}

	if d.HasChange("protocol_block_config") {
		protocolBlockConfigs := d.Get("protocol_block_config").([]interface{})
		for _, protocolBlockConfigItem := range protocolBlockConfigs {
			protocolBlockConfig := protocolBlockConfigItem.(map[string]interface{})
			dropIcmp := protocolBlockConfig["drop_icmp"].(int)
			dropTcp := protocolBlockConfig["drop_tcp"].(int)
			dropUdp := protocolBlockConfig["drop_udp"].(int)
			dropOther := protocolBlockConfig["drop_other"].(int)
			tmpProtocolBlockConfig := antiddos.ProtocolBlockConfig{
				DropIcmp:  helper.IntInt64(dropIcmp),
				DropTcp:   helper.IntInt64(dropTcp),
				DropUdp:   helper.IntInt64(dropUdp),
				DropOther: helper.IntInt64(dropOther),
			}
			err := antiddosService.CreateProtocolBlockConfig(ctx, resourceId, tmpProtocolBlockConfig)
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("ddos_connect_limit") {
		ddosConnectLimits := d.Get("ddos_connect_limit").([]interface{})
		for _, ddosConnectLimitItem := range ddosConnectLimits {
			ddosConnectLimit := ddosConnectLimitItem.(map[string]interface{})
			sdNewLimit := ddosConnectLimit["sd_new_limit"].(int)
			sdConnLimit := ddosConnectLimit["sd_conn_limit"].(int)
			dstNewLimit := ddosConnectLimit["dst_new_limit"].(int)
			dstConnLimit := ddosConnectLimit["dst_conn_limit"].(int)
			badConnThreshold := ddosConnectLimit["bad_conn_threshold"].(int)
			synRate := ddosConnectLimit["syn_rate"].(int)
			synLimit := ddosConnectLimit["syn_limit"].(int)
			connTimeout := ddosConnectLimit["conn_timeout"].(int)
			nullConnEnable := ddosConnectLimit["null_conn_enable"].(int)
			tmpConnectLimitConfig := antiddos.ConnectLimitConfig{
				SdNewLimit:       helper.IntUint64(sdNewLimit),
				SdConnLimit:      helper.IntUint64(sdConnLimit),
				DstNewLimit:      helper.IntUint64(dstNewLimit),
				DstConnLimit:     helper.IntUint64(dstConnLimit),
				BadConnThreshold: helper.IntUint64(badConnThreshold),
				SynRate:          helper.IntUint64(synRate),
				SynLimit:         helper.IntUint64(synLimit),
				ConnTimeout:      helper.IntUint64(connTimeout),
				NullConnEnable:   helper.IntUint64(nullConnEnable),
			}
			err := antiddosService.CreateDDoSConnectLimit(ctx, resourceId, tmpConnectLimitConfig)
			if err != nil {
				return err
			}
		}

	}
	if d.HasChange("ddos_ai") {
		err := antiddosService.CreateDDoSAI(ctx, resourceId, d.Get("ddos_ai").(string))
		if err != nil {
			return err
		}

	}

	if d.HasChange("black_white_ips") {
		blackWhiteIps := d.Get("black_white_ips").([]interface{})
		oldUniqIpMap := make(map[string]int)
		oldBlackWhiteIps, err := antiddosService.DescribeListBlackWhiteIpList(ctx, resourceId)
		if err != nil {
			return err
		}
		for _, oldBlackWhiteIp := range oldBlackWhiteIps {
			ip := oldBlackWhiteIp.Ip
			ipType := oldBlackWhiteIp.Type
			key := *ipType + "_" + *ip
			oldUniqIpMap[key] = 1
		}
		newUniqIpMap := make(map[string]int)
		for _, blackWhiteIpItem := range blackWhiteIps {
			blackWhiteIp := blackWhiteIpItem.(map[string]interface{})
			key := blackWhiteIp["ip_type"].(string) + "_" + blackWhiteIp["ip"].(string)
			newUniqIpMap[key] = 1
			if oldUniqIpMap["key"] == 0 {
				err := antiddosService.CreateDDoSBlackWhiteIpList(ctx, resourceId, []string{blackWhiteIp["ip"].(string)}, blackWhiteIp["ip_type"].(string))
				if err != nil {
					return err
				}
			}
		}
		for _, oldBlackWhiteIp := range oldBlackWhiteIps {
			ip := oldBlackWhiteIp.Ip
			ipType := oldBlackWhiteIp.Type
			key := *ipType + "_" + *ip
			if newUniqIpMap[key] == 0 {
				_ = antiddosService.DeleteDDoSBlackWhiteIpList(ctx, resourceId, []string{*ip}, *ipType)
			}
		}

	}

	if d.HasChange("acls") {
		oldAclConfigRelationList, err := antiddosService.DescribeListPortAclList(ctx, resourceId)
		if err != nil {
			return err
		}
		oldAclConfigs := make([]interface{}, 0)
		for _, oldAclConfigRelation := range oldAclConfigRelationList {
			oldAclConfigs = append(oldAclConfigs, *oldAclConfigRelation.AclConfig)
		}
		newAclConfigs := make([]interface{}, 0)
		acls := d.Get("acls").([]interface{})
		for _, aclItem := range acls {
			acl := aclItem.(map[string]interface{})
			action := acl["action"].(string)
			dPortStart := acl["d_port_start"].(int)
			dPortEnd := acl["d_port_end"].(int)
			sPortStart := acl["s_port_start"].(int)
			sPortEnd := acl["s_port_end"].(int)
			priority := acl["priority"].(int)
			forwardProtocol := acl["forward_protocol"].(string)
			tmpAclConfig := antiddos.AclConfig{
				Action:          &action,
				DPortStart:      helper.IntUint64(dPortStart),
				DPortEnd:        helper.IntUint64(dPortEnd),
				SPortStart:      helper.IntUint64(sPortStart),
				SPortEnd:        helper.IntUint64(sPortEnd),
				Priority:        helper.IntUint64(priority),
				ForwardProtocol: &forwardProtocol,
			}
			newAclConfigs = append(newAclConfigs, tmpAclConfig)
		}
		increments, decrements := DeltaList(oldAclConfigs, newAclConfigs)
		for _, decrementItem := range decrements {
			var decrement antiddos.AclConfig
			_ = json.Unmarshal([]byte(decrementItem), &decrement)
			err := antiddosService.DeletePortAclConfig(ctx, resourceId, decrement)
			if err != nil {
				return err
			}
		}
		for _, incrementItems := range increments {
			var increment antiddos.AclConfig
			_ = json.Unmarshal([]byte(incrementItems), &increment)
			err := antiddosService.CreatePortAclConfig(ctx, resourceId, increment)
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("water_print_config.0.offset") || d.HasChange("water_print_config.0.listeners") || d.HasChange("water_print_config.0.verify") {
		oldWaterPrintConfigList, err := antiddosService.DescribeListWaterPrintConfig(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(oldWaterPrintConfigList) > 0 {
			err := antiddosService.DeleteWaterPrintConfig(ctx, resourceId)
			if err != nil {
				return err
			}
		}

		waterPrintConfigs := d.Get("water_print_config").([]interface{})
		if len(waterPrintConfigs) > 0 {
			waterPrintConfigItem := waterPrintConfigs[0]
			waterPrintConfigItemMap := waterPrintConfigItem.(map[string]interface{})
			offset := waterPrintConfigItemMap["offset"].(int)
			openStatus := waterPrintConfigItemMap["open_status"].(int)
			verify := waterPrintConfigItemMap["verify"].(string)
			listeners := waterPrintConfigItemMap["listeners"].([]interface{})
			listenerList := make([]*antiddos.ForwardListener, 0)
			for _, listenerItem := range listeners {
				listenerMap := listenerItem.(map[string]interface{})
				frontendPort := listenerMap["frontend_port"].(int)
				forwardProtocol := listenerMap["forward_protocol"].(string)
				frontendPortEnd := listenerMap["frontend_port_end"].(int)
				listenerList = append(listenerList, &antiddos.ForwardListener{
					FrontendPort:    helper.IntInt64(frontendPort),
					ForwardProtocol: helper.String(forwardProtocol),
					FrontendPortEnd: helper.IntInt64(frontendPortEnd),
				})
			}
			tmpWaterPrintConfig := antiddos.WaterPrintConfig{
				Offset:     helper.IntInt64(offset),
				OpenStatus: helper.IntInt64(openStatus),
				Verify:     helper.String(verify),
				Listeners:  listenerList,
			}
			err := antiddosService.CreateWaterPrintConfig(ctx, resourceId, tmpWaterPrintConfig)
			if err != nil {
				return err
			}

		}

	}

	if d.HasChange("water_print_config.0.open_status") {
		waterPrintConfigs := d.Get("water_print_config").([]interface{})
		if len(waterPrintConfigs) > 0 {
			waterPrintConfigItem := waterPrintConfigs[0]
			waterPrintConfigItemMap := waterPrintConfigItem.(map[string]interface{})
			openStatus := waterPrintConfigItemMap["open_status"].(int)
			err := antiddosService.SwitchWaterPrintConfig(ctx, resourceId, openStatus)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("ddos_geo_ip_block_config") {
		oldDDoSGeoIPBlockConfigRelations, err := antiddosService.DescribeListDDoSGeoIPBlockConfig(ctx, resourceId)
		if err != nil {
			return err
		}

		oldDDoSGeoIPBlockConfigs := make([]interface{}, 0)
		for _, oldDDoSGeoIPBlockConfigRelation := range oldDDoSGeoIPBlockConfigRelations {
			ddosGeoIPBlockConfig := *oldDDoSGeoIPBlockConfigRelation.GeoIPBlockConfig
			oldDDoSGeoIPBlockConfigs = append(oldDDoSGeoIPBlockConfigs, ddosGeoIPBlockConfig)
		}

		newDDoSGeoIPBlockConfigs := make([]interface{}, 0)
		for _, ddosGeoIpBlockConfigItem := range d.Get("ddos_geo_ip_block_config").([]interface{}) {
			ddosGeoIpBlockConfig := ddosGeoIpBlockConfigItem.(map[string]interface{})
			action := ddosGeoIpBlockConfig["action"].(string)
			areaList := ddosGeoIpBlockConfig["area_list"].([]interface{})
			regionType := ddosGeoIpBlockConfig["region_type"].(string)
			areaInt64List := make([]*int64, 0)
			for _, area := range areaList {
				areaInt64List = append(areaInt64List, helper.IntInt64(area.(int)))
			}
			tmpDDoSGeoIPBlockConfig := antiddos.DDoSGeoIPBlockConfig{
				Action:     &action,
				Id:         &resourceId,
				RegionType: &regionType,
				AreaList:   areaInt64List,
			}
			newDDoSGeoIPBlockConfigs = append(newDDoSGeoIPBlockConfigs, tmpDDoSGeoIPBlockConfig)
		}
		increments, decrements := DeltaList(oldDDoSGeoIPBlockConfigs, newDDoSGeoIPBlockConfigs)
		for _, decrementItem := range decrements {
			var decrement antiddos.DDoSGeoIPBlockConfig
			_ = json.Unmarshal([]byte(decrementItem), &decrement)
			err := antiddosService.DeleteDDoSGeoIPBlockConfig(ctx, resourceId, decrement)
			if err != nil {
				return err
			}
		}
		for _, incrementItems := range increments {
			var increment antiddos.DDoSGeoIPBlockConfig
			_ = json.Unmarshal([]byte(incrementItems), &increment)
			err := antiddosService.CreateDDoSGeoIPBlockConfig(ctx, resourceId, increment)
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("ddos_speed_limit_config") {
		oldDDoSSpeedLimitConfigRelations, err := antiddosService.DescribeListDDoSSpeedLimitConfig(ctx, resourceId)
		if err != nil {
			return err
		}

		oldDDoSSpeedLimitConfigs := make([]interface{}, 0)
		for _, oldDDoSSpeedLimitConfigRelation := range oldDDoSSpeedLimitConfigRelations {
			oldDDoSSpeedLimitConfig := *oldDDoSSpeedLimitConfigRelation.SpeedLimitConfig
			oldDDoSSpeedLimitConfigs = append(oldDDoSSpeedLimitConfigs, oldDDoSSpeedLimitConfig)
		}

		newDDoSSpeedLimitConfigs := make([]interface{}, 0)
		for _, ddosConnectLimitItem := range d.Get("ddos_speed_limit_config").([]interface{}) {
			ddosConnectLimit := ddosConnectLimitItem.(map[string]interface{})
			protocolList := ddosConnectLimit["protocol_list"].(string)
			dstPortList := ddosConnectLimit["dst_port_list"].(string)
			mode := ddosConnectLimit["mode"].(int)
			packetRate := ddosConnectLimit["packet_rate"].(int)
			bandwidth := ddosConnectLimit["bandwidth"].(int)

			tmpDDosConnectLimit := antiddos.DDoSSpeedLimitConfig{
				ProtocolList: &protocolList,
				DstPortList:  &dstPortList,
				Mode:         helper.IntUint64(mode),
				SpeedValues: []*antiddos.SpeedValue{
					{
						Type:  helper.IntUint64(1),
						Value: helper.IntUint64(packetRate),
					},
					{
						Type:  helper.IntUint64(2),
						Value: helper.IntUint64(bandwidth),
					},
				},
			}
			newDDoSSpeedLimitConfigs = append(newDDoSSpeedLimitConfigs, tmpDDosConnectLimit)
		}

		increments, decrements := DeltaList(oldDDoSSpeedLimitConfigs, newDDoSSpeedLimitConfigs)
		for _, decrementItem := range decrements {
			var decrement antiddos.DDoSSpeedLimitConfig
			_ = json.Unmarshal([]byte(decrementItem), &decrement)
			err := antiddosService.DeleteDDoSSpeedLimitConfig(ctx, resourceId, decrement)
			if err != nil {
				return err
			}
		}
		for _, incrementItems := range increments {
			var increment antiddos.DDoSSpeedLimitConfig
			_ = json.Unmarshal([]byte(incrementItems), &increment)
			err := antiddosService.CreateDDoSSpeedLimitConfig(ctx, resourceId, increment)
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("packet_filters") {
		oldPacketFilterRelationList, err := antiddosService.DescribeListPacketFilterConfig(ctx, resourceId)
		if err != nil {
			return err
		}
		oldPacketFilters := make([]interface{}, 0)
		for _, packetFilterRelation := range oldPacketFilterRelationList {
			tmpPacketFilter := packetFilterRelation.PacketFilterConfig
			oldPacketFilters = append(oldPacketFilters, tmpPacketFilter)
		}

		newPacketFilters := make([]interface{}, 0)

		for _, packetFilterItem := range d.Get("packet_filters").([]interface{}) {
			packetFilter := packetFilterItem.(map[string]interface{})
			action := packetFilter["action"].(string)
			protocol := packetFilter["protocol"].(string)
			sPortStart := packetFilter["s_port_start"].(int)
			sPortEnd := packetFilter["s_port_end"].(int)
			dPortStart := packetFilter["d_port_start"].(int)
			dPortEnd := packetFilter["d_port_end"].(int)
			pktlenMin := packetFilter["pktlen_min"].(int)
			pktlenMax := packetFilter["pktlen_max"].(int)
			str := packetFilter["str"].(string)
			str2 := packetFilter["str2"].(string)
			matchLogic := packetFilter["match_logic"].(string)
			matchType := packetFilter["match_type"].(string)
			matchType2 := packetFilter["match_type2"].(string)
			matchBegin := packetFilter["match_begin"].(string)
			matchBegin2 := packetFilter["match_begin2"].(string)
			depth := packetFilter["depth"].(int)
			depth2 := packetFilter["depth2"].(int)
			offset := packetFilter["offset"].(int)
			offset2 := packetFilter["offset2"].(int)
			isNot := packetFilter["is_not"].(int)
			isNot2 := packetFilter["is_not2"].(int)

			tmpPacketFilter := antiddos.PacketFilterConfig{
				Id:          &resourceId,
				Action:      &action,
				Protocol:    &protocol,
				SportStart:  helper.IntInt64(sPortStart),
				SportEnd:    helper.IntInt64(sPortEnd),
				DportStart:  helper.IntInt64(dPortStart),
				DportEnd:    helper.IntInt64(dPortEnd),
				PktlenMin:   helper.IntInt64(pktlenMin),
				PktlenMax:   helper.IntInt64(pktlenMax),
				Str:         &str,
				Str2:        &str2,
				MatchLogic:  &matchLogic,
				MatchType:   &matchType,
				MatchType2:  &matchType2,
				MatchBegin:  &matchBegin,
				MatchBegin2: &matchBegin2,
				Depth:       helper.IntInt64(depth),
				Offset:      helper.IntInt64(offset),
				IsNot:       helper.IntInt64(isNot),
				Depth2:      helper.IntInt64(depth2),
				Offset2:     helper.IntInt64(offset2),
				IsNot2:      helper.IntInt64(isNot2),
			}
			newPacketFilters = append(newPacketFilters, tmpPacketFilter)
		}
		increments, decrements := DeltaList(oldPacketFilters, newPacketFilters)

		for _, decrementItem := range decrements {
			var decrement antiddos.PacketFilterConfig
			_ = json.Unmarshal([]byte(decrementItem), &decrement)
			err := antiddosService.DeletePacketFilterConfig(ctx, resourceId, decrement)
			if err != nil {
				return err
			}
		}
		for _, incrementItem := range increments {
			var increment antiddos.PacketFilterConfig
			_ = json.Unmarshal([]byte(incrementItem), &increment)
			err := antiddosService.CreatePacketFilterConfig(ctx, resourceId, increment)
			if err != nil {
				return err
			}
		}

	}

	d.Partial(false)

	return resourceTencentCloudDayuDdosPolicyV2Read(d, meta)
}

func resourceTencentCloudDayuDdosPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceId := items[0]
	business := items[1]

	_ = antiddosService.DeleteyDDoSLevel(ctx, business, resourceId)
	_ = antiddosService.DeleteDDoSThreshold(ctx, business, resourceId)

	blackWhiteIpRelationList, err := antiddosService.DescribeListBlackWhiteIpList(ctx, resourceId)
	if err != nil {
		return err
	}
	blacks := make([]string, 0)
	whites := make([]string, 0)
	for _, blackWhiteIpRelation := range blackWhiteIpRelationList {
		ip := blackWhiteIpRelation.Ip
		ipType := blackWhiteIpRelation.Type
		if *ipType == DDOS_BLACK_WHITE_IP_TYPE_BLACK {
			blacks = append(blacks, *ip)
		}
		if *ipType == DDOS_BLACK_WHITE_IP_TYPE_WHITE {
			whites = append(whites, *ip)
		}
	}
	if len(blacks) > 0 {
		_ = antiddosService.DeleteDDoSBlackWhiteIpList(ctx, resourceId, blacks, DDOS_BLACK_WHITE_IP_TYPE_BLACK)
	}
	if len(whites) > 0 {
		_ = antiddosService.DeleteDDoSBlackWhiteIpList(ctx, resourceId, whites, DDOS_BLACK_WHITE_IP_TYPE_WHITE)
	}

	aclConfigRelationList, err := antiddosService.DescribeListPortAclList(ctx, resourceId)
	if err != nil {
		return err
	}
	for _, aclConfigRelation := range aclConfigRelationList {
		deleteAclConfigRelation := aclConfigRelation
		_ = antiddosService.DeletePortAclConfig(ctx, resourceId, *deleteAclConfigRelation.AclConfig)
	}

	_ = antiddosService.DeleteProtocolBlockConfig(ctx, resourceId)
	_ = antiddosService.DeleteDDoSConnectLimit(ctx, resourceId)
	_ = antiddosService.DeleteDDoSAI(ctx, resourceId)

	ddosGeoIPBlockConfigRelations, err := antiddosService.DescribeListDDoSGeoIPBlockConfig(ctx, resourceId)
	if err != nil {
		return err
	}
	for _, ddosGeoIPBlockConfigRelation := range ddosGeoIPBlockConfigRelations {
		_ = antiddosService.DeleteDDoSGeoIPBlockConfig(ctx, resourceId, *ddosGeoIPBlockConfigRelation.GeoIPBlockConfig)
	}

	ddosSpeedLimitConfigRelations, err := antiddosService.DescribeListDDoSSpeedLimitConfig(ctx, resourceId)
	if err != nil {
		return err
	}
	for _, ddosSpeedLimitConfigRelation := range ddosSpeedLimitConfigRelations {
		_ = antiddosService.DeleteDDoSSpeedLimitConfig(ctx, resourceId, *ddosSpeedLimitConfigRelation.SpeedLimitConfig)

	}

	packetFilterRelationList, err := antiddosService.DescribeListPacketFilterConfig(ctx, resourceId)
	if err != nil {
		return err
	}

	for _, packetFilterRelation := range packetFilterRelationList {
		_ = antiddosService.DeletePacketFilterConfig(ctx, resourceId, *packetFilterRelation.PacketFilterConfig)
	}

	oldWaterPrintConfigList, err := antiddosService.DescribeListWaterPrintConfig(ctx, resourceId)
	if err != nil {
		return err
	}
	if len(oldWaterPrintConfigList) > 0 {
		err := antiddosService.DeleteWaterPrintConfig(ctx, resourceId)
		if err != nil {
			return err
		}
	}
	return nil
}
