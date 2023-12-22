package vpc

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var routeTypeApiMap = map[string]int{
	"public_gateway":     0,
	"vpn_gateway":        1,
	"dc_gateway":         3,
	"peering_connection": 4,
	"havip":              5,
	"sslvpn_gateway":     7,
	"nat_gateway":        8,
	"instance":           9,
	"eip":                10,
	"local_gateway":      11,
}

var RouteTypeApiMap = routeTypeApiMap

var routeTypeNewMap = map[string]string{
	"public_gateway":     "CVM",
	"vpn_gateway":        "VPN",
	"dc_gateway":         "DIRECTCONNECT",
	"peering_connection": "PEERCONNECTION",
	"havip":              "HAVIP",
	"sslvpn_gateway":     "SSLVPN",
	"nat_gateway":        "NAT",
	"instance":           "NORMAL_CVM",
	"eip":                "EIP",
	"local_gateway":      "LOCAL_GATEWAY",
}

var RouteTypeNewMap = routeTypeNewMap

func ResourceTencentCloudRouteEntry() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_route_table_entry' instead.",
		Create:             resourceTencentCloudRouteEntryCreate,
		Read:               resourceTencentCloudRouteEntryRead,
		Delete:             resourceTencentCloudRouteEntryDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC ID.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the route table.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCIDRNetworkAddress,
				Description:  "The RouteEntry's target network segment.",
			},
			"next_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					_, ok := routeTypeNewMap[value]
					if !ok {
						var nextHubDesc []string
						for vgwKey := range routeTypeNewMap {
							nextHubDesc = append(nextHubDesc, vgwKey)
						}
						errors = append(errors, fmt.Errorf("%s Only one of %s is allowed", k, strings.Join(nextHubDesc, ",")))
					}
					return
				},
				Description: "The next hop type. Valid values: `public_gateway`,`vpn_gateway`,`sslvpn_gateway`,`dc_gateway`,`peering_connection`,`nat_gateway`,`havip`,`local_gateway` and `instance`. `instance` points to CVM Instance.",
			},
			"next_hub": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The route entry's next hub. CVM instance ID or VPC router interface ID.",
			},
		},
	}
}

func resourceTencentCloudRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_entry.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	vpcId := d.Get("vpc_id").(string)
	routeTableId := d.Get("route_table_id").(string)
	destinationCidrBlock := d.Get("cidr_block").(string)
	nextType := d.Get("next_type").(string)
	nextHub := d.Get("next_hub").(string)

	if nextType == GATE_WAY_TYPE_EIP && nextHub != "0" {
		return fmt.Errorf("if next_type is %s, next_hub can only be \"0\" ", GATE_WAY_TYPE_EIP)
	}

	if _, ok := routeTypeNewMap[nextType]; !ok {
		return fmt.Errorf("The value of next_type is invalid")
	}

	_, err := service.CreateRoutes(ctx, routeTableId, destinationCidrBlock, routeTypeNewMap[nextType], nextHub, "", true)
	if err != nil {
		return err
	}

	route := map[string]string{
		"vpcId":                vpcId,
		"routeTableId":         routeTableId,
		"destinationCidrBlock": destinationCidrBlock,
		"nextType":             fmt.Sprintf("%d", routeTypeApiMap[nextType]),
		"nextHub":              nextHub,
	}
	uniqRouteEntryId, ok := routeIdEncode(route)
	if !ok {
		return fmt.Errorf("Failed to encode route entry")
	}

	d.SetId(uniqRouteEntryId)

	return resourceTencentCloudRouteEntryRead(d, meta)
}

func resourceTencentCloudRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_entry.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	route, ok := RouteIdDecode(d.Id())
	if !ok {
		return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", d.Id())
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		if has != 1 {
			e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return resource.NonRetryableError(e)
		}
		for _, v := range info.entryInfos {
			var nextType string
			var nextTypeId string
			for kk, vv := range routeTypeNewMap {
				if vv == v.nextType {
					nextType = kk
				}
			}
			if _, ok := routeTypeApiMap[nextType]; ok {
				nextTypeId = fmt.Sprintf("%d", routeTypeApiMap[nextType])
			}
			if v.destinationCidr == route["destinationCidrBlock"] &&
				nextTypeId == route["nextType"] &&
				v.nextBub == route["nextHub"] &&
				v.description == route["description"] {
				_ = d.Set("vpc_id", route["vpcId"])
				_ = d.Set("route_table_id", route["routeTableId"])
				_ = d.Set("cidr_block", route["destinationCidrBlock"])
				_ = d.Set("next_type", nextType)
				_ = d.Set("next_hub", route["nextHub"])
				return nil
			}
		}

		d.SetId("")
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_entry.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	route, ok := RouteIdDecode(d.Id())
	if !ok {
		return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", d.Id())
	}

	var routeEntryId int64
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		if has != 1 {
			e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return resource.NonRetryableError(e)
		}

		for _, v := range info.entryInfos {
			var nextType string
			var nextTypeId string
			for kk, vv := range routeTypeNewMap {
				if vv == v.nextType {
					nextType = kk
				}
			}
			if _, ok := routeTypeApiMap[nextType]; ok {
				nextTypeId = fmt.Sprintf("%d", routeTypeApiMap[nextType])
			}
			if v.destinationCidr == route["destinationCidrBlock"] &&
				nextTypeId == route["nextType"] &&
				v.nextBub == route["nextHub"] &&
				v.description == route["description"] {
				routeEntryId = v.routeEntryId
				return nil
			}
		}

		d.SetId("")
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteRoutes(ctx, route["routeTableId"], uint64(routeEntryId)); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}

func routeIdEncode(route map[string]string) (routeId string, ok bool) {
	vpcId, ok0 := route["vpcId"]
	rtbId, ok1 := route["routeTableId"]
	cidrBlock, ok2 := route["destinationCidrBlock"]
	nextType, ok3 := route["nextType"]
	nextHub, ok4 := route["nextHub"]
	if ok0 && ok1 && ok2 && ok3 && ok4 {
		return fmt.Sprintf("%v::%v::%v::%v::%v", vpcId, rtbId, cidrBlock, nextType, nextHub), true
	}
	return "", false
}

func RouteIdDecode(routeId string) (route map[string]string, ok bool) {
	route = map[string]string{}
	routeArray := strings.Split(routeId, "::")
	if len(routeArray) != 5 {
		return route, false
	}
	route["vpcId"] = routeArray[0]
	route["routeTableId"] = routeArray[1]
	route["destinationCidrBlock"] = routeArray[2]
	route["nextType"] = routeArray[3]
	route["nextHub"] = routeArray[4]
	return route, true
}
