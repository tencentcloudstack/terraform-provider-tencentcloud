/*
Use this data source to query detailed information of tdcpg instances

Example Usage

```hcl
data "tencentcloud_tdcpg_instances" "instances" {
  cluster_id = ""
  instance_id = ""
  instance_name = ""
  status = ""
  instance_type = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdcpgInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdcpgInstancesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance id.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance name.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance status.",
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance type.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance name.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster id.",
						},
						"endpoint_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "endpoint id.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "zone.",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db version.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "pay mode.",
						},
						"pay_period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "pay period expired time.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "cpu cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "memory size, unit is GiB.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance type.",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db major version.",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db kernel version.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdcpgInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdcpg_instances.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId *string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["instance_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		paramMap["instance_type"] = helper.String(v.(string))
	}

	tdcpgService := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceSet []*tdcpg.Instance
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := tdcpgService.DescribeTdcpgInstancesByFilter(ctx, clusterId, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s read Tdcpg instanceSet failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(instanceSet))
	instanceList := make([]map[string]interface{}, 0, len(instanceSet))
	if instanceSet != nil {
		for _, instance := range instanceSet {
			instanceSetMap := map[string]interface{}{}
			if instance.InstanceId != nil {
				instanceSetMap["instance_id"] = instance.InstanceId
			}
			if instance.InstanceName != nil {
				instanceSetMap["instance_name"] = instance.InstanceName
			}
			if instance.ClusterId != nil {
				instanceSetMap["cluster_id"] = instance.ClusterId
			}
			if instance.EndpointId != nil {
				instanceSetMap["endpoint_id"] = instance.EndpointId
			}
			if instance.Region != nil {
				instanceSetMap["region"] = instance.Region
			}
			if instance.Zone != nil {
				instanceSetMap["zone"] = instance.Zone
			}
			if instance.DBVersion != nil {
				instanceSetMap["db_version"] = instance.DBVersion
			}
			if instance.Status != nil {
				instanceSetMap["status"] = instance.Status
			}
			if instance.StatusDesc != nil {
				instanceSetMap["status_desc"] = instance.StatusDesc
			}
			if instance.CreateTime != nil {
				instanceSetMap["create_time"] = instance.CreateTime
			}
			if instance.PayMode != nil {
				instanceSetMap["pay_mode"] = instance.PayMode
			}
			if instance.PayPeriodEndTime != nil {
				instanceSetMap["pay_period_end_time"] = instance.PayPeriodEndTime
			}
			if instance.CPU != nil {
				instanceSetMap["cpu"] = instance.CPU
			}
			if instance.Memory != nil {
				instanceSetMap["memory"] = instance.Memory
			}
			if instance.InstanceType != nil {
				instanceSetMap["instance_type"] = instance.InstanceType
			}
			if instance.DBMajorVersion != nil {
				instanceSetMap["db_major_version"] = instance.DBMajorVersion
			}
			if instance.DBKernelVersion != nil {
				instanceSetMap["db_kernel_version"] = instance.DBKernelVersion
			}
			ids = append(ids, *instance.ClusterId+FILED_SP+*instance.InstanceId)
			instanceList = append(instanceList, instanceSetMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		err := d.Set("list", instanceList)
		if err != nil {
			log.Printf("[CRITICAL]%s set tdcpg instanceList failed, reason:%+v", logId, err)
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
