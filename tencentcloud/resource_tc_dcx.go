/*
Provides a resource to creating dedicated tunnels instances.

~> **NOTE:** 1. ID of the DC is queried, can only apply for this resource offline.

Example Usage

```hcl
variable "dc_id" {
  default = "dc-kax48sg7"
}

variable "dcg_id" {
  default = "dcg-dmbhf7jf"
}

variable "vpc_id" {
  default = "vpc-4h9v4mo3"
}

resource "tencentcloud_dcx" "bgp_main" {
  bandwidth    = 900
  dc_id        = var.dc_id
  dcg_id       = var.dcg_id
  name         = "bgp_main"
  network_type = "VPC"
  route_type   = "BGP"
  vlan         = 306
  vpc_id       = var.vpc_id
}

resource "tencentcloud_dcx" "static_main" {
  bandwidth             = 900
  dc_id                 = var.dc_id
  dcg_id                = var.dcg_id
  name                  = "static_main"
  network_type          = "VPC"
  route_type            = "STATIC"
  vlan                  = 301
  vpc_id                = var.vpc_id
  route_filter_prefixes = ["10.10.10.101/32"]
  tencent_address       = "100.93.46.1/30"
  customer_address      = "100.93.46.2/30"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description:  "ID of the DC to be queried, application deployment offline.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the dedicated tunnel.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_NETWORK_TYPE_VPC,
				ValidateFunc: validateAllowedStringValue(DC_NETWORK_TYPES),
				Description:  "Type of the network. Valid value: `VPC`, `BMVPC` and `CCN`. The default value is `VPC`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC or BMVPC.",
			},
			"route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DC_ROUTE_TYPE_BGP,
				ValidateFunc: validateAllowedStringValue(DC_ROUTE_TYPES),
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
					return hashcode.String(v.(string))
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
	defer logElapsed("resource.tencentcloud_dcx.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcId                      = d.Get("dc_id").(string)
		name                      = d.Get("name").(string)
		networkType               = d.Get("network_type").(string)
		networkRegion             = service.client.Region
		vpcId                     = d.Get("vpc_id").(string)
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
	defer logElapsed("resource.tencentcloud_dcx.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcxId = d.Id()
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		item, has, e := service.DescribeDirectConnectTunnel(ctx, dcxId)
		if e != nil {
			return retryError(e)
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
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudDcxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcx.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_dcx.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcxId = d.Id()
	)

	return service.DeleteDirectConnectTunnel(ctx, dcxId)
}
