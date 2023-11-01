/*
Provides a resource to create a elasticsearch diagnose

Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose" "diagnose" {
  instance_id = "es-xxxxxx"
  cron_time = "15:00:00"
}
```

Import

es diagnose can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_diagnose.diagnose diagnose_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudElasticsearchDiagnose() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchDiagnoseCreate,
		Read:   resourceTencentCloudElasticsearchDiagnoseRead,
		Update: resourceTencentCloudElasticsearchDiagnoseUpdate,
		Delete: resourceTencentCloudElasticsearchDiagnoseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"cron_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Intelligent operation and maintenance staff regularly patrol the inspection time every day, the time format is HH:00:00, such as 15:00:00.",
			},

			"diagnose_job_metas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Diagnostic items and meta-information of intelligent operation and maintenance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English name of diagnosis item for intelligent operation and maintenance.",
						},
						"job_zh_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese name of intelligent operation and maintenance diagnosis item.",
						},
						"job_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intelligent operation and maintenance diagnostic item description.",
						},
					},
				},
			},

			"max_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of manual triggers per day for intelligent operation and maintenance staff.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchDiagnoseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	params := make(map[string]interface{})
	params["Status"] = 0
	if v, ok := d.GetOk("cron_time"); ok {
		params["CronTime"] = v.(string)
	}
	err := service.UpdateDiagnoseSettings(ctx, instanceId, params)
	if err != nil {
		return err
	}

	d.SetId(instanceId)
	return resourceTencentCloudElasticsearchDiagnoseRead(d, meta)
}

func resourceTencentCloudElasticsearchDiagnoseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	diagnoseSettings, err := service.GetDiagnoseSettingsById(ctx, instanceId)
	if err != nil {
		return err
	}

	if diagnoseSettings == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ElasticsearchDiagnose` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if diagnoseSettings.CronTime != nil {
		_ = d.Set("cron_time", diagnoseSettings.CronTime)
	}

	if diagnoseSettings.Count != nil {
		_ = d.Set("max_count", diagnoseSettings.Count)
	}

	if diagnoseSettings.DiagnoseJobMetas != nil {
		diagnoseJobMetas := []interface{}{}
		for _, item := range diagnoseSettings.DiagnoseJobMetas {
			diagnoseJobMetaMap := make(map[string]interface{})
			diagnoseJobMetaMap["job_name"] = item.JobName
			diagnoseJobMetaMap["job_zh_name"] = item.JobZhName
			diagnoseJobMetaMap["job_description"] = item.JobDescription
			diagnoseJobMetas = append(diagnoseJobMetas, diagnoseJobMetaMap)
		}
		_ = d.Set("diagnose_job_metas", diagnoseJobMetas)
	}
	return nil
}

func resourceTencentCloudElasticsearchDiagnoseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cron_time") {
		service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
		params := make(map[string]interface{})
		cronTime := d.Get("cron_time").(string)
		params["CronTime"] = cronTime
		err := service.UpdateDiagnoseSettings(ctx, instanceId, params)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudElasticsearchDiagnoseRead(d, meta)
}

func resourceTencentCloudElasticsearchDiagnoseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	params := make(map[string]interface{})
	params["Status"] = -1
	err := service.UpdateDiagnoseSettings(ctx, instanceId, params)
	if err != nil {
		return err
	}
	return nil
}
