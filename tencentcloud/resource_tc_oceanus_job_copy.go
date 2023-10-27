/*
Provides a resource to create a oceanus job_copy

Example Usage

```hcl
resource "tencentcloud_oceanus_job_copy" "example" {
  source_id         = "cql-0nob2hx8"
  target_cluster_id = "cluster-1kcd524h"
  source_name       = "keep_jar"
  target_name       = "tf_copy_example"
  target_folder_id  = "folder-7ctl246z"
  job_type          = 2
  work_space_id     = "space-2idq8wbr"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusJobCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusJobCopyCreate,
		Read:   resourceTencentCloudOceanusJobCopyRead,
		Delete: resourceTencentCloudOceanusJobCopyDelete,

		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The serial ID of the job to be copied.",
			},
			"source_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the job to be copied.",
			},
			"target_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster serial ID of the target cluster.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the new job.",
			},
			"target_folder_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The directory ID of the new job.",
			},
			"job_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of the source job.",
			},
			"work_space_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"job_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Copy Job ID.",
			},
		},
	}
}

func resourceTencentCloudOceanusJobCopyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job_copy.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = oceanus.NewCopyJobsRequest()
		response = oceanus.NewCopyJobsResponse()
		jobId    string
	)

	copyJobItem := oceanus.CopyJobItem{}
	if v, ok := d.GetOk("source_id"); ok {
		copyJobItem.SourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_cluster_id"); ok {
		copyJobItem.TargetClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_name"); ok {
		copyJobItem.SourceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_name"); ok {
		copyJobItem.TargetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_folder_id"); ok {
		copyJobItem.TargetFolderId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("job_type"); ok {
		copyJobItem.JobType = helper.IntInt64(v.(int))
	}

	request.JobItems = append(request.JobItems, &copyJobItem)

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().CopyJobs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus JobCopy not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus JobCopy failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.CopyJobsResults[0].TargetJobId
	d.SetId(jobId)

	return resourceTencentCloudOceanusJobCopyRead(d, meta)
}

func resourceTencentCloudOceanusJobCopyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job_copy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId   = d.Id()
	)

	JobCopy, err := service.DescribeOceanusJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if JobCopy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusJobCopy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if JobCopy.JobId != nil {
		_ = d.Set("job_id", JobCopy.JobId)
	}

	if JobCopy.ClusterId != nil {
		_ = d.Set("target_cluster_id", JobCopy.ClusterId)
	}

	if JobCopy.Name != nil {
		_ = d.Set("target_name", JobCopy.Name)
	}

	if JobCopy.WorkSpaceId != nil {
		_ = d.Set("work_space_id", JobCopy.WorkSpaceId)
	}

	return nil
}

func resourceTencentCloudOceanusJobCopyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job_copy.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId   = d.Id()
	)

	if err := service.DeleteOceanusJobById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
