/*

Use this data source to query tcaplus applications

Example Usage

```hcl
data "tencentcloud_tcaplus_applications" "name" {
  app_name = "app"
}
data "tencentcloud_tcaplus_applications" "id" {
  app_id = tencentcloud_tcaplus_application.test.id
}
data "tencentcloud_tcaplus_applications" "idname" {
  app_id   = tencentcloud_tcaplus_application.test.id
  app_name = "app"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudTcaplusApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusApplicationsRead,
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the tcapplus application to be query.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the tcapplus application to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus application. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the tcapplus application.",
						},
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Id of the tcapplus application.",
						},
						"idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the tcapplus application.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id of the tcapplus application.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id of the tcapplus application.",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Password of the tcapplus application.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "network type of the tcapplus application.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time of the tcapplus application.",
						},
						"password_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "password status of the tcapplus application.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.",
						},
						"api_access_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access id of the tcapplus application.For TcaplusDB SDK connect.",
						},
						"api_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access ip of the tcapplus application.For TcaplusDB SDK connect.",
						},
						"api_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "access port of the tcapplus application.For TcaplusDB SDK connect.",
						},
						"old_password_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_applications.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	applicationId := d.Get("app_id").(string)
	applicationName := d.Get("app_name").(string)

	apps, err := service.DescribeApps(ctx, applicationId, applicationName)
	if err != nil {
		apps, err = service.DescribeApps(ctx, applicationId, applicationName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(apps))

	for _, app := range apps {
		listItem := make(map[string]interface{})
		listItem["app_name"] = *app.AppName
		listItem["app_id"] = *app.ApplicationId
		listItem["idl_type"] = *app.IdlType
		listItem["vpc_id"] = *app.VpcId
		listItem["subnet_id"] = *app.SubnetId
		listItem["password"] = *app.Password
		listItem["network_type"] = *app.NetworkType
		listItem["create_time"] = *app.CreatedTime
		listItem["password_status"] = *app.PasswordStatus
		listItem["api_access_id"] = *app.ApiAccessId
		listItem["api_access_ip"] = *app.ApiAccessIp
		listItem["api_access_port"] = *app.ApiAccessPort
		listItem["old_password_expire_time"] = *app.OldPasswordExpireTime
		list = append(list, listItem)
	}

	d.SetId("app." + applicationId + "." + applicationName)
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
