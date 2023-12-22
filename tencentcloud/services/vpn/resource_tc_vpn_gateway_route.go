package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnGatewayRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnGatewayRouteCreate,
		Read:   resourceTencentCloudVpnGatewayRouteRead,
		Update: resourceTencentCloudVpnGatewayRouteUpdate,
		Delete: resourceTencentCloudVpnGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: VpnGatewayRoutePara(),
	}
}

func VpnGatewayRoutePara() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpn_gateway_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "VPN gateway ID.",
		},
		"destination_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Destination IDC IP range.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Next hop type (type of the associated instance). Valid values: VPNCONN (VPN tunnel) and CCN (CCN instance).",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Instance ID of the next hop.",
		},
		"priority": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Priority. Valid values: 0 and 100.",
		},
		"status": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Status. Valid values: ENABLE and DISABLE.",
		},
		"route_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route ID.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route type. Default value: Static.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Create time.",
		},
		"update_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Update time.",
		},
	}
}

func resourceTencentCloudVpnGatewayRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_route.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vpnGatewayId := d.Get("vpn_gateway_id").(string)
	priority := int64(d.Get("priority").(int))
	route := &vpc.VpnGatewayRoute{
		DestinationCidrBlock: helper.String(d.Get("destination_cidr_block").(string)),
		InstanceType:         helper.String(d.Get("instance_type").(string)),
		InstanceId:           helper.String(d.Get("instance_id").(string)),
		Priority:             &priority,
		Status:               helper.String(d.Get("status").(string)),
	}
	if routeType, ok := d.GetOk("type"); ok {
		route.Type = helper.String(routeType.(string))
	}

	vpcService := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err, routeList := vpcService.CreateVpnGatewayRoute(ctx, vpnGatewayId, []*vpc.VpnGatewayRoute{route})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN gateway route failed, reason:%s\n", logId, err.Error())
		return err
	}

	if len(routeList) == 0 {
		return fmt.Errorf("VPN gateway route id is nil")
	}
	d.SetId(helper.IdFormat(vpnGatewayId, *(routeList[0].RouteId)))

	//setRouteInfo(d, vpnGatewayId, route)
	return resourceTencentCloudVpnGatewayRouteRead(d, meta)
}

func resourceTencentCloudVpnGatewayRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_route.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id      = d.Id()
	)
	compositeId := helper.IdParse(id)
	if len(compositeId) != 2 {
		return errors.New("the id format must be '{vpn_gateway_id}#{route_id}'")
	}

	err, routeList := service.DescribeVpnGatewayRoutes(ctx, compositeId[0], nil)
	if err != nil {
		log.Printf("[CRITAL]%s read VPN gateway routes failed, reason:%s\n", logId, err.Error())
		return err
	}
	var route *vpc.VpnGatewayRoute
	for _, r := range routeList {
		if compositeId[1] == *r.RouteId {
			route = r
		}
	}
	if route == nil {
		d.SetId("")
		return nil
	}

	setRouteInfo(d, compositeId[0], route)
	return nil
}

func resourceTencentCloudVpnGatewayRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_route.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id      = d.Id()
	)
	compositeId := helper.IdParse(id)
	if len(compositeId) != 2 {
		return errors.New("the id format must be '{vpn_gateway_id}#{route_id}'")
	}

	if !d.HasChange("status") {
		return nil
	}
	status := d.Get("status").(string)

	// update
	err, route := service.ModifyVpnGatewayRoute(ctx, compositeId[0], compositeId[1], status)
	if err != nil {
		log.Printf("[CRITAL]%s modify VPN gateway route failed, reason:%s\n", logId, err.Error())
		return err
	}

	setRouteInfo(d, compositeId[0], route)
	return nil
}

func resourceTencentCloudVpnGatewayRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_route.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id      = d.Id()
	)
	compositeId := helper.IdParse(id)
	if len(compositeId) != 2 {
		return errors.New("the id format must be '{vpn_gateway_id}#{route_id}'")
	}

	err := service.DeleteVpnGatewayRoutes(ctx, compositeId[0], []*string{&compositeId[1]})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN gateway routes failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func setRouteInfo(d *schema.ResourceData, vpnGatewayId string, route *vpc.VpnGatewayRoute) {
	m := ConverterVpnGatewayRouteToMap(vpnGatewayId, route)
	for k, v := range m {
		_ = d.Set(k, v)
	}
}

func ConverterVpnGatewayRouteToMap(vpnGatewayId string, route *vpc.VpnGatewayRoute) map[string]interface{} {
	return map[string]interface{}{
		"vpn_gateway_id":         vpnGatewayId,
		"destination_cidr_block": route.DestinationCidrBlock,
		"instance_type":          route.InstanceType,
		"instance_id":            route.InstanceId,
		"priority":               route.Priority,
		"status":                 route.Status,
		"route_id":               route.RouteId,
		"type":                   route.Type,
		"create_time":            route.CreateTime,
		"update_time":            route.UpdateTime,
	}
}
