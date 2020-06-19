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
							Description: "ID of the postgresql instance speccode.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size(in GB).",
						},
						"storage_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum volume size(in GB).",
						},
						"storage_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum volume size(in GB).",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The CPU number of the postgresql instance.",
						},
						"qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The QPS of the postgresql instance.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the postgresql database engine.",
						},
						"engine_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version name of the postgresql database engine.",
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
		listItem["memory"] = *v.Memory / 1024
		listItem["storage_min"] = v.MinStorage
		listItem["storage_max"] = v.MaxStorage
		listItem["cpu"] = v.Cpu
		listItem["qps"] = v.Qps
		listItem["engine_version"] = v.Version
		listItem["engine_version_name"] = v.VersionName
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
