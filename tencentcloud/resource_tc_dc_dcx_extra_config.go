/*
Provides a resource to create a dc dcx_extra_config

Example Usage

```hcl
resource "tencentcloud_dc_dcx_extra_config" "dcx_extra_config" {
  direct_connect_tunnel_id = "dcx-test123"
  vlan = 123
  bgp_peer {
		asn = 65101
		auth_key = "test123"

  }
  route_filter_prefixes {
		cidr = "192.168.0.0/24"

  }
  tencent_address = "192.168.1.1"
  tencent_backup_address = "192.168.1.2"
  customer_address = "192.168.1.4"
  bandwidth = 10M
  enable_b_g_p_community = false
  bfd_enable = false
  nqa_enable = false
  bfd_info {
		probe_failed_times = 3
		interval = 100

  }
  nqa_info {
		probe_failed_times = 3
		interval = 100
		destination_ip = "192.168.2.2"

  }
  i_pv6_enable = 0
  customer_i_d_c_routes {
		cidr = ""

  }
  jumbo_enable = 0
}
```

Import

dc dcx_extra_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_dcx_extra_config.dcx_extra_config dcx_extra_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcDcxExtraConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcDcxExtraConfigCreate,
		Read:   resourceTencentCloudDcDcxExtraConfigRead,
		Update: resourceTencentCloudDcDcxExtraConfigUpdate,
		Delete: resourceTencentCloudDcDcxExtraConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_tunnel_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Direct connect tunnel id.",
			},

			"vlan": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Direct connect tunnel vlan id.",
			},

			"bgp_peer": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Idc BGP, Asn, AuthKey.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "User idc BGP Asn.",
						},
						"auth_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User bgp key.",
						},
					},
				},
			},

			"route_filter_prefixes": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "User filter network prefixes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User network prefixes.",
						},
					},
				},
			},

			"tencent_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Direct connect tunnel tencent cloud connect ip.",
			},

			"tencent_backup_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Direct connect tunnel tencent cloud backup connect ip.",
			},

			"customer_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Direct connect tunnel user idc connect ip.",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Direct connect tunnel bandwidth.",
			},

			"enable_b_g_p_community": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "BGP community attribute.",
			},

			"bfd_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Be enabled BFD.",
			},

			"nqa_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Be enabled NQA.",
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
							Description: "Detect times.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Detect interval.",
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
							Description: "Detect times.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Detect interval.",
						},
						"destination_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detect ip.",
						},
					},
				},
			},

			"i_pv6_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0: disable IPv61: enable IPv6.",
			},

			"customer_i_d_c_routes": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Idc route.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User idc network prefix.",
						},
					},
				},
			},

			"jumbo_enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Direct connect tunnel support jumbo frame1: enable direct connect tunnel jumbo frame0: disable direct connect tunnel jumbo frame.",
			},
		},
	}
}

func resourceTencentCloudDcDcxExtraConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_dcx_extra_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request               = dc.NewModifyDirectConnectTunnelExtraRequest()
		response              = dc.NewModifyDirectConnectTunnelExtraResponse()
		directConnectTunnelId string
	)
	if v, ok := d.GetOk("direct_connect_tunnel_id"); ok {
		directConnectTunnelId = v.(string)
		request.DirectConnectTunnelId = helper.String(v.(string))
	}

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

	if v, ok := d.GetOkExists("enable_b_g_p_community"); ok {
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

	if v, ok := d.GetOkExists("i_pv6_enable"); ok {
		request.IPv6Enable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("customer_i_d_c_routes"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			routeFilterPrefix := dc.RouteFilterPrefix{}
			if v, ok := dMap["cidr"]; ok {
				routeFilterPrefix.Cidr = helper.String(v.(string))
			}
			request.CustomerIDCRoutes = append(request.CustomerIDCRoutes, &routeFilterPrefix)
		}
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dc dcxExtraConfig failed, reason:%+v", logId, err)
		return err
	}

	directConnectTunnelId = *response.Response.DirectConnectTunnelId
	d.SetId(directConnectTunnelId)

	return resourceTencentCloudDcDcxExtraConfigRead(d, meta)
}

func resourceTencentCloudDcDcxExtraConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_dcx_extra_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	dcxExtraConfigId := d.Id()

	dcxExtraConfig, err := service.DescribeDcDcxExtraConfigById(ctx, directConnectTunnelId)
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

		if dcxExtraConfig.RouteFilterPrefixes.Cidr != nil {
			routeFilterPrefixesMap["cidr"] = dcxExtraConfig.RouteFilterPrefixes.Cidr
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
		_ = d.Set("enable_b_g_p_community", dcxExtraConfig.EnableBGPCommunity)
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
		_ = d.Set("i_pv6_enable", dcxExtraConfig.IPv6Enable)
	}

	if dcxExtraConfig.CustomerIDCRoutes != nil {
		customerIDCRoutesList := []interface{}{}
		for _, customerIDCRoutes := range dcxExtraConfig.CustomerIDCRoutes {
			customerIDCRoutesMap := map[string]interface{}{}

			if dcxExtraConfig.CustomerIDCRoutes.Cidr != nil {
				customerIDCRoutesMap["cidr"] = dcxExtraConfig.CustomerIDCRoutes.Cidr
			}

			customerIDCRoutesList = append(customerIDCRoutesList, customerIDCRoutesMap)
		}

		_ = d.Set("customer_i_d_c_routes", customerIDCRoutesList)

	}

	if dcxExtraConfig.JumboEnable != nil {
		_ = d.Set("jumbo_enable", dcxExtraConfig.JumboEnable)
	}

	return nil
}

func resourceTencentCloudDcDcxExtraConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_dcx_extra_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dc.NewModifyDirectConnectTunnelExtraRequest()

	dcxExtraConfigId := d.Id()

	request.DirectConnectTunnelId = &directConnectTunnelId

	immutableArgs := []string{"direct_connect_tunnel_id", "vlan", "bgp_peer", "route_filter_prefixes", "tencent_address", "tencent_backup_address", "customer_address", "bandwidth", "enable_b_g_p_community", "bfd_enable", "nqa_enable", "bfd_info", "nqa_info", "i_pv6_enable", "customer_i_d_c_routes", "jumbo_enable"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
		log.Printf("[CRITAL]%s update dc dcxExtraConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcDcxExtraConfigRead(d, meta)
}

func resourceTencentCloudDcDcxExtraConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_dcx_extra_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}
	dcxExtraConfigId := d.Id()

	if err := service.DeleteDcDcxExtraConfigById(ctx, directConnectTunnelId); err != nil {
		return err
	}

	return nil
}
