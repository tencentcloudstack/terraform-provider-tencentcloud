package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
)

func resourceTencentCloudGaapLayer7Listener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapLayer7ListenerCreate,
		Read:   resourceTencentCloudGaapLayer7ListenerRead,
		Update: resourceTencentCloudGaapLayer7ListenerUpdate,
		Delete: resourceTencentCloudGaapLayer7ListenerDelete,
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePort,
				ForceNew:     true,
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
			},
			"auth_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				ForceNew:     true,
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// computed
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudGaapLayer7ListenerCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_layer7_listener.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)
	port := d.Get("port").(int)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		id  string
		err error
	)

	switch protocol {
	case "HTTP":
		id, err = service.CreateHTTPListener(ctx, name, proxyId, port)

	case "HTTPS":
		var (
			certificateId       string
			forwardProtocol     string
			authType            int
			clientCertificateId *string
		)

		if raw, ok := d.GetOk("certificate_id"); ok {
			certificateId = raw.(string)
		} else {
			return errors.New("when protocol is HTTPS, certificate_id is required")
		}

		if raw, ok := d.GetOk("forward_protocol"); ok {
			forwardProtocol = raw.(string)
		} else {
			return errors.New("when protocol is HTTPS, forward_protocol is required")
		}

		if raw, ok := d.GetOkExists("auth_type"); ok {
			authType = raw.(int)
		} else {
			return errors.New("when protocol is HTTPS, auth_type is required")
		}

		if raw, ok := d.GetOk("client_certificate_id"); ok {
			clientCertificateId = stringToPointer(raw.(string))
		}

		id, err = service.CreateHTTPSListener(ctx, name, certificateId, forwardProtocol, proxyId, port, authType, clientCertificateId)
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapLayer7ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer7ListenerRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_layer7_listener.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	proxyId := d.Get("proxy_id").(string)
	protocol := d.Get("protocol").(string)
	var (
		name                string
		port                int
		certificateId       *string
		forwardProtocol     *string
		authType            *int
		clientCertificateId *string
		status              int
		createTime          int
	)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	switch protocol {
	case "HTTP":
		listeners, err := service.DescribeHTTPListeners(ctx, &proxyId, &id, nil, nil)
		if err != nil {
			return err
		}

		if len(listeners) == 0 {
			d.SetId("")
			return nil
		}

		var listener *gaap.HTTPListener
		for _, l := range listeners {
			if l.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if *l.ListenerId == id {
				listener = l
				break
			}
		}

		if listener == nil {
			d.SetId("")
			return nil
		}

		if listener.ListenerName == nil {
			return errors.New("listener name is nil")
		}
		name = *listener.ListenerName

		if listener.Port == nil {
			return errors.New("listener port is nil")
		}
		port = int(*listener.Port)

		if listener.ListenerStatus == nil {
			return errors.New("listener status is nil")
		}
		status = int(*listener.ListenerStatus)

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = int(*listener.CreateTime)

	case "HTTPS":
		listeners, err := service.DescribeHTTPSListeners(ctx, &proxyId, &id, nil, nil)
		if err != nil {
			return err
		}

		if len(listeners) == 0 {
			d.SetId("")
			return nil
		}

		var listener *gaap.HTTPSListener
		for _, l := range listeners {
			if l.ListenerId == nil {
				return errors.New("listener id is nil")
			}
			if *l.ListenerId == id {
				listener = l
				break
			}
		}

		if listener == nil {
			d.SetId("")
			return nil
		}

		if listener.ListenerName == nil {
			return errors.New("listener name is nil")
		}
		name = *listener.ListenerName

		if listener.Port == nil {
			return errors.New("listener port is nil")
		}
		port = int(*listener.Port)

		if listener.CertificateId == nil {
			return errors.New("listener certificate id is nil")
		}
		certificateId = listener.CertificateId

		if listener.ForwardProtocol == nil {
			return errors.New("listener forward protocol is nil")
		}
		forwardProtocol = listener.ForwardProtocol

		if listener.AuthType == nil {
			return errors.New("listener auth type is nil")
		}
		authType = common.IntPtr(int(*listener.AuthType))

		clientCertificateId = listener.ClientCertificateId

		if listener.ListenerStatus == nil {
			return errors.New("listener status is nil")
		}
		status = int(*listener.ListenerStatus)

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = int(*listener.CreateTime)
	}

	d.Set("name", name)
	d.Set("port", port)
	if certificateId != nil {
		d.Set("certificate_id", certificateId)
	}
	if forwardProtocol != nil {
		d.Set("forward_protocol", forwardProtocol)
	}
	if authType != nil {
		d.Set("auth_type", authType)
	}
	if clientCertificateId != nil {
		d.Set("client_certificate_id", clientCertificateId)
	}
	d.Set("status", status)
	d.Set("create_time", createTime)

	return nil
}

func resourceTencentCloudGaapLayer7ListenerUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_layer7_listener.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	switch protocol {
	case "HTTP":
		if d.HasChange("name") {
			name := d.Get("name").(string)
			if err := service.ModifyHTTPListener(ctx, id, proxyId, name); err != nil {
				return err
			}
		}

	case "HTTPS":
		var (
			name                *string
			certificateId       *string
			forwardProtocol     *string
			clientCertificateId *string
		)

		if d.HasChange("name") {
			name = stringToPointer(d.Get("name").(string))
		}
		if d.HasChange("certificate_id") {
			certificateId = stringToPointer(d.Get("certificate_id").(string))
		}
		if d.HasChange("forward_protocol") {
			forwardProtocol = stringToPointer(d.Get("forward_protocol").(string))
		}
		if d.HasChange("client_certificate_id") {
			clientCertificateId = stringToPointer(d.Get("client_certificate_id").(string))
		}

		if err := service.ModifyHTTPSListener(ctx, proxyId, id, name, forwardProtocol, certificateId, clientCertificateId); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapLayer7ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer7ListenerDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_layer7_listener.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	proxyId := d.Get("proxy_id").(string)
	protocol := d.Get("protocol").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteLayer7Listener(ctx, id, proxyId, protocol)
}
