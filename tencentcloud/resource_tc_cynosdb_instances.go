/*
Provides a resource to create a cynosdb instances

Example Usage

```hcl
resource "tencentcloud_cynosdb_instances" "instances" {
  cluster_id = "cynosdbmysql-ins-bzkxxrmt"
  instance_id_list =
}
```

Import

cynosdb instances can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_instances.instances instances_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstancesCreate,
		Read:   resourceTencentCloudCynosdbInstancesRead,
		Delete: resourceTencentCloudCynosdbInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id_list": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID array.",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewIsolateInstanceRequest()
		response   = cynosdb.NewIsolateInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id_list"); ok {
		instanceIdListSet := v.(*schema.Set).List()
		for i := range instanceIdListSet {
			instanceIdList := instanceIdListSet[i].(string)
			request.InstanceIdList = append(request.InstanceIdList, &instanceIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().IsolateInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb instances failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbInstancesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbInstancesRead(d, meta)
}

func resourceTencentCloudCynosdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instancesId := d.Id()

	instances, err := service.DescribeCynosdbInstancesById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instances == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbInstances` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instances.ClusterId != nil {
		_ = d.Set("cluster_id", instances.ClusterId)
	}

	if instances.InstanceIdList != nil {
		_ = d.Set("instance_id_list", instances.InstanceIdList)
	}

	return nil
}

func resourceTencentCloudCynosdbInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instances.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instancesId := d.Id()

	if err := service.DeleteCynosdbInstancesById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
