package gaap

import (
	"context"
	"errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapLayer7Listener() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				ForceNew:     true,
				Description:  "Protocol of the layer7 listener. Valid value: `HTTP` and `HTTPS`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "Name of the layer7 listener, the maximum length is 30.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				ForceNew:     true,
				Description:  "Port of the layer7 listener.",
			},
			"proxy_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_id"},
				AtLeastOneOf:  []string{"group_id"},
				Description:   "ID of the GAAP proxy.",
			},
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"proxy_id"},
				AtLeastOneOf:  []string{"proxy_id"},
				Description:   "Group ID.",
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
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				ForceNew:     true,
				Description:  "Protocol type of the forwarding. Valid value: `HTTP` and `HTTPS`. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"auth_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				ForceNew:     true,
				Description:  "Authentication type of the layer7 listener. `0` is one-way authentication and `1` is mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"client_certificate_id": {
				Deprecated:    "It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"client_certificate_ids"},
				Description:   "ID of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"client_certificate_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"client_certificate_id"},
				Description:   "ID list of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.",
			},
			"tls_support_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "TLS version, optional TLSv1, TLSv1.1, TLSv1.2, TLSv1.3.",
			},
			"tls_ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Password Suite, optional GAAP_TLS_CIPHERS_STRICT, GAAP_TLS_CIPHERS_GENERAL, GAAP_TLS_CIPHERS_WIDE(default).",
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
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer7_listener.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)
	port := d.Get("port").(int)
	var proxyId, groupId string
	if v, ok := d.GetOk("proxy_id"); ok {
		proxyId = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id  string
		err error
	)

	switch protocol {
	case "HTTP":
		id, err = service.CreateHTTPListener(ctx, name, proxyId, groupId, port)

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
			name, certificateId, forwardProtocol, proxyId, groupId, polyClientCertificateIds, port, authType,
		)
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	vTlsCiphers, okTlsCiphers := d.GetOk("tls_ciphers")
	vTlsSupportVersions, okTlsSupportVersions := d.GetOk("tls_support_versions")
	if okTlsCiphers && okTlsSupportVersions {
		if protocol != "HTTPS" {
			return errors.New("Only https can set tls")
		}
		if proxyId != "" {
			proxyDetail, err := service.DescribeGaapProxyDetail(ctx, proxyId)
			if err != nil {
				return err
			}
			if proxyDetail.IsSupportTLSChoice != nil && int(*proxyDetail.IsSupportTLSChoice) != 1 {
				return fmt.Errorf("proxy(%s) not support TLS Choice", proxyId)
			}
		}
		if groupId != "" {
			proxyGroup, err := service.DescribeGaapProxyGroupById(ctx, groupId)
			if err != nil {
				return err
			}
			if proxyGroup.IsSupportTLSChoice != nil && int(*proxyGroup.IsSupportTLSChoice) != 1 {
				return fmt.Errorf("group(%s) not support TLS Choice", groupId)
			}
		}
		err := service.SetTlsVersion(ctx, id, vTlsCiphers.(string), helper.InterfacesStrings(vTlsSupportVersions.(*schema.Set).List()))
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapLayer7ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer7ListenerRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer7_listener.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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
		proxyId                  *string
		tlsCiphers               *string
		tlsSupportVersion        []*string
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

LOOP:
	for {
		switch protocol {
		case "":
			// import mode, need check protocol
			httpListeners, err := service.DescribeHTTPListeners(ctx, nil, nil, &id, nil, nil)
			if err != nil {
				return err
			}
			if len(httpListeners) > 0 {
				protocol = "HTTP"
				continue
			}

			httpsListeners, err := service.DescribeHTTPSListeners(ctx, nil, nil, &id, nil, nil)
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
			listeners, err := service.DescribeHTTPListeners(ctx, nil, nil, &id, nil, nil)
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
			proxyId = listener.ProxyId

			if listener.CreateTime == nil {
				return errors.New("listener create time is nil")
			}
			createTime = helper.FormatUnixTime(*listener.CreateTime)

			break LOOP

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(ctx, nil, nil, &id, nil, nil)
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
			proxyId = listener.ProxyId
			tlsCiphers = listener.TLSCiphers
			tlsSupportVersion = listener.TLSSupportVersion

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
			createTime = helper.FormatUnixTime(*listener.CreateTime)

			break LOOP
		}
	}

	_ = d.Set("protocol", protocol)
	_ = d.Set("name", name)
	_ = d.Set("port", port)
	_ = d.Set("certificate_id", certificateId)
	_ = d.Set("forward_protocol", forwardProtocol)
	_ = d.Set("auth_type", authType)
	_ = d.Set("client_certificate_id", clientCertificateId)
	_ = d.Set("client_certificate_ids", polyClientCertificateIds)
	_ = d.Set("status", status)
	_ = d.Set("create_time", createTime)
	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("tls_ciphers", tlsCiphers)
	_ = d.Set("tls_support_versions", tlsSupportVersion)

	return nil
}

func resourceTencentCloudGaapLayer7ListenerUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer7_listener.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	var proxyId, groupId string
	if v, ok := d.GetOk("proxy_id"); ok {
		proxyId = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	switch protocol {
	case "HTTP":
		if d.HasChange("name") {
			name := d.Get("name").(string)
			if err := service.ModifyHTTPListener(ctx, id, proxyId, groupId, name); err != nil {
				return err
			}

		}
		if d.HasChange("tls_support_versions") || d.HasChange("tls_ciphers") {
			return errors.New("http listener not support change tls_support_versions or tls_ciphers")
		}

	case "HTTPS":
		var (
			name                     *string
			certificateId            *string
			forwardProtocol          *string
			polyClientCertificateIds []string
			isModifyHTTPSListener    bool
		)

		name = helper.String(d.Get("name").(string))
		if d.HasChange("certificate_id") {
			certificateId = helper.String(d.Get("certificate_id").(string))
			isModifyHTTPSListener = true
		}
		forwardProtocol = helper.String(d.Get("forward_protocol").(string))
		if d.HasChange("client_certificate_id") {
			if raw, ok := d.GetOk("client_certificate_id"); ok {
				polyClientCertificateIds = append(polyClientCertificateIds, raw.(string))
			}
			isModifyHTTPSListener = true
		}

		if d.HasChange("client_certificate_ids") {
			if raw, ok := d.GetOk("client_certificate_ids"); ok {
				set := raw.(*schema.Set)
				polyClientCertificateIds = make([]string, 0, set.Len())

				for _, polyId := range set.List() {
					polyClientCertificateIds = append(polyClientCertificateIds, polyId.(string))
				}
			}
			isModifyHTTPSListener = true
		}

		if isModifyHTTPSListener {
			if err := service.ModifyHTTPSListener(ctx, proxyId, groupId, id, name, forwardProtocol, certificateId, polyClientCertificateIds); err != nil {
				return err
			}
		}

		if d.HasChange("tls_support_versions") || d.HasChange("tls_ciphers") {
			var (
				tlsCiphers        string
				tlsSupportVersion []string
			)
			if v, ok := d.GetOk("tls_ciphers"); ok {
				tlsCiphers = v.(string)
			}
			if v, ok := d.GetOk("tls_support_versions"); ok {
				tlsSupportVersion = helper.InterfacesStrings(v.(*schema.Set).List())
			}
			if proxyId != "" {
				proxyDetail, err := service.DescribeGaapProxyDetail(ctx, proxyId)
				if err != nil {
					return err
				}
				if proxyDetail.IsSupportTLSChoice != nil && int(*proxyDetail.IsSupportTLSChoice) != 1 {
					return fmt.Errorf("proxy(%s) not support TLS Choice", proxyId)
				}
			}
			if groupId != "" {
				proxyGroup, err := service.DescribeGaapProxyGroupById(ctx, groupId)
				if err != nil {
					return err
				}
				if proxyGroup.IsSupportTLSChoice != nil && int(*proxyGroup.IsSupportTLSChoice) != 1 {
					return fmt.Errorf("group(%s) not support TLS Choice", groupId)
				}
			}
			err := service.SetTlsVersion(ctx, id, tlsCiphers, tlsSupportVersion)
			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudGaapLayer7ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer7ListenerDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer7_listener.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	proxyId := d.Get("proxy_id").(string)
	groupId := d.Get("group_id").(string)
	protocol := d.Get("protocol").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteLayer7Listener(ctx, id, proxyId, groupId, protocol)
}
