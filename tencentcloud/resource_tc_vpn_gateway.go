/*
Provides a resource to create a VPN gateway.

-> **NOTE:** The prepaid VPN gateway do not support renew operation or delete operation with terraform.

Example Usage

POSTPAID_BY_HOUR VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  vpc_id    = "vpc-dk8zmwuf"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}
```

PREPAID VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name           = "test"
  vpc_id         = "vpc-dk8zmwuf"
  bandwidth      = 5
  zone           = "ap-guangzhou-3"
  charge_type    = "PREPAID"
  prepaid_period = 1

  tags = {
    test = "test"
  }
}
```

Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.foo vpngw-8ccsnclt
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnGatewayCreate,
		Read:   resourceTencentCloudVpnGatewayRead,
		Update: resourceTencentCloudVpnGatewayUpdate,
		Delete: resourceTencentCloudVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the VPN gateway. The length of character is limited to 1-60.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the VPC. Required if vpn gateway is not in `CCN` type, and doesn't make sense if vpn gateway is in `CCN` type.",
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
				Description: "Public ip of the VPN gateway.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Type of gateway instance, valid values are `IPSEC`, `SSL` and `CCN`. Note: CCN type is only for whitelist customer now.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the VPN gateway, valid values are `PENDING`, `DELETING`, `AVAILABLE`.",
			},
			"prepaid_renew_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     VPN_PERIOD_PREPAID_RENEW_FLAG_AUTO_NOTIFY,
				Description: "Flag indicates whether to renew or not, valid values are `NOTIFY_AND_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `NOT_NOTIFY_AND_NOT_RENEW`. This para can only be set to take effect in create operation.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateAllowedIntValue([]int{1, 2, 3, 4, 6, 7, 8, 9, 12, 24, 36}),
				Description:  "Period of instance to be prepaid. Valid values are 1, 2, 3, 4, 6, 7, 8, 9, 12, 24, 36 and unit is month. Caution: when this para and renew_flag para are valid, the request means to renew several months more pre-paid period. This para can only be set to take effect in create operation.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     VPN_CHARGE_TYPE_POSTPAID_BY_HOUR,
				Description: "Charge Type of the VPN gateway, valid values are `PREPAID`, `POSTPAID_BY_HOUR` and default is `POSTPAID_BY_HOUR`.",
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
				Description: "The plan of new purchase, valid value is `PREPAID_TO_POSTPAID`.",
			},
			"restrict_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Restrict state of gateway, valid values are `PRETECIVELY_ISOLATED`, `NORMAL`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone of the VPN gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
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
	defer logElapsed("resource.tencentcloud_vpn_gateway.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := vpc.NewCreateVpnGatewayRequest()
	request.VpnGatewayName = helper.String(d.Get("name").(string))
	bandwidth := d.Get("bandwidth").(int)
	bandwidth64 := uint64(bandwidth)
	request.InternetMaxBandwidthOut = &bandwidth64
	request.Zone = helper.String(d.Get("zone").(string))
	chargeType := d.Get("charge_type").(string)
	//only support change renew_flag when charge type is pre-paid
	if chargeType == VPN_CHARGE_TYPE_PREPAID {
		var preChargePara vpc.InstanceChargePrepaid
		preChargePara.Period = helper.IntUint64(d.Get("prepaid_period").(int))
		preChargePara.RenewFlag = helper.String(d.Get("prepaid_renew_flag").(string))
		request.InstanceChargePrepaid = &preChargePara
	}
	request.InstanceChargeType = &chargeType
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
		if v.(string) != "CCN" {
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
			return fmt.Errorf("[CRITAL] vpc_id is required for vpn gateway in %s type", v.(string))
		}
		request.VpcId = helper.String(d.Get("vpc_id").(string))
	}
	var response *vpc.CreateVpnGatewayResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateVpnGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.VpnGateway == nil {
		return fmt.Errorf("VPN gateway id is nil")
	}
	gatewayId := *response.Response.VpnGateway.VpnGatewayId
	d.SetId(gatewayId)

	// must wait for creating gateway finished
	statRequest := vpc.NewDescribeVpnGatewaysRequest()
	statRequest.VpnGatewayIds = []*string{helper.String(gatewayId)}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, statRequest.GetAction(), statRequest.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			//if not, quit
			if len(result.Response.VpnGatewaySet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			} else {
				if *result.Response.VpnGatewaySet[0].State == VPN_STATE_AVAILABLE {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("State is not available: %s, wait for state to be AVAILABLE.", *result.Response.VpnGatewaySet[0].State))
				}
			}
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "vpngw", region, gatewayId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnGatewayRead(d, meta)
}

func resourceTencentCloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	gatewayId := d.Id()
	request := vpc.NewDescribeVpnGatewaysRequest()
	request.VpnGatewayIds = []*string{&gatewayId}
	var response *vpc.DescribeVpnGatewaysResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(request)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(e)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(response.Response.VpnGatewaySet) < 1 {
		d.SetId("")
		return nil
	}

	gateway := response.Response.VpnGatewaySet[0]

	_ = d.Set("name", gateway.VpnGatewayName)
	_ = d.Set("public_ip_address", gateway.PublicIpAddress)
	_ = d.Set("bandwidth", int(*gateway.InternetMaxBandwidthOut))
	_ = d.Set("type", gateway.Type)
	_ = d.Set("create_time", gateway.CreatedTime)
	_ = d.Set("state", gateway.State)
	_ = d.Set("prepaid_renew_flag", gateway.RenewFlag)
	_ = d.Set("charge_type", gateway.InstanceChargeType)
	_ = d.Set("expired_time", gateway.ExpiredTime)
	_ = d.Set("is_address_blocked", gateway.IsAddressBlocked)
	_ = d.Set("new_purchase_plan", gateway.NewPurchasePlan)
	_ = d.Set("restrict_state", gateway.RestrictState)
	_ = d.Set("zone", gateway.Zone)

	//tags
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpngw", region, gatewayId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	d.Partial(true)
	gatewayId := d.Id()

	unsupportedUpdateFields := []string{
		"prepaid_period",
		"prepaid_renew_flag",
		"type",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("Template resource_tc_vpn_gateway update on %s is not supportted yet. Please renew it on controller web page.", field)
		}
	}

	if d.HasChange("name") || d.HasChange("charge_type") {
		//check that the charge type change is valid
		//only pre-paid --> post-paid is valid
		oldInterface, newInterface := d.GetChange("charge_type")
		oldChargeType := oldInterface.(string)
		newChargeType := newInterface.(string)
		request := vpc.NewModifyVpnGatewayAttributeRequest()
		request.VpnGatewayId = &gatewayId
		request.VpnGatewayName = helper.String(d.Get("name").(string))
		if oldChargeType == VPN_CHARGE_TYPE_PREPAID && newChargeType == VPN_CHARGE_TYPE_POSTPAID_BY_HOUR {
			request.InstanceChargeType = &newChargeType
		} else if oldChargeType == VPN_CHARGE_TYPE_POSTPAID_BY_HOUR && newChargeType == VPN_CHARGE_TYPE_PREPAID {
			return fmt.Errorf("Invalid charge type change. Only support pre-paid to post-paid way.")
		}
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyVpnGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway name failed, reason:%s\n", logId, err.Error())
			return err
		}
		if d.HasChange("name") {
			d.SetPartial("name")
		}
		if d.HasChange("charge_type") {
			d.SetPartial("charge_type")
		}
	}

	//bandwidth
	if d.HasChange("bandwidth") {
		request := vpc.NewResetVpnGatewayInternetMaxBandwidthRequest()
		request.VpnGatewayId = &gatewayId
		bandwidth := d.Get("bandwidth").(int)
		bandwidth64 := uint64(bandwidth)
		request.InternetMaxBandwidthOut = &bandwidth64
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ResetVpnGatewayInternetMaxBandwidth(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN gateway bandwidth failed, reason:%s\n", logId, err.Error())
			return err
		}
		d.SetPartial("bandwidth")
	}

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "vpngw", region, gatewayId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudVpnGatewayRead(d, meta)
}

func resourceTencentCloudVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway.delete")()

	logId := getLogId(contextNil)

	gatewayId := d.Id()

	//prepaid instances or instances which attached to ccn can not be deleted
	//to get expire_time of the VPN gateway
	//to get the status of gateway
	//to get the type and networkinstanceid of gateway
	vpngwRequest := vpc.NewDescribeVpnGatewaysRequest()
	vpngwRequest.VpnGatewayIds = []*string{&gatewayId}
	vpngwErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(vpngwRequest)
		if e != nil {
			return retryError(e)
		} else {
			//if deleted, quit
			if len(result.Response.VpnGatewaySet) == 0 {
				return nil
			}
			if result.Response.VpnGatewaySet[0].ExpiredTime != nil && *result.Response.VpnGatewaySet[0].InstanceChargeType == VPN_CHARGE_TYPE_PREPAID {
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
			if *result.Response.VpnGatewaySet[0].Type == GATE_WAY_TYPE_CCN && *result.Response.VpnGatewaySet[0].NetworkInstanceId != "" {
				return resource.NonRetryableError(fmt.Errorf("Delete operation is unsupported when VPN gateway is attached to CCN instance."))
			}
			return nil
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

	tErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(tRequest)

		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, tRequest.GetAction(), tRequest.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			if len(result.Response.VpnConnectionSet) == 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("There is associated tunnel exists, please delete associated tunnels first."))
			}
		}
	})
	if tErr != nil {
		log.Printf("[CRITAL]%s describe VPN connection failed, reason:%s\n", logId, tErr.Error())
		return tErr
	}

	request := vpc.NewDeleteVpnGatewayRequest()
	request.VpnGatewayId = &gatewayId

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteVpnGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
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
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(statRequest)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(e)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
		} else {
			//if not, quit
			if len(result.Response.VpnGatewaySet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
