/*
Use this data source to query user appid, uin and ownerUin.

Example Usage

```hcl
data "tencentcloud_user_info" "foo" {}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceTencentCloudUserInfo() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudUserInfoRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account App ID.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account UIN.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account OwnerUIN.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
		},
	}
}

func datasourceTencentCloudUserInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("datasource.tencentcloud_user_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn

	service := CamService{client: client}

	response, err := service.GetUserAppId(ctx)
	if err != nil {
		return err
	}

	result := response.Response

	if result == nil {
		return fmt.Errorf("get user appid error: empty response")
	}

	appId := strconv.FormatUint(*result.AppId, 10)
	uin := *result.Uin
	ownerUin := *result.OwnerUin

	d.SetId(fmt.Sprintf("user-%s-%s-%d", uin, appId, rand.Intn(10000)))

	_ = d.Set("app_id", appId)
	_ = d.Set("uin", uin)
	_ = d.Set("owner_uin", ownerUin)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"app_id":   appId,
			"uin":      uin,
			"ownerUin": ownerUin,
		}); e != nil {
			return e
		}
	}

	return nil
}
