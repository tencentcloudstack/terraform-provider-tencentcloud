package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafSaasDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafSaasDomainCreate,
		Read:   resourceTencentCloudWafSaasDomainRead,
		Update: resourceTencentCloudWafSaasDomainUpdate,
		Delete: resourceTencentCloudWafSaasDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique ID of Instance.",
			},
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain names that require defense.",
			},
			"cert_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      CERT_TYPE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CERT_TYPES),
				Description:  "Certificate type, 0 represents no certificate, CertType=1 represents self owned certificate, and 2 represents managed certificate.",
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
			"ssl_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID, When CertType=2, this parameter needs to be filled.",
			},
			"is_cdn": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      ISCDN_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ISCDN_STSTUS),
				Description:  "Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.",
			},
			"upstream_scheme": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(UPSTREAM_SCHEMES),
				Description:  "Upstream scheme for https, http or https.",
			},
			"https_upstream_port": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upstream port for https, When listen ports has https port and UpstreamScheme is HTTP, the current field needs to be filled.",
			},
			"upstream_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      UP_STREAM_TYPE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(UP_STREAM_TYPES),
				Description:  "Upstream type, 0 represents IP, 1 represents domain name.",
			},
			"upstream_domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upstream domain, When UpstreamType=1, this parameter needs to be filled.",
			},
			"src_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Upstream IP List, When UpstreamType=0, this parameter needs to be filled.",
			},
			"is_http2": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      IS_HTTP2_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(IS_HTTP2_STATUS),
				Description:  "Whether enable HTTP2, Enabling HTTP2 requires HTTPS support, 1 means enabled, 0 does not.",
			},
			"is_websocket": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      IS_WEBSOCKET_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(IS_WEBSOCKET_STATUS),
				Description:  "Is WebSocket support enabled. 1 means enabled, 0 does not.",
			},
			"load_balance": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      LOAD_BALANCE_0,
				ValidateFunc: tccommon.ValidateAllowedStringValue(LOAD_BALANCE_STATUS),
				Description:  "Load balancing strategy, where 0 represents polling and 1 represents IP hash and 2 weighted round robin.",
			},
			"https_rewrite": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      HTTPS_REWRITE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(HTTPS_REWRITE_STATUS),
				Description:  "Whether redirect to https, 1 will redirect and 0 will not.",
			},
			"ports": {
				Required:    true,
				Type:        schema.TypeSet,
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
							Computed:    true,
							Description: "Nginx server ID.",
						},
					},
				},
			},
			"is_keep_alive": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      IS_KEEP_ALIVE_1,
				ValidateFunc: tccommon.ValidateAllowedStringValue(IS_KEEP_ALIVE_STATUS),
				Description:  "Whether to enable keep-alive, 0 disable, 1 enable.",
			},
			"weights": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Weight of each upstream.",
			},
			"active_check": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      ACTIVE_CHECK_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ACTIVE_CHECK_STATUS),
				Description:  "Whether to enable active health detection, 0 represents disable and 1 represents enable.",
			},
			"tls_version": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     TLS_VERSION_STATUS_3,
				Description: "Version of TLS Protocol.",
			},
			"ciphers": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Encryption Suite Information.",
			},
			"cipher_template": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      CIPHER_TEMPLATE_1,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CIPHER_TEMPLATES),
				Description:  "Encryption Suite Template, 0:default  1:Universal template 2:Security template 3:Custom template.",
			},
			"proxy_read_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     PROXY_READ_TIMEOUT,
				Description: "300s.",
			},
			"proxy_send_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     PROXY_SEND_TIMEOUT,
				Description: "300s.",
			},
			"sni_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      SNI_TYPE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(SNI_TYPES),
				Description:  "Sni type fo upstream, 0:disable SNI; 1:enable SNI and SNI equal original request host; 2:and SNI equal upstream host 3:enable SNI and equal customize host.",
			},
			"sni_host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "When SniType=3, this parameter needs to be filled in to represent a custom host.",
			},
			"ip_headers": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "When is_cdn=3, this parameter needs to be filled in to indicate a custom header.",
			},
			"xff_reset": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      XFF_RESET_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(XFF_RESET_STATUS),
				Description:  "0:disable xff reset; 1:ensable xff reset.",
			},
			"bot_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      BOT_STATUS_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(BOT_STATUS),
				Description:  "Whether to enable bot, 1 enable, 0 disable.",
			},
			"api_safe_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      API_SAFE_STATUS_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(API_SAFE_STATUS),
				Description:  "Whether to enable api safe, 1 enable, 0 disable.",
			},
			"cls_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      CLS_STATUS_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CLS_STATUS),
				Description:  "Whether to enable access logs, 1 enable, 0 disable.",
			},
			//"ipv6_status": {
			//	Type:         schema.TypeInt,
			//	Optional:     true,
			//	Default:      IPV6_STATUS_0,
			//	ValidateFunc: tccommon.ValidateAllowedIntValue(IPV6_STATUS),
			//	Description:  "Whether to enable ipv6, 1 enable, 0 disable.",
			//},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      PROTECTION_STATUS_1,
				ValidateFunc: tccommon.ValidateAllowedIntValue(PROTECTION_STATUS),
				Description:  "Binding status between waf and LB, 0:not bind, 1:binding.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain id.",
			},
		},
	}
}

func resourceTencentCloudWafSaasDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service          = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		verifyRequest    = waf.NewDescribeDomainVerifyResultRequest()
		request          = waf.NewAddSpartaProtectionRequest()
		instanceID       string
		domain           string
		domainId         string
		loadBalance      string
		botStatus        uint64
		apiSafeStatus    uint64
		clsStatus        uint64
		protectionStatus uint64
		isCdn            int
		//ipv6Status    int64

	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = helper.String(v.(string))
		instanceID = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	// check domain legal
	verifyRequest.InstanceID = &instanceID
	verifyRequest.Domain = &domain
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().DescribeDomainVerifyResult(verifyRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.VerifyCode == DescribeDomainVerifyResultSUCCESS {
			return nil
		}

		e = fmt.Errorf("The current domain %s is illegal, errMsg: %s.", domain, *result.Response.Msg)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s verify waf saasDomain failed, reason:%+v", logId, err)
		return err
	}

	if v, ok := d.GetOkExists("cert_type"); ok {
		request.CertType = helper.IntInt64(v.(int))

		cert := d.Get("cert").(string)
		privateKey := d.Get("private_key").(string)
		sslId := d.Get("ssl_id").(string)

		certType := v.(int)
		if certType == CERT_TYPE_0 {
			if cert != "" || privateKey != "" || sslId != "" {
				return fmt.Errorf("If `cert_type` is 0, not support setting `cert`, `private_key`, `ssl_id`.")
			}

		} else if certType == CERT_TYPE_1 {
			if sslId != "" {
				return fmt.Errorf("If `cert_type` is 1, not support setting `ssl_id`.")
			}

			if cert == "" || privateKey == "" {
				return fmt.Errorf("If `cert_type` is 1, `cert`, `private_key` is required.")
			}

			request.Cert = &cert
			request.PrivateKey = &privateKey

		} else {
			if cert != "" || privateKey != "" {
				return fmt.Errorf("If `cert_type` is 2, not support setting `cert`, `private_key`.")
			}

			if sslId == "" {
				return fmt.Errorf("If `cert_type` is 2, `ssl_id` is required.")
			}

			request.SSLId = &sslId
		}
	}

	if v, ok := d.GetOkExists("is_cdn"); ok {
		request.IsCdn = helper.IntInt64(v.(int))
		isCdn = v.(int)
	}

	if v, ok := d.GetOk("load_balance"); ok {
		request.LoadBalance = helper.String(v.(string))
		loadBalance = v.(string)
	}

	if v, ok := d.GetOk("upstream_scheme"); ok {
		request.UpstreamScheme = helper.String(v.(string))

		httpsUpstreamPort := d.Get("https_upstream_port").(string)

		upstreamScheme := v.(string)
		if upstreamScheme == UPSTREAM_SCHEME_HTTP {
			if httpsUpstreamPort == "" {
				return fmt.Errorf("If `upstream_scheme` is http, `https_upstream_port` is required.")
			}

			request.HttpsUpstreamPort = &httpsUpstreamPort
		}
	}

	if v, ok := d.GetOkExists("upstream_type"); ok {
		request.UpstreamType = helper.IntInt64(v.(int))

		upstreamType := v.(int)
		if upstreamType == UP_STREAM_TYPE_0 {
			if _, ok := d.GetOk("upstream_domain"); ok {
				return fmt.Errorf("If `upstream_type` is 0, not support setting `upstream_domain`.")
			}

			if v, ok := d.GetOk("src_list"); ok {
				srcListSet := v.([]interface{})
				for i := range srcListSet {
					srcList := srcListSet[i].(string)
					request.SrcList = append(request.SrcList, &srcList)
				}

				if len(srcListSet) == 1 {
					if _, ok := d.GetOk("weights"); ok {
						return fmt.Errorf("If `src_list` length is 1, not support setting `weights`.")
					}
				} else {
					if loadBalance != LOAD_BALANCE_2 {
						return fmt.Errorf("If `load_balance` is 0 or 1, not support setting `weights`.")
					}

					if v, ok := d.GetOk("weights"); ok {
						weightsSet := v.([]interface{})
						if len(weightsSet) != len(srcListSet) {
							return fmt.Errorf("The lengths of `weights` and `src_list` are not equal.")
						}

						for i := range weightsSet {
							weight := int64(weightsSet[i].(int))
							request.Weights = append(request.Weights, &weight)
						}
					}
				}

			} else {
				return fmt.Errorf("If `upstream_type` is 0, `src_list` is required.")
			}

			if v, ok := d.GetOk("is_keep_alive"); ok {
				request.IsKeepAlive = helper.String(v.(string))
			}

		} else {
			if _, ok := d.GetOk("src_list"); ok {
				return fmt.Errorf("If `upstream_type` is 1, not support setting `src_list`.")
			}

			if _, ok := d.GetOk("weights"); ok {
				return fmt.Errorf("If `upstream_type` is 1, not support setting `weights`.")
			}

			if v, ok := d.GetOk("is_keep_alive"); ok {
				request.IsKeepAlive = helper.String(v.(string))
			}

			if v, ok := d.GetOk("upstream_domain"); ok {
				request.UpstreamDomain = helper.String(v.(string))
			} else {
				return fmt.Errorf("If `upstream_type` is 1, `upstream_domain` is required.")
			}
		}
	}

	if v, ok := d.GetOkExists("is_http2"); ok {
		request.IsHttp2 = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_websocket"); ok {
		request.IsWebsocket = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("https_rewrite"); ok {
		request.HttpsRewrite = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ports"); ok {
		for _, item := range v.(*schema.Set).List() {
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

			portItem.NginxServerId = helper.String("0")

			request.Ports = append(request.Ports, &portItem)
		}
	}

	if v, ok := d.GetOkExists("active_check"); ok {
		request.ActiveCheck = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("tls_version"); ok {
		request.TLSVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("cipher_template"); ok {
		request.CipherTemplate = helper.IntInt64(v.(int))

		cipherTemplate := v.(int)

		if cipherTemplate != CIPHER_TEMPLATE_3 {
			if _, ok := d.GetOk("ciphers"); ok {
				return fmt.Errorf("If `cipher_template` is 1 or 2, not support setting `ciphers`.")
			}
		} else {
			if v, ok := d.GetOk("ciphers"); ok {
				ciphersSet := v.([]interface{})
				for i := range ciphersSet {
					ciphers := ciphersSet[i].(int)
					request.Ciphers = append(request.Ciphers, helper.IntInt64(ciphers))
				}
			}
		}
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
		if isCdn == ISCDN_3 {
			ipHeadersSet := v.([]interface{})
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				request.IpHeaders = append(request.IpHeaders, &ipHeaders)
			}
		} else {
			return fmt.Errorf("If `is_cdn` is %d, not supported setting `ip_headers`.", isCdn)
		}
	}

	if v, ok := d.GetOkExists("xff_reset"); ok {
		request.XFFReset = helper.IntInt64(v.(int))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().AddSpartaProtection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf saasDomain failed, reason:%+v", logId, err)
		return err
	}

	// get domain id
	domainInfo, err := service.DescribeDomainsById(ctx, instanceID, domain)
	if err != nil {
		return err
	}

	if domainInfo == nil {
		log.Printf("[WARN]%s resource `DescribeDomains` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainInfo.DomainId != nil {
		domainId = *domainInfo.DomainId
	}

	d.SetId(strings.Join([]string{instanceID, domain, domainId}, tccommon.FILED_SP))

	// wait domain state
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDomainsById(ctx, instanceID, domain)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if *result.State == 0 || *result.State == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("domain is still in state %d", *result.State))
	})

	if err != nil {
		return err
	}

	// set bot
	if v, ok := d.GetOkExists("bot_status"); ok {
		tmpBotStatus := v.(int)

		if tmpBotStatus != BOT_STATUS_0 {
			botStatus = uint64(tmpBotStatus)
			modifyBotStatusRequest := waf.NewModifyBotStatusRequest()
			modifyBotStatusRequest.Domain = &domain
			modifyBotStatusRequest.InstanceID = &instanceID
			tmpStatus := strconv.FormatUint(botStatus, 10)
			modifyBotStatusRequest.Status = &tmpStatus
			modifyBotStatusRequest.Category = common.StringPtr("bot")

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyBotStatus(modifyBotStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyBotStatusRequest.GetAction(), modifyBotStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf saasDomain bot_status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set api safe
	if v, ok := d.GetOkExists("api_safe_status"); ok {
		tmpApiSafeStatus := v.(int)

		if tmpApiSafeStatus != API_SAFE_STATUS_0 {
			apiSafeStatus = uint64(tmpApiSafeStatus)
			modifyApiAnalyzeStatusRequest := waf.NewModifyApiAnalyzeStatusRequest()
			modifyApiAnalyzeStatusRequest.Domain = &domain
			modifyApiAnalyzeStatusRequest.InstanceId = &instanceID
			modifyApiAnalyzeStatusRequest.Status = &apiSafeStatus

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyApiAnalyzeStatus(modifyApiAnalyzeStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyApiAnalyzeStatusRequest.GetAction(), modifyApiAnalyzeStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf saasDomain api_safe_status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set cls
	if v, ok := d.GetOkExists("cls_status"); ok {
		tmpClsStatus := v.(int)

		if tmpClsStatus != CLS_STATUS_0 {
			clsStatus = uint64(tmpClsStatus)
			modifyDomainsCLSStatusRequest := waf.NewModifyDomainsCLSStatusRequest()
			modifyDomainsCLSStatusRequest.Domains = []*waf.DomainURI{
				{
					Domain:     common.StringPtr(domain),
					Edition:    common.StringPtr("sparta-waf"),
					InstanceID: common.StringPtr(instanceID),
				},
			}
			modifyDomainsCLSStatusRequest.Status = &clsStatus

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyDomainsCLSStatus(modifyDomainsCLSStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyDomainsCLSStatusRequest.GetAction(), modifyDomainsCLSStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain cls_status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set ipv6
	//if v, ok := d.GetOkExists("ipv6_status"); ok {
	//	tmpIpv6Status := v.(int)
	//
	//	if tmpIpv6Status != IPV6_STATUS_0 {
	//		ipv6Status = int64(IPV6_ON)
	//		modifyDomainIpv6StatusRequest := waf.NewModifyDomainIpv6StatusRequest()
	//		modifyDomainIpv6StatusRequest.Domain = &domain
	//		modifyDomainIpv6StatusRequest.DomainId = &domainId
	//		modifyDomainIpv6StatusRequest.InstanceId = &instanceID
	//		modifyDomainIpv6StatusRequest.Status = &ipv6Status
	//
	//		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
	//			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyDomainIpv6Status(modifyDomainIpv6StatusRequest)
	//			if e != nil {
	//				return tccommon.RetryError(e)
	//			} else {
	//				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyDomainIpv6StatusRequest.GetAction(), modifyDomainIpv6StatusRequest.ToJsonString(), result.ToJsonString())
	//			}
	//
	//			return nil
	//		})
	//
	//		if err != nil {
	//			log.Printf("[CRITAL]%s modify waf saasDomain ipv6_status failed, reason:%+v", logId, err)
	//			return err
	//		}
	//	}
	//}

	// set status
	if v, ok := d.GetOkExists("status"); ok {
		tmpProtectionStatus := v.(int)

		if tmpProtectionStatus != PROTECTION_STATUS_1 {
			protectionStatus = uint64(tmpProtectionStatus)
			modifyProtectionStatusRequest := waf.NewModifyProtectionStatusRequest()
			modifyProtectionStatusRequest.Domain = &domain
			modifyProtectionStatusRequest.Status = &protectionStatus

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyProtectionStatus(modifyProtectionStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyProtectionStatusRequest.GetAction(), modifyProtectionStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf saasDomain status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafSaasDomainRead(d, meta)
}

func resourceTencentCloudWafSaasDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]
	domainId := idSplit[2]

	saasDomain, err := service.DescribeWafSaasDomainById(ctx, instanceID, domain, domainId)
	if err != nil {
		return err
	}

	if saasDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafSaasDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	domainInfo, err := service.DescribeDomainsById(ctx, instanceID, domain)
	if err != nil {
		return err
	}

	if domainInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DescribeDomains` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceID)

	_ = d.Set("domain", domain)

	if saasDomain.CertType != nil {
		_ = d.Set("cert_type", saasDomain.CertType)
	}

	if saasDomain.Cert != nil {
		_ = d.Set("cert", saasDomain.Cert)
	}

	if saasDomain.PrivateKey != nil {
		_ = d.Set("private_key", saasDomain.PrivateKey)
	}

	if saasDomain.SSLId != nil {
		_ = d.Set("ssl_id", saasDomain.SSLId)
	}

	if saasDomain.IsCdn != nil {
		_ = d.Set("is_cdn", saasDomain.IsCdn)
	}

	if saasDomain.UpstreamScheme != nil {
		_ = d.Set("upstream_scheme", saasDomain.UpstreamScheme)
	}

	if saasDomain.HttpsUpstreamPort != nil {
		_ = d.Set("https_upstream_port", saasDomain.HttpsUpstreamPort)
	}

	if saasDomain.UpstreamType != nil {
		_ = d.Set("upstream_type", saasDomain.UpstreamType)
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

	if saasDomain.IsWebsocket != nil {
		_ = d.Set("is_websocket", saasDomain.IsWebsocket)
	}

	if saasDomain.LoadBalance != nil {
		tmpLoadBalance := *saasDomain.LoadBalance
		loadBalance := strconv.FormatUint(tmpLoadBalance, 10)
		_ = d.Set("load_balance", loadBalance)
	}

	if saasDomain.HttpsRewrite != nil {
		_ = d.Set("https_rewrite", saasDomain.HttpsRewrite)
	}

	if saasDomain.Ports != nil {
		portsList := []interface{}{}
		for _, ports := range saasDomain.Ports {
			portsMap := map[string]interface{}{}

			if ports.Port != nil {
				portsMap["port"] = ports.Port
			}

			if ports.Protocol != nil {
				portsMap["protocol"] = ports.Protocol
			}

			if ports.UpstreamPort != nil {
				portsMap["upstream_port"] = ports.UpstreamPort
			}

			if ports.UpstreamProtocol != nil {
				portsMap["upstream_protocol"] = ports.UpstreamProtocol
			}

			if ports.NginxServerId != nil {
				tmpNginxServerId := *ports.NginxServerId
				nginxServerId := strconv.FormatUint(tmpNginxServerId, 10)
				portsMap["nginx_server_id"] = nginxServerId
			}

			portsList = append(portsList, portsMap)
		}

		_ = d.Set("ports", portsList)

	}

	if saasDomain.IsKeepAlive != nil {
		tmpIsKeepAlive := *saasDomain.IsKeepAlive
		isKeepAlive := strconv.FormatUint(tmpIsKeepAlive, 10)
		_ = d.Set("is_keep_alive", isKeepAlive)
	}

	if saasDomain.Weights != nil {
		tmpList := make([]int, 0, len(saasDomain.Weights))
		for _, v := range saasDomain.Weights {
			item, _ := strconv.Atoi(*v)
			tmpList = append(tmpList, item)
		}

		_ = d.Set("weights", tmpList)
	}

	if saasDomain.ActiveCheck != nil {
		_ = d.Set("active_check", saasDomain.ActiveCheck)
	}

	if saasDomain.TLSVersion != nil {
		_ = d.Set("tls_version", saasDomain.TLSVersion)
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
		_ = d.Set("xff_reset", saasDomain.XFFReset)
	}

	if domainInfo.BotStatus != nil {
		_ = d.Set("bot_status", domainInfo.BotStatus)
	}

	if domainInfo.ApiStatus != nil {
		_ = d.Set("api_safe_status", domainInfo.ApiStatus)
	}

	if domainInfo.ClsStatus != nil {
		_ = d.Set("cls_status", domainInfo.ClsStatus)
	}

	if domainInfo.Status != nil {
		_ = d.Set("status", domainInfo.Status)
	}

	//if domainInfo.Ipv6Status != nil {
	//	_ = d.Set("ipv6_status", domainInfo.Ipv6Status)
	//}

	return nil
}

func resourceTencentCloudWafSaasDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = waf.NewModifySpartaProtectionRequest()
		botStatus     uint64
		apiSafeStatus uint64
		clsStatus     uint64
		//ipv6Status    int64
		loadBalance string
		isCdn       int
	)

	immutableArgs := []string{"instance_id", "domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]
	domainId := idSplit[2]

	// set waf status
	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			tmpProtectionStatus := v.(int)
			// open first
			if tmpProtectionStatus == PROTECTION_STATUS_1 {
				protectionStatus := uint64(tmpProtectionStatus)
				modifyProtectionStatusRequest := waf.NewModifyProtectionStatusRequest()
				modifyProtectionStatusRequest.Domain = &domain
				modifyProtectionStatusRequest.Status = &protectionStatus

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyProtectionStatus(modifyProtectionStatusRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyProtectionStatusRequest.GetAction(), modifyProtectionStatusRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify waf saasDomain status failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	// get ports by api
	saasDomain, err := service.DescribeWafSaasDomainById(ctx, instanceID, domain, domainId)
	if err != nil {
		return err
	}

	if saasDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafSaasDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	portsList := []interface{}{}
	if saasDomain.Ports != nil {
		for _, ports := range saasDomain.Ports {
			portsMap := map[string]interface{}{}

			if ports.Port != nil {
				portsMap["port"] = *ports.Port
			}

			if ports.Protocol != nil {
				portsMap["protocol"] = *ports.Protocol
			}

			if ports.UpstreamPort != nil {
				portsMap["upstream_port"] = *ports.UpstreamPort
			}

			if ports.UpstreamProtocol != nil {
				portsMap["upstream_protocol"] = *ports.UpstreamProtocol
			}

			if ports.NginxServerId != nil {
				tmpNginxServerId := *ports.NginxServerId
				nginxServerId := strconv.FormatUint(tmpNginxServerId, 10)
				portsMap["nginx_server_id"] = nginxServerId
			}

			portsList = append(portsList, portsMap)
		}
	}

	request.Domain = &domain
	request.DomainId = &domainId
	request.InstanceID = &instanceID

	if v, ok := d.GetOkExists("cert_type"); ok {
		request.CertType = helper.IntInt64(v.(int))

		cert := d.Get("cert").(string)
		privateKey := d.Get("private_key").(string)
		sslId := d.Get("ssl_id").(string)

		certType := v.(int)
		if certType == CERT_TYPE_0 {
			if cert != "" || privateKey != "" || sslId != "" {
				return fmt.Errorf("If `cert_type` is 0, not support setting `cert`, `private_key`, `ssl_id`.")
			}

		} else if certType == CERT_TYPE_1 {
			if sslId != "" {
				return fmt.Errorf("If `cert_type` is 1, not support setting `ssl_id`.")
			}

			if cert == "" || privateKey == "" {
				return fmt.Errorf("If `cert_type` is 1, `cert`, `private_key` is required.")
			}

			request.Cert = &cert
			request.PrivateKey = &privateKey

		} else {
			if cert != "" || privateKey != "" {
				return fmt.Errorf("If `cert_type` is 2, not support setting `cert`, `private_key`.")
			}

			if sslId == "" {
				return fmt.Errorf("If `cert_type` is 2, `ssl_id` is required.")
			}

			request.SSLId = &sslId
		}
	}

	if v, ok := d.GetOkExists("is_cdn"); ok {
		request.IsCdn = helper.IntInt64(v.(int))
		isCdn = v.(int)
	}

	if v, ok := d.GetOk("load_balance"); ok {
		loadBalance = v.(string)
		loadBalanceInt64, _ := strconv.ParseInt(loadBalance, 10, 64)
		request.LoadBalance = &loadBalanceInt64

	}

	if v, ok := d.GetOk("upstream_scheme"); ok {
		request.UpstreamScheme = helper.String(v.(string))

		httpsUpstreamPort := d.Get("https_upstream_port").(string)

		upstreamScheme := v.(string)
		if upstreamScheme == UPSTREAM_SCHEME_HTTP {
			if httpsUpstreamPort == "" {
				return fmt.Errorf("If `upstream_scheme` is http, `https_upstream_port` is required.")
			}

			request.HttpsUpstreamPort = &httpsUpstreamPort
		}
	}

	if v, ok := d.GetOkExists("upstream_type"); ok {
		request.UpstreamType = helper.IntInt64(v.(int))

		upstreamType := v.(int)
		if upstreamType == UP_STREAM_TYPE_0 {
			if _, ok := d.GetOk("upstream_domain"); ok {
				return fmt.Errorf("If `upstream_type` is 0, not support setting `upstream_domain`.")
			}

			if v, ok := d.GetOk("src_list"); ok {
				srcListSet := v.([]interface{})
				for i := range srcListSet {
					srcList := srcListSet[i].(string)
					request.SrcList = append(request.SrcList, &srcList)
				}

				if len(srcListSet) == 1 {
					if _, ok := d.GetOk("weights"); ok {
						return fmt.Errorf("If `src_list` length is 1, not support setting `weights`.")
					}
				} else {
					if loadBalance != LOAD_BALANCE_2 {
						return fmt.Errorf("If `load_balance` is 0 or 1, not support setting `weights`.")
					}

					if v, ok := d.GetOk("weights"); ok {
						weightsSet := v.([]interface{})
						if len(weightsSet) != len(srcListSet) {
							return fmt.Errorf("The lengths of `weights` and `src_list` are not equal.")
						}

						for i := range weightsSet {
							weight := int64(weightsSet[i].(int))
							request.Weights = append(request.Weights, &weight)
						}
					}
				}

			} else {
				return fmt.Errorf("If `upstream_type` is 0, `src_list` is required.")
			}

			if v, ok := d.GetOk("is_keep_alive"); ok {
				request.IsKeepAlive = helper.String(v.(string))
			}

		} else {
			if _, ok := d.GetOk("src_list"); ok {
				return fmt.Errorf("If `upstream_type` is 1, not support setting `src_list`.")
			}

			if _, ok := d.GetOk("weights"); ok {
				return fmt.Errorf("If `upstream_type` is 1, not support setting `weights`.")
			}

			if v, ok := d.GetOk("is_keep_alive"); ok {
				request.IsKeepAlive = helper.String(v.(string))
			}

			if v, ok := d.GetOk("upstream_domain"); ok {
				request.UpstreamDomain = helper.String(v.(string))
			} else {
				return fmt.Errorf("If `upstream_type` is 1, `upstream_domain` is required.")
			}
		}
	}

	if v, ok := d.GetOkExists("is_http2"); ok {
		request.IsHttp2 = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_websocket"); ok {
		request.IsWebsocket = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("https_rewrite"); ok {
		request.HttpsRewrite = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("ports"); ok {
		// get ports by user
		tmpPortsList := []interface{}{}
		for _, item := range v.(*schema.Set).List() {
			portsMap := map[string]interface{}{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["port"]; ok {
				portsMap["port"] = v.(string)
			}

			if v, ok := dMap["protocol"]; ok {
				portsMap["protocol"] = v.(string)
			}

			if v, ok := dMap["upstream_port"]; ok {
				portsMap["upstream_port"] = v.(string)
			}

			if v, ok := dMap["upstream_protocol"]; ok {
				portsMap["upstream_protocol"] = v.(string)
			}

			tmpPortsList = append(tmpPortsList, portsMap)
		}

		// check ports
		resPort := checkPorts(portsList, tmpPortsList)

		for _, item := range resPort {
			dMap := item.(map[string]interface{})
			portItem := waf.SpartaProtectionPort{}
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
				tmpNginxServerId := v.(string)
				intNginxServerId, _ := strconv.Atoi(tmpNginxServerId)
				portItem.NginxServerId = helper.IntUint64(intNginxServerId)
			}

			request.Ports = append(request.Ports, &portItem)
		}
	}

	if v, ok := d.GetOkExists("active_check"); ok {
		request.ActiveCheck = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("tls_version"); ok {
		request.TLSVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("cipher_template"); ok {
		request.CipherTemplate = helper.IntInt64(v.(int))

		if v, ok := d.GetOk("ciphers"); ok {
			ciphersSet := v.([]interface{})
			for i := range ciphersSet {
				ciphers := ciphersSet[i].(int)
				request.Ciphers = append(request.Ciphers, helper.IntInt64(ciphers))
			}
		}
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
		if isCdn == ISCDN_3 {
			ipHeadersSet := v.([]interface{})
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				request.IpHeaders = append(request.IpHeaders, &ipHeaders)
			}
		} else {
			return fmt.Errorf("If `is_cdn` is %d, not supported setting `ip_headers`.", isCdn)
		}
	}

	if v, ok := d.GetOkExists("xff_reset"); ok {
		request.XFFReset = helper.IntInt64(v.(int))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifySpartaProtection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf saasDomain failed, reason:%+v", logId, err)
		return err
	}

	// set bot
	if d.HasChange("bot_status") {
		if v, ok := d.GetOkExists("bot_status"); ok {
			botStatus = uint64(v.(int))
			modifyBotStatusRequest := waf.NewModifyBotStatusRequest()
			modifyBotStatusRequest.Domain = &domain
			modifyBotStatusRequest.InstanceID = &instanceID
			tmpStatus := strconv.FormatUint(botStatus, 10)
			modifyBotStatusRequest.Status = &tmpStatus
			modifyBotStatusRequest.Category = common.StringPtr("bot")

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyBotStatus(modifyBotStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyBotStatusRequest.GetAction(), modifyBotStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf saasDomain bot_status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set api safe
	if d.HasChange("api_safe_status") {
		if v, ok := d.GetOkExists("api_safe_status"); ok {
			apiSafeStatus = uint64(v.(int))
			modifyApiAnalyzeStatusRequest := waf.NewModifyApiAnalyzeStatusRequest()
			modifyApiAnalyzeStatusRequest.Domain = &domain
			modifyApiAnalyzeStatusRequest.InstanceId = &instanceID
			modifyApiAnalyzeStatusRequest.Status = &apiSafeStatus

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyApiAnalyzeStatus(modifyApiAnalyzeStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyApiAnalyzeStatusRequest.GetAction(), modifyApiAnalyzeStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf saasDomain api_safe_status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set cls
	if v, ok := d.GetOkExists("cls_status"); ok {
		clsStatus = uint64(v.(int))
		modifyDomainsCLSStatusRequest := waf.NewModifyDomainsCLSStatusRequest()
		modifyDomainsCLSStatusRequest.Domains = []*waf.DomainURI{
			{
				Domain:     common.StringPtr(domain),
				Edition:    common.StringPtr("sparta-waf"),
				InstanceID: common.StringPtr(instanceID),
			},
		}
		modifyDomainsCLSStatusRequest.Status = &clsStatus

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyDomainsCLSStatus(modifyDomainsCLSStatusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyDomainsCLSStatusRequest.GetAction(), modifyDomainsCLSStatusRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify waf clbDomain cls_status failed, reason:%+v", logId, err)
			return err
		}
	}

	// set ipv6
	//if d.HasChange("ipv6_status") {
	//	if v, ok := d.GetOkExists("ipv6_status"); ok {
	//		tmpIpv6Status := v.(int)
	//		if tmpIpv6Status == IPV6_STATUS_0 {
	//			ipv6Status = int64(IPV6_OFF)
	//		} else {
	//			ipv6Status = int64(IPV6_ON)
	//		}
	//		modifyDomainIpv6StatusRequest := waf.NewModifyDomainIpv6StatusRequest()
	//		modifyDomainIpv6StatusRequest.Domain = &domain
	//		modifyDomainIpv6StatusRequest.DomainId = &domainId
	//		modifyDomainIpv6StatusRequest.InstanceId = &instanceID
	//		modifyDomainIpv6StatusRequest.Status = &ipv6Status
	//
	//		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
	//			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyDomainIpv6Status(modifyDomainIpv6StatusRequest)
	//			if e != nil {
	//				return tccommon.RetryError(e)
	//			} else {
	//				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyDomainIpv6StatusRequest.GetAction(), modifyDomainIpv6StatusRequest.ToJsonString(), result.ToJsonString())
	//			}
	//
	//			return nil
	//		})
	//
	//		if err != nil {
	//			log.Printf("[CRITAL]%s modify waf saasDomain ipv6_status failed, reason:%+v", logId, err)
	//			return err
	//		}
	//	}
	//}

	// set waf status
	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			tmpProtectionStatus := v.(int)
			// close end
			if tmpProtectionStatus == PROTECTION_STATUS_0 {
				protectionStatus := uint64(tmpProtectionStatus)
				modifyProtectionStatusRequest := waf.NewModifyProtectionStatusRequest()
				modifyProtectionStatusRequest.Domain = &domain
				modifyProtectionStatusRequest.Status = &protectionStatus

				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyProtectionStatus(modifyProtectionStatusRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyProtectionStatusRequest.GetAction(), modifyProtectionStatusRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify waf saasDomain status failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudWafSaasDomainRead(d, meta)
}

func resourceTencentCloudWafSaasDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]

	if err := service.DeleteWafSaasDomainById(ctx, instanceID, domain); err != nil {
		return err
	}

	return nil
}

func checkPorts(portsList []interface{}, tmpPortsList []interface{}) (res []interface{}) {
	for _, tmpItem := range tmpPortsList {
		loopComplete := true
		portsMap := map[string]interface{}{}
		tmpMap := tmpItem.(map[string]interface{})
		tmpPort := tmpMap["port"]
		tmpProtocol := tmpMap["protocol"]
		tmpUpstreamPort := tmpMap["upstream_port"]
		tmpUpstreamProtocol := tmpMap["upstream_protocol"]

		for _, item := range portsList {
			dMap := item.(map[string]interface{})
			dPort := dMap["port"]
			dProtocol := dMap["protocol"]
			dUpstreamPort := dMap["upstream_port"]
			dUpstreamProtocol := dMap["upstream_protocol"]
			dNginxServerId := dMap["nginx_server_id"]

			if tmpPort == dPort && tmpProtocol == dProtocol && tmpUpstreamPort == dUpstreamPort && tmpUpstreamProtocol == dUpstreamProtocol {
				portsMap["port"] = tmpPort
				portsMap["protocol"] = tmpProtocol
				portsMap["upstream_port"] = tmpUpstreamPort
				portsMap["upstream_protocol"] = tmpUpstreamProtocol
				portsMap["nginx_server_id"] = dNginxServerId
				res = append(res, portsMap)
				loopComplete = false
				break
			}
		}

		if loopComplete {
			portsMap["port"] = tmpPort
			portsMap["protocol"] = tmpProtocol
			portsMap["upstream_port"] = tmpUpstreamPort
			portsMap["upstream_protocol"] = tmpUpstreamProtocol
			portsMap["nginx_server_id"] = "0"
			res = append(res, portsMap)
		}
	}

	return
}
