/*
Provides a resource to create a CLB listener.

Example Usage

```hcl
resource "tencentcloud_clb_listener" "clb_listener" {
  clb_id                     = "lb-k2zjp9lv"
  listener_name              = "mylistener"
  port                       = 80
  protocol                   = "HTTP"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
    certificate_ssl_mode       = "MUTUAL"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "mycert server ID "
  certificate_ca_id          = "mycert ca ID"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```

Import

CLB listener can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_listener.foo lbl-qckdffns#lb-p7nlgs4t

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerCreate,
		Read:   resourceTencentCloudClbListenerRead,
		Update: resourceTencentCloudClbListenerUpdate,
		Delete: resourceTencentCloudClbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description:  "Type of protocol within the listener, and available values include TCP, UDP, HTTP, HTTPS and TCP_SSL. NOTES: TCP_SSL is testing internally, please apply if you need to use.",
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
				Description:  "Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(2, 10),
				Description:  "Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is 2-10.",
			},
			"certificate_ssl_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CERT_SSL_MODE),
				Description:  "Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' protocol.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: " ID of the server certificate. NOTES: Only supports listeners of 'HTTPS' protocol.",
			},
			"certificate_ca_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the client certificate. NOTES: Only supports listeners of 'HTTPS' protocol.",
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(30, 300),
				Description:  "Time of session persistence within the CLB listener. NOTES: Only supports listeners of 'WRR' scheduler.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'. NOTES: The listener of HTTP and 'HTTPS' protocol additionally supports the 'IP Hash' method.",
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
	defer LogElapsed("resource.tencentcloud_clb_listener.create")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := GetLogId(nil)
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
	if protocol == "TCP_SSL" {
		return fmt.Errorf("TCP_SSL protocol type needs manual application")
	} else {
		request.Protocol = stringToPointer(protocol)
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag == true {
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d)

	if certErr != nil {
		return certErr
	}
	if certificateSetFlag == true {
		request.Certificate = certificateInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_HTTPS {
			return fmt.Errorf("certificated need to be set when protocol is HTTPS")
		}
	}
	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("Scheduler can only be set with listener protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
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
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateListener(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId := *response.Response.RequestId

		retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}
	if len(response.Response.ListenerIds) < 1 {
		return fmt.Errorf("load balancer id is nil")
	}
	listenerId := *response.Response.ListenerIds[0]
	d.SetId(listenerId + "#" + clbId)
	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("resource.tencentcloud_clb_listener.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instance, err := clbService.DescribeListenerById(ctx, d.Id())
	if err != nil {
		return err
	}

	d.Set("clb_id", items[1])
	d.Set("listener_name", instance.ListenerName)
	d.Set("port", instance.Port)
	d.Set("protocol", instance.Protocol)
	d.Set("session_expire_time", instance.SessionExpireTime)
	d.Set("scheduler", instance.Scheduler)
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
	defer LogElapsed("resource.tencentcloud_clb_listener.update")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
	}

	listenerId := items[0]
	clbId := items[1]
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
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("Scheduler can only be set with listener protocol TCP/UDP or rule of listener HTTP/HTTPS")
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
	if healthSetFlag == true {
		changed = true
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag == true {
		changed = true
		request.Certificate = certificateInput
	}

	if changed {

		response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyListener(request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return err
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId := *response.Response.RequestId
			retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return retryErr
			}
		}
	}

	return nil
}

func resourceTencentCloudClbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("resource.tencentcloud_clb_listener.delete")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
	}
	listenerId := items[0]
	clbId := items[1]

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := clbService.DeleteListenerById(ctx, clbId, listenerId)
	if err != nil {
		log.Printf("[CRITAL]%s reason[%s]\n", logId, err.Error())
		return err
	}

	return nil
}
