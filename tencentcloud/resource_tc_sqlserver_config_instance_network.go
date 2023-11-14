/*
Provides a resource to create a sqlserver config_instance_network

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_network" "config_instance_network" {
  instance_id = "mssql-i1z41iwd"
  new_vpc_id = "vpc-j90ok"
  new_subnet_id = "sub-ja891"
  old_ip_retain_time = 0
  vip = "10.1.200.11"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
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

			"old_ip_retain_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Retention period (in hours) of the original VIP. Value range: 0-168. Default value: 0, indicating the original VIP is released immediately.",
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

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigInstanceNetworkUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configInstanceNetworkId := d.Id()

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

	if configInstanceNetwork.NewVpcId != nil {
		_ = d.Set("new_vpc_id", configInstanceNetwork.NewVpcId)
	}

	if configInstanceNetwork.NewSubnetId != nil {
		_ = d.Set("new_subnet_id", configInstanceNetwork.NewSubnetId)
	}

	if configInstanceNetwork.OldIpRetainTime != nil {
		_ = d.Set("old_ip_retain_time", configInstanceNetwork.OldIpRetainTime)
	}

	if configInstanceNetwork.Vip != nil {
		_ = d.Set("vip", configInstanceNetwork.Vip)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDBInstanceNetworkRequest()

	configInstanceNetworkId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "new_vpc_id", "new_subnet_id", "old_ip_retain_time", "vip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceNetwork(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceNetworkRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_network.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
