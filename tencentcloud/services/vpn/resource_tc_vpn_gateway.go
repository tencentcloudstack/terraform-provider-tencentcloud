package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnGatewayCreate,
		Read:   resourceTencentCloudVpnGatewayRead,
		Update: resourceTencentCloudVpnGatewayUpdate,
		Delete: resourceTencentCloudVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"prepaid_period": 1,
			}),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the VPN gateway. The length of character is limited to 1-60.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("type"); ok && (v.(string) == "CCN" || v.(string) == "SSL_CCN") {
						return true
					}
					return old == new
				},
				Description: "ID of the VPC. Required if vpn gateway is not in `CCN` or `SSL_CCN` type, and doesn't make sense for `CCN` or `SSL_CCN` vpn gateway.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100,200,500,1000. Default is 5. When charge type is `PREPAID`, bandwidth degradation operation is unsupported.",
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP of the VPN gateway.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Type of gateway instance, Default is `IPSEC`. Valid value: `IPSEC`, `SSL`, `CCN` and `SSL_CCN`.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the VPN gateway. Valid value: `PENDING`, `DELETING`, `AVAILABLE`.",
			},
			"prepaid_renew_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_PERIOD_PREPAID_RENEW_FLAG_AUTO_NOTIFY,
				Description: "Flag indicates whether to renew or not. Valid value: `NOTIFY_AND_AUTO_RENEW`, `NOTIFY_AND_MANUAL_RENEW`.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3, 4, 6, 7, 8, 9, 12, 24, 36}),
				Description:  "Period of instance to be prepaid. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. The unit is month. Caution: when this para and renew_flag para are valid, the request means to renew several months more pre-paid period. This para can only be changed on `IPSEC` vpn gateway.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_CHARGE_TYPE_POSTPAID_BY_HOUR,
				Description: "Charge Type of the VPN gateway. Valid value: `PREPAID`, `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.",
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CDC instance ID.",
			},
			"max_connection": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of connected clients allowed for the SSL VPN gateway. Valid values: [5, 10, 20, 50, 100]. This parameter is only required for SSL VPN gateways.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expired time of the VPN gateway when charge type is `PREPAID`.",
			},
			"is_address_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether ip address is blocked.",
			},
			"new_purchase_plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The plan of new purchase. Valid value: `PREPAID_TO_POSTPAID`.",
			},
			"restrict_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Restrict state of gateway. Valid value: `PRETECIVELY_ISOLATED`, `NORMAL`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Zone of the VPN gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},
			"bgp_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "BGP ASN. Value range: 1 - 4294967295. Using BGP requires configuring ASN.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the VPN gateway.",
			},
		},
	}
}

func resourceTencentCloudVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := vpc.NewCreateVpnGatewayRequest()
	request.VpnGatewayName = helper.String(d.Get("name").(string))
	bandwidth := d.Get("bandwidth").(int)
	bandwidth64 := uint64(bandwidth)
	request.InternetMaxBandwidthOut = &bandwidth64
	chargeType := d.Get("charge_type").(string)

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	//only support change renew_flag when charge type is pre-paid
	if chargeType == svcvpc.VPN_CHARGE_TYPE_PREPAID {
		var preChargePara vpc.InstanceChargePrepaid
		preChargePara.Period = helper.IntUint64(d.Get("prepaid_period").(int))
		preChargePara.RenewFlag = helper.String(d.Get("prepaid_renew_flag").(string))
		request.InstanceChargePrepaid = &preChargePara
	}
	request.InstanceChargeType = &chargeType
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
		if v.(string) != "CCN" && v.(string) != "SSL_CCN" {
			if _, ok := d.GetOk("vpc_id"); !ok {
				return fmt.Errorf("[CRITAL] vpc_id is required for vpn gateway in %s type", v.(string))
			}
			request.VpcId = helper.String(d.Get("vpc_id").(string))
		} else {
			if _, ok := d.GetOk("vpc_id"); ok {
				return fmt.Errorf("[CRITAL] vpc_id doesn't make sense when vpn gateway is in CCN type")
			}
			request.VpcId = helper.String("")
		}
	} else {
		if _, ok := d.GetOk("vpc_id"); !ok {
			return fmt.Errorf("[CRITAL] vpc_id is required for vpn gateway in %s type", "IPSEC")
		}
		request.VpcId = helper.String(d.Get("vpc_id").(string))
	}

	if v, ok := d.GetOk("cdc_id"); ok {
		request.CdcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_connection"); ok {
		request.MaxConnection = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("bgp_asn"); ok {
		request.BgpAsn = helper.IntUint64(v.(int))
	}

	var response *vpc.CreateVpnGatewayResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpnGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.VpnGateway == nil {
			return resource.NonRetryableError(fmt.Errorf("create VPN gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.VpnGateway.VpnGatewayId == nil {
		return fmt.Errorf("VPN gateway id is nil")
	}

	gatewayId := *response.Response.VpnGateway.VpnGatewayId
	d.SetId(gatewayId)

	// must wait for creating gateway finished
	statRequest := vpc.NewDescribeVpnGatewaysRequest()
	statRequest.VpnGatewayIds = []*string{helper.String(gatewayId)}
	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, statRequest.GetAction(), statRequest.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			//if not, quit
			if result != nil && result.Response != nil && result.Response.VpnGatewaySet != nil {
				if len(result.Response.VpnGatewaySet) != 1 {
					return resource.NonRetryableError(fmt.Errorf("creating error"))
				} else {
					if *result.Response.VpnGatewaySet[0].State == svcvpc.VPN_STATE_AVAILABLE {
						return nil
					} else {
						return resource.RetryableError(fmt.Errorf("State is not available: %s, wait for state to be AVAILABLE.", *result.Response.VpnGatewaySet[0].State))
					}
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("Describe Vpn Gateways failed, Response is nil."))
			}
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpngw", region, gatewayId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnGatewayRead(d, meta)
}

func resourceTencentCloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		gatewayId = d.Id()
	)

	has, gateway, err := service.DescribeVpngwById(ctx, gatewayId)
	if err != nil {
		log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", gateway.VpnGatewayName)
	_ = d.Set("public_ip_address", gateway.PublicIpAddress)
	_ = d.Set("bandwidth", int(*gateway.InternetMaxBandwidthOut))
	_ = d.Set("type", gateway.Type)
	_ = d.Set("create_time", gateway.CreatedTime)
	_ = d.Set("state", gateway.State)
	if gateway.RenewFlag != nil {
		_ = d.Set("prepaid_renew_flag", *gateway.RenewFlag)
	} else {
		_ = d.Set("prepaid_renew_flag", svcvpc.VPN_PERIOD_PREPAID_RENEW_FLAG_AUTO_NOTIFY)
	}
	_ = d.Set("charge_type", gateway.InstanceChargeType)
	_ = d.Set("expired_time", gateway.ExpiredTime)
	_ = d.Set("is_address_blocked", gateway.IsAddressBlocked)
	_ = d.Set("new_purchase_plan", gateway.NewPurchasePlan)
	_ = d.Set("restrict_state", gateway.RestrictState)
	_ = d.Set("zone", gateway.Zone)
	_ = d.Set("cdc_id", gateway.CdcId)
	_ = d.Set("max_connection", gateway.MaxConnection)
	if gateway.BgpAsn != nil && *gateway.BgpAsn != 0 {
		_ = d.Set("bgp_asn", gateway.BgpAsn)
	}

	//tags
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpngw", region, gatewayId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	d.Partial(true)
	gatewayId := d.Id()

	unsupportedUpdateFields := []string{
		"type",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("Template resource_tc_vpn_gateway update on %s is not supportted yet. Please renew it on controller web page.", field)
		}
	}

	if d.HasChange("prepaid_period") {
		chargeType := d.Get("charge_type").(string)
		period := d.Get("prepaid_period").(int)
		if chargeType != svcvpc.VPN_CHARGE_TYPE_PREPAID {
			return fmt.Errorf("Invalid renew flag change. Only support pre-paid vpn.")
		}
		request := vpc.NewRenewVpnGatewayRequest()
		request.VpnGatewayId = &gatewayId
		var preChargePara vpc.InstanceChargePrepaid
		preChargePara.Period = helper.IntUint64(period)
		request.InstanceChargePrepaid = &preChargePara

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().RenewVpnGateway(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway prepaid period failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("prepaid_renew_flag") {
		chargeType := d.Get("charge_type").(string)
		renewFlag := d.Get("prepaid_renew_flag").(string)
		if chargeType != svcvpc.VPN_CHARGE_TYPE_PREPAID {
			return fmt.Errorf("Invalid renew flag change. Only support pre-paid vpn.")
		}
		request := vpc.NewSetVpnGatewaysRenewFlagRequest()
		request.VpnGatewayIds = []*string{&gatewayId}
		if renewFlag == "NOTIFY_AND_AUTO_RENEW" {
			request.AutoRenewFlag = helper.IntInt64(1)
		} else {
			request.AutoRenewFlag = helper.IntInt64(0)
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().SetVpnGatewaysRenewFlag(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway renewflag failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("name") || d.HasChange("charge_type") || d.HasChange("bgp_asn") {
		//check that the charge type change is valid
		//only pre-paid --> post-paid is valid
		oldInterface, newInterface := d.GetChange("charge_type")
		oldChargeType := oldInterface.(string)
		newChargeType := newInterface.(string)
		request := vpc.NewModifyVpnGatewayAttributeRequest()
		request.VpnGatewayId = &gatewayId
		request.VpnGatewayName = helper.String(d.Get("name").(string))
		if v, ok := d.GetOkExists("bgp_asn"); ok {
			request.BgpAsn = helper.IntUint64(v.(int))
		}
		if oldChargeType == svcvpc.VPN_CHARGE_TYPE_PREPAID && newChargeType == svcvpc.VPN_CHARGE_TYPE_POSTPAID_BY_HOUR {
			request.InstanceChargeType = &newChargeType
		} else if oldChargeType == svcvpc.VPN_CHARGE_TYPE_POSTPAID_BY_HOUR && newChargeType == svcvpc.VPN_CHARGE_TYPE_PREPAID {
			return fmt.Errorf("Invalid charge type change. Only support pre-paid to post-paid way.")
		}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpnGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway name failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	//bandwidth
	if d.HasChange("bandwidth") {
		request := vpc.NewResetVpnGatewayInternetMaxBandwidthRequest()
		request.VpnGatewayId = &gatewayId
		bandwidth := d.Get("bandwidth").(int)
		bandwidth64 := uint64(bandwidth)
		request.InternetMaxBandwidthOut = &bandwidth64
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetVpnGatewayInternetMaxBandwidth(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway bandwidth failed, reason:%s\n", logId, err.Error())
			return err
		}

	}

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpngw", region, gatewayId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	if d.HasChange("cdc_id") || d.HasChange("max_connection") {
		return fmt.Errorf("cdc_id and max_connection do not support change now.")
	}

	d.Partial(false)

	return resourceTencentCloudVpnGatewayRead(d, meta)
}

func resourceTencentCloudVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	gatewayId := d.Id()

	//prepaid instances or instances which attached to ccn can not be deleted
	//to get expire_time of the VPN gateway
	//to get the status of gateway
	//to get the type and networkinstanceid of gateway
	vpngwRequest := vpc.NewDescribeVpnGatewaysRequest()
	vpngwRequest.VpnGatewayIds = []*string{&gatewayId}
	vpngwErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnGateways(vpngwRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			if result != nil && result.Response != nil && result.Response.VpnGatewaySet != nil {
				//if deleted, quit
				if len(result.Response.VpnGatewaySet) == 0 {
					return nil
				}
				if result.Response.VpnGatewaySet[0].ExpiredTime != nil && *result.Response.VpnGatewaySet[0].InstanceChargeType == svcvpc.VPN_CHARGE_TYPE_PREPAID {
					expiredTime := *result.Response.VpnGatewaySet[0].ExpiredTime
					if expiredTime != "0000-00-00 00:00:00" {
						t, err := time.Parse("2006-01-02 15:04:05", expiredTime)
						if err != nil {
							return resource.NonRetryableError(fmt.Errorf("Error format expired time.%x %s", expiredTime, err))
						}
						if time.Until(t) > 0 {
							return resource.NonRetryableError(fmt.Errorf("Delete operation is unsupport when VPN gateway is not expired."))
						}
					}
				}
				if *result.Response.VpnGatewaySet[0].Type == svcvpc.GATE_WAY_TYPE_CCN && *result.Response.VpnGatewaySet[0].NetworkInstanceId != "" {
					return resource.NonRetryableError(fmt.Errorf("Delete operation is unsupported when VPN gateway is attached to CCN instance."))
				}
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Describe Vpn Gateways failed, Response is nil."))
			}
		}
	})
	if vpngwErr != nil {
		log.Printf("[CRITAL]%s describe VPN gateway failed, reason:%s\n", logId, vpngwErr.Error())
		return vpngwErr
	}

	//check the vpn gateway is not related with any tunnel
	tRequest := vpc.NewDescribeVpnConnectionsRequest()
	tRequest.Filters = make([]*vpc.Filter, 0, 2)
	params := make(map[string]string)
	params["vpn-gateway-id"] = gatewayId

	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}

	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		tRequest.Filters = append(tRequest.Filters, filter)
	}
	offset := uint64(0)
	tRequest.Offset = &offset

	tErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(tRequest)

		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, tRequest.GetAction(), tRequest.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			if result != nil && result.Response != nil && result.Response.VpnConnectionSet != nil {
				if len(result.Response.VpnConnectionSet) == 0 {
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("There is associated tunnel exists, please delete associated tunnels first."))
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("Describe Vpn Connections failed, Response is nil."))
			}
		}
	})
	if tErr != nil {
		log.Printf("[CRITAL]%s describe VPN connection failed, reason:%s\n", logId, tErr.Error())
		return tErr
	}

	request := vpc.NewDeleteVpnGatewayRequest()
	request.VpnGatewayId = &gatewayId

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteVpnGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	//to get the status of gateway
	statRequest := vpc.NewDescribeVpnGatewaysRequest()
	statRequest.VpnGatewayIds = []*string{&gatewayId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnGateways(statRequest)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(e)
			}
			if ee.Code == svcvpc.VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
		} else {
			if result != nil && result.Response != nil && result.Response.VpnGatewaySet != nil {
				//if not, quit
				if len(result.Response.VpnGatewaySet) == 0 {
					return nil
				}
				//else consider delete fail
				return resource.RetryableError(fmt.Errorf("deleting retry"))
			} else {
				return resource.NonRetryableError(fmt.Errorf("Describe Vpn Gateways failed, Response is nil."))
			}
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
