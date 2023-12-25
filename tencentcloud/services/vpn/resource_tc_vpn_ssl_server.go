package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
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
				Description: "The protocol of ssl vpn. Default value: UDP.",
			},
			"ssl_vpn_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of ssl vpn. Default value: 1194.",
			},
			"integrity_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The integrity algorithm. Valid values: SHA1, MD5 and NONE. Default value: NONE.",
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The encrypt algorithm. Valid values: AES-128-CBC, AES-192-CBC, AES-256-CBC, NONE." +
					"Default value: NONE.",
			},
			"compress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     svccvm.FALSE,
				Description: "need compressed. Default value: False.",
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
	if v, ok := d.GetOk("ssl_vpn_port"); ok {
		request.SslVpnPort = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("integrity_algorithm"); ok {
		request.IntegrityAlgorithm = helper.String(v.(string))
	}
	if v, ok := d.GetOk("encrypt_algorithm"); ok {
		request.EncryptAlgorithm = helper.String(v.(string))
	}
	if v, ok := d.GetOk("compress"); ok {
		request.Compress = helper.Bool(v.(bool))
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
		if v, ok := d.GetOk("ssl_vpn_port"); ok {
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
