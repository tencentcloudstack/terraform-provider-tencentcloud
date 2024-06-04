package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstancesRead,
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CVM instance locates at.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instances to be queried.",
			},

			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cvm instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocate_public_ip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether public ip is assigned.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CVM instance locates at.",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CAM role name authorized to access.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CPU cores of the instance.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance os name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the instance.",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of data disk. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_disk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image ID of the data disk.",
									},
									"data_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size of the data disk.",
									},
									"data_disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the data disk.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the data disk is destroyed with the instance.",
									},
								},
							},
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expired time of the instance.",
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
						"instance_charge_type_prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The way that CVM instance will be renew automatically or not when it reach the end of the prepaid tenancy.",
						},
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
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity, unit in GB.",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP of the instance.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project CVM belongs to.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP of the instance.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Security groups of the instance.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the instance.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of a vpc subnetwork.",
						},
						"system_disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID of the system disk.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the system disk.",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the system disk.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the instance.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc.",
						},
					},
				},
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the instances to be queried.",
			},

			"instance_set_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      100,
				ConflictsWith: []string{"instance_id", "instance_name", "availability_zone", "project_id", "vpc_id", "subnet_id", "tags"},
				Description:   "Instance set ids, max length is 100, conflict with other field.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The project CVM belongs to.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a vpc subnetwork.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the instance.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vpc to be queried.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var filtersList []*cvm.Filter
	filtersMap := map[string]*cvm.Filter{}
	filter := cvm.Filter{}
	name := "instance-id"
	filter.Name = &name
	if v, ok := d.GetOk("instance_id"); ok {
		filter.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp0"] = &filter
	if v, ok := filtersMap["Temp0"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter2 := cvm.Filter{}
	name2 := "instance-name"
	filter2.Name = &name2
	if v, ok := d.GetOk("instance_name"); ok {
		filter2.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp1"] = &filter2
	if v, ok := filtersMap["Temp1"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter3 := cvm.Filter{}
	name3 := "zone"
	filter3.Name = &name3
	if v, ok := d.GetOk("availability_zone"); ok {
		filter3.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp2"] = &filter3
	if v, ok := filtersMap["Temp2"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter4 := cvm.Filter{}
	name4 := "project-id"
	filter4.Name = &name4
	if v, ok := d.GetOk("project_id"); ok {
		filter4.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp3"] = &filter4
	if v, ok := filtersMap["Temp3"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter5 := cvm.Filter{}
	name5 := "vpc-id"
	filter5.Name = &name5
	if v, ok := d.GetOk("vpc_id"); ok {
		filter5.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp4"] = &filter5
	if v, ok := filtersMap["Temp4"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter6 := cvm.Filter{}
	name6 := "subnet-id"
	filter6.Name = &name6
	if v, ok := d.GetOk("subnet_id"); ok {
		filter6.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp5"] = &filter6
	if v, ok := filtersMap["Temp5"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	paramMap["Filters"] = filtersList

	if err := dataSourceTencentCloudInstancesReadPostFillRequest0(ctx, paramMap); err != nil {
		return err
	}

	var respData []*cvm.Instance
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudInstancesReadPostHandleResponse0(ctx, paramMap, &respData); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudInstancesReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
