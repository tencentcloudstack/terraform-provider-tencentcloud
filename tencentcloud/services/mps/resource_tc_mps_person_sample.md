Provides a resource to create a mps person_sample

Example Usage

```hcl
resource "tencentcloud_mps_person_sample" "person_sample" {
  name          = "test"
  usages        = [
    "Review.Face"
  ]
  description   = "test"
  face_contents = [
    filebase64("./person.png")
  ]
}
```

Import

mps person_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_person_sample.person_sample person_sample_id
```