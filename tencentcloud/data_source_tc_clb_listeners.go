package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenersRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the CLB to be queried.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the listener to be queried.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description:  "Type of protocol within the listener, and available values are `TCP`, `UDP`, `HTTP`, `HTTPS` and `TCP_SSL`.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Description:  "Port of the CLB listener.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"listener_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of listeners of cloud load balancers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CLB.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the listener.",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CLB listener.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the listener. Available values are `HTTP`, `HTTPS`, `TCP`, `UDP`, `TCP_SSL`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the CLB listener.",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enabled.",
						},
						"health_check_time_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Response timeout of health check. The value range is 2-60 sec, and the default is `2` sec. Response timeout needs to be less than check interval. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration.",
						},
						"health_check_interval_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval time of health check. The value range is 2-300 sec, and the default is `5` sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold of health check, and the default is `3`. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealthy threshold of health check, and the default is `3`. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
						},
						"health_check_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol used for health check.",
						},
						"health_check_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check port is the port of the backend service.",
						},
						"health_check_http_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The HTTP version of the backend service.",
						},
						"health_check_http_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP health check code of TCP listener.",
						},
						"health_check_http_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP health check path of TCP listener.",
						},
						"health_check_http_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP health check domain of TCP listener.",
						},
						"health_check_http_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP health check method of TCP listener.",
						},
						"health_check_context_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health check protocol.",
						},
						"health_check_send_context": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "It represents the content of the request sent by the health check.",
						},
						"health_check_recv_context": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "It represents the result returned by the health check.",
						},
						"certificate_ssl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of certificate, and available values inclue `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the server certificate. It must be set when protocol is `HTTPS` or `TCP_SSL`. NOTES: only supported by listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
						},
						"certificate_ca_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the client certificate. It must be set when SSLMode is `mutual`. NOTES: only supported by listeners of `HTTPS` and `TCP_SSL` protocol.",
						},
						"session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time of session persistence within the CLB listener. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling method of the CLB listener, and available values are `WRR` and `LEAST_CONN`. The default is `WRR`. NOTES: The listener of 'HTTP' and `HTTPS` protocol additionally supports the `IP HASH` method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.",
						},
						"sni_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether SNI is enabled. NOTES: Only supported by `HTTPS` protocol.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbListenersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_listeners.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Get("clb_id").(string)

	params := make(map[string]interface{})
	params["clb_id"] = clbId
	if v, ok := d.GetOk("listener_id"); ok {
		listenerId := v.(string)
		params["listener_id"] = listenerId
		checkErr := ListenerIdCheck(listenerId)
		if checkErr != nil {
			return checkErr
		}
	}
	if v, ok := d.GetOk("port"); ok {
		params["port"] = v.(int)
	}
	if v, ok := d.GetOk("protocol"); ok {
		params["protocol"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var listeners []*clb.Listener
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeListenersByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		listeners = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listeners failed, reason:%+v", logId, err)
		return err
	}
	listenerList := make([]map[string]interface{}, 0, len(listeners))
	ids := make([]string, 0, len(listeners))
	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"clb_id":        clbId,
			"listener_id":   listener.ListenerId,
			"listener_name": listener.ListenerName,
			"protocol":      listener.Protocol,
			"port":          listener.Port,
		}
		if listener.SessionExpireTime != nil {
			mapping["session_expire_time"] = listener.SessionExpireTime
		}
		if listener.SniSwitch != nil {
			sniSwitch := false
			if *listener.SniSwitch == int64(1) {
				sniSwitch = true
			}
			mapping["sni_switch"] = sniSwitch
		}
		mapping["scheduler"] = listener.Scheduler
		if listener.HealthCheck != nil {
			health_check_switch := false
			if *listener.HealthCheck.HealthSwitch == int64(1) {
				health_check_switch = true
			}
			mapping["health_check_switch"] = health_check_switch
			mapping["health_check_time_out"] = listener.HealthCheck.TimeOut
			mapping["health_check_interval_time"] = listener.HealthCheck.IntervalTime
			mapping["health_check_health_num"] = listener.HealthCheck.HealthNum
			mapping["health_check_unhealth_num"] = listener.HealthCheck.UnHealthNum
			mapping["health_check_http_code"] = listener.HealthCheck.HttpCode
			mapping["health_check_http_path"] = listener.HealthCheck.HttpCheckPath
			mapping["health_check_http_domain"] = listener.HealthCheck.HttpCheckDomain
			mapping["health_check_http_method"] = listener.HealthCheck.HttpCheckMethod
			mapping["health_check_http_version"] = listener.HealthCheck.HttpVersion
			mapping["health_check_context_type"] = listener.HealthCheck.ContextType
			mapping["health_check_send_context"] = listener.HealthCheck.SendContext
			mapping["health_check_recv_context"] = listener.HealthCheck.RecvContext
			mapping["health_check_type"] = listener.HealthCheck.CheckType
			mapping["health_check_port"] = listener.HealthCheck.CheckPort
		}
		if listener.Certificate != nil {
			mapping["certificate_ssl_mode"] = listener.Certificate.SSLMode
			mapping["certificate_id"] = listener.Certificate.CertId
			if listener.Certificate.CertCaId != nil {
				mapping["certificate_ca_id"] = listener.Certificate.CertCaId
			}
		}
		listenerList = append(listenerList, mapping)
		ids = append(ids, *listener.ListenerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("listener_list", listenerList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB listener list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), listenerList); e != nil {
			return e
		}
	}

	return nil
}
