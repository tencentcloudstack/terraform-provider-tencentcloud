package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

func dataSourceTencentCloudStsCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudStsCallerIdentityRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current caller ARN.",
			},

			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The primary account Uin to which the current caller belongs.",
			},

			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity:- When the caller is a cloud account, the current account `Uin` is returned.- When the caller is a role, it returns `roleId:roleSessionName`- When the caller is a federated identity, it returns `uin:federatedUserName`.",
			},

			"principal_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account Uin to which the key belongs:- The caller is a cloud account, and the returned current account Uin- The caller is a role, and the returned account Uin that applies for the role key.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity type.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudStsCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sts_caller_identity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	stsService := StsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var callerIdentity *sts.GetCallerIdentityResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := stsService.DescribeStsCallerIdentityByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		callerIdentity = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Sts instances failed, reason:%+v", logId, err)
		return err
	}

	if callerIdentity.Arn != nil {
		_ = d.Set("arn", callerIdentity.Arn)
	}

	if callerIdentity.AccountId != nil {
		_ = d.Set("account_id", callerIdentity.AccountId)
	}

	if callerIdentity.UserId != nil {
		_ = d.Set("user_id", callerIdentity.UserId)
	}

	if callerIdentity.PrincipalId != nil {
		_ = d.Set("principal_id", callerIdentity.PrincipalId)
	}

	if callerIdentity.Type != nil {
		_ = d.Set("type", callerIdentity.Type)
	}

	d.SetId(*callerIdentity.Arn)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"arn":          callerIdentity.Arn,
			"account_id":   callerIdentity.AccountId,
			"user_id":      callerIdentity.UserId,
			"principal_id": callerIdentity.PrincipalId,
			"type":         callerIdentity.Type,
		}); e != nil {
			return e
		}
	}

	return nil
}
