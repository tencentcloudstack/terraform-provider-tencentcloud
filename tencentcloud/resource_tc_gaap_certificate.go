package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapCertificateCreate,
		Read:   resourceTencentCloudGaapCertificateRead,
		Update: resourceTencentCloudGaapCertificateUpdate,
		Delete: resourceTencentCloudGaapCertificateDelete,
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
	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	certificateType := d.Get("type").(int)
	content := d.Get("content").(string)

	var (
		name *string
		key  *string
	)

	if rawName, ok := d.GetOk("name"); ok {
		name = stringToPointer(rawName.(string))
	}

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
	logId := getLogId(nil)
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

	if certificate.CertificateId == nil {
		err := fmt.Errorf("certificate id is nil")
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}

	if certificate.CertificateType == nil {
		err := fmt.Errorf("certificate type is nil")
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}
	d.Set("type", certificate.CertificateType)

	if certificate.CertificateContent == nil {
		err := fmt.Errorf("certificate content is nil")
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}
	d.Set("content", certificate.CertificateContent)

	if _, ok := d.GetOk("name"); ok {
		if certificate.CertificateAlias == nil {
			err := fmt.Errorf("certificate name is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}
		d.Set("name", certificate.CertificateAlias)
	}

	if _, ok := d.GetOk("key"); ok {
		if certificate.CertificateKey == nil {
			err := fmt.Errorf("certificate key is nil")
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}
		d.Set("key", certificate.CertificateKey)
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
	logId := getLogId(nil)
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
	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteCertificate(ctx, id)
}
