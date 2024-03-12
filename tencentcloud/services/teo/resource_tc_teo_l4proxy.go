package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoL4proxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoL4proxyCreate,
		Read:   resourceTencentCloudTeoL4proxyRead,
		Update: resourceTencentCloudTeoL4proxyUpdate,
		Delete: resourceTencentCloudTeoL4proxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Zone ID.",
			},

			"proxy_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Layer 4 proxy instance name. You can enter 1-50 characters. Valid characters are a-z, 0-9, and hyphens (-). However, hyphens (-) cannot be used individually or consecutively and should not be placed at the beginning or end of the name. Modifications are not allowed after creation.",
			},

			"area": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Acceleration zone of the Layer 4 proxy instance.<li>mainland: Availability zone in the Chinese mainland;</li><li>overseas: Global availability zone (excluding the Chinese mainland);</li><li>global: Global availability zone.</li>",
			},

			"ipv6": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies whether to enable IPv6 access. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values:<li>on: Enable;</li>\n<li>off: Disable.</li>\n",
			},

			"static_ip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies whether to enable the fixed IP address. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values:<li>on: Enable;</li>\n<li>off: Disable.</li>\n",
			},

			"accelerate_mainland": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies whether to enable network optimization in the Chinese mainland. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values:<li>on: Enable;</li>\n<li>off: Disable.</li>\n",
			},

			"ddos_protection_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Layer 3/Layer 4 DDoS protection. The default protection option of the platform will be used if it is left empty. For details, see [Exclusive DDoS Protection Usage](https://intl.cloud.tencent.com/document/product/1552/95994?from_cn_redirect=1).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level_mainland": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exclusive DDoS protection specifications in the Chinese mainland. For details, see [Dedicated DDoS Mitigation Fee (Pay-as-You-Go)] (https://intl.cloud.tencent.com/document/product/1552/94162?from_cn_redirect=1).<li>PLATFORM: Default protection of the platform, i.e., Exclusive DDoS protection is not enabled;</li>\n<li>BASE30_MAX300: Exclusive DDoS protection enabled, providing a baseline protection bandwidth of 30 Gbps and an elastic protection bandwidth of up to 300 Gbps;</li><li>BASE60_MAX600: Exclusive DDoS protection enabled, providing a baseline protection bandwidth of 60 Gbps and an elastic protection bandwidth of up to 600 Gbps.</li>If no parameters are filled, the default value PLATFORM is used.",
						},
						"max_bandwidth_mainland": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Configuration of elastic protection bandwidth for exclusive DDoS protection in the Chinese mainland.Valid only when exclusive DDoS protection in the Chinese mainland is enabled (refer to the LevelMainland parameter configuration), and the value has the following limitations:<li>When exclusive DDoS protection is enabled in the Chinese mainland and the 30 Gbps baseline protection bandwidth is used (the LevelMainland parameter value is BASE30_MAX300): the value range is 30 to 300 in Gbps;</li><li>When exclusive DDoS protection is enabled in the Chinese mainland and the 60 Gbps baseline protection bandwidth is used (the LevelMainland parameter value is BASE60_MAX600): the value range is 60 to 600 in Gbps;</li><li>When the default protection of the platform is used (the LevelMainland parameter value is PLATFORM): configuration is not supported, and the value of this parameter is invalid.</li>",
						},
						"level_overseas": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exclusive DDoS protection specifications in the worldwide region (excluding the Chinese mainland).<li>PLATFORM: Default protection of the platform, i.e., Exclusive DDoS protection is not enabled;</li><li>ANYCAST300: Exclusive DDoS protection enabled, offering a total maximum protection bandwidth of 300 Gbps;</li>\n<li>ANYCAST_ALLIN: Exclusive DDoS protection enabled, utilizing all available protection resources for protection.</li>When no parameters are filled, the default value PLATFORM is used.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoL4proxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = teo.NewCreateL4ProxyRequest()
		response = teo.NewCreateL4ProxyResponse()
		zoneId   string
		proxyId  string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_name"); ok {
		request.ProxyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ipv6"); ok {
		request.Ipv6 = helper.String(v.(string))
	}

	if v, ok := d.GetOk("static_ip"); ok {
		request.StaticIp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("accelerate_mainland"); ok {
		request.AccelerateMainland = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ddos_protection_config"); ok {
		dDosProtectionConfig := teo.DDosProtectionConfig{}
		if v, ok := dMap["level_mainland"]; ok {
			dDosProtectionConfig.LevelMainland = helper.String(v.(string))
		}
		if v, ok := dMap["max_bandwidth_mainland"]; ok {
			dDosProtectionConfig.MaxBandwidthMainland = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["level_overseas"]; ok {
			dDosProtectionConfig.LevelOverseas = helper.String(v.(string))
		}
		request.DDosProtectionConfig = &dDosProtectionConfig
	}


	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateL4Proxy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo l4proxy failed, reason:%+v", logId, err)
		return err
	}

	proxyId = *response.Response.ProxyId
	d.SetId(strings.Join([]string{zoneId, proxyId}, tccommon.FILED_SP))

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"online"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TeoL4proxyStateRefreshFunc(zoneId, proxyId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoL4proxyRead(d, meta)
}

func resourceTencentCloudTeoL4proxyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	l4proxy, err := service.DescribeTeoL4proxyById(ctx, zoneId, proxyId)
	if err != nil {
		return err
	}

	if l4proxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoL4proxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if l4proxy.ZoneId != nil {
		_ = d.Set("zone_id", l4proxy.ZoneId)
	}

	if l4proxy.ProxyName != nil {
		_ = d.Set("proxy_name", l4proxy.ProxyName)
	}

	if l4proxy.Area != nil {
		_ = d.Set("area", l4proxy.Area)
	}

	if l4proxy.Ipv6 != nil {
		_ = d.Set("ipv6", l4proxy.Ipv6)
	}

	if l4proxy.StaticIp != nil {
		_ = d.Set("static_ip", l4proxy.StaticIp)
	}

	if l4proxy.AccelerateMainland != nil {
		_ = d.Set("accelerate_mainland", l4proxy.AccelerateMainland)
	}

	if l4proxy.DDosProtectionConfig != nil {
		dDosProtectionConfigMap := map[string]interface{}{}

		if l4proxy.DDosProtectionConfig.LevelMainland != nil {
			dDosProtectionConfigMap["level_mainland"] = l4proxy.DDosProtectionConfig.LevelMainland
		}

		if l4proxy.DDosProtectionConfig.MaxBandwidthMainland != nil {
			dDosProtectionConfigMap["max_bandwidth_mainland"] = l4proxy.DDosProtectionConfig.MaxBandwidthMainland
		}

		if l4proxy.DDosProtectionConfig.LevelOverseas != nil {
			dDosProtectionConfigMap["level_overseas"] = l4proxy.DDosProtectionConfig.LevelOverseas
		}

		_ = d.Set("ddos_protection_config", []interface{}{dDosProtectionConfigMap})
	}

	return nil
}

func resourceTencentCloudTeoL4proxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = teo.NewModifyL4ProxyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	immutableArgs := []string{"proxy_name", "area", "static_ip", "ddos_protection_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
		}
	}

	if d.HasChange("ipv6") {
		if v, ok := d.GetOk("ipv6"); ok {
			request.Ipv6 = helper.String(v.(string))
		}
	}

	if d.HasChange("accelerate_mainland") {
		if v, ok := d.GetOk("accelerate_mainland"); ok {
			request.AccelerateMainland = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyL4Proxy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo l4proxy failed, reason:%+v", logId, err)
		return err
	}

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"online"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TeoL4proxyStateRefreshFunc(zoneId, proxyId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoL4proxyRead(d, meta)
}

func resourceTencentCloudTeoL4proxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l4proxy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	if err := service.DeleteTeoL4proxyById(ctx, zoneId, proxyId); err != nil {
		return err
	}

	return nil
}
