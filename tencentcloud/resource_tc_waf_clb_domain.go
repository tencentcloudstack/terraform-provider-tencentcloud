package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafClbDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafClbDomainCreate,
		Read:   resourceTencentCloudWafClbDomainRead,
		Update: resourceTencentCloudWafClbDomainUpdate,
		Delete: resourceTencentCloudWafClbDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance unique ID.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name.",
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      CLB_DOMAIN_STATUS_1,
				ValidateFunc: validateAllowedIntValue(CLB_DOMAIN_STATUS),
				Description:  "Binding status between waf and LB, 0:not bind, 1:binding.",
			},
			"engine": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      CLB_DOMAIN_ENGINE_20,
				ValidateFunc: validateAllowedIntValue(CLB_DOMAIN_ENGINE),
				Description:  "Protection Status: 10: Rule Observation&&AI Off Mode, 11: Rule Observation&&AI Observation Mode, 12: Rule Observation&&AI Interception Mode, 20: Rule Interception&&AI Off Mode, 21: Rule Interception&&AI Observation Mode, 22: Rule Interception&&AI Interception Mode, Default 20.",
			},
			"is_cdn": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      ISCDN_0,
				ValidateFunc: validateAllowedIntValue(ISCDN_STSTUS),
				Description:  "Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.",
			},
			"load_balancer_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of bound LB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LoadBalancer unique ID.",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LoadBalancer name.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of listener in LB.",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listener name.",
						},
						"vip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LoadBalancer IP.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "LoadBalancer port.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LoadBalancer region.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol of listener, http or https.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LoadBalancer zone.",
						},
						"numerical_vpc_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VPCID for load balancer, public network is -1, and internal network is filled in according to actual conditions.",
						},
						"load_balancer_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Network type for load balancer.",
						},
					},
				},
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Regions of LB bound by domain.",
			},
			"flow_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      FLOW_MODE_0,
				ValidateFunc: validateAllowedIntValue(FLOW_MODE_STATUS),
				Description:  "WAF traffic mode, 1 cleaning mode, 0 mirroring mode.",
			},
			"cls_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      CLS_STATUS_0,
				ValidateFunc: validateAllowedIntValue(CLS_STATUS),
				Description:  "Whether to enable access logs, 1 enable, 0 disable.",
			},
			"bot_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      BOT_STATUS_0,
				ValidateFunc: validateAllowedIntValue(BOT_STATUS),
				Description:  "Whether to enable bot, 1 enable, 0 disable.",
			},
			"api_safe_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      API_SAFE_STATUS_0,
				ValidateFunc: validateAllowedIntValue(API_SAFE_STATUS),
				Description:  "Whether to enable api safe, 1 enable, 0 disable.",
			},
			"alb_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ALB_TYPE_CLB,
				ValidateFunc: validateAllowedStringValue(ALB_TYPES),
				Description:  "Load balancer type: clb, apisix or tsegw, default clb.",
			},
			"ip_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "When is_cdn=3, this parameter needs to be filled in to indicate a custom header.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain id.",
			},
		},
	}
}

func resourceTencentCloudWafClbDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		request       = waf.NewCreateHostRequest()
		instanceID    string
		domain        string
		domainId      string
		wafStatus     uint64
		engine        uint64
		botStatus     uint64
		apiSafeStatus uint64
		clsStatus     uint64
		isCdn         int
		albType       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = helper.String(v.(string))
		instanceID = v.(string)
	}

	hostRecord := waf.HostRecord{}
	if v, ok := d.GetOk("domain"); ok {
		hostRecord.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOkExists("is_cdn"); ok {
		hostRecord.IsCdn = helper.IntUint64(v.(int))
		isCdn = v.(int)
	}

	if v, ok := d.GetOk("alb_type"); ok {
		hostRecord.AlbType = helper.String(v.(string))
		albType = v.(string)

		if albType == ALB_TYPE_CLB {
			if v, ok := d.GetOk("load_balancer_set"); ok {
				for _, item := range v.([]interface{}) {
					loadBalancerSetMap := item.(map[string]interface{})
					loadBalancer := waf.LoadBalancer{}
					if v, ok := loadBalancerSetMap["load_balancer_id"]; ok {
						loadBalancer.LoadBalancerId = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["load_balancer_name"]; ok {
						loadBalancer.LoadBalancerName = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["listener_id"]; ok {
						loadBalancer.ListenerId = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["listener_name"]; ok {
						loadBalancer.ListenerName = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["vip"]; ok {
						loadBalancer.Vip = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["vport"]; ok {
						loadBalancer.Vport = helper.IntUint64(v.(int))
					}

					if v, ok := loadBalancerSetMap["region"]; ok {
						loadBalancer.Region = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["protocol"]; ok {
						loadBalancer.Protocol = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["zone"]; ok {
						loadBalancer.Zone = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["numerical_vpc_id"]; ok {
						loadBalancer.NumericalVpcId = helper.IntInt64(v.(int))
					}

					if v, ok := loadBalancerSetMap["load_balancer_type"]; ok {
						loadBalancer.LoadBalancerType = helper.String(v.(string))
					}

					hostRecord.LoadBalancerSet = append(hostRecord.LoadBalancerSet, &loadBalancer)
				}
			} else {
				return fmt.Errorf("If `alb_type` is clb, `load_balancer_set` is required.")
			}
		} else {
			if _, ok := d.GetOk("load_balancer_set"); ok {
				return fmt.Errorf("If `alb_type` is apisix or tsegw, `load_balancer_set` is not supported.")
			}
		}
	}

	if v, ok := d.GetOk("region"); ok {
		hostRecord.Region = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("flow_mode"); ok {
		hostRecord.FlowMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("ip_headers"); ok {
		if isCdn == ISCDN_3 {
			ipHeadersSet := v.([]interface{})
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				hostRecord.IpHeaders = append(hostRecord.IpHeaders, &ipHeaders)
			}
		} else {
			return fmt.Errorf("If `is_cdn` is %d, not supported setting `ip_headers`.", isCdn)
		}
	}

	// check domain legal
	describeHostLimitRequest := waf.NewDescribeHostLimitRequest()
	describeHostLimitRequest.Domain = &domain
	describeHostLimitRequest.InstanceID = &instanceID
	describeHostLimitRequest.AlbType = &albType
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().DescribeHostLimit(describeHostLimitRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result.Response.Success != nil && *result.Response.Success.Code == DescribeHostLimitSuccess {
			return nil
		}

		e = fmt.Errorf("The current domain %s is illegal.", domain)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf clbDomain failed, reason:%+v", logId, err)
		return err
	}

	request.Host = &hostRecord
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().CreateHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf clbDomain failed, reason:%+v", logId, err)
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

	d.SetId(strings.Join([]string{instanceID, domain, domainId}, FILED_SP))

	// wait domain state
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDomainsById(ctx, instanceID, domain)
		if e != nil {
			return retryError(e)
		}

		if *result.State == 0 || *result.State == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("domain is still in state %d", *result.State))
	})

	if err != nil {
		return err
	}

	// set engine
	if v, ok := d.GetOkExists("engine"); ok {
		tmpEngine := v.(int)

		if tmpEngine != CLB_DOMAIN_ENGINE_20 {
			engine = uint64(tmpEngine)
			modifyHostModeRequest := waf.NewModifyHostModeRequest()
			modifyHostModeRequest.Domain = &domain
			modifyHostModeRequest.DomainId = &domainId
			modifyHostModeRequest.InstanceID = &instanceID
			modifyHostModeRequest.Mode = &engine

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostMode(modifyHostModeRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostModeRequest.GetAction(), modifyHostModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain engine failed, reason:%+v", logId, err)
				return err
			}
		}
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

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyBotStatus(modifyBotStatusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyBotStatusRequest.GetAction(), modifyBotStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain bot_status failed, reason:%+v", logId, err)
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

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyApiAnalyzeStatus(modifyApiAnalyzeStatusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyApiAnalyzeStatusRequest.GetAction(), modifyApiAnalyzeStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain api_safe_status failed, reason:%+v", logId, err)
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
					Edition:    common.StringPtr("clb-waf"),
					InstanceID: common.StringPtr(instanceID),
				},
			}
			modifyDomainsCLSStatusRequest.Status = &clsStatus

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyDomainsCLSStatus(modifyDomainsCLSStatusRequest)
				if e != nil {
					return retryError(e)
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

	// set waf status
	if v, ok := d.GetOkExists("status"); ok {
		tmpWafStatus := v.(int)

		if tmpWafStatus != CLB_DOMAIN_STATUS_1 {
			wafStatus = uint64(tmpWafStatus)
			modifyHostStatusRequest := waf.NewModifyHostStatusRequest()
			modifyHostStatusRequest.HostsStatus = []*waf.HostStatus{
				{
					Domain:     common.StringPtr(domain),
					DomainId:   common.StringPtr(domainId),
					Status:     common.Uint64Ptr(wafStatus),
					InstanceID: common.StringPtr(instanceID),
				},
			}

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostStatus(modifyHostStatusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostStatusRequest.GetAction(), modifyHostStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafClbDomainRead(d, meta)
}

func resourceTencentCloudWafClbDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]
	domainId := idSplit[2]

	domainInfo, err := service.DescribeDomainsById(ctx, instanceID, domain)
	if err != nil {
		return err
	}

	if domainInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DescribeDomains` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainInfo.InstanceId != nil {
		_ = d.Set("instance_id", domainInfo.InstanceId)
	}

	if domainInfo.Domain != nil {
		_ = d.Set("domain", domainInfo.Domain)
	}

	if domainInfo.Status != nil {
		_ = d.Set("status", domainInfo.Status)
	}

	if domainInfo.Engine != nil {
		_ = d.Set("engine", domainInfo.Engine)
	}

	if domainInfo.AlbType != nil {
		_ = d.Set("alb_type", domainInfo.AlbType)

		if *domainInfo.AlbType == ALB_TYPE_CLB {
			if domainInfo.LoadBalancerSet != nil {
				loadBalancerSetList := []interface{}{}
				for _, loadBalancerSet := range domainInfo.LoadBalancerSet {
					loadBalancerSetMap := map[string]interface{}{}

					if loadBalancerSet.LoadBalancerId != nil {
						loadBalancerSetMap["load_balancer_id"] = loadBalancerSet.LoadBalancerId
					}

					if loadBalancerSet.LoadBalancerName != nil {
						loadBalancerSetMap["load_balancer_name"] = loadBalancerSet.LoadBalancerName
					}

					if loadBalancerSet.ListenerId != nil {
						loadBalancerSetMap["listener_id"] = loadBalancerSet.ListenerId
					}

					if loadBalancerSet.ListenerName != nil {
						loadBalancerSetMap["listener_name"] = loadBalancerSet.ListenerName
					}

					if loadBalancerSet.Vip != nil {
						loadBalancerSetMap["vip"] = loadBalancerSet.Vip
					}

					if loadBalancerSet.Vport != nil {
						loadBalancerSetMap["vport"] = loadBalancerSet.Vport
					}

					if loadBalancerSet.Region != nil {
						loadBalancerSetMap["region"] = loadBalancerSet.Region
					}

					if loadBalancerSet.Protocol != nil {
						loadBalancerSetMap["protocol"] = loadBalancerSet.Protocol
					}

					if loadBalancerSet.Zone != nil {
						loadBalancerSetMap["zone"] = loadBalancerSet.Zone
					}

					if loadBalancerSet.NumericalVpcId != nil {
						loadBalancerSetMap["numerical_vpc_id"] = loadBalancerSet.NumericalVpcId
					}

					if loadBalancerSet.LoadBalancerType != nil {
						loadBalancerSetMap["load_balancer_type"] = loadBalancerSet.LoadBalancerType
					}

					loadBalancerSetList = append(loadBalancerSetList, loadBalancerSetMap)
				}

				_ = d.Set("load_balancer_set", loadBalancerSetList)
			}
		}
	}

	if domainInfo.Region != nil {
		_ = d.Set("region", domainInfo.Region)
	}

	if domainInfo.FlowMode != nil {
		_ = d.Set("flow_mode", domainInfo.FlowMode)
	}

	if domainInfo.ClsStatus != nil {
		_ = d.Set("cls_status", domainInfo.ClsStatus)
	}

	if domainInfo.BotStatus != nil {
		_ = d.Set("bot_status", domainInfo.BotStatus)
	}

	if domainInfo.ApiStatus != nil {
		_ = d.Set("api_safe_status", domainInfo.ApiStatus)
	}

	if domainInfo.DomainId != nil {
		_ = d.Set("domain_id", domainInfo.DomainId)
	}

	clbInfo, err := service.DescribeWafClbDomainById(ctx, instanceID, domain, domainId)
	if err != nil {
		return err
	}

	if clbInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DescribeDomainDetailsClb` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clbInfo.IsCdn != nil {
		_ = d.Set("is_cdn", clbInfo.IsCdn)
	}

	if clbInfo.IpHeaders != nil {
		_ = d.Set("ip_headers", clbInfo.IpHeaders)
	}

	return nil
}

func resourceTencentCloudWafClbDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		modifyHostRequest = waf.NewModifyHostRequest()
		isCdn             int
		wafStatus         uint64
		flowMode          uint64
		engine            uint64
		botStatus         uint64
		apiSafeStatus     uint64
		clsStatus         uint64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]
	domainId := idSplit[2]

	immutableArgs := []string{"instance_id", "domain", "alb_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	// set waf status
	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			tmpWafStatus := v.(int)
			// open first
			if tmpWafStatus == CLB_DOMAIN_STATUS_1 {
				wafStatus = uint64(v.(int))
				modifyHostStatusRequest := waf.NewModifyHostStatusRequest()
				modifyHostStatusRequest.HostsStatus = []*waf.HostStatus{
					{
						Domain:     common.StringPtr(domain),
						DomainId:   common.StringPtr(domainId),
						Status:     common.Uint64Ptr(wafStatus),
						InstanceID: common.StringPtr(instanceID),
					},
				}

				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostStatus(modifyHostStatusRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostStatusRequest.GetAction(), modifyHostStatusRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify waf clbDomain status failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	modifyHostRequest.InstanceID = &instanceID

	hostRecord := waf.HostRecord{}
	hostRecord.Domain = &domain
	hostRecord.DomainId = &domainId

	if v, ok := d.GetOkExists("is_cdn"); ok {
		hostRecord.IsCdn = helper.IntUint64(v.(int))
		isCdn = v.(int)
	}

	if v, ok := d.GetOk("alb_type"); ok {
		hostRecord.AlbType = helper.String(v.(string))
		albType := v.(string)

		if albType == ALB_TYPE_CLB {
			if v, ok := d.GetOk("load_balancer_set"); ok {
				for _, item := range v.([]interface{}) {
					loadBalancerSetMap := item.(map[string]interface{})
					loadBalancer := waf.LoadBalancer{}
					if v, ok := loadBalancerSetMap["load_balancer_id"]; ok {
						loadBalancer.LoadBalancerId = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["load_balancer_name"]; ok {
						loadBalancer.LoadBalancerName = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["listener_id"]; ok {
						loadBalancer.ListenerId = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["listener_name"]; ok {
						loadBalancer.ListenerName = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["vip"]; ok {
						loadBalancer.Vip = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["vport"]; ok {
						loadBalancer.Vport = helper.IntUint64(v.(int))
					}

					if v, ok := loadBalancerSetMap["region"]; ok {
						loadBalancer.Region = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["protocol"]; ok {
						loadBalancer.Protocol = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["zone"]; ok {
						loadBalancer.Zone = helper.String(v.(string))
					}

					if v, ok := loadBalancerSetMap["numerical_vpc_id"]; ok {
						loadBalancer.NumericalVpcId = helper.IntInt64(v.(int))
					}

					if v, ok := loadBalancerSetMap["load_balancer_type"]; ok {
						loadBalancer.LoadBalancerType = helper.String(v.(string))
					}

					hostRecord.LoadBalancerSet = append(hostRecord.LoadBalancerSet, &loadBalancer)
				}
			} else {
				return fmt.Errorf("If `alb_type` is clb, `load_balancer_set` is required.")
			}
		} else {
			if _, ok := d.GetOk("load_balancer_set"); ok {
				return fmt.Errorf("If `alb_type` is apisix or tsegw, `load_balancer_set` is not supported.")
			}
		}
	}

	if v, ok := d.GetOk("region"); ok {
		hostRecord.Region = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_headers"); ok {
		if isCdn == ISCDN_3 {
			ipHeadersSet := v.([]interface{})
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				hostRecord.IpHeaders = append(hostRecord.IpHeaders, &ipHeaders)
			}
		} else {
			return fmt.Errorf("If `is_cdn` is %d, not supported setting `ip_headers`.", isCdn)
		}
	}

	modifyHostRequest.Host = &hostRecord
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHost(modifyHostRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostRequest.GetAction(), modifyHostRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify waf clbDomain failed, reason:%+v", logId, err)
		return err
	}

	// set flow mode
	if d.HasChange("flow_mode") {
		if v, ok := d.GetOkExists("flow_mode"); ok {
			flowMode = uint64(v.(int))
			modifyHostFlowModeRequest := waf.NewModifyHostFlowModeRequest()
			modifyHostFlowModeRequest.Domain = &domain
			modifyHostFlowModeRequest.DomainId = &domainId
			modifyHostFlowModeRequest.InstanceID = &instanceID
			modifyHostFlowModeRequest.FlowMode = &flowMode

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostFlowMode(modifyHostFlowModeRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostFlowModeRequest.GetAction(), modifyHostFlowModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain flow_mode failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	// set engine
	if d.HasChange("engine") {
		if v, ok := d.GetOkExists("engine"); ok {
			engine = uint64(v.(int))
			modifyHostModeRequest := waf.NewModifyHostModeRequest()
			modifyHostModeRequest.Domain = &domain
			modifyHostModeRequest.DomainId = &domainId
			modifyHostModeRequest.InstanceID = &instanceID
			modifyHostModeRequest.Mode = &engine

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostMode(modifyHostModeRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostModeRequest.GetAction(), modifyHostModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain engine failed, reason:%+v", logId, err)
				return err
			}
		}
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

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyBotStatus(modifyBotStatusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyBotStatusRequest.GetAction(), modifyBotStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain bot_status failed, reason:%+v", logId, err)
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

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyApiAnalyzeStatus(modifyApiAnalyzeStatusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyApiAnalyzeStatusRequest.GetAction(), modifyApiAnalyzeStatusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf clbDomain api_safe_status failed, reason:%+v", logId, err)
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
				Edition:    common.StringPtr("clb-waf"),
				InstanceID: common.StringPtr(instanceID),
			},
		}
		modifyDomainsCLSStatusRequest.Status = &clsStatus

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyDomainsCLSStatus(modifyDomainsCLSStatusRequest)
			if e != nil {
				return retryError(e)
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

	// set waf status
	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			tmpWafStatus := v.(int)
			// close end
			if tmpWafStatus == CLB_DOMAIN_STATUS_0 {
				wafStatus = uint64(v.(int))
				modifyHostStatusRequest := waf.NewModifyHostStatusRequest()
				modifyHostStatusRequest.HostsStatus = []*waf.HostStatus{
					{
						Domain:     common.StringPtr(domain),
						DomainId:   common.StringPtr(domainId),
						Status:     common.Uint64Ptr(wafStatus),
						InstanceID: common.StringPtr(instanceID),
					},
				}

				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHostStatus(modifyHostStatusRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyHostStatusRequest.GetAction(), modifyHostStatusRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify waf clbDomain status failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudWafClbDomainRead(d, meta)
}

func resourceTencentCloudWafClbDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceID := idSplit[0]
	domain := idSplit[1]
	domainId := idSplit[2]

	if err := service.DeleteWafClbDomainById(ctx, instanceID, domain, domainId); err != nil {
		return err
	}

	return nil
}
