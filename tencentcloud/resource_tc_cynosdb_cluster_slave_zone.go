/*
Provides a resource to create a cynosdb cluster_slave_zone

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = "cynosdbmysql-xxxxxxxx"
  slave_zone = "ap-guangzhou-3"
}
```

Import

cynosdb cluster_slave_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone cluster_slave_zone_id
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

func resourceTencentCloudCynosdbCluster_slave_zone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbCluster_slave_zoneCreate,
		Read:   resourceTencentCloudCynosdbCluster_slave_zoneRead,
		Update: resourceTencentCloudCynosdbCluster_slave_zoneUpdate,
		Delete: resourceTencentCloudCynosdbCluster_slave_zoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"slave_zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slave zone.",
			},
		},
	}
}

func resourceTencentCloudCynosdbCluster_slave_zoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewAddClusterSlaveZoneRequest()
		response  = cynosdb.NewAddClusterSlaveZoneResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("slave_zone"); ok {
		request.SlaveZone = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().AddClusterSlaveZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb cluster_slave_zone failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbCluster_slave_zoneStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbCluster_slave_zoneRead(d, meta)
}

func resourceTencentCloudCynosdbCluster_slave_zoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	cluster_slave_zoneId := d.Id()

	cluster_slave_zone, err := service.DescribeCynosdbCluster_slave_zoneById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster_slave_zone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbCluster_slave_zone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster_slave_zone.ClusterId != nil {
		_ = d.Set("cluster_id", cluster_slave_zone.ClusterId)
	}

	if cluster_slave_zone.SlaveZone != nil {
		_ = d.Set("slave_zone", cluster_slave_zone.SlaveZone)
	}

	return nil
}

func resourceTencentCloudCynosdbCluster_slave_zoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyClusterSlaveZoneRequest()

	cluster_slave_zoneId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "slave_zone"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyClusterSlaveZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb cluster_slave_zone failed, reason:%+v", logId, err)
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbCluster_slave_zoneStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbCluster_slave_zoneRead(d, meta)
}

func resourceTencentCloudCynosdbCluster_slave_zoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	cluster_slave_zoneId := d.Id()

	if err := service.DeleteCynosdbCluster_slave_zoneById(ctx, clusterId); err != nil {
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbCluster_slave_zoneStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
