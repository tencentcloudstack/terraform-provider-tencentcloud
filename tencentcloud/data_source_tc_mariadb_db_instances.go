/*
Use this data source to query detailed information of mariadb dbInstances

Example Usage

```hcl
data "tencentcloud_mariadb_db_instances" "dbInstances" {
  instance_ids = ""
  project_ids = ""
  search_name = ""
  vpc_id = ""
  subnet_id = ""
  }
```
*/

package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "instance ids.",
			},

			"project_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "project ids.",
			},

			"search_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance name or vip.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id.",
			},

			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "instances info.",
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
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "avaliable zone.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "meory of instance.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "storage of instance.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet id.",
						},
						"db_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db version id.",
						},
						"resource_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "resource tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag value.",
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

func dataSourceTencentCloudMariadbDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_db_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	param := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		param["instance_ids"] = make([]*string, 0)
		instance_idsSet := v.(*schema.Set).List()
		for i := range instance_idsSet {
			instance_ids := instance_idsSet[i].(string)
			param["instance_ids"] = append(param["instance_ids"].([]*string), &instance_ids)
		}
	}

	if v, ok := d.GetOk("project_ids"); ok {
		param["project_ids"] = make([]*int64, 0)
		project_idsSet := v.(*schema.Set).List()
		for i := range project_idsSet {
			project_ids := project_idsSet[i].(int)
			param["project_ids"] = append(param["project_ids"].([]*int64), helper.IntInt64(project_ids))
		}
	}

	if v, ok := d.GetOk("search_name"); ok {
		param["search_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		param["vpc_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		param["subnet_id"] = helper.String(v.(string))
	}

	mariadbService := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instances []*mariadb.DBInstance
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := mariadbService.DescribeMariadbDbInstancesByFilter(ctx, param)
		if e != nil {
			return retryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Mariadb instances failed, reason:%+v", logId, err)
		return err
	}

	if instances != nil {
		instanceList := []interface{}{}
		for _, instance := range instances {
			instanceMap := map[string]interface{}{}
			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}
			if instance.InstanceName != nil {
				instanceMap["instance_name"] = instance.InstanceName
			}
			if instance.ProjectId != nil {
				instanceMap["project_id"] = instance.ProjectId
			}
			if instance.Region != nil {
				instanceMap["region"] = instance.Region
			}
			if instance.Zone != nil {
				instanceMap["zone"] = instance.Zone
			}
			if instance.Memory != nil {
				instanceMap["memory"] = instance.Memory
			}
			if instance.Storage != nil {
				instanceMap["storage"] = instance.Storage
			}
			if instance.VpcId != nil {
				instanceMap["vpc_id"] = instance.VpcId
			}
			if instance.SubnetId != nil {
				instanceMap["subnet_id"] = instance.SubnetId
			}
			if instance.DbVersionId != nil {
				instanceMap["db_version_id"] = instance.DbVersionId
			}
			if instance.ResourceTags != nil {
				resouceTagsList := []interface{}{}
				for _, resourceTags := range instance.ResourceTags {
					resouceTagsMap := map[string]interface{}{}
					if resourceTags.TagKey != nil {
						resouceTagsMap["tag_key"] = resourceTags.TagKey
					}
					if resourceTags.TagValue != nil {
						resouceTagsMap["tag_value"] = resourceTags.TagValue
					}

					resouceTagsList = append(resouceTagsList, resouceTagsMap)
				}
				instanceMap["resource_tags"] = resouceTagsList
			}

			instanceList = append(instanceList, instanceMap)
		}
		_ = d.Set("instances", instanceList)
	}

	return nil
}
