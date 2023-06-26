/*
Provides a resource to create a cynosdb binlog_save_days

Example Usage

```hcl
resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id       = "cynosdbmysql-123"
  binlog_save_days = 7
}
```

Import

cynosdb binlog_save_days can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_binlog_save_days.binlog_save_days binlog_save_days_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbBinlogSaveDays() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbBinlogSaveDaysCreate,
		Read:   resourceTencentCloudCynosdbBinlogSaveDaysRead,
		Update: resourceTencentCloudCynosdbBinlogSaveDaysUpdate,
		Delete: resourceTencentCloudCynosdbBinlogSaveDaysDelete,
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

			"binlog_save_days": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Binlog retention days.",
			},
		},
	}
}

func resourceTencentCloudCynosdbBinlogSaveDaysCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_binlog_save_days.create")()
	defer inconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbBinlogSaveDaysUpdate(d, meta)
}

func resourceTencentCloudCynosdbBinlogSaveDaysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_binlog_save_days.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()
	binlogSaveDays, err := service.DescribeCynosdbBinlogSaveDaysById(ctx, clusterId)
	if err != nil {
		return err
	}

	if binlogSaveDays == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbBinlogSaveDays` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if binlogSaveDays != nil {
		_ = d.Set("binlog_save_days", binlogSaveDays)
	}

	return nil
}

func resourceTencentCloudCynosdbBinlogSaveDaysUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_binlog_save_days.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyBinlogSaveDaysRequest()

	clusterId := d.Id()
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("binlog_save_days"); ok {
		request.BinlogSaveDays = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyBinlogSaveDays(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb binlogSaveDays failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbBinlogSaveDaysRead(d, meta)
}

func resourceTencentCloudCynosdbBinlogSaveDaysDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_binlog_save_days.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
