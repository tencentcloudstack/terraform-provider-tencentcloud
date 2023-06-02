/*
Use this data source to query detailed information of vpc cvm_instances

Example Usage

```hcl
data "tencentcloud_vpc_cvm_instances" "cvm_instances" {
  filters {
    name   = "vpc-id"
    values = ["vpc-lh4nqig9"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcCvmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcCvmInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Filter condition. `RouteTableIds` and `Filters` cannot be specified at the same time. vpc-id - String - (Filter condition) VPC instance ID, such as `vpc-f49l6u0z`;instance-type - String - (Filter condition) CVM instance ID;instance-name - String - (Filter condition) CVM name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The attribute name. If more than one Filter exists, the logical relation between these Filters is `AND`.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`. For a `bool` parameter, the valid values include `TRUE` and `FALSE`.",
						},
					},
				},
			},

			"instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of CVM instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC instance ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet instance ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM Name.",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM status.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CPU cores in an instance (in core).",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance's memory capacity. Unit: GB.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type.",
						},
						"eni_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance ENI quota (including primary ENIs).",
						},
						"eni_ip_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Private IP quoata for instance ENIs (including primary ENIs).",
						},
						"instance_eni_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of ENIs (including primary ENIs) bound to a instance.",
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

func dataSourceTencentCloudVpcCvmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_cvm_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceSet []*vpc.CvmInstance

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCvmInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceSet))

	if instanceSet != nil {
		for _, cvmInstance := range instanceSet {
			cvmInstanceMap := map[string]interface{}{}

			if cvmInstance.VpcId != nil {
				cvmInstanceMap["vpc_id"] = cvmInstance.VpcId
			}

			if cvmInstance.SubnetId != nil {
				cvmInstanceMap["subnet_id"] = cvmInstance.SubnetId
			}

			if cvmInstance.InstanceId != nil {
				cvmInstanceMap["instance_id"] = cvmInstance.InstanceId
			}

			if cvmInstance.InstanceName != nil {
				cvmInstanceMap["instance_name"] = cvmInstance.InstanceName
			}

			if cvmInstance.InstanceState != nil {
				cvmInstanceMap["instance_state"] = cvmInstance.InstanceState
			}

			if cvmInstance.CPU != nil {
				cvmInstanceMap["cpu"] = cvmInstance.CPU
			}

			if cvmInstance.Memory != nil {
				cvmInstanceMap["memory"] = cvmInstance.Memory
			}

			if cvmInstance.CreatedTime != nil {
				cvmInstanceMap["created_time"] = cvmInstance.CreatedTime
			}

			if cvmInstance.InstanceType != nil {
				cvmInstanceMap["instance_type"] = cvmInstance.InstanceType
			}

			if cvmInstance.EniLimit != nil {
				cvmInstanceMap["eni_limit"] = cvmInstance.EniLimit
			}

			if cvmInstance.EniIpLimit != nil {
				cvmInstanceMap["eni_ip_limit"] = cvmInstance.EniIpLimit
			}

			if cvmInstance.InstanceEniCount != nil {
				cvmInstanceMap["instance_eni_count"] = cvmInstance.InstanceEniCount
			}

			ids = append(ids, *cvmInstance.InstanceId)
			tmpList = append(tmpList, cvmInstanceMap)
		}

		_ = d.Set("instance_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
