/*
Provides a resource to create a cynosdb isolate_instance

Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  account_name         = "terraform_test"
  account_password     = "Password@1234"
  host                 = "%"
  description          = "testx"
  max_user_connections = 2
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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
)

func resourceTencentCloudCynosdbIsolateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbIsolateInstanceCreate,
		Read:   resourceTencentCloudCynosdbIsolateInstanceRead,
		Update: resourceTencentCloudCynosdbIsolateInstanceUpdate,
		Delete: resourceTencentCloudCynosdbIsolateInstanceDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "isolate, activate.",
			},
		},
	}
}

func resourceTencentCloudCynosdbIsolateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_isolate_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		clusterId  string
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(clusterId + FILED_SP + instanceId)

	return resourceTencentCloudCynosdbIsolateInstanceUpdate(d, meta)
}

func resourceTencentCloudCynosdbIsolateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_isolate_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbIsolateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_isolate_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	var flowId int64
	if operate == "isolate" {
		request := cynosdb.NewIsolateInstanceRequest()
		request.ClusterId = &clusterId
		request.InstanceIdList = []*string{&instanceId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().IsolateInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = *result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s isolate cynosdb instance failed, reason:%+v", logId, err)
			return err
		}
	} else if operate == "activate" {
		request := cynosdb.NewActivateInstanceRequest()
		request.ClusterId = &clusterId
		request.InstanceIdList = []*string{&instanceId}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ActivateInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = *result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s activate cynosdb instance failed, reason:%+v", logId, err)
			return err
		}
	} else {
		return fmt.Errorf("[CRITAL]%s Operation type error", logId)
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("isolate or activate cynosdb instance is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s isolate or activate cynosdb instance fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbIsolateInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbIsolateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_isolate_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
