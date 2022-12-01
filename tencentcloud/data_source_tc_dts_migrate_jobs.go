/*
Use this data source to query detailed information of dts migrateJobs

Example Usage

```hcl
data "tencentcloud_dts_migrate_jobs" "migrateJobs" {
  job_id = ""
  job_name = ""
  status = ""
  src_instance_id = ""
  src_region = ""
  src_database_type = ""
  src_access_type = ""
  dst_instance_id = ""
  dst_region = ""
  dst_database_type = ""
  dst_access_type = ""
  run_mode = ""
  order_seq = ""
  tag_filters {
			tag_key = ""
			tag_value = ""

  }
}
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDtsMigrateJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsMigrateJobsRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job id.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job name.",
			},

			"status": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "migrate status.",
			},

			"src_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance id.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source region.",
			},

			"src_database_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "source database type.",
			},

			"src_access_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "source access type.",
			},

			"dst_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance id.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "destination region.",
			},

			"dst_database_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "destination database type.",
			},

			"dst_access_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "destination access type.",
			},

			"run_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "run mode.",
			},

			"order_seq": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "order by, default by create time .",
			},

			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "tag filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDtsMigrateJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dts_migrate_jobs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["job_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		paramMap["job_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		for i := range statusSet {
			status := statusSet[i].(string)
			paramMap["status"] = append(paramMap["status"], &status)
		}
	}

	if v, ok := d.GetOk("src_instance_id"); ok {
		paramMap["src_instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_region"); ok {
		paramMap["src_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_database_type"); ok {
		src_database_typeSet := v.(*schema.Set).List()
		for i := range src_database_typeSet {
			src_database_type := src_database_typeSet[i].(string)
			paramMap["src_database_type"] = append(paramMap["src_database_type"], &src_database_type)
		}
	}

	if v, ok := d.GetOk("src_access_type"); ok {
		src_access_typeSet := v.(*schema.Set).List()
		for i := range src_access_typeSet {
			src_access_type := src_access_typeSet[i].(string)
			paramMap["src_access_type"] = append(paramMap["src_access_type"], &src_access_type)
		}
	}

	if v, ok := d.GetOk("dst_instance_id"); ok {
		paramMap["dst_instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_region"); ok {
		paramMap["dst_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_database_type"); ok {
		dst_database_typeSet := v.(*schema.Set).List()
		for i := range dst_database_typeSet {
			dst_database_type := dst_database_typeSet[i].(string)
			paramMap["dst_database_type"] = append(paramMap["dst_database_type"], &dst_database_type)
		}
	}

	if v, ok := d.GetOk("dst_access_type"); ok {
		dst_access_typeSet := v.(*schema.Set).List()
		for i := range dst_access_typeSet {
			dst_access_type := dst_access_typeSet[i].(string)
			paramMap["dst_access_type"] = append(paramMap["dst_access_type"], &dst_access_type)
		}
	}

	if v, ok := d.GetOk("run_mode"); ok {
		paramMap["run_mode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_seq"); ok {
		paramMap["order_seq"] = helper.String(v.(string))
	}

	dtsService := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var jobList []*dts.JobItem
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dtsService.DescribeDtsMigrateJobsByFilter(ctx, param)
		if e != nil {
			return retryError(e)
		}
		jobList = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dts jobList failed, reason:%+v", logId, err)
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), List); e != nil {
			return e
		}
	}

	return nil
}
