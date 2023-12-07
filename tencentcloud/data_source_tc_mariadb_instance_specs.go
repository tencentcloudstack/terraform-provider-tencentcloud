package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbInstanceSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbInstanceSpecsRead,
		Schema: map[string]*schema.Schema{
			"specs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of instance specifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "machine type.",
						},
						"spec_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of machine specifications.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"machine": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "machine type.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "memory, in GB.",
									},
									"min_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "minimum storage size, in GB.",
									},
									"max_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "maximum storage size, in GB.",
									},
									"suit_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "recommended usage scenarios.",
									},
									"qps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "maximum QPS.",
									},
									"pid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "product price id.",
									},
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "node count.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CPU cores.",
									},
								},
							},
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

func dataSourceTencentCloudMariadbInstanceSpecsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_instance_specs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		specs   []*mariadb.InstanceSpec
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbInstanceSpecsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}

		specs = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(specs))
	tmpList := make([]map[string]interface{}, 0, len(specs))
	if specs != nil {
		for _, instanceSpec := range specs {
			instanceSpecMap := map[string]interface{}{}

			if instanceSpec.Machine != nil {
				instanceSpecMap["machine"] = instanceSpec.Machine
			}

			if instanceSpec.SpecInfos != nil {
				specInfosList := []interface{}{}
				for _, specInfos := range instanceSpec.SpecInfos {
					specInfosMap := map[string]interface{}{}

					if specInfos.Machine != nil {
						specInfosMap["machine"] = specInfos.Machine
					}

					if specInfos.Memory != nil {
						specInfosMap["memory"] = specInfos.Memory
					}

					if specInfos.MinStorage != nil {
						specInfosMap["min_storage"] = specInfos.MinStorage
					}

					if specInfos.MaxStorage != nil {
						specInfosMap["max_storage"] = specInfos.MaxStorage
					}

					if specInfos.SuitInfo != nil {
						specInfosMap["suit_info"] = specInfos.SuitInfo
					}

					if specInfos.Qps != nil {
						specInfosMap["qps"] = specInfos.Qps
					}

					if specInfos.Pid != nil {
						specInfosMap["pid"] = specInfos.Pid
					}

					if specInfos.NodeCount != nil {
						specInfosMap["node_count"] = specInfos.NodeCount
					}

					if specInfos.Cpu != nil {
						specInfosMap["cpu"] = specInfos.Cpu
					}

					specInfosList = append(specInfosList, specInfosMap)
				}

				instanceSpecMap["spec_infos"] = specInfosList
			}

			ids = append(ids, *instanceSpec.Machine)
			tmpList = append(tmpList, instanceSpecMap)
		}
		_ = d.Set("specs", tmpList)
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
