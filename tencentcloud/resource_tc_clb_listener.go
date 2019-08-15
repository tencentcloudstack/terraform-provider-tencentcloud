/*
Provides a resource to create a CLB listener.

Example Usage

HTTP Listener

```hcl
resource "tencentcloud_clb_listener" "HTTP_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = 80
  protocol                   = "HTTP"
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
}
```

HTTPS Listener

```hcl
resource "tencentcloud_clb_listener" "HTTPS_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = "80"
  protocol                   = "HTTPS"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjAYq9xc"
  certificate_ca_id          = "VfqcL1ME"
  sni_switch                 = true
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
  certificate_id             = "VjAYq9xc"
  certificate_ca_id          = "VfqcL1ME"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerCreate,
		Read:   resourceTencentCloudClbListenerRead,
		Update: resourceTencentCloudClbListenerUpdate,
		Delete: resourceTencentCloudClbListenerDelete,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Id of the CLB.",
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
				Description:  "Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'.",
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
				Description:  "Response timeout of health check. The value range is 2-60 sec, and the default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of 'TCP','UDP','TCP_SSL' protocol.",
			},
			"health_check_interval_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(5, 300),
				Description:  "Interval time of health check. The value range is 5-300 sec, and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"certificate_ssl_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CERT_SSL_MODE),
				Description:  "Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the server certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.",
			},
			"certificate_ca_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when the ssl mode is 'MUTUAL'.",
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(30, 3600),
				Description:  "Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR', and not available when listener protocol is 'TCP_SSL'. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Default:      CLB_LISTENER_SCHEDULER_WRR,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of HTTP and 'HTTPS' protocol additionally supports the 'IP Hash' method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
			},
			"sni_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicates whether SNI is enabled, and only supported with protocol 'HTTPS'.",
			},
		},
	}
}

func resourceTencentCloudClbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Get("clb_id").(string)
	listenerName := d.Get("listener_name").(string)
	request := clb.NewCreateListenerRequest()

	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerNames = []*string{&listenerName}

	port := int64(d.Get("port").(int))
	ports := []*int64{&port}
	request.Ports = ports
	protocol := d.Get("protocol").(string)
	request.Protocol = stringToPointer(protocol)

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d)

	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		request.Certificate = certificateInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_HTTPS || protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("certificated need to be set when protocol is HTTPS/TCPSSL")
		}
	}
	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if v == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS")
		}
		scheduler = v.(string)
		request.Scheduler = stringToPointer(scheduler)
	}

	if v, ok := d.GetOk("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		request.SessionExpireTime = &vv
	}
	if v, ok := d.GetOk("sni_switch"); ok {
		if protocol != CLB_LISTENER_PROTOCOL_HTTPS {
			return fmt.Errorf("sni_switch can only be set with protocol HTTPS ")
		} else {
			vv := v.(bool)
			vvv := int64(0)
			if vv {
				vvv = 1
			}
			request.SniSwitch = &vvv
		}
	}
	var response *clb.CreateListenerResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateListener(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId

			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(retryErr)
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb listener failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if len(response.Response.ListenerIds) < 1 {
		return fmt.Errorf("listener id is wrong")
	}
	listenerId := *response.Response.ListenerIds[0]
	d.SetId(listenerId)

	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	clbId := d.Get("clb_id").(string)
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *clb.Listener
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeListenerById(ctx, d.Id(), clbId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read clb listener failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("clb_id", clbId)
	d.Set("listener_name", instance.ListenerName)
	d.Set("port", instance.Port)
	d.Set("protocol", instance.Protocol)
	d.Set("session_expire_time", instance.SessionExpireTime)
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL || *instance.Protocol == CLB_LISTENER_PROTOCOL_UDP {
		d.Set("scheduler", instance.Scheduler)
	}
	d.Set("sni_switch", instance.SniSwitch)

	//health check
	if instance.HealthCheck != nil {
		health_check_switch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			health_check_switch = true
		}
		d.Set("health_check_switch", health_check_switch)
		d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
		d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		d.Set("health_check_health_num ", instance.HealthCheck.HealthNum)
		d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
	}

	if instance.Certificate != nil {
		d.Set("certificate_ssl_mode", instance.Certificate.SSLMode)
		d.Set("certificate_id", instance.Certificate.CertId)
		d.Set("certificate_ca_id", instance.Certificate.CertCaId)
	}

	return nil
}

func resourceTencentCloudClbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	listenerId := d.Id()
	clbId := d.Get("clb_id").(string)
	changed := false
	scheduler := ""
	listenerName := ""
	sessionExpireTime := 0
	protocol := d.Get("protocol").(string)

	request := clb.NewModifyListenerRequest()
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)

	if d.HasChange("listener_name") {
		listenerName = d.Get("listener_name").(string)
		request.ListenerName = stringToPointer(listenerName)
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP || protocol == CLB_LISTENER_PROTOCOL_TCPSSL) {
			return fmt.Errorf("Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS")
		}
		if scheduler == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS")
		}
		request.Scheduler = stringToPointer(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("session_expire_time can only be set when scheduler is WRR")
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

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		changed = true
		request.Certificate = certificateInput
	}

	if changed {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyListener(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update clb listener failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return nil
}

func resourceTencentCloudClbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_listener.delete")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	listenerId := d.Id()
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteListenerById(ctx, clbId, listenerId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete clb listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
