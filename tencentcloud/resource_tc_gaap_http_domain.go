package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func resourceTencentCloudGaapHttpDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapHttpDomainCreate,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realserver_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"realserver_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realserver_certificate_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"basic_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"basic_auth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gaap_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gaap_auth_id": {
				Type:     schema.TypeString,
				Optional: true,
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

	var (
		certificateId               *string
		clientCertificateId         *string
		realserverAuth              *bool
		realserverCertificateId     *string
		realserverCertificateDomain *string
		basicAuth                   *bool
		basicAuthId                 *string
		gaapAuth                    *bool
		gaapCertificateId           *string
		advanceAuth                 bool
	)

	if raw, ok := d.GetOk("certificate_id"); ok {
		certificateId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("client_certificate_id"); ok {
		clientCertificateId = stringToPointer(raw.(string))
	}

	if raw, ok := d.GetOkExists("realserver_auth"); ok {
		advanceAuth = true
		realserverAuth = common.BoolPtr(raw.(bool))
	}
	if raw, ok := d.GetOk("realserver_certificate_id"); ok {
		advanceAuth = true
		realserverCertificateId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("realserver_certificate_domain"); ok {
		advanceAuth = true
		realserverCertificateDomain = stringToPointer(raw.(string))
	}
	if realserverAuth != nil && *realserverAuth {
		if realserverCertificateId == nil {
			return errors.New("when use realserver auth, realserver auth id should be set")
		}
	}

	if raw, ok := d.GetOkExists("basic_auth"); ok {
		advanceAuth = true
		basicAuth = common.BoolPtr(raw.(bool))
	}
	if raw, ok := d.GetOk("basic_auth_id"); ok {
		advanceAuth = true
		basicAuthId = stringToPointer(raw.(string))
	}
	if basicAuth != nil && *basicAuth {
		if basicAuthId == nil {
			return errors.New("when use basic auth, basic auth id should be set")
		}
	}

	if raw, ok := d.GetOkExists("gaap_auth"); ok {
		advanceAuth = true
		gaapAuth = common.BoolPtr(raw.(bool))
	}
	if raw, ok := d.GetOk("gaap_auth_id"); ok {
		advanceAuth = true
		gaapCertificateId = stringToPointer(raw.(string))
	}
	if gaapAuth != nil && *gaapAuth {
		if gaapCertificateId == nil {
			return errors.New("when use gaap auth, gaap auth id should be set")
		}
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateHttpDomain(ctx, listenerId, domain, certificateId, clientCertificateId)
	if err != nil {
		return err
	}
	// set id early to create resource so that can destroy it if set advanced auth failed
	d.SetId(id)

	if advanceAuth {
		if err := service.SetAdvancedAuth(
			ctx,
			listenerId, domain,
			realserverAuth, basicAuth, gaapAuth,
			realserverCertificateId, realserverCertificateDomain, basicAuthId, gaapCertificateId,
		); err != nil {
			return err
		}
	}
}
