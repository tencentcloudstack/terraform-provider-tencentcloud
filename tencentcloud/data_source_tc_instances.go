/*
Use this data source to query cvm instances.

Example Usage

```hcl
data "tencentcloud_instances" "foo" {
  instance_id = "ins-da412f5a"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instances to be queried.",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the instances to be queried.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CVM instance locates at.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The project CVM belongs to.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vpc to be queried.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a vpc subnetwork.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cvm instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instances.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the instances.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the instance.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CPU cores of the instance.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity, unit in GB.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CVM instance locates at.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project CVM belongs to.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the image.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of the instance.",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the system disk.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the system disk.",
						},
						"system_disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID of the system disk.",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of data disk. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the data disk.",
									},
									"data_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size of the data disk.",
									},
									"data_disk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image ID of the data disk.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the data disk is destroyed with the instance.",
									},
								},
							},
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of a vpc subnetwork.",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of the instance.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public network maximum output bandwidth of the instance.",
						},
						"allocate_public_ip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether public ip is assigned.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the instance.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public ip of the instance.",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private ip of the instance.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Security groups of the instance.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the instance.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expired time of the instance.",
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The way that CVM instance will be renew automatically or not when it reach the end of the prepaid tenancy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_instances.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("instance_id"); ok {
		filter["instance-id"] = v.(string)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		filter["instance-name"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		filter["project-id"] = fmt.Sprintf("%d", v.(int))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		filter["vpc-id"] = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		filter["subnet-id"] = v.(string)
	}

	var instances []*cvm.Instance
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeInstanceByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, InternalError)
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
			"instance_id":                instance.InstanceId,
			"instance_name":              instance.InstanceName,
			"instance_type":              instance.InstanceType,
			"cpu":                        instance.CPU,
			"memory":                     instance.Memory,
			"availability_zone":          instance.Placement.Zone,
			"project_id":                 instance.Placement.ProjectId,
			"image_id":                   instance.ImageId,
			"instance_charge_type":       instance.InstanceChargeType,
			"system_disk_type":           instance.SystemDisk.DiskType,
			"system_disk_size":           instance.SystemDisk.DiskSize,
			"system_disk_id":             instance.SystemDisk.DiskId,
			"vpc_id":                     instance.VirtualPrivateCloud.VpcId,
			"subnet_id":                  instance.VirtualPrivateCloud.SubnetId,
			"internet_charge_type":       instance.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": instance.InternetAccessible.InternetMaxBandwidthOut,
			"allocate_public_ip":         instance.InternetAccessible.PublicIpAssigned,
			"status":                     instance.InstanceState,
			"security_groups":            helper.StringsInterfaces(instance.SecurityGroupIds),
			"tags":                       flattenCvmTagsMapping(instance.Tags),
			"create_time":                instance.CreatedTime,
			"expired_time":               instance.ExpiredTime,
			"instance_charge_type_prepaid_renew_flag": instance.RenewFlag,
		}
		if len(instance.PublicIpAddresses) > 0 {
			mapping["public_ip"] = *instance.PublicIpAddresses[0]
		}
		if len(instance.PrivateIpAddresses) > 0 {
			mapping["private_ip"] = *instance.PrivateIpAddresses[0]
		}
		dataDisks := make([]map[string]interface{}, 0, len(instance.DataDisks))
		for _, v := range instance.DataDisks {
			dataDisk := map[string]interface{}{
				"data_disk_type":       v.DiskType,
				"data_disk_size":       v.DiskSize,
				"data_disk_id":         v.DiskId,
				"delete_with_instance": v.DeleteWithInstance,
			}
			dataDisks = append(dataDisks, dataDisk)
		}
		mapping["data_disks"] = dataDisks
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
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
