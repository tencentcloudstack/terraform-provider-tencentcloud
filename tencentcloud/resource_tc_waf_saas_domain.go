/*
Provides a resource to create a waf saas_domain

Example Usage

```hcl
resource "tencentcloud_waf_saas_domain" "saas_domain" {
  domain = ""
  cert_type =
  is_cdn =
  upstream_type =
  is_websocket =
  load_balance = ""
  cert = ""
  private_key = ""
  s_s_l_id = ""
  resource_id = ""
  upstream_scheme = ""
  https_upstream_port = ""
  is_gray =
  gray_areas =
  upstream_domain = ""
  src_list =
  is_http2 =
  https_rewrite =
  ports {
		port = ""
		protocol = ""
		upstream_port = ""
		upstream_protocol = ""
		nginx_server_id = ""

  }
  edition = ""
  is_keep_alive = ""
  instance_i_d = ""
  anycast =
  weights =
  active_check =
  t_l_s_version =
  ciphers =
  cipher_template =
  proxy_read_timeout =
  proxy_send_timeout =
  sni_type =
  sni_host = ""
  ip_headers =
  x_f_f_reset =
}
```

Import

waf saas_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_saas_domain.saas_domain saas_domain_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWafSaasDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafSaasDomainCreate,
		Read:   resourceTencentCloudWafSaasDomainRead,
		Update: resourceTencentCloudWafSaasDomainUpdate,
		Delete: resourceTencentCloudWafSaasDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain names that require defense.",
			},

			"cert_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Certificate type, 0 represents no certificate, CertType=1 represents self owned certificate, and 2 represents managed certificate.",
			},

			"is_cdn": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.",
			},

			"upstream_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Upstream type, 0 represents IP, 1 represents domain name.",
			},

			"is_websocket": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Is WebSocket support enabled. 1 means enabled, 0 does not.",
			},

			"load_balance": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Load balancing strategy, where 0 represents polling and 1 represents IP hash.",
			},

			"cert": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate content, When CertType=1, this parameter needs to be filled.",
			},

			"private_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate key, When CertType=1, this parameter needs to be filled.",
			},

			"s_s_l_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID, When CertType=2, this parameter needs to be filled.",
			},

			"resource_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Waf&amp;amp;#39;s resource ID.",
			},

			"upstream_scheme": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upstream scheme for https, http or https.",
			},

			"https_upstream_port": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upstream port for https, When listen ports has https port and UpstreamScheme is HTTP, the current field needs to be filled.",
			},

			"is_gray": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to turn on grayscale, 0 indicates not to turn on grayscale.",
			},

			"gray_areas": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Grayscale region.",
			},

			"upstream_domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upstream domain, When UpstreamType=1, this parameter needs to be filled.",
			},

			"src_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Upstream IP List, When UpstreamType=0, this parameter needs to be filled.",
			},

			"is_http2": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether enable HTTP2, Enabling HTTP2 requires HTTPS support.",
			},

			"https_rewrite": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether redirect to https, 1 will redirect and 0 will not.",
			},

			"ports": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "This field needs to be set for multiple ports in the upstream server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listening port.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The listening protocol of listening port.",
						},
						"upstream_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The upstream port for listening port.",
						},
						"upstream_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The upstream protocol for listening port.",
						},
						"nginx_server_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Nginx&amp;#39;s server ID.",
						},
					},
				},
			},

			"edition": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.",
			},

			"is_keep_alive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable keep-alive, 0 disable, 1 enable.",
			},

			"instance_i_d": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Unique ID of Instance.",
			},

			"anycast": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Anycast IP switch, 0 represents off and 1 represents on.",
			},

			"weights": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Weight of each upstream.",
			},

			"active_check": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable active health detection, 0 represents disable and 1 represents enable.",
			},

			"t_l_s_version": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Version of TLS Protocol.",
			},

			"ciphers": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Encryption Suite Information.",
			},

			"cipher_template": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Encryption Suite Template, 0:default  1:Universal template 2:Security template 3:Custom template.",
			},

			"proxy_read_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "300s.",
			},

			"proxy_send_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "300s.",
			},

			"sni_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sni type fo upstream, 0:disable SNI；1:enable SNI and SNI equal original request host；2:and SNI equal upstream host 3：enable SNI and equal customize host.",
			},

			"sni_host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "When SniType=3, this parameter needs to be filled in to represent a custom host.",
			},

			"ip_headers": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "When is_cdn=3, this parameter needs to be filled in to indicate a custom header.",
			},

			"x_f_f_reset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0:disable xff reset；1:ensable xff reset.",
			},
		},
	}
}

func resourceTencentCloudWafSaasDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = waf.NewAddSpartaProtectionRequest()
		response = waf.NewAddSpartaProtectionResponse()
		domainId string
	)
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cert_type"); ok {
		request.CertType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_cdn"); ok {
		request.IsCdn = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("upstream_type"); ok {
		request.UpstreamType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_websocket"); ok {
		request.IsWebsocket = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("load_balance"); ok {
		request.LoadBalance = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert"); ok {
		request.Cert = helper.String(v.(string))
	}

	if v, ok := d.GetOk("private_key"); ok {
		request.PrivateKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("s_s_l_id"); ok {
		request.SSLId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upstream_scheme"); ok {
		request.UpstreamScheme = helper.String(v.(string))
	}

	if v, ok := d.GetOk("https_upstream_port"); ok {
		request.HttpsUpstreamPort = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_gray"); ok {
		request.IsGray = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("gray_areas"); ok {
		grayAreasSet := v.(*schema.Set).List()
		for i := range grayAreasSet {
			grayAreas := grayAreasSet[i].(string)
			request.GrayAreas = append(request.GrayAreas, &grayAreas)
		}
	}

	if v, ok := d.GetOk("upstream_domain"); ok {
		request.UpstreamDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_list"); ok {
		srcListSet := v.(*schema.Set).List()
		for i := range srcListSet {
			srcList := srcListSet[i].(string)
			request.SrcList = append(request.SrcList, &srcList)
		}
	}

	if v, ok := d.GetOkExists("is_http2"); ok {
		request.IsHttp2 = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("https_rewrite"); ok {
		request.HttpsRewrite = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ports"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			portItem := waf.PortItem{}
			if v, ok := dMap["port"]; ok {
				portItem.Port = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				portItem.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["upstream_port"]; ok {
				portItem.UpstreamPort = helper.String(v.(string))
			}
			if v, ok := dMap["upstream_protocol"]; ok {
				portItem.UpstreamProtocol = helper.String(v.(string))
			}
			if v, ok := dMap["nginx_server_id"]; ok {
				portItem.NginxServerId = helper.String(v.(string))
			}
			request.Ports = append(request.Ports, &portItem)
		}
	}

	if v, ok := d.GetOk("edition"); ok {
		request.Edition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_keep_alive"); ok {
		request.IsKeepAlive = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("anycast"); ok {
		request.Anycast = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("weights"); ok {
		weightsSet := v.(*schema.Set).List()
		for i := range weightsSet {
			weights := weightsSet[i].(int)
			request.Weights = append(request.Weights, helper.IntInt64(weights))
		}
	}

	if v, ok := d.GetOkExists("active_check"); ok {
		request.ActiveCheck = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("t_l_s_version"); ok {
		request.TLSVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ciphers"); ok {
		ciphersSet := v.(*schema.Set).List()
		for i := range ciphersSet {
			ciphers := ciphersSet[i].(int)
			request.Ciphers = append(request.Ciphers, helper.IntInt64(ciphers))
		}
	}

	if v, ok := d.GetOkExists("cipher_template"); ok {
		request.CipherTemplate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("proxy_read_timeout"); ok {
		request.ProxyReadTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("proxy_send_timeout"); ok {
		request.ProxySendTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("sni_type"); ok {
		request.SniType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sni_host"); ok {
		request.SniHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_headers"); ok {
		ipHeadersSet := v.(*schema.Set).List()
		for i := range ipHeadersSet {
			ipHeaders := ipHeadersSet[i].(string)
			request.IpHeaders = append(request.IpHeaders, &ipHeaders)
		}
	}

	if v, ok := d.GetOkExists("x_f_f_reset"); ok {
		request.XFFReset = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().AddSpartaProtection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create waf saasDomain failed, reason:%+v", logId, err)
		return err
	}

	domainId = *response.Response.DomainId
	d.SetId(domainId)

	return resourceTencentCloudWafSaasDomainRead(d, meta)
}

func resourceTencentCloudWafSaasDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	saasDomainId := d.Id()

	saasDomain, err := service.DescribeWafSaasDomainById(ctx, domainId)
	if err != nil {
		return err
	}

	if saasDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafSaasDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if saasDomain.Domain != nil {
		_ = d.Set("domain", saasDomain.Domain)
	}

	if saasDomain.CertType != nil {
		_ = d.Set("cert_type", saasDomain.CertType)
	}

	if saasDomain.IsCdn != nil {
		_ = d.Set("is_cdn", saasDomain.IsCdn)
	}

	if saasDomain.UpstreamType != nil {
		_ = d.Set("upstream_type", saasDomain.UpstreamType)
	}

	if saasDomain.IsWebsocket != nil {
		_ = d.Set("is_websocket", saasDomain.IsWebsocket)
	}

	if saasDomain.LoadBalance != nil {
		_ = d.Set("load_balance", saasDomain.LoadBalance)
	}

	if saasDomain.Cert != nil {
		_ = d.Set("cert", saasDomain.Cert)
	}

	if saasDomain.PrivateKey != nil {
		_ = d.Set("private_key", saasDomain.PrivateKey)
	}

	if saasDomain.SSLId != nil {
		_ = d.Set("s_s_l_id", saasDomain.SSLId)
	}

	if saasDomain.ResourceId != nil {
		_ = d.Set("resource_id", saasDomain.ResourceId)
	}

	if saasDomain.UpstreamScheme != nil {
		_ = d.Set("upstream_scheme", saasDomain.UpstreamScheme)
	}

	if saasDomain.HttpsUpstreamPort != nil {
		_ = d.Set("https_upstream_port", saasDomain.HttpsUpstreamPort)
	}

	if saasDomain.IsGray != nil {
		_ = d.Set("is_gray", saasDomain.IsGray)
	}

	if saasDomain.GrayAreas != nil {
		_ = d.Set("gray_areas", saasDomain.GrayAreas)
	}

	if saasDomain.UpstreamDomain != nil {
		_ = d.Set("upstream_domain", saasDomain.UpstreamDomain)
	}

	if saasDomain.SrcList != nil {
		_ = d.Set("src_list", saasDomain.SrcList)
	}

	if saasDomain.IsHttp2 != nil {
		_ = d.Set("is_http2", saasDomain.IsHttp2)
	}

	if saasDomain.HttpsRewrite != nil {
		_ = d.Set("https_rewrite", saasDomain.HttpsRewrite)
	}

	if saasDomain.Ports != nil {
		portsList := []interface{}{}
		for _, ports := range saasDomain.Ports {
			portsMap := map[string]interface{}{}

			if saasDomain.Ports.Port != nil {
				portsMap["port"] = saasDomain.Ports.Port
			}

			if saasDomain.Ports.Protocol != nil {
				portsMap["protocol"] = saasDomain.Ports.Protocol
			}

			if saasDomain.Ports.UpstreamPort != nil {
				portsMap["upstream_port"] = saasDomain.Ports.UpstreamPort
			}

			if saasDomain.Ports.UpstreamProtocol != nil {
				portsMap["upstream_protocol"] = saasDomain.Ports.UpstreamProtocol
			}

			if saasDomain.Ports.NginxServerId != nil {
				portsMap["nginx_server_id"] = saasDomain.Ports.NginxServerId
			}

			portsList = append(portsList, portsMap)
		}

		_ = d.Set("ports", portsList)

	}

	if saasDomain.Edition != nil {
		_ = d.Set("edition", saasDomain.Edition)
	}

	if saasDomain.IsKeepAlive != nil {
		_ = d.Set("is_keep_alive", saasDomain.IsKeepAlive)
	}

	if saasDomain.InstanceID != nil {
		_ = d.Set("instance_i_d", saasDomain.InstanceID)
	}

	if saasDomain.Anycast != nil {
		_ = d.Set("anycast", saasDomain.Anycast)
	}

	if saasDomain.Weights != nil {
		_ = d.Set("weights", saasDomain.Weights)
	}

	if saasDomain.ActiveCheck != nil {
		_ = d.Set("active_check", saasDomain.ActiveCheck)
	}

	if saasDomain.TLSVersion != nil {
		_ = d.Set("t_l_s_version", saasDomain.TLSVersion)
	}

	if saasDomain.Ciphers != nil {
		_ = d.Set("ciphers", saasDomain.Ciphers)
	}

	if saasDomain.CipherTemplate != nil {
		_ = d.Set("cipher_template", saasDomain.CipherTemplate)
	}

	if saasDomain.ProxyReadTimeout != nil {
		_ = d.Set("proxy_read_timeout", saasDomain.ProxyReadTimeout)
	}

	if saasDomain.ProxySendTimeout != nil {
		_ = d.Set("proxy_send_timeout", saasDomain.ProxySendTimeout)
	}

	if saasDomain.SniType != nil {
		_ = d.Set("sni_type", saasDomain.SniType)
	}

	if saasDomain.SniHost != nil {
		_ = d.Set("sni_host", saasDomain.SniHost)
	}

	if saasDomain.IpHeaders != nil {
		_ = d.Set("ip_headers", saasDomain.IpHeaders)
	}

	if saasDomain.XFFReset != nil {
		_ = d.Set("x_f_f_reset", saasDomain.XFFReset)
	}

	return nil
}

func resourceTencentCloudWafSaasDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifySpartaProtectionRequest  = waf.NewModifySpartaProtectionRequest()
		modifySpartaProtectionResponse = waf.NewModifySpartaProtectionResponse()
	)

	saasDomainId := d.Id()

	request.DomainId = &domainId

	immutableArgs := []string{"domain", "cert_type", "is_cdn", "upstream_type", "is_websocket", "load_balance", "cert", "private_key", "s_s_l_id", "resource_id", "upstream_scheme", "https_upstream_port", "is_gray", "gray_areas", "upstream_domain", "src_list", "is_http2", "https_rewrite", "ports", "edition", "is_keep_alive", "instance_i_d", "anycast", "weights", "active_check", "t_l_s_version", "ciphers", "cipher_template", "proxy_read_timeout", "proxy_send_timeout", "sni_type", "sni_host", "ip_headers", "x_f_f_reset"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("domain") {
		if v, ok := d.GetOk("domain"); ok {
			request.Domain = helper.String(v.(string))
		}
	}

	if d.HasChange("cert_type") {
		if v, ok := d.GetOkExists("cert_type"); ok {
			request.CertType = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("is_cdn") {
		if v, ok := d.GetOkExists("is_cdn"); ok {
			request.IsCdn = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("upstream_type") {
		if v, ok := d.GetOkExists("upstream_type"); ok {
			request.UpstreamType = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("is_websocket") {
		if v, ok := d.GetOkExists("is_websocket"); ok {
			request.IsWebsocket = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("load_balance") {
		if v, ok := d.GetOk("load_balance"); ok {
			request.LoadBalance = helper.String(v.(string))
		}
	}

	if d.HasChange("cert") {
		if v, ok := d.GetOk("cert"); ok {
			request.Cert = helper.String(v.(string))
		}
	}

	if d.HasChange("private_key") {
		if v, ok := d.GetOk("private_key"); ok {
			request.PrivateKey = helper.String(v.(string))
		}
	}

	if d.HasChange("s_s_l_id") {
		if v, ok := d.GetOk("s_s_l_id"); ok {
			request.SSLId = helper.String(v.(string))
		}
	}

	if d.HasChange("upstream_scheme") {
		if v, ok := d.GetOk("upstream_scheme"); ok {
			request.UpstreamScheme = helper.String(v.(string))
		}
	}

	if d.HasChange("https_upstream_port") {
		if v, ok := d.GetOk("https_upstream_port"); ok {
			request.HttpsUpstreamPort = helper.String(v.(string))
		}
	}

	if d.HasChange("is_gray") {
		if v, ok := d.GetOkExists("is_gray"); ok {
			request.IsGray = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("upstream_domain") {
		if v, ok := d.GetOk("upstream_domain"); ok {
			request.UpstreamDomain = helper.String(v.(string))
		}
	}

	if d.HasChange("src_list") {
		if v, ok := d.GetOk("src_list"); ok {
			srcListSet := v.(*schema.Set).List()
			for i := range srcListSet {
				srcList := srcListSet[i].(string)
				request.SrcList = append(request.SrcList, &srcList)
			}
		}
	}

	if d.HasChange("is_http2") {
		if v, ok := d.GetOkExists("is_http2"); ok {
			request.IsHttp2 = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("https_rewrite") {
		if v, ok := d.GetOkExists("https_rewrite"); ok {
			request.HttpsRewrite = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("ports") {
		if v, ok := d.GetOk("ports"); ok {
			for _, item := range v.([]interface{}) {
				portItem := waf.PortItem{}
				if v, ok := dMap["port"]; ok {
					portItem.Port = helper.String(v.(string))
				}
				if v, ok := dMap["protocol"]; ok {
					portItem.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["upstream_port"]; ok {
					portItem.UpstreamPort = helper.String(v.(string))
				}
				if v, ok := dMap["upstream_protocol"]; ok {
					portItem.UpstreamProtocol = helper.String(v.(string))
				}
				if v, ok := dMap["nginx_server_id"]; ok {
					portItem.NginxServerId = helper.String(v.(string))
				}
				request.Ports = append(request.Ports, &portItem)
			}
		}
	}

	if d.HasChange("edition") {
		if v, ok := d.GetOk("edition"); ok {
			request.Edition = helper.String(v.(string))
		}
	}

	if d.HasChange("is_keep_alive") {
		if v, ok := d.GetOk("is_keep_alive"); ok {
			request.IsKeepAlive = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_i_d") {
		if v, ok := d.GetOk("instance_i_d"); ok {
			request.InstanceID = helper.String(v.(string))
		}
	}

	if d.HasChange("anycast") {
		if v, ok := d.GetOkExists("anycast"); ok {
			request.Anycast = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("weights") {
		if v, ok := d.GetOk("weights"); ok {
			weightsSet := v.(*schema.Set).List()
			for i := range weightsSet {
				weights := weightsSet[i].(int)
				request.Weights = append(request.Weights, helper.IntInt64(weights))
			}
		}
	}

	if d.HasChange("active_check") {
		if v, ok := d.GetOkExists("active_check"); ok {
			request.ActiveCheck = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("t_l_s_version") {
		if v, ok := d.GetOkExists("t_l_s_version"); ok {
			request.TLSVersion = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("ciphers") {
		if v, ok := d.GetOk("ciphers"); ok {
			ciphersSet := v.(*schema.Set).List()
			for i := range ciphersSet {
				ciphers := ciphersSet[i].(int)
				request.Ciphers = append(request.Ciphers, helper.IntInt64(ciphers))
			}
		}
	}

	if d.HasChange("cipher_template") {
		if v, ok := d.GetOkExists("cipher_template"); ok {
			request.CipherTemplate = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("proxy_read_timeout") {
		if v, ok := d.GetOkExists("proxy_read_timeout"); ok {
			request.ProxyReadTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("proxy_send_timeout") {
		if v, ok := d.GetOkExists("proxy_send_timeout"); ok {
			request.ProxySendTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("sni_type") {
		if v, ok := d.GetOkExists("sni_type"); ok {
			request.SniType = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("sni_host") {
		if v, ok := d.GetOk("sni_host"); ok {
			request.SniHost = helper.String(v.(string))
		}
	}

	if d.HasChange("ip_headers") {
		if v, ok := d.GetOk("ip_headers"); ok {
			ipHeadersSet := v.(*schema.Set).List()
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				request.IpHeaders = append(request.IpHeaders, &ipHeaders)
			}
		}
	}

	if d.HasChange("x_f_f_reset") {
		if v, ok := d.GetOkExists("x_f_f_reset"); ok {
			request.XFFReset = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifySpartaProtection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update waf saasDomain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafSaasDomainRead(d, meta)
}

func resourceTencentCloudWafSaasDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	saasDomainId := d.Id()

	if err := service.DeleteWafSaasDomainById(ctx, domainId); err != nil {
		return err
	}

	return nil
}
