load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "cat_dario_mergo",
        importpath = "dario.cat/mergo",
        sum = "h1:85+piFYR1tMbRrLcDwR18y4UKJ3aH1Tbzi24VRW1TK8=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_apparentlymart_go_textseg_v13",
        importpath = "github.com/apparentlymart/go-textseg/v13",
        sum = "h1:Y+KvPE1NYz0xl601PVImeQfFyEy6iT90AvPUL1NNfNw=",
        version = "v13.0.0",
    )
    go_repository(
        name = "com_github_apparentlymart_go_textseg_v15",
        importpath = "github.com/apparentlymart/go-textseg/v15",
        sum = "h1:uYvfpb3DyLSCGWnctWKGj857c6ew1u1fNQOlOtuGxQY=",
        version = "v15.0.0",
    )
    go_repository(
        name = "com_github_armon_go_radix",
        importpath = "github.com/armon/go-radix",
        sum = "h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_bgentry_speakeasy",
        importpath = "github.com/bgentry/speakeasy",
        sum = "h1:tgObeVOf8WAvtuAX6DhJ4xks4CFNwPDZiqzGqIHE51E=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_bmatcuk_doublestar_v4",
        importpath = "github.com/bmatcuk/doublestar/v4",
        sum = "h1:X8jg9rRZmJd4yRy7ZeNDRnM+T3ZfHv15JiBJ/avrEXE=",
        version = "v4.9.1",
    )
    go_repository(
        name = "com_github_bufbuild_protocompile",
        importpath = "github.com/bufbuild/protocompile",
        sum = "h1:iA73zAf/fyljNjQKwYzUHD6AD4R8KMasmwa/FBatYVw=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:W5quZX/G/csjUnuI8SUYlsHs9M38FC7znL0lIO+DvMg=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_bwesterb_go_ristretto",
        importpath = "github.com/bwesterb/go-ristretto",
        sum = "h1:1w53tCkGhCQ5djbat3+MH0BAQ5Kfgbt56UZQ/JMzngw=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_clipperhouse_stringish",
        importpath = "github.com/clipperhouse/stringish",
        sum = "h1:+NSqMOr3GR6k1FdRhhnXrLfztGzuG+VuFDfatpWHKCs=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_clipperhouse_uax29_v2",
        importpath = "github.com/clipperhouse/uax29/v2",
        sum = "h1:SNdx9DVUqMoBuBoW3iLOj4FQv3dN5mDtuqwuhIGpJy4=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_cloudflare_circl",
        importpath = "github.com/cloudflare/circl",
        sum = "h1:zqIqSPIndyBh1bjLVVDHMPpVKqp8Su/V+6MeDzzQBQ0=",
        version = "v1.6.1",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:aQ3y1lwWyqYPiWZThqv1aFbZMiM9vblcSArJRf2Irls=",
        version = "v0.0.0-20250501225837-2ac532fd4443",
    )
    go_repository(
        name = "com_github_cyphar_filepath_securejoin",
        importpath = "github.com/cyphar/filepath-securejoin",
        sum = "h1:JyxxyPEaktOD+GAnqIqTf9A8tHyAG22rowi7HkoSU1s=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=",
        version = "v1.1.2-0.20180830191138-d8f796af33cc",
    )
    go_repository(
        name = "com_github_emirpasic_gods",
        importpath = "github.com/emirpasic/gods",
        sum = "h1:FXtiHYKDGKCW2KzwZKx0iC0PQmdlorYgdFG9jPXJ1Bc=",
        version = "v1.18.1",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:zEqyPVyku6IvWCFwux4x9RxkLOMUL+1vC9xUFv5l2/M=",
        version = "v0.13.4",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:jb83lalDRZSpPWW2Z7Mck/8kXZ5CQAFYVjQcdVIr83A=",
        version = "v1.32.4",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_ratelimit",
        importpath = "github.com/envoyproxy/go-control-plane/ratelimit",
        sum = "h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:DEo3O99U8j4hBFwbJfrz9VtgcDfUKS7KJ7spH3d86P8=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_fatih_color",
        importpath = "github.com/fatih/color",
        sum = "h1:S8gINlzdQ840/4pfAwic/ZE0djQEH3wM94VfqLTZcOM=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_frankban_quicktest",
        importpath = "github.com/frankban/quicktest",
        sum = "h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=",
        version = "v1.14.6",
    )
    go_repository(
        name = "com_github_go_git_gcfg",
        importpath = "github.com/go-git/gcfg",
        sum = "h1:+zs/tPmkDkHx3U66DAb0lQFJrpS6731Oaa12ikc+DiI=",
        version = "v1.5.1-0.20230307220236-3a3c6141e376",
    )
    go_repository(
        name = "com_github_go_git_go_billy_v5",
        importpath = "github.com/go-git/go-billy/v5",
        sum = "h1:6Q86EsPXMa7c3YZ3aLAQsMA0VlWmy43r6FHqa/UNbRM=",
        version = "v5.6.2",
    )
    go_repository(
        name = "com_github_go_git_go_git_v5",
        importpath = "github.com/go-git/go-git/v5",
        sum = "h1:/MD3lCrGjCen5WfEAzKg00MJJffKhC8gzS80ycmCi60=",
        version = "v5.14.0",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:TK/7NqRQZfgAh+Td8AlsrvtPoUyiHh0LqVvokh+1vHI=",
        version = "v4.1.2",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=",
        version = "v1.4.3",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:DrW6hGnjIhtvhOIiAKT6Psh/Kd/ldepEa81DKeiRJ5I=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:f+oWsMOmNPc8JmEHVZIycC7hBoQxHH9pNKQORJNozsQ=",
        version = "v0.0.0-20241129210726-2c02b8208cf8",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_detectors_gcp",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp",
        sum = "h1:UQUsRi8WTzhZntp5313l+CHIAT95ojUI2lpP/ExlZa4=",
        version = "v1.29.0",
    )
    go_repository(
        name = "com_github_hashicorp_cli",
        importpath = "github.com/hashicorp/cli",
        sum = "h1:/fZJ+hNdwfTSfsxMBa9WWMlfjUZbX8/LnUxgAd7lCVU=",
        version = "v1.1.7",
    )
    go_repository(
        name = "com_github_hashicorp_errwrap",
        importpath = "github.com/hashicorp/errwrap",
        sum = "h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_checkpoint",
        importpath = "github.com/hashicorp/go-checkpoint",
        sum = "h1:MFYpPZCnQqQTE18jFwSII6eUQrD/oxMFp3mlgcqk5mU=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_cleanhttp",
        importpath = "github.com/hashicorp/go-cleanhttp",
        sum = "h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_hclog",
        importpath = "github.com/hashicorp/go-hclog",
        sum = "h1:Qr2kF+eVWjTiYmU7Y31tYlP1h0q/X3Nl3tPGdaB11/k=",
        version = "v1.6.3",
    )
    go_repository(
        name = "com_github_hashicorp_go_multierror",
        importpath = "github.com/hashicorp/go-multierror",
        sum = "h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_plugin",
        importpath = "github.com/hashicorp/go-plugin",
        sum = "h1:YghfQH/0QmPNc/AZMTFE3ac8fipZyZECHdDPshfk+mA=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_retryablehttp",
        importpath = "github.com/hashicorp/go-retryablehttp",
        sum = "h1:ylXZWnqa7Lhqpk0L1P1LzDtGcCR0rPVUrx/c8Unxc48=",
        version = "v0.7.8",
    )
    go_repository(
        name = "com_github_hashicorp_go_uuid",
        importpath = "github.com/hashicorp/go-uuid",
        sum = "h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_hashicorp_go_version",
        importpath = "github.com/hashicorp/go-version",
        sum = "h1:5tqGy27NaOTB8yJKUZELlFAS/LTKJkrmONwQKeRZfjY=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_hashicorp_hc_install",
        importpath = "github.com/hashicorp/hc-install",
        sum = "h1:v80EtNX4fCVHqzL9Lg/2xkp62bbvQMnvPQ0G+OmtO24=",
        version = "v0.9.2",
    )
    go_repository(
        name = "com_github_hashicorp_logutils",
        importpath = "github.com/hashicorp/logutils",
        sum = "h1:dLEQVugN8vlakKOUE3ihGLTZJRB4j+M2cdTm/ORI65Y=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_exec",
        importpath = "github.com/hashicorp/terraform-exec",
        sum = "h1:mL0xlk9H5g2bn0pPF6JQZk5YlByqSqrO5VoaNtAf8OE=",
        version = "v0.24.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_json",
        importpath = "github.com/hashicorp/terraform-json",
        sum = "h1:BwGuzM6iUPqf9JYM/Z4AF1OJ5VVJEEzoKST/tRDBJKU=",
        version = "v0.27.2",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_plugin_docs",
        importpath = "github.com/hashicorp/terraform-plugin-docs",
        sum = "h1:YNZYd+8cpYclQyXbl1EEngbld8w7/LPOm99GD5nikIU=",
        version = "v0.24.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_plugin_framework",
        importpath = "github.com/hashicorp/terraform-plugin-framework",
        sum = "h1:1+zwFm3MEqd/0K3YBB2v9u9DtyYHyEuhVOfeIXbteWA=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_plugin_go",
        importpath = "github.com/hashicorp/terraform-plugin-go",
        sum = "h1:1nXKl/nSpaYIUBU1IG/EsDOX0vv+9JxAltQyDMpq5mU=",
        version = "v0.29.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_plugin_log",
        importpath = "github.com/hashicorp/terraform-plugin-log",
        sum = "h1:i7hOA+vdAItN1/7UrfBqBwvYPQ9TFvymaRGZED3FCV0=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_registry_address",
        importpath = "github.com/hashicorp/terraform-registry-address",
        sum = "h1:S1yCGomj30Sao4l5BMPjTGZmCNzuv7/GDTDX99E9gTk=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_hashicorp_terraform_svchost",
        importpath = "github.com/hashicorp/terraform-svchost",
        sum = "h1:EZZimZ1GxdqFRinZ1tpJwVxxt49xc/S52uzrw4x0jKQ=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_hashicorp_yamux",
        importpath = "github.com/hashicorp/yamux",
        sum = "h1:XtB8kyFOyHXYVFnwT5C3+Bdo8gArse7j2AQ0DA0Uey8=",
        version = "v0.1.2",
    )
    go_repository(
        name = "com_github_huandu_xstrings",
        importpath = "github.com/huandu/xstrings",
        sum = "h1:2ag3IFq9ZDANvthTwTiqSSZLjDc+BedvHPAp5tJy2TI=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_imdario_mergo",
        importpath = "github.com/imdario/mergo",
        sum = "h1:M8XP7IuFNsqUx6VPK2P9OSmsYsI/YFaGil0uD21V3dM=",
        version = "v0.3.15",
    )
    go_repository(
        name = "com_github_jbenet_go_context",
        importpath = "github.com/jbenet/go-context",
        sum = "h1:BQSFePA1RWJOlocH6Fxy8MmwDt+yVQYULKfN0RoTN8A=",
        version = "v0.0.0-20150711004518-d14ea06fba99",
    )
    go_repository(
        name = "com_github_jhump_protoreflect",
        importpath = "github.com/jhump/protoreflect",
        sum = "h1:qOEr613fac2lOuTgWN4tPAtLL7fUSbuJL5X5XumQh94=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_github_kevinburke_ssh_config",
        importpath = "github.com/kevinburke/ssh_config",
        sum = "h1:x584FjTGwHzMwvHx18PXxbBVzfnxogHaAReU4gf13a4=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_kunde21_markdownfmt_v3",
        importpath = "github.com/Kunde21/markdownfmt/v3",
        sum = "h1:KiZu9LKs+wFFBQKhrZJrFZwtLnCCWJahL+S+E/3VnM0=",
        version = "v3.1.0",
    )
    go_repository(
        name = "com_github_masterminds_goutils",
        importpath = "github.com/Masterminds/goutils",
        sum = "h1:5nUrii3FMTL5diU80unEVvNevw1nH4+ZV4DSLVJLSYI=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_masterminds_semver_v3",
        importpath = "github.com/Masterminds/semver/v3",
        sum = "h1:Zog+i5UMtVoCU8oKka5P7i9q9HgrJeGzI9SA1Xbatp0=",
        version = "v3.4.0",
    )
    go_repository(
        name = "com_github_masterminds_sprig_v3",
        importpath = "github.com/Masterminds/sprig/v3",
        sum = "h1:mQh0Yrg1XPo6vjYXgtf5OtijNAKJRNcTdOOGZe3tPhs=",
        version = "v3.3.0",
    )
    go_repository(
        name = "com_github_mattn_go_colorable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:9A9LHSqF/7dyVVX6g0U9cwm9pG3kP9gSzcuIPHPsaIE=",
        version = "v0.1.14",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=",
        version = "v0.0.20",
    )
    go_repository(
        name = "com_github_mattn_go_runewidth",
        importpath = "github.com/mattn/go-runewidth",
        sum = "h1:v++JhqYnZuu5jSKrk9RbgF5v4CGUjqRfBm05byFGLdw=",
        version = "v0.0.19",
    )
    go_repository(
        name = "com_github_microsoft_go_winio",
        importpath = "github.com/Microsoft/go-winio",
        sum = "h1:F2VQgta7ecxGYO8k3ZZz3RS8fVIXVxONVUPlNERoyfY=",
        version = "v0.6.2",
    )
    go_repository(
        name = "com_github_mitchellh_copystructure",
        importpath = "github.com/mitchellh/copystructure",
        sum = "h1:vpKXTN4ewci03Vljg/q9QvCGUDttBOGBIa15WveJJGw=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_mitchellh_go_testing_interface",
        importpath = "github.com/mitchellh/go-testing-interface",
        sum = "h1:jrgshOhYAUVNMAJiKbEu7EqAwgJJ2JqpQmpLJOu07cU=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_github_mitchellh_reflectwalk",
        importpath = "github.com/mitchellh/reflectwalk",
        sum = "h1:G2LzWKi524PWgd3mLHV8Y5k7s6XUvT0Gef6zxSIeXaQ=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_oklog_run",
        importpath = "github.com/oklog/run",
        sum = "h1:O8x3yXwah4A73hJdlrwo/2X6J62gE5qTMusH0dvz60E=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_pjbgf_sha1cd",
        importpath = "github.com/pjbgf/sha1cd",
        sum = "h1:a9wb0bp1oC2TGwStyn0Umc/IGKQnEgF0vVaZ8QF8eo4=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_pkg_diff",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )
    go_repository(
        name = "com_github_planetscale_vtprotobuf",
        importpath = "github.com/planetscale/vtprotobuf",
        sum = "h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=",
        version = "v0.6.1-0.20240319094008-0393e58bdf10",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=",
        version = "v1.0.1-0.20181226105442-5d4384ee4fb2",
    )
    go_repository(
        name = "com_github_posener_complete",
        importpath = "github.com/posener/complete",
        sum = "h1:NP0eAhjcjImqslEwo/1hq7gpajME0fTLTezBKDqfXqo=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_protonmail_go_crypto",
        importpath = "github.com/ProtonMail/go-crypto",
        sum = "h1:ILq8+Sf5If5DCpHQp4PbZdS1J7HDFRXz/+xKBiRGFrw=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:UQB4HGPB6osV0SQTLymcB4TgvyWu6ZyliaW0tI/otEQ=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_github_sebdah_goldie",
        importpath = "github.com/sebdah/goldie",
        sum = "h1:9GNhIat69MSlz/ndaBg48vl9dF5fI+NBB6kfOxgfkMc=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_sergi_go_diff",
        importpath = "github.com/sergi/go-diff",
        sum = "h1:n661drycOFuPLCN3Uc8sB6B/s6Z4t2xvBgU1htSHuq8=",
        version = "v1.3.2-0.20230802210424-5b0b94c5c0d3",
    )
    go_repository(
        name = "com_github_shopspring_decimal",
        importpath = "github.com/shopspring/decimal",
        sum = "h1:bxl37RwXBklmTi0C79JfXCEBD1cqqHt0bbgBAGFp81k=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_skeema_knownhosts",
        importpath = "github.com/skeema/knownhosts",
        sum = "h1:X2osQ+RAjK76shCbvhHHHVl3ZlgDm8apHEHFqRjnBY8=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:h2x0u2shc1QuLHfxi+cTJvs30+ZAHOGRic8uyGTDWxY=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:N2I01KCUkv1FAjZXJMwh95KK1ZIQLYbPfhaxw8WS0hE=",
        version = "v2.5.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:Xv5erBjTwe/5IxqUQTdXv5kgmIvbHo3QQyRwhJsOfJA=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_vmihailenco_msgpack_v5",
        importpath = "github.com/vmihailenco/msgpack/v5",
        sum = "h1:cQriyiUvjTwOHg8QZaPihLWeRAAVoCpE00IUPn0Bjt8=",
        version = "v5.4.1",
    )
    go_repository(
        name = "com_github_vmihailenco_tagparser_v2",
        importpath = "github.com/vmihailenco/tagparser/v2",
        sum = "h1:y09buUbR+b5aycVFQs/g70pqKVZNBmxwAhO7/IwNM9g=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_xanzy_ssh_agent",
        importpath = "github.com/xanzy/ssh-agent",
        sum = "h1:+/15pJfg/RsTxqYcX6fHqOXZwwMP+2VyYWJeWM2qQFM=",
        version = "v0.3.3",
    )
    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:GPddIs617DnBLFFVJFgpo1aBfe/4xcvMc3SB5t/D0pA=",
        version = "v1.7.13",
    )
    go_repository(
        name = "com_github_yuin_goldmark_meta",
        importpath = "github.com/yuin/goldmark-meta",
        sum = "h1:pWw+JLHGZe8Rk0EGsMVssiNb/AaPMHfSRszZeUeiOUc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_zclconf_go_cty",
        importpath = "github.com/zclconf/go-cty",
        sum = "h1:seZvECve6XX4tmnvRzWtJNHdscMtYEx5R7bnnVyd/d0=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_github_zclconf_go_cty_debug",
        importpath = "github.com/zclconf/go-cty-debug",
        sum = "h1:FosyBZYxY34Wul7O/MSKey3txpPYyCqVO5ZyceuQJEI=",
        version = "v0.0.0-20191215020915-b22d67c1ba0b",
    )
    go_repository(
        name = "com_github_zeebo_errs",
        importpath = "github.com/zeebo/errs",
        sum = "h1:XNdoD/RRMKP7HD0UhJnIzUy74ISdGGxURlYG8HSWSfM=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:PBWF+iiAerVNe8UCHxdOt6eHLVc3ydFeOCw78U8ytSU=",
        version = "v0.7.0",
    )
    go_repository(
        name = "dev_abhg_go_goldmark_frontmatter",
        importpath = "go.abhg.dev/goldmark/frontmatter",
        sum = "h1:P8kPG0YkL12+aYk2yU3xHv4tcXzeVnN+gU0tJ5JnxRw=",
        version = "v0.2.0",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:56OvJKSH3hDGL0ml5uSxZmz3/3Pq4tJ+fb1unVLAFcY=",
        version = "v0.24.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
        version = "v0.0.0-20161208181325-20d25e280405",
    )
    go_repository(
        name = "in_gopkg_warnings_v0",
        importpath = "gopkg.in/warnings.v0",
        sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
        version = "v0.1.2",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:cH53jehLUN6UFLY71z+NDOiNJqDdPRaXzTel0sJySYA=",
        version = "v1.1.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:F7q2tNlCaHY9nMKHR6XH9/qkp8FktLnIcy6jJNyOCQw=",
        version = "v1.36.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:9zhNfelUvx0KBfu/gb+ZgeAfAgtWrfHJZcAqFC228wQ=",
        version = "v1.37.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:mvwbQS5m0tbmqML4NqK+e3aDiO02vsf/WgbsdpcPoZE=",
        version = "v1.37.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:ItB0QUqnjesGRvNcmAcU0LyvkVyGJ2xftD29bWdDvKI=",
        version = "v1.37.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:90lI228XrB9jCMuSdA0673aubgRobVZFhbjxHHspCPc=",
        version = "v1.37.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:HLdcFNbRQBE2imdSEgm/kwqmQj1Or1l/7bW6mxVK7z4=",
        version = "v1.37.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:ULiyYQ0FdsJhwwZUwbaXpZF5yUE3h+RA+gxvBu37ucc=",
        version = "v0.0.0-20250804133106-a7a43d27e69b",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:tRPGkdGHuewF4UisLzzHHr1spKw92qLM98nIzxbC0wY=",
        version = "v0.0.0-20251103181224-f26f9409b101",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:UnVkv1+uMLYXoIz6o7chp59WfQUYA2ex/BXQ9rHZu7A=",
        version = "v1.76.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=",
        version = "v1.36.10",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:dduJYIi3A3KOfdGOHX8AVZ/jGiyPa3IbBozJ5kNuE04=",
        version = "v0.43.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:mgKeJMpvi0yx/sU5GsxQ7p6s2wtOnGAHZWCHUM4KGzY=",
        version = "v0.0.0-20251023183803-a4bb9ffd2546",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:HV8lRxZC4l2cr3Zq1LvtOsi/ThTgWnUk/y64QSs8GwA=",
        version = "v0.29.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:giFlY12I07fugqwPuWJi68oOnpfqFnJIJzaIIm2JVV4=",
        version = "v0.46.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:dnDm7JmhM45NNpd8FDDeLhK6FwqbOf4MLCM9zb1BOHI=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:l60nONMj9l5drqw6jlhIELNv9I0A4OFgRsG9k2oT9Ug=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=",
        version = "v0.37.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:zMPR+aF8gfksFprF/Nc/rd1wRS1EI6nDBGyWAvDzx2Q=",
        version = "v0.36.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:Hx2Xv8hISq8Lm16jvBZ2VQf+RLmbd7wVUsALibYI/IQ=",
        version = "v0.38.0",
    )
    go_repository(
        name = "org_golang_x_tools_go_expect",
        importpath = "golang.org/x/tools/go/expect",
        sum = "h1:jpBZDwmgPhXsKZC6WhL20P4b/wmnpsEAGHaNy0n/rJM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_golang_x_tools_go_packages_packagestest",
        importpath = "golang.org/x/tools/go/packages/packagestest",
        sum = "h1:1h2MnaIAIXISqTFKdENegdpAgUXz6NrPEsbIeWaBRvM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=",
        version = "v0.16.0",
    )
