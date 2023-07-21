resource "tencentcloud_sms_sign" "example" {
  sign_name     = "tf_example_sms_sign"
  sign_type     = 1 # 1：APP,  DocumentType can be chosen（0，1，2，3，4）
  document_type = 4 # Screenshot of application background management (personally developed APP)
  international = 0 # Mainland China SMS
  sign_purpose  = 0 # personal use
  proof_image   = "your_proof_image"
}
