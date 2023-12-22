package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnGatewaySslClientCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnGatewaySslClientCertCreate,
		Read:   resourceTencentCloudVpnGatewaySslClientCertRead,
		Update: resourceTencentCloudVpnGatewaySslClientCertUpdate,
		Delete: resourceTencentCloudVpnGatewaySslClientCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ssl_vpn_client_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SSL-VPN-CLIENT Instance ID.",
			},

			"switch": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      "on",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
				Description:  "`on`: Enable, `off`: Disable.",
			},
		},
	}
}

func resourceTencentCloudVpnGatewaySslClientCertCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ssl_client_cert.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	sslVpnClientId := d.Get("ssl_vpn_client_id").(string)
	d.SetId(sslVpnClientId)

	return resourceTencentCloudVpnGatewaySslClientCertUpdate(d, meta)
}

func resourceTencentCloudVpnGatewaySslClientCertRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ssl_client_cert.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	sslVpnClientId := d.Id()

	_, vpnGatewaySslClientCert, err := service.DescribeVpnSslClientById(ctx, sslVpnClientId)
	if err != nil {
		return err
	}

	if vpnGatewaySslClientCert == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpnGatewaySslClientCert` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if vpnGatewaySslClientCert.SslVpnClientId != nil {
		_ = d.Set("ssl_vpn_client_id", vpnGatewaySslClientCert.SslVpnClientId)
	}

	if vpnGatewaySslClientCert.CertStatus != nil {
		if *vpnGatewaySslClientCert.CertStatus == 1 {
			_ = d.Set("switch", "on")
		}
		if *vpnGatewaySslClientCert.CertStatus == 2 {
			_ = d.Set("switch", "off")
		}
	}

	return nil
}

func resourceTencentCloudVpnGatewaySslClientCertUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ssl_client_cert.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var taskId *uint64

	logId := tccommon.GetLogId(tccommon.ContextNil)

	sslVpnClientId := d.Id()

	certSwitch := d.Get("switch").(string)

	if certSwitch == "on" {

		var (
			request = vpc.NewEnableVpnGatewaySslClientCertRequest()
		)

		request.SslVpnClientId = &sslVpnClientId

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableVpnGatewaySslClientCert(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			taskId = result.Response.TaskId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s enable vpc vpnGatewaySslClientCert failed, reason:%+v", logId, err)
			return err
		}

	} else {

		var (
			request = vpc.NewDisableVpnGatewaySslClientCertRequest()
		)

		request.SslVpnClientId = &sslVpnClientId

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisableVpnGatewaySslClientCert(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			taskId = result.Response.TaskId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s disable vpc vpnGatewaySslClientCert failed, reason:%+v", logId, err)
			return err
		}
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := service.DescribeVpcTaskResult(ctx, helper.String(helper.UInt64ToStr(*taskId)))
	if err != nil {
		return err
	}

	return resourceTencentCloudVpnGatewaySslClientCertRead(d, meta)
}

func resourceTencentCloudVpnGatewaySslClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ssl_client_cert.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
