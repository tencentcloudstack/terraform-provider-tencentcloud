/*
Use this resource to create dayu DDoS policy

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy" "test_policy" {
  resource_type         = "bgpip"
  name                  = "tf_test_policy"
  black_ips = ["1.1.1.1"]
  white_ips = ["2.2.2.2]

  drop_options{
    drop_tcp  = true
	drop_udp  = true
	drop_icmp  = true
	drop_other  = true
	drop_abroad  = true
	check_sync_conn = true
	s_new_limit = 100
	d_new_limit = 100
	s_conn_limit = 100
	d_conn_limit = 100
	tcp_mbps_limit = 100
	udp_mbps_limit = 100
	icmp_mbps_limit = 100
	other_mbps_limit = 100
	bad_conn_threshold = 100
	null_conn_enable = true
	conn_timeout = 500
	syn_rate = 50
	syn_limit = 100
  }

  port_limits{
	start_port = "2000"
	end_port = "2500"
	protocol = "all"
  	action = "drop"
	kind = 1
  }

  packet_filters{
	protocol = "tcp"
	action = "drop"
	d_start_port = 1000
	d_end_port = 1500
	s_start_port = 2000
	s_end_port = 2500
	pkt_length_max = 1400
	pkt_length_min = 1000
	is_include = true
	match_begin = "begin_l5"
	match_type = "pcre"
	depth = 1000
	offset = 500
  }

  watermark_filters{
  	tcp_port_list = ["2000-3000", "3500-4000"]
	udp_port_list = ["5000-6000"]
	offset = 50
	auto_remove = true
	open_switch = true
  }
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func resourceTencentCloudDayuDdosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyCreate,
		Read:   resourceTencentCloudDayuDdosPolicyRead,
		Update: resourceTencentCloudDayuDdosPolicyUpdate,
		Delete: resourceTencentCloudDayuDdosPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 32),
				Description:  "Name of the DDoS policy. Length should between 1 and 32.",
			},
			"drop_options": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drop_tcp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop TCP protocol or not.",
						},
						"drop_udp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate to drop UDP protocol or not.",
						},
						"drop_icmp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop ICMP protocol or not.",
						},
						"drop_other": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop other protocols(exclude TCP/UDP/ICMP) or not.",
						},
						"drop_abroad": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop abroad traffic or not.",
						},
						"check_sync_conn": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to check null connection or not.",
						},
						"d_new_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of new connections based on destination IP, and valid value is range from 0 to 4294967295.",
						},
						"d_conn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of concurrent connections based on destination IP, and valid value is range from 0 to 4294967295.",
						},
						"s_new_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of new connections based on source IP, and valid value is range from 0 to 4294967295.",
						},
						"s_conn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of concurrent connections based on source IP, and valid value is range from 0 to 4294967295.",
						},
						"bad_conn_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The number of new connections based on destination IP that trigger suppression of connections, and valid value is range from 0 to 4294967295.",
						},
						"null_conn_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate to enable null connection or not.",
						},
						"conn_timeout": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 65535),
							Description:  "Connection timeout of abnormal connection check, and valid value is range from 0 to 65535.",
						},
						"syn_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
							Description:  "The percentage of syn in ack of abnormal connection check, and valid value is range from 0 to 100.",
						},
						"syn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
							Description:  "The limit of syn of abnormal connection check, and valid value is range from 0 to 100.",
						},
						"tcp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of TCP traffic, and valid value is range from 0 to 4294967295(Mbps).",
						},
						"udp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of UDP traffic rate, and valid value is range from 0 to 4294967295(Mbps).",
						},
						"icmp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of ICMP traffic rate, and valid value is range from 0 to 4294967295(Mbps).",
						},
						"other_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 4294967295),
							Description:  "The limit of other protocols(exclude TCP/UDP/ICMP) traffic rate, and valid value is range from 0 to 4294967295(Mbps).",
						},
					},
				},
				Description: "Option list of abnormal check of the DDos policy, should set at least one policy.",
			},
			"port_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_PROTOCOL),
							Description:  "Protocol, valid values are `tcp`, `udp`, `icmp`, `all`.",
						},
						"start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validatePort,
							Description:  "Start port, valid value is range from 0 to 65535.",
						},
						"end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validatePort,
							Description:  "End port, valid value is range from 0 to 65535. It must be greater than `d_start_port`.",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_PORT_ACTION),
							Description:  "Action of port to take, valid values area `drop`, `transmit`.",
						},
						"kind": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
							Description:  "The type of forbidden port, and valid values are 0, 1, 2. 0 for destination ports make effect, 1 for source ports make effect. 2 for both destination and source ports.",
						},
					},
				},
				Description: "Port limits of abnormal check of the DDos policy.",
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_PROTOCOL),
							Description:  "Protocol, valid values are `tcp`, `udp`, `icmp`, `all`.",
						},
						"d_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "Start port of the destination, valid value is range from 0 to 65535.",
						},
						"d_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "End port of the destination, valid value is range from 0 to 65535. It must be greater than `d_start_port`.",
						},
						"s_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "Start port of the source, valid value is range from 0 to 65535.",
						},
						"s_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "End port of the source, valid value is range from 0 to 65535. It must be greater than `s_start_port`.",
						},
						"pkt_length_min": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The minimum length of the packet, and valid value is range from 0 to 1500(Mbps).",
						},
						"pkt_length_max": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The max length of the packet, and valid value is range from 0 to 1500(Mbps). It must be greater than `pkt_length_min`.",
						},
						"match_begin": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_MATCH_SWITCH),
							Description:  "Indicate whether to check load or not, `begin_l5` means to match and `no_match` means not.",
						},
						"match_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_MATCH_TYPE),
							Description:  "Match type, valid values are `sunday` and `pcre`, `sunday` means key word match while `pcre` means regular match.",
						},
						"match_str": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key word or regular expression.",
						},
						"depth": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The depth of match, and valid value is range from 0 to 1500.",
						},
						"offset": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The offset of match, and valid value is range from 0 to 1500.",
						},
						"is_include": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicate whether to include the key word/regular expression or not.",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_PACKET_ACTION),
							Description:  "Action of port to take, valid values area `drop`(drop the packet), `drop_black`(drop the packet and black the ip),`drop_rst`(drop the packet and disconnect),`drop_black_rst`(drop the packet, black the ip and disconnect),`transmit`(transmit the packet).",
						},
					},
				},
				Description: "Message filter options list.",
			},
			"watermark_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp_port_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validatePortRange,
							},
							Description: "Port range of TCP, the format is like `2000-3000`.",
						},
						"udp_port_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validatePortRange,
							},
							Description: "Port range of TCP, the format is like `2000-3000`.",
						},
						"offset": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
							Description:  "The offset of water print, and valid value is range from 0 to 100.",
						},
						"auto_remove": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate whether to auto-remove the water print or not.",
						},
						"open_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate whether to open water print or not.",
						},
					},
				},
				Description: "Water print policy options, and only support one water print policy at most.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the DDos policy.",
			},
			"scene_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of scene that the DDos policy works for.",
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
							Description: "Indicate whether to auto-remove the water print or not.",
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
	}
}

func resourceTencentCloudDayuDdosPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	resourceType := d.Get("resource_type").(string)
	name := d.Get("name").(string)
	//set DDosPolicyDropOption
	dropMapping := d.Get("drop_options").([]interface{})
	ddosPolicyDropOption, _ := setDdosPolicyDropOption(dropMapping)

	//set DDoSPolicyPortLimit
	portMapping := d.Get("port_filters").([]interface{})
	ddosPolicyPortLimit, lErr := setDdosPolicyPortLimit(portMapping)

	if lErr != nil {
		return lErr
	}

	//set IpBlackWhite
	blackIps := d.Get("black_ips").(*schema.Set).List()
	whiteIps := d.Get("white_ips").(*schema.Set).List()
	ipBlackWhite, ipErr := setIpBlackWhite(blackIps, whiteIps)

	if ipErr != nil {
		return ipErr
	}
	//set DDoSPolicyPacketFilter
	packetFilterMapping := d.Get("packet_filters").([]interface{})
	ddosPacketFilter, pErr := setDdosPolicyPacketFilter(packetFilterMapping)
	if pErr != nil {
		return pErr
	}

	//set WaterPrintPolicy
	waterPrintMapping := d.Get("watermark_filters").([]interface{})
	waterPrintPolicy, _ := setWaterPrintPolicy(waterPrintMapping)

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	policyId := ""

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateDdosPolicy(ctx, resourceType, name, ddosPolicyDropOption, ddosPolicyPortLimit, ipBlackWhite, ddosPacketFilter, waterPrintPolicy)
		if e != nil {
			return retryError(e)
		}
		policyId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + FILED_SP + policyId)

	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy")
	}
	resourceType := items[0]
	policyId := items[1]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	ddosPolicy, has, err := dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ddosPolicy, has, err = dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("drop_options", flattenDdosDropOptionList([]*dayu.DDoSPolicyDropOption{ddosPolicy.DropOptions}))
	_ = d.Set("port_filters", flattenDdosPortLimitList(ddosPolicy.PortLimits))
	_ = d.Set("packet_filters", flattenDdosPacketFilterList(ddosPolicy.PacketFilters))
	blackIps, whiteIps := flattenIpBlackWhiteList(ddosPolicy.IpBlackWhiteLists)
	_ = d.Set("black_ips", blackIps)
	_ = d.Set("white_ips", whiteIps)
	_ = d.Set("watermark_filters", flattenWaterPrintPolicyList(ddosPolicy.WaterPrint))
	_ = d.Set("create_time", ddosPolicy.CreateTime)
	_ = d.Set("name", ddosPolicy.PolicyName)
	_ = d.Set("scene_id", ddosPolicy.SceneId)
	_ = d.Set("policy_id", ddosPolicy.PolicyId)
	_ = d.Set("watermark_key", flattenWaterPrintKeyList(ddosPolicy.WaterKey))

	return nil
}

func resourceTencentCloudDayuDdosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy")
	}
	resourceType := items[0]
	policyId := items[1]
	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	d.Partial(true)

	if d.HasChange("name") {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyDdosPolicyName(ctx, resourceType, policyId, d.Get("name").(string))
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("watermark_filters") || d.HasChange("ip_filters") || d.HasChange("packet_filters") || d.HasChange("port_filters") || d.HasChange("drop_options") {

		//set DDosPolicyDropOption
		dropMapping := d.Get("drop_options").([]interface{})
		ddosPolicyDropOption, _ := setDdosPolicyDropOption(dropMapping)

		//set DDoSPolicyPortLimit
		portMapping := d.Get("port_filters").([]interface{})
		ddosPolicyPortLimit, lErr := setDdosPolicyPortLimit(portMapping)

		if lErr != nil {
			return lErr
		}

		//set IpBlackWhite
		blackIps := d.Get("black_ips").(*schema.Set).List()
		whiteIps := d.Get("white_ips").(*schema.Set).List()
		ipBlackWhite, ipErr := setIpBlackWhite(blackIps, whiteIps)

		if ipErr != nil {
			return ipErr
		}
		//set DDoSPolicyPacketFilter
		packetFilterMapping := d.Get("packet_filters").([]interface{})
		ddosPacketFilter, pErr := setDdosPolicyPacketFilter(packetFilterMapping)
		if pErr != nil {
			return pErr
		}

		//set WaterPrintPolicy
		waterPrintMapping := d.Get("watermark_filters").([]interface{})
		waterPrintPolicy, _ := setWaterPrintPolicy(waterPrintMapping)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyDdosPolicy(ctx, resourceType, policyId, ddosPolicyDropOption, ddosPolicyPortLimit, ipBlackWhite, ddosPacketFilter, waterPrintPolicy)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("watermark_filters")
		d.SetPartial("white_ips")
		d.SetPartial("black_ips")
		d.SetPartial("drop_options")
		d.SetPartial("port_filters")
		d.SetPartial("packet_filters")
	}

	d.Partial(false)

	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy")
	}
	resourceType := items[0]
	policyId := items[1]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteDdosPolicy(ctx, resourceType, policyId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
			if err != nil {
				return retryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy fail, DDoS policy still exist from sdk DescribeDDosPolicy")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete DDoS policy fail, DDoS policy still exist from sdk DescribeDDosPolicy")
	}
}
