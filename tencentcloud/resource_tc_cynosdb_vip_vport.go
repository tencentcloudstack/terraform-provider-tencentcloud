/*
Provides a resource to create a cynosdb vip_vport

Example Usage

```hcl
resource "tencentcloud_cynosdb_vip_vport" "vip_vport" {
  cluster_id = "xxx"
  instance_grp_id = "xxx"
  vip = "xx.xx.xx.xx"
  vport = 5432
  db_type = "MYSQL"
  old_ip_reserve_hours = 0
}
```

Import

cynosdb vip_vport can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_vip_vport.vip_vport vip_vport_id
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

func resourceTencentCloudCynosdbVipVport() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbVipVportCreate,
		Read:   resourceTencentCloudCynosdbVipVportRead,
		Delete: resourceTencentCloudCynosdbVipVportDelete,
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

			"instance_grp_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance Group ID.",
			},

			"vip": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Purpose IP that needs to be modified.",
			},

			"vport": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Destination port that needs to be modified.",
			},

			"db_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database type, value range:&amp;amp;lt;li&amp;amp;gt;MYSQL&amp;amp;lt;/li&amp;amp;gt;.",
			},

			"old_ip_reserve_hours": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The retention time of the old IP before recycling, in hours, 0 means immediate recycling.",
			},
		},
	}
}

func resourceTencentCloudCynosdbVipVportCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_vip_vport.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewModifyVipVportRequest()
		response  = cynosdb.NewModifyVipVportResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_grp_id"); ok {
		request.InstanceGrpId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	if v, _ := d.GetOk("vport"); v != nil {
		request.Vport = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_type"); ok {
		request.DbType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("old_ip_reserve_hours"); v != nil {
		request.OldIpReserveHours = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyVipVport(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb vipVport failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbVipVportStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbVipVportRead(d, meta)
}

func resourceTencentCloudCynosdbVipVportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_vip_vport.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbVipVportDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_vip_vport.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
