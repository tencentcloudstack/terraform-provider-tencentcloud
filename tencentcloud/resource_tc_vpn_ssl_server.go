/*
Provide a resource to create a VPN SSL Server.

Example Usage

```hcl
resource "tencentcloud_vpn_ssl_server" "server" {
  local_address       = [
    "10.0.0.0/17",
  ]
  remote_address      = "11.0.0.0/16"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-335lwf7d"
  ssl_vpn_protocol = "UDP"
  ssl_vpn_port = 1194
  integrity_algorithm = "MD5"
  encrypt_algorithm = "AES-128-CBC"
  compress = true
}
```

Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.server vpn-server-id
```
*/
package tencentcloud

import (
	"context"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudVpnSslServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnSslServerCreate,
		Read:   resourceTencentCloudVpnSslServerRead,
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
				ForceNew:    true,
				Description: "The name of ssl vpn server to be created.",
			},
			"local_address": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of local CIDR.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Remote CIDR for client.",
			},
			"ssl_vpn_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The protocol of ssl vpn. Default value: UDP.",
			},
			"ssl_vpn_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The port of ssl vpn. Default value: 1194.",
			},
			"integrity_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The integrity algorithm. Valid values: SHA1, MD5 and NONE. Default value: NONE.",
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The encrypt algorithm. Valid values: AES-128-CBC, AES-192-CBC, AES-256-CBC, NONE." +
					"Default value: NONE.",
			},
			"compress": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     FALSE,
				Description: "need compressed. Default value: False.",
			},
		},
	}
}

func resourceTencentCloudVpnSslServerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_ssl_server.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		vpcService   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
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
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := vpcService.client.UseVpcClient().CreateVpnGatewaySslServer(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
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
	defer logElapsed("resource.tencentcloud_vpn_ssl_server.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sslServerId := d.Id()
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		has, info, e := vpcService.DescribeVpnSslServerById(ctx, sslServerId)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("vpn_gateway_id", info.VpnGatewayId)
		_ = d.Set("ssl_vpn_server_name", info.SslVpnServerName)
		_ = d.Set("local_address", helper.StringsInterfaces(info.LocalAddress))
		_ = d.Set("remote_address", info.RemoteAddress)
		if _, ok := d.GetOk("ssl_vpn_protocol"); ok {
			_ = d.Set("ssl_vpn_protocol", info.SslVpnProtocol)
		}
		if _, ok := d.GetOk("ssl_vpn_port"); ok {
			_ = d.Set("ssl_vpn_port", info.SslVpnPort)
		}
		if _, ok := d.GetOk("integrity_algorithm"); ok {
			_ = d.Set("integrity_algorithm", info.IntegrityAlgorithm)
		}
		if _, ok := d.GetOk("encrypt_algorithm"); ok {
			_ = d.Set("encrypt_algorithm", info.EncryptAlgorithm)
		}
		if _, ok := d.GetOk("compress"); ok {
			_ = d.Set("compress", info.Compress)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudVpnSslServerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_ssl_server.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

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
