package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudVpnSslServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnSslServerCreate,
		Read:   resourceTencentCloudVpnSslServerRead,
		Update: resourceTencentCloudVpnSslServerUpdate,
		Delete: resourceTencentCloudVpnSslServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPN gateway ID.",
			},
			"ssl_vpn_server_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of ssl vpn server to be created.",
			},
			"local_address": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of local CIDR.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote CIDR for client.",
			},
			"ssl_vpn_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The protocol of ssl vpn. Default value: UDP.",
			},
			"ssl_vpn_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The port of ssl vpn. Currently only supports UDP. Default value: 1194.",
			},
			"integrity_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The integrity algorithm. Valid values: SHA1. Default value: SHA1.",
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The encrypt algorithm. Valid values: AES-128-CBC, AES-192-CBC, AES-256-CBC." +
					"Default value: AES-128-CBC.",
			},
			"compress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     svccvm.FALSE,
				Description: "Need compressed. Currently is not supports compress. Default value: False.",
			},
			"sso_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSO authentication. Default: false. This feature requires whitelist approval.",
			},
			"access_policy_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Enable access policy control. Default: false.",
			},
			"saml_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML-DATA. Required when sso_enabled is true.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags for resource management.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "DNS server configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_dns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary DNS server address.",
						},
						"secondary_dns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Secondary DNS server address.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpnSslServerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_server.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		vpcService   = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request      = vpc.NewCreateVpnGatewaySslServerRequest()
		vpnGatewayId string
	)

	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(vpnGatewayId)
	}
	if v, ok := d.GetOk("ssl_vpn_server_name"); ok {
		request.SslVpnServerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("local_address"); ok {
		address := v.([]interface{})
		request.LocalAddress = helper.InterfacesStringsPoint(address)
	}

	if v, ok := d.GetOk("remote_address"); ok {
		request.RemoteAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ssl_vpn_protocol"); ok {
		request.SslVpnProtocol = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("ssl_vpn_port"); ok {
		request.SslVpnPort = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("integrity_algorithm"); ok {
		request.IntegrityAlgorithm = helper.String(v.(string))
	}
	if v, ok := d.GetOk("encrypt_algorithm"); ok {
		request.EncryptAlgorithm = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("compress"); ok {
		request.Compress = helper.Bool(v.(bool))
	}

	// SSO authentication
	if v, ok := d.GetOkExists("sso_enabled"); ok {
		request.SsoEnabled = helper.Bool(v.(bool))
	}

	// Access policy control
	if v, ok := d.GetOkExists("access_policy_enabled"); ok {
		request.AccessPolicyEnabled = helper.Bool(v.(bool))
	}

	// SAML data
	if v, ok := d.GetOk("saml_data"); ok {
		request.SamlData = helper.String(v.(string))
	}

	// Note: Tags will be handled by Tag Service after creation

	// DNS servers
	if v, ok := d.GetOk("dns_servers"); ok {
		dnsServersList := v.([]interface{})
		if len(dnsServersList) > 0 {
			dnsServersMap := dnsServersList[0].(map[string]interface{})
			request.DnsServers = &vpc.DnsServers{}
			if primaryDns, ok := dnsServersMap["primary_dns"]; ok && primaryDns.(string) != "" {
				request.DnsServers.PrimaryDns = helper.String(primaryDns.(string))
			}
			if secondaryDns, ok := dnsServersMap["secondary_dns"]; ok && secondaryDns.(string) != "" {
				request.DnsServers.SecondaryDns = helper.String(secondaryDns.(string))
			}
		}
	}

	var (
		taskId      *int64
		sslServerId *string
	)
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpnGatewaySslServer(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		taskId = response.Response.TaskId
		sslServerId = response.Response.SslVpnServerId
		return nil
	}); err != nil {
		return err
	}

	err := vpcService.DescribeVpcTaskResult(ctx, helper.String(helper.Int64ToStr(*taskId)))
	if err != nil {
		return err
	}

	d.SetId(*sslServerId)

	// Handle tags using Tag Service
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpns", region, *sslServerId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnSslServerRead(d, meta)
}

func resourceTencentCloudVpnSslServerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_server.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sslServerId := d.Id()
	vpcService := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, info, e := vpcService.DescribeVpnSslServerById(ctx, sslServerId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("vpn_gateway_id", info.VpnGatewayId)
		_ = d.Set("ssl_vpn_server_name", info.SslVpnServerName)
		_ = d.Set("local_address", helper.StringsInterfaces(info.LocalAddress))
		_ = d.Set("remote_address", info.RemoteAddress)
		_ = d.Set("ssl_vpn_protocol", info.SslVpnProtocol)
		_ = d.Set("ssl_vpn_port", info.SslVpnPort)
		_ = d.Set("integrity_algorithm", info.IntegrityAlgorithm)

		_ = d.Set("encrypt_algorithm", info.EncryptAlgorithm)

		compress := *info.Compress
		_ = d.Set("compress", false)
		if compress != 0 {
			_ = d.Set("compress", true)
		}

		// SSO authentication - convert uint64 to bool
		if info.SsoEnabled != nil {
			_ = d.Set("sso_enabled", *info.SsoEnabled != 0)
		}

		// Access policy control - convert uint64 to bool
		if info.AccessPolicyEnabled != nil {
			_ = d.Set("access_policy_enabled", *info.AccessPolicyEnabled != 0)
		}

		// Note: SamlData is not returned by DescribeVpnGatewaySslServers API
		// It will remain as configured by user (Computed attribute)

		// DNS servers
		if info.DnsServers != nil {
			dnsServersMap := map[string]interface{}{}
			if info.DnsServers.PrimaryDns != nil {
				dnsServersMap["primary_dns"] = *info.DnsServers.PrimaryDns
			}
			if info.DnsServers.SecondaryDns != nil {
				dnsServersMap["secondary_dns"] = *info.DnsServers.SecondaryDns
			}
			if len(dnsServersMap) > 0 {
				_ = d.Set("dns_servers", []interface{}{dnsServersMap})
			}
		}

		// Tags - use Tag Service to read
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpns", region, sslServerId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		_ = d.Set("tags", tags)

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudVpnSslServerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_server.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		vpcService = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request    = vpc.NewModifyVpnGatewaySslServerRequest()
	)

	sslServerId := d.Id()
	request.SslVpnServerId = helper.String(sslServerId)

	needChange := false
	mutableArgs := []string{
		"ssl_vpn_server_name", "local_address", "remote_address", "ssl_vpn_protocol",
		"ssl_vpn_port", "integrity_algorithm", "encrypt_algorithm", "compress",
		"sso_enabled", "saml_data", "dns_servers",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("ssl_vpn_server_name"); ok {
			request.SslVpnServerName = helper.String(v.(string))
		}
		if v, ok := d.GetOk("local_address"); ok {
			address := v.([]interface{})
			request.LocalAddress = helper.InterfacesStringsPoint(address)
		}
		if v, ok := d.GetOk("remote_address"); ok {
			request.RemoteAddress = helper.String(v.(string))
		}
		if v, ok := d.GetOk("ssl_vpn_protocol"); ok {
			request.SslVpnProtocol = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("ssl_vpn_port"); ok {
			request.SslVpnPort = helper.IntInt64(v.(int))
		}
		if v, ok := d.GetOk("integrity_algorithm"); ok {
			request.IntegrityAlgorithm = helper.String(v.(string))
		}
		if v, ok := d.GetOk("encrypt_algorithm"); ok {
			request.EncryptAlgorithm = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("compress"); ok {
			request.Compress = helper.Bool(v.(bool))
		}

		// SSO authentication
		if v, ok := d.GetOkExists("sso_enabled"); ok {
			request.SsoEnabled = helper.Bool(v.(bool))
		}

		// SAML data
		if v, ok := d.GetOk("saml_data"); ok {
			request.SamlData = helper.String(v.(string))
		}

		// Note: AccessPolicyEnabled is not supported by ModifyVpnGatewaySslServer API
		// It can only be set during creation

		// Note: Tags are handled separately using Tag Service

		// DNS servers
		if v, ok := d.GetOk("dns_servers"); ok {
			dnsServersList := v.([]interface{})
			if len(dnsServersList) > 0 {
				dnsServersMap := dnsServersList[0].(map[string]interface{})
				request.DnsServers = &vpc.DnsServers{}
				if primaryDns, ok := dnsServersMap["primary_dns"]; ok && primaryDns.(string) != "" {
					request.DnsServers.PrimaryDns = helper.String(primaryDns.(string))
				}
				if secondaryDns, ok := dnsServersMap["secondary_dns"]; ok && secondaryDns.(string) != "" {
					request.DnsServers.SecondaryDns = helper.String(secondaryDns.(string))
				}
			}
		}
	}

	var (
		taskId *int64
	)
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpnGatewaySslServer(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		taskId = response.Response.TaskId
		return nil
	}); err != nil {
		return err
	}

	err := vpcService.DescribeVpcTaskResult(ctx, helper.String(helper.Int64ToStr(*taskId)))
	if err != nil {
		return err
	}

	// Handle tags using Tag Service
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpns", region, sslServerId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnSslServerRead(d, meta)
}

func resourceTencentCloudVpnSslServerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_server.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	serverId := d.Id()

	taskId, err := service.DeleteVpnGatewaySslServer(ctx, serverId)
	if err != nil {
		return err
	}

	err = service.DescribeVpcTaskResult(ctx, helper.String(helper.UInt64ToStr(taskId)))
	if err != nil {
		return err
	}

	return nil
}
