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
  proof_image = ""
  commission_image = ""
  remark = ""
}

```
Import

sms sign can be imported using the id, e.g.
```
$ terraform import tencentcloud_sms_sign.sign sign_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSmsSign() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSmsSignRead,
		Create: resourceTencentCloudSmsSignCreate,

		Update: resourceTencentCloudSmsSignUpdate,
		Delete: resourceTencentCloudSmsSignDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sign_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sms sign name, unique.",
			},

			"sign_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Sms sign type: 0、1、2、3、4、5、6.",
			},

			"document_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "DocumentType is used for enterprise authentication, or website, app authentication, etc. DocumentType: 0, 1, 2 3, 4, 5, 6, 7, 8.",
			},

			"international": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.",
			},

			"sign_purpose": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Signature purpose: 0: for personal use.   1: for others.",
			},

			"proof_image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You should Base64-encode the image of the identity certificate corresponding to the signature first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter.",
			},

			"commission_image": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Power of attorney, which should be submitted if SignPurpose is for use by others. You should Base64-encode the image first, remove the prefix data:image/jpeg;base64, from the resulted string, and then use it as the value of this parameter. Note: this field will take effect only when SignPurpose is 1 (for user by others).",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Signature application remarks.",
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
		response *sms.AddSmsSignResponse
		signId   uint64
	)

	if v, ok := d.GetOk("sign_name"); ok {
		request.SignName = helper.String(v.(string))
	}

	// uint64 当值为0时，d.GetOk()  ok = false, 参数无法传入
	// if v, ok := d.GetOk("sign_type"); ok {
	// 	request.SignType = helper.IntUint64(v.(int))
	// }

	// if v, ok := d.GetOk("document_type"); ok {
	// 	request.DocumentType = helper.IntUint64(v.(int))
	// }

	// if v, ok := d.GetOk("international"); ok {
	// 	request.International = helper.IntUint64(v.(int))
	// }

	// if v, ok := d.GetOk("sign_purpose"); ok {
	// 	request.SignPurpose = helper.IntUint64(v.(int))
	// }
	sign_type := d.Get("sign_type")
	request.SignType = helper.IntUint64(sign_type.(int))
	document_type := d.Get("document_type")
	request.DocumentType = helper.IntUint64(document_type.(int))
	international := d.Get("international")
	request.International = helper.IntUint64(international.(int))
	sign_purpose := d.Get("sign_purpose")
	request.SignPurpose = helper.IntUint64(sign_purpose.(int))

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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sms sign failed, reason:%+v", logId, err)
		return err
	}

	signId = *response.Response.AddSignStatus.SignId //数据结构修改
	signId_string := strconv.FormatUint(signId, 10) //数据转换
	d.Set("international", helper.IntUint64(international.(int)))
	d.SetId(signId_string)  //id字符串类型

	return resourceTencentCloudSmsSignRead(d, meta)
}

func resourceTencentCloudSmsSignRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_sign.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	signId := d.Id()
	international := helper.IntUint64(d.Get("international").(int))
	sign, err := service.DescribeSmsSign(ctx, signId, international)

	if err != nil {
		return err
	}

	if sign == nil {
		d.SetId("")
		return fmt.Errorf("resource `sign` %s does not exist", signId)
	}

	if sign.SignName != nil {
		_ = d.Set("sign_name", sign.SignName)
	}

	if sign.SignId != nil {
		_ = d.Set("sign_id", sign.SignId)
	}

	if sign.International != nil {
		_ = d.Set("international", sign.International)
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
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := sms.NewModifySmsSignRequest()

	// 类型转换
	signId_string := d.Id()

	temp, _ := strconv.ParseUint(signId_string, 10, 64)
	request.SignId = &temp

	if d.HasChange("sign_name") {
		if v, ok := d.GetOk("sign_name"); ok {
			request.SignName = helper.String(v.(string))
		}

	}

	if d.HasChange("sign_type") {
		if v, ok := d.GetOk("sign_type"); ok {
			request.SignType = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("document_type") {
		if v, ok := d.GetOk("document_type"); ok {
			request.DocumentType = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("international") {
		if v, ok := d.GetOk("international"); ok {
			request.International = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("sign_purpose") {
		if v, ok := d.GetOk("sign_purpose"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sms sign failed, reason:%+v", logId, err)
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
