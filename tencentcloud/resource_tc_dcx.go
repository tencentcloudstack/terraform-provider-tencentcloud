package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceTencentCloudDcxInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcxInstanceCreate,
		Read:   resourceTencentCloudDcxInstanceRead,
		Update: resourceTencentCloudDcxInstanceUpdate,
		Delete: resourceTencentCloudDcxInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_NETWORK_TYPE_VPC,
				ValidateFunc: validateAllowedStringValue(DC_NETWORK_TYPES),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_ROUTE_TYPE_BGP,
				ValidateFunc: validateAllowedStringValue(DC_ROUTE_TYPES),
			},
			"dcg_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bgp_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bgp_auth_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"route_filter_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"vlan": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"tencent_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"customer_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Computed values
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudDcxInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "resource.tencentcloud_dcx.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcId                      = d.Get("dc_id").(string)
		name                      = d.Get("name").(string)
		networkType               = d.Get("network_type").(string)
		networkRegion             = service.client.Region
		vpcId                     = d.Get("vpc_id").(string)
		routeType                 = d.Get("route_type").(string)
		bgpAsn              int64 = -1
		bgpAuthKey                = ""
		vlan                      = int64(d.Get("vlan").(int))
		tencentAddress            = ""
		customerAddress           = ""
		bandwidth           int64 = -1
		routeFilterPrefixes []string
		dcgId               = d.Get("dcg_id").(string)
	)

	if temp, ok := d.GetOk("network_region"); ok {
		networkRegion = temp.(string)
	}

	bgpAsnTemp, bgpAsnOk := d.GetOkExists("bgp_asn");
	bgpKeyTemp, bgpKeyOk := d.GetOkExists("bgp_auth_key");
	if (!bgpAsnOk && bgpKeyOk) || (bgpAsnOk && !bgpKeyOk) {
		return fmt.Errorf("bgp_asn and bgp_auth_key should both set or both unset")
	}
	if bgpAsnOk {
		bgpAsn = int64(bgpAsnTemp.(int))
		bgpAuthKey = bgpKeyTemp.(string)
	}
	if temp, ok := d.GetOk("tencent_address"); ok {
		tencentAddress = temp.(string)
	}
	if temp, ok := d.GetOk("customer_address"); ok {
		customerAddress = temp.(string)
	}
	if temp, ok := d.GetOk("bandwidth"); ok {
		bandwidth = int64(temp.(int))
	}

	if temp, ok := d.GetOk("route_filter_prefixes"); ok {
		for _, v := range temp.(*schema.Set).List() {
			routeFilterPrefixes = append(routeFilterPrefixes, v.(string))
		}
	}
	dcxId, err := service.CreateDirectConnectTunnel(ctx, dcId, name, networkType,
		networkRegion, vpcId, routeType,
		bgpAuthKey, tencentAddress,
		customerAddress, dcgId,
		bgpAsn, vlan,
		bandwidth, routeFilterPrefixes)
	if err != nil {
		return err
	}
	d.SetId(dcxId)
	return resourceTencentCloudDcxInstanceRead(d, meta)
}

func resourceTencentCloudDcxInstanceRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "resource.tencentcloud_dcx.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcxId = d.Id()
	)

	item, has, err := service.DescribeDirectConnectTunnel(ctx, dcxId)

	if err != nil {
		return err
	}
	if has == 0 {
		d.SetId("")
	}
	d.Set("dc_id", service.strPt2str(item.DirectConnectId))
	d.Set("name", *item.DirectConnectTunnelName)
	d.Set("network_type", strings.ToUpper(service.strPt2str(item.NetworkType)))
	d.Set("network_region", service.strPt2str(item.NetworkRegion))
	d.Set("vpc_id", service.strPt2str(item.VpcId))
	d.Set("bandwidth", service.int64Pt2int64(item.Bandwidth))
	d.Set("route_type", strings.ToUpper(service.strPt2str(item.RouteType)))

	if item.BgpPeer == nil {
		d.Set("bgp_asn", 0)
		d.Set("bgp_auth_key", "")
	} else {
		d.Set("bgp_asn", service.int64Pt2int64(item.BgpPeer.Asn))
		d.Set("bgp_auth_key", service.strPt2str(item.BgpPeer.AuthKey))
	}
	var routeFilterPrefixes = make([]string, 0, len(item.RouteFilterPrefixes))
	for _, v := range item.RouteFilterPrefixes {
		if v.Cidr != nil {
			routeFilterPrefixes = append(routeFilterPrefixes, *v.Cidr)
		}
	}
	d.Set("route_filter_prefixes", routeFilterPrefixes)

	d.Set("vlan", service.int64Pt2int64(item.Vlan))
	d.Set("tencent_address", service.strPt2str(item.TencentAddress))
	d.Set("customer_address", service.strPt2str(item.CustomerAddress))
	d.Set("dcg_id", service.strPt2str(item.DirectConnectGatewayId))

	d.Set("state", strings.ToUpper(service.strPt2str(item.State)))
	d.Set("create_time", service.strPt2str(item.CreatedTime))

	return nil
}
func resourceTencentCloudDcxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "resource.tencentcloud_dcx.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcxId                     = d.Id()
		name                      = ""
		bandwidth           int64 = -1
		routeFilterPrefixes       = make([]string, 0, 10)
		bgpAsn              int64 = -1
		bgpAuthKey                = ""
		tencentAddress            = ""
		customerAddress           = ""
	)

	if d.HasChange("name") {
		name = d.Get("name").(string)
	}
	if d.HasChange("bandwidth") {
		if temp, ok := d.GetOk("bandwidth"); ok {
			bandwidth = int64(temp.(int))
		}
	}

	if d.HasChange("route_filter_prefixes") {
		if temp, ok := d.GetOk("route_filter_prefixes"); ok {
			for _, v := range temp.(*schema.Set).List() {
				routeFilterPrefixes = append(routeFilterPrefixes, v.(string))
			}
		}
	}

	if d.HasChange("bgp_asn") || d.HasChange("bgp_auth_key") {
		bgpAsnTemp, bgpAsnOk := d.GetOkExists("bgp_asn");
		bgpKeyTemp, bgpKeyOk := d.GetOkExists("bgp_auth_key");
		if (!bgpAsnOk && bgpKeyOk) || (bgpAsnOk && !bgpKeyOk) {
			return fmt.Errorf("bgp_asn and bgp_auth_key should both set or both unset")
		}
		if bgpAsnOk {
			bgpAsn = int64(bgpAsnTemp.(int))
			bgpAuthKey = bgpKeyTemp.(string)
		}
	}

	if d.HasChange("tencent_address") {
		if temp, ok := d.GetOk("tencent_address"); ok {
			tencentAddress = temp.(string)
		}
	}
	if d.HasChange("customer_address") {
		if temp, ok := d.GetOk("customer_address"); ok {
			customerAddress = temp.(string)
		}
	}
	err := service.ModifyDirectConnectTunnelAttribute(ctx, dcxId,
		name, bgpAuthKey, tencentAddress, customerAddress,
		bandwidth, bgpAsn,
		routeFilterPrefixes)

	if err != nil {
		return err
	}
	d.SetId(dcxId)
	return resourceTencentCloudDcxInstanceRead(d, meta)
}
func resourceTencentCloudDcxInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "resource.tencentcloud_dcx.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcxId = d.Id()
	)
	return service.DeleteDirectConnectTunnel(ctx, dcxId)
}
