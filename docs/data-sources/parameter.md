---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nullplatform_parameter Data Source - nullplatform"
subcategory: ""
description: |-
  Provides information about the Parameter by ID.
---

# nullplatform_parameter (Data Source)

Provides information about the Parameter by ID.

## Example Usage

```terraform
terraform {
  required_providers {
    nullplatform = {
      source = "nullplatform/nullplatform"
    }
  }
}

data "nullplatform_parameter" "example" {
  id = "123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) A system-wide unique ID representing the resource.

### Read-Only

- `destination_path` (String) The full path for file. Required when `type = file`.
- `encoding` (String) Possible values: [`plaintext`, `base64`]
- `name` (String) Definition name of the variable.
- `nrn` (String) The NRN of the application to which the parameter belongs to.
- `read_only` (Boolean) `true` if the value is a secret, `false` otherwise
- `secret` (Boolean) `true` if the value is a secret, `false` otherwise
- `type` (String) Possible values: [`environment`, `file`]
- `variable` (String) The name of the environment variable. Required when `type = environment`.