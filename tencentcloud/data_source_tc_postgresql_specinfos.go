/*
Use this data source to get the available product configs of the postgresql instance.

Example Usage

```hcl
data "tencentcloud_postgresql_specinfos" "foo"{
	availability_zone = "ap-shanghai-2"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudPostgresqlSpecinfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlSpecinfosRead,
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone of the postgresql instance to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of zones will be exported and its every element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the speccode of the postgresql instance. This parameter is used as `spec_code` for the creation of postgresql instance.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size(in MB).",
						},
						"storage_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum volume size(in GB).",
						},
						"storage_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The max volume size(in GB).",
						},
						"cpu_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The CPU number of the postgresql instance.",
						},
						"qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The QPS of the postgresql instance.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the postgresql instance.",
						},
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version name of the postgresql instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlSpecinfosRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgresql_specinfos.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	zone := d.Get("availability_zone").(string)
	speccodes, err := service.DescribeSpecinfos(ctx, zone)
	if err != nil {
		speccodes, err = service.DescribeSpecinfos(ctx, zone)
	}
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(speccodes))

	for _, v := range speccodes {
		listItem := make(map[string]interface{})
		listItem["id"] = v.SpecCode
		listItem["memory"] = v.Memory
		listItem["storage_min"] = v.MinStorage
		listItem["storage_max"] = v.MaxStorage
		listItem["cpu_number"] = v.Cpu
		listItem["qps"] = v.Qps
		listItem["version"] = v.Version
		listItem["version_name"] = v.VersionName
		list = append(list, listItem)
	}
	d.SetId("speccode." + zone)
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
