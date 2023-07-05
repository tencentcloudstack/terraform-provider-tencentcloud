/*
Provides a resource to create a cynosdb upgrade_proxy_version

Example Usage

Specify `proxy_group_id` modification
```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id        = "cynosdbmysql-bws8h88b"
  proxy_group_id    = "cynosdbmysql-proxy-laz8hd6c"
  src_proxy_version = "1.3.5"
  dst_proxy_version = "1.3.7"
}
```

Modify all proxy database versions in the current cluster
```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id        = "cynosdbmysql-bws8h88b"
  src_proxy_version = "1.3.5"
  dst_proxy_version = "1.3.7"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbUpgradeProxyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbUpgradeProxyVersionCreate,
		Read:   resourceTencentCloudCynosdbUpgradeProxyVersionRead,
		Delete: resourceTencentCloudCynosdbUpgradeProxyVersionDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"proxy_group_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Group ID.",
			},
			"src_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Current Version.",
			},
			"dst_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Upgrade Version.",
			},
		},
	}
}

func resourceTencentCloudCynosdbUpgradeProxyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request   = cynosdb.NewUpgradeProxyVersionRequest()
		response  = cynosdb.NewUpgradeProxyVersionResponse()
		clusterId string
		flowId    int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("src_proxy_version"); ok {
		request.SrcProxyVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_proxy_version"); ok {
		request.DstProxyVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		request.ProxyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_in_maintain_period"); ok {
		request.IsInMaintainPeriod = helper.String(v.(string))
	}

	request.IsInMaintainPeriod = helper.String("no")

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeProxyVersion(request)
		if e != nil {
			return retryError(e)

		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb upgradeProxyVersion failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("operate cynosdb upgradeProxyVersion is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb upgradeProxyVersion fail, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(clusterId)
	return resourceTencentCloudCynosdbUpgradeProxyVersionRead(d, meta)
}

func resourceTencentCloudCynosdbUpgradeProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbUpgradeProxyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
