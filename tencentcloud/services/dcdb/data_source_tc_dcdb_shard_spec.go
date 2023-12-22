package dcdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbShardSpec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbShardSpecRead,
		Schema: map[string]*schema.Schema{
			"spec_config": {
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
						"spec_config_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of machine specifications.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "node count.",
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
										Description: "maximum storage size, inGB.",
									},
									"suit_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "recommended usage scenarios.",
									},
									"pid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "product price id.",
									},
									"qps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "maximum QPS.",
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

func dataSourceTencentCloudDcdbShardSpecRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_shard_spec.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var specConfig []*dcdb.SpecConfig

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbShardSpecByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		specConfig = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(specConfig))
	tmpList := make([]map[string]interface{}, 0, len(specConfig))

	if specConfig != nil {
		for _, specConfig := range specConfig {
			specConfigMap := map[string]interface{}{}

			if specConfig.Machine != nil {
				specConfigMap["machine"] = specConfig.Machine
			}

			if specConfig.SpecConfigInfos != nil {
				specConfigInfosList := []interface{}{}
				for _, specConfigInfos := range specConfig.SpecConfigInfos {
					specConfigInfosMap := map[string]interface{}{}

					if specConfigInfos.NodeCount != nil {
						specConfigInfosMap["node_count"] = specConfigInfos.NodeCount
					}

					if specConfigInfos.Memory != nil {
						specConfigInfosMap["memory"] = specConfigInfos.Memory
					}

					if specConfigInfos.MinStorage != nil {
						specConfigInfosMap["min_storage"] = specConfigInfos.MinStorage
					}

					if specConfigInfos.MaxStorage != nil {
						specConfigInfosMap["max_storage"] = specConfigInfos.MaxStorage
					}

					if specConfigInfos.SuitInfo != nil {
						specConfigInfosMap["suit_info"] = specConfigInfos.SuitInfo
					}

					if specConfigInfos.Pid != nil {
						specConfigInfosMap["pid"] = specConfigInfos.Pid
					}

					if specConfigInfos.Qps != nil {
						specConfigInfosMap["qps"] = specConfigInfos.Qps
					}

					if specConfigInfos.Cpu != nil {
						specConfigInfosMap["cpu"] = specConfigInfos.Cpu
					}

					specConfigInfosList = append(specConfigInfosList, specConfigInfosMap)
				}

				specConfigMap["spec_config_infos"] = specConfigInfosList
			}

			ids = append(ids, *specConfig.Machine)
			tmpList = append(tmpList, specConfigMap)
		}

		_ = d.Set("spec_config", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
