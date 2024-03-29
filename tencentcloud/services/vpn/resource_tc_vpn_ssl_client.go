package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudVpnSslClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnSslClientCreate,
		Read:   resourceTencentCloudVpnSslClientRead,
		Delete: resourceTencentCloudVpnSslClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPN ssl server id.",
			},
			"ssl_vpn_client_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of ssl vpn client to be created.",
			},
		},
	}
}

func resourceTencentCloudVpnSslClientCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_client.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		vpcService       = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request          = vpc.NewCreateVpnGatewaySslClientRequest()
		sslVpnServerId   string
		sslVpnClientName string
	)

	if v, ok := d.GetOk("ssl_vpn_server_id"); ok {
		sslVpnServerId = v.(string)
		request.SslVpnServerId = helper.String(sslVpnServerId)
	}
	if v, ok := d.GetOk("ssl_vpn_client_name"); ok {
		sslVpnClientName = v.(string)
		request.SslVpnClientName = helper.String(sslVpnClientName)
	}

	var (
		taskId      *uint64
		sslClientId *string
	)
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpnGatewaySslClient(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		taskId = response.Response.TaskId
		sslClientId = response.Response.SslVpnClientId
		return nil
	}); err != nil {
		return err
	}

	err := vpcService.DescribeVpcTaskResult(ctx, helper.String(helper.UInt64ToStr(*taskId)))
	if err != nil {
		return err
	}

	d.SetId(*sslClientId)

	return resourceTencentCloudVpnSslClientRead(d, meta)
}

func resourceTencentCloudVpnSslClientRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_client.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sslClientId := d.Id()
	vpcService := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, info, e := vpcService.DescribeVpnSslClientById(ctx, sslClientId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("ssl_vpn_server_id", info.SslVpnServerId)
		_ = d.Set("ssl_vpn_client_name", info.Name)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudVpnSslClientDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_ssl_client.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	sslClientId := d.Id()

	taskId, err := service.DeleteVpnGatewaySslClient(ctx, sslClientId)
	if err != nil {
		return err
	}

	err = service.DescribeVpcTaskResult(ctx, helper.String(helper.UInt64ToStr(*taskId)))
	if err != nil {
		return err
	}

	return err
}
