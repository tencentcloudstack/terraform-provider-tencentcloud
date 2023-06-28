/*
Provides a resource to create a cynosdb reload_proxy_node

Example Usage

```hcl
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = "cynosdbmysql-cgd2gpwr"
  proxy_group_id = "cynosdbmysql-proxy-8lqtl8pk"
}
```

Import

cynosdb reload_proxy_node can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node reload_proxy_node_id
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

func resourceTencentCloudCynosdbReloadProxyNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReloadProxyNodeCreate,
		Read:   resourceTencentCloudCynosdbReloadProxyNodeRead,
		Delete: resourceTencentCloudCynosdbReloadProxyNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "cluster id.",
			},
			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "proxy group id.",
			},
		},
	}
}

func resourceTencentCloudCynosdbReloadProxyNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request   = cynosdb.NewReloadBalanceProxyNodeRequest()
		response  = cynosdb.NewReloadBalanceProxyNodeResponse()
		clusterId string
		flowId    int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		request.ProxyGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ReloadBalanceProxyNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("cynosdb proxyNode not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb proxyNode failed, reason:%+v", logId, err)
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
			return resource.RetryableError(fmt.Errorf("cynosdb proxy is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxy fail, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbReloadProxyNodeRead(d, meta)
}

func resourceTencentCloudCynosdbReloadProxyNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbReloadProxyNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
