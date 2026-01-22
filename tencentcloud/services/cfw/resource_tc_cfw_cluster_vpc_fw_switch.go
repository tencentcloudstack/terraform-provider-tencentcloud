package cfw

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwClusterVpcFwSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwClusterVpcFwSwitchCreate,
		Read:   resourceTencentCloudCfwClusterVpcFwSwitchRead,
		Update: resourceTencentCloudCfwClusterVpcFwSwitchUpdate,
		Delete: resourceTencentCloudCfwClusterVpcFwSwitchDelete,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CCN ID.",
			},
			"switch_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Switch access mode, 1: automatic access, 2: manual access.",
			},
			"routing_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Traffic steering routing method, 0: multi-route table, 1: policy routing.",
			},
			"region_cidr_configs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Regional level CIDR configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Traffic steering region.",
						},
						"cidr_mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "CIDR mode: 0-skip, 1-automatic, 2-custom.",
						},
						"custom_cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom CIDR (required when CidrMode=2), empty string otherwise.",
						},
					},
				},
			},
			"interconnect_pairs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Interconnect pair list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_a": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Group A.",
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
										Type:        schema.TypeSet,
										Required:    true,
										Description: "List of network segments for accessing firewall.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"group_b": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Group B.",
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
										Type:        schema.TypeSet,
										Required:    true,
										Description: "List of network segments for accessing firewall.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"interconnect_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interconnect mode: `CrossConnect`: cross interconnect (each instance in group A interconnects with each instance in group B), `FullMesh`: full mesh (group A content is identical to group B, equivalent to pairwise interconnection within the group).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCfwClusterVpcFwSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_vpc_fw_switch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewModifyClusterVpcFwSwitchRequest()
		ccnId   string
	)

	ccnSwitchInfo := cfwv20190904.CcnSwitchInfo{}
	if v, ok := d.GetOk("ccn_id"); ok {
		ccnSwitchInfo.CcnId = helper.String(v.(string))
		ccnId = v.(string)
	}

	if v, ok := d.GetOk("switch_mode"); ok {
		ccnSwitchInfo.SwitchMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("routing_mode"); ok {
		ccnSwitchInfo.RoutingMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("region_cidr_configs"); ok {
		for _, item := range v.([]interface{}) {
			regionCidrConfigsMap := item.(map[string]interface{})
			regionCidrConfig := cfwv20190904.RegionCidrConfig{}
			if v, ok := regionCidrConfigsMap["region"].(string); ok && v != "" {
				regionCidrConfig.Region = helper.String(v)
			}

			if v, ok := regionCidrConfigsMap["cidr_mode"].(int); ok {
				regionCidrConfig.CidrMode = helper.IntInt64(v)
			}

			if v, ok := regionCidrConfigsMap["custom_cidr"].(string); ok && v != "" {
				regionCidrConfig.CustomCidr = helper.String(v)
			}

			ccnSwitchInfo.RegionCidrConfigs = append(ccnSwitchInfo.RegionCidrConfigs, &regionCidrConfig)
		}
	}

	if v, ok := d.GetOk("interconnect_pairs"); ok {
		for _, item := range v.([]interface{}) {
			interconnectPairsMap := item.(map[string]interface{})
			interconnectPair := cfwv20190904.InterconnectPair{}
			if v, ok := interconnectPairsMap["group_a"]; ok {
				for _, item := range v.([]interface{}) {
					groupAMap := item.(map[string]interface{})
					accessInstanceInfo := cfwv20190904.AccessInstanceInfo{}
					if v, ok := groupAMap["instance_id"].(string); ok && v != "" {
						accessInstanceInfo.InstanceId = helper.String(v)
					}

					if v, ok := groupAMap["instance_type"].(string); ok && v != "" {
						accessInstanceInfo.InstanceType = helper.String(v)
					}

					if v, ok := groupAMap["instance_region"].(string); ok && v != "" {
						accessInstanceInfo.InstanceRegion = helper.String(v)
					}

					if v, ok := groupAMap["access_cidr_mode"].(int); ok {
						accessInstanceInfo.AccessCidrMode = helper.IntInt64(v)
					}

					if v, ok := groupAMap["access_cidr_list"]; ok {
						accessCidrListSet := v.(*schema.Set).List()
						for i := range accessCidrListSet {
							accessCidrList := accessCidrListSet[i].(string)
							accessInstanceInfo.AccessCidrList = append(accessInstanceInfo.AccessCidrList, helper.String(accessCidrList))
						}
					}

					interconnectPair.GroupA = append(interconnectPair.GroupA, &accessInstanceInfo)
				}
			}

			if v, ok := interconnectPairsMap["group_b"]; ok {
				for _, item := range v.([]interface{}) {
					groupBMap := item.(map[string]interface{})
					accessInstanceInfo := cfwv20190904.AccessInstanceInfo{}
					if v, ok := groupBMap["instance_id"].(string); ok && v != "" {
						accessInstanceInfo.InstanceId = helper.String(v)
					}

					if v, ok := groupBMap["instance_type"].(string); ok && v != "" {
						accessInstanceInfo.InstanceType = helper.String(v)
					}

					if v, ok := groupBMap["instance_region"].(string); ok && v != "" {
						accessInstanceInfo.InstanceRegion = helper.String(v)
					}

					if v, ok := groupBMap["access_cidr_mode"].(int); ok {
						accessInstanceInfo.AccessCidrMode = helper.IntInt64(v)
					}

					if v, ok := groupBMap["access_cidr_list"]; ok {
						accessCidrListSet := v.(*schema.Set).List()
						for i := range accessCidrListSet {
							accessCidrList := accessCidrListSet[i].(string)
							accessInstanceInfo.AccessCidrList = append(accessInstanceInfo.AccessCidrList, helper.String(accessCidrList))
						}
					}

					interconnectPair.GroupB = append(interconnectPair.GroupB, &accessInstanceInfo)
				}
			}

			if v, ok := interconnectPairsMap["interconnect_mode"].(string); ok && v != "" {
				interconnectPair.InterconnectMode = helper.String(v)
			}

			ccnSwitchInfo.InterconnectPairs = append(ccnSwitchInfo.InterconnectPairs, &interconnectPair)
		}
	}

	request.CcnSwitch = append(request.CcnSwitch, &ccnSwitchInfo)
	request.Enable = helper.IntInt64(1)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyClusterVpcFwSwitchWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := cfwv20190904.NewDescribeClusterVpcFwSwitchsRequest()
	waitReq.Filters = []*cfwv20190904.CommonFilter{
		{
			Name:         helper.String("InsObj"),
			OperatorType: helper.IntInt64(1),
			Values:       helper.Strings([]string{ccnId}),
		},
	}
	waitReq.Offset = helper.IntUint64(0)
	waitReq.Limit = helper.IntUint64(20)
	reqErr = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().DescribeClusterVpcFwSwitchsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster vpc fw switchs ailed, Response is nil."))
		}

		if len(result.Response.Data) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Data is empty."))
		}

		obj := result.Response.Data[0]
		if obj != nil && obj.Status != nil {
			if *obj.Status == 1 {
				return nil
			}

			// create error
			if *obj.Status == 0 {
				service := CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
				respData, e := service.DescribeCfwVpcSwitchErrorById(ctx, ccnId, "ERR_VPC_FW_OPEN_FAILED")
				if e != nil {
					return resource.NonRetryableError(e)
				}

				if respData == nil || respData.ErrMsg == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe switch error failed. Response is nil."))
				}

				errMsg := *respData.ErrMsg
				return resource.NonRetryableError(fmt.Errorf("Cluster vpc fw switch create failed. Reason:%s", errMsg))
			}
		}

		return resource.RetryableError(fmt.Errorf("wait for cluster vpc fw switch create."))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(ccnId)
	return resourceTencentCloudCfwClusterVpcFwSwitchRead(d, meta)
}

func resourceTencentCloudCfwClusterVpcFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_vpc_fw_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   = d.Id()
	)

	respData, err := service.DescribeCfwClusterVpcFwSwitchById(ctx, ccnId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_cluster_vpc_fw_switch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	respData1, err := service.DescribeCfwClusterVpcFwSwitchsById(ctx, ccnId)
	if err != nil {
		return err
	}

	if respData1 == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_cluster_vpc_fw_switch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("ccn_id", ccnId)

	interconnectPairsList := make([]map[string]interface{}, 0, len(respData))
	for _, interconnectPairs := range respData {
		interconnectPairsMap := map[string]interface{}{}
		groupAList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupA))
		if interconnectPairs.GroupA != nil {
			for _, groupA := range interconnectPairs.GroupA {
				groupAMap := map[string]interface{}{}
				if groupA.InstanceId != nil {
					groupAMap["instance_id"] = groupA.InstanceId
				}

				if groupA.InstanceType != nil {
					groupAMap["instance_type"] = groupA.InstanceType
				}

				if groupA.InstanceRegion != nil {
					groupAMap["instance_region"] = groupA.InstanceRegion
				}

				if groupA.AccessCidrMode != nil {
					groupAMap["access_cidr_mode"] = groupA.AccessCidrMode
				}

				if groupA.AccessCidrList != nil {
					groupAMap["access_cidr_list"] = groupA.AccessCidrList
				}

				groupAList = append(groupAList, groupAMap)
			}

			interconnectPairsMap["group_a"] = groupAList
		}

		groupBList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupB))
		if interconnectPairs.GroupB != nil {
			for _, groupB := range interconnectPairs.GroupB {
				groupBMap := map[string]interface{}{}
				if groupB.InstanceId != nil {
					groupBMap["instance_id"] = groupB.InstanceId
				}

				if groupB.InstanceType != nil {
					groupBMap["instance_type"] = groupB.InstanceType
				}

				if groupB.InstanceRegion != nil {
					groupBMap["instance_region"] = groupB.InstanceRegion
				}

				if groupB.AccessCidrMode != nil {
					groupBMap["access_cidr_mode"] = groupB.AccessCidrMode
				}

				if groupB.AccessCidrList != nil {
					groupBMap["access_cidr_list"] = groupB.AccessCidrList
				}

				groupBList = append(groupBList, groupBMap)
			}

			interconnectPairsMap["group_b"] = groupBList
		}

		if interconnectPairs.InterconnectMode != nil {
			interconnectPairsMap["interconnect_mode"] = interconnectPairs.InterconnectMode
		}

		interconnectPairsList = append(interconnectPairsList, interconnectPairsMap)
		_ = d.Set("interconnect_pairs", interconnectPairsList)
	}

	if respData1.SwitchMode != nil {
		_ = d.Set("switch_mode", respData1.SwitchMode)
	}

	if respData1.RoutingMode != nil {
		_ = d.Set("routing_mode", respData1.RoutingMode)
	}

	return nil
}

