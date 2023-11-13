/*
Provides a resource to create a cynosdb cluster_version

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_version" "cluster_version" {
  cluster_id = "xxx"
  cynos_version = "2.0.0"
  upgrade_type = "upgradeImmediate"
}
```

Import

cynosdb cluster_version can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_version.cluster_version cluster_version_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbClusterVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterVersionCreate,
		Read:   resourceTencentCloudCynosdbClusterVersionRead,
		Delete: resourceTencentCloudCynosdbClusterVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cynos_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Kernel version.",
			},

			"upgrade_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Upgrade time type, optional: upgradeImmediate, upgradeInMaintain.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewUpgradeClusterVersionRequest()
		response  = cynosdb.NewUpgradeClusterVersionResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cynos_version"); ok {
		request.CynosVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upgrade_type"); ok {
		request.UpgradeType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeClusterVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb clusterVersion failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterVersionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbClusterVersionRead(d, meta)
}

func resourceTencentCloudCynosdbClusterVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbClusterVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
