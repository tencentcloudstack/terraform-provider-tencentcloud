/*
Provides a resource to create a layer7 listener of GAAP.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
```

Import

GAAP layer7 listener can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_layer7_listener.foo listener-11112222
```
*/
package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
)

func resourceTencentCloudGaapLayer7Listener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapLayer7ListenerCreate,
		Read:   resourceTencentCloudGaapLayer7ListenerRead,
		Update: resourceTencentCloudGaapLayer7ListenerUpdate,
		Delete: resourceTencentCloudGaapLayer7ListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				ForceNew:     true,
				Description:  "Protocol of the layer7 listener, and the available values include `HTTP` and `HTTPS`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the layer7 listener, the maximum length is 30.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePort,
				ForceNew:     true,
				Description:  "Port of the layer7 listener.",
			},
			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the GAAP proxy.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Certificate ID of the layer7 listener. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"forward_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				ForceNew:     true,
				Description:  "Protocol type of the forwarding, the available values include `HTTP` and `HTTPS`. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"auth_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				ForceNew:     true,
				Description:  "Authentication type of the layer7 listener. `0` is one-way authentication and `1` is mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"client_certificate_id": {
				Deprecated:    "It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"client_certificate_ids"},
				Description:   "ID of the client certificate. Set only when `auth_type` is specified as mutual authentication.  NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"client_certificate_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"client_certificate_id"},
				Description:   "ID list of the client certificate.  Set only when `auth_type` is specified as mutual authentication.  NOTES: Only supports listeners of `HTTPS` protocol.",
			},

			// computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the layer7 listener.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the layer7 listener.",
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
			forwardProtocol string
			authType        int
		)

		certificateId := d.Get("certificate_id").(string)
		if certificateId == "" {
			return errors.New("when protocol is HTTPS, certificate_id can't be empty")
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

		var polyClientCertificateIds []string

		if raw, ok := d.GetOk("client_certificate_id"); ok {
			polyClientCertificateIds = append(polyClientCertificateIds, raw.(string))
		}
		if raw, ok := d.GetOk("client_certificate_ids"); ok {
			set := raw.(*schema.Set)
			polyClientCertificateIds = make([]string, 0, set.Len())
			for _, polyId := range set.List() {
				polyClientCertificateIds = append(polyClientCertificateIds, polyId.(string))
			}
		}

		if authType == 1 && len(polyClientCertificateIds) == 0 {
			return errors.New("when protocol is HTTPS and auth type is 1, client_certificate_ids can't be empty")
		}

		id, err = service.CreateHTTPSListener(
			ctx,
			name, certificateId, forwardProtocol, proxyId, polyClientCertificateIds, port, authType,
		)
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
	protocol := d.Get("protocol").(string)
	var (
		name                     *string
		port                     *uint64
		certificateId            *string
		forwardProtocol          *string
		authType                 *int64
		clientCertificateId      *string
		status                   *uint64
		createTime               string
		polyClientCertificateIds []*string
	)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

LOOP:
	for {
		switch protocol {
		case "":
			// import mode, need check protocol
			httpListeners, err := service.DescribeHTTPListeners(ctx, nil, &id, nil, nil)
			if err != nil {
				return err
			}
			if len(httpListeners) > 0 {
				protocol = "HTTP"
				continue
			}

			httpsListeners, err := service.DescribeHTTPSListeners(ctx, nil, &id, nil, nil)
			if err != nil {
				return err
			}
			if len(httpsListeners) > 0 {
				protocol = "HTTPS"
				continue
			}

			// layer7 listener is not found
			d.SetId("")
			return nil

		case "HTTP":
			listeners, err := service.DescribeHTTPListeners(ctx, nil, &id, nil, nil)
			if err != nil {
				return err
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

			name = listener.ListenerName
			port = listener.Port
			status = listener.ListenerStatus

			if listener.CreateTime == nil {
				return errors.New("listener create time is nil")
			}
			createTime = formatUnixTime(*listener.CreateTime)

			break LOOP

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(ctx, nil, &id, nil, nil)
			if err != nil {
				return err
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

			name = listener.ListenerName
			port = listener.Port
			certificateId = listener.CertificateId
			forwardProtocol = listener.ForwardProtocol
			authType = listener.AuthType

			// mutual authentication
			if *authType == 1 {
				clientCertificateId = listener.PolyClientCertificateAliasInfo[0].CertificateId
				polyClientCertificateIds = make([]*string, 0, len(listener.PolyClientCertificateAliasInfo))
				for _, polyCc := range listener.PolyClientCertificateAliasInfo {
					polyClientCertificateIds = append(polyClientCertificateIds, polyCc.CertificateId)
				}
			}

			status = listener.ListenerStatus

			if listener.CreateTime == nil {
				return errors.New("listener create time is nil")
			}
			createTime = formatUnixTime(*listener.CreateTime)

			break LOOP
		}
	}

	d.Set("protocol", protocol)
	d.Set("name", name)
	d.Set("port", port)
	d.Set("certificate_id", certificateId)
	d.Set("forward_protocol", forwardProtocol)
	d.Set("auth_type", authType)
	d.Set("client_certificate_id", clientCertificateId)
	d.Set("client_certificate_ids", polyClientCertificateIds)
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
			name                     *string
			certificateId            *string
			forwardProtocol          *string
			polyClientCertificateIds []string
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
			if raw, ok := d.GetOk("client_certificate_id"); ok {
				polyClientCertificateIds = append(polyClientCertificateIds, raw.(string))
			}
		}

		if d.HasChange("client_certificate_ids") {
			if raw, ok := d.GetOk("client_certificate_ids"); ok {
				set := raw.(*schema.Set)
				polyClientCertificateIds = make([]string, 0, set.Len())

				for _, polyId := range set.List() {
					polyClientCertificateIds = append(polyClientCertificateIds, polyId.(string))
				}
			}
		}

		if err := service.ModifyHTTPSListener(ctx, proxyId, id, name, forwardProtocol, certificateId, polyClientCertificateIds); err != nil {
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