func resourceTencentCloudCfwClusterVpcFwSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_vpc_fw_switch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		ccnId = d.Id()
	)

	if d.HasChange("region_cidr_configs") || d.HasChange("interconnect_pairs") {
		request := cfwv20190904.NewUpdateClusterVpcFwRequest()
		ccnSwitchInfo := cfwv20190904.CcnSwitchInfo{}
		if v, ok := d.GetOk("switch_mode"); ok {
			ccnSwitchInfo.SwitchMode = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("routing_mode"); ok {
			ccnSwitchInfo.RoutingMode = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("region_cidr_configs"); ok {
			for _, item := range v.([]interface{}) {
				regionCidrConfigsMap := item.(map[string]interface{})
				regionCidrConfig := cfwv20190904.RegionCidrConfig{}
				if v, ok := regionCidrConfigsMap["region"].(string); ok && v != "" {
					regionCidrConfig.Region = helper.String(v)
				}

				if v, ok := regionCidrConfigsMap["cidr_mode"].(int); ok {
					regionCidrConfig.CidrMode = helper.IntInt64(v)
				}

				if v, ok := regionCidrConfigsMap["custom_cidr"].(string); ok && v != "" {
					regionCidrConfig.CustomCidr = helper.String(v)
				}

				ccnSwitchInfo.RegionCidrConfigs = append(ccnSwitchInfo.RegionCidrConfigs, &regionCidrConfig)
			}
		}

		if v, ok := d.GetOk("interconnect_pairs"); ok {
			for _, item := range v.([]interface{}) {
				interconnectPairsMap := item.(map[string]interface{})
				interconnectPair := cfwv20190904.InterconnectPair{}
				if v, ok := interconnectPairsMap["group_a"]; ok {
					for _, item := range v.([]interface{}) {
						groupAMap := item.(map[string]interface{})
						accessInstanceInfo := cfwv20190904.AccessInstanceInfo{}
						if v, ok := groupAMap["instance_id"].(string); ok && v != "" {
							accessInstanceInfo.InstanceId = helper.String(v)
						}

						if v, ok := groupAMap["instance_type"].(string); ok && v != "" {
							accessInstanceInfo.InstanceType = helper.String(v)
						}

						if v, ok := groupAMap["instance_region"].(string); ok && v != "" {
							accessInstanceInfo.InstanceRegion = helper.String(v)
						}

						if v, ok := groupAMap["access_cidr_mode"].(int); ok {
							accessInstanceInfo.AccessCidrMode = helper.IntInt64(v)
						}

						if v, ok := groupAMap["access_cidr_list"]; ok {
							accessCidrListSet := v.(*schema.Set).List()
							for i := range accessCidrListSet {
								accessCidrList := accessCidrListSet[i].(string)
								accessInstanceInfo.AccessCidrList = append(accessInstanceInfo.AccessCidrList, helper.String(accessCidrList))
							}
						}

						interconnectPair.GroupA = append(interconnectPair.GroupA, &accessInstanceInfo)
					}
				}

				if v, ok := interconnectPairsMap["group_b"]; ok {
					for _, item := range v.([]interface{}) {
						groupBMap := item.(map[string]interface{})
						accessInstanceInfo := cfwv20190904.AccessInstanceInfo{}
						if v, ok := groupBMap["instance_id"].(string); ok && v != "" {
							accessInstanceInfo.InstanceId = helper.String(v)
						}

						if v, ok := groupBMap["instance_type"].(string); ok && v != "" {
							accessInstanceInfo.InstanceType = helper.String(v)
						}

						if v, ok := groupBMap["instance_region"].(string); ok && v != "" {
							accessInstanceInfo.InstanceRegion = helper.String(v)
						}

						if v, ok := groupBMap["access_cidr_mode"].(int); ok {
							accessInstanceInfo.AccessCidrMode = helper.IntInt64(v)
						}

						if v, ok := groupBMap["access_cidr_list"]; ok {
							accessCidrListSet := v.(*schema.Set).List()
							for i := range accessCidrListSet {
								accessCidrList := accessCidrListSet[i].(string)
								accessInstanceInfo.AccessCidrList = append(accessInstanceInfo.AccessCidrList, helper.String(accessCidrList))
							}
						}

						interconnectPair.GroupB = append(interconnectPair.GroupB, &accessInstanceInfo)
					}
				}

				if v, ok := interconnectPairsMap["interconnect_mode"].(string); ok && v != "" {
					interconnectPair.InterconnectMode = helper.String(v)
				}

				ccnSwitchInfo.InterconnectPairs = append(ccnSwitchInfo.InterconnectPairs, &interconnectPair)
			}
		}

		ccnSwitchInfo.CcnId = &ccnId
		request.CcnSwitch = &ccnSwitchInfo
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().UpdateClusterVpcFwWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		time.Sleep(60 * time.Second)
		waitReq := cfwv20190904.NewDescribeClusterVpcFwSwitchsRequest()
		waitReq.Filters = []*cfwv20190904.CommonFilter{
			{
				Name:         helper.String("InsObj"),
				OperatorType: helper.IntInt64(1),
				Values:       helper.Strings([]string{ccnId}),
			},
		}
		waitReq.Offset = helper.IntUint64(0)
		waitReq.Limit = helper.IntUint64(20)
		reqErr = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().DescribeClusterVpcFwSwitchsWithContext(ctx, waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe cluster vpc fw switchs ailed, Response is nil."))
			}

			if len(result.Response.Data) == 0 {
				return resource.NonRetryableError(fmt.Errorf("Data is empty."))
			}

			obj := result.Response.Data[0]
			if obj != nil && obj.Status != nil {
				if *obj.Status == 1 {
					return nil
				}

				// update error
				if *obj.Status == 0 {
					service := CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
					respData, e := service.DescribeCfwVpcSwitchErrorById(ctx, ccnId, "ERR_VPC_FW_UPDATE_FAILED")
					if e != nil {
						return resource.NonRetryableError(e)
					}

					if respData == nil || respData.ErrMsg == nil {
						return resource.NonRetryableError(fmt.Errorf("Describe switch error failed. Response is nil."))
					}

					errMsg := *respData.ErrMsg
					return resource.NonRetryableError(fmt.Errorf("Cluster vpc fw switch update failed. Reason:%s", errMsg))
				}
			}

			return resource.RetryableError(fmt.Errorf("wait for cluster vpc fw switch update."))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudCfwClusterVpcFwSwitchRead(d, meta)
}

func resourceTencentCloudCfwClusterVpcFwSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_vpc_fw_switch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewModifyClusterVpcFwSwitchRequest()
		ccnId   = d.Id()
	)

	ccnSwitchInfo := cfwv20190904.CcnSwitchInfo{}
	ccnSwitchInfo.CcnId = &ccnId
	if v, ok := d.GetOk("switch_mode"); ok {
		ccnSwitchInfo.SwitchMode = helper.IntUint64(v.(int))
	}

	request.CcnSwitch = append(request.CcnSwitch, &ccnSwitchInfo)
	request.Enable = helper.IntInt64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyClusterVpcFwSwitchWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := cfwv20190904.NewDescribeClusterVpcFwSwitchsRequest()
	waitReq.Filters = []*cfwv20190904.CommonFilter{
		{
			Name:         helper.String("InsObj"),
			OperatorType: helper.IntInt64(1),
			Values:       helper.Strings([]string{ccnId}),
		},
	}
	waitReq.Offset = helper.IntUint64(0)
	waitReq.Limit = helper.IntUint64(20)
	reqErr = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().DescribeClusterVpcFwSwitchsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster vpc fw switchs ailed, Response is nil."))
		}

		if len(result.Response.Data) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Data is empty."))
		}

		obj := result.Response.Data[0]
		if obj != nil && obj.Status != nil {
			if *obj.Status == 0 {
				return nil
			}

			// delete error
			if *obj.Status == 1 {
				service := CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
				respData, e := service.DescribeCfwVpcSwitchErrorById(ctx, ccnId, "ERR_VPC_FW_CLOSE_FAILED")
				if e != nil {
					return resource.NonRetryableError(e)
				}

				if respData == nil || respData.ErrMsg == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe switch error failed. Response is nil."))
				}

				errMsg := *respData.ErrMsg
				return resource.NonRetryableError(fmt.Errorf("Cluster vpc fw switch delete failed. Reason:%s", errMsg))
			}
		}

		return resource.RetryableError(fmt.Errorf("wait for cluster vpc fw switch delete."))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cfw cluster vpc fw switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
