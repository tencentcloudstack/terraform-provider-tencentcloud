/*
Provides a resource to create a forward domain of layer7 listener.

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

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}
```

Import

GAAP http domain can be imported using the id, e.g.

-> **NOTE:** The format of tencentcloud_gaap_http_domain id is `[listener-id]+[protocol]+[domain]`.

```
  $ terraform import tencentcloud_gaap_http_domain.foo listener-11112222+HTTP+www.qq.com
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapHttpDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapHttpDomainCreate,
		Read:   resourceTencentCloudGaapHttpDomainRead,
		Update: resourceTencentCloudGaapHttpDomainUpdate,
		Delete: resourceTencentCloudGaapHttpDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the layer7 listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Forward domain of the layer7 listener.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "ID of the server certificate, default value is `default`.",
			},
			"client_certificate_id": {
				Deprecated:    "It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"client_certificate_ids"},
				Description:   "ID of the client certificate, default value is `default`.",
			},
			"client_certificate_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				Set:           schema.HashString,
				ConflictsWith: []string{"client_certificate_id"},
				Description:   "ID list of the poly client certificate.",
			},
			"basic_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether basic authentication is enable, default is `false`.",
			},
			"basic_auth_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the basic authentication.",
			},
			"realserver_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether realserver authentication is enable, default is `false`.",
			},
			"realserver_certificate_id": {
				Deprecated:    "It has been deprecated from version 1.28.0. Set `client_certificate_ids` instead.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"realserver_certificate_ids"},
				Description:   "CA certificate ID of the realserver.",
			},
			"realserver_certificate_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"realserver_certificate_id"},
				Description:   "CA certificate ID list of the realserver.",
			},
			"realserver_certificate_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CA certificate domain of the realserver.",
			},
			"gaap_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether SSL certificate authentication is enable, default is `false`.",
			},
			"gaap_auth_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the SSL certificate.",
			},
		},
	}
}

func resourceTencentCloudGaapHttpDomainCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_domain.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		protocol        string
		forwardProtocol string
	)
	httpListeners, err := service.DescribeHTTPListeners(ctx, nil, &listenerId, nil, nil)
	if err != nil {
		return err
	}
	if len(httpListeners) > 0 {
		protocol = "HTTP"
	}

	httpsListeners, err := service.DescribeHTTPSListeners(ctx, nil, &listenerId, nil, nil)
	if err != nil {
		return err
	}
	if len(httpsListeners) > 0 {
		protocol = "HTTPS"
		if httpsListeners[0].ForwardProtocol == nil {
			return errors.New("https listener forward protocol is nil")
		}
		forwardProtocol = *httpsListeners[0].ForwardProtocol
	}

	switch protocol {
	case "HTTP":
		if err := service.CreateHTTPDomain(ctx, listenerId, domain); err != nil {
			return err
		}

		id := fmt.Sprintf("%s+%s+%s", listenerId, protocol, domain)
		d.SetId(id)

	case "HTTPS":
		var (
			basicAuthId                 *string
			realserverCertificateDomain *string
			gaapCertificateId           *string
			realserverCertificateIds    []string
			polyClientCertificateIds    []string
		)

		certificateId := d.Get("certificate_id").(string)

		if raw, ok := d.GetOk("client_certificate_ids"); ok {
			set := raw.(*schema.Set)
			polyClientCertificateIds = make([]string, 0, set.Len())
			for _, polyIdRaw := range set.List() {
				polyId := polyIdRaw.(string)
				if polyId == "default" {
					return errors.New("client_certificate_ids can't have `default`")
				}

				polyClientCertificateIds = append(polyClientCertificateIds, polyId)
			}
		} else {
			if raw, ok := d.GetOk("client_certificate_id"); ok {
				ccId := raw.(string)
				if ccId != "default" {
					polyClientCertificateIds = append(polyClientCertificateIds, ccId)
				}
			}
		}

		// basic auth
		basicAuth := d.Get("basic_auth").(bool)

		if raw, ok := d.GetOk("basic_auth_id"); ok {
			basicAuthId = stringToPointer(raw.(string))
		}

		if basicAuth && basicAuthId == nil {
			return errors.New("when use basic auth, basic auth id should be set")
		}

		// realserver certification
		realserverAuth := d.Get("realserver_auth").(bool)
		if forwardProtocol == "HTTP" && realserverAuth {
			return errors.New("when listener forward protocol is http, realserver_auth can't be true")
		}

		if raw, ok := d.GetOk("realserver_certificate_ids"); ok {
			set := raw.(*schema.Set)
			realserverCertificateIds = make([]string, 0, set.Len())

			for _, id := range set.List() {
				realserverCertificateIds = append(realserverCertificateIds, id.(string))
			}
		} else if raw, ok := d.GetOk("realserver_certificate_id"); ok {
			realserverCertificateIds = []string{raw.(string)}
		}

		if raw, ok := d.GetOk("realserver_certificate_domain"); ok {
			realserverCertificateDomain = stringToPointer(raw.(string))
		}

		if realserverAuth && (len(realserverCertificateIds) == 0 || realserverCertificateDomain == nil) {
			return errors.New("when use realserver auth, realserver_certificate_ids and domain should be set")
		}

		// gaap certification
		gaapAuth := d.Get("gaap_auth").(bool)
		if forwardProtocol == "HTTP" && gaapAuth {
			return errors.New("when listener forward protocol is http, gaap_auth can't be set")
		}

		if raw, ok := d.GetOk("gaap_auth_id"); ok {
			gaapCertificateId = stringToPointer(raw.(string))
		}

		if gaapAuth && gaapCertificateId == nil {
			return errors.New("when use gaap auth, gaap auth id should be set")
		}

		if err := service.CreateHTTPSDomain(ctx, listenerId, domain, certificateId, polyClientCertificateIds); err != nil {
			return err
		}

		id := fmt.Sprintf("%s+%s+%s", listenerId, protocol, domain)
		d.SetId(id)

		if err := service.SetAdvancedAuth(
			ctx,
			listenerId, domain,
			realserverAuth, basicAuth, gaapAuth,
			realserverCertificateIds,
			realserverCertificateDomain, basicAuthId, gaapCertificateId,
		); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapHttpDomainRead(d, m)
}

func resourceTencentCloudGaapHttpDomainRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_domain.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	var (
		listenerId string
		domain     string
	)
	split := strings.Split(id, "+")

	if len(split) != 3 {
		log.Printf("[CRITAL]%s id is broken", logId)
		d.SetId("")
		return nil
	}

	listenerId, _, domain = split[0], split[1], split[2]

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	httpDomain, err := service.DescribeDomain(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	if httpDomain == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", httpDomain.Domain)
	_ = d.Set("listener_id", listenerId)

	if httpDomain.CertificateId == nil {
		httpDomain.CertificateId = stringToPointer("default")
	}
	_ = d.Set("certificate_id", httpDomain.CertificateId)

	clientCertificateIds := make([]*string, 0, len(httpDomain.PolyClientCertificateAliasInfo))
	for _, info := range httpDomain.PolyClientCertificateAliasInfo {
		clientCertificateIds = append(clientCertificateIds, info.CertificateId)
	}

	_ = d.Set("client_certificate_id", clientCertificateIds[0])
	_ = d.Set("client_certificate_ids", clientCertificateIds)

	if httpDomain.BasicAuth == nil {
		httpDomain.BasicAuth = int64ToPointer(0)
	}
	_ = d.Set("basic_auth", *httpDomain.BasicAuth == 1)
	_ = d.Set("basic_auth_id", httpDomain.BasicAuthConfId)

	if httpDomain.RealServerAuth == nil {
		httpDomain.RealServerAuth = int64ToPointer(0)
	}
	_ = d.Set("realserver_auth", *httpDomain.RealServerAuth == 1)

	realserverCertificateIds := make([]*string, 0, len(httpDomain.PolyRealServerCertificateAliasInfo))
	for _, info := range httpDomain.PolyRealServerCertificateAliasInfo {
		realserverCertificateIds = append(realserverCertificateIds, info.CertificateId)
	}

	_ = d.Set("realserver_certificate_ids", realserverCertificateIds)

	var realserverCertificateId *string
	if len(realserverCertificateIds) > 0 {
		realserverCertificateId = realserverCertificateIds[0]
	}
	_ = d.Set("realserver_certificate_id", realserverCertificateId)

	_ = d.Set("realserver_certificate_domain", httpDomain.RealServerCertificateDomain)

	if httpDomain.GaapAuth == nil {
		httpDomain.GaapAuth = int64ToPointer(0)
	}
	_ = d.Set("gaap_auth", *httpDomain.GaapAuth == 1)
	_ = d.Set("gaap_auth_id", httpDomain.GaapCertificateId)

	return nil
}

func resourceTencentCloudGaapHttpDomainUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_domain.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	var (
		listenerId string
		protocol   string
		domain     string
	)
	split := strings.Split(id, "+")

	if len(split) != 3 {
		log.Printf("[CRITAL]%s id is broken", logId)
		return resourceTencentCloudGaapHttpDomainRead(d, m)
	}

	listenerId, protocol, domain = split[0], split[1], split[2]

	switch protocol {
	default:
		log.Printf("[CRITAL]%s id is broken, protocol %s is invalid", logId, protocol)
		return resourceTencentCloudGaapHttpDomainRead(d, m)

	case "HTTP":
		// when protocol is http, nothing can be updated
		return errors.New("http listener can't set auth")

	case "HTTPS":
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	listeners, err := service.DescribeHTTPSListeners(ctx, nil, &listenerId, nil, nil)
	if err != nil {
		return err
	}

	// if listener doesn't exist, domain won't exist
	if len(listeners) == 0 {
		log.Printf("[DEBUG]%s domain listener doesn't exist", logId)
		d.SetId("")
		return nil
	}

	if listeners[0].ForwardProtocol == nil {
		return errors.New("listener forward protocol is nil")
	}
	forwardProtocol := *listeners[0].ForwardProtocol

	d.Partial(true)

	if d.HasChange("certificate_id") || d.HasChange("client_certificate_id") || d.HasChange("client_certificate_ids") {
		certificateId := d.Get("certificate_id").(string)

		var polyClientCertificateIds []string

		if d.HasChange("client_certificate_id") {
			if raw, ok := d.GetOk("client_certificate_id"); ok {
				if ccId := raw.(string); ccId != "default" {
					polyClientCertificateIds = []string{ccId}
				}
			}
		}

		if d.HasChange("client_certificate_ids") {
			if raw, ok := d.GetOk("client_certificate_ids"); ok {
				set := raw.(*schema.Set)
				polyClientCertificateIds = make([]string, 0, set.Len())
				for _, polyIdRaw := range set.List() {
					if polyId := polyIdRaw.(string); polyId != "default" {
						polyClientCertificateIds = append(polyClientCertificateIds, polyId)
					}
				}
			}
		}

		if err := service.ModifyDomainCertificate(
			ctx, listenerId, domain, certificateId,
			polyClientCertificateIds,
		); err != nil {
			return err
		}

		if d.HasChange("certificate_id") {
			d.SetPartial("certificate_id")
		}
		if d.HasChange("client_certificate_id") {
			d.SetPartial("client_certificate_id")
		}
		if d.HasChange("client_certificate_ids") {
			d.SetPartial("client_certificate_ids")
		}
	}

	if forwardProtocol == "HTTP" {
		if d.HasChange("realserver_auth") ||
			d.HasChange("realserver_certificate_id") ||
			d.HasChange("realserver_certificate_ids") ||
			d.HasChange("realserver_certificate_domain") ||
			d.HasChange("gaap_auth") ||
			d.HasChange("gaap_auth_id") {
			return errors.New("when listener forward protocol is HTTP, only can set basic auth")
		}
	}

	basicAuth := d.Get("basic_auth").(bool)
	realserverAuth := d.Get("realserver_auth").(bool)
	gaapAuth := d.Get("gaap_auth").(bool)

	var (
		updateAdvancedAttr []string

		basicAuthId                 *string
		realserverCertificateDomain *string
		gaapCertificateId           *string
		realserverCertificateIds    []string
	)

	if d.HasChange("basic_auth") {
		updateAdvancedAttr = append(updateAdvancedAttr, "basic_auth")
	}

	if d.HasChange("realserver_auth") {
		updateAdvancedAttr = append(updateAdvancedAttr, "realserver_auth")
	}

	if d.HasChange("gaap_auth") {
		updateAdvancedAttr = append(updateAdvancedAttr, "gaap_auth")
	}

	if raw, ok := d.GetOk("basic_auth_id"); ok {
		basicAuthId = stringToPointer(raw.(string))
	}

	if d.HasChange("basic_auth_id") {
		updateAdvancedAttr = append(updateAdvancedAttr, "basic_auth_id")
	}

	// only happen on modify basic auth
	if basicAuth && basicAuthId == nil {
		return errors.New("when enable basic auth, basic_auth_id must be set")
	}

	if raw, ok := d.GetOk("realserver_certificate_id"); ok {
		realserverCertificateIds = []string{raw.(string)}
	}

	if raw, ok := d.GetOk("realserver_certificate_ids"); ok {
		set := raw.(*schema.Set)
		realserverCertificateIds = make([]string, 0, set.Len())
		for _, id := range set.List() {
			realserverCertificateIds = append(realserverCertificateIds, id.(string))
		}
	}

	if d.HasChange("realserver_certificate_id") {
		updateAdvancedAttr = append(updateAdvancedAttr, "realserver_certificate_id")
	}

	if d.HasChange("realserver_certificate_ids") {
		updateAdvancedAttr = append(updateAdvancedAttr, "realserver_certificate_ids")
	}

	if raw, ok := d.GetOk("realserver_certificate_domain"); ok {
		realserverCertificateDomain = stringToPointer(raw.(string))
	}

	if d.HasChange("realserver_certificate_domain") {
		updateAdvancedAttr = append(updateAdvancedAttr, "realserver_certificate_domain")
	}

	// only happen on modify realserver auth
	if realserverAuth {
		if len(realserverCertificateIds) == 0 {
			return errors.New("when enable realserver auth, realserver_certificate_ids must be set")
		}

		if realserverCertificateDomain == nil {
			return errors.New("when enable realserver auth, realserver_certificate_domain must be set")
		}
	}

	if raw, ok := d.GetOk("gaap_auth_id"); ok {
		gaapCertificateId = stringToPointer(raw.(string))
	}

	if d.HasChange("gaap_auth_id") {
		updateAdvancedAttr = append(updateAdvancedAttr, "gaap_auth_id")
	}

	// only happen on modify gaap auth
	if gaapAuth && gaapCertificateId == nil {
		return errors.New("when enable gaap auth, gaap_auth_id must be set")
	}

	if len(updateAdvancedAttr) > 0 {
		if err := service.SetAdvancedAuth(
			ctx,
			listenerId, domain, realserverAuth, basicAuth, gaapAuth,
			realserverCertificateIds, realserverCertificateDomain, basicAuthId, gaapCertificateId,
		); err != nil {
			return err
		}

		for _, attr := range updateAdvancedAttr {
			d.SetPartial(attr)
		}
	}

	d.Partial(false)

	return resourceTencentCloudGaapHttpDomainRead(d, m)
}

func resourceTencentCloudGaapHttpDomainDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_domain.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	var (
		listenerId string
		domain     string
	)
	split := strings.Split(id, "+")

	if len(split) != 3 {
		log.Printf("[CRITAL]%s id is broken", logId)
		return nil
	}

	listenerId, domain = split[0], split[2]

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteDomain(ctx, listenerId, domain)
}
