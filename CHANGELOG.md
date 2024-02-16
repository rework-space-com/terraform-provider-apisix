## 2.0.0 (Unreleased)

## 1.1.2 (13 Feb, 2024)

BUG FIXES:

- resource/apisix_upstream: Updating upstream's "service_name" requires recreation of upstream ([#5](https://github.com/rework-space-com/terraform-provider-apisix/issues/5))

## 1.1.1 (13 Feb, 2024)

NOTES:

- resource/apisix_ssl_certificate: The `validity_start` and `validity_end` attributes have been deprecated and will be removed in the next major version of the provider ([#7](https://github.com/rework-space-com/terraform-provider-apisix/issues/7))

## 1.1.0 (4 Jul, 2023)

FEATURES:

- resource/apisix_ssl_certificate: Added calculation of notAfter and notBefore dates for the SSL certificate ([#4](https://github.com/rework-space-com/terraform-provider-apisix/pull/4))

## 1.0.0 (28 Jun, 2023)
Initial release
