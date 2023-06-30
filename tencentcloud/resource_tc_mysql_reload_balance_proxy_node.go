/*
Provides a resource to create a mysql reload_balance_proxy_node

Example Usage

```hcl
resource "tencentcloud_mysql_reload_balance_proxy_node" "reload_balance_proxy_node" {
  proxy_group_id = "proxy-gmi1f78l"
  proxy_address_id = "proxyaddr-4wc4y1pq"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlReloadBalanceProxyNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlReloadBalanceProxyNodeCreate,
		Read:   resourceTencentCloudMysqlReloadBalanceProxyNodeRead,
		Delete: resourceTencentCloudMysqlReloadBalanceProxyNodeDelete,

		Schema: map[string]*schema.Schema{
			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Proxy id.",
			},

			"proxy_address_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Proxy address id.",
			},
		},
	}
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewReloadBalanceProxyNodeRequest()
		instanceId string
	)
	if v, ok := d.GetOk("proxy_group_id"); ok {
		instanceId = v.(string)
		request.ProxyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_address_id"); ok {
		request.ProxyAddressId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ReloadBalanceProxyNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql reloadBalanceProxyNode failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlReloadBalanceProxyNodeRead(d, meta)
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
