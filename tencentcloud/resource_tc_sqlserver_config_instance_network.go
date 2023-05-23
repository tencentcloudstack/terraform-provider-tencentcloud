/*
Provides a resource to create a sqlserver config_instance_network

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_network" "config_instance_network" {
  instance_id = "mssql-qelbzgwf"
  new_vpc_id = "vpc-4owdpnwr"
  new_subnet_id = "sub-ahv6swf2"
  vip = "172.16.16.48"
}
```

Import

sqlserver config_instance_network can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_network.config_instance_network config_instance_network_id
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
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverConfigInstanceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceNetworkCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceNetworkRead,
		Update: resourceTencentCloudSqlserverConfigInstanceNetworkUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"new_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the new VPC.",
			},
			"new_subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the new subnet.",
			},
			"vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "New VIP.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
		vip        string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("vip"); ok {
		vip = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, vip}, FILED_SP))

	return resourceTencentCloudSqlserverConfigInstanceNetworkUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	vip := idSplit[1]

	configInstanceNetwork, err := service.DescribeSqlserverConfigInstanceNetworkById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceNetwork == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceNetwork` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configInstanceNetwork.InstanceId != nil {
		_ = d.Set("instance_id", configInstanceNetwork.InstanceId)
	}

	if configInstanceNetwork.UniqVpcId != nil {
		_ = d.Set("new_vpc_id", configInstanceNetwork.UniqVpcId)
	}

	if configInstanceNetwork.UniqSubnetId != nil {
		_ = d.Set("new_subnet_id", configInstanceNetwork.UniqSubnetId)
	}

	if vip != "" {
		if configInstanceNetwork.Vip != nil {
			_ = d.Set("vip", configInstanceNetwork.Vip)
		}
	} else {
		_ = d.Set("vip", vip)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewModifyDBInstanceNetworkRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	vip := idSplit[1]

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("new_vpc_id"); ok {
		request.NewVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("new_subnet_id"); ok {
		request.NewSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("old_ip_retain_time"); ok {
		request.OldIpRetainTime = helper.IntInt64(v.(int))
	}

	if vip != "" {
		request.Vip = &vip
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceNetwork(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver configInstanceNetwork not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
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
			return resource.RetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is fail"))
		} else {
			e = fmt.Errorf("sqlserver configInstanceNetwork status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceNetworkRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
