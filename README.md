[![test](https://github.com/ks6088ts/terraform-provider-soracom/workflows/test/badge.svg)](https://github.com/ks6088ts/terraform-provider-soracom/actions/workflows/test.yml)
[![release](https://github.com/ks6088ts/terraform-provider-soracom/workflows/release/badge.svg)](https://github.com/ks6088ts/terraform-provider-soracom/actions/workflows/release.yml)

# Terraform provider for SORACOM

_**Please note:** We take Terraform's security and our users' trust very seriously. If you believe you have found a security issue in the Terraform provider for SORACOM, please responsibly disclose it by contacting us at ks6088ts@gmail.com._

## Documentation

Documentation is available on the Terraform website: https://registry.terraform.io/providers/ks6088ts/soracom/latest/docs


## Get started

1. See [SORACOM CLI Getting Started Guide](https://developers.soracom.io/en/start/soracom/soracom-cli-guide/) to create configuration profiles for SORACOM CLI. By default, the SORACOM CLI configuration file will be stored in `~/.soracom/myprofile.json` on macOS/Linux or `C:\Users\<username>\.soracom\myprofile.json` on Windows.

2. Create following usage example.

```hcl
# Specify the version of the SORACOM Provider to use
terraform {
  required_providers {
    soracom = {
      source = "ks6088ts/soracom"
      version = "=0.0.2"
    }
  }
}

# Configure the SORACOM Provider
provider "soracom" {
  # set your own profile
  profile = "myprofile"
}

# Create resources (e.g. Create SIM group named as "my_sim_group")
resource "soracom_group" "group" {
  tags = {
    name = "my_sim_group"
    team = "soracom"
  }
}
```

3. run the following commands to manage resources

```bash
# Initialize a new or existing Terraform working directory
terraform init

# Creates or updates infrastructure according to Terraform configuration files in the current directory.
terraform apply

# Confirm SIM group is created via SORACOM CLI
soracom groups list
[
	{
		"configuration": {},
		"createdAt": 1666314306643,
		"createdTime": 1666314306643,
		"groupId": "3d3cc5f7-9924-41d6-bc5f-f5b73e6ffe05",
		"lastModifiedAt": 1666314306643,
		"lastModifiedTime": 1666314306643,
		"operatorId": "OP0076153716",
		"tags": {
			"name": "my_sim_group",
			"team": "soracom"
		}
	}
]

# Destroy Terraform-managed infrastructure.
terraform destroy

# Confirm SIM group is deleted via SORACOM CLI
soracom groups list
[]
```

## Contributing

### Prerequisites

- [Go](https://golang.org/) 1.19+
- [Git](https://git-scm.com/)
- [Terraform](https://developer.hashicorp.com/terraform/downloads)
- [GNU Make](https://www.gnu.org/software/make/)

### Build

Run the following from repository root to generate SORACOM provider in `plugins` directory.

```bash
make build
```

### Test

All the tests are done on GitHub Actions. Please see details here: [.github/workflows/test.yml](.github/workflows/test.yml).

#### Provider

To run tasks related to provider, see [Makefile](./Makefile) for details.
For example, you can run all the tests for provider by the following command.

```bash
make ci-test
```

#### Examples

To run tasks related to HCL examples, see [terraform.mk](./terraform.mk) for details.
For example, you can run all the tests for examples by the following command.

```bash
make -f terraform.mk ci-test-examples
```

#### Documents

To run tasks related to documents, see [terraform.mk](./terraform.mk) for details.
For example, you can run all the tests for documents by the following command.

```bash
make -f terraform.mk ci-test-docs
```

# References

- To understand what IaC is: [[AWS Black Belt Online Seminar] AWS Cloud Development Kit (CDK)](https://d1.awsstatic.com/webinars/jp/pdf/services/20200303_BlackBelt_CDK.pdf)
- To understand what Terraform is: [What is Terraform?](https://developer.hashicorp.com/terraform/intro)
- To understand how to implement Terraform Custom Provider: [Call APIs with Custom SDK Providers
](https://developer.hashicorp.com/terraform/tutorials/providers)
