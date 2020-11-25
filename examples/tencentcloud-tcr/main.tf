resource "tencentcloud_tcr_instance" "example" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_token" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  description       = "test"
  enable   = false
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "test"
  is_public   = false
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  namespace_name        = tencentcloud_tcr_namespace.example.name
  name = "test"
  brief_desc = "example"
  description = "long example"
}

data "tencentcloud_tcr_instances" "example" {
  name = tencentcloud_tcr_instance.example.name
}

data "tencentcloud_tcr_tokens" "example" {
  instance_id = tencentcloud_tcr_token.example.instance_id
}

data "tencentcloud_tcr_namespaces" "example" {
  instance_id = tencentcloud_tcr_namespace.example.instance_id
}

data "tencentcloud_tcr_repositories" "example" {
  instance_id = tencentcloud_tcr_repository.example.instance_id
  namespace_name = tencentcloud_tcr_namespace.example.name
}