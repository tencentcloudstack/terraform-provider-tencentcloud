package gaap

import (
	"context"
	"errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapLayer4Listener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapLayer4ListenerCreate,
		Read:   resourceTencentCloudGaapLayer4ListenerRead,
		Update: resourceTencentCloudGaapLayer4ListenerUpdate,
		Delete: resourceTencentCloudGaapLayer4ListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"TCP", "UDP"}),
				ForceNew:     true,
				Description:  "Protocol of the layer4 listener. Valid value: `TCP` and `UDP`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "Name of the layer4 listener, the maximum length is 30.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				ForceNew:     true,
				Description:  "Port of the layer4 listener.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rr",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"rr", "wrr", "lc"}),
				Description:  "Scheduling policy of the layer4 listener, default value is `rr`. Valid value: `rr`, `wrr` and `lc`.",
			},
			"realserver_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"IP", "DOMAIN"}),
				ForceNew:     true,
				Description:  "Type of the realserver. Valid value: `IP` and `DOMAIN`. NOTES: when the `protocol` is specified as `TCP` and the `scheduler` is specified as `wrr`, the item can only be set to `IP`.",
			},
			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the GAAP proxy.",
			},
			"health_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether health check is enable, default value is `false`. NOTES: Only supports listeners of `TCP` protocol.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: tccommon.ValidateIntegerInRange(5, 300),
				Description:  "Interval of the health check, default value is 5s. NOTES: Only supports listeners of `TCP` protocol.",
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description:  "Timeout of the health check response, should less than interval, default value is 2s. NOTES: Only supports listeners of `TCP` protocol and require less than `interval`.",
			},
			"client_ip_method": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "The way the listener gets the client IP, 0 for TOA, 1 for Proxy Protocol, default value is 0. NOTES: Only supports listeners of `TCP` protocol.",
			},
			"realserver_bind_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return helper.HashString(fmt.Sprintf("%s-%s-%d-%d", m["id"].(string), m["ip"].(string), m["port"].(int), m["weight"].(int)))
				},
				Description: "An information list of GAAP realserver.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the GAAP realserver.",
						},
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP of the GAAP realserver.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "Port of the GAAP realserver.",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: tccommon.ValidateIntegerInRange(1, 100),
							Description:  "Scheduling weight, default value is `1`. The range of values is [1,100].",
						},
					},
				},
			},

			// computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the layer4 listener.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the layer4 listener.",
			},
		},
	}
}

