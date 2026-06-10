package cfw

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwClusterNatFwSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwClusterNatFwSwitchCreate,
		Read:   resourceTencentCloudCfwClusterNatFwSwitchRead,
		Update: resourceTencentCloudCfwClusterNatFwSwitchUpdate,
		Delete: resourceTencentCloudCfwClusterNatFwSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_ccn_switch": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "NAT CCN switch configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ins_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "NAT firewall instance ID.",
						},
						"ccn_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "CCN instance ID.",
						},
						"switch_mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Switch access mode, 1: automatic access, 2: manual access.",
						},
						"routing_mode": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Traffic steering routing method, 0: multi-route table, 1: policy routing. Automatic access mode only supports policy routing (1); manual access mode supports both multi-route table (0) and policy routing (1).",
						},
						"access_instance_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of access instances.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Instance ID.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Instance type such as VPC or DIRECTCONNECT.",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Region where the instance is located.",
									},
									"access_cidr_mode": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.",
									},
									"access_cidr_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of network segments for accessing firewall.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"lead_vpc_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CIDR of the lead VPC.",
						},
					},
				},
			},

			// computed
			"switch_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Switch access mode, 1: automatic access, 2: manual access.",
			},
			"routing_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Traffic steering routing method, 0: multi-route table, 1: policy routing.",
			},
			"bypass": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Bypass status, 0: disabled, 1: enabled.",
			},
			"ccn_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CCN instance ID.",
			},
		},
	}
}

func resourceTencentCloudCfwClusterNatFwSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_nat_fw_switch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cfwv20190904.NewOpenClusterNatFwSwitchRequest()
		natInsId string
		ccnId    string
	)

	if v, ok := d.GetOk("nat_ccn_switch"); ok {
		natCcnSwitchList := v.([]interface{})
		if len(natCcnSwitchList) > 0 {
			natCcnSwitchMap := natCcnSwitchList[0].(map[string]interface{})
			natCcnSwitch := cfwv20190904.NatCcnSwitchConfig{}

			if v, ok := natCcnSwitchMap["nat_ins_id"].(string); ok && v != "" {
				natCcnSwitch.NatInsId = helper.String(v)
				natInsId = v
			}

			if v, ok := natCcnSwitchMap["ccn_id"].(string); ok && v != "" {
				natCcnSwitch.CcnId = helper.String(v)
				ccnId = v
			}

			if v, ok := natCcnSwitchMap["switch_mode"].(int); ok {
				natCcnSwitch.SwitchMode = helper.IntInt64(v)
			}

			if v, ok := natCcnSwitchMap["routing_mode"].(int); ok {
				natCcnSwitch.RoutingMode = helper.IntInt64(v)
			}

			if v, ok := natCcnSwitchMap["access_instance_list"]; ok {
				for _, item := range v.([]interface{}) {
					accessInstanceMap := item.(map[string]interface{})
					accessInstance := cfwv20190904.AccessInstanceInfo{}
					if v, ok := accessInstanceMap["instance_id"].(string); ok && v != "" {
						accessInstance.InstanceId = helper.String(v)
					}

					if v, ok := accessInstanceMap["instance_type"].(string); ok && v != "" {
						accessInstance.InstanceType = helper.String(v)
					}

					if v, ok := accessInstanceMap["instance_region"].(string); ok && v != "" {
						accessInstance.InstanceRegion = helper.String(v)
					}

					if v, ok := accessInstanceMap["access_cidr_mode"].(int); ok {
						accessInstance.AccessCidrMode = helper.IntInt64(v)
					}

					if v, ok := accessInstanceMap["access_cidr_list"]; ok {
						for _, cidr := range v.([]interface{}) {
							accessInstance.AccessCidrList = append(accessInstance.AccessCidrList, helper.String(cidr.(string)))
						}
					}

					natCcnSwitch.AccessInstanceList = append(natCcnSwitch.AccessInstanceList, &accessInstance)
				}
			}

			if v, ok := natCcnSwitchMap["lead_vpc_cidr"].(string); ok && v != "" {
				natCcnSwitch.LeadVpcCidr = helper.String(v)
			}

			request.NatCcnSwitch = &natCcnSwitch
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().OpenClusterNatFwSwitchWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("open cfw cluster_nat_fw_switch failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cfw cluster_nat_fw_switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create cfw cluster_nat_fw_switch, nat_ins_id=%s, ccn_id=%s", logId, natInsId, ccnId)
	d.SetId(strings.Join([]string{natInsId, ccnId}, tccommon.FILED_SP))
	return resourceTencentCloudCfwClusterNatFwSwitchRead(d, meta)
}

func resourceTencentCloudCfwClusterNatFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_nat_fw_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	natInsId := idSplit[0]
	ccnId := idSplit[1]

	respData, err := service.DescribeNatCcnFwSwitchById(ctx, natInsId, ccnId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_cluster_nat_fw_switch` nat_ins_id=%s ccn_id=%s not found, please check if it has been deleted.\n", logId, natInsId, ccnId)
		d.SetId("")
		return nil
	}

	if respData.SwitchMode != nil {
		_ = d.Set("switch_mode", respData.SwitchMode)
	}

	if respData.RoutingMode != nil {
		_ = d.Set("routing_mode", respData.RoutingMode)
	}

	if respData.Bypass != nil {
		_ = d.Set("bypass", respData.Bypass)
	}

	if respData.CcnId != nil {
		_ = d.Set("ccn_id", respData.CcnId)
	}

	// Update nat_ccn_switch block with current values
	if v, ok := d.GetOk("nat_ccn_switch"); ok {
		natCcnSwitchList := v.([]interface{})
		if len(natCcnSwitchList) > 0 {
			natCcnSwitchMap := natCcnSwitchList[0].(map[string]interface{})
			natCcnSwitchMap["nat_ins_id"] = natInsId
			natCcnSwitchMap["ccn_id"] = ccnId
			if respData.SwitchMode != nil {
				natCcnSwitchMap["switch_mode"] = int(*respData.SwitchMode)
			}

			if respData.RoutingMode != nil {
				natCcnSwitchMap["routing_mode"] = int(*respData.RoutingMode)
			}

			if respData.AccessInstanceList != nil {
				accessInstanceListResult := make([]map[string]interface{}, 0, len(respData.AccessInstanceList))
				for _, accessInstance := range respData.AccessInstanceList {
					accessInstanceMap := map[string]interface{}{}
					if accessInstance.InstanceId != nil {
						accessInstanceMap["instance_id"] = accessInstance.InstanceId
					}

					if accessInstance.InstanceType != nil {
						accessInstanceMap["instance_type"] = accessInstance.InstanceType
					}

					if accessInstance.InstanceRegion != nil {
						accessInstanceMap["instance_region"] = accessInstance.InstanceRegion
					}

					if accessInstance.AccessCidrMode != nil {
						accessInstanceMap["access_cidr_mode"] = accessInstance.AccessCidrMode
					}

					if accessInstance.AccessCidrList != nil {
						accessInstanceMap["access_cidr_list"] = accessInstance.AccessCidrList
					}

					accessInstanceListResult = append(accessInstanceListResult, accessInstanceMap)
				}

				natCcnSwitchMap["access_instance_list"] = accessInstanceListResult
			}

			_ = d.Set("nat_ccn_switch", []interface{}{natCcnSwitchMap})
		}
	} else {
		natCcnSwitchMap := map[string]interface{}{
			"nat_ins_id": natInsId,
			"ccn_id":     ccnId,
		}
		if respData.SwitchMode != nil {
			natCcnSwitchMap["switch_mode"] = int(*respData.SwitchMode)
		}

		if respData.RoutingMode != nil {
			natCcnSwitchMap["routing_mode"] = int(*respData.RoutingMode)
		}

		if respData.AccessInstanceList != nil {
			accessInstanceListResult := make([]map[string]interface{}, 0, len(respData.AccessInstanceList))
			for _, accessInstance := range respData.AccessInstanceList {
				accessInstanceMap := map[string]interface{}{}
				if accessInstance.InstanceId != nil {
					accessInstanceMap["instance_id"] = accessInstance.InstanceId
				}

				if accessInstance.InstanceType != nil {
					accessInstanceMap["instance_type"] = accessInstance.InstanceType
				}

				if accessInstance.InstanceRegion != nil {
					accessInstanceMap["instance_region"] = accessInstance.InstanceRegion
				}

				if accessInstance.AccessCidrMode != nil {
					accessInstanceMap["access_cidr_mode"] = accessInstance.AccessCidrMode
				}

				if accessInstance.AccessCidrList != nil {
					accessInstanceMap["access_cidr_list"] = accessInstance.AccessCidrList
				}

				accessInstanceListResult = append(accessInstanceListResult, accessInstanceMap)
			}

			natCcnSwitchMap["access_instance_list"] = accessInstanceListResult
		}

		_ = d.Set("nat_ccn_switch", []interface{}{natCcnSwitchMap})
	}

	return nil
}

func resourceTencentCloudCfwClusterNatFwSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_nat_fw_switch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	natInsId := idSplit[0]
	ccnId := idSplit[1]

	if d.HasChange("nat_ccn_switch") {
		request := cfwv20190904.NewModifyClusterNatFwSwitchRequest()
		natCcnSwitch := cfwv20190904.NatCcnSwitchConfig{}
		natCcnSwitch.NatInsId = &natInsId
		natCcnSwitch.CcnId = &ccnId

		if v, ok := d.GetOk("nat_ccn_switch"); ok {
			natCcnSwitchList := v.([]interface{})
			if len(natCcnSwitchList) > 0 {
				natCcnSwitchMap := natCcnSwitchList[0].(map[string]interface{})

				if v, ok := natCcnSwitchMap["switch_mode"].(int); ok {
					natCcnSwitch.SwitchMode = helper.IntInt64(v)
				}

				if v, ok := natCcnSwitchMap["routing_mode"].(int); ok {
					natCcnSwitch.RoutingMode = helper.IntInt64(v)
				}

				if v, ok := natCcnSwitchMap["access_instance_list"]; ok {
					for _, item := range v.([]interface{}) {
						accessInstanceMap := item.(map[string]interface{})
						accessInstance := cfwv20190904.AccessInstanceInfo{}
						if v, ok := accessInstanceMap["instance_id"].(string); ok && v != "" {
							accessInstance.InstanceId = helper.String(v)
						}

						if v, ok := accessInstanceMap["instance_type"].(string); ok && v != "" {
							accessInstance.InstanceType = helper.String(v)
						}

						if v, ok := accessInstanceMap["instance_region"].(string); ok && v != "" {
							accessInstance.InstanceRegion = helper.String(v)
						}

						if v, ok := accessInstanceMap["access_cidr_mode"].(int); ok {
							accessInstance.AccessCidrMode = helper.IntInt64(v)
						}

						if v, ok := accessInstanceMap["access_cidr_list"]; ok {
							for _, cidr := range v.([]interface{}) {
								accessInstance.AccessCidrList = append(accessInstance.AccessCidrList, helper.String(cidr.(string)))
							}
						}

						natCcnSwitch.AccessInstanceList = append(natCcnSwitch.AccessInstanceList, &accessInstance)
					}
				}

				if v, ok := natCcnSwitchMap["lead_vpc_cidr"].(string); ok && v != "" {
					natCcnSwitch.LeadVpcCidr = helper.String(v)
				}
			}
		}

		request.NatCcnSwitch = &natCcnSwitch
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyClusterNatFwSwitchWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("modify cfw cluster_nat_fw_switch failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cfw cluster_nat_fw_switch failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudCfwClusterNatFwSwitchRead(d, meta)
}

func resourceTencentCloudCfwClusterNatFwSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_nat_fw_switch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewCloseClusterNatFwSwitchRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	natInsId := idSplit[0]
	ccnId := idSplit[1]

	request.NatInsId = &natInsId
	request.CcnId = &ccnId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().CloseClusterNatFwSwitchWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("close cfw cluster_nat_fw_switch failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cfw cluster_nat_fw_switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
