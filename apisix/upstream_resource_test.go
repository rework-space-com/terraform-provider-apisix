package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUpstreamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	name = "Example"
	desc = "Example of the upstream resource usage"
	type = "roundrobin"
	labels = {
		version = "v1"
	}
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 1980
			weight = 1
		},
		{
			host   = "127.0.0.1"
			port   = 1970
			weight = 1
		},
	]
	keepalive_pool = {
		idle_timeout = 5
		requests     = 10
		size         = 15
	}
	checks = {
		active = {
			host      = "example.com"
			port      = 8888
			timeout   = 5
			http_path = "/status"
			healthy = {
				interval  = 2,
				successes = 1
			}
			unhealthy = {
				interval      = 1
				http_failures = 2
			}
		}
		passive = {
			healthy = {
				http_statuses = [200, 201]
			}
			unhealthy = {
				http_statuses = [500]
				http_failures = 3
				tcp_failures  = 3
			}
		}
	}
	tls = {
		client_cert = "-----BEGIN CERTIFICATE-----\nMIICqzCCAZMCFFFD4RIOMzAQMEkV5Y9kzSxighFxMA0GCSqGSIb3DQEBCwUAMBMx\nETAPBgNVBAMMCE15Um9vdENBMB4XDTI1MDcxNzE0NDAyMVoXDTI2MDcxNzE0NDAy\nMVowETEPMA0GA1UEAwwGY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAw9aiY+b8Fpwsz36f/YHpx9v1t6o0qrZa/og9ecWaQ4di7w6LgxohBMK3\n6YO2cdaP3s7eQDiu3tujc0fZLLgcSmnavI3XW9bEQo0d2u/Y9a9CBr2ZNC8tTME5\nJB45cl/xyfEgHeMm7o1g+5g08qN5sCdJBEd+at9rWeRZ+xv/0KFeY0wcOO/0HZq7\nW6q3Ra0juBbsT5iNEes6eCFfrgg7PkjgGirSJJ/VlLEt6Tdbouvn1+ymxiRiWdGq\n7nresTryXYA7FqLHhQa7WBlTr6VzWNYEJ0n8qrCf3BbomscJhwCg7Hgu5e1jQwRg\nVLfJ7rmIgpBjuWAcu1aRKUU8EG/1jQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAa\nz/HAF4EhYBsrgLklqMbn4+1+7KL7rIsNoxApRWU14jIk/LVuLBbukuguzbe3xHgG\ncsQDa9Lt/KN48ZQ3NWwV/cjxpIMxIey8VOhezu2t+wj4BZOCvU3wHL8u32T/0Mg3\nbXyqJ8Wjs0bnNEZ70WOjBqT/I5vZgltxCBtdXo6c5VvlQVMEMXzQGP2pFjM5AQ54\nc4CpOJsIMyVU/h+Ou3I9baCYb/xAEG8cEm8EGSaY1LemOPaLhFt1yKqyMFNzTm1c\nK+1bfq9GQe/VUF7DHMI9PZAd30tQMqTUZet80UKmuVV6hxwazeMsYLeuvxx1NxxV\n3EvsyCUd1rolsBkp/6fw\n-----END CERTIFICATE-----\n"
		client_key = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDD1qJj5vwWnCzP\nfp/9genH2/W3qjSqtlr+iD15xZpDh2LvDouDGiEEwrfpg7Zx1o/ezt5AOK7e26Nz\nR9ksuBxKadq8jddb1sRCjR3a79j1r0IGvZk0Ly1MwTkkHjlyX/HJ8SAd4ybujWD7\nmDTyo3mwJ0kER35q32tZ5Fn7G//QoV5jTBw47/QdmrtbqrdFrSO4FuxPmI0R6zp4\nIV+uCDs+SOAaKtIkn9WUsS3pN1ui6+fX7KbGJGJZ0aruet6xOvJdgDsWoseFBrtY\nGVOvpXNY1gQnSfyqsJ/cFuiaxwmHAKDseC7l7WNDBGBUt8nuuYiCkGO5YBy7VpEp\nRTwQb/WNAgMBAAECggEAKWyQm+4je4rcZaWEpQRiVW6e+pMLoeKBu95IlqXoHAma\nsTNT6k7QFigz67Z7FHhMpVX/p/j1cFloKP3VH8Lv5QOgC4s7NwdmKyebXZCnRUyl\nfDSFoAasn9QtSIkGIL3PsKYK45eFSCdqkL1g0cQnfM3KgZe301Zf6DtHlziUc2YY\nDikg/zeQBmNheb2SehazA2Pry4L0ph9akDyEZfvBePtOzzxfjP4o4GHuYZBCn7tF\nTjJkmESINrDgHtvO2FO2cLf++XbDi38wRgKKKDol1JvuIGfyxCbTcu7kB+55HP6u\nqvJduEnhC3NDWJ8tYt0jgACg87sCJLbX2TYaugwdNwKBgQDgI28ISh95Mj9L2Cc5\niqYT+Ldm35YoFXUf3yIAJ0PBrdhAOfRccl39rZ6aCPwHnJR/jaDpqD+gTXiDriYX\nsD9AUeRiPK0sv2HaZ2PFr0fSymrPFogOo5v/4c9eFOBJMdqfVNdshNotvnwRt2Ou\nUMF193a3vkFkv9tfVQkMCko6MwKBgQDfrVehvDkCNcy706Q19unAmu32tOfKKBu2\nmx6pw50alWcEKRlU9OaqHzdQI4OTxfeUu1tf/lcY6c3N3cOv2lScDWRb8NXxTdNF\nsWqtQuNFvVt4Ds7RD2hTO551BXsuAAAgkYYiywJ4wnALi1KcpVZB0QYGDoeWHA6z\nCHOgfrjRPwKBgQDQKPhYaYabZ0gTpzaeoR6mk6m419O7fFofdHo+TDkIKe0ZkPlZ\n1jlmfJU8lzWB2DCt2ZnlBwW4Wdqf6N+lxmCn2qZReeqXEVLOpJCrqqL4qFbT5ygK\n+HXMCiotRRQbxjo1GXVMaoG6VBsj1P61iHhzl5ThBsfyyp/xBKd3UCMpswKBgQDC\nRHNP7YI2ATQIDhEZLZJnzifPld+bHKq1NpSzLUpNxGTsoCV4PBv6tZH88FtfBRm6\n+96oyOYspSQyIOHM4fuKbbc0gz/NjKJqbWURhn6OG6BN7c6ClLcvUyCU0mXh1e5G\nWx39KgTDjVxzKlZd5tu73ic1K3lnTocVx8llI6qxUQKBgDo8cYzssIb50vwkwpQS\nViriPkG7wKBZeFJDFO4zy4OSuXuMtsryizdHHdwV7ZNMDz42foOhmLuK+TnsNB6Q\nXJbTpM/2hXA0Lhw8vlvtj5VRuDMgB5YVWt5qJc0GxOn2GEDWHXMQ+JA4pSfPVwGD\ng4c5Tdvr8S6VcdVRaE+QGQdS\n-----END PRIVATE KEY-----\n"
	}
}		
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_upstream.test", "id"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "type", "roundrobin"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "pass_host", "pass"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "scheme", "http"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "hash_on", "vars"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_upstream.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore client_key during import as APISIX uses base64 format for it
				ImportStateVerifyIgnore: []string{"tls.client_key"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	name = "Example"
	desc = "Example of the upstream resource usage"
	type = "roundrobin"
	labels = {
		version = "v2"
	}
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 1980
			weight = 1
		},
		{
			host   = "127.0.0.1"
			port   = 1970
			weight = 5
		},
	]
	keepalive_pool = {
		idle_timeout = 10
		requests     = 10
		size         = 15
	}
	checks = {
		active = {
			host      = "example.com"
			port      = 8888
			timeout   = 5
			http_path = "/status"
			healthy = {
				interval  = 3,
				successes = 1
			}
			unhealthy = {
				interval      = 1
				http_failures = 2
			}
		}
	}
	tls = {
		client_cert = "-----BEGIN CERTIFICATE-----\nMIICqzCCAZMCFFFD4RIOMzAQMEkV5Y9kzSxighFxMA0GCSqGSIb3DQEBCwUAMBMx\nETAPBgNVBAMMCE15Um9vdENBMB4XDTI1MDcxNzE0NDAyMVoXDTI2MDcxNzE0NDAy\nMVowETEPMA0GA1UEAwwGY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAw9aiY+b8Fpwsz36f/YHpx9v1t6o0qrZa/og9ecWaQ4di7w6LgxohBMK3\n6YO2cdaP3s7eQDiu3tujc0fZLLgcSmnavI3XW9bEQo0d2u/Y9a9CBr2ZNC8tTME5\nJB45cl/xyfEgHeMm7o1g+5g08qN5sCdJBEd+at9rWeRZ+xv/0KFeY0wcOO/0HZq7\nW6q3Ra0juBbsT5iNEes6eCFfrgg7PkjgGirSJJ/VlLEt6Tdbouvn1+ymxiRiWdGq\n7nresTryXYA7FqLHhQa7WBlTr6VzWNYEJ0n8qrCf3BbomscJhwCg7Hgu5e1jQwRg\nVLfJ7rmIgpBjuWAcu1aRKUU8EG/1jQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAa\nz/HAF4EhYBsrgLklqMbn4+1+7KL7rIsNoxApRWU14jIk/LVuLBbukuguzbe3xHgG\ncsQDa9Lt/KN48ZQ3NWwV/cjxpIMxIey8VOhezu2t+wj4BZOCvU3wHL8u32T/0Mg3\nbXyqJ8Wjs0bnNEZ70WOjBqT/I5vZgltxCBtdXo6c5VvlQVMEMXzQGP2pFjM5AQ54\nc4CpOJsIMyVU/h+Ou3I9baCYb/xAEG8cEm8EGSaY1LemOPaLhFt1yKqyMFNzTm1c\nK+1bfq9GQe/VUF7DHMI9PZAd30tQMqTUZet80UKmuVV6hxwazeMsYLeuvxx1NxxV\n3EvsyCUd1rolsBkp/6fw\n-----END CERTIFICATE-----\n"
		client_key = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDD1qJj5vwWnCzP\nfp/9genH2/W3qjSqtlr+iD15xZpDh2LvDouDGiEEwrfpg7Zx1o/ezt5AOK7e26Nz\nR9ksuBxKadq8jddb1sRCjR3a79j1r0IGvZk0Ly1MwTkkHjlyX/HJ8SAd4ybujWD7\nmDTyo3mwJ0kER35q32tZ5Fn7G//QoV5jTBw47/QdmrtbqrdFrSO4FuxPmI0R6zp4\nIV+uCDs+SOAaKtIkn9WUsS3pN1ui6+fX7KbGJGJZ0aruet6xOvJdgDsWoseFBrtY\nGVOvpXNY1gQnSfyqsJ/cFuiaxwmHAKDseC7l7WNDBGBUt8nuuYiCkGO5YBy7VpEp\nRTwQb/WNAgMBAAECggEAKWyQm+4je4rcZaWEpQRiVW6e+pMLoeKBu95IlqXoHAma\nsTNT6k7QFigz67Z7FHhMpVX/p/j1cFloKP3VH8Lv5QOgC4s7NwdmKyebXZCnRUyl\nfDSFoAasn9QtSIkGIL3PsKYK45eFSCdqkL1g0cQnfM3KgZe301Zf6DtHlziUc2YY\nDikg/zeQBmNheb2SehazA2Pry4L0ph9akDyEZfvBePtOzzxfjP4o4GHuYZBCn7tF\nTjJkmESINrDgHtvO2FO2cLf++XbDi38wRgKKKDol1JvuIGfyxCbTcu7kB+55HP6u\nqvJduEnhC3NDWJ8tYt0jgACg87sCJLbX2TYaugwdNwKBgQDgI28ISh95Mj9L2Cc5\niqYT+Ldm35YoFXUf3yIAJ0PBrdhAOfRccl39rZ6aCPwHnJR/jaDpqD+gTXiDriYX\nsD9AUeRiPK0sv2HaZ2PFr0fSymrPFogOo5v/4c9eFOBJMdqfVNdshNotvnwRt2Ou\nUMF193a3vkFkv9tfVQkMCko6MwKBgQDfrVehvDkCNcy706Q19unAmu32tOfKKBu2\nmx6pw50alWcEKRlU9OaqHzdQI4OTxfeUu1tf/lcY6c3N3cOv2lScDWRb8NXxTdNF\nsWqtQuNFvVt4Ds7RD2hTO551BXsuAAAgkYYiywJ4wnALi1KcpVZB0QYGDoeWHA6z\nCHOgfrjRPwKBgQDQKPhYaYabZ0gTpzaeoR6mk6m419O7fFofdHo+TDkIKe0ZkPlZ\n1jlmfJU8lzWB2DCt2ZnlBwW4Wdqf6N+lxmCn2qZReeqXEVLOpJCrqqL4qFbT5ygK\n+HXMCiotRRQbxjo1GXVMaoG6VBsj1P61iHhzl5ThBsfyyp/xBKd3UCMpswKBgQDC\nRHNP7YI2ATQIDhEZLZJnzifPld+bHKq1NpSzLUpNxGTsoCV4PBv6tZH88FtfBRm6\n+96oyOYspSQyIOHM4fuKbbc0gz/NjKJqbWURhn6OG6BN7c6ClLcvUyCU0mXh1e5G\nWx39KgTDjVxzKlZd5tu73ic1K3lnTocVx8llI6qxUQKBgDo8cYzssIb50vwkwpQS\nViriPkG7wKBZeFJDFO4zy4OSuXuMtsryizdHHdwV7ZNMDz42foOhmLuK+TnsNB6Q\nXJbTpM/2hXA0Lhw8vlvtj5VRuDMgB5YVWt5qJc0GxOn2GEDWHXMQ+JA4pSfPVwGD\ng4c5Tdvr8S6VcdVRaE+QGQdS\n-----END PRIVATE KEY-----\n"
	}		
}				
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_upstream.test", "id"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "type", "roundrobin"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "pass_host", "pass"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "scheme", "http"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "hash_on", "vars"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
