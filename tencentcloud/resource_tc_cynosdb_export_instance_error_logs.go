/*
Provides a resource to create a cynosdb export_instance_error_logs

Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_error_logs" "export_instance_error_logs" {
  instance_id = "cynosdbmysql-ins-123"
  start_time = "2022-01-01 12:00:00"
  end_time = "2022-01-01 14:00:00"
  log_levels =
  key_words =
  file_type = "csv"
  order_by = "Timestamp"
  order_by_type = "ASC"
}
```

Import

cynosdb export_instance_error_logs can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_export_instance_error_logs.export_instance_error_logs export_instance_error_logs_id
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

func resourceTencentCloudCynosdbExportInstanceErrorLogs() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbExportInstanceErrorLogsCreate,
		Read:   resourceTencentCloudCynosdbExportInstanceErrorLogsRead,
		Delete: resourceTencentCloudCynosdbExportInstanceErrorLogsDelete,
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
				Description: "Log earliest time.",
			},

			"end_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Latest log time.",
			},

			"log_levels": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Log level.",
			},

			"key_words": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Keyword.",
			},

			"file_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File type, optional values: csv, original.",
			},

			"order_by": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Optional value Timestamp.",
			},

			"order_by_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ASC or DESC.",
			},
		},
	}
}

func resourceTencentCloudCynosdbExportInstanceErrorLogsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_error_logs.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewExportInstanceErrorLogsRequest()
		response   = cynosdb.NewExportInstanceErrorLogsResponse()
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

	if v, ok := d.GetOk("log_levels"); ok {
		logLevelsSet := v.(*schema.Set).List()
		for i := range logLevelsSet {
			logLevels := logLevelsSet[i].(string)
			request.LogLevels = append(request.LogLevels, &logLevels)
		}
	}

	if v, ok := d.GetOk("key_words"); ok {
		keyWordsSet := v.(*schema.Set).List()
		for i := range keyWordsSet {
			keyWords := keyWordsSet[i].(string)
			request.KeyWords = append(request.KeyWords, &keyWords)
		}
	}

	if v, ok := d.GetOk("file_type"); ok {
		request.FileType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		request.OrderByType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ExportInstanceErrorLogs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb exportInstanceErrorLogs failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCynosdbExportInstanceErrorLogsRead(d, meta)
}

func resourceTencentCloudCynosdbExportInstanceErrorLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_error_logs.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbExportInstanceErrorLogsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_export_instance_error_logs.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
