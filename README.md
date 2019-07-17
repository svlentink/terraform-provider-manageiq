# Terraform-provider-ManageIQ

Status: development

Please fork this repo. and expand/improve it.

## build

```
docker-compose up
```

## run

required environment variables:
```
MANAGEIQ_USERNAME=ldapuser01
MANAGEIQ_PASSWORD=example123
MANAGEIQ_CONFIGFILE=~/manageiq_terraform_config.yaml
```

[configfile format](https://github.com/svlentink/terraform-provider-manageiq/blob/master/config_file_format.go)
with example values:
```yaml
cat << EOF > ~/manageiq_terraform_config.yaml
api_hostname: manageiq.example.com
order_resource_parameters:
  examplekey: example value
  tag_0_backup: yes
  tag_0_department: IT services
  subnet_name: prod_dmz_EU_01
  vm_memory: 2048
EOF
```

optional:
```
MANAGEIQ_INSECURE=true #ignores TLS warnings
```

## Links

- How to get logs: https://github.com/hashicorp/terraform/issues/16752
- https://github.com/ManageIQ/manageiq_docs/blob/master/doc-REST_API/topics/Available_Actions.adoc
- https://github.com/ManageIQ/manageiq_docs/tree/master/api/examples
- https://access.redhat.com/documentation/en-us/red_hat_cloudforms/4.7/html-single/red_hat_cloudforms_rest_api/index
- http://manageiq.org/docs/reference/fine/api/examples/queries
- https://liquidat.wordpress.com/2015/08/27/howto-accessing-cloudforms-3-2-via-rest-with-python/

## Tags

- CloudForms Management Engine - terraform-provider-cloudformsmanagementengine
- CFME - terraform-provider-cfme
- ManageIQ - terraform-provider-manageiq
