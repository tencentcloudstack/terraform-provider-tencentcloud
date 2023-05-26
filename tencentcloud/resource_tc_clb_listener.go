/*
Provides a resource to create a CLB listener.

Example Usage

HTTP Listener

```hcl
resource "tencentcloud_clb_listener" "HTTP_listener" {
  clb_id        = "lb-0lh5au7v"
  listener_name = "test_listener"
  port          = 80
  protocol      = "HTTP"
}
```

TCP/UDP Listener

```hcl
resource "tencentcloud_clb_listener" "TCP_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = 80
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_port          = 200
  health_check_type          = "HTTP"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
}
```

TCP/UDP Listener with tcp health check
```hcl
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "TCP"
  health_check_port          = 200
}
```

TCP/UDP Listener with http health check
```hcl
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "HTTP"
  health_check_http_domain   = "www.tencent.com"
  health_check_http_code     = 16
  health_check_http_version  = "HTTP/1.1"
  health_check_http_method   = "HEAD"
  health_check_http_path     = "/"
}
```

TCP/UDP Listener with customer health check
```hcl
resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "CUSTOM"
  health_check_context_type  = "HEX"
  health_check_send_context  = "0123456789ABCDEF"
  health_check_recv_context  = "ABCD"
  target_type                = "TARGETGROUP"
}
```

HTTPS Listener

```hcl
resource "tencentcloud_clb_listener" "HTTPS_listener" {
  clb_id               = "lb-0lh5au7v"
  listener_name        = "test_listener"
  port                 = "80"
  protocol             = "HTTPS"
  certificate_ssl_mode = "MUTUAL"
  certificate_id       = "VjANRdz8"
  certificate_ca_id    = "VfqO4zkB"
  sni_switch           = true
}
```

TCP SSL Listener

```hcl
resource "tencentcloud_clb_listener" "TCPSSL_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = "80"
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
  target_type                = "TARGETGROUP"
}
```
Import

CLB listener can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener.foo lb-7a0t6zqb#lbl-hh141sn9
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerCreate,
		Read:   resourceTencentCloudClbListenerRead,
		Update: resourceTencentCloudClbListenerUpdate,
		Delete: resourceTencentCloudClbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"scheduler": CLB_LISTENER_SCHEDULER_WRR,
			}),
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "ID of the CLB.",
			},
			"listener_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the CLB listener, and available values can only be Chinese characters, English letters, numbers, underscore and hyphen '-'.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Description:  "Port of the CLB listener.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description:  "Type of protocol within the listener. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS` and `TCP_SSL`.",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether health check is enabled.",
			},
			"health_check_time_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 60),
				Description:  "Response timeout of health check. Valid value ranges: [2~60] sec. Default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of `TCP`,`UDP`,`TCP_SSL` protocol.",
			},
			"health_check_interval_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(5, 300),
				Description:  "Interval time of health check. Valid value ranges: [5~300] sec. and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is `3`. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description: "Unhealthy threshold of health check, and the default is `3`. " +
					"If a success result is returned for the health check 3 consecutive times, " +
					"the CVM is identified as unhealthy. The value range is [2-10]. " +
					"NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, " +
					"HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(HEALTH_CHECK_TYPE),
				Description:  "Protocol used for health check. Valid values: `CUSTOM`, `TCP`, `HTTP`.",
			},
			"health_check_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Description: "The health check port is the port of the backend service by default. " +
					"Unless you want to specify a specific port, it is recommended to leave it blank. " +
					"Only applicable to TCP/UDP listener.",
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(HTTP_VERSION),
				Description: "The HTTP version of the backend service. When the value of `health_check_type` of " +
					"the health check protocol is `HTTP`, this field is required. " +
					"Valid values: `HTTP/1.0`, `HTTP/1.1`.",
			},
			"health_check_http_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 31),
				Description: "HTTP health check code of TCP listener, Valid value ranges: [1~31]. When the value of `health_check_type` of " +
					"the health check protocol is `HTTP`, this field is required. Valid values: `1`, `2`, `4`, `8`, `16`. " +
					"`1` means http_1xx, `2` means http_2xx, `4` means http_3xx, `8` means http_4xx, `16` means http_5xx." +
					"If you want multiple return codes to indicate health, need to add the corresponding values.",
			},
			"health_check_http_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP health check path of TCP listener.",
			},
			"health_check_http_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP health check domain of TCP listener.",
			},
			"health_check_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(CLB_HTTP_METHOD),
				Description:  "HTTP health check method of TCP listener. Valid values: `HEAD`, `GET`.",
			},
			"health_check_context_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CONTEX_TYPE),
				Description: "Health check protocol. When the value of `health_check_type` of the health check protocol is `CUSTOM`, " +
					"this field is required, which represents the input format of the health check. " +
					"Valid values: `HEX`, `TEXT`.",
			},
			"health_check_send_context": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 500),
				Description: "It represents the content of the request sent by the health check. " +
					"When the value of `health_check_type` of the health check protocol is `CUSTOM`, " +
					"this field is required. Only visible ASCII characters are allowed and the maximum length is 500. " +
					"When `health_check_context_type` value is `HEX`, " +
					"the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` " +
					"and the length must be even digits.",
			},
			"health_check_recv_context": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 500),
				Description: "It represents the result returned by the health check. " +
					"When the value of `health_check_type` of the health check protocol is `CUSTOM`, " +
					"this field is required. Only ASCII visible characters are allowed and the maximum length is 500. " +
					"When `health_check_context_type` value is `HEX`, " +
					"the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` " +
					"and the length must be even digits.",
			},
			"certificate_ssl_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CERT_SSL_MODE),
				Description:  "Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the server certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
			},
			"certificate_ca_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the client certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when the ssl mode is `MUTUAL`.",
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(30, 3600),
				Description:  "Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as `WRR`, and not available when listener protocol is `TCP_SSL`. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Default:      CLB_LISTENER_SCHEDULER_WRR,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method of the CLB listener, and available values are 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of `HTTP` and `HTTPS` protocol additionally supports the `IP Hash` method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"sni_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicates whether SNI is enabled, and only supported with protocol `HTTPS`. If enabled, you can set a certificate for each rule in `tencentcloud_clb_listener_rule`, otherwise all rules have a certificate.",
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{CLB_TARGET_TYPE_NODE, CLB_TARGET_TYPE_TARGETGROUP}),
				Description:  "Backend target type. Valid values: `NODE`, `TARGETGROUP`. `NODE` means to bind ordinary nodes, `TARGETGROUP` means to bind target group. NOTES: TCP/UDP/TCP_SSL listener must configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			//computed
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of this CLB listener.",
			},
		},
	}
}

func resourceTencentCloudClbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Get("clb_id").(string)
	listenerName := d.Get("listener_name").(string)
	request := clb.NewCreateListenerRequest()

	request.LoadBalancerId = helper.String(clbId)
	request.ListenerNames = []*string{&listenerName}

	port := int64(d.Get("port").(int))
	ports := []*int64{&port}
	request.Ports = ports
	protocol := d.Get("protocol").(string)
	request.Protocol = helper.String(protocol)

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)

	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		request.Certificate = certificateInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is TCPSSL")
		}
	}
	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if v == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS")
		}
		scheduler = v.(string)
		request.Scheduler = helper.String(scheduler)
	}
	if v, ok := d.GetOk("target_type"); ok {
		targetType := v.(string)
		request.TargetType = &targetType
	} else if protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP || protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
		targetType := CLB_TARGET_TYPE_NODE
		request.TargetType = &targetType
	}

	if v, ok := d.GetOk("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		request.SessionExpireTime = &vv
	}
	if v, ok := d.GetOkExists("sni_switch"); ok {
		if protocol != CLB_LISTENER_PROTOCOL_HTTPS {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: sni_switch can only be set with protocol HTTPS ")
		} else {
			vv := v.(bool)
			vvv := int64(0)
			if vv {
				vvv = 1
			} else {
				if !certificateSetFlag {
					return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is HTTPS")
				}
			}
			request.SniSwitch = &vvv
		}
	}
	var response *clb.CreateListenerResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateListener(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB listener failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.ListenerIds) < 1 {
		return fmt.Errorf("[CHECK][CLB listener][Create] check: Response error, listener id is null")
	}
	listenerId := *response.Response.ListenerIds[0]

	//this ID style changes since terraform 1.47.0
	d.SetId(clbId + FILED_SP + listenerId)
	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	resourceId := d.Id()
	var listenerId = resourceId
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)

	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("the old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 2 && clbId == "" {
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 2 && clbId != "" {
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 2 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	var instance *clb.Listener
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", instance.ListenerId)
	_ = d.Set("port", instance.Port)
	_ = d.Set("protocol", instance.Protocol)
	if instance.ListenerName != nil {
		_ = d.Set("listener_name", instance.ListenerName)
	}
	if instance.TargetType != nil {
		_ = d.Set("target_type", instance.TargetType)
	}
	if instance.SessionExpireTime != nil {
		_ = d.Set("session_expire_time", instance.SessionExpireTime)
	}
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL || *instance.Protocol == CLB_LISTENER_PROTOCOL_UDP {
		_ = d.Set("scheduler", instance.Scheduler)
	}
	_ = d.Set("sni_switch", *instance.SniSwitch > 0)

	//health check
	if instance.HealthCheck != nil {
		healthCheckSwitch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			healthCheckSwitch = true
		}
		_ = d.Set("health_check_switch", healthCheckSwitch)
		if instance.HealthCheck.IntervalTime != nil {
			_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		}
		if instance.HealthCheck.TimeOut != nil {
			_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
		}
		if instance.HealthCheck.HealthNum != nil {
			_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		}
		if instance.HealthCheck.UnHealthNum != nil {
			_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		}
		if instance.HealthCheck.CheckPort != nil {
			_ = d.Set("health_check_port", instance.HealthCheck.CheckPort)
		}
		if instance.HealthCheck.CheckType != nil {
			_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		}
		if instance.HealthCheck.HttpCode != nil {
			_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		}
		if instance.HealthCheck.HttpCheckPath != nil {
			_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		}
		if instance.HealthCheck.HttpCheckDomain != nil {
			_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		}
		if instance.HealthCheck.HttpCheckMethod != nil {
			_ = d.Set("health_check_http_method", instance.HealthCheck.HttpCheckMethod)
		}
		if instance.HealthCheck.HttpVersion != nil {
			_ = d.Set("health_check_http_version", instance.HealthCheck.HttpVersion)
		}
		if instance.HealthCheck.ContextType != nil {
			_ = d.Set("health_check_context_type", instance.HealthCheck.ContextType)
		}
		if instance.HealthCheck.SendContext != nil {
			_ = d.Set("health_check_send_context", instance.HealthCheck.SendContext)
		}
		if instance.HealthCheck.RecvContext != nil {
			_ = d.Set("health_check_recv_context", instance.HealthCheck.RecvContext)
		}
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

func resourceTencentCloudClbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)
	changed := false
	scheduler := ""
	listenerName := ""
	sessionExpireTime := 0
	protocol := d.Get("protocol").(string)

	request := clb.NewModifyListenerRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)

	if d.HasChange("listener_name") {
		listenerName = d.Get("listener_name").(string)
		request.ListenerName = helper.String(listenerName)
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP || protocol == CLB_LISTENER_PROTOCOL_TCPSSL) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS", listenerId)
		}
		request.Scheduler = helper.String(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set when scheduler is WRR", listenerId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		changed = true
		request.Certificate = certificateInput
	}

	if d.HasChange("target_type") {
		changed = true
		targetType := d.Get("target_type").(string)
		request.TargetType = helper.String(targetType)
	}

	if changed {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyListener(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.delete")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteListenerById(ctx, clbId, listenerId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
