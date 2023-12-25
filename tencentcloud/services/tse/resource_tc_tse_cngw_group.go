package tse

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTseCngwGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwGroupCreate,
		Read:   resourceTencentCloudTseCngwGroupRead,
		Update: resourceTencentCloudTseCngwGroupUpdate,
		Delete: resourceTencentCloudTseCngwGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway IDonly postpaid gateway supported.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway group name.",
			},

			"node_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "group node configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "group specification, 1c2g|2c4g|4c8g|8c16g.",
						},
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "group node number, 2-50.",
						},
					},
				},
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value:- subnet-ahde9me9.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "description information of group.",
			},

			"internet_max_bandwidth_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "public network outbound traffic bandwidth,[1,2048]Mbps.",
			},

			"internet_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "internet configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_address_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "internet type. Reference value:- IPV4 (default value)- IPV6.",
						},
						"internet_pay_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "trade type of internet. Reference value:- BANDWIDTH- TRAFFIC (default value).",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "public network bandwidth.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description of clb.",
						},
						"sla_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "specification type of clb. Default shared type when this parameter is empty. Reference value:- SLA LCU-supported.",
						},
						"multi_zone_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether load balancing has multiple availability zones.",
						},
						"master_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "primary availability zone.",
						},
						"slave_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "alternate availability zone.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request   = tse.NewCreateNativeGatewayServerGroupRequest()
		response  = tse.NewCreateNativeGatewayServerGroupResponse()
		gatewayId string
		groupId   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		cloudNativeAPIGatewayNodeConfig := tse.CloudNativeAPIGatewayNodeConfig{}
		if v, ok := dMap["specification"]; ok {
			cloudNativeAPIGatewayNodeConfig.Specification = helper.String(v.(string))
		}
		if v, ok := dMap["number"]; ok {
			cloudNativeAPIGatewayNodeConfig.Number = helper.IntInt64(v.(int))
		}
		request.NodeConfig = &cloudNativeAPIGatewayNodeConfig
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_config"); ok {
		internetConfig := tse.InternetConfig{}
		if v, ok := dMap["internet_address_version"]; ok {
			internetConfig.InternetAddressVersion = helper.String(v.(string))
		}
		if v, ok := dMap["internet_pay_mode"]; ok {
			internetConfig.InternetPayMode = helper.String(v.(string))
		}
		if v, ok := dMap["internet_max_bandwidth_out"]; ok {
			internetConfig.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["description"]; ok {
			internetConfig.Description = helper.String(v.(string))
		}
		if v, ok := dMap["sla_type"]; ok {
			internetConfig.SlaType = helper.String(v.(string))
		}
		if v, ok := dMap["multi_zone_flag"]; ok {
			internetConfig.MultiZoneFlag = helper.Bool(v.(bool))
		}
		if v, ok := dMap["master_zone_id"]; ok {
			internetConfig.MasterZoneId = helper.String(v.(string))
		}
		if v, ok := dMap["slave_zone_id"]; ok {
			internetConfig.SlaveZoneId = helper.String(v.(string))
		}
		request.InternetConfig = &internetConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CreateNativeGatewayServerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.Result.GroupId
	d.SetId(gatewayId + tccommon.FILED_SP + groupId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.CheckTseNativeAPIGatewayGroupStatusById(ctx, gatewayId, groupId, "create"); err != nil {
		return err
	}

	return resourceTencentCloudTseCngwGroupRead(d, meta)
}

func resourceTencentCloudTseCngwGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	cngwGroup, err := service.DescribeTseCngwGroupById(ctx, gatewayId, groupId)
	if err != nil {
		return err
	}

	if cngwGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwGroup.GatewayId != nil {
		_ = d.Set("gateway_id", cngwGroup.GatewayId)
	}

	if cngwGroup.Name != nil {
		_ = d.Set("name", cngwGroup.Name)
	}

	if cngwGroup.NodeConfig != nil {
		nodeConfigMap := map[string]interface{}{}

		if cngwGroup.NodeConfig.Specification != nil {
			nodeConfigMap["specification"] = cngwGroup.NodeConfig.Specification
		}

		if cngwGroup.NodeConfig.Number != nil {
			nodeConfigMap["number"] = cngwGroup.NodeConfig.Number
		}

		_ = d.Set("node_config", []interface{}{nodeConfigMap})
	}

	if cngwGroup.SubnetIds != nil {
		_ = d.Set("subnet_id", cngwGroup.SubnetIds)
	}

	if cngwGroup.Description != nil {
		_ = d.Set("description", cngwGroup.Description)
	}

	if cngwGroup.InternetMaxBandwidthOut != nil {
		_ = d.Set("internet_max_bandwidth_out", cngwGroup.InternetMaxBandwidthOut)
	}

	// if cngwGroup.InternetConfig != nil {
	// 	internetConfigMap := map[string]interface{}{}

	// 	if cngwGroup.InternetConfig.InternetAddressVersion != nil {
	// 		internetConfigMap["internet_address_version"] = cngwGroup.InternetConfig.InternetAddressVersion
	// 	}

	// 	if cngwGroup.InternetConfig.InternetPayMode != nil {
	// 		internetConfigMap["internet_pay_mode"] = cngwGroup.InternetConfig.InternetPayMode
	// 	}

	// 	if cngwGroup.InternetConfig.InternetMaxBandwidthOut != nil {
	// 		internetConfigMap["internet_max_bandwidth_out"] = cngwGroup.InternetConfig.InternetMaxBandwidthOut
	// 	}

	// 	if cngwGroup.InternetConfig.Description != nil {
	// 		internetConfigMap["description"] = cngwGroup.InternetConfig.Description
	// 	}

	// 	if cngwGroup.InternetConfig.SlaType != nil {
	// 		internetConfigMap["sla_type"] = cngwGroup.InternetConfig.SlaType
	// 	}

	// 	if cngwGroup.InternetConfig.MultiZoneFlag != nil {
	// 		internetConfigMap["multi_zone_flag"] = cngwGroup.InternetConfig.MultiZoneFlag
	// 	}

	// 	if cngwGroup.InternetConfig.MasterZoneId != nil {
	// 		internetConfigMap["master_zone_id"] = cngwGroup.InternetConfig.MasterZoneId
	// 	}

	// 	if cngwGroup.InternetConfig.SlaveZoneId != nil {
	// 		internetConfigMap["slave_zone_id"] = cngwGroup.InternetConfig.SlaveZoneId
	// 	}

	// 	_ = d.Set("internet_config", []interface{}{internetConfigMap})
	// }

	return nil
}

func resourceTencentCloudTseCngwGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tse.NewModifyNativeGatewayServerGroupRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	request.GatewayId = &gatewayId
	request.GroupId = &groupId

	immutableArgs := []string{"gateway_id", "node_config", "subnet_id", "internet_max_bandwidth_out", "internet_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().ModifyNativeGatewayServerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwGroup failed, reason:%+v", logId, err)
		return err
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.CheckTseNativeAPIGatewayGroupStatusById(ctx, gatewayId, groupId, "update"); err != nil {
		return err
	}

	return resourceTencentCloudTseCngwGroupRead(d, meta)
}

func resourceTencentCloudTseCngwGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	if err := service.DeleteTseCngwGroupById(ctx, gatewayId, groupId); err != nil {
		return err
	}
	if err := service.CheckTseNativeAPIGatewayGroupStatusById(ctx, gatewayId, groupId, "delete"); err != nil {
		return err
	}

	return nil
}
