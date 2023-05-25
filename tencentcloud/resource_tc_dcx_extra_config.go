/*
Provides a resource to create a dc dcx_extra_config

Example Usage

```hcl
resource "tencentcloud_dcx_extra_config" "dcx_extra_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  vlan                     = 123
  bgp_peer {
    asn      = 65101
    auth_key = "test123"

  }
  route_filter_prefixes {
    cidr = "192.168.0.0/24"
  }
  tencent_address        = "192.168.1.1"
  tencent_backup_address = "192.168.1.2"
  customer_address       = "192.168.1.4"
  bandwidth              = 10
  enable_bgp_community   = false
  bfd_enable             = 0
  nqa_enable             = 1
  bfd_info {
    probe_failed_times = 3
    interval           = 100

  }
  nqa_info {
    probe_failed_times = 3
    interval           = 100
    destination_ip     = "192.168.2.2"

  }
  ipv6_enable = 0
  jumbo_enable = 0
}
```

Import

dc dcx_extra_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcx_extra_config.dcx_extra_config dcx_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcxExtraConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcxExtraConfigCreate,
		Read:   resourceTencentCloudDcxExtraConfigRead,
		Update: resourceTencentCloudDcxExtraConfigUpdate,
		Delete: resourceTencentCloudDcxExtraConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_tunnel_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "direct connect tunnel id.",
			},

			"vlan": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "direct connect tunnel vlan id.",
			},

			"bgp_peer": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "idc BGP, Asn, AuthKey.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "user idc BGP Asn.",
						},
						"auth_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "user bgp key.",
						},
					},
				},
			},

			"route_filter_prefixes": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "user filter network prefixes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "user network prefixes.",
						},
					},
				},
			},

			"tencent_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "direct connect tunnel tencent cloud connect ip.",
			},

			"tencent_backup_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "direct connect tunnel tencent cloud backup connect ip.",
			},

			"customer_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "direct connect tunnel user idc connect ip.",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "direct connect tunnel bandwidth.",
			},

			"enable_bgp_community": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "BGP community attribute.",
			},

			"bfd_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "be enabled BFD.",
			},

			"nqa_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "be enabled NQA.",
			},

			"bfd_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "BFD config info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"probe_failed_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "detect times.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "detect interval.",
						},
					},
				},
			},

			"nqa_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "NQA config info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"probe_failed_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "detect times.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "detect interval.",
						},
						"destination_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "detect ip.",
						},
					},
				},
			},

			"ipv6_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0: disable IPv61: enable IPv6.",
			},

			"jumbo_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "direct connect tunnel support jumbo frame1: enable direct connect tunnel jumbo frame0: disable direct connect tunnel jumbo frame.",
			},
		},
	}
}

func resourceTencentCloudDcxExtraConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcx_extra_config.create")()
	defer inconsistentCheck(d, meta)()

	directConnectTunnelId := d.Get("direct_connect_tunnel_id").(string)

	d.SetId(directConnectTunnelId)

	return resourceTencentCloudDcxExtraConfigUpdate(d, meta)
}

func resourceTencentCloudDcxExtraConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcx_extra_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	directConnectTunnelId := d.Id()

	dcxExtraConfig, err := service.DescribeDcxExtraConfigById(ctx, directConnectTunnelId)
	if err != nil {
		return err
	}

	if dcxExtraConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcDcxExtraConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dcxExtraConfig.DirectConnectTunnelId != nil {
		_ = d.Set("direct_connect_tunnel_id", dcxExtraConfig.DirectConnectTunnelId)
	}

	if dcxExtraConfig.Vlan != nil {
		_ = d.Set("vlan", dcxExtraConfig.Vlan)
	}

	if dcxExtraConfig.BgpPeer != nil {
		bgpPeerMap := map[string]interface{}{}

		if dcxExtraConfig.BgpPeer.Asn != nil {
			bgpPeerMap["asn"] = dcxExtraConfig.BgpPeer.Asn
		}

		if dcxExtraConfig.BgpPeer.AuthKey != nil {
			bgpPeerMap["auth_key"] = dcxExtraConfig.BgpPeer.AuthKey
		}

		_ = d.Set("bgp_peer", []interface{}{bgpPeerMap})
	}

	if dcxExtraConfig.RouteFilterPrefixes != nil {
		routeFilterPrefixesMap := map[string]interface{}{}

		if dcxExtraConfig.RouteFilterPrefixes != nil {
			if len(dcxExtraConfig.RouteFilterPrefixes) > 0 {
				routeFilterPrefixesMap["cidr"] = dcxExtraConfig.RouteFilterPrefixes[0].Cidr
			}
		}
		_ = d.Set("route_filter_prefixes", []interface{}{routeFilterPrefixesMap})
	}

	if dcxExtraConfig.TencentAddress != nil {
		_ = d.Set("tencent_address", dcxExtraConfig.TencentAddress)
	}

	if dcxExtraConfig.TencentBackupAddress != nil {
		_ = d.Set("tencent_backup_address", dcxExtraConfig.TencentBackupAddress)
	}

	if dcxExtraConfig.CustomerAddress != nil {
		_ = d.Set("customer_address", dcxExtraConfig.CustomerAddress)
	}

	if dcxExtraConfig.Bandwidth != nil {
		_ = d.Set("bandwidth", dcxExtraConfig.Bandwidth)
	}

	if dcxExtraConfig.EnableBGPCommunity != nil {
		_ = d.Set("enable_bgp_community", dcxExtraConfig.EnableBGPCommunity)
	}

	if dcxExtraConfig.BfdEnable != nil {
		_ = d.Set("bfd_enable", dcxExtraConfig.BfdEnable)
	}

	if dcxExtraConfig.NqaEnable != nil {
		_ = d.Set("nqa_enable", dcxExtraConfig.NqaEnable)
	}

	if dcxExtraConfig.BfdInfo != nil {
		bfdInfoMap := map[string]interface{}{}

		if dcxExtraConfig.BfdInfo.ProbeFailedTimes != nil {
			bfdInfoMap["probe_failed_times"] = dcxExtraConfig.BfdInfo.ProbeFailedTimes
		}

		if dcxExtraConfig.BfdInfo.Interval != nil {
			bfdInfoMap["interval"] = dcxExtraConfig.BfdInfo.Interval
		}

		_ = d.Set("bfd_info", []interface{}{bfdInfoMap})
	}

	if dcxExtraConfig.NqaInfo != nil {
		nqaInfoMap := map[string]interface{}{}

		if dcxExtraConfig.NqaInfo.ProbeFailedTimes != nil {
			nqaInfoMap["probe_failed_times"] = dcxExtraConfig.NqaInfo.ProbeFailedTimes
		}

		if dcxExtraConfig.NqaInfo.Interval != nil {
			nqaInfoMap["interval"] = dcxExtraConfig.NqaInfo.Interval
		}

		if dcxExtraConfig.NqaInfo.DestinationIp != nil {
			nqaInfoMap["destination_ip"] = dcxExtraConfig.NqaInfo.DestinationIp
		}

		_ = d.Set("nqa_info", []interface{}{nqaInfoMap})
	}

	if dcxExtraConfig.IPv6Enable != nil {
		_ = d.Set("ipv6_enable", dcxExtraConfig.IPv6Enable)
	}

	if dcxExtraConfig.JumboEnable != nil {
		_ = d.Set("jumbo_enable", dcxExtraConfig.JumboEnable)
	}

	return nil
}

func resourceTencentCloudDcxExtraConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcx_extra_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request               = dc.NewModifyDirectConnectTunnelExtraRequest()
		directConnectTunnelId string
	)

	directConnectTunnelId = d.Id()

	request.DirectConnectTunnelId = helper.String(directConnectTunnelId)

	if v, ok := d.GetOkExists("vlan"); ok {
		request.Vlan = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "bgp_peer"); ok {
		bgpPeer := dc.BgpPeer{}
		if v, ok := dMap["asn"]; ok {
			bgpPeer.Asn = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["auth_key"]; ok {
			bgpPeer.AuthKey = helper.String(v.(string))
		}
		request.BgpPeer = &bgpPeer
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "route_filter_prefixes"); ok {
		routeFilterPrefix := dc.RouteFilterPrefix{}
		if v, ok := dMap["cidr"]; ok {
			routeFilterPrefix.Cidr = helper.String(v.(string))
		}
		request.RouteFilterPrefixes = &routeFilterPrefix
	}

	if v, ok := d.GetOk("tencent_address"); ok {
		request.TencentAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tencent_backup_address"); ok {
		request.TencentBackupAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_address"); ok {
		request.CustomerAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_bgp_community"); ok {
		request.EnableBGPCommunity = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("bfd_enable"); ok {
		request.BfdEnable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("nqa_enable"); ok {
		request.NqaEnable = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "bfd_info"); ok {
		bFDInfo := dc.BFDInfo{}
		if v, ok := dMap["probe_failed_times"]; ok {
			bFDInfo.ProbeFailedTimes = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["interval"]; ok {
			bFDInfo.Interval = helper.IntInt64(v.(int))
		}
		request.BfdInfo = &bFDInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "nqa_info"); ok {
		nQAInfo := dc.NQAInfo{}
		if v, ok := dMap["probe_failed_times"]; ok {
			nQAInfo.ProbeFailedTimes = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["interval"]; ok {
			nQAInfo.Interval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["destination_ip"]; ok {
			nQAInfo.DestinationIp = helper.String(v.(string))
		}
		request.NqaInfo = &nQAInfo
	}

	if v, ok := d.GetOkExists("ipv6_enable"); ok {
		request.IPv6Enable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("jumbo_enable"); ok {
		request.JumboEnable = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().ModifyDirectConnectTunnelExtra(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dc dcxExtraConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcxExtraConfigRead(d, meta)
}

func resourceTencentCloudDcxExtraConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcx_extra_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
