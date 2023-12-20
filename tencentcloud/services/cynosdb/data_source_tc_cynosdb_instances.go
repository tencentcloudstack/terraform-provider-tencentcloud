package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbInstancesRead,

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of CynosDB, and available values include `MYSQL`, `POSTGRESQL`.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the cluster.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the Cynosdb instance to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project to be queried.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Cynosdb instance to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of instances. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the cluster.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CynosDB instance.",
						},
						"instance_cpu_core": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of CPU cores of the Cynosdb instance.",
						},
						"instance_memory_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Memory capacity of the Cynosdb instance, unit in GB.",
						},
						"instance_storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage size of the Cynosdb instance, unit in GB.",
						},
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the Cynosdb instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the CynosDB instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. `ro` for readonly instance, `rw` for read and write instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCynosdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]string)
	if v, ok := d.GetOk("instance_id"); ok {
		params["InstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		params["InstanceName"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["ProjectId"] = fmt.Sprintf("%d", v.(int))
	}
	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	cynosdbService := CynosdbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, e := cynosdbService.DescribeInstances(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ids := make([]string, 0, len(instances))
		instanceList := make([]map[string]interface{}, 0, len(instances))
		for _, instance := range instances {
			if clusterId != "" && *instance.ClusterId != clusterId {
				continue
			}

			mapping := map[string]interface{}{
				"cluster_id":            instance.ClusterId,
				"instance_id":           instance.InstanceId,
				"instance_name":         instance.InstanceName,
				"instance_memory_size":  instance.Memory,
				"instance_storage_size": instance.Storage,
				"instance_cpu_core":     instance.Cpu,
				"create_time":           instance.CreateTime,
				"instance_status":       instance.Status,
				"instance_type":         instance.InstanceType,
			}

			instanceList = append(instanceList, mapping)
			ids = append(ids, *instance.ClusterId)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		if e = d.Set("instance_list", instanceList); e != nil {
			log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, e.Error())
			return resource.NonRetryableError(e)
		}

		output, ok := d.GetOk("result_output_file")
		if ok && output.(string) != "" {
			if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
				return resource.NonRetryableError(e)
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cynosdb instances failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
