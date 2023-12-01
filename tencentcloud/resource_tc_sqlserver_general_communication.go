/*
Provides a resource to create a sqlserver general_communication

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_general_communication" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}
```

Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.example mssql-hlh6yka1
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverGeneralCommunication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCommunicationCreate,
		Read:   resourceTencentCloudSqlserverGeneralCommunicationRead,
		Delete: resourceTencentCloudSqlserverGeneralCommunicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instances.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCommunicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewOpenInterCommunicationRequest()
		response    = sqlserver.NewOpenInterCommunicationResponse()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIdSet = append(request.InstanceIdSet, helper.String(v.(string)))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().OpenInterCommunication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver generalCommunication not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *response.Response.InterInstanceFlowSet[0].InstanceId
	flowId := *response.Response.InterInstanceFlowSet[0].FlowId
	flowRequest.FlowId = &flowId
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return retryError(e)
		}
		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver generalCommunication status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver bgeneralCommunication status is fail"))
		} else {
			e = fmt.Errorf("sqlserver generalCommunication status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCommunicationRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCommunicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	generalCommunication, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCommunication == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCommunication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCommunication.InstanceId != nil {
		_ = d.Set("instance_id", generalCommunication.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCommunicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		instanceId  = d.Id()
	)

	flowId, err := service.DeleteSqlserverGeneralCommunicationById(ctx, instanceId)
	if err != nil {
		log.Printf("[CRITAL]%s delete sqlserver generalCommunication failed, reason:%+v", logId, err)
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
			return resource.RetryableError(fmt.Errorf("sqlserver generalCommunication status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver bgeneralCommunication status is fail"))
		} else {
			e = fmt.Errorf("sqlserver generalCommunication status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete sqlserver generalCommunication status failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