func resourceTencentCloudGaapLayer4ListenerCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)
	port := d.Get("port").(int)
	scheduler := d.Get("scheduler").(string)
	realserverType := d.Get("realserver_type").(string)

	if protocol == "TCP" && realserverType == "DOMAIN" && scheduler == "wrr" {
		return errors.New("TCP listener DOMAIN realserver type doesn't support wrr scheduler")
	}

	proxyId := d.Get("proxy_id").(string)
	healthCheck := d.Get("health_check").(bool)

	if protocol == "UDP" && healthCheck {
		return errors.New("UDP listener can't use health check")
	}

	interval := d.Get("interval").(int)
	connectTimeout := d.Get("connect_timeout").(int)

	// only check for TCP listener
	if protocol == "TCP" && connectTimeout >= interval {
		return errors.New("connect_timeout must be less than interval")
	}
	clientIPMethod := d.Get("client_ip_method").(int)

	var realservers []gaapRealserverBind
	if raw, ok := d.GetOk("realserver_bind_set"); ok {
		list := raw.(*schema.Set).List()
		realservers = make([]gaapRealserverBind, 0, len(list))
		for _, v := range list {
			m := v.(map[string]interface{})

			if scheduler == "rr" && m["weight"].(int) != 1 {
				return errors.New("when scheduler is rr, realserver weight should be 1 or null")
			}

			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}
	}

	var (
		id  string
		err error
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	switch protocol {
	case "TCP":
		id, err = service.CreateTCPListener(ctx, name, scheduler, realserverType, proxyId, port, interval, connectTimeout, clientIPMethod, healthCheck)
		if err != nil {
			return err
		}

	case "UDP":
		id, err = service.CreateUDPListener(ctx, name, scheduler, realserverType, proxyId, port)
		if err != nil {
			return err
		}
	}

	d.SetId(id)

	if len(realservers) > 0 {
		if err := service.BindLayer4ListenerRealservers(ctx, id, protocol, proxyId, realservers); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapLayer4ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer4ListenerRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	var (
		protocol       string
		name           *string
		port           *uint64
		scheduler      *string
		realServerType *string
		healthCheck    *bool
		interval       *uint64
		connectTimeout *uint64
		status         *uint64
		createTime     string
		realservers    []map[string]interface{}
		clientIpMethod *uint64
		proxyId        *string
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	tcpListeners, err := service.DescribeTCPListeners(ctx, nil, &id, nil, nil)
	if err != nil {
		return err
	}

	udpListeners, err := service.DescribeUDPListeners(ctx, nil, &id, nil, nil)
	if err != nil {
		return err
	}

	switch {
	case len(tcpListeners) > 0:
		protocol = "TCP"

		listener := tcpListeners[0]
		clientIpMethod = listener.ClientIPMethod
		name = listener.ListenerName
		port = listener.Port
		scheduler = listener.Scheduler
		realServerType = listener.RealServerType
		proxyId = listener.ProxyId

		if listener.HealthCheck == nil {
			return errors.New("listener health check is nil")
		}
		healthCheck = helper.Bool(*listener.HealthCheck == 1)

		interval = listener.DelayLoop
		connectTimeout = listener.ConnectTimeout

		if len(listener.RealServerSet) > 0 {
			realservers = make([]map[string]interface{}, 0, len(listener.RealServerSet))
			for _, rs := range listener.RealServerSet {
				realservers = append(realservers, map[string]interface{}{
					"id":     rs.RealServerId,
					"ip":     rs.RealServerIP,
					"port":   rs.RealServerPort,
					"weight": rs.RealServerWeight,
				})
			}
		}

		status = listener.ListenerStatus

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = helper.FormatUnixTime(*listener.CreateTime)

	case len(udpListeners) > 0:
		protocol = "UDP"

		listener := udpListeners[0]

		name = listener.ListenerName
		port = listener.Port
		scheduler = listener.Scheduler
		realServerType = listener.RealServerType
		proxyId = listener.ProxyId

		healthCheck = helper.Bool(false)
		connectTimeout = helper.IntUint64(2)
		interval = helper.IntUint64(5)
		clientIpMethod = helper.IntUint64(0)

		if len(listener.RealServerSet) > 0 {
			realservers = make([]map[string]interface{}, 0, len(listener.RealServerSet))
			for _, rs := range listener.RealServerSet {
				realservers = append(realservers, map[string]interface{}{
					"id":     rs.RealServerId,
					"ip":     rs.RealServerIP,
					"port":   rs.RealServerPort,
					"weight": rs.RealServerWeight,
				})
			}
		}

		status = listener.ListenerStatus

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = helper.FormatUnixTime(*listener.CreateTime)

	default:
		d.SetId("")
		return nil
	}

	_ = d.Set("protocol", protocol)
	_ = d.Set("name", name)
	_ = d.Set("port", port)
	_ = d.Set("scheduler", scheduler)
	_ = d.Set("realserver_type", realServerType)
	_ = d.Set("health_check", healthCheck)
	_ = d.Set("interval", interval)
	_ = d.Set("connect_timeout", connectTimeout)
	_ = d.Set("client_ip_method", clientIpMethod)
	_ = d.Set("realserver_bind_set", realservers)
	_ = d.Set("status", status)
	_ = d.Set("create_time", createTime)
	_ = d.Set("proxy_id", proxyId)

	return nil
}

func resourceTencentCloudGaapLayer4ListenerUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)
	var (
		name           *string
		scheduler      *string
		healthCheck    *bool
		interval       int
		connectTimeout int
		attrChange     []string
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("name") {
		attrChange = append(attrChange, "name")
		name = helper.String(d.Get("name").(string))
	}

	if d.HasChange("scheduler") {
		attrChange = append(attrChange, "scheduler")
		scheduler = helper.String(d.Get("scheduler").(string))
	}

	if d.HasChange("health_check") {
		attrChange = append(attrChange, "health_check")
		healthCheck = helper.Bool(d.Get("health_check").(bool))
	}
	if protocol == "UDP" && healthCheck != nil && *healthCheck {
		return errors.New("UDP listener can't enable health check")
	}

	if d.HasChange("interval") {
		attrChange = append(attrChange, "interval")
	}
	interval = d.Get("interval").(int)

	if d.HasChange("connect_timeout") {
		attrChange = append(attrChange, "connect_timeout")
	}
	connectTimeout = d.Get("connect_timeout").(int)

	// only check for TCP listener
	if protocol == "TCP" && connectTimeout >= interval {
		return errors.New("connect_timeout must be less than interval")
	}

	if len(attrChange) > 0 {
		switch protocol {
		case "TCP":
			if err := service.ModifyTCPListenerAttribute(ctx, proxyId, id, name, scheduler, healthCheck, interval, connectTimeout); err != nil {
				return err
			}

		case "UDP":
			if err := service.ModifyUDPListenerAttribute(ctx, proxyId, id, name, scheduler); err != nil {
				return err
			}
		}
	}

	if d.HasChange("realserver_bind_set") {
		list := d.Get("realserver_bind_set").(*schema.Set).List()
		realservers := make([]gaapRealserverBind, 0, len(list))
		for _, v := range list {
			m := v.(map[string]interface{})
			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}

		if err := service.BindLayer4ListenerRealservers(ctx, id, protocol, proxyId, realservers); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapLayer4ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer4ListenerDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteLayer4Listener(ctx, id, proxyId, protocol)
}
