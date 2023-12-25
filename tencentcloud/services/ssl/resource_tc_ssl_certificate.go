package ssl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCertificateCreate,
		Read:   resourceTencentCloudSslCertificateRead,
		Update: resourceTencentCloudSslCertificateUpdate,
		Delete: resourceTencentCloudSslCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the SSL certificate.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SSL_CERT_TYPE),
				ForceNew:     true,
				Description:  "Type of the SSL certificate. Valid values: `CA` and `SVR`.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID of the SSL certificate. Default is `0`.",
			},
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of the SSL certificate. Not allowed newline at the start and end.",
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n suffix"))
					}
					return
				},
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Key of the SSL certificate and required when certificate type is `SVR`. Not allowed newline at the start and end.",
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n suffix"))
					}
					return
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tags of the SSL certificate.",
			},
			// computed
			"product_zh_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate authority.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the SSL certificate.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the SSL certificate.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Beginning time of the SSL certificate.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending time of the SSL certificate.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the SSL certificate.",
			},
			"subject_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "ALL domains included in the SSL certificate. Including the primary domain name.",
			},
		},
	}
}

func resourceTencentCloudSslCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_certificate.create")()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sslService       = SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr    error
		id               string
		describeResponse *ssl.DescribeCertificateDetailResponse
	)

	request := ssl.NewUploadCertificateRequest()
	request.CertificatePublicKey = helper.String(d.Get("cert").(string))
	request.CertificateType = helper.String(d.Get("type").(string))
	request.ProjectId = helper.Uint64(uint64(d.Get("project_id").(int)))
	request.Alias = helper.String(d.Get("name").(string))
	if raw, ok := d.GetOk("key"); ok {
		request.CertificatePrivateKey = helper.String(raw.(string))
	}
	if *request.CertificateType == "SVR" && (request.CertificatePrivateKey == nil || *request.CertificatePrivateKey == "") {
		return errors.New("when type is SVR, key can't be empty")
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		id, inErr = sslService.UploadCertificate(ctx, request)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, outErr)
		return outErr
	}

	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &id
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if describeResponse == nil || describeResponse.Response == nil {
			err := fmt.Errorf("TencentCloud SDK %s return empty response", describeRequest.GetAction())
			return tccommon.RetryError(err)
		}
		if describeResponse.Response.Status == nil {
			err := fmt.Errorf("api[%s] certificate status is nil", describeRequest.GetAction())
			return resource.NonRetryableError(err)
		}

		if *describeResponse.Response.Status != SSL_STATUS_AVAILABLE {
			err := fmt.Errorf("certificate is not available, status is %d", *describeResponse.Response.Status)
			return resource.RetryableError(err)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s create certificate failed, reason: %v", logId, outErr)
		return outErr
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagClient := m.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tagClient)
		resourceName := tccommon.BuildTagResourceName("ssl", "certificate", tagClient.Region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(id)

	return resourceTencentCloudSslCertificateRead(d, m)
}

func resourceTencentCloudSslCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_certificate.read")()
	defer tccommon.InconsistentCheck(d, m)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sslService       = SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr    error
		id               = d.Id()
		describeResponse *ssl.DescribeCertificateDetailResponse
	)

	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &id
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
		if inErr != nil {
			if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == CertificateNotFound {
					return nil
				}
			}
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s read certificate failed, reason: %v", logId, outErr)
		return outErr
	}

	if describeResponse == nil || describeResponse.Response == nil || describeResponse.Response.CertificateId == nil {
		d.SetId("")
		return nil
	}

	certificate := describeResponse.Response
	if nilNames := tccommon.CheckNil(certificate, map[string]string{
		"Alias":                "name",
		"CertificateType":      "type",
		"ProjectId":            "project id",
		"CertificatePublicKey": "cert",
		"ProductZhName":        "product zh name",
		"Domain":               "domain",
		"Status":               "status",
		"CertBeginTime":        "begin time",
		"CertEndTime":          "end time",
		"InsertTime":           "create time",
	}); len(nilNames) > 0 {
		return fmt.Errorf("certificate %v are nil", nilNames)
	}

	_ = d.Set("name", certificate.Alias)
	_ = d.Set("type", certificate.CertificateType)
	projectId, err := strconv.Atoi(*certificate.ProjectId)
	if err != nil {
		return err
	}
	_ = d.Set("project_id", projectId)
	_ = d.Set("cert", strings.TrimRight(*certificate.CertificatePublicKey, "\n"))
	_ = d.Set("product_zh_name", certificate.ProductZhName)
	_ = d.Set("domain", certificate.Domain)
	_ = d.Set("status", certificate.Status)
	_ = d.Set("begin_time", certificate.CertBeginTime)
	_ = d.Set("end_time", certificate.CertEndTime)
	_ = d.Set("create_time", certificate.InsertTime)

	subjectAltNames := make([]string, 0, len(certificate.SubjectAltName))
	for _, subjectAltName := range certificate.SubjectAltName {
		subjectAltNames = append(subjectAltNames, *subjectAltName)
	}
	_ = d.Set("subject_names", subjectAltNames)

	tagClient := m.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tagClient)

	tags, err := tagService.DescribeResourceTags(ctx, "ssl", "certificate", tagClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSslCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_certificate.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id         = d.Id()
		sslService = SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	d.Partial(true)
	if d.HasChange("name") {
		aliasRequest := ssl.NewModifyCertificateAliasRequest()
		aliasRequest.CertificateId = helper.String(id)
		_, alias := d.GetChange("name")
		aliasRequest.Alias = helper.String(alias.(string))

		if outErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			if inErr := sslService.ModifyCertificateAlias(ctx, aliasRequest); inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					code := sdkErr.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkErr)
					}
				}
				return tccommon.RetryError(inErr)
			}
			return nil
		}); outErr != nil {
			return outErr
		}

	}
	if d.HasChange("project_id") {
		projectRequest := ssl.NewModifyCertificateProjectRequest()
		projectRequest.CertificateIdList = []*string{
			helper.String(id),
		}
		_, projectId := d.GetChange("project_id")
		projectRequest.ProjectId = helper.Uint64(uint64(projectId.(int)))

		if outErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			if inErr := sslService.ModifyCertificateProject(ctx, projectRequest); inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					code := sdkErr.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkErr)
					}
				}
				return tccommon.RetryError(inErr)
			}
			return nil
		}); outErr != nil {
			return outErr
		}

	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagClient := m.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tagClient)
		resourceName := tccommon.BuildTagResourceName("ssl", "certificate", tagClient.Region, id)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}
	d.Partial(false)
	return resourceTencentCloudSslCertificateRead(d, m)
}

func resourceTencentCloudSslCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_certificate.delete")()
	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sslService    = SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr error
		id            = d.Id()
		deleteResult  bool
	)
	request := ssl.NewDeleteCertificateRequest()
	request.CertificateId = helper.String(id)

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		deleteResult, inErr = sslService.DeleteCertificate(ctx, request)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if !deleteResult {
			return resource.NonRetryableError(errors.New("failed to delete certificate"))
		}
		return nil
	})

	if outErr != nil {
		log.Printf("[CRITAL]%s delete SSL certificate failed, reason:%+v", logId, outErr)
		return outErr
	}
	return nil
}
