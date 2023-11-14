/*
Provides a resource to create a cynosdb upgrade_proxy_version

Example Usage

```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id = "cynosdbmysql-xxxxxx"
  src_proxy_version = "无"
  dst_proxy_version = "无"
  proxy_group_id = "无"
  is_in_maintain_period = "no"
}
```

Import

cynosdb upgrade_proxy_version can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version upgrade_proxy_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbUpgradeProxyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbUpgradeProxyVersionCreate,
		Read:   resourceTencentCloudCynosdbUpgradeProxyVersionRead,
		Update: resourceTencentCloudCynosdbUpgradeProxyVersionUpdate,
		Delete: resourceTencentCloudCynosdbUpgradeProxyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"src_proxy_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Current Version.",
			},

			"dst_proxy_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Upgrade Version.",
			},

			"proxy_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Group ID.",
			},

			"is_in_maintain_period": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upgrade time: no (when the upgrade is completed) yes (instance maintenance time).",
			},
		},
	}
}

func resourceTencentCloudCynosdbUpgradeProxyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.create")()
	defer inconsistentCheck(d, meta)()

	var proxyGroupId string
	if v, ok := d.GetOk("proxy_group_id"); ok {
		proxyGroupId = v.(string)
	}

	d.SetId(proxyGroupId)

	return resourceTencentCloudCynosdbUpgradeProxyVersionUpdate(d, meta)
}

func resourceTencentCloudCynosdbUpgradeProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	upgradeProxyVersionId := d.Id()

	upgradeProxyVersion, err := service.DescribeCynosdbUpgradeProxyVersionById(ctx, proxyGroupId)
	if err != nil {
		return err
	}

	if upgradeProxyVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbUpgradeProxyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if upgradeProxyVersion.ClusterId != nil {
		_ = d.Set("cluster_id", upgradeProxyVersion.ClusterId)
	}

	if upgradeProxyVersion.SrcProxyVersion != nil {
		_ = d.Set("src_proxy_version", upgradeProxyVersion.SrcProxyVersion)
	}

	if upgradeProxyVersion.DstProxyVersion != nil {
		_ = d.Set("dst_proxy_version", upgradeProxyVersion.DstProxyVersion)
	}

	if upgradeProxyVersion.ProxyGroupId != nil {
		_ = d.Set("proxy_group_id", upgradeProxyVersion.ProxyGroupId)
	}

	if upgradeProxyVersion.IsInMaintainPeriod != nil {
		_ = d.Set("is_in_maintain_period", upgradeProxyVersion.IsInMaintainPeriod)
	}

	return nil
}

func resourceTencentCloudCynosdbUpgradeProxyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewUpgradeProxyVersionRequest()

	upgradeProxyVersionId := d.Id()

	request.ProxyGroupId = &proxyGroupId

	immutableArgs := []string{"cluster_id", "src_proxy_version", "dst_proxy_version", "proxy_group_id", "is_in_maintain_period"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("src_proxy_version") {
		if v, ok := d.GetOk("src_proxy_version"); ok {
			request.SrcProxyVersion = helper.String(v.(string))
		}
	}

	if d.HasChange("dst_proxy_version") {
		if v, ok := d.GetOk("dst_proxy_version"); ok {
			request.DstProxyVersion = helper.String(v.(string))
		}
	}

	if d.HasChange("is_in_maintain_period") {
		if v, ok := d.GetOk("is_in_maintain_period"); ok {
			request.IsInMaintainPeriod = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeProxyVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb upgradeProxyVersion failed, reason:%+v", logId, err)
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbUpgradeProxyVersionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbUpgradeProxyVersionRead(d, meta)
}

func resourceTencentCloudCynosdbUpgradeProxyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
