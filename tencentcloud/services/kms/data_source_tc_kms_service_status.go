package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsServiceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsServiceStatusRead,
		Schema: map[string]*schema.Schema{
			"service_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the KMS service has been activated. true: activated.",
			},

			"invalid_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Service unavailability type. 0: not purchased; 1: normal; 2: suspended due to arrears; 3: resource released.",
			},

			"user_level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "0: Basic Edition, 1: Ultimate Edition.",
			},

			"pro_expire_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Expiration time of the KMS Ultimate edition. It's represented in a Unix Epoch timestamp.\nNote: This field may return null, indicating that no valid values can be obtained.",
			},

			"pro_renew_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether to automatically renew Ultimate Edition. 0: no, 1: yes\nNote: this field may return null, indicating that no valid values can be obtained.",
			},

			"pro_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Ultimate Edition purchase record. If the Ultimate Edition is not activated, the returned value will be null.\nNote: this field may return null, indicating that no valid values can be obtained.",
			},

			"exclusive_vsm_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to activate Managed KMS\nNote: This field may return `null`, indicating that no valid value can be obtained.",
			},

			"exclusive_hsm_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to activate Exclusive KMS\nNote: This field may return `null`, indicating that no valid value can be obtained.",
			},

			"subscription_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "KMS subscription information.\nNote: This field may return null, indicating that no valid values can be obtained.",
			},

			"cmk_user_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Return the number of KMS user key usage.",
			},

			"cmk_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Return KMS user key specification quantity.",
			},

			"exclusive_hsm_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Return to Exclusive Cluster Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hsm_cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Exclusive cluster ID.",
						},
						"hsm_cluster_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Exclusive cluster name.",
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

func dataSourceTencentCloudKmsServiceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_service_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	var respData *kms.GetServiceStatusResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsServiceStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData.ServiceEnabled != nil {
		_ = d.Set("service_enabled", respData.ServiceEnabled)
	}

	if respData.InvalidType != nil {
		_ = d.Set("invalid_type", respData.InvalidType)
	}

	if respData.UserLevel != nil {
		_ = d.Set("user_level", respData.UserLevel)
	}

	if respData.ProExpireTime != nil {
		_ = d.Set("pro_expire_time", respData.ProExpireTime)
	}

	if respData.ProRenewFlag != nil {
		_ = d.Set("pro_renew_flag", respData.ProRenewFlag)
	}

	if respData.ProResourceId != nil {
		_ = d.Set("pro_resource_id", respData.ProResourceId)
	}

	if respData.ExclusiveVSMEnabled != nil {
		_ = d.Set("exclusive_vsm_enabled", respData.ExclusiveVSMEnabled)
	}

	if respData.ExclusiveHSMEnabled != nil {
		_ = d.Set("exclusive_hsm_enabled", respData.ExclusiveHSMEnabled)
	}

	if respData.SubscriptionInfo != nil {
		_ = d.Set("subscription_info", respData.SubscriptionInfo)
	}

	if respData.CmkUserCount != nil {
		_ = d.Set("cmk_user_count", respData.CmkUserCount)
	}

	if respData.CmkLimit != nil {
		_ = d.Set("cmk_limit", respData.CmkLimit)
	}

	if respData.ExclusiveHSMList != nil {
		tmpList := make([]map[string]interface{}, 0, len(respData.ExclusiveHSMList))
		for _, item := range respData.ExclusiveHSMList {
			dMap := make(map[string]interface{})
			if item.HsmClusterId != nil {
				dMap["hsm_cluster_id"] = item.HsmClusterId
			}

			if item.HsmClusterName != nil {
				dMap["hsm_cluster_name"] = item.HsmClusterName
			}

			tmpList = append(tmpList, dMap)
		}

		_ = d.Set("exclusive_hsm_list", tmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
