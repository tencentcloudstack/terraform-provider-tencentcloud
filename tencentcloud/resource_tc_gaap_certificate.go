package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapCertificateCreate,
		Read:   resourceTencentCloudGaapCertificateRead,
		Update: resourceTencentCloudGaapCertificateUpdate,
		Delete: resourceTencentCloudGaapCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(0, 4),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			// computed
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"issuer_cn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_cn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudGaapCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	certificateType := d.Get("type").(int)
	content := d.Get("content").(string)

	name := d.Get("name").(string)

	var key *string
	if rawKey, ok := d.GetOk("key"); ok {
		key = stringToPointer(rawKey.(string))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.createCertificate(ctx, certificateType, content, name, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	certificate, err := service.DescribeCertificateById(ctx, id)
	if err != nil {
		return err
	}

	if certificate == nil {
		d.SetId("")
		return nil
	}

	if certificate.CertificateType == nil {
		return errors.New("certificate type is nil")
	}
	d.Set("type", certificate.CertificateType)

	if certificate.CertificateContent == nil {
		return errors.New("certificate content is nil")
	}
	d.Set("content", certificate.CertificateContent)

	if certificate.CertificateAlias == nil {
		return errors.New("certificate name is nil")
	}
	d.Set("name", certificate.CertificateAlias)

	if _, ok := d.GetOk("key"); ok {
		if certificate.CertificateKey == nil {
			return errors.New("certificate key is nil")
		}
		d.Set("key", certificate.CertificateKey)
	}

	if certificate.CreateTime == nil {
		return errors.New("certificate create time is nil")
	}
	d.Set("create_time", certificate.CreateTime)

	if certificate.BeginTime != nil {
		d.Set("begin_time", certificate.BeginTime)
	}
	if certificate.EndTime != nil {
		d.Set("end_time", certificate.EndTime)
	}
	if certificate.IssuerCN != nil {
		d.Set("issuer_cn", certificate.IssuerCN)
	}
	if certificate.SubjectCN != nil {
		d.Set("subject_cn", certificate.SubjectCN)
	}

	return nil
}

func resourceTencentCloudGaapCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	name := d.Get("name").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.ModifyCertificateName(ctx, id, name); err != nil {
		return err
	}

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteCertificate(ctx, id)
}
