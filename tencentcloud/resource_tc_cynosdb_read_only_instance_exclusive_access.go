/*
Provides a resource to create a cynosdb read_only_instance_exclusive_access

Example Usage

```hcl
resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id = "cynosdbmysql-12345678"
  instance_id = "cynosdbmysql-ins-12345678"
  vpc_id = "vpc-12345678"
  subnet_id = "subnet-12345678"
  port = 1234
  security_group_ids =
}
```

Import

cynosdb read_only_instance_exclusive_access can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access read_only_instance_exclusive_access_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessCreate,
		Read:   resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessRead,
		Delete: resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessDelete,
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

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Need to activate a read-only instance ID with unique access.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Specified VPC ID.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The specified subnet ID.",
			},

			"port": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Port.",
			},

			"security_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security Group.",
			},
		},
	}
}

func resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_read_only_instance_exclusive_access.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewOpenReadOnlyInstanceExclusiveAccessRequest()
		response  = cynosdb.NewOpenReadOnlyInstanceExclusiveAccessResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("port"); v != nil {
		request.Port = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().OpenReadOnlyInstanceExclusiveAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb readOnlyInstanceExclusiveAccess failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbReadOnlyInstanceExclusiveAccessStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessRead(d, meta)
}

func resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_read_only_instance_exclusive_access.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_read_only_instance_exclusive_access.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
