package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSesEmailIdentities() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesEmailIdentitiesRead,
		Schema: map[string]*schema.Schema{
			"email_identities": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Sending domain name list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sending domain name.",
						},
						"identity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authentication type, fixed as DOMAIN.",
						},
						"sending_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it verified.",
						},
						"current_reputation_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current credit rating.",
						},
						"daily_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Highest number of letters of the day.",
						},
					},
				},
			},

			"max_reputation_level": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum credit rating.",
			},

			"max_daily_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum daily sending volume for a single domain name.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSesEmailIdentitiesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ses_email_identities.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	var emailIdentities *ses.ListEmailIdentitiesResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesEmailIdentitiesByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		emailIdentities = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(emailIdentities.EmailIdentities))
	tmpList := make([]map[string]interface{}, 0, len(emailIdentities.EmailIdentities))

	if emailIdentities.EmailIdentities != nil {
		for _, emailIdentity := range emailIdentities.EmailIdentities {
			emailIdentityMap := map[string]interface{}{}

			if emailIdentity.IdentityName != nil {
				emailIdentityMap["identity_name"] = emailIdentity.IdentityName
			}

			if emailIdentity.IdentityType != nil {
				emailIdentityMap["identity_type"] = emailIdentity.IdentityType
			}

			if emailIdentity.SendingEnabled != nil {
				emailIdentityMap["sending_enabled"] = emailIdentity.SendingEnabled
			}

			if emailIdentity.CurrentReputationLevel != nil {
				emailIdentityMap["current_reputation_level"] = emailIdentity.CurrentReputationLevel
			}

			if emailIdentity.DailyQuota != nil {
				emailIdentityMap["daily_quota"] = emailIdentity.DailyQuota
			}

			ids = append(ids, *emailIdentity.IdentityName)
			tmpList = append(tmpList, emailIdentityMap)
		}

		_ = d.Set("email_identities", tmpList)
	}

	if emailIdentities.MaxReputationLevel != nil {
		_ = d.Set("max_reputation_level", emailIdentities.MaxReputationLevel)
	}

	if emailIdentities.MaxDailyQuota != nil {
		_ = d.Set("max_daily_quota", emailIdentities.MaxDailyQuota)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
