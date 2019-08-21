# Terraform-provider-ManageIQ

Status: alpha

Ordering and deleting of a VM works, tested on CFME.

Please fork this repo. and expand/improve it.

## Running it

Run the following from the dir with your `.tf` files:
```
docker run -it --rm -v $PWD:/data \
  svlentink/terraform-provider-manageiq [init|plan|apply|etc.]
```
for the format see
[example.tf](https://github.com/svlentink/terraform-provider-manageiq/blob/master/example.tf)

## configuration

optional environment variables:
```
MANAGEIQ_USERNAME=ldapuser01
MANAGEIQ_PASSWORD=example123
```
which will set them as basic-auth values while doing the requests.

other optional:
```
MANAGEIQ_INSECURE=true #ignores TLS warnings
```
we might changed this to
[configuration options](https://learn.hashicorp.com/terraform/getting-started/variables#from-environment-variables)
as well.


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
