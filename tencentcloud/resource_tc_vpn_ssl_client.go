/*
Provide a resource to create a VPN SSL Client.

Example Usage

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id = "vpns-aog5xcjj"
  ssl_vpn_client_name = "hello"
}

```

Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudVpnSslClient() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_vpn_ssl_client.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		vpcService   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		request      = vpc.NewCreateVpnGatewaySslClientRequest()
		sslVpnServerId string
	)

	if v, ok := d.GetOk("ssl_vpn_server_id"); ok {
		sslVpnServerId = v.(string)
		request.SslVpnServerId = helper.String(sslVpnServerId)
	}
	if v, ok := d.GetOk("ssl_vpn_client_name"); ok {
		request.SslVpnClientName = helper.String(v.(string))
	}

	var taskId *uint64
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := vpcService.client.UseVpcClient().CreateVpnGatewaySslClient(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		taskId = response.Response.TaskId
		return nil
	}); err != nil {
		return err
	}

	err := vpcService.DescribeTaskResult(ctx, helper.Uint64(*taskId))
	if err != nil {
		return err
	}

	// add protect
	time.Sleep(3)

	filter := make(map[string]string)
	filter["ssl-vpn-server-id"] = sslVpnServerId

	instances, err := vpcService.DescribeVpnGwSslClientByFilter(ctx, filter)

	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	sslClient := instances[0]
	d.SetId(*sslClient.SslVpnClientId)

	return resourceTencentCloudVpnSslClientRead(d, meta)
}

func resourceTencentCloudVpnSslClientRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_ssl_client.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sslClientId := d.Id()
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		has, info, e := vpcService.DescribeVpnSslClientById(ctx, sslClientId)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_vpn_ssl_client.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	sslClientId := d.Id()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteVpnGatewaySslClient(ctx, sslClientId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
