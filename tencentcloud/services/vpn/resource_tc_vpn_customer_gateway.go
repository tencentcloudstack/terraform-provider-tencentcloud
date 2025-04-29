package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnCustomerGatewayCreate,
		Read:   resourceTencentCloudVpnCustomerGatewayRead,
		Update: resourceTencentCloudVpnCustomerGatewayUpdate,
		Delete: resourceTencentCloudVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the customer gateway. The length of character is limited to 1-60.",
			},
			"public_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Public IP of the customer gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},
			"bgp_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "BGP ASN. Value range: 1 - 4294967295. Using BGP requires configuring ASN. 139341, 45090, and 58835 are not available.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the customer gateway.",
			},
		},
	}
}

func resourceTencentCloudVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_customer_gateway.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request  = vpc.NewCreateCustomerGatewayRequest()
		response = vpc.NewCreateCustomerGatewayResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.CustomerGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("public_ip_address"); ok {
		request.IpAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bgp_asn"); ok {
		request.BgpAsn = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateCustomerGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.CustomerGateway == nil {
			return resource.RetryableError(fmt.Errorf("Create VPN customer gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.CustomerGateway.CustomerGatewayId == nil {
		return fmt.Errorf("VPN customer gateway id is nil")
	}

	customerGatewayId := *response.Response.CustomerGateway.CustomerGatewayId
	d.SetId(customerGatewayId)

	// must wait for finishing creating customer gateway
	statRequest := vpc.NewDescribeCustomerGatewaysRequest()
	statRequest.CustomerGatewayIds = []*string{response.Response.CustomerGateway.CustomerGatewayId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeCustomerGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			if result == nil || result.Response == nil || result.Response.CustomerGatewaySet == nil {
				return resource.NonRetryableError(fmt.Errorf("response is nil"))
			}

			//if not, quit
			if len(result.Response.CustomerGatewaySet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			}

			//else consider created, cos there is no status of gateway
			return nil
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "cgw", region, customerGatewayId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnCustomerGatewayRead(d, meta)
}

func resourceTencentCloudVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_customer_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request           = vpc.NewDescribeCustomerGatewaysRequest()
		customerGatewayId = d.Id()
	)

	request.CustomerGatewayIds = []*string{&customerGatewayId}
	var response *vpc.DescribeCustomerGatewaysResponse
	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = customerGatewayId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient(iacExtInfo).DescribeCustomerGateways(request)
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
		}

		if result == nil || result.Response == nil || result.Response.CustomerGatewaySet == nil {
			return resource.NonRetryableError(fmt.Errorf("Read VPN customer gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response == nil || response.Response == nil || len(response.Response.CustomerGatewaySet) < 1 {
		d.SetId("")
		return nil
	}

	gateway := response.Response.CustomerGatewaySet[0]
	if gateway.CustomerGatewayName != nil {
		_ = d.Set("name", gateway.CustomerGatewayName)
	}

	if gateway.IpAddress != nil {
		_ = d.Set("public_ip_address", gateway.IpAddress)
	}

	if gateway.BgpAsn != nil && *gateway.BgpAsn != 0 {
		_ = d.Set("bgp_asn", gateway.BgpAsn)
	}

	if gateway.CreatedTime != nil {
		_ = d.Set("create_time", gateway.CreatedTime)
	}

	//tags
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "cgw", region, customerGatewayId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_customer_gateway.update")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		customerGatewayId = d.Id()
	)

	d.Partial(true)
	if d.HasChange("name") || d.HasChange("bgp_asn") {
		request := vpc.NewModifyCustomerGatewayAttributeRequest()
		request.CustomerGatewayId = &customerGatewayId
		if v, ok := d.GetOk("name"); ok {
			request.CustomerGatewayName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("bgp_asn"); ok {
			request.BgpAsn = helper.IntUint64(v.(int))
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyCustomerGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify VPN customer gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "cgw", region, customerGatewayId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudVpnCustomerGatewayRead(d, meta)
}

func resourceTencentCloudVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_customer_gateway.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	customerGatewayId := d.Id()

	//check the customer gateway is not related with any tunnel
	tRequest := vpc.NewDescribeVpnConnectionsRequest()
	tRequest.Filters = make([]*vpc.Filter, 0, 1)
	params := make(map[string]string)
	params["customer-gateway-id"] = customerGatewayId

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
			if result == nil || result.Response == nil || result.Response.VpnConnectionSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Read VPN connections failed, Response is nil."))
			}

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

	request := vpc.NewDeleteCustomerGatewayRequest()
	request.CustomerGatewayId = &customerGatewayId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteCustomerGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	//to get the status of customer gateway
	statRequest := vpc.NewDescribeCustomerGatewaysRequest()
	statRequest.CustomerGatewayIds = []*string{&customerGatewayId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeCustomerGateways(statRequest)
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
			if result == nil || result.Response == nil || result.Response.CustomerGatewaySet == nil {
				return resource.NonRetryableError(fmt.Errorf("Read VPN customer gateways failed, Response is nil."))
			}

			//if not, quit
			if len(result.Response.CustomerGatewaySet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
