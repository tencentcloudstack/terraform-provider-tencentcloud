/*
Provides a resource to create a mysql db_import_job_operation

Example Usage

```hcl
resource "tencentcloud_mysql_db_import_job_operation" "db_import_job_operation" {
  instance_id = "cdb-c1nl9rpv"
  user = "root"
  file_name = "mysql.sql"
  password = "ABCabc123"
  db_name = "t_test"
  cos_url = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mysql/mysql.sql?q-sign-algorithm=sha1&q-ak=AKIDRnMWiUNr14F29GvCwOSHu9l_FdCdORqAxblrE10nDaO6mVI701oXTe-gL1QpClgW&q-sign-time=1684921483;1684925083&q-key-time=1684921483;1684925083&q-header-list=host&q-url-param-list=&q-signature=7410be4ef93075aebca459af4e617f8bcaa36f48&x-cos-security-token=EzDm9S6aRDwBLQcaxUNfb0TA30PqhOTa7d82a06a36e94b66bdbc6d09064a397bZypr0mD3oVkbJR9bRYix6BSDVYncX3Y2VCGYK6V2jFWZqIuEHoWJCe-2pDvJDNbMjF3ttWfLMqEouOkxNk28ay9NPHtMXrJgEEMb95BMAhGwi38oA2LjYfQRkk7AHesg2toSf11hiTAjVv-alf5uEidWGnFKe_6BgmnADYvtPptgXHNtsUZCxc33PF6tGBqX"
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlDbImportJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDbImportJobOperationCreate,
		Read:   resourceTencentCloudMysqlDbImportJobOperationRead,
		Delete: resourceTencentCloudMysqlDbImportJobOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The username of the cloud database.",
			},

			"file_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "file name. This file refers to the file that the user has uploaded to Tencent Cloud, and only .sql files are supported.",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The password of the user account of the cloud database instance.",
			},

			"db_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the imported target database, if it is not passed, it means that no database is specified.",
			},

			"cos_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the imported target database, if it is not passed, it means that no database is specified.",
			},

			"async_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The request ID of the asynchronous task.",
			},
		},
	}
}

func resourceTencentCloudMysqlDbImportJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_db_import_job_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mysql.NewCreateDBImportJobRequest()
		response   = mysql.NewCreateDBImportJobResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_name"); ok {
		request.FileName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_url"); ok {
		request.CosUrl = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateDBImportJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql dbImportJob failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create dbImportJob  status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create dbImportJob status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dbImportJob fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(instanceId + FILED_SP + asyncRequestId)

	return resourceTencentCloudMysqlDbImportJobOperationRead(d, meta)
}

func resourceTencentCloudMysqlDbImportJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_db_import_job_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	asyncRequestId := idSplit[1]

	dbImportJob, err := service.DescribeMysqlDbImportJobById(ctx, instanceId, asyncRequestId)
	if err != nil {
		return err
	}

	if dbImportJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlDbImportJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbImportJob.InstanceId != nil {
		_ = d.Set("instance_id", dbImportJob.InstanceId)
	}

	if dbImportJob.FileName != nil {
		_ = d.Set("file_name", dbImportJob.FileName)
	}

	if dbImportJob.DbName != nil {
		_ = d.Set("db_name", dbImportJob.DbName)
	}

	if dbImportJob.AsyncRequestId != nil {
		_ = d.Set("async_request_id", dbImportJob.AsyncRequestId)
	}

	return nil
}

func resourceTencentCloudMysqlDbImportJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_db_import_job_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
