/*
Use this data source to query detailed information of CLB listener

Example Usage

```hcl
data "tencentcloud_clb_listeners" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  protocol    = "TCP"
  port        = 80
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudClbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenersRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the CLB to be queried.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the listener to be queried",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description:  "Protocol of the listener. Available values are 'HTTP', 'HTTPS', 'TCP', 'UDP'('TCP_SSL' is in the internal test, please apply if you need to use). ",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Description:  "Port of the listener.",
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
							Description: "Name of the listener.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the listener. Available values are 'HTTP', 'HTTPS', 'TCP', 'UDP'.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the listener.",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enabled.",
						},
						"health_check_time_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Response timeout of health check. The value range is 2-60 sec, and the default is 2 sec. Response timeout needs to be less than check interval.",
						},
						"health_check_interval_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10.",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10.",
						},
						"certificate_ssl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of SSL Mode, and available values inclue 'UNIDRECTIONAL', 'MUTUAL'.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the server certificate. If not specified, the content, key, and name of the server certificate must be set. NOTES: only supported by listeners of protocol 'HTTPS'.",
						},
						"certificate_ca_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the client certificate. If not specified, the content, key, name of client certificate must be set when SSLMode is 'mutual'. NOTES: only supported by listeners of protocol 'HTTPS'.",
						},
						"session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time of session persistence within the CLB listener.",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'.",
						},
						"sni_switch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates whether SNI is enabled, and only supported with protocol 'HTTPS'",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbListenersRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_clb_listeners.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Get("clb_id").(string)

	params := make(map[string]interface{})
	params["clb_id"] = clbId
	if v, ok := d.GetOk("listener_id"); ok {
		params["listener_id"] = v.(string)
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
	listeners, err := clbService.DescribeListenersByFilter(ctx, params)
	if err != nil {
		return err
	}

	listenerList := make([]map[string]interface{}, 0, len(listeners))
	ids := make([]string, 0, len(listeners))
	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"clb_id":        clbId,
			"listener_id":   *listener.ListenerId,
			"listener_name": *listener.ListenerName,
			"protocol":      *listener.Protocol,
			"port":          *listener.Port,
		}
		if listener.SessionExpireTime != nil {
			mapping["session_expire_time"] = *listener.SessionExpireTime
		}
		if listener.SniSwitch != nil {
			mapping["sni_switch"] = *listener.SniSwitch
		}
		if listener.Scheduler != nil {
			mapping["scheduler"] = *listener.Scheduler
		}
		if listener.HealthCheck != nil {
			health_check_switch := false
			if *listener.HealthCheck.HealthSwitch == int64(1) {
				health_check_switch = true
			}
			mapping["health_check_switch"] = health_check_switch
			mapping["health_check_time_out"] = *listener.HealthCheck.TimeOut
			mapping["health_check_interval_time"] = *listener.HealthCheck.IntervalTime
			mapping["health_check_health_num"] = *listener.HealthCheck.HealthNum
			mapping["health_check_unhealth_num"] = *listener.HealthCheck.UnHealthNum
		}
		if listener.Certificate != nil {
			mapping["certificate_ssl_mode"] = *listener.Certificate.SSLMode
			mapping["certificate_id"] = *listener.Certificate.CertId
			mapping["certificate_ca_id"] = *listener.Certificate.CertCaId
		}
		listenerList = append(listenerList, mapping)
		ids = append(ids, *listener.ListenerId+"#"+clbId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("listener_list", listenerList); err != nil {
		log.Printf("[CRITAL]%s provider set clb listener list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), listenerList); err != nil {
			return err
		}
	}

	return nil
}
