package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAntiddosPacketFilterConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosPacketFilterConfigCreate,
		Read:   resourceTencentCloudAntiddosPacketFilterConfigRead,
		Delete: resourceTencentCloudAntiddosPacketFilterConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "resource id.",
			},

			"packet_filter_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Feature filtering configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol, value [TCP udp icmp all].",
						},
						"sport_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Starting source port, ranging from 0 to 65535.",
						},
						"sport_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "End source port, values range from 1 to 65535, must be greater than or equal to the start source port.",
						},
						"dport_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Starting destination port, ranging from 0 to 65535.",
						},
						"dport_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "end destination port, ranging from 0 to 65535.",
						},
						"pktlen_min": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minimum message length, ranging from 1 to 1500.",
						},
						"pktlen_max": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The maximum message length, ranging from 1 to 1500, must be greater than or equal to the minimum message length.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action, value [drop (discard) transmit (release) drop_black (discard and pull black) drop_rst (intercept) drop_black_rst (intercept and pull black) forward (continue protection)].",
						},
						"match_begin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detection position, value [begin_l3 (IP header) begin_l4 (TCP/UDP header) begin_l5 (T payload) no_match (mismatch)].",
						},
						"match_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detection type, value [Sunday (keyword) pcre (regular expression)].",
						},
						"str": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detection value, key string or regular expression, value [When the detection type is Sunday, please fill in the string or hexadecimal bytecode, for example, x313233 corresponds to the hexadecimal word&gt;section code of the string &#39;123&#39;; when the detection type is pcre, please fill in the regular expression character string;].",
						},
						"depth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The detection depth starting from the detection position, with a value of [0-1500].",
						},
						"offset": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The offset from the detection position, with a value range of [0, Depth].",
						},
						"is_not": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to include detection values, with a value of [0 (inclusive) and 1 (exclusive)].",
						},
						"match_logic": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "When there is a second detection condition, the AND or relationship with the first detection condition, with the value [and (and relationship) none (fill in this value when there is no second detection condition)].",
						},
						"match_begin2": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Second detection position, value [begin_l5 (load) no_match (mismatch)].",
						},
						"match_type2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The second detection type, with a value of [Sunday (keyword) pcre (regular expression)].",
						},
						"str2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "key string or regular expression, value [When the detection type is Sunday, please fill in the string or hexadecimal bytecode, for example, x313233 corresponds to the hexadecimal word&gt;section code of the string &#39;123&#39;; when the detection type is pcre, please fill in the regular expression character string;].",
						},
						"depth2": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The second detection depth starting from the second detection position, with a value of [01500].",
						},
						"offset2": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The offset from the second detection position, with a value range of [0, Depth2].",
						},
						"is_not2": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether the second detection includes detection values, with a value of [0 (inclusive) and 1 (exclusive)].",
						},
						"pkt_len_gt": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Greater than message length, value 1+.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosPacketFilterConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_packet_filter_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = antiddos.NewCreatePacketFilterConfigRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	packetFilterConfig := antiddos.PacketFilterConfig{}
	if dMap, ok := helper.InterfacesHeadMap(d, "packet_filter_config"); ok {
		if v, ok := dMap["protocol"]; ok {
			packetFilterConfig.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["sport_start"]; ok {
			packetFilterConfig.SportStart = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sport_end"]; ok {
			packetFilterConfig.SportEnd = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["dport_start"]; ok {
			packetFilterConfig.DportStart = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["dport_end"]; ok {
			packetFilterConfig.DportEnd = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["pktlen_min"]; ok {
			packetFilterConfig.PktlenMin = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["pktlen_max"]; ok {
			packetFilterConfig.PktlenMax = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["action"]; ok {
			packetFilterConfig.Action = helper.String(v.(string))
		}
		if v, ok := dMap["match_begin"]; ok {
			packetFilterConfig.MatchBegin = helper.String(v.(string))
		}
		if v, ok := dMap["match_type"]; ok {
			packetFilterConfig.MatchType = helper.String(v.(string))
		}
		if v, ok := dMap["str"]; ok {
			packetFilterConfig.Str = helper.String(v.(string))
		}
		if v, ok := dMap["depth"]; ok {
			packetFilterConfig.Depth = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["offset"]; ok {
			packetFilterConfig.Offset = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["is_not"]; ok {
			packetFilterConfig.IsNot = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["match_logic"]; ok {
			packetFilterConfig.MatchLogic = helper.String(v.(string))
		}
		if v, ok := dMap["match_begin2"]; ok {
			packetFilterConfig.MatchBegin2 = helper.String(v.(string))
		}
		if v, ok := dMap["match_type2"]; ok {
			packetFilterConfig.MatchType2 = helper.String(v.(string))
		}
		if v, ok := dMap["str2"]; ok {
			packetFilterConfig.Str2 = helper.String(v.(string))
		}
		if v, ok := dMap["depth2"]; ok {
			packetFilterConfig.Depth2 = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["offset2"]; ok {
			packetFilterConfig.Offset2 = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["is_not2"]; ok {
			packetFilterConfig.IsNot2 = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["pkt_len_gt"]; ok {
			packetFilterConfig.PktLenGT = helper.IntInt64(v.(int))
		}
		request.PacketFilterConfig = &packetFilterConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreatePacketFilterConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos packetFilterConfig failed, reason:%+v", logId, err)
		return err
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	packetFilterConfigs, err := service.DescribeAntiddosPacketFilterConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	var configId *string

	for _, config := range packetFilterConfigs {
		if *config.InstanceDetailList[0].InstanceId != instanceId {
			continue
		}

		if *packetFilterConfig.Protocol != *config.PacketFilterConfig.Protocol {
			continue
		}
		if *packetFilterConfig.SportStart != *config.PacketFilterConfig.SportStart {
			continue
		}
		if *packetFilterConfig.SportEnd != *config.PacketFilterConfig.SportEnd {
			continue
		}
		if *packetFilterConfig.DportStart != *config.PacketFilterConfig.DportStart {
			continue
		}
		if *packetFilterConfig.DportEnd != *config.PacketFilterConfig.DportEnd {
			continue
		}
		if *packetFilterConfig.PktlenMin != *config.PacketFilterConfig.PktlenMin {
			continue
		}
		if *packetFilterConfig.PktlenMax != *config.PacketFilterConfig.PktlenMax {
			continue
		}
		if *packetFilterConfig.Action != *config.PacketFilterConfig.Action {
			continue
		}
		if *packetFilterConfig.MatchBegin != *config.PacketFilterConfig.MatchBegin {
			continue
		}
		if *packetFilterConfig.MatchType != *config.PacketFilterConfig.MatchType {
			continue
		}
		if *packetFilterConfig.Str != *config.PacketFilterConfig.Str {
			continue
		}
		if *packetFilterConfig.Depth != *config.PacketFilterConfig.Depth {
			continue
		}
		if *packetFilterConfig.Offset != *config.PacketFilterConfig.Offset {
			continue
		}
		if *packetFilterConfig.IsNot != *config.PacketFilterConfig.IsNot {
			continue
		}
		if *packetFilterConfig.MatchLogic != "" && *packetFilterConfig.MatchLogic != *config.PacketFilterConfig.MatchLogic {
			continue
		}
		if *packetFilterConfig.MatchBegin2 != "" && *packetFilterConfig.MatchBegin2 != *config.PacketFilterConfig.MatchBegin2 {
			continue
		}
		if *packetFilterConfig.MatchType2 != *config.PacketFilterConfig.MatchType2 {
			continue
		}
		if *packetFilterConfig.Str2 != *config.PacketFilterConfig.Str2 {
			continue
		}
		if *packetFilterConfig.Depth2 != *config.PacketFilterConfig.Depth2 {
			continue
		}
		if *packetFilterConfig.Offset2 != *config.PacketFilterConfig.Offset2 {
			continue
		}
		if *packetFilterConfig.IsNot2 != *config.PacketFilterConfig.IsNot2 {
			continue
		}
		if *packetFilterConfig.PktLenGT != *config.PacketFilterConfig.PktLenGT {
			continue
		}
		configId = config.PacketFilterConfig.Id
	}
	if configId != nil {
		d.SetId(instanceId + FILED_SP + *configId)
	} else {
		return fmt.Errorf("can not find config ids")
	}

	return resourceTencentCloudAntiddosPacketFilterConfigRead(d, meta)
}

func resourceTencentCloudAntiddosPacketFilterConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_packet_filter_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	configId := idSplit[1]
	packetFilterConfigs, err := service.DescribeAntiddosPacketFilterConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if len(packetFilterConfigs) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosPacketFilterConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	var packetFilterConfig *antiddos.PacketFilterRelation
	for _, config := range packetFilterConfigs {
		if *config.PacketFilterConfig.Id == configId {
			packetFilterConfig = config
		}
	}
	if packetFilterConfig != nil && packetFilterConfig.PacketFilterConfig != nil {
		packetFilterConfigMap := map[string]interface{}{}

		if packetFilterConfig.PacketFilterConfig.Protocol != nil {
			packetFilterConfigMap["protocol"] = packetFilterConfig.PacketFilterConfig.Protocol
		}

		if packetFilterConfig.PacketFilterConfig.SportStart != nil {
			packetFilterConfigMap["sport_start"] = packetFilterConfig.PacketFilterConfig.SportStart
		}

		if packetFilterConfig.PacketFilterConfig.SportEnd != nil {
			packetFilterConfigMap["sport_end"] = packetFilterConfig.PacketFilterConfig.SportEnd
		}

		if packetFilterConfig.PacketFilterConfig.DportStart != nil {
			packetFilterConfigMap["dport_start"] = packetFilterConfig.PacketFilterConfig.DportStart
		}

		if packetFilterConfig.PacketFilterConfig.DportEnd != nil {
			packetFilterConfigMap["dport_end"] = packetFilterConfig.PacketFilterConfig.DportEnd
		}

		if packetFilterConfig.PacketFilterConfig.PktlenMin != nil {
			packetFilterConfigMap["pktlen_min"] = packetFilterConfig.PacketFilterConfig.PktlenMin
		}

		if packetFilterConfig.PacketFilterConfig.PktlenMax != nil {
			packetFilterConfigMap["pktlen_max"] = packetFilterConfig.PacketFilterConfig.PktlenMax
		}

		if packetFilterConfig.PacketFilterConfig.Action != nil {
			packetFilterConfigMap["action"] = packetFilterConfig.PacketFilterConfig.Action
		}

		if packetFilterConfig.PacketFilterConfig.MatchBegin != nil {
			packetFilterConfigMap["match_begin"] = packetFilterConfig.PacketFilterConfig.MatchBegin
		}

		if packetFilterConfig.PacketFilterConfig.MatchType != nil {
			packetFilterConfigMap["match_type"] = packetFilterConfig.PacketFilterConfig.MatchType
		}

		if packetFilterConfig.PacketFilterConfig.Str != nil {
			packetFilterConfigMap["str"] = packetFilterConfig.PacketFilterConfig.Str
		}

		if packetFilterConfig.PacketFilterConfig.Depth != nil {
			packetFilterConfigMap["depth"] = packetFilterConfig.PacketFilterConfig.Depth
		}

		if packetFilterConfig.PacketFilterConfig.Offset != nil {
			packetFilterConfigMap["offset"] = packetFilterConfig.PacketFilterConfig.Offset
		}

		if packetFilterConfig.PacketFilterConfig.IsNot != nil {
			packetFilterConfigMap["is_not"] = packetFilterConfig.PacketFilterConfig.IsNot
		}

		if packetFilterConfig.PacketFilterConfig.MatchLogic != nil {
			packetFilterConfigMap["match_logic"] = packetFilterConfig.PacketFilterConfig.MatchLogic
		}

		if packetFilterConfig.PacketFilterConfig.MatchBegin2 != nil {
			packetFilterConfigMap["match_begin2"] = packetFilterConfig.PacketFilterConfig.MatchBegin2
		}

		if packetFilterConfig.PacketFilterConfig.MatchType2 != nil {
			packetFilterConfigMap["match_type2"] = packetFilterConfig.PacketFilterConfig.MatchType2
		}

		if packetFilterConfig.PacketFilterConfig.Str2 != nil {
			packetFilterConfigMap["str2"] = packetFilterConfig.PacketFilterConfig.Str2
		}

		if packetFilterConfig.PacketFilterConfig.Depth2 != nil {
			packetFilterConfigMap["depth2"] = packetFilterConfig.PacketFilterConfig.Depth2
		}

		if packetFilterConfig.PacketFilterConfig.Offset2 != nil {
			packetFilterConfigMap["offset2"] = packetFilterConfig.PacketFilterConfig.Offset2
		}

		if packetFilterConfig.PacketFilterConfig.IsNot2 != nil {
			packetFilterConfigMap["is_not2"] = packetFilterConfig.PacketFilterConfig.IsNot2
		}

		if packetFilterConfig.PacketFilterConfig.PktLenGT != nil {
			packetFilterConfigMap["pkt_len_gt"] = packetFilterConfig.PacketFilterConfig.PktLenGT
		}

		_ = d.Set("packet_filter_config", []interface{}{packetFilterConfigMap})
	}

	return nil
}

func resourceTencentCloudAntiddosPacketFilterConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_packet_filter_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	configId := idSplit[1]
	packetFilterConfigs, err := service.DescribeAntiddosPacketFilterConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	var packetFilterConfig *antiddos.PacketFilterRelation
	for _, config := range packetFilterConfigs {
		if *config.PacketFilterConfig.Id == configId {
			packetFilterConfig = config
		}
	}

	if packetFilterConfig != nil {
		if err := service.DeleteAntiddosPacketFilterConfigById(ctx, instanceId, packetFilterConfig.PacketFilterConfig); err != nil {
			return err
		}
	}

	return nil
}
