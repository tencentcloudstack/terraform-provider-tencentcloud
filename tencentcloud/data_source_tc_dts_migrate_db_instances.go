package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDtsMigrateDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsMigrateDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"database_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database type.",
			},

			"migrate_role": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether the instance is the migration source or destination,src(for source), dst(for destination).",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database instance id.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database instance name.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Limit.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset.",
			},

			"account_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The owning account of the resource is null or self(resources in the self account), other(resources in the other account).",
			},

			"tmp_secret_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary secret id, used across account.",
			},

			"tmp_secret_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary secret key, used across account.",
			},

			"tmp_token": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary token, used across account.",
			},

			"instances": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database instance name.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance vip.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance port.",
						},
						"usable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Can used in migration, 1-yes, 0-no.",
						},
						"hint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason of can&#39;t used in migration.",
						},
					},
				},
			},

			"request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique request id, provide this when encounter a problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDtsMigrateDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dts_migrate_db_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("database_type"); ok {
		paramMap["DatabaseType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("migrate_role"); ok {
		paramMap["MigrateRole"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_mode"); ok {
		paramMap["AccountMode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_secret_id"); ok {
		paramMap["TmpSecretId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_secret_key"); ok {
		paramMap["TmpSecretKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_token"); ok {
		paramMap["TmpToken"] = helper.String(v.(string))
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instances []*dts.MigrateDBItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDtsMigrateDbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instances = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instances))
	tmpList := make([]map[string]interface{}, 0, len(instances))

	if instances != nil {
		for _, migrateDBItem := range instances {
			migrateDBItemMap := map[string]interface{}{}

			if migrateDBItem.InstanceId != nil {
				migrateDBItemMap["instance_id"] = migrateDBItem.InstanceId
			}

			if migrateDBItem.InstanceName != nil {
				migrateDBItemMap["instance_name"] = migrateDBItem.InstanceName
			}

			if migrateDBItem.Vip != nil {
				migrateDBItemMap["vip"] = migrateDBItem.Vip
			}

			if migrateDBItem.Vport != nil {
				migrateDBItemMap["vport"] = migrateDBItem.Vport
			}

			if migrateDBItem.Usable != nil {
				migrateDBItemMap["usable"] = migrateDBItem.Usable
			}

			if migrateDBItem.Hint != nil {
				migrateDBItemMap["hint"] = migrateDBItem.Hint
			}

			ids = append(ids, *migrateDBItem.InstanceId)
			tmpList = append(tmpList, migrateDBItemMap)
		}

		_ = d.Set("instances", tmpList)
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
