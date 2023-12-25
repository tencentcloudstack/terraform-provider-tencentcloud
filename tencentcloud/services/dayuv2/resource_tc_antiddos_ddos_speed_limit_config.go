package dayuv2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosDdosSpeedLimitConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosDdosSpeedLimitConfigCreate,
		Read:   resourceTencentCloudAntiddosDdosSpeedLimitConfigRead,
		Update: resourceTencentCloudAntiddosDdosSpeedLimitConfigUpdate,
		Delete: resourceTencentCloudAntiddosDdosSpeedLimitConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"ddos_speed_limit_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Accessing speed limit configuration, the configuration ID cannot be empty when filling in parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Speed limit mode, value [1 (based on source IP speed limit) 2 (based on destination port speed limit)].",
						},
						"speed_values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Speed limit values, each type of speed limit value can support up to 1; This field array has at least one speed limit value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Speed limit value type, value [1 (packet rate pps) 2 (bandwidth bps)].",
									},
									"value": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "value.",
									},
								},
							},
						},
						"dst_port_scopes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "This field has been deprecated. Please fill in the new field DstPortList.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Starting port, ranging from 1 to 65535.",
									},
									"end_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "end  port, ranging from 1 to 65535.",
									},
								},
							},
						},
						"protocol_list": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP protocol numbers, values [ALL (all protocols) TCP (tcp protocol) UDP (udp protocol) SMP (smp protocol) 1; 2-100 (custom protocol number range, up to 8)] Note: When customizing the protocol number range, only the protocol number can be filled in, multiple ranges; Separation; When filling in ALL, no other agreements or agreements can be filled inNumber.",
						},
						"dst_port_list": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "List of port ranges, up to 8, multiple; Separate and indicate the range with -; This port range must be filled in; Fill in style 1:0-65535, style 2: 80; 443; 1000-2000.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosDdosSpeedLimitConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_speed_limit_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request      = antiddos.NewCreateDDoSSpeedLimitConfigRequest()
		instanceId   string
		dstPortList  string
		protocolList string
		mode         uint64
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	speedValues := make([]*antiddos.SpeedValue, 0)
	dstPortScopes := make([]*antiddos.PortSegment, 0)
	if dMap, ok := helper.InterfacesHeadMap(d, "ddos_speed_limit_config"); ok {
		dDoSSpeedLimitConfig := antiddos.DDoSSpeedLimitConfig{}
		if v, ok := dMap["mode"]; ok {
			mode = uint64(v.(int))
			dDoSSpeedLimitConfig.Mode = &mode
		}
		if v, ok := dMap["speed_values"]; ok {
			for _, item := range v.([]interface{}) {
				speedValuesMap := item.(map[string]interface{})
				speedValue := antiddos.SpeedValue{}
				if v, ok := speedValuesMap["type"]; ok {
					speedValue.Type = helper.IntUint64(v.(int))
				}
				if v, ok := speedValuesMap["value"]; ok {
					speedValue.Value = helper.IntUint64(v.(int))
				}
				speedValues = append(speedValues, &speedValue)
			}
			dDoSSpeedLimitConfig.SpeedValues = speedValues
		}
		if v, ok := dMap["dst_port_scopes"]; ok {
			for _, item := range v.([]interface{}) {
				dstPortScopesMap := item.(map[string]interface{})
				portSegment := antiddos.PortSegment{}
				if v, ok := dstPortScopesMap["begin_port"]; ok {
					portSegment.BeginPort = helper.IntUint64(v.(int))
				}
				if v, ok := dstPortScopesMap["end_port"]; ok {
					portSegment.EndPort = helper.IntUint64(v.(int))
				}
				dstPortScopes = append(dstPortScopes, &portSegment)
			}
			dDoSSpeedLimitConfig.DstPortScopes = dstPortScopes
		}
		if v, ok := dMap["protocol_list"]; ok {
			protocolList = v.(string)
			dDoSSpeedLimitConfig.ProtocolList = helper.String(protocolList)
		}
		if v, ok := dMap["dst_port_list"]; ok {
			dstPortList = v.(string)
			dDoSSpeedLimitConfig.DstPortList = helper.String(dstPortList)
		}
		request.DDoSSpeedLimitConfig = &dDoSSpeedLimitConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateDDoSSpeedLimitConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos ddosSpeedLimitConfig failed, reason:%+v", logId, err)
		return err
	}

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	configList, err := service.DescribeAntiddosDdosSpeedLimitConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDdosSpeedLimitConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var targetConfig *antiddos.DDoSSpeedLimitConfig
	for _, item := range configList {
		if *item.SpeedLimitConfig.Mode != mode {
			log.Println("==>SpeedLimitConfig.Mode")
			continue
		}
		if *item.SpeedLimitConfig.ProtocolList != protocolList {
			log.Println("==>SpeedLimitConfig.ProtocolList", *item.SpeedLimitConfig.ProtocolList, protocolList, "===")
			continue
		}
		if *item.SpeedLimitConfig.DstPortList != dstPortList {
			log.Println("==>SpeedLimitConfig.DstPortList")
			continue
		}
		tmpSpeedValues := item.SpeedLimitConfig.SpeedValues
		sort.Slice(speedValues, func(i, j int) bool {
			return *speedValues[i].Type < *speedValues[j].Type
		})
		sort.Slice(tmpSpeedValues, func(i, j int) bool {
			return *tmpSpeedValues[i].Type < *tmpSpeedValues[j].Type
		})
		if !reflect.DeepEqual(speedValues, tmpSpeedValues) {
			log.Println("==>SpeedLimitConfig.SpeedValues")
			continue
		}
		targetConfig = item.SpeedLimitConfig
	}

	if targetConfig == nil {
		return fmt.Errorf("can not find speed limit config")
	}

	d.SetId(instanceId + tccommon.FILED_SP + *targetConfig.Id)

	return resourceTencentCloudAntiddosDdosSpeedLimitConfigRead(d, meta)
}

