package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSesBlackEmailAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesBlackEmailAddressRead,
		Schema: map[string]*schema.Schema{
			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start date in the format of `YYYY-MM-DD`.",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End date in the format of `YYYY-MM-DD`.",
			},

			"email_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "You can specify an email address to query.",
			},

			"task_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "You can specify a task ID to query.",
			},

			"black_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of blocklisted addresses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bounce_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the email address is blocklisted.",
						},
						"email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blocklisted email address.",
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

func dataSourceTencentCloudSesBlackEmailAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ses_black_email_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email_address"); ok {
		paramMap["EmailAddress"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskID"] = helper.String(v.(string))
	}

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	var blackList []*ses.BlackEmailAddress

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesBlackEmailAddressByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		blackList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(blackList))
	tmpList := make([]map[string]interface{}, 0, len(blackList))

	if blackList != nil {
		for _, blackEmailAddress := range blackList {
			blackEmailAddressMap := map[string]interface{}{}

			if blackEmailAddress.BounceTime != nil {
				blackEmailAddressMap["bounce_time"] = blackEmailAddress.BounceTime
			}

			if blackEmailAddress.EmailAddress != nil {
				blackEmailAddressMap["email_address"] = blackEmailAddress.EmailAddress
			}

			ids = append(ids, *blackEmailAddress.EmailAddress)
			tmpList = append(tmpList, blackEmailAddressMap)
		}

		_ = d.Set("black_list", tmpList)
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
