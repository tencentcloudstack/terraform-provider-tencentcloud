/*
Provide a resource to create a Free Certificate.

~> **NOTE:** Once certificat created, it cannot be removed within 1 hours.

Example Usage

Currently, `package_type` only support type 2. 2=TrustAsia TLS RSA CA.

```hcl
resource "tencentcloud_ssl_free_certificate" "example" {
  dv_auth_method    = "DNS_AUTO"
  domain            = "example.com"
  package_type      = "2"
  contact_email     = "test@example.com"
  contact_phone     = "18352458901"
  validity_period   = 12
  csr_encrypt_algo  = "RSA"
  csr_key_parameter = "2048"
  csr_key_password  = "csr_pwd"
  alias             = "example_free_cert"
}
```

Import

FreeCertificate instance can be imported, e.g.
```
$ terraform import tencentcloud_ssl_free_certificate.test free_certificate-id
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"

	ssl2 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSSLFreeCertificate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudFreeCertificateRead,
		Create: resourceTencentCloudFreeCertificateCreate,
		Update: resourceTencentCloudFreeCertificateUpdate,
		Delete: resourceTencentCloudFreeCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dv_auth_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specify DV authorize method. Available values: `DNS_AUTO` - automatic DNS auth, `DNS` - manual DNS auth, `FILE` - auth by file.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specify domain name.",
			},
			"package_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of package. Only support `\"2\"` (TrustAsia TLS RSA CA).",
			},
			"contact_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address.",
			},
			"contact_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number.",
			},
			"validity_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify validity period in month, only support `\"12\"` months for now.",
			},
			"csr_encrypt_algo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify CSR encrypt algorithm, only support `RSA` for now.",
			},
			"csr_key_parameter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify CSR key parameter, only support `\"2048\"` for now.",
			},
			"csr_key_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Specify CSR key password.",
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify alias for remark.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of projects which this certification belong to.",
			},
			"old_certificate_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Specify old certificate ID, used for re-apply.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate status. 0 = Approving, 1 = Approved, 2 = Approve failed, 3 = expired, 4 = DNS record added, 5 = OV/EV Certificate and confirm letter needed, 6 = Order canceling, 7 = Order canceled, 8 = Submitted and confirm letter needed, 9 = Revoking, 10 = Revoked, 11 = re-applying, 12 = Revoke and confirm letter needed, 13 = Free SSL and confirm letter needed.",
			},
			"status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate status name.",
			},
			"status_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate status message.",
			},
			"product_zh_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product zh name.",
			},
			"vulnerability_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vulnerability status.",
			},
			"cert_begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate begin time.",
			},
			"cert_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate end time.",
			},
			"insert_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate insert time.",
			},
			"certificate_private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate private key.",
			},
			"certificate_public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate public key.",
			},
			"renewable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate renewable.",
			},
			"deployable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate deployable.",
			},
			"dv_auths": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DV certification information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dv_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication key.",
						},
						"dv_auth_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication value.",
						},
						"dv_auth_verify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication type.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudFreeCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_free_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := SSLService{client}
	id := d.Id()

	request := ssl2.NewDescribeCertificateDetailRequest()
	request.CertificateId = &id

	response, err := service.DescribeCertificateDetail(ctx, request)

	if err != nil {
		d.SetId("")
		return err
	}

	detail := response.Response

	d.SetId(id)

	if detail.VerifyType != nil {
		_ = d.Set("dv_auth_method", detail.VerifyType)
	}

	if detail.Domain != nil {
		_ = d.Set("domain", detail.Domain)
	}
	if detail.PackageType != nil {
		_ = d.Set("package_type", detail.PackageType)
	}
	if detail.ValidityPeriod != nil {
		_ = d.Set("validity_period", detail.ValidityPeriod)
	}
	if detail.Alias != nil {
		_ = d.Set("alias", detail.Alias)
	}

	if detail.ProductZhName != nil {
		_ = d.Set("product_zh_name", detail.ProductZhName)
	}
	if detail.Status != nil {
		_ = d.Set("status", detail.Status)
	}
	if detail.StatusMsg != nil {
		_ = d.Set("status_msg", detail.StatusMsg)
	}
	if detail.VulnerabilityStatus != nil {
		_ = d.Set("vulnerability_status", detail.VulnerabilityStatus)
	}
	if detail.CertBeginTime != nil {
		_ = d.Set("cert_begin_time", detail.CertBeginTime)
	}
	if detail.CertEndTime != nil {
		_ = d.Set("cert_end_time", detail.CertEndTime)
	}
	if detail.InsertTime != nil {
		_ = d.Set("insert_time", detail.InsertTime)
	}
	if detail.CertificatePrivateKey != nil {
		_ = d.Set("certificate_private_key", detail.CertificatePrivateKey)
	}
	if detail.CertificatePublicKey != nil {
		_ = d.Set("certificate_public_key", detail.CertificatePublicKey)
	}
	if detail.StatusName != nil {
		_ = d.Set("status_name", detail.StatusName)
	}
	if detail.RenewAble != nil {
		_ = d.Set("renewable", detail.RenewAble)
	}
	if detail.Deployable != nil {
		_ = d.Set("deployable", detail.Deployable)
	}

	if detail.ProjectId != nil {
		pid, err := strconv.Atoi(*detail.ProjectId)
		if err != nil {
			return err
		}
		_ = d.Set("project_id", pid)
	}
	if detail.DvAuthDetail != nil && len(detail.DvAuthDetail.DvAuths) != 0 {
		dvAuths := make([]map[string]string, 0)
		for _, item := range detail.DvAuthDetail.DvAuths {
			dvAuth := make(map[string]string)
			dvAuth["dv_auth_key"] = *item.DvAuthKey
			dvAuth["dv_auth_value"] = *item.DvAuthValue
			dvAuth["dv_auth_verify_type"] = *item.DvAuthVerifyType
			dvAuths = append(dvAuths, dvAuth)
		}

		_ = d.Set("dv_auths", dvAuths)
	}

	return nil
}

func resourceTencentCloudFreeCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_free_certificate.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := SSLService{client}

	var (
		authMethod = d.Get("dv_auth_method").(string)
		domain     = d.Get("domain").(string)
	)

	request := ssl2.NewApplyCertificateRequest()

	request.DvAuthMethod = &authMethod
	request.DomainName = &domain

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("contact_email"); ok {
		request.ContactEmail = helper.String(v.(string))
	}
	if v, ok := d.GetOk("contact_phone"); ok {
		request.ContactPhone = helper.String(v.(string))
	}
	if v, ok := d.GetOk("validity_period"); ok {
		request.ValidityPeriod = helper.String(v.(string))
	}
	if v, ok := d.GetOk("csr_encrypt_algo"); ok {
		request.CsrEncryptAlgo = helper.String(v.(string))
	}
	if v, ok := d.GetOk("csr_key_parameter"); ok {
		request.CsrKeyParameter = helper.String(v.(string))
	}
	if v, ok := d.GetOk("csr_key_password"); ok {
		request.CsrKeyPassword = helper.String(v.(string))
	}
	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok && v.(int) > 0 {
		request.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("old_certificate_id"); ok {
		request.OldCertificateId = helper.String(v.(string))
	}

	id, err := service.ApplyCertificate(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudFreeCertificateRead(d, meta)
}

func resourceTencentCloudFreeCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_free_certificate.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := SSLService{client}

	id := d.Id()

	if d.HasChange("alias") {
		request := ssl2.NewModifyCertificateAliasRequest()
		alias, ok := d.GetOk("alias")
		if ok {
			request.Alias = helper.String(alias.(string))
		} else {
			request.Alias = helper.String("")
		}
		request.CertificateId = &id
		if err := service.ModifyCertificateAlias(ctx, request); err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		request := ssl2.NewModifyCertificateProjectRequest()
		pid, ok := d.GetOk("project_id")
		if ok {
			request.ProjectId = helper.IntUint64(pid.(int))
		} else {
			request.ProjectId = helper.IntUint64(-1)
		}
		request.CertificateIdList = []*string{&id}
		if err := service.ModifyCertificateProject(ctx, request); err != nil {
			return err
		}
	}

	immutableFields := []string{
		"dv_auth_method",
		"package_type",
		"contact_email",
		"contact_phone",
		"validity_period",
		"csr_encrypt_algo",
		"csr_key_parameter",
		"csr_key_password",
	}
	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("cannot update argument `%s`, please reset to previous value or replace by creating a new resource", f)
		}
	}

	return resourceTencentCloudFreeCertificateRead(d, meta)
}

func resourceTencentCloudFreeCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_free_certificate.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := SSLService{client}
	id := d.Id()

	request := ssl2.NewDeleteCertificateRequest()
	request.CertificateId = &id

	_, err := service.DeleteCertificate(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
