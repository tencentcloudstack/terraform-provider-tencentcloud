package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func dataSourceTencentCloudGaapLayer7Listeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapLayer7ListenersRead,
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
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
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"forward_protocol": {
							Type:     schema.TypeString,
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

func dataSourceTencentCloudGaapLayer7ListenersRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_layer7_listeners.read")()
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
	case "HTTP":
		httpListeners, err := service.DescribeHTTPListeners(ctx, &proxyId, listenerId, name, port)
		if err != nil {
			return err
		}

		ids = make([]string, 0, len(httpListeners))
		listeners = make([]map[string]interface{}, 0, len(httpListeners))

		for _, ls := range httpListeners {
			if ls.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if ls.ListenerName == nil {
				return errors.New("listener name is nil")
			}
			if ls.Port == nil {
				return errors.New("listener port is nil")
			}
			if ls.ListenerStatus == nil {
				return errors.New("listener status is nil")
			}
			if ls.CreateTime == nil {
				return errors.New("listener create time is nil")
			}

			ids = append(ids, *ls.ListenerId)

			listeners = append(listeners, map[string]interface{}{
				"protocol":    "HTTP",
				"id":          *ls.ListenerId,
				"name":        *ls.ListenerName,
				"port":        *ls.Port,
				"status":      *ls.ListenerStatus,
				"create_time": *ls.CreateTime,
			})
		}

	case "HTTPS":
		httpsListeners, err := service.DescribeHTTPSListeners(ctx, &proxyId, listenerId, name, port)
		if err != nil {
			return err
		}

		ids = make([]string, 0, len(httpsListeners))
		listeners = make([]map[string]interface{}, 0, len(httpsListeners))

		for _, ls := range httpsListeners {
			if ls.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if ls.ListenerName == nil {
				return errors.New("listener name is nil")
			}
			if ls.Port == nil {
				return errors.New("listener port is nil")
			}
			if ls.ListenerStatus == nil {
				return errors.New("listener status is nil")
			}
			if ls.CertificateId == nil {
				return errors.New("listener certificate id is nil")
			}
			if ls.AuthType == nil {
				return errors.New("listener auth type is nil")
			}
			if ls.ForwardProtocol == nil {
				return errors.New("listener forward protocol is nil")
			}
			if ls.CreateTime == nil {
				return errors.New("listener create time is nil")
			}

			ids = append(ids, *ls.ListenerId)

			m := map[string]interface{}{
				"protocol":         "HTTPS",
				"id":               *ls.ListenerId,
				"name":             *ls.ListenerName,
				"port":             *ls.Port,
				"status":           *ls.ListenerStatus,
				"certificate_id":   *ls.CertificateId,
				"auth_type":        *ls.AuthType,
				"forward_protocol": *ls.ForwardProtocol,
				"create_time":      *ls.CreateTime,
			}

			if ls.ClientCertificateId != nil {
				m["client_certificate_id"] = *ls.ClientCertificateId
			}

			listeners = append(listeners, m)
		}
	}

	d.Set("listeners", listeners)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
