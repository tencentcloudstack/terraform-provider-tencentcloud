/*
Provides a resource to create a VPN gateway.

Example Usage

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

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	errors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
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
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{5, 10, 20, 50, 100}),
				Description:  "The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100. Default is 5.",
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public ip of the VPN gateway.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of gateway instance, valid values are `IPSEC`, `SSL`.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the VPN gateway, valid values are `PENDING`, `DELETING`, `AVAILABLE`.",
			},
			"prepaid_renew_flag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Flag indicates whether to renew or not, valid values are `NOTIFY_AND_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `NOT_NOTIFY_AND_NOT_RENEW`.",
			},
			"prepaid_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Period of instance to be prepaid. Valid values are 1, 2, 3, 4, 6, 7, 8, 9, 12, 24, 36 and unit is month.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Computed:    true,
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
	ctx := context.WithValue(context.TODO(), "logId", logId)

	request := vpc.NewCreateVpnGatewayRequest()
	request.VpnGatewayName = stringToPointer(d.Get("name").(string))
	bandwidth := d.Get("bandwidth").(int)
	bandwidth64 := uint64(bandwidth)
	request.InternetMaxBandwidthOut = &bandwidth64
	request.Zone = stringToPointer(d.Get("zone").(string))
	request.VpcId = stringToPointer(d.Get("vpc_id").(string))

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
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if response.Response.VpnGateway == nil {
		d.SetId("")
		return fmt.Errorf("VPN gateway id is nil")
	}
	gatewayId := *response.Response.VpnGateway.VpnGatewayId
	d.SetId(gatewayId)

	// must wait for creating gateway finished
	statRequest := vpc.NewDescribeVpnGatewaysRequest()
	statRequest.VpnGatewayIds = []*string{stringToPointer(gatewayId)}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
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
		log.Printf("[CRITAL]%s create VPN gateway failed, reason:%s\n ", logId, err.Error())
		return err
	}

	//modify tags
	if tags := getTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:vpngw/%s", region, gatewayId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnGatewayRead(d, meta)
}

func resourceTencentCloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	gatewayId := d.Id()
	request := vpc.NewDescribeVpnGatewaysRequest()
	request.VpnGatewayIds = []*string{&gatewayId}
	var response *vpc.DescribeVpnGatewaysResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnGateways(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if len(response.Response.VpnGatewaySet) < 1 {
		return fmt.Errorf("VPN gateway id is nil")
	}

	gateway := response.Response.VpnGatewaySet[0]

	d.Set("name", *gateway.VpnGatewayName)
	d.Set("public_ip_address", *gateway.PublicIpAddress)
	d.Set("bandwidth", int(*gateway.InternetMaxBandwidthOut))
	d.Set("type", *gateway.Type)
	d.Set("create_time", *gateway.CreatedTime)
	d.Set("state", *gateway.Type)
	d.Set("prepaid_renew_flag", *gateway.RenewFlag)
	d.Set("charge_type", *gateway.InstanceChargeType)
	d.Set("expired_time", *gateway.ExpiredTime)
	d.Set("is_address_blocked", *gateway.IsAddressBlocked)
	d.Set("new_purchase_plan", *gateway.NewPurchasePlan)
	d.Set("restrict_state", *gateway.RestrictState)
	d.Set("zone", *gateway.Zone)
	//tags
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpngw", region, gatewayId)
	if err != nil {
		return err
	}
	d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	d.Partial(true)
	gatewayId := d.Id()

	if d.HasChange("name") {
		request := vpc.NewModifyVpnGatewayAttributeRequest()
		request.VpnGatewayId = &gatewayId
		request.VpnGatewayName = stringToPointer(d.Get("name").(string))
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
			log.Printf("[CRITAL]%s modify VPN gateway name failed, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("name")
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
			log.Printf("[CRITAL]%s modify VPN gateway bandwidth failed, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("bandwidth")
	}

	//tag
	if d.HasChange("tags") {
		old, new := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(old.(map[string]interface{}), new.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:vpngw/%s", region, gatewayId)
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

	//check the vpn gateway is not related with any tunnels
	tRequest := vpc.NewDescribeVpnConnectionsRequest()
	tRequest.Filters = make([]*vpc.Filter, 0, 2)
	params := make(map[string]string)
	params["vpn-gateway-id"] = gatewayId

	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}

	for k, v := range params {
		filter := &vpc.Filter{
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
		}
		tRequest.Filters = append(tRequest.Filters, filter)
	}
	offset := uint64(0)
	tRequest.Offset = &offset

	tErr := resource.Retry(3*time.Minute, func() *resource.RetryError {
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
		log.Printf("[CRITAL]%s create VPN connection failed, reason:%s\n", logId, tErr.Error())
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
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
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
