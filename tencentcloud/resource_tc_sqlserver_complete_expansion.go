/*
Provides a resource to create a sqlserver complete_expansion

Example Usage

First, Create a basic SQL instance

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "example-sg"
  description = "desc."
}

resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [tencentcloud_security_group.security_group.id]
  project_id             = 0
  memory                 = 2
  storage                = 20
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  tags                   = {
    "createBy" = "tfExample"
  }
}
```

Expand the current instance, storage: 20->40, wait_switch = 1

```hcl
resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [tencentcloud_security_group.security_group.id]
  project_id             = 0
  memory                 = 2
  storage                = 40
  wait_switch            = 1
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  tags                   = {
    "createBy" = "tfExample"
  }
}
```

Complete the expansion task immediately

```hcl
resource "tencentcloud_sqlserver_complete_expansion" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverCompleteExpansion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverCompleteExpansionCreate,
		Read:   resourceTencentCloudSqlserverCompleteExpansionRead,
		Delete: resourceTencentCloudSqlserverCompleteExpansionDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},
		},
	}
}

func resourceTencentCloudSqlserverCompleteExpansionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewCompleteExpansionRequest()
		response    = sqlserver.NewCompleteExpansionResponse()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		instanceId  string
		flowId      int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CompleteExpansion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver complete expansion %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		response = result
		flowId = int64(*response.Response.FlowId)
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver completeExpansion failed, reason:%+v", logId, err)
		return err
	}

	flowRequest.FlowId = &flowId
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return retryError(e)
		}

		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver completeExpansion status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver completeExpansion status is fail"))
		} else {
			e = fmt.Errorf("sqlserver completeExpansion status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver completeExpansion failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverCompleteExpansionRead(d, meta)
}

func resourceTencentCloudSqlserverCompleteExpansionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverCompleteExpansionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
