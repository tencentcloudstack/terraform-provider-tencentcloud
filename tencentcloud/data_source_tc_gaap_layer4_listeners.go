/*
Use this data source to query gaap layer4 listeners.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = tencentcloud_gaap_proxy.foo.id
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}

data "tencentcloud_gaap_layer4_listeners" "foo" {
  protocol    = "TCP"
  proxy_id    = tencentcloud_gaap_proxy.foo.id
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapLayer4Listeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapLayer4ListenersRead,
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP"}),
				Description:  "Protocol of the layer4 listener to be queried. Valid values: `TCP` and `UDP`.",
			},
			"proxy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the GAAP proxy to be queried.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the layer4 listener to be queried.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the layer4 listener to be queried.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  "Port of the layer4 listener to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of layer4 listeners. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the layer4 listener.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the layer4 listener.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the layer4 listener.",
						},
						"proxy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the GAAP proxy.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the layer4 listener.",
						},
						"realserver_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the realserver.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the layer4 listener.",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling policy of the layer4 listener.",
						},
						"health_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enable.",
						},
						"connect_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout of the health check response.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval of the health check.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the layer4 listener.",
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	protocol := d.Get("protocol").(string)

	var (
		proxyId    *string
		listenerId *string
		name       *string
		port       *int
		ids        []string
		listeners  []map[string]interface{}
	)

	if raw, ok := d.GetOk("proxy_id"); ok {
		proxyId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("listener_id"); ok {
		listenerId = helper.String(raw.(string))
	}

	if proxyId == nil && listenerId == nil {
		return errors.New("proxy_id or listener_id must be set")
	}

	if raw, ok := d.GetOk("listener_name"); ok {
		name = helper.String(raw.(string))
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
			if ls.HealthCheck == nil {
				ls.HealthCheck = helper.IntUint64(0)
			}

			ids = append(ids, *ls.ListenerId)

			m := map[string]interface{}{
				"protocol":        "TCP",
				"id":              ls.ListenerId,
				"name":            ls.ListenerName,
				"proxy_id":        ls.ProxyId,
				"port":            ls.Port,
				"realserver_type": ls.RealServerType,
				"status":          ls.ListenerStatus,
				"scheduler":       ls.Scheduler,
				"health_check":    *ls.HealthCheck == 1,
				"create_time":     helper.FormatUnixTime(*ls.CreateTime),
				"connect_timeout": ls.ConnectTimeout,
				"interval":        ls.DelayLoop,
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
			ids = append(ids, *ls.ListenerId)

			m := map[string]interface{}{
				"protocol":        "UDP",
				"id":              ls.ListenerId,
				"name":            ls.ListenerName,
				"proxy_id":        ls.ProxyId,
				"port":            ls.Port,
				"realserver_type": ls.RealServerType,
				"status":          ls.ListenerStatus,
				"scheduler":       ls.Scheduler,
				"create_time":     helper.FormatUnixTime(*ls.CreateTime),
			}

			listeners = append(listeners, m)
		}
	}

	_ = d.Set("listeners", listeners)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), listeners); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
