package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func dataSourceTencentCloudGaapLayer4Listeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapLayer4ListenersRead,
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP"}),
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePort,
			},

			// computed
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"realserver_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapLayer4ListenersRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_layer4_listeners.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)

	var (
		listenerId *string
		name       *string
		port       *int
		ids        []string
		listeners  []map[string]interface{}
	)

	if raw, ok := d.GetOk("listener_id"); ok {
		listenerId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("listener_name"); ok {
		name = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("port"); ok {
		port = common.IntPtr(raw.(int))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	switch protocol {
	case "TCP":
		tcpListeners, err := service.DescribeTCPListeners(ctx, proxyId, listenerId, name, port)
		if err != nil {
			return err
		}

		ids = make([]string, 0, len(tcpListeners))
		listeners = make([]map[string]interface{}, 0, len(tcpListeners))

		for _, ls := range tcpListeners {
			if ls.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if ls.ListenerName == nil {
				return errors.New("listener name is nil")
			}
			if ls.Port == nil {
				return errors.New("listener port is nil")
			}
			if ls.RealServerType == nil {
				return errors.New("listener realserver type is nil")
			}
			if ls.ListenerStatus == nil {
				return errors.New("listener realserver type is nil")
			}
			if ls.Scheduler == nil {
				return errors.New("listener scheduler is nil")
			}
			if ls.HealthCheck == nil {
				return errors.New("listener health check is nil")
			}
			if ls.CreateTime == nil {
				return errors.New("listener create time is nil")
			}

			ids = append(ids, *ls.ListenerId)

			m := map[string]interface{}{
				"protocol":        "TCP",
				"id":              *ls.ListenerId,
				"name":            *ls.ListenerName,
				"port":            *ls.Port,
				"realserver_type": *ls.RealServerType,
				"status":          *ls.ListenerStatus,
				"scheduler":       *ls.Scheduler,
				"health_check":    *ls.HealthCheck == 1,
				"create_time":     *ls.CreateTime,
			}

			if ls.ConnectTimeout != nil {
				m["connect_timeout"] = *ls.ConnectTimeout
			}
			if ls.DelayLoop != nil {
				m["delay_loop"] = *ls.DelayLoop
			}

			listeners = append(listeners, m)
		}

	case "UDP":
		udpListeners, err := service.DescribeUDPListeners(ctx, proxyId, listenerId, name, port)
		if err != nil {
			return err
		}

		ids = make([]string, 0, len(udpListeners))
		listeners = make([]map[string]interface{}, 0, len(udpListeners))

		for _, ls := range udpListeners {
			if ls.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if ls.ListenerName == nil {
				return errors.New("listener name is nil")
			}
			if ls.Port == nil {
				return errors.New("listener port is nil")
			}
			if ls.RealServerType == nil {
				return errors.New("listener realserver type is nil")
			}
			if ls.ListenerStatus == nil {
				return errors.New("listener realserver type is nil")
			}
			if ls.Scheduler == nil {
				return errors.New("listener scheduler is nil")
			}
			if ls.CreateTime == nil {
				return errors.New("listener create time is nil")
			}

			ids = append(ids, *ls.ListenerId)

			m := map[string]interface{}{
				"protocol":        "UDP",
				"id":              *ls.ListenerId,
				"name":            *ls.ListenerName,
				"port":            *ls.Port,
				"realserver_type": *ls.RealServerType,
				"status":          *ls.ListenerStatus,
				"scheduler":       *ls.Scheduler,
				"create_time":     *ls.CreateTime,
			}

			listeners = append(listeners, m)
		}
	}

	d.Set("listeners", listeners)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
