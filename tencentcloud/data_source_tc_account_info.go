/*
Use this data source to query account info from tencentcloud
Example Usage

```hcl
data "tencentcloud_account_info" "info" {
}
```
*/
package tencentcloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

func dataSourceTencentCloudAccountInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAccountInfoRead,

		Schema: map[string]*schema.Schema{
			"appid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Appid in tencentcloud.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAccountInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_account_info.read")()

	var (
		client = meta.(*TencentCloudClient).apiV3Conn.UseBaseClient()
	)

	type Response struct {
		*tchttp.BaseResponse
		Response *struct {
			AppId *int64
		}
	}

	var res = &Response{
		BaseResponse: &tchttp.BaseResponse{},
	}
	var req = &tchttp.BaseRequest{}

	req = req.Init().WithApiInfo("tic", "2019-12-13", "DescribeAccount")

	err := client.Send(req, res)
	if err != nil {
		err = client.Send(req, res)
	}
	if err != nil {
		return err
	}

	if res.Response == nil || res.Response.AppId == nil || *res.Response.AppId == 0 {
		return errors.New("api tic.DescribeAccount return empty appid")
	}
	_ = d.Set("appid", res.Response.AppId)
	d.SetId(fmt.Sprintf("%d", *res.Response.AppId))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"appid":res.Response.AppId,
		}); e != nil {
			return e
		}
	}

	return nil
}
