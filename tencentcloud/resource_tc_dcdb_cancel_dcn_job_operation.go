/*
Provides a resource to create a dcdb cancel_dcn_job_operation

Example Usage

```hcl
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance_dcn" {
	instance_name     = "test_dcdb_db_hourdb_instance_dcn"
	zones             = [var.default_az]
	shard_memory      = "2"
	shard_storage     = "10"
	shard_node_count  = "2"
	shard_count       = "2"
	vpc_id            = local.vpc_id
	subnet_id         = local.subnet_id
	security_group_id = local.sg_id
	db_version_id     = "8.0"
	dcn_region        = "ap-guangzhou"
	dcn_instance_id   = local.dcdb_id  //master_instance
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
}

locals {
	dcn_dcdb_id = tencentcloud_dcdb_hourdb_instance.hourdb_instance_dcn.id
}

resource "tencentcloud_dcdb_cancel_dcn_job_operation" "cancel_operation" {
	instance_id = local.dcn_dcdb_id // cancel the dcn/readonly instance
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbCancelDcnJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbCancelDcnJobOperationCreate,
		Read:   resourceTencentCloudDcdbCancelDcnJobOperationRead,
		Delete: resourceTencentCloudDcdbCancelDcnJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudDcdbCancelDcnJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCancelDcnJobRequest()
		instanceId string
		flowId     *int64
		service    = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CancelDcnJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb cancelDcnJobOperation failed, reason:%+v", logId, err)
		return err
	}

	// need to wait flow success
	// 0:success; 1:failed, 2:running
	conf := BuildStateChangeConf([]string{}, []string{"0"}, 3*readRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbCancelDcnJobOperationRead(d, meta)
}

func resourceTencentCloudDcdbCancelDcnJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbCancelDcnJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
