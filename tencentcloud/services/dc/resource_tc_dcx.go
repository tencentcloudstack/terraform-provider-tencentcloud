package dc

import (
	"context"
	"fmt"
	"log"
	"strings"

	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
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
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_NETWORK_TYPE_VPC,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DC_NETWORK_TYPES),
				Description:  "Type of the network. Valid value: `VPC`, `BMVPC` and `CCN`. The default value is `VPC`.",
			},
			"network_region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Network region.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the VPC or BMVPC.",
			},
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the DC Gateway. Currently only new in the console.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Bandwidth of the DC.",
			},
			"route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_ROUTE_TYPE_BGP,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DC_ROUTE_TYPES),
				Description:  "Type of the route, and available values include BGP and STATIC. The default value is `BGP`.",
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

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = dc.NewCreateDirectConnectTunnelRequest()
		response = dc.NewCreateDirectConnectTunnelResponse()
	)

	if v, ok := d.GetOk("dc_id"); ok {
		request.DirectConnectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.DirectConnectTunnelName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dc_owner_account"); ok {
		request.DirectConnectOwnerAccount = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_region"); ok {
		request.NetworkRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcg_id"); ok {
		request.DirectConnectGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("route_type"); ok {
		request.RouteType = helper.String(v.(string))
	}

	var bgpPeer dc.BgpPeer
	if v, ok := d.GetOkExists("bgp_asn"); ok {
		bgpPeer.Asn = helper.IntInt64(v.(int))
		request.BgpPeer = &bgpPeer
	}

	if v, ok := d.GetOk("bgp_auth_key"); ok {
		bgpPeer.AuthKey = helper.String(v.(string))
		request.BgpPeer = &bgpPeer
	}

	if v, ok := d.GetOk("route_filter_prefixes"); ok {
		for _, item := range v.(*schema.Set).List() {
			var dcPrefix dc.RouteFilterPrefix
			dcPrefix.Cidr = helper.String(item.(string))
			request.RouteFilterPrefixes = append(request.RouteFilterPrefixes, &dcPrefix)
		}
	}

	if v, ok := d.GetOkExists("vlan"); ok {
		request.Vlan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tencent_address"); ok {
		request.TencentAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_address"); ok {
		request.CustomerAddress = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcClient().CreateDirectConnectTunnel(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DirectConnectTunnelIdSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Create direct connect tunnel failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Create direct connect tunnel failed, reason:%s\n", logId, err.Error())
		return err
	}

	if len(response.Response.DirectConnectTunnelIdSet) < 1 {
		return fmt.Errorf("DirectConnectTunnelIdSet is nil.")
	}

	dcxId := *response.Response.DirectConnectTunnelIdSet[0]
	d.SetId(dcxId)

	err = service.waitCreateDirectConnectTunnelAvailable(ctx, dcxId)
	if err != nil {
		return err
	}

	return resourceTencentCloudDcxInstanceRead(d, meta)
}

func resourceTencentCloudDcxInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dcxId   = d.Id()
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

		if item.DirectConnectId != nil {
			_ = d.Set("dc_id", service.strPt2str(item.DirectConnectId))
		}

		if item.DirectConnectTunnelName != nil {
			_ = d.Set("name", item.DirectConnectTunnelName)
		}

		if item.DirectConnectOwnerAccount != nil {
			_ = d.Set("dc_owner_account", service.strPt2str(item.DirectConnectOwnerAccount))
		}

		if item.NetworkType != nil {
			_ = d.Set("network_type", strings.ToUpper(service.strPt2str(item.NetworkType)))
		}

		if item.NetworkRegion != nil {
			_ = d.Set("network_region", service.strPt2str(item.NetworkRegion))
		}

		if item.VpcId != nil {
			_ = d.Set("vpc_id", service.strPt2str(item.VpcId))
		}

		if item.DirectConnectGatewayId != nil {
			_ = d.Set("dcg_id", service.strPt2str(item.DirectConnectGatewayId))
		}

		if item.Bandwidth != nil {
			_ = d.Set("bandwidth", service.int64Pt2int64(item.Bandwidth))
		}

		if item.RouteType != nil {
			var routeType = strings.ToUpper(service.strPt2str(item.RouteType))
			_ = d.Set("route_type", routeType)

			if routeType == DC_ROUTE_TYPE_BGP {
				if item.BgpPeer == nil {
					_ = d.Set("bgp_asn", 0)
					_ = d.Set("bgp_auth_key", "")
				} else {
					if item.BgpPeer.Asn != nil {
						_ = d.Set("bgp_asn", service.int64Pt2int64(item.BgpPeer.Asn))
					}

					if item.BgpPeer.AuthKey != nil {
						_ = d.Set("bgp_auth_key", service.strPt2str(item.BgpPeer.AuthKey))
					}
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
		}

		if item.Vlan != nil {
			_ = d.Set("vlan", service.int64Pt2int64(item.Vlan))
		}

		if item.TencentAddress != nil {
			_ = d.Set("tencent_address", service.strPt2str(item.TencentAddress))
		}

		if item.CustomerAddress != nil {
			_ = d.Set("customer_address", service.strPt2str(item.CustomerAddress))
		}

		if item.State != nil {
			_ = d.Set("state", strings.ToUpper(service.strPt2str(item.State)))
		}

		if item.CreatedTime != nil {
			_ = d.Set("create_time", service.strPt2str(item.CreatedTime))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudDcxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.update")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		dcxId = d.Id()
	)

	if d.HasChange("name") {
		request := dc.NewModifyDirectConnectTunnelAttributeRequest()
		request.DirectConnectTunnelId = &dcxId
		if v, ok := d.GetOk("name"); ok {
			request.DirectConnectTunnelName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcClient().ModifyDirectConnectTunnelAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Modify direct connect tunnel failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudDcxInstanceRead(d, meta)
}

func resourceTencentCloudDcxInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcx.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dc.NewDeleteDirectConnectTunnelRequest()
		dcxId   = d.Id()
	)

	request.DirectConnectTunnelId = &dcxId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcClient().DeleteDirectConnectTunnel(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Delete direct connect tunnel failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
