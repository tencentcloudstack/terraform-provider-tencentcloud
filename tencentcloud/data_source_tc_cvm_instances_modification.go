/*
Use this data source to query cvm instances modification.

Example Usage

```hcl
data "tencentcloud_cvm_instances_modification" "foo" {
  instance_ids = ["ins-xxxxxxx"]
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmInstancesModification() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmInstancesModificationRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "One or more instance ID to be queried. It can be obtained from the InstanceId in the returned value of API DescribeInstances. The maximum number of instances in batch for each request is 20.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The upper limit of Filters for each request is 10 and the upper limit for Filter.Values is 2.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Value of the field.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"instance_type_config_status_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of model configurations that can be adjusted by the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State description.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status description information.",
						},
						"instance_type_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type.",
									},
									"instance_family": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance family.",
									},
									"gpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of GPU kernels, in cores.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of CPU kernels, in cores.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory capacity (in GB).",
									},
									"fpga": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of FPGA kernels, in cores.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCvmInstancesModificationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_instances_modification.read")()
	logId := getLogId(contextNil)

	var (
		request  = cvm.NewDescribeInstancesModificationRequest()
		response = cvm.NewDescribeInstancesModificationResponse()
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		request.InstanceIds = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("filters"); ok {
		filters := make([]*cvm.Filter, 0)
		for _, item := range v.(*schema.Set).List() {
			filter := item.(map[string]interface{})
			name := filter["name"].(string)
			filters = append(filters, &cvm.Filter{
				Name:   &name,
				Values: helper.StringsStringsPoint(filter["values"].([]string)),
			})
		}
		request.Filters = filters
	}

	instanceTypeConfigStatusList := make([]map[string]interface{}, 0)

	var innerErr error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, innerErr = meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().DescribeInstancesModification(request)
		if innerErr != nil {
			return retryError(innerErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	for _, instanceTypeConfigStatusSetItem := range response.Response.InstanceTypeConfigStatusSet {
		instanceTypeConfigStatus := make(map[string]interface{})
		instanceTypeConfigStatus["status"] = instanceTypeConfigStatusSetItem.Status
		instanceTypeConfigStatus["message"] = instanceTypeConfigStatusSetItem.Message

		instanceTypeConfigMaps := make([]map[string]interface{}, 0)
		instanceTypeConfigMap := make(map[string]interface{})
		instanceTypeConfig := instanceTypeConfigStatusSetItem.InstanceTypeConfig
		instanceTypeConfigMap["zone"] = instanceTypeConfig.Zone
		ids = append(ids, *instanceTypeConfig.InstanceType)
		instanceTypeConfigMap["instance_type"] = instanceTypeConfig.InstanceType
		instanceTypeConfigMap["instance_family"] = instanceTypeConfig.InstanceFamily
		instanceTypeConfigMap["gpu"] = instanceTypeConfig.GPU
		instanceTypeConfigMap["cpu"] = instanceTypeConfig.CPU
		instanceTypeConfigMap["memory"] = instanceTypeConfig.Memory
		instanceTypeConfigMap["fpga"] = instanceTypeConfig.FPGA
		instanceTypeConfigMaps = append(instanceTypeConfigMaps, instanceTypeConfigMap)
		instanceTypeConfigStatus["instance_type_config"] = instanceTypeConfigMaps

		instanceTypeConfigStatusList = append(instanceTypeConfigStatusList, instanceTypeConfigStatus)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_type_config_status_list", instanceTypeConfigStatusList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceTypeConfigStatusList); err != nil {
			return err
		}
	}
	return nil

}
