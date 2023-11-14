/*
Provides a resource to create a cynosdb cluster_zone

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_zone" "cluster_zone" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  old_zone = "ap-guangzhou-2"
  new_zone = "ap-guangzhou-3"
  is_in_maintain_period = "false"
}
```

Import

cynosdb cluster_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_zone.cluster_zone cluster_zone_id
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

func resourceTencentCloudCynosdbClusterZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterZoneCreate,
		Read:   resourceTencentCloudCynosdbClusterZoneRead,
		Delete: resourceTencentCloudCynosdbClusterZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster Id.",
			},

			"old_zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Active area.",
			},

			"new_zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The available zone to switch to.",
			},

			"is_in_maintain_period": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Execute - yes during maintenance, execute - no immediately.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_zone.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewSwitchClusterZoneRequest()
		response  = cynosdb.NewSwitchClusterZoneResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("old_zone"); ok {
		request.OldZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("new_zone"); ok {
		request.NewZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_in_maintain_period"); ok {
		request.IsInMaintainPeriod = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().SwitchClusterZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb clusterZone failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterZoneStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbClusterZoneRead(d, meta)
}

func resourceTencentCloudCynosdbClusterZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_zone.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbClusterZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_zone.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
