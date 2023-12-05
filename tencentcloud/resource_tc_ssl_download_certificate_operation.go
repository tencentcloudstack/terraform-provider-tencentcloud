package tencentcloud

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslDownloadCertificateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslDownloadCertificateCreate,
		Read:   resourceTencentCloudSslDownloadCertificateRead,
		Delete: resourceTencentCloudSslDownloadCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},
			"output_path": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},
		},
	}
}

func resourceTencentCloudSslDownloadCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_download_certificate_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ssl.NewDownloadCertificateRequest()
		response      = ssl.NewDownloadCertificateResponse()
		certificateId uint64
		outputPath    string
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = helper.StrToUInt64(v.(string))
		request.CertificateId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("output_path"); ok {
		outputPath = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().DownloadCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl downloadCertificate failed, reason:%+v", logId, err)
		return err
	}

	if response != nil && response.Response.Content != nil {
		zipBytes, err := base64.StdEncoding.DecodeString(*response.Response.Content)
		if err != nil {
			return err
		}

		// 保存为Zip文件
		err = ioutil.WriteFile("output.zip", zipBytes, 0644)
		if err != nil {
			return err
		}

		log.Printf("Zip file saved successfully. certificateId[%v] path[%s]", certificateId, outputPath)
	}
	d.SetId(helper.UInt64ToStr(certificateId))

	return resourceTencentCloudSslDownloadCertificateRead(d, meta)
}

func resourceTencentCloudSslDownloadCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_download_certificate_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslDownloadCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_download_certificate_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
