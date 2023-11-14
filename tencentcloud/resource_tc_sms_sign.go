/*
Provides a resource to create a sms sign

Example Usage

```hcl
resource "tencentcloud_sms_sign" "sign" {
  sign_name = "SignName"
  sign_type = 0
  document_type = 0
  international = 0
  sign_purpose = 0
  proof_image = &lt;nil&gt;
  commission_image = &lt;nil&gt;
  remark = &lt;nil&gt;
      }
```

Import

sms sign can be imported using the id, e.g.

```
terraform import tencentcloud_sms_sign.sign sign_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSmsSign() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSmsSignCreate,
		Read:   resourceTencentCloudSmsSignRead,
		Update: resourceTencentCloudSmsSignUpdate,
		Delete: resourceTencentCloudSmsSignDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sign_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sms sign name, unique.",
			},

			"sign_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Signature type. Each of these types is followed by their DocumentType (identity certificate type) option: 0: company. Valid values of DocumentType include 0 and 1. 1: app. Valid values of DocumentType include 0, 1, 2, 3, and 4. 2: website. Valid values of DocumentType include 0, 1, 2, 3, and 5. 3: WeChat Official Account. Valid values of DocumentType include 0, 1, 2, 3, and 8. 4: trademark. Valid values of DocumentType include 7. 5: government/public institution/other. Valid values of DocumentType include 2 and 3. 6: WeChat Mini Program. Valid values of DocumentType include 0, 1, 2, 3, and 6. Note: the identity certificate type must be selected according to the correspondence; otherwise, the review will fail.",
			},

			"document_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Identity certificate type: 0: three-in-one licence. 1: business license. 2: organization code certificate. 3: social credit code certificate. 4: screenshot of application backend management (for personal app). 5: screenshot of website ICP filing backend (for personal website). 6: screenshot of WeChat Mini Program settings page (for personal WeChat Mini Program). 7: trademark registration certificate. 8: screenshot of WeChat Official Account settings page (for personal WeChat Official Account).",
			},

			"international": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.",
			},

			"sign_purpose": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Signature purpose: 0: for personal use.   1: for others.",
			},

			"proof_image": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "You should Base64-encode the image of the identity certificate corresponding to the signature first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter.",
			},

			"commission_image": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Power of attorney, which should be submitted if SignPurpose is for use by others. You should Base64-encode the image first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter. Note: this field will take effect only when SignPurpose is 1 (for user by others).",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Signature application remarks.",
			},

			"review_reply": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Review reply, i.e., response given by the reviewer, which is usually the reason for rejection.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Application submission time in the format of UNIX timestamp in seconds.",
			},

			"status_code": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Signature application status. Valid values: 0: approved; 1: under review; -1: application rejected or failed.",
			},
		},
	}
}

func resourceTencentCloudSmsSignCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_sign.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = sms.NewAddSmsSignRequest()
		response = sms.NewAddSmsSignResponse()
		signId   int
	)
	if v, ok := d.GetOk("sign_name"); ok {
		request.SignName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sign_type"); ok {
		request.SignType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("document_type"); ok {
		request.DocumentType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("international"); ok {
		request.International = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("sign_purpose"); ok {
		request.SignPurpose = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("proof_image"); ok {
		request.ProofImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("commission_image"); ok {
		request.CommissionImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSmsClient().AddSmsSign(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sms sign failed, reason:%+v", logId, err)
		return err
	}

	signId = *response.Response.SignId
	d.SetId(helper.Int64ToStr(int64(signId)))

	return resourceTencentCloudSmsSignRead(d, meta)
}

func resourceTencentCloudSmsSignRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_sign.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	signId := d.Id()

	sign, err := service.DescribeSmsSignById(ctx, signId)
	if err != nil {
		return err
	}

	if sign == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SmsSign` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if sign.SignName != nil {
		_ = d.Set("sign_name", sign.SignName)
	}

	if sign.SignType != nil {
		_ = d.Set("sign_type", sign.SignType)
	}

	if sign.DocumentType != nil {
		_ = d.Set("document_type", sign.DocumentType)
	}

	if sign.International != nil {
		_ = d.Set("international", sign.International)
	}

	if sign.SignPurpose != nil {
		_ = d.Set("sign_purpose", sign.SignPurpose)
	}

	if sign.ProofImage != nil {
		_ = d.Set("proof_image", sign.ProofImage)
	}

	if sign.CommissionImage != nil {
		_ = d.Set("commission_image", sign.CommissionImage)
	}

	if sign.Remark != nil {
		_ = d.Set("remark", sign.Remark)
	}

	if sign.ReviewReply != nil {
		_ = d.Set("review_reply", sign.ReviewReply)
	}

	if sign.CreateTime != nil {
		_ = d.Set("create_time", sign.CreateTime)
	}

	if sign.StatusCode != nil {
		_ = d.Set("status_code", sign.StatusCode)
	}

	return nil
}

func resourceTencentCloudSmsSignUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_sign.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sms.NewModifySmsSignRequest()

	signId := d.Id()

	request.SignId = &signId

	immutableArgs := []string{"sign_name", "sign_type", "document_type", "international", "sign_purpose", "proof_image", "commission_image", "remark", "review_reply", "create_time", "status_code"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("sign_name") {
		if v, ok := d.GetOk("sign_name"); ok {
			request.SignName = helper.String(v.(string))
		}
	}

	if d.HasChange("sign_type") {
		if v, ok := d.GetOkExists("sign_type"); ok {
			request.SignType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("document_type") {
		if v, ok := d.GetOkExists("document_type"); ok {
			request.DocumentType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("international") {
		if v, ok := d.GetOkExists("international"); ok {
			request.International = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("sign_purpose") {
		if v, ok := d.GetOkExists("sign_purpose"); ok {
			request.SignPurpose = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("proof_image") {
		if v, ok := d.GetOk("proof_image"); ok {
			request.ProofImage = helper.String(v.(string))
		}
	}

	if d.HasChange("commission_image") {
		if v, ok := d.GetOk("commission_image"); ok {
			request.CommissionImage = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSmsClient().ModifySmsSign(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sms sign failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSmsSignRead(d, meta)
}

func resourceTencentCloudSmsSignDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_sign.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}
	signId := d.Id()

	if err := service.DeleteSmsSignById(ctx, signId); err != nil {
		return err
	}

	return nil
}
