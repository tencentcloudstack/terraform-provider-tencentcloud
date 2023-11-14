/*
Provides a resource to create a cynosdb export_instance_slow_queries

Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_slow_queries" "export_instance_slow_queries" {
  instance_id = "cynosdbmysql-ins-123"
  start_time = "2022-01-01 12:00:00"
  end_time = "2022-01-01 14:00:00"
  username = "root"
  host = "10.10.10.10"
  database = "db1"
  file_type = "csv"
}
```

Import

cynosdb export_instance_slow_queries can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_export_instance_slow_queries.export_instance_slow_queries export_instance_slow_queries_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbExportInstanceSlowQueries() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbExportInstanceSlowQueriesCreate,
		Read:   resourceTencentCloudCynosdbExportInstanceSlowQueriesRead,
		Delete: resourceTencentCloudCynosdbExportInstanceSlowQueriesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"start_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Earliest transaction start time.",
			},

			"end_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Latest transaction start time.",
			},

			"username": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "User name.",
			},

			"host": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Client host.",
			},

			"database": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"file_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File type, optional values: csv, original.",
			},
		},
	}
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewExportInstanceSlowQueriesRequest()
		response   = cynosdb.NewExportInstanceSlowQueriesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		request.Database = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_type"); ok {
		request.FileType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ExportInstanceSlowQueries(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb exportInstanceSlowQueries failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCynosdbExportInstanceSlowQueriesRead(d, meta)
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
