/*
Use this data source to query CDH instances.

Example Usage

```hcl
data "tencentcloud_cdh_instances" "list" {
  availability_zone = "ap-guangzhou-3"
  host_id = "host-d6s7i5q4"
  host_name = "test"
  host_state = "RUNNING"
  project_id = 1154137
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdhInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdhInstancesRead,

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CDH instances to be queried.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CDH instances to be queried.",
			},
			"host_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State of the CDH instances to be queried. Valid values: `PENDING`, `LAUNCH_FAILURE`, `RUNNING`, `EXPIRED`.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CDH instance locates at.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The project CDH belongs to.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"cdh_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cdh instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CDH instance.",
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CDH instance.",
						},
						"host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the CDH instance.",
						},
						"host_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the CDH instance.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of the CDH instance.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CDH instance locates at.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project CDH belongs to.",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto renewal flag.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the CDH instance.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expired time of the CDH instance.",
						},
						"cvm_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Id of CVM instances that have been created on the CDH instance.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cage ID of the CDH instance. This parameter is only valid for CDH instances in the cages of finance availability zones.",
						},
						"host_resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of host resource. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_total_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of total CPU cores of the instance.",
									},
									"cpu_available_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of available CPU cores of the instance.",
									},
									"memory_total_size": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Instance memory total capacity, unit in GiB.",
									},
									"memory_available_size": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Instance memory available capacity, unit in GiB.",
									},
									"disk_total_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance disk total capacity, unit in GiB.",
									},
									"disk_available_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance disk available capacity, unit in GiB.",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the disk.",
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

func dataSourceTencentCloudCdhInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdh_instances.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cdhService := CdhService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("host_id"); ok {
		filter["host-id"] = v.(string)
	}
	if v, ok := d.GetOk("host_name"); ok {
		filter["host-name"] = v.(string)
	}
	if v, ok := d.GetOk("host_state"); ok {
		filter["host-state"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		filter["project-id"] = strconv.FormatInt(int64(v.(int)), 10)
	}

	var instances []*cvm.HostItem
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = cdhService.DescribeCdhInstanceByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"host_id":            instance.HostId,
			"host_name":          instance.HostName,
			"host_type":          instance.HostType,
			"host_state":         instance.HostState,
			"charge_type":        instance.HostChargeType,
			"availability_zone":  instance.Placement.Zone,
			"project_id":         instance.Placement.ProjectId,
			"prepaid_renew_flag": instance.RenewFlag,
			"create_time":        instance.CreatedTime,
			"expired_time":       instance.ExpiredTime,
			"cvm_instance_ids":   helper.StringsInterfaces(instance.InstanceIds),
			"cage_id":            instance.CageId,
		}
		hostResource := map[string]interface{}{
			"cpu_total_num":         instance.HostResource.CpuTotal,
			"cpu_available_num":     instance.HostResource.CpuAvailable,
			"memory_total_size":     instance.HostResource.MemTotal,
			"memory_available_size": instance.HostResource.MemAvailable,
			"disk_total_size":       instance.HostResource.DiskTotal,
			"disk_available_size":   instance.HostResource.DiskAvailable,
			"disk_type":             instance.HostResource.DiskType,
		}
		mapping["host_resource"] = []map[string]interface{}{hostResource}
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.HostId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("cdh_instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cdh instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			return err
		}
	}
	return nil
}
