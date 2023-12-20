package cfw

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwNatInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatInstanceCreate,
		Read:   resourceTencentCloudCfwNatInstanceRead,
		Update: resourceTencentCloudCfwNatInstanceUpdate,
		Delete: resourceTencentCloudCfwNatInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Firewall instance name.",
			},
			"width": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerMin(BAND_WIDTH),
				Description:  "Bandwidth.",
			},
			"mode": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(MODE),
				Description:  "Mode 1: access mode; 0: new mode.",
			},
			"new_mode_items": {
				Optional:     true,
				Type:         schema.TypeList,
				MaxItems:     1,
				ExactlyOneOf: []string{"nat_gw_list"},
				Description:  "New mode passing parameters are added, at least one of new_mode_items and nat_gw_list is passed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_list": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "List of vpcs connected in new mode.",
						},
						"eips": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "List of egress elastic public network IPs bound in the new mode.",
						},
					},
				},
			},
			"nat_gw_list": {
				Optional:     true,
				Type:         schema.TypeSet,
				ExactlyOneOf: []string{"new_mode_items"},
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  "A list of nat gateways connected to the access mode, at least one of NewModeItems and NatgwList is passed.",
			},
			"zone_set": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    2,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Zone list.",
			},
			"cross_a_zone": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      CROSS_A_ZONE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CROSS_A_ZONE),
				Description:  "Off-site disaster recovery 1: use off-site disaster recovery; 0: do not use off-site disaster recovery; if empty, the default is not to use off-site disaster recovery.",
			},
			//"domain": {
			//	Optional:    true,
			//	Type:        schema.TypeString,
			//	Description: "Required if you want to create a domain name.",
			//},
			//"fw_cidr_info": {
			//	Optional:    true,
			//	Type:        schema.TypeList,
			//	MaxItems:    1,
			//	Description: "Specify the network segment information used by the firewall.",
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"fw_cidr_type": {
			//				Type:        schema.TypeString,
			//				Required:    true,
			//				Description: "The type of network segment used by the firewall. The values VpcSelf/Assis/Custom respectively represent own network segment priority/extended network segment priority/custom.",
			//			},
			//			"fw_cidr_lst": {
			//				Type:        schema.TypeList,
			//				Optional:    true,
			//				Description: "Specify the network segment of the firewall for each vpc.",
			//				Elem: &schema.Resource{
			//					Schema: map[string]*schema.Schema{
			//						"vpc_id": {
			//							Type:        schema.TypeString,
			//							Required:    true,
			//							Description: "Vpc id.",
			//						},
			//						"fw_cidr": {
			//							Type:        schema.TypeString,
			//							Required:    true,
			//							Description: "Firewall network segment, at least /24 network segment.",
			//						},
			//					},
			//				},
			//			},
			//			"com_fw_cidr": {
			//				Type:        schema.TypeString,
			//				Optional:    true,
			//				Description: "Other firewalls occupy the network segment, which is usually the network segment specified when the firewall needs to exclusively occupy the vpc.",
			//			},
			//		},
			//	},
			//},
		},
	}
}

func resourceTencentCloudCfwNatInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = cfw.NewCreateNatFwInstanceWithDomainRequest()
		response   = cfw.NewCreateNatFwInstanceWithDomainResponse()
		instanceId string
		mode       int
		crossAZone int
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("mode"); ok {
		request.Mode = helper.IntInt64(v.(int))
		mode = v.(int)
	}

	if mode == MODE_0 {
		if v, ok := d.GetOk("new_mode_items"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				newModeItems := cfw.NewModeItems{}
				newModeItems.AddCount = helper.IntInt64(0)
				if v, ok := dMap["vpc_list"]; ok {
					vpcList := v.(*schema.Set).List()
					tmqVpcList := make([]*string, 0, len(vpcList))
					for i := range vpcList {
						vpc := vpcList[i].(string)
						tmqVpcList = append(tmqVpcList, &vpc)
					}
					newModeItems.VpcList = tmqVpcList
				}

				if v, ok := dMap["eips"]; ok {
					eipList := v.(*schema.Set).List()
					tmqEipList := make([]*string, 0, len(eipList))
					for i := range eipList {
						eip := eipList[i].(string)
						tmqEipList = append(tmqEipList, &eip)
					}
					newModeItems.Eips = tmqEipList
				}

				request.NewModeItems = &newModeItems
			}

		} else {
			return fmt.Errorf("If `mode` is 0, `new_mode_items` is required.")
		}

	} else {
		if v, ok := d.GetOk("nat_gw_list"); ok {
			gwList := v.(*schema.Set).List()
			tmqGwList := make([]*string, 0, len(gwList))
			for i := range gwList {
				gw := gwList[i].(string)
				tmqGwList = append(tmqGwList, &gw)
			}

			request.NatGwList = tmqGwList

		} else {
			return fmt.Errorf("If `mode` is 1, `nat_gw_list` is required.")
		}
	}

	if v, ok := d.GetOkExists("cross_a_zone"); ok {
		crossAZone = v.(int)
		if v, ok := d.GetOk("zone_set"); ok {
			zoneList := v.(*schema.Set).List()
			if crossAZone == CROSS_A_ZONE_0 {
				if len(zoneList) != 1 {
					return fmt.Errorf("If `cross_a_zone` is 0, `zone_set` only can be set one zone.")
				}

				request.Zone = helper.String(zoneList[0].(string))

			} else {
				if len(zoneList) != 2 {
					return fmt.Errorf("If `cross_a_zone` is 1, `zone_set` must be set tow zones.")
				}

				request.Zone = helper.String(zoneList[0].(string))
				request.ZoneBak = helper.String(zoneList[1].(string))
			}
		}

		request.CrossAZone = helper.IntInt64(v.(int))
	}

	//if v, ok := d.GetOk("domain"); ok {
	//	request.Domain = helper.String(v.(string))
	//	request.IsCreateDomain = helper.IntInt64(1)
	//} else {
	//	request.IsCreateDomain = helper.IntInt64(0)
	//}
	//
	//if v, ok := d.GetOk("fw_cidr_info"); ok {
	//	for _, item := range v.([]interface{}) {
	//		dMap := item.(map[string]interface{})
	//		fwCidrInfo := cfw.FwCidrInfo{}
	//		if v, ok := dMap["fw_cidr_type"]; ok {
	//			fwCidrInfo.FwCidrType = helper.String(v.(string))
	//		}
	//
	//		if v, ok := dMap["com_fw_cidr"]; ok {
	//			fwCidrInfo.ComFwCidr = helper.String(v.(string))
	//		}
	//
	//		if v, ok := dMap["fw_cidr_lst"]; ok {
	//			for _, cidr := range v.([]interface{}) {
	//				iMap := cidr.(map[string]interface{})
	//				fwCidr := cfw.FwVpcCidr{}
	//				if v, ok := iMap["vpc_id"]; ok {
	//					fwCidr.VpcId = helper.String(v.(string))
	//				}
	//
	//				if v, ok := iMap["fw_cidr"]; ok {
	//					fwCidr.FwCidr = helper.String(v.(string))
	//				}
	//
	//				fwCidrInfo.FwCidrLst = append(fwCidrInfo.FwCidrLst, &fwCidr)
	//			}
	//
	//		}
	//
	//		request.FwCidrInfo = &fwCidrInfo
	//	}
	//}

	fwCidrInfo := cfw.FwCidrInfo{}
	fwCidrInfo.FwCidrType = helper.String("VpcSelf")
	fwCidrInfo.ComFwCidr = helper.String("")
	request.FwCidrInfo = &fwCidrInfo

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().CreateNatFwInstanceWithDomain(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw natInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.CfwInsId
	d.SetId(instanceId)

	// wait
	err = resource.Retry(tccommon.WriteRetryTimeout*3, func() *resource.RetryError {
		natInstance, e := service.DescribeCfwNatInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if natInstance == nil {
			e = fmt.Errorf("cfw nat instance %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		if *natInstance.Status == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("create cfw natInstance status is %d", *natInstance.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw natInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatInstanceRead(d, meta)
}

func resourceTencentCloudCfwNatInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
		zone       string
		zoneBak    string
	)

	natInstance, err := service.DescribeCfwNatInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if natInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwNatInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if natInstance.NatinsName != nil {
		_ = d.Set("name", natInstance.NatinsName)
	}

	if natInstance.BandWidth != nil {
		_ = d.Set("width", natInstance.BandWidth)
	}

	if natInstance.FwMode != nil {
		_ = d.Set("mode", natInstance.FwMode)

		if *natInstance.FwMode == MODE_0 {
			newModeItems := []interface{}{}
			newModeItemsMap := map[string]interface{}{}
			vpcList, err := service.DescribeNatFwVpcDnsLstById(ctx, instanceId)
			if err != nil {
				return err
			}

			if vpcList == nil {
				d.SetId("")
				log.Printf("[WARN]%s resource `Cfw VpcList` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
				return nil
			}

			if vpcList != nil {
				newModeItemsMap["vpc_list"] = vpcList
			}

			if natInstance.EipAddress != nil {
				newModeItemsMap["eips"] = natInstance.EipAddress
			}

			newModeItems = append(newModeItems, newModeItemsMap)
			_ = d.Set("new_mode_items", newModeItems)
		} else {
			natGwList, err := service.DescribeCfwEipsById(ctx, instanceId)
			if err != nil {
				return err
			}

			if natGwList == nil {
				d.SetId("")
				log.Printf("[WARN]%s resource `CfwEips` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
				return nil
			}

			if natGwList != nil {
				_ = d.Set("nat_gw_list", natGwList)
			}
		}
	}

	if natInstance.ZoneZh != nil {
		zoneZh := *natInstance.ZoneZh
		zone = ZONE_MAP_CN2EN[zoneZh]
	}

	if natInstance.ZoneZhBak != nil {
		zoneBakZh := *natInstance.ZoneZhBak
		zoneBak = ZONE_MAP_CN2EN[zoneBakZh]
	}

	if zone == zoneBak {
		_ = d.Set("cross_a_zone", CROSS_A_ZONE_0)
		zoneList := []string{
			zone,
		}
		_ = d.Set("zone_set", zoneList)
	} else {
		_ = d.Set("cross_a_zone", CROSS_A_ZONE_1)
		zoneList := []string{
			zone,
			zoneBak,
		}
		_ = d.Set("zone_set", zoneList)
	}

	//if natInstance.Domain != nil {
	//	_ = d.Set("domain", natInstance.Domain)
	//}
	//
	//if natInstance.FwCidrInfo != nil {
	//	fwCidrInfoMap := map[string]interface{}{}
	//
	//	if natInstance.FwCidrInfo.FwCidrType != nil {
	//		fwCidrInfoMap["fw_cidr_type"] = natInstance.FwCidrInfo.FwCidrType
	//	}
	//
	//	if natInstance.FwCidrInfo.FwCidrLst != nil {
	//		fwCidrLstList := []interface{}{}
	//		for _, fwCidrLst := range natInstance.FwCidrInfo.FwCidrLst {
	//			fwCidrLstMap := map[string]interface{}{}
	//
	//			if fwCidrLst.VpcId != nil {
	//				fwCidrLstMap["vpc_id"] = fwCidrLst.VpcId
	//			}
	//
	//			if fwCidrLst.FwCidr != nil {
	//				fwCidrLstMap["fw_cidr"] = fwCidrLst.FwCidr
	//			}
	//
	//			fwCidrLstList = append(fwCidrLstList, fwCidrLstMap)
	//		}
	//
	//		fwCidrInfoMap["fw_cidr_lst"] = []interface{}{fwCidrLstList}
	//	}
	//
	//	if natInstance.FwCidrInfo.ComFwCidr != nil {
	//		fwCidrInfoMap["com_fw_cidr"] = natInstance.FwCidrInfo.ComFwCidr
	//	}
	//
	//	_ = d.Set("fw_cidr_info", []interface{}{fwCidrInfoMap})
	//}

	return nil
}

func resourceTencentCloudCfwNatInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = cfw.NewModifyNatInstanceRequest()
		instanceId = d.Id()
	)

	immutableArgs := []string{"width", "mode", "new_mode_items", "nat_gw_list", "zone", "zone_bak", "cross_a_zone", "domain", "fw_cidr_info"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.NatInstanceId = &instanceId

	if v, ok := d.GetOk("name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyNatInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw natInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatInstanceRead(d, meta)
}

func resourceTencentCloudCfwNatInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.DeleteCfwNatInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