func resourceTencentCloudAntiddosDdosSpeedLimitConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_speed_limit_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := idSplit[0]
	configId := idSplit[1]

	configList, err := service.DescribeAntiddosDdosSpeedLimitConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDdosSpeedLimitConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var targetConfig *antiddos.DDoSSpeedLimitConfig
	for _, item := range configList {
		if *item.SpeedLimitConfig.Id == configId {
			targetConfig = item.SpeedLimitConfig
			break
		}
	}

	_ = d.Set("instance_id", instanceId)

	if targetConfig != nil {
		configMap := map[string]interface{}{}

		if targetConfig.Mode != nil {
			configMap["mode"] = targetConfig.Mode
		}

		if targetConfig.SpeedValues != nil {
			speedValuesList := []interface{}{}
			for _, speedValues := range targetConfig.SpeedValues {
				speedValuesMap := map[string]interface{}{}

				if speedValues.Type != nil {
					speedValuesMap["type"] = speedValues.Type
				}

				if speedValues.Value != nil {
					speedValuesMap["value"] = speedValues.Value
				}

				speedValuesList = append(speedValuesList, speedValuesMap)
			}

			configMap["speed_values"] = speedValuesList
		}

		if targetConfig.DstPortScopes != nil {
			dstPortScopesList := []interface{}{}
			for _, dstPortScopes := range targetConfig.DstPortScopes {
				dstPortScopesMap := map[string]interface{}{}

				if dstPortScopes.BeginPort != nil {
					dstPortScopesMap["begin_port"] = dstPortScopes.BeginPort
				}

				if dstPortScopes.EndPort != nil {
					dstPortScopesMap["end_port"] = dstPortScopes.EndPort
				}

				dstPortScopesList = append(dstPortScopesList, dstPortScopesMap)
			}

			configMap["dst_port_scopes"] = dstPortScopesList
		}

		if targetConfig.ProtocolList != nil {
			configMap["protocol_list"] = targetConfig.ProtocolList
		}

		if targetConfig.DstPortList != nil {
			configMap["dst_port_list"] = targetConfig.DstPortList
		}

		_ = d.Set("ddos_speed_limit_config", []interface{}{configMap})
	}

	return nil
}

func resourceTencentCloudAntiddosDdosSpeedLimitConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_speed_limit_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewModifyDDoSSpeedLimitConfigRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	request.InstanceId = helper.String(idSplit[0])

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("ddos_speed_limit_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ddos_speed_limit_config"); ok {
			dDoSSpeedLimitConfig := antiddos.DDoSSpeedLimitConfig{}
			if v, ok := dMap["mode"]; ok {
				dDoSSpeedLimitConfig.Mode = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["speed_values"]; ok {
				for _, item := range v.([]interface{}) {
					speedValuesMap := item.(map[string]interface{})
					speedValue := antiddos.SpeedValue{}
					if v, ok := speedValuesMap["type"]; ok {
						speedValue.Type = helper.IntUint64(v.(int))
					}
					if v, ok := speedValuesMap["value"]; ok {
						speedValue.Value = helper.IntUint64(v.(int))
					}
					dDoSSpeedLimitConfig.SpeedValues = append(dDoSSpeedLimitConfig.SpeedValues, &speedValue)
				}
			}
			if v, ok := dMap["dst_port_scopes"]; ok {
				for _, item := range v.([]interface{}) {
					dstPortScopesMap := item.(map[string]interface{})
					portSegment := antiddos.PortSegment{}
					if v, ok := dstPortScopesMap["begin_port"]; ok {
						portSegment.BeginPort = helper.IntUint64(v.(int))
					}
					if v, ok := dstPortScopesMap["end_port"]; ok {
						portSegment.EndPort = helper.IntUint64(v.(int))
					}
					dDoSSpeedLimitConfig.DstPortScopes = append(dDoSSpeedLimitConfig.DstPortScopes, &portSegment)
				}
			}
			dDoSSpeedLimitConfig.Id = helper.String(idSplit[1])
			if v, ok := dMap["protocol_list"]; ok {
				dDoSSpeedLimitConfig.ProtocolList = helper.String(v.(string))
			}
			if v, ok := dMap["dst_port_list"]; ok {
				dDoSSpeedLimitConfig.DstPortList = helper.String(v.(string))
			}
			request.DDoSSpeedLimitConfig = &dDoSSpeedLimitConfig
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().ModifyDDoSSpeedLimitConfig(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos ddosSpeedLimitConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosDdosSpeedLimitConfigRead(d, meta)
}

func resourceTencentCloudAntiddosDdosSpeedLimitConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_speed_limit_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	configId := idSplit[1]

	configList, err := service.DescribeAntiddosDdosSpeedLimitConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDdosSpeedLimitConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var targetConfig *antiddos.DDoSSpeedLimitConfig
	for _, item := range configList {
		if *item.SpeedLimitConfig.Id == configId {
			targetConfig = item.SpeedLimitConfig
			break
		}
	}

	if targetConfig != nil {
		if err := service.DeleteAntiddosDdosSpeedLimitConfigById(ctx, instanceId, targetConfig); err != nil {
			return err
		}
	}

	return nil
}
