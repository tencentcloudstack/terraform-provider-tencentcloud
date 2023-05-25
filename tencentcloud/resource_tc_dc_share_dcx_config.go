/*
Provides a resource to create a dc share_dcx_config

Example Usage

```hcl
resource "tencentcloud_dc_share_dcx_config" "share_dcx_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  enable = false
}
```

Import

dc share_dcx_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_share_dcx_config.share_dcx_config dcx_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
)

func resourceTencentCloudDcShareDcxConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcShareDcxConfigCreate,
		Read:   resourceTencentCloudDcShareDcxConfigRead,
		Update: resourceTencentCloudDcShareDcxConfigUpdate,
		Delete: resourceTencentCloudDcShareDcxConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_tunnel_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "the direct connect owner accept or reject the apply of direct connect tunnel.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "if accept or reject direct connect tunnel.",
			},
		},
	}
}

func resourceTencentCloudDcShareDcxConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_share_dcx_config.create")()
	defer inconsistentCheck(d, meta)()

	dcxId := d.Get("direct_connect_tunnel_id").(string)

	d.SetId(dcxId)

	return resourceTencentCloudDcShareDcxConfigUpdate(d, meta)
}

func resourceTencentCloudDcShareDcxConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_share_dcx_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	dcxId := d.Id()

	ShareDcxConfig, err := service.DescribeDcShareDcxConfigById(ctx, dcxId)
	if err != nil {
		return err
	}

	if ShareDcxConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcShareDcxConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ShareDcxConfig.DirectConnectTunnelId != nil {
		_ = d.Set("direct_connect_tunnel_id", ShareDcxConfig.DirectConnectTunnelId)
	}

	if *ShareDcxConfig.State == "AVAILABLE" {
		_ = d.Set("enable", true)
	}

	if *ShareDcxConfig.State == "REJECTED" {
		_ = d.Set("enable", false)
	}

	return nil
}

func resourceTencentCloudDcShareDcxConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_share_dcx_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enable        bool
		acceptRequest = dc.NewAcceptDirectConnectTunnelRequest()
		rejectRequest = dc.NewRejectDirectConnectTunnelRequest()
	)

	dcxId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		acceptRequest.DirectConnectTunnelId = &dcxId

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().AcceptDirectConnectTunnel(acceptRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, acceptRequest.GetAction(), acceptRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc ShareDcxConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		rejectRequest.DirectConnectTunnelId = &dcxId

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().RejectDirectConnectTunnel(rejectRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, rejectRequest.GetAction(), rejectRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc ShareDcxConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDcShareDcxConfigRead(d, meta)
}

func resourceTencentCloudDcShareDcxConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_share_dcx_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
