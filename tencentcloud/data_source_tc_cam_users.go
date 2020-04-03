/*
Use this data source to query detailed information of CAM users

Example Usage

```hcl
# query by name
data "tencentcloud_cam_users" "foo" {
  name      = "cam-user-test"
}

# query by email
data "tencentcloud_cam_users" "bar" {
  email     = "hello@test.com"
}

# query by phone
data "tencentcloud_cam_users" "far" {
  phone_num = "12345678910"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamUsersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of CAM user to be queried.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark of the CAM user to be queried.",
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone num of the CAM user to be queried.",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country code of the CAM user to be queried.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email of the CAM user to be queried.",
			},
			"uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Uin of the CAM user to be queried.",
			},
			"uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Uid of the CAM user to be queried.",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate whether the user can login in.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"user_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM users. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of CAM user. Its value equals to `name` argument.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM user.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark of the CAM user.",
						},
						"phone_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone num of the CAM user.",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country code of the CAM user.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email of the CAM user.",
						},
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of the CAM user.",
						},
						"uid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uid of the CAM user.",
						},
						"console_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate whether the user can login in.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamUsersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_users.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("uin"); ok {
		params["uin"] = v.(int)
	}
	if v, ok := d.GetOk("remark"); ok {
		params["remark"] = v.(string)
	}
	if v, ok := d.GetOk("uid"); ok {
		params["uid"] = v.(int)
	}
	if v, ok := d.GetOk("phone_num"); ok {
		params["phone_num"] = v.(string)
	}
	if v, ok := d.GetOk("country_code"); ok {
		params["country_code"] = v.(string)
	}
	if v, ok := d.GetOk("email"); ok {
		params["email"] = v.(string)
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		consoleLogin := v.(bool)
		if consoleLogin {
			params["console_login"] = 1
		} else {
			params["console_login"] = 0
		}
	}

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var users []*cam.SubAccountInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeUsersByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		users = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM users failed, reason:%s\n", logId, err.Error())
		return err
	}
	userList := make([]map[string]interface{}, 0, len(users))
	ids := make([]string, 0, len(users))
	for _, user := range users {
		mapping := map[string]interface{}{
			"uin":          int(*user.Uin),
			"uid":          int(*user.Uid),
			"name":         *user.Name,
			"remark":       *user.Remark,
			"phone_num":    *user.PhoneNum,
			"country_code": *user.CountryCode,
			"email":        *user.Email,
			"user_id":      *user.Name,
		}
		if int(*user.ConsoleLogin) == 1 {
			mapping["console_login"] = true
		} else {
			mapping["console_login"] = false
		}
		userList = append(userList, mapping)
		ids = append(ids, *user.Name)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("user_list", userList); e != nil {
		log.Printf("[CRITAL]%s provider set CAM user list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), userList); e != nil {
			return e
		}
	}

	return nil
}
