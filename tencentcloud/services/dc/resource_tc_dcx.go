package dc

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcxInstance() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "ID of the DC to be queried, application deployment offline.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the dedicated tunnel.",
			},
			"dc_owner_account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Connection owner, who is the current customer by default. The developer account ID should be entered for shared connections.",
			},
			"network_region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Network region.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_NETWORK_TYPE_VPC,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DC_NETWORK_TYPES),
				Description:  "Type of the network. Valid value: `VPC`, `BMVPC` and `CCN`. The default value is `VPC`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the VPC or BMVPC.",
			},
			"route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_ROUTE_TYPE_BGP,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DC_ROUTE_TYPES),
				Description:  "Type of the route, and available values include BGP and STATIC. The default value is `BGP`.",
			},
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the DC Gateway. Currently only new in the console.",
			},
			"bgp_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "BGP ASN of the user. A required field within BGP.",
			},
			"bgp_auth_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "BGP key of the user.",
			},
			"route_filter_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return helper.HashString(v.(string))
				},
				Description: "Static route, the network address of the user IDC. It can be modified after setting but cannot be deleted. AN unable field within BGP.",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Vlan of the dedicated tunnels. Valid value ranges: (0~3000). `0` means that only one tunnel can be created for the physical connect.",
			},
			"tencent_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Interconnect IP of the DC within Tencent.",
			},
			"customer_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Interconnect IP of the DC within client.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Bandwidth of the DC.",
			},
			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the dedicated tunnels. Valid value: `PENDING`, `ALLOCATING`, `ALLOCATED`, `ALTERING`, `DELETING`, `DELETED`, `COMFIRMING` and `REJECTED`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudDcxInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		dcId                      = d.Get("dc_id").(string)
		name                      = d.Get("name").(string)
		dcOwnerAccount            = ""
		networkType               = d.Get("network_type").(string)
		networkRegion             = service.client.Region
		vpcId                     = ""
		routeType                 = strings.ToUpper(d.Get("route_type").(string))
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

	bgpAsnTemp, bgpAsnOk := d.GetOkExists("bgp_asn")
	bgpKeyTemp, bgpKeyOk := d.GetOkExists("bgp_auth_key")
	if bgpKeyOk && !bgpAsnOk {
		return fmt.Errorf("bgp_auth_key need bgp_asn set")
	}
	if bgpAsnOk {
		bgpAsn = int64(bgpAsnTemp.(int))
	}
	if bgpKeyOk {
		bgpAuthKey = bgpKeyTemp.(string)
	}

	if temp, ok := d.GetOk("dc_owner_account"); ok {
		dcOwnerAccount = temp.(string)
	}

	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
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

	if routeType == DC_ROUTE_TYPE_BGP && len(routeFilterPrefixes) > 0 {
		return fmt.Errorf("can not set `route_filter_prefixes` if `route_type` is '%s'", DC_ROUTE_TYPE_BGP)
	}

	if routeType != DC_ROUTE_TYPE_BGP && bgpAsn != -1 {
		return fmt.Errorf("can not set `bgp_asn`,`bgp_auth_key` if  `route_type` is not '%s'", DC_ROUTE_TYPE_BGP)
	}

	dcxId, err := service.CreateDirectConnectTunnel(ctx, dcId, name, networkType,
		networkRegion, vpcId, routeType,
		bgpAuthKey, tencentAddress,
		customerAddress, dcgId, dcOwnerAccount,
		bgpAsn, vlan,
		bandwidth, routeFilterPrefixes)
	if err != nil {
		return err
	}

	err = service.waitCreateDirectConnectTunnelAvailable(ctx, dcxId)
	if err != nil {
		return err
	}
	d.SetId(dcxId)

	return resourceTencentCloudDcxInstanceRead(d, meta)
}

func resourceTencentCloudDcxInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		dcxId = d.Id()
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		item, has, e := service.DescribeDirectConnectTunnel(ctx, dcxId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("dc_id", service.strPt2str(item.DirectConnectId))
		_ = d.Set("name", *item.DirectConnectTunnelName)
		_ = d.Set("network_type", strings.ToUpper(service.strPt2str(item.NetworkType)))
		_ = d.Set("network_region", service.strPt2str(item.NetworkRegion))
		_ = d.Set("vpc_id", service.strPt2str(item.VpcId))
		_ = d.Set("bandwidth", service.int64Pt2int64(item.Bandwidth))

		var routeType = strings.ToUpper(service.strPt2str(item.RouteType))
		_ = d.Set("route_type", routeType)

		if routeType == DC_ROUTE_TYPE_BGP {
			if item.BgpPeer == nil {
				_ = d.Set("bgp_asn", 0)
				_ = d.Set("bgp_auth_key", "")
			} else {
				_ = d.Set("bgp_asn", service.int64Pt2int64(item.BgpPeer.Asn))
				_ = d.Set("bgp_auth_key", service.strPt2str(item.BgpPeer.AuthKey))
			}
		} else {
			var routeFilterPrefixes = make([]string, 0, len(item.RouteFilterPrefixes))
			for _, v := range item.RouteFilterPrefixes {
				if v.Cidr != nil {
					routeFilterPrefixes = append(routeFilterPrefixes, *v.Cidr)
				}
			}
			_ = d.Set("route_filter_prefixes", routeFilterPrefixes)
		}

		_ = d.Set("vlan", service.int64Pt2int64(item.Vlan))
		_ = d.Set("tencent_address", service.strPt2str(item.TencentAddress))
		_ = d.Set("customer_address", service.strPt2str(item.CustomerAddress))
		_ = d.Set("dcg_id", service.strPt2str(item.DirectConnectGatewayId))

		_ = d.Set("state", strings.ToUpper(service.strPt2str(item.State)))
		_ = d.Set("create_time", service.strPt2str(item.CreatedTime))
		_ = d.Set("dc_owner_account", service.strPt2str(item.DirectConnectOwnerAccount))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudDcxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		dcxId = d.Id()
		name  = ""
	)

	if d.HasChange("name") {
		name = d.Get("name").(string)
	}

	err := service.ModifyDirectConnectTunnelAttribute(ctx, dcxId, name, "", "", "", -1, -1, nil)
	if err != nil {
		return err
	}

	return resourceTencentCloudDcxInstanceRead(d, meta)
}

func resourceTencentCloudDcxInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		dcxId = d.Id()
	)

	return service.DeleteDirectConnectTunnel(ctx, dcxId)
}
