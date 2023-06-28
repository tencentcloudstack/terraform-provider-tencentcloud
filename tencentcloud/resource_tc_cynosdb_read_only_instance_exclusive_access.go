/*
Provides a resource to create a cynosdb read_only_instance_exclusive_access

Example Usage

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}
variable "cynosdb_cluster_instance_id" {
  default = "default_cluster_instance"
}

variable "cynosdb_cluster_security_group_id" {
  default = "default_security_group_id"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id         = var.cynosdb_cluster_id
  instance_id        = var.cynosdb_cluster_instance_id
  vpc_id             = local.vpc_id
  subnet_id          = local.subnet_id
  port               = 1234
  security_group_ids = [var.cynosdb_cluster_security_group_id]
}
```
*/
package tencentcloud

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessCreate,
		Read:   resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessRead,
		Delete: resourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccessDelete,
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
				Description: "port.",
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
		request    = cynosdb.NewOpenReadOnlyInstanceExclusiveAccessRequest()
		response   = cynosdb.NewOpenReadOnlyInstanceExclusiveAccessResponse()
		flowId     *int64
		clusterId  string
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
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
			if securityGroupIdsSet[i] != nil {
				securityGroupIds := securityGroupIdsSet[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().OpenReadOnlyInstanceExclusiveAccess(request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				// repeat to execute this cmd can be ignored
				if sdkErr.Code == "FailedOperation.OperationFailedError" {
					return nil
				}
			}
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

	if response.Response == nil || response.Response.FlowId == nil {
		log.Printf("[CRITAL]%s FlowId is null. Ingnore this operation.", logId)
	} else {
		flowId = response.Response.FlowId

		service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		conf := BuildStateChangeConf([]string{}, []string{CYNOSDB_FLOW_STATUS_SUCCESSFUL}, 10*readRetryTimeout, time.Second, service.CynosdbClusterSlaveZoneStateRefreshFunc(*flowId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(strings.Join([]string{clusterId, instanceId}, FILED_SP))

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
