package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbExportInstanceErrorLogs() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbExportInstanceErrorLogsCreate,
		Read:   resourceTencentCloudCynosdbExportInstanceErrorLogsRead,
		Delete: resourceTencentCloudCynosdbExportInstanceErrorLogsDelete,

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
				Description: "keyword.",
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

			"error_log_item_export": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of instances in the read-write instance group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "time.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log level, optional values note, warning, error.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "log content.",
						},
					},
				},
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

	d.SetId(instanceId)

	items := response.Response.ErrorLogItems
	if items != nil {
		logItemList := []interface{}{}
		for _, logItem := range items {
			logItemMap := map[string]interface{}{}
			if logItem.Timestamp != nil {
				logItemMap["timestamp"] = logItem.Timestamp
			}
			if logItem.Level != nil {
				logItemMap["level"] = logItem.Level
			}
			if logItem.Content != nil {
				logItemMap["content"] = logItem.Content
			}

			logItemList = append(logItemList, logItemMap)
		}
		_ = d.Set("error_log_item_export", logItemList)
	}

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
