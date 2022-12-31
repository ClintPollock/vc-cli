# Overview

Help for using this with the pipeline scan functionality

Configure to do a mock scan with real time flaws:

```
veracode configure set urls.host stage-mock-pipeline-scan-system.era1.dev.vnext.veracode.io
veracode configure set urls.scheme https
veracode configure set urls.ssl-ignore true

veracode configure set urls.pipeline_api_path ''

veracode configure set mock.number-of-scanner-instances 1
veracode configure set mock.number-of-findings-to-throw 10
veracode configure set mock.delay-in-ms-between-findings: 2000
veracode configure set mock.run-scanners-in-parallel: false
veracode configure set mock.skip-upload: true

veracode configure set sast.use_realtime_flaw_apis true

```

( Need to also set the `eventing_mode` value to `LENINENT` as well somehow ... )

Run it ...
```
./veracode/veracode sast roller-orig.veracodegen.war  --debug
```

In general the data model for the config.yaml right now looks like this:

```
credentials:
    veracode_api_key_id: [ ... ]
    veracode_api_key_secret: [ ... ]
ignoreauth: "false"
mock:
    delay-in-ms-between-findings: 2000
    number-of-findings-to-throw: 7
    number-of-scanner-instances: 5
    run-scanners-in-parallel: "false"
    skip-upload: "true"
red: cat
sast:
    use_realtime_flaw_apis: "true"
ui:
    prettyprint: "true"
    verbose: "false"
urls:
    host: stage-mock-pipeline-scan-system.era1.dev.vnext.veracode.io
    pipeline_api_path: ""
    policy_api_path: /appsec/v1/policies
    scheme: https
    skip-upload: "false"
    ssl-ignore: "true"
```

and review results ...
