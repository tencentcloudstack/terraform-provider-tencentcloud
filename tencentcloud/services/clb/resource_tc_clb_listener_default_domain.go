package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbListenerDefaultDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerDefaultDomainCreate,
		Read:   resourceTencentCloudClbListenerDefaultDomainRead,
		Update: resourceTencentCloudClbListenerDefaultDomainUpdate,
		Delete: resourceTencentCloudClbListenerDefaultDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB instance.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    true,
				Description: "Domain name of the listener rule. Single domain rules are passed to `domain`, and multi domain rules are passed to `domains`.",
			},
		},
	}
}

func resourceTencentCloudClbListenerDefaultDomainCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = clb.NewModifyDomainAttributesRequest()
		response *clb.ModifyDomainAttributesResponse
	)

	if v, ok := d.GetOk("clb_id"); ok {
		request.LoadBalancerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listener_id"); ok {
		request.ListenerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyDomainAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cls topic failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.RequestId

	d.SetId(id)
	return resourceTencentCloudClbListenerDefaultDomainRead(d, meta)
}

func resourceTencentCloudClbListenerDefaultDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()
	var locationId = resourceId
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("The old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 3 && clbId == "" {
		locationId = items[2]
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 3 && clbId != "" {
		locationId = items[2]
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 3 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	//this function is not supported by api, need to be travelled
	filter := map[string]string{"rule_id": locationId, "listener_id": listenerId, "clb_id": clbId}
	var instances []*clb.RuleOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRulesByFilter(ctx, filter)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	if len(instances) == 0 {
		d.SetId("")
		return nil
	}

	instance := instances[0]
	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	if instance.Domain != nil {
		_ = d.Set("domain", instance.Domain)
	}

	if instance.Domains != nil {
		_ = d.Set("domains", helper.StringsInterfaces(instance.Domains))
	}

	_ = d.Set("rule_id", instance.LocationId)
	_ = d.Set("url", instance.Url)
	_ = d.Set("scheduler", instance.Scheduler)
	_ = d.Set("session_expire_time", instance.SessionExpireTime)
	_ = d.Set("target_type", instance.TargetType)
	_ = d.Set("forward_type", instance.ForwardType)
	_ = d.Set("http2_switch", instance.Http2)

	if instance.QuicStatus != nil {
		if *instance.QuicStatus == "QUIC_ACTIVE" {
			_ = d.Set("quic", true)
		} else {
			_ = d.Set("quic", false)
		}
	}

	//health check
	if instance.HealthCheck != nil {
		health_check_switch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			health_check_switch = true
		}
		_ = d.Set("health_check_switch", health_check_switch)
		_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		_ = d.Set("health_check_http_method", helper.String(strings.ToUpper(*instance.HealthCheck.HttpCheckMethod)))
		_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
	}

	if instance.Certificate != nil {
		_ = d.Set("certificate_ssl_mode", instance.Certificate.SSLMode)
		_ = d.Set("certificate_id", instance.Certificate.CertId)
		if instance.Certificate.CertCaId != nil {
			_ = d.Set("certificate_ca_id", instance.Certificate.CertCaId)
		}
	}

	return nil
}

func resourceTencentCloudClbListenerDefaultDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	protocol := ""
	//get listener protocol
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instance, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		protocol = *(instance.Protocol)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s get CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	changed := false
	url := ""
	scheduler := ""
	sessionExpireTime := 0

	request := clb.NewModifyRuleRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)
	request.LocationId = helper.String(locationId)
	if d.HasChange("url") {
		changed = true
		url = d.Get("url").(string)
		request.Url = helper.String(url)
	}

	if d.HasChange("forward_type") {
		changed = true
		request.ForwardType = helper.String(d.Get("forward_type").(string))
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS", locationId)
		}
		request.Scheduler = helper.String(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", locationId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set when scheduler is WRR", locationId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_RULE)
	if healthErr != nil {
		return healthErr
	}

	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyRule(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	//modify domain and ssl
	domainChanged := false
	domainRequest := clb.NewModifyDomainAttributesRequest()
	if d.HasChange("domain") {
		old, new := d.GetChange("domain")
		domainChanged = true
		domainRequest.Domain = helper.String(old.(string))
		domainRequest.NewDomain = helper.String(new.(string))
	} else {
		domainRequest.Domain = helper.String(d.Get("domain").(string))
	}

	if d.HasChange("certificate_id") || d.HasChange("certificate_ca_id ") || d.HasChange("certificate_ssl_mode") {
		domainChanged = true
		certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
		if certErr != nil {
			return certErr
		}
		if certificateSetFlag {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainRequest.Certificate = certificateInput
		}
	}

	if d.HasChange("http2_switch") {
		if v, ok := d.GetOkExists("http2_switch"); ok {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainChanged = true
			domainRequest.Http2 = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("quic") {
		domainChanged = true
		if v, ok := d.GetOkExists("quic"); ok {
			domainRequest.Quic = helper.Bool(v.(bool))
		}
	}

	if domainChanged {
		domainRequest.ListenerId = &listenerId
		domainRequest.LoadBalancerId = &clbId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyDomainAttributes(domainRequest)
			if e != nil {
				if err := processRetryErrMsg(e); err != nil {
					return err
				}
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}

func resourceTencentCloudClbListenerDefaultDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_domain_default.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteRuleById(ctx, clbId, listenerId, locationId)
		if e != nil {
			if err := processRetryErrMsg(e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener rule failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
