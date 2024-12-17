## x.x.x (Unreleased)

## 1.2.2 (17 Dec, 2024)

BUG FIXES:

- provider: Authentication failed on APISIX 3.11.0 ([#20](https://github.com/rework-space-com/terraform-provider-apisix/issues/20))

ENHANCEMENTS:

- provider: Upgrade provider to support APISIX v3.11.0 ([#20](https://github.com/rework-space-com/terraform-provider-apisix/issues/20))
- provider: Fix vulnerability issues by updating the terraform-provider-framework ([#21](https://github.com/rework-space-com/terraform-provider-apisix/issues/22))

## 1.2.1 (22 Mar, 2024)

BUG FIXES:

- resource/apisix_route: Can't create a route that uses `plugins` and `plugin_config_id` attributest simultaneously ([#15](https://github.com/rework-space-com/terraform-provider-apisix/issues/15))

## 1.2.0 (29 Feb, 2024)

BREAKING CHANGES:

- resource/apisix_ssl_certificate: The `validity_end` and `validity_start` attributes have been removed due to the APISIX v3.6.0 compatibility issue. (https://github.com/rework-space-com/terraform-provider-apisix/issues/10)

BUG FIXES:

- resource/apisix_ssl_certificate: Fix compatibility issue with resource `apisix_ssl_certificate` (https://github.com/rework-space-com/terraform-provider-apisix/issues/10)

ENHANCEMENTS:

- provider: Upgrade provider to support APISIX v3.8.0 (https://github.com/rework-space-com/terraform-provider-apisix/issues/6)
- provider: Fix vulnerability issues by updating the terraform-provider-framework (https://github.com/rework-space-com/terraform-provider-apisix/issues/12)

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
