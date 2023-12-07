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

func resourceTencentCloudOceanusJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusJobCreate,
		Read:   resourceTencentCloudOceanusJobRead,
		Update: resourceTencentCloudOceanusJobUpdate,
		Delete: resourceTencentCloudOceanusJobDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of the job. It can be composed of Chinese, English, numbers, hyphens (-), underscores (_), and periods (.), and the length cannot exceed 50 characters. Note that the job name cannot be the same as an existing job.",
			},
			"job_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(JOB_TYPE),
				Description:  "The type of the job. 1 indicates SQL job, and 2 indicates JAR job.",
			},
			"cluster_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(CLUSTER_TYPE),
				Description:  "The type of the cluster. 1 indicates shared cluster, and 2 indicates exclusive cluster.",
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "When ClusterType=2, it is required to specify the ID of the exclusive cluster to which the job is submitted.",
			},
			"cu_mem": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      CU_MEM_4,
				ValidateFunc: validateAllowedIntValue(CU_MEM),
				Description:  "Set the memory specification of each CU, in GB. It supports 2, 4, 8, and 16 (which needs to apply for the whitelist before use). The default is 4, that is, 1 CU corresponds to 4 GB of running memory.",
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The remark information of the job. It can be set arbitrarily.",
			},
			"folder_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The folder ID to which the job name belongs. The root directory is root.",
			},
			"flink_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The Flink version that the job runs.",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The workspace SerialId.",
			},
		},
	}
}

func resourceTencentCloudOceanusJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = oceanus.NewCreateJobRequest()
		response = oceanus.NewCreateJobResponse()
		jobId    string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("job_type"); ok {
		request.JobType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("cluster_type"); ok {
		request.ClusterType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cu_mem"); ok {
		request.CuMem = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("folder_id"); ok {
		request.FolderId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flink_version"); ok {
		request.FlinkVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().CreateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus Job not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus Job failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudOceanusJobRead(d, meta)
}

func resourceTencentCloudOceanusJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId   = d.Id()
	)

	Job, err := service.DescribeOceanusJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if Job == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Job.Name != nil {
		_ = d.Set("name", Job.Name)
	}

	if Job.JobType != nil {
		_ = d.Set("job_type", Job.JobType)
	}

	if Job.ClusterId != nil {
		_ = d.Set("cluster_id", Job.ClusterId)
	}

	if Job.CuMem != nil {
		_ = d.Set("cu_mem", Job.CuMem)
	}

	if Job.Remark != nil {
		_ = d.Set("remark", Job.Remark)
	}

	if Job.FlinkVersion != nil {
		_ = d.Set("flink_version", Job.FlinkVersion)
	}

	if Job.WorkSpaceId != nil {
		_ = d.Set("work_space_id", Job.WorkSpaceId)
	}

	return nil
}

func resourceTencentCloudOceanusJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = oceanus.NewModifyJobRequest()
		jobId   = d.Id()
	)

	immutableArgs := []string{"job_type", "cluster_type", "cluster_id", "cu_mem", "folder_id", "flink_version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.JobId = &jobId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("work_space_id") {
		if v, ok := d.GetOk("work_space_id"); ok {
			request.WorkSpaceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().ModifyJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update oceanus Job failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOceanusJobRead(d, meta)
}

func resourceTencentCloudOceanusJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_job.delete")()
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
