# N8N Terraform Provider - Supported Nodes

**Generated**: 2025-11-25T12:59:33.936Z **Provider Version**: Latest **N8N Version**: unknown **Last Sync**: 2025-11-18T01:36:46.907Z

## Overview

This document lists all **296 n8n nodes** currently supported by the Terraform provider.

### Statistics

- **Total Nodes**: 296
- **Core Nodes**: 5
- **Trigger Nodes**: 25
- **Integration Nodes**: 266
- **Require Credentials**: 252
- **No Credentials**: 44

### Testing Status

All 296 nodes have been tested with `terraform init`, `terraform validate`, `terraform apply`, and `terraform destroy`:

- ✅ **296/296 workflows passed** (100% success rate)
- Each node has a complete example workflow in `{category}/{node-slug}/` (relative to this README)
- Full test results available in root `COVERAGE.MD`

---

## Quick Navigation

- [Core Nodes](#core-nodes) (5)
- [Trigger Nodes](#trigger-nodes) (25)
- [Integration Nodes](#integration-nodes) (266)
- [Credential Requirements](#credential-requirements)
- [Usage Examples](#usage-examples)

---

## Core Nodes

Essential workflow building blocks for data manipulation, flow control, and logic.

**Total**: 5 nodes

| Node       | Type     | Description                                                            | Credentials | Example                        |
| ---------- | -------- | ---------------------------------------------------------------------- | ----------- | ------------------------------ |
| **Code**   | `code`   | Run custom JavaScript or Python code                                   | ✅ None     | [`core/code/`](core/code/)     |
| **If**     | `if`     | Route items to different branches (true/false)                         | ✅ None     | [`core/if/`](core/if/)         |
| **Merge**  | `merge`  | Merges data of multiple streams once data from both is available       | ✅ None     | [`core/merge/`](core/merge/)   |
| **Set**    | `set`    | Add or edit fields on an input item and optionally remove other fields | ✅ None     | [`core/set/`](core/set/)       |
| **Switch** | `switch` | Route items depending on defined expression or rules                   | ✅ None     | [`core/switch/`](core/switch/) |

## Trigger Nodes

Event-based nodes that initiate workflow execution.

**Total**: 25 nodes

| Node                          | Type                      | Description                                                                   | Credentials                | Example                                                                    |
| ----------------------------- | ------------------------- | ----------------------------------------------------------------------------- | -------------------------- | -------------------------------------------------------------------------- |
| **Acuity Scheduling Trigger** | `acuitySchedulingTrigger` | Handle Acuity Scheduling events via webhooks                                  | ✅ None                    | [`trigger/acuity-scheduling-trigger/`](trigger/acuity-scheduling-trigger/) |
| **Bitbucket Trigger**         | `bitbucketTrigger`        | Handle Bitbucket events via webhooks                                          | ✅ None                    | [`trigger/bitbucket-trigger/`](trigger/bitbucket-trigger/)                 |
| **Cal.com Trigger**           | `calTrigger`              | Handle Cal.com events via webhooks                                            | ✅ None                    | [`trigger/cal-com-trigger/`](trigger/cal-com-trigger/)                     |
| **Calendly Trigger**          | `calendlyTrigger`         | Starts the workflow when Calendly events occur                                | ✅ None                    | [`trigger/calendly-trigger/`](trigger/calendly-trigger/)                   |
| **Email Trigger (IMAP)**      | `emailReadImap`           | Triggers the workflow when a new email is received                            | ✅ None                    | [`trigger/email-trigger-imap/`](trigger/email-trigger-imap/)               |
| **Error Trigger**             | `errorTrigger`            | Triggers the workflow when another workflow has an error                      | ✅ None                    | [`trigger/error-trigger/`](trigger/error-trigger/)                         |
| **Eventbrite Trigger**        | `eventbriteTrigger`       | Handle Eventbrite events via webhooks                                         | ✅ None                    | [`trigger/eventbrite-trigger/`](trigger/eventbrite-trigger/)               |
| **Facebook Lead Ads Trigger** | `facebookLeadAdsTrigger`  | Handle Facebook Lead Ads events via webhooks                                  | ⚠️ Authentication Required | [`trigger/facebook-lead-ads-trigger/`](trigger/facebook-lead-ads-trigger/) |
| **Figma Trigger (Beta)**      | `figmaTrigger`            | Starts the workflow when Figma events occur                                   | ✅ None                    | [`trigger/figma-trigger-beta/`](trigger/figma-trigger-beta/)               |
| **Form.io Trigger**           | `formIoTrigger`           | Handle form.io events via webhooks                                            | ✅ None                    | [`trigger/form-io-trigger/`](trigger/form-io-trigger/)                     |
| **Formstack Trigger**         | `formstackTrigger`        | Starts the workflow on a Formstack form submission.                           | ✅ None                    | [`trigger/formstack-trigger/`](trigger/formstack-trigger/)                 |
| **Gumroad Trigger**           | `gumroadTrigger`          | Handle Gumroad events via webhooks                                            | ✅ None                    | [`trigger/gumroad-trigger/`](trigger/gumroad-trigger/)                     |
| **Jotform Trigger**           | `jotFormTrigger`          | Handle Jotform events via webhooks                                            | ✅ None                    | [`trigger/jotform-trigger/`](trigger/jotform-trigger/)                     |
| **Local File Trigger**        | `localFileTrigger`        | Triggers a workflow on file system changes                                    | ✅ None                    | [`trigger/local-file-trigger/`](trigger/local-file-trigger/)               |
| **Manual Trigger**            | `manualTrigger`           | Runs the flow on clicking a button in n8n                                     | ✅ None                    | [`trigger/manual-trigger/`](trigger/manual-trigger/)                       |
| **n8n Trigger**               | `n8nTrigger`              | Handle events and perform actions on your n8n instance                        | ✅ None                    | [`trigger/n8n-trigger/`](trigger/n8n-trigger/)                             |
| **Postmark Trigger**          | `postmarkTrigger`         | Starts the workflow when Postmark events occur                                | ✅ None                    | [`trigger/postmark-trigger/`](trigger/postmark-trigger/)                   |
| **SSE Trigger**               | `sseTrigger`              | Triggers the workflow when Server-Sent Events occur                           | ✅ None                    | [`trigger/sse-trigger/`](trigger/sse-trigger/)                             |
| **SurveyMonkey Trigger**      | `surveyMonkeyTrigger`     | Starts the workflow when Survey Monkey events occur                           | ⚠️ Authentication Required | [`trigger/surveymonkey-trigger/`](trigger/surveymonkey-trigger/)           |
| **Toggl Trigger**             | `togglTrigger`            | Starts the workflow when Toggl events occur                                   | ✅ None                    | [`trigger/toggl-trigger/`](trigger/toggl-trigger/)                         |
| **Typeform Trigger**          | `typeformTrigger`         | Starts the workflow on a Typeform form submission                             | ✅ None                    | [`trigger/typeform-trigger/`](trigger/typeform-trigger/)                   |
| **Webhook**                   | `webhook`                 | Starts the workflow when a webhook is called                                  | ✅ None                    | [`trigger/webhook/`](trigger/webhook/)                                     |
| **Workable Trigger**          | `workableTrigger`         | Starts the workflow when Workable events occur                                | ✅ None                    | [`trigger/workable-trigger/`](trigger/workable-trigger/)                   |
| **Workflow Trigger**          | `workflowTrigger`         | Triggers based on various lifecycle events, like when a workflow is activated | ✅ None                    | [`trigger/workflow-trigger/`](trigger/workflow-trigger/)                   |
| **Wufoo Trigger**             | `wufooTrigger`            | Handle Wufoo events via webhooks                                              | ✅ None                    | [`trigger/wufoo-trigger/`](trigger/wufoo-trigger/)                         |

## Integration Nodes

Third-party service integrations for connecting to external platforms.

**Total**: 266 nodes

| Node                                           | Type                            | Description                                                                                          | Credentials                | Example                                                                                                              |
| ---------------------------------------------- | ------------------------------- | ---------------------------------------------------------------------------------------------------- | -------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| **Action Network**                             | `actionNetwork`                 | Consume the Action Network API                                                                       | ⚠️ Authentication Required | [`integration/action-network/`](integration/action-network/)                                                         |
| **ActiveCampaign**                             | `activeCampaign`                | Create and edit data in ActiveCampaign                                                               | ⚠️ Authentication Required | [`integration/activecampaign/`](integration/activecampaign/)                                                         |
| **Adalo**                                      | `adalo`                         | Consume Adalo API                                                                                    | ⚠️ Authentication Required | [`integration/adalo/`](integration/adalo/)                                                                           |
| **Affinity**                                   | `affinity`                      | Consume Affinity API                                                                                 | ⚠️ Authentication Required | [`integration/affinity/`](integration/affinity/)                                                                     |
| **Agile CRM**                                  | `agileCrm`                      | Consume Agile CRM API                                                                                | ⚠️ Authentication Required | [`integration/agile-crm/`](integration/agile-crm/)                                                                   |
| **AI Transform**                               | `aiTransform`                   | Modify data based on instructions written in plain english                                           | ⚠️ Authentication Required | [`integration/ai-transform/`](integration/ai-transform/)                                                             |
| **Airtable**                                   | `airtable`                      | Read, update, write and delete data from Airtable                                                    | ⚠️ Authentication Required | [`integration/airtable/`](integration/airtable/)                                                                     |
| **Airtop**                                     | `airtop`                        | Scrape and control any site with Airtop                                                              | ⚠️ Authentication Required | [`integration/airtop/`](integration/airtop/)                                                                         |
| **AMQP Sender**                                | `amqp`                          | Sends a raw-message via AMQP 1.0, executed once per item                                             | ⚠️ Authentication Required | [`integration/amqp-sender/`](integration/amqp-sender/)                                                               |
| **APITemplate.io**                             | `apiTemplateIo`                 | Consume the APITemplate.io API                                                                       | ⚠️ API Key                 | [`integration/apitemplate-io/`](integration/apitemplate-io/)                                                         |
| **Asana**                                      | `asana`                         | Consume Asana REST API                                                                               | ⚠️ Authentication Required | [`integration/asana/`](integration/asana/)                                                                           |
| **Automizy**                                   | `automizy`                      | Consume Automizy API                                                                                 | ⚠️ Authentication Required | [`integration/automizy/`](integration/automizy/)                                                                     |
| **Autopilot**                                  | `autopilot`                     | Consume Autopilot API                                                                                | ⚠️ Authentication Required | [`integration/autopilot/`](integration/autopilot/)                                                                   |
| **AWS Lambda**                                 | `awsLambda`                     | Invoke functions on AWS Lambda                                                                       | ⚠️ Authentication Required | [`integration/aws-lambda/`](integration/aws-lambda/)                                                                 |
| **Background Color**                           | `Blur`                          | Adds a blur to the image and so makes it less sharp                                                  | ⚠️ Authentication Required | [`integration/background-color/`](integration/background-color/)                                                     |
| **BambooHr**                                   | `n8n-nodes-base.bamboohr`       | N/A                                                                                                  | ⚠️ Authentication Required | [`integration/bamboohr/`](integration/bamboohr/)                                                                     |
| **Bannerbear**                                 | `bannerbear`                    | Consume Bannerbear API                                                                               | ⚠️ Authentication Required | [`integration/bannerbear/`](integration/bannerbear/)                                                                 |
| **Baserow**                                    | `baserow`                       | Consume the Baserow API                                                                              | ⚠️ Authentication Required | [`integration/baserow/`](integration/baserow/)                                                                       |
| **Beeminder**                                  | `beeminder`                     | Consume Beeminder API                                                                                | ⚠️ Authentication Required | [`integration/beeminder/`](integration/beeminder/)                                                                   |
| **Bitly**                                      | `bitly`                         | Consume Bitly API                                                                                    | ⚠️ Authentication Required | [`integration/bitly/`](integration/bitly/)                                                                           |
| **Bitwarden**                                  | `bitwarden`                     | Consume the Bitwarden API                                                                            | ⚠️ Authentication Required | [`integration/bitwarden/`](integration/bitwarden/)                                                                   |
| **Box**                                        | `box`                           | Consume Box API                                                                                      | ⚠️ Authentication Required | [`integration/box/`](integration/box/)                                                                               |
| **Brandfetch**                                 | `Brandfetch`                    | Consume Brandfetch API                                                                               | ⚠️ Authentication Required | [`integration/brandfetch/`](integration/brandfetch/)                                                                 |
| **Brevo**                                      | `sendInBlue`                    | Consume Brevo API                                                                                    | ⚠️ Authentication Required | [`integration/brevo/`](integration/brevo/)                                                                           |
| **Bubble**                                     | `bubble`                        | Consume the Bubble Data API                                                                          | ⚠️ Authentication Required | [`integration/bubble/`](integration/bubble/)                                                                         |
| **Chargebee**                                  | `chargebee`                     | Retrieve data from Chargebee API                                                                     | ⚠️ Authentication Required | [`integration/chargebee/`](integration/chargebee/)                                                                   |
| **CircleCI**                                   | `circleCi`                      | Consume CircleCI API                                                                                 | ⚠️ Authentication Required | [`integration/circleci/`](integration/circleci/)                                                                     |
| **Clearbit**                                   | `clearbit`                      | Consume Clearbit API                                                                                 | ⚠️ Authentication Required | [`integration/clearbit/`](integration/clearbit/)                                                                     |
| **ClickUp**                                    | `clickUp`                       | Consume ClickUp API (Beta)                                                                           | ⚠️ Authentication Required | [`integration/clickup/`](integration/clickup/)                                                                       |
| **Clockify**                                   | `clockify`                      | Consume Clockify REST API                                                                            | ✅ None                    | [`integration/clockify/`](integration/clockify/)                                                                     |
| **Cloudflare**                                 | `cloudflare`                    | Consume Cloudflare API                                                                               | ⚠️ Authentication Required | [`integration/cloudflare/`](integration/cloudflare/)                                                                 |
| **Cockpit**                                    | `cockpit`                       | Consume Cockpit API                                                                                  | ⚠️ Authentication Required | [`integration/cockpit/`](integration/cockpit/)                                                                       |
| **Coda**                                       | `coda`                          | Consume Coda API                                                                                     | ⚠️ Authentication Required | [`integration/coda/`](integration/coda/)                                                                             |
| **CoinGecko**                                  | `coinGecko`                     | Consume CoinGecko API                                                                                | ⚠️ Authentication Required | [`integration/coingecko/`](integration/coingecko/)                                                                   |
| **Compare Datasets**                           | `compareDatasets`               | Compare two inputs for changes                                                                       | ✅ None                    | [`integration/compare-datasets/`](integration/compare-datasets/)                                                     |
| **Compression**                                | `compression`                   | Compress and decompress files                                                                        | ⚠️ Authentication Required | [`integration/compression/`](integration/compression/)                                                               |
| **Contentful**                                 | `contentful`                    | Consume Contentful API                                                                               | ⚠️ Authentication Required | [`integration/contentful/`](integration/contentful/)                                                                 |
| **Convert to/from binary data**                | `moveBinaryData`                | Move data between binary and JSON properties                                                         | ⚠️ Authentication Required | [`integration/convert-to-from-binary-data/`](integration/convert-to-from-binary-data/)                               |
| **ConvertKit**                                 | `convertKit`                    | Consume ConvertKit API                                                                               | ⚠️ Authentication Required | [`integration/convertkit/`](integration/convertkit/)                                                                 |
| **Copper**                                     | `copper`                        | Consume the Copper API                                                                               | ⚠️ Authentication Required | [`integration/copper/`](integration/copper/)                                                                         |
| **Cortex**                                     | `cortex`                        | Apply the Cortex analyzer/responder on the given entity                                              | ⚠️ Authentication Required | [`integration/cortex/`](integration/cortex/)                                                                         |
| **CrateDB**                                    | `crateDb`                       | Add and update data in CrateDB                                                                       | ⚠️ Authentication Required | [`integration/cratedb/`](integration/cratedb/)                                                                       |
| **Cron**                                       | `cron`                          | Triggers the workflow at a specific time                                                             | ✅ None                    | [`integration/cron/`](integration/cron/)                                                                             |
| **crowd.dev**                                  | `crowdDev`                      | crowd.dev is an open-source suite of community and data tools built to unlock community-led growth f | ⚠️ Authentication Required | [`integration/crowd-dev/`](integration/crowd-dev/)                                                                   |
| **Crypto**                                     | `crypto`                        | Provide cryptographic utilities                                                                      | ⚠️ Authentication Required | [`integration/crypto/`](integration/crypto/)                                                                         |
| **Customer Datastore (n8n training)**          | `Jay Gatsby`                    | Dummy node used for n8n training                                                                     | ⚠️ Authentication Required | [`integration/customer-datastore-n8n-training/`](integration/customer-datastore-n8n-training/)                       |
| **Customer Messenger (n8n training)**          | `n8nTrainingCustomerMessenger`  | Dummy node used for n8n training                                                                     | ⚠️ Authentication Required | [`integration/customer-messenger-n8n-training/`](integration/customer-messenger-n8n-training/)                       |
| **Customer.io**                                | `customerIo`                    | Consume Customer.io API                                                                              | ⚠️ Authentication Required | [`integration/customer-io/`](integration/customer-io/)                                                               |
| **Data table**                                 | `dataTable`                     | Permanently save data across workflow executions in a table                                          | ⚠️ Authentication Required | [`integration/data-table/`](integration/data-table/)                                                                 |
| **Date & Time**                                | `dateTime`                      | Allows you to manipulate date and time values                                                        | ⚠️ Authentication Required | [`integration/date-time/`](integration/date-time/)                                                                   |
| **DebugHelper**                                | `debugHelper`                   | Causes problems intentionally and generates useful data for debugging                                | ⚠️ Authentication Required | [`integration/debughelper/`](integration/debughelper/)                                                               |
| **DeepL**                                      | `deepL`                         | Translate data using DeepL                                                                           | ⚠️ Authentication Required | [`integration/deepl/`](integration/deepl/)                                                                           |
| **Demio**                                      | `demio`                         | Consume the Demio API                                                                                | ⚠️ Authentication Required | [`integration/demio/`](integration/demio/)                                                                           |
| **DHL**                                        | `dhl`                           | Consume DHL API                                                                                      | ⚠️ Authentication Required | [`integration/dhl/`](integration/dhl/)                                                                               |
| **Discord**                                    | `discord`                       | Sends data to Discord                                                                                | ⚠️ Authentication Required | [`integration/discord/`](integration/discord/)                                                                       |
| **Discourse**                                  | `discourse`                     | Consume Discourse API                                                                                | ⚠️ Authentication Required | [`integration/discourse/`](integration/discourse/)                                                                   |
| **Disqus**                                     | `disqus`                        | Access data on Disqus                                                                                | ⚠️ Authentication Required | [`integration/disqus/`](integration/disqus/)                                                                         |
| **Drift**                                      | `drift`                         | Consume Drift API                                                                                    | ✅ None                    | [`integration/drift/`](integration/drift/)                                                                           |
| **Dropbox**                                    | `dropbox`                       | Access data on Dropbox                                                                               | ⚠️ Authentication Required | [`integration/dropbox/`](integration/dropbox/)                                                                       |
| **Dropcontact**                                | `dropcontact`                   | Find B2B emails and enrich contacts                                                                  | ⚠️ Authentication Required | [`integration/dropcontact/`](integration/dropcontact/)                                                               |
| **E-goi**                                      | `egoi`                          | Consume E-goi API                                                                                    | ⚠️ Authentication Required | [`integration/e-goi/`](integration/e-goi/)                                                                           |
| **E2E Test**                                   | `e2eTest`                       | Dummy node used for e2e testing                                                                      | ⚠️ Authentication Required | [`integration/e2e-test/`](integration/e2e-test/)                                                                     |
| **Emelia**                                     | `emelia`                        | Consume the Emelia API                                                                               | ⚠️ Authentication Required | [`integration/emelia/`](integration/emelia/)                                                                         |
| **ERPNext**                                    | `erpNext`                       | Consume ERPNext API                                                                                  | ⚠️ Authentication Required | [`integration/erpnext/`](integration/erpnext/)                                                                       |
| **Execute Command**                            | `executeCommand`                | Executes a command on the host                                                                       | ⚠️ Authentication Required | [`integration/execute-command/`](integration/execute-command/)                                                       |
| **Execution Data**                             | `executionData`                 | Add execution data for search                                                                        | ⚠️ Authentication Required | [`integration/execution-data/`](integration/execution-data/)                                                         |
| **Extraction Values**                          | `extractionValues`              | The key under which the extracted value should be saved                                              | ⚠️ Authentication Required | [`integration/extraction-values/`](integration/extraction-values/)                                                   |
| **Facebook Graph API**                         | `facebookGraphApi`              | Interacts with Facebook using the Graph API                                                          | ⚠️ API Key                 | [`integration/facebook-graph-api/`](integration/facebook-graph-api/)                                                 |
| **FileMaker**                                  | `filemaker`                     | Retrieve data from the FileMaker data API                                                            | ⚠️ Authentication Required | [`integration/filemaker/`](integration/filemaker/)                                                                   |
| **Filter**                                     | `filter`                        | Remove items matching a condition                                                                    | ✅ None                    | [`integration/filter/`](integration/filter/)                                                                         |
| **Flow**                                       | `flow`                          | Consume Flow API                                                                                     | ⚠️ Authentication Required | [`integration/flow/`](integration/flow/)                                                                             |
| **Freshdesk**                                  | `freshdesk`                     | Consume Freshdesk API                                                                                | ⚠️ Authentication Required | [`integration/freshdesk/`](integration/freshdesk/)                                                                   |
| **Freshservice**                               | `freshservice`                  | Consume the Freshservice API                                                                         | ⚠️ Authentication Required | [`integration/freshservice/`](integration/freshservice/)                                                             |
| **Freshworks CRM**                             | `freshworksCrm`                 | Consume the Freshworks CRM API                                                                       | ⚠️ Authentication Required | [`integration/freshworks-crm/`](integration/freshworks-crm/)                                                         |
| **FTP**                                        | `ftp`                           | Transfer files via FTP or SFTP                                                                       | ⚠️ Authentication Required | [`integration/ftp/`](integration/ftp/)                                                                               |
| **Function**                                   | `function`                      | Run custom function code which gets executed once and allows you to add, remove, change and replace  | ✅ None                    | [`integration/function/`](integration/function/)                                                                     |
| **Function Item**                              | `functionItem`                  | Run custom function code which gets executed once per item                                           | ✅ None                    | [`integration/function-item/`](integration/function-item/)                                                           |
| **GetResponse**                                | `getResponse`                   | Consume GetResponse API                                                                              | ⚠️ Authentication Required | [`integration/getresponse/`](integration/getresponse/)                                                               |
| **Ghost**                                      | `ghost`                         | Consume Ghost API                                                                                    | ⚠️ Authentication Required | [`integration/ghost/`](integration/ghost/)                                                                           |
| **Git**                                        | `git`                           | Control git.                                                                                         | ⚠️ Authentication Required | [`integration/git/`](integration/git/)                                                                               |
| **GitHub**                                     | `github`                        | Consume GitHub API                                                                                   | ⚠️ Authentication Required | [`integration/github/`](integration/github/)                                                                         |
| **GitLab**                                     | `gitlab`                        | Retrieve data from GitLab API                                                                        | ⚠️ Authentication Required | [`integration/gitlab/`](integration/gitlab/)                                                                         |
| **Gong**                                       | `gong`                          | Interact with Gong API                                                                               | ⚠️ Authentication Required | [`integration/gong/`](integration/gong/)                                                                             |
| **Gotify**                                     | `gotify`                        | Consume Gotify API                                                                                   | ✅ None                    | [`integration/gotify/`](integration/gotify/)                                                                         |
| **GoToWebinar**                                | `goToWebinar`                   | Consume the GoToWebinar API                                                                          | ⚠️ Authentication Required | [`integration/gotowebinar/`](integration/gotowebinar/)                                                               |
| **Grafana**                                    | `grafana`                       | Consume the Grafana API                                                                              | ⚠️ Authentication Required | [`integration/grafana/`](integration/grafana/)                                                                       |
| **GraphQL**                                    | `graphql`                       | Makes a GraphQL request and returns the received data                                                | ⚠️ Authentication Required | [`integration/graphql/`](integration/graphql/)                                                                       |
| **Grist**                                      | `grist`                         | Consume the Grist API                                                                                | ⚠️ Authentication Required | [`integration/grist/`](integration/grist/)                                                                           |
| **Hacker News**                                | `hackerNews`                    | Consume Hacker News API                                                                              | ⚠️ Authentication Required | [`integration/hacker-news/`](integration/hacker-news/)                                                               |
| **HaloPSA**                                    | `haloPSA`                       | Consume HaloPSA API                                                                                  | ⚠️ Authentication Required | [`integration/halopsa/`](integration/halopsa/)                                                                       |
| **Harvest**                                    | `harvest`                       | Access data on Harvest                                                                               | ⚠️ Authentication Required | [`integration/harvest/`](integration/harvest/)                                                                       |
| **Help Scout**                                 | `helpScout`                     | Consume Help Scout API                                                                               | ⚠️ Authentication Required | [`integration/help-scout/`](integration/help-scout/)                                                                 |
| **HighLevel**                                  | `highLevel`                     | Consume HighLevel API                                                                                | ⚠️ Authentication Required | [`integration/highlevel/`](integration/highlevel/)                                                                   |
| **Home Assistant**                             | `homeAssistant`                 | Consume Home Assistant API                                                                           | ⚠️ Authentication Required | [`integration/home-assistant/`](integration/home-assistant/)                                                         |
| **HTML Extract**                               | `htmlExtract`                   | Extracts data from HTML                                                                              | ⚠️ Authentication Required | [`integration/html-extract/`](integration/html-extract/)                                                             |
| **HTTP Request**                               | `httpRequest`                   | Makes an HTTP request and returns the response data                                                  | ⚠️ Authentication Required | [`integration/http-request/`](integration/http-request/)                                                             |
| **HubSpot**                                    | `hubspot`                       | Consume HubSpot API                                                                                  | ⚠️ Authentication Required | [`integration/hubspot/`](integration/hubspot/)                                                                       |
| **Humantic AI**                                | `humanticAi`                    | Consume Humantic AI API                                                                              | ⚠️ Authentication Required | [`integration/humantic-ai/`](integration/humantic-ai/)                                                               |
| **Hunter**                                     | `hunter`                        | Consume Hunter API                                                                                   | ⚠️ Authentication Required | [`integration/hunter/`](integration/hunter/)                                                                         |
| **iCalendar**                                  | `iCal`                          | Create iCalendar file                                                                                | ⚠️ Authentication Required | [`integration/icalendar/`](integration/icalendar/)                                                                   |
| **Interact with Telegram using our pre-built** | `preBuiltAgentsCalloutTelegram` | Sends data to Telegram                                                                               | ⚠️ Authentication Required | [`integration/interact-with-telegram-using-our-pre-built/`](integration/interact-with-telegram-using-our-pre-built/) |
| **Intercom**                                   | `intercom`                      | Consume Intercom API                                                                                 | ⚠️ Authentication Required | [`integration/intercom/`](integration/intercom/)                                                                     |
| **Interval**                                   | `interval`                      | Triggers the workflow in a given interval                                                            | ⚠️ Authentication Required | [`integration/interval/`](integration/interval/)                                                                     |
| **Invoice Ninja**                              | `invoiceNinja`                  | Consume Invoice Ninja API                                                                            | ⚠️ Authentication Required | [`integration/invoice-ninja/`](integration/invoice-ninja/)                                                           |
| **Item Lists**                                 | `itemLists`                     | Helper for working with lists of items and transforming arrays                                       | ⚠️ Authentication Required | [`integration/item-lists/`](integration/item-lists/)                                                                 |
| **Iterable**                                   | `iterable`                      | Consume Iterable API                                                                                 | ⚠️ Authentication Required | [`integration/iterable/`](integration/iterable/)                                                                     |
| **Jenkins**                                    | `jenkins`                       | Consume Jenkins API                                                                                  | ⚠️ Authentication Required | [`integration/jenkins/`](integration/jenkins/)                                                                       |
| **Jina AI**                                    | `jinaAi`                        | Interact with Jina AI API                                                                            | ⚠️ Authentication Required | [`integration/jina-ai/`](integration/jina-ai/)                                                                       |
| **Jira Software**                              | `jira`                          | Consume Jira Software API                                                                            | ⚠️ Authentication Required | [`integration/jira-software/`](integration/jira-software/)                                                           |
| **JWT**                                        | `jwt`                           | Be sure to add a valid JWT token to the                                                              | ⚠️ Authentication Required | [`integration/jwt/`](integration/jwt/)                                                                               |
| **Kafka**                                      | `kafka`                         | Sends messages to a Kafka topic                                                                      | ⚠️ Authentication Required | [`integration/kafka/`](integration/kafka/)                                                                           |
| **Keap**                                       | `keap`                          | Consume Keap API                                                                                     | ⚠️ Authentication Required | [`integration/keap/`](integration/keap/)                                                                             |
| **Kitemaker**                                  | `kitemaker`                     | Consume the Kitemaker GraphQL API                                                                    | ⚠️ Authentication Required | [`integration/kitemaker/`](integration/kitemaker/)                                                                   |
| **KoBoToolbox**                                | `koBoToolbox`                   | Work with KoBoToolbox forms and submissions                                                          | ⚠️ Authentication Required | [`integration/kobotoolbox/`](integration/kobotoolbox/)                                                               |
| **Ldap**                                       | `ldap`                          | Interact with LDAP servers                                                                           | ⚠️ Authentication Required | [`integration/ldap/`](integration/ldap/)                                                                             |
| **Lemlist**                                    | `lemlist`                       | Consume the Lemlist API                                                                              | ⚠️ Authentication Required | [`integration/lemlist/`](integration/lemlist/)                                                                       |
| **Limit Wait Time**                            | `limitWaitTime`                 | Whether to limit the time this node should wait for a user response before execution resumes         | ✅ None                    | [`integration/limit-wait-time/`](integration/limit-wait-time/)                                                       |
| **Line**                                       | `line`                          | Consume Line API                                                                                     | ⚠️ Authentication Required | [`integration/line/`](integration/line/)                                                                             |
| **Linear**                                     | `linear`                        | Consume Linear API                                                                                   | ⚠️ Authentication Required | [`integration/linear/`](integration/linear/)                                                                         |
| **LingvaNex**                                  | `lingvaNex`                     | Consume LingvaNex API                                                                                | ⚠️ Authentication Required | [`integration/lingvanex/`](integration/lingvanex/)                                                                   |
| **LinkedIn**                                   | `linkedIn`                      | Consume LinkedIn API                                                                                 | ⚠️ Authentication Required | [`integration/linkedin/`](integration/linkedin/)                                                                     |
| **LoneScale**                                  | `loneScale`                     | Create List, add / delete items                                                                      | ⚠️ Authentication Required | [`integration/lonescale/`](integration/lonescale/)                                                                   |
| **Magento 2**                                  | `magento2`                      | Consume Magento API                                                                                  | ⚠️ Authentication Required | [`integration/magento-2/`](integration/magento-2/)                                                                   |
| **Mailcheck**                                  | `mailcheck`                     | Consume Mailcheck API                                                                                | ⚠️ Authentication Required | [`integration/mailcheck/`](integration/mailcheck/)                                                                   |
| **Mailchimp**                                  | `mailchimp`                     | Consume Mailchimp API                                                                                | ⚠️ Authentication Required | [`integration/mailchimp/`](integration/mailchimp/)                                                                   |
| **MailerLite**                                 | `mailerLite`                    | Consume MailerLite API                                                                               | ⚠️ Authentication Required | [`integration/mailerlite/`](integration/mailerlite/)                                                                 |
| **Mailgun**                                    | `mailgun`                       | Sends an email via Mailgun                                                                           | ⚠️ Authentication Required | [`integration/mailgun/`](integration/mailgun/)                                                                       |
| **Mailjet**                                    | `mailjet`                       | Consume Mailjet API                                                                                  | ⚠️ Authentication Required | [`integration/mailjet/`](integration/mailjet/)                                                                       |
| **Mandrill**                                   | `mandrill`                      | Consume Mandrill API                                                                                 | ⚠️ Authentication Required | [`integration/mandrill/`](integration/mandrill/)                                                                     |
| **Markdown**                                   | `markdown`                      | Convert data between Markdown and HTML                                                               | ⚠️ Authentication Required | [`integration/markdown/`](integration/markdown/)                                                                     |
| **Marketstack**                                | `marketstack`                   | Consume Marketstack API                                                                              | ⚠️ Authentication Required | [`integration/marketstack/`](integration/marketstack/)                                                               |
| **Matrix**                                     | `matrix`                        | Consume Matrix API                                                                                   | ⚠️ Authentication Required | [`integration/matrix/`](integration/matrix/)                                                                         |
| **Mattermost**                                 | `mattermost`                    | Sends data to Mattermost                                                                             | ⚠️ Authentication Required | [`integration/mattermost/`](integration/mattermost/)                                                                 |
| **Mautic**                                     | `mautic`                        | Consume Mautic API                                                                                   | ⚠️ Authentication Required | [`integration/mautic/`](integration/mautic/)                                                                         |
| **Medium**                                     | `medium`                        | Consume Medium API                                                                                   | ⚠️ Authentication Required | [`integration/medium/`](integration/medium/)                                                                         |
| **MessageBird**                                | `messageBird`                   | Sends SMS via MessageBird                                                                            | ⚠️ Authentication Required | [`integration/messagebird/`](integration/messagebird/)                                                               |
| **Metabase**                                   | `metabase`                      | Use the Metabase API                                                                                 | ⚠️ Authentication Required | [`integration/metabase/`](integration/metabase/)                                                                     |
| **Mindee**                                     | `mindee`                        | Consume Mindee API                                                                                   | ⚠️ Authentication Required | [`integration/mindee/`](integration/mindee/)                                                                         |
| **MISP**                                       | `misp`                          | Consume the MISP API                                                                                 | ⚠️ Authentication Required | [`integration/misp/`](integration/misp/)                                                                             |
| **Mistral AI**                                 | `mistralAi`                     | Consume Mistral AI API                                                                               | ⚠️ Authentication Required | [`integration/mistral-ai/`](integration/mistral-ai/)                                                                 |
| **Mocean**                                     | `mocean`                        | Send SMS and voice messages via Mocean                                                               | ⚠️ Authentication Required | [`integration/mocean/`](integration/mocean/)                                                                         |
| **Monday.com**                                 | `mondayCom`                     | Consume Monday.com API                                                                               | ⚠️ Authentication Required | [`integration/monday-com/`](integration/monday-com/)                                                                 |
| **MongoDB**                                    | `mongoDb`                       | Find, insert and update documents in MongoDB                                                         | ⚠️ Authentication Required | [`integration/mongodb/`](integration/mongodb/)                                                                       |
| **Monica CRM**                                 | `monicaCrm`                     | Consume the Monica CRM API                                                                           | ⚠️ Authentication Required | [`integration/monica-crm/`](integration/monica-crm/)                                                                 |
| **MQTT**                                       | `mqtt`                          | Push messages to MQTT                                                                                | ⚠️ Authentication Required | [`integration/mqtt/`](integration/mqtt/)                                                                             |
| **MSG91**                                      | `msg91`                         | Sends transactional SMS via MSG91                                                                    | ⚠️ Authentication Required | [`integration/msg91/`](integration/msg91/)                                                                           |
| **MySQL**                                      | `mySql`                         | Get, add and update data in MySQL                                                                    | ⚠️ Authentication Required | [`integration/mysql/`](integration/mysql/)                                                                           |
| **n8n**                                        | `n8n`                           | Handle events and perform actions on your n8n instance                                               | ⚠️ Authentication Required | [`integration/n8n/`](integration/n8n/)                                                                               |
| **NASA**                                       | `nasa`                          | Retrieve data from the NASA API                                                                      | ⚠️ Authentication Required | [`integration/nasa/`](integration/nasa/)                                                                             |
| **Netlify**                                    | `netlify`                       | Consume Netlify API                                                                                  | ✅ None                    | [`integration/netlify/`](integration/netlify/)                                                                       |
| **Nextcloud**                                  | `nextCloud`                     | Access data on Nextcloud                                                                             | ⚠️ Authentication Required | [`integration/nextcloud/`](integration/nextcloud/)                                                                   |
| **No Operation, do nothing**                   | `noOp`                          | No Operation                                                                                         | ⚠️ Authentication Required | [`integration/no-operation-do-nothing/`](integration/no-operation-do-nothing/)                                       |
| **NocoDB**                                     | `nocoDb`                        | Read, update, write and delete data from NocoDB                                                      | ⚠️ Authentication Required | [`integration/nocodb/`](integration/nocodb/)                                                                         |
| **Notion**                                     | `notion`                        | Consume Notion API                                                                                   | ⚠️ Authentication Required | [`integration/notion/`](integration/notion/)                                                                         |
| **Npm**                                        | `npm`                           | Consume NPM registry API                                                                             | ⚠️ Authentication Required | [`integration/npm/`](integration/npm/)                                                                               |
| **Odoo**                                       | `odoo`                          | Consume Odoo API                                                                                     | ⚠️ Authentication Required | [`integration/odoo/`](integration/odoo/)                                                                             |
| **Okta**                                       | `okta`                          | Use the Okta API                                                                                     | ⚠️ Authentication Required | [`integration/okta/`](integration/okta/)                                                                             |
| **One Simple API**                             | `oneSimpleApi`                  | A toolbox of no-code utilities                                                                       | ⚠️ API Key                 | [`integration/one-simple-api/`](integration/one-simple-api/)                                                         |
| **Onfleet**                                    | `onfleet`                       | Consume Onfleet API                                                                                  | ⚠️ Authentication Required | [`integration/onfleet/`](integration/onfleet/)                                                                       |
| **OpenAI**                                     | `openAi`                        | Consume Open AI                                                                                      | ⚠️ Authentication Required | [`integration/openai/`](integration/openai/)                                                                         |
| **OpenThesaurus**                              | `openThesaurus`                 | Get synonmns for German words using the OpenThesaurus API                                            | ⚠️ Authentication Required | [`integration/openthesaurus/`](integration/openthesaurus/)                                                           |
| **OpenWeatherMap**                             | `openWeatherMap`                | Gets current and future weather information                                                          | ⚠️ Authentication Required | [`integration/openweathermap/`](integration/openweathermap/)                                                         |
| **Orbit**                                      | `orbit`                         | Consume Orbit API                                                                                    | ⚠️ Authentication Required | [`integration/orbit/`](integration/orbit/)                                                                           |
| **Oura**                                       | `oura`                          | Consume Oura API                                                                                     | ⚠️ Authentication Required | [`integration/oura/`](integration/oura/)                                                                             |
| **Paddle**                                     | `paddle`                        | Consume Paddle API                                                                                   | ⚠️ Authentication Required | [`integration/paddle/`](integration/paddle/)                                                                         |
| **PagerDuty**                                  | `pagerDuty`                     | Consume PagerDuty API                                                                                | ⚠️ Authentication Required | [`integration/pagerduty/`](integration/pagerduty/)                                                                   |
| **PayPal**                                     | `payPal`                        | Consume PayPal API                                                                                   | ⚠️ Authentication Required | [`integration/paypal/`](integration/paypal/)                                                                         |
| **Peekalink**                                  | `peekalink`                     | Consume the Peekalink API                                                                            | ⚠️ Authentication Required | [`integration/peekalink/`](integration/peekalink/)                                                                   |
| **Perplexity**                                 | `perplexity`                    | Interact with the Perplexity API to generate AI responses with citations                             | ⚠️ Authentication Required | [`integration/perplexity/`](integration/perplexity/)                                                                 |
| **Phantombuster**                              | `phantombuster`                 | Consume Phantombuster API                                                                            | ⚠️ Authentication Required | [`integration/phantombuster/`](integration/phantombuster/)                                                           |
| **Philips Hue**                                | `philipsHue`                    | Consume Philips Hue API                                                                              | ⚠️ Authentication Required | [`integration/philips-hue/`](integration/philips-hue/)                                                               |
| **Pipedrive**                                  | `pipedrive`                     | Create and edit data in Pipedrive                                                                    | ⚠️ Authentication Required | [`integration/pipedrive/`](integration/pipedrive/)                                                                   |
| **Plivo**                                      | `plivo`                         | Send SMS/MMS messages or make phone calls                                                            | ⚠️ Authentication Required | [`integration/plivo/`](integration/plivo/)                                                                           |
| **PostBin**                                    | `postBin`                       | Consume PostBin API                                                                                  | ⚠️ Authentication Required | [`integration/postbin/`](integration/postbin/)                                                                       |
| **Postgres**                                   | `postgres`                      | Get, add and update data in Postgres                                                                 | ⚠️ Authentication Required | [`integration/postgres/`](integration/postgres/)                                                                     |
| **PostHog**                                    | `postHog`                       | Consume PostHog API                                                                                  | ⚠️ Authentication Required | [`integration/posthog/`](integration/posthog/)                                                                       |
| **ProfitWell**                                 | `profitWell`                    | Consume ProfitWell API                                                                               | ⚠️ Authentication Required | [`integration/profitwell/`](integration/profitwell/)                                                                 |
| **Pushbullet**                                 | `pushbullet`                    | Consume Pushbullet API                                                                               | ⚠️ Authentication Required | [`integration/pushbullet/`](integration/pushbullet/)                                                                 |
| **Pushcut**                                    | `pushcut`                       | Consume Pushcut API                                                                                  | ⚠️ Authentication Required | [`integration/pushcut/`](integration/pushcut/)                                                                       |
| **Pushover**                                   | `pushover`                      | Consume Pushover API                                                                                 | ⚠️ Authentication Required | [`integration/pushover/`](integration/pushover/)                                                                     |
| **QuestDB**                                    | `questDb`                       | Get, add and update data in QuestDB                                                                  | ⚠️ Authentication Required | [`integration/questdb/`](integration/questdb/)                                                                       |
| **Quick Base**                                 | `quickbase`                     | Integrate with the Quick Base RESTful API                                                            | ⚠️ Authentication Required | [`integration/quick-base/`](integration/quick-base/)                                                                 |
| **QuickBooks Online**                          | `quickbooks`                    | Consume the QuickBooks Online API                                                                    | ⚠️ Authentication Required | [`integration/quickbooks-online/`](integration/quickbooks-online/)                                                   |
| **QuickChart**                                 | `quickChart`                    | Create a chart via QuickChart                                                                        | ⚠️ Authentication Required | [`integration/quickchart/`](integration/quickchart/)                                                                 |
| **RabbitMQ**                                   | `rabbitmq`                      | Sends messages to a RabbitMQ topic                                                                   | ⚠️ Authentication Required | [`integration/rabbitmq/`](integration/rabbitmq/)                                                                     |
| **Raindrop**                                   | `raindrop`                      | Consume the Raindrop API                                                                             | ⚠️ Authentication Required | [`integration/raindrop/`](integration/raindrop/)                                                                     |
| **Read Binary File**                           | `readBinaryFile`                | Reads a binary file from disk                                                                        | ⚠️ Authentication Required | [`integration/read-binary-file/`](integration/read-binary-file/)                                                     |
| **Read Binary Files**                          | `readBinaryFiles`               | Reads binary files from disk                                                                         | ⚠️ Authentication Required | [`integration/read-binary-files/`](integration/read-binary-files/)                                                   |
| **Read PDF**                                   | `readPDF`                       | Reads a PDF and extracts its content                                                                 | ⚠️ Authentication Required | [`integration/read-pdf/`](integration/read-pdf/)                                                                     |
| **Reddit**                                     | `reddit`                        | Consume the Reddit API                                                                               | ⚠️ Authentication Required | [`integration/reddit/`](integration/reddit/)                                                                         |
| **Redis**                                      | `redis`                         | Get, send and update data in Redis                                                                   | ⚠️ Authentication Required | [`integration/redis/`](integration/redis/)                                                                           |
| **Rename Keys**                                | `renameKeys`                    | Update item field names                                                                              | ⚠️ Authentication Required | [`integration/rename-keys/`](integration/rename-keys/)                                                               |
| **Respond With**                               | `respondWith`                   | Respond with all input JSON items                                                                    | ⚠️ Authentication Required | [`integration/respond-with/`](integration/respond-with/)                                                             |
| **RocketChat**                                 | `rocketchat`                    | Consume RocketChat API                                                                               | ⚠️ Authentication Required | [`integration/rocketchat/`](integration/rocketchat/)                                                                 |
| **RSS Read**                                   | `rssFeedRead`                   | Reads data from an RSS Feed                                                                          | ⚠️ Authentication Required | [`integration/rss-read/`](integration/rss-read/)                                                                     |
| **Rundeck**                                    | `rundeck`                       | Manage Rundeck API                                                                                   | ⚠️ Authentication Required | [`integration/rundeck/`](integration/rundeck/)                                                                       |
| **S3**                                         | `s3`                            | Sends data to any S3-compatible service                                                              | ⚠️ Authentication Required | [`integration/s3/`](integration/s3/)                                                                                 |
| **Salesforce**                                 | `salesforce`                    | Consume Salesforce API                                                                               | ⚠️ Authentication Required | [`integration/salesforce/`](integration/salesforce/)                                                                 |
| **Salesmate**                                  | `salesmate`                     | Consume Salesmate API                                                                                | ⚠️ Authentication Required | [`integration/salesmate/`](integration/salesmate/)                                                                   |
| **Schedule Trigger**                           | `scheduleTrigger`               | Triggers the workflow on a given schedule                                                            | ✅ None                    | [`integration/schedule-trigger/`](integration/schedule-trigger/)                                                     |
| **SeaTable**                                   | `seaTable`                      | Read, update, write and delete data from SeaTable                                                    | ⚠️ Authentication Required | [`integration/seatable/`](integration/seatable/)                                                                     |
| **SecurityScorecard**                          | `securityScorecard`             | Consume SecurityScorecard API                                                                        | ⚠️ Authentication Required | [`integration/securityscorecard/`](integration/securityscorecard/)                                                   |
| **Segment**                                    | `segment`                       | Consume Segment API                                                                                  | ⚠️ Authentication Required | [`integration/segment/`](integration/segment/)                                                                       |
| **Send Email**                                 | `emailSend`                     | Sends an email using SMTP protocol                                                                   | ⚠️ Authentication Required | [`integration/send-email/`](integration/send-email/)                                                                 |
| **SendGrid**                                   | `sendGrid`                      | Consume SendGrid API                                                                                 | ⚠️ Authentication Required | [`integration/sendgrid/`](integration/sendgrid/)                                                                     |
| **Sendy**                                      | `sendy`                         | Consume Sendy API                                                                                    | ⚠️ Authentication Required | [`integration/sendy/`](integration/sendy/)                                                                           |
| **Sentry.io**                                  | `sentryIo`                      | Consume Sentry.io API                                                                                | ⚠️ Authentication Required | [`integration/sentry-io/`](integration/sentry-io/)                                                                   |
| **ServiceNow**                                 | `serviceNow`                    | Consume ServiceNow API                                                                               | ⚠️ Authentication Required | [`integration/servicenow/`](integration/servicenow/)                                                                 |
| **seven**                                      | `sms77`                         | Send SMS and make text-to-speech calls                                                               | ⚠️ Authentication Required | [`integration/seven/`](integration/seven/)                                                                           |
| **Shopify**                                    | `shopify`                       | Consume Shopify API                                                                                  | ✅ None                    | [`integration/shopify/`](integration/shopify/)                                                                       |
| **SIGNL4**                                     | `signl4`                        | Consume SIGNL4 API                                                                                   | ⚠️ Authentication Required | [`integration/signl4/`](integration/signl4/)                                                                         |
| **Simulate**                                   | `simulate`                      | Simulate a node                                                                                      | ⚠️ Authentication Required | [`integration/simulate/`](integration/simulate/)                                                                     |
| **Slack**                                      | `slack`                         | Consume Slack API                                                                                    | ⚠️ Authentication Required | [`integration/slack/`](integration/slack/)                                                                           |
| **Snowflake**                                  | `snowflake`                     | Get, add and update data in Snowflake                                                                | ⚠️ Authentication Required | [`integration/snowflake/`](integration/snowflake/)                                                                   |
| **Split In Batches**                           | `splitInBatches`                | Split data into batches and iterate over each batch                                                  | ✅ None                    | [`integration/split-in-batches/`](integration/split-in-batches/)                                                     |
| **Splunk**                                     | `splunk`                        | Consume the Splunk Enterprise API                                                                    | ⚠️ Authentication Required | [`integration/splunk/`](integration/splunk/)                                                                         |
| **Spontit**                                    | `spontit`                       | Consume Spontit API                                                                                  | ⚠️ Authentication Required | [`integration/spontit/`](integration/spontit/)                                                                       |
| **Spotify**                                    | `spotify`                       | Access public song data via the Spotify API                                                          | ✅ None                    | [`integration/spotify/`](integration/spotify/)                                                                       |
| **Spreadsheet File**                           | `spreadsheetFile`               | Reads and writes data from a spreadsheet file like CSV, XLS, ODS, etc                                | ⚠️ Authentication Required | [`integration/spreadsheet-file/`](integration/spreadsheet-file/)                                                     |
| **SSH**                                        | `ssh`                           | Execute commands via SSH                                                                             | ⚠️ Authentication Required | [`integration/ssh/`](integration/ssh/)                                                                               |
| **Stackby**                                    | `stackby`                       | Read, write, and delete data in Stackby                                                              | ⚠️ Authentication Required | [`integration/stackby/`](integration/stackby/)                                                                       |
| **Start**                                      | `start`                         | Starts the workflow execution from this node                                                         | ⚠️ Authentication Required | [`integration/start/`](integration/start/)                                                                           |
| **Sticky Note**                                | `stickyNote`                    | Make your workflow easier to understand                                                              | ⚠️ Authentication Required | [`integration/sticky-note/`](integration/sticky-note/)                                                               |
| **Stop and Error**                             | `stopAndError`                  | Throw an error in the workflow                                                                       | ✅ None                    | [`integration/stop-and-error/`](integration/stop-and-error/)                                                         |
| **Storyblok**                                  | `storyblok`                     | Consume Storyblok API                                                                                | ⚠️ Authentication Required | [`integration/storyblok/`](integration/storyblok/)                                                                   |
| **Strapi**                                     | `strapi`                        | Consume Strapi API                                                                                   | ⚠️ API Key                 | [`integration/strapi/`](integration/strapi/)                                                                         |
| **Strava**                                     | `strava`                        | Consume Strava API                                                                                   | ⚠️ Authentication Required | [`integration/strava/`](integration/strava/)                                                                         |
| **Stripe**                                     | `stripe`                        | Consume the Stripe API                                                                               | ⚠️ Authentication Required | [`integration/stripe/`](integration/stripe/)                                                                         |
| **Supabase**                                   | `supabase`                      | Add, get, delete and update data in a table                                                          | ⚠️ Authentication Required | [`integration/supabase/`](integration/supabase/)                                                                     |
| **SyncroMSP**                                  | `syncroMsp`                     | Manage contacts, tickets and more from Syncro MSP                                                    | ⚠️ Authentication Required | [`integration/syncromsp/`](integration/syncromsp/)                                                                   |
| **Taiga**                                      | `taiga`                         | Consume Taiga API                                                                                    | ⚠️ Authentication Required | [`integration/taiga/`](integration/taiga/)                                                                           |
| **Tapfiliate**                                 | `tapfiliate`                    | Consume Tapfiliate API                                                                               | ⚠️ Authentication Required | [`integration/tapfiliate/`](integration/tapfiliate/)                                                                 |
| **TheHive**                                    | `theHive`                       | Consume TheHive API                                                                                  | ⚠️ Authentication Required | [`integration/thehive/`](integration/thehive/)                                                                       |
| **TheHiveProject**                             | `n8n-nodes-base.thehiveproject` | N/A                                                                                                  | ⚠️ Authentication Required | [`integration/thehiveproject/`](integration/thehiveproject/)                                                         |
| **TimescaleDB**                                | `timescaleDb`                   | Add and update data in TimescaleDB                                                                   | ⚠️ Authentication Required | [`integration/timescaledb/`](integration/timescaledb/)                                                               |
| **Todoist**                                    | `todoist`                       | Consume Todoist API                                                                                  | ⚠️ Authentication Required | [`integration/todoist/`](integration/todoist/)                                                                       |
| **TOTP**                                       | `totp`                          | Generate a time-based one-time password                                                              | ⚠️ Authentication Required | [`integration/totp/`](integration/totp/)                                                                             |
| **TravisCI**                                   | `travisCi`                      | Consume TravisCI API                                                                                 | ⚠️ Authentication Required | [`integration/travisci/`](integration/travisci/)                                                                     |
| **Trello**                                     | `trello`                        | Create, change and delete boards and cards                                                           | ⚠️ Authentication Required | [`integration/trello/`](integration/trello/)                                                                         |
| **Twake**                                      | `twake`                         | Consume Twake API                                                                                    | ⚠️ Authentication Required | [`integration/twake/`](integration/twake/)                                                                           |
| **Twilio**                                     | `twilio`                        | Send SMS and WhatsApp messages or make phone calls                                                   | ⚠️ Authentication Required | [`integration/twilio/`](integration/twilio/)                                                                         |
| **Twist**                                      | `twist`                         | Consume Twist API                                                                                    | ⚠️ Authentication Required | [`integration/twist/`](integration/twist/)                                                                           |
| **Unleashed Software**                         | `unleashedSoftware`             | Consume Unleashed Software API                                                                       | ⚠️ Authentication Required | [`integration/unleashed-software/`](integration/unleashed-software/)                                                 |
| **Uplead**                                     | `uplead`                        | Consume Uplead API                                                                                   | ⚠️ Authentication Required | [`integration/uplead/`](integration/uplead/)                                                                         |
| **uProc**                                      | `uproc`                         | Consume uProc API                                                                                    | ⚠️ Authentication Required | [`integration/uproc/`](integration/uproc/)                                                                           |
| **UptimeRobot**                                | `uptimeRobot`                   | Consume UptimeRobot API                                                                              | ⚠️ Authentication Required | [`integration/uptimerobot/`](integration/uptimerobot/)                                                               |
| **urlscan.io**                                 | `urlScanIo`                     | Provides various utilities for monitoring websites like health checks or screenshots                 | ⚠️ Authentication Required | [`integration/urlscan-io/`](integration/urlscan-io/)                                                                 |
| **Vero**                                       | `vero`                          | Consume Vero API                                                                                     | ⚠️ Authentication Required | [`integration/vero/`](integration/vero/)                                                                             |
| **Vonage**                                     | `vonage`                        | Consume Vonage API                                                                                   | ⚠️ Authentication Required | [`integration/vonage/`](integration/vonage/)                                                                         |
| **Wait Amount**                                | `amount`                        | The time to wait                                                                                     | ✅ None                    | [`integration/wait-amount/`](integration/wait-amount/)                                                               |
| **Webflow**                                    | `webflow`                       | Consume the Webflow API                                                                              | ⚠️ Authentication Required | [`integration/webflow/`](integration/webflow/)                                                                       |
| **Wekan**                                      | `wekan`                         | Consume Wekan API                                                                                    | ⚠️ Authentication Required | [`integration/wekan/`](integration/wekan/)                                                                           |
| **WhatsApp Business Cloud**                    | `whatsApp`                      | Access WhatsApp API                                                                                  | ⚠️ Authentication Required | [`integration/whatsapp-business-cloud/`](integration/whatsapp-business-cloud/)                                       |
| **Wise**                                       | `wise`                          | Consume the Wise API                                                                                 | ⚠️ Authentication Required | [`integration/wise/`](integration/wise/)                                                                             |
| **WooCommerce**                                | `wooCommerce`                   | Consume WooCommerce API                                                                              | ⚠️ Authentication Required | [`integration/woocommerce/`](integration/woocommerce/)                                                               |
| **Wordpress**                                  | `wordpress`                     | Consume Wordpress API                                                                                | ⚠️ Authentication Required | [`integration/wordpress/`](integration/wordpress/)                                                                   |
| **Write Binary File**                          | `writeBinaryFile`               | Writes a binary file to disk                                                                         | ⚠️ Authentication Required | [`integration/write-binary-file/`](integration/write-binary-file/)                                                   |
| **X (Formerly Twitter)**                       | `twitter`                       | Consume the X API                                                                                    | ⚠️ Authentication Required | [`integration/x-formerly-twitter/`](integration/x-formerly-twitter/)                                                 |
| **Xero**                                       | `xero`                          | Consume Xero API                                                                                     | ⚠️ Authentication Required | [`integration/xero/`](integration/xero/)                                                                             |
| **XML**                                        | `xml`                           | Convert data from and to XML                                                                         | ⚠️ Authentication Required | [`integration/xml/`](integration/xml/)                                                                               |
| **Yourls**                                     | `yourls`                        | Consume Yourls API                                                                                   | ⚠️ Authentication Required | [`integration/yourls/`](integration/yourls/)                                                                         |
| **Zammad**                                     | `zammad`                        | Consume the Zammad API                                                                               | ⚠️ Authentication Required | [`integration/zammad/`](integration/zammad/)                                                                         |
| **Zendesk**                                    | `zendesk`                       | Consume Zendesk API                                                                                  | ⚠️ Authentication Required | [`integration/zendesk/`](integration/zendesk/)                                                                       |
| **Zoho CRM**                                   | `zohoCrm`                       | Consume Zoho CRM API                                                                                 | ⚠️ Authentication Required | [`integration/zoho-crm/`](integration/zoho-crm/)                                                                     |
| **Zoom**                                       | `zoom`                          | Consume Zoom API                                                                                     | ⚠️ Authentication Required | [`integration/zoom/`](integration/zoom/)                                                                             |
| **Zulip**                                      | `zulip`                         | Consume Zulip API                                                                                    | ⚠️ Authentication Required | [`integration/zulip/`](integration/zulip/)                                                                           |

---

## Credential Requirements

Some nodes require external service credentials (API keys, OAuth tokens, etc.).

### Nodes Requiring Credentials (252)

#### Authentication Required (248)

- **Action Network** - [`actionNetwork`](integration/action-network/)
- **ActiveCampaign** - [`activeCampaign`](integration/activecampaign/)
- **Adalo** - [`adalo`](integration/adalo/)
- **Affinity** - [`affinity`](integration/affinity/)
- **Agile CRM** - [`agileCrm`](integration/agile-crm/)
- **AI Transform** - [`aiTransform`](integration/ai-transform/)
- **Airtable** - [`airtable`](integration/airtable/)
- **Airtop** - [`airtop`](integration/airtop/)
- **AMQP Sender** - [`amqp`](integration/amqp-sender/)
- **Asana** - [`asana`](integration/asana/)
- **Automizy** - [`automizy`](integration/automizy/)
- **Autopilot** - [`autopilot`](integration/autopilot/)
- **AWS Lambda** - [`awsLambda`](integration/aws-lambda/)
- **Background Color** - [`Blur`](integration/background-color/)
- **BambooHr** - [`n8n-nodes-base.bamboohr`](integration/bamboohr/)
- **Bannerbear** - [`bannerbear`](integration/bannerbear/)
- **Baserow** - [`baserow`](integration/baserow/)
- **Beeminder** - [`beeminder`](integration/beeminder/)
- **Bitly** - [`bitly`](integration/bitly/)
- **Bitwarden** - [`bitwarden`](integration/bitwarden/)
- **Box** - [`box`](integration/box/)
- **Brandfetch** - [`Brandfetch`](integration/brandfetch/)
- **Brevo** - [`sendInBlue`](integration/brevo/)
- **Bubble** - [`bubble`](integration/bubble/)
- **Chargebee** - [`chargebee`](integration/chargebee/)
- **CircleCI** - [`circleCi`](integration/circleci/)
- **Clearbit** - [`clearbit`](integration/clearbit/)
- **ClickUp** - [`clickUp`](integration/clickup/)
- **Cloudflare** - [`cloudflare`](integration/cloudflare/)
- **Cockpit** - [`cockpit`](integration/cockpit/)
- **Coda** - [`coda`](integration/coda/)
- **CoinGecko** - [`coinGecko`](integration/coingecko/)
- **Compression** - [`compression`](integration/compression/)
- **Contentful** - [`contentful`](integration/contentful/)
- **Convert to/from binary data** - [`moveBinaryData`](integration/convert-to-from-binary-data/)
- **ConvertKit** - [`convertKit`](integration/convertkit/)
- **Copper** - [`copper`](integration/copper/)
- **Cortex** - [`cortex`](integration/cortex/)
- **CrateDB** - [`crateDb`](integration/cratedb/)
- **crowd.dev** - [`crowdDev`](integration/crowd-dev/)
- **Crypto** - [`crypto`](integration/crypto/)
- **Customer Datastore (n8n training)** - [`Jay Gatsby`](integration/customer-datastore-n8n-training/)
- **Customer Messenger (n8n training)** - [`n8nTrainingCustomerMessenger`](integration/customer-messenger-n8n-training/)
- **Customer.io** - [`customerIo`](integration/customer-io/)
- **Data table** - [`dataTable`](integration/data-table/)
- **Date & Time** - [`dateTime`](integration/date-time/)
- **DebugHelper** - [`debugHelper`](integration/debughelper/)
- **DeepL** - [`deepL`](integration/deepl/)
- **Demio** - [`demio`](integration/demio/)
- **DHL** - [`dhl`](integration/dhl/)
- **Discord** - [`discord`](integration/discord/)
- **Discourse** - [`discourse`](integration/discourse/)
- **Disqus** - [`disqus`](integration/disqus/)
- **Dropbox** - [`dropbox`](integration/dropbox/)
- **Dropcontact** - [`dropcontact`](integration/dropcontact/)
- **E-goi** - [`egoi`](integration/e-goi/)
- **E2E Test** - [`e2eTest`](integration/e2e-test/)
- **Emelia** - [`emelia`](integration/emelia/)
- **ERPNext** - [`erpNext`](integration/erpnext/)
- **Execute Command** - [`executeCommand`](integration/execute-command/)
- **Execution Data** - [`executionData`](integration/execution-data/)
- **Extraction Values** - [`extractionValues`](integration/extraction-values/)
- **Facebook Lead Ads Trigger** - [`facebookLeadAdsTrigger`](trigger/facebook-lead-ads-trigger/)
- **FileMaker** - [`filemaker`](integration/filemaker/)
- **Flow** - [`flow`](integration/flow/)
- **Freshdesk** - [`freshdesk`](integration/freshdesk/)
- **Freshservice** - [`freshservice`](integration/freshservice/)
- **Freshworks CRM** - [`freshworksCrm`](integration/freshworks-crm/)
- **FTP** - [`ftp`](integration/ftp/)
- **GetResponse** - [`getResponse`](integration/getresponse/)
- **Ghost** - [`ghost`](integration/ghost/)
- **Git** - [`git`](integration/git/)
- **GitHub** - [`github`](integration/github/)
- **GitLab** - [`gitlab`](integration/gitlab/)
- **Gong** - [`gong`](integration/gong/)
- **GoToWebinar** - [`goToWebinar`](integration/gotowebinar/)
- **Grafana** - [`grafana`](integration/grafana/)
- **GraphQL** - [`graphql`](integration/graphql/)
- **Grist** - [`grist`](integration/grist/)
- **Hacker News** - [`hackerNews`](integration/hacker-news/)
- **HaloPSA** - [`haloPSA`](integration/halopsa/)
- **Harvest** - [`harvest`](integration/harvest/)
- **Help Scout** - [`helpScout`](integration/help-scout/)
- **HighLevel** - [`highLevel`](integration/highlevel/)
- **Home Assistant** - [`homeAssistant`](integration/home-assistant/)
- **HTML Extract** - [`htmlExtract`](integration/html-extract/)
- **HTTP Request** - [`httpRequest`](integration/http-request/)
- **HubSpot** - [`hubspot`](integration/hubspot/)
- **Humantic AI** - [`humanticAi`](integration/humantic-ai/)
- **Hunter** - [`hunter`](integration/hunter/)
- **iCalendar** - [`iCal`](integration/icalendar/)
- **Interact with Telegram using our pre-built** - [`preBuiltAgentsCalloutTelegram`](integration/interact-with-telegram-using-our-pre-built/)
- **Intercom** - [`intercom`](integration/intercom/)
- **Interval** - [`interval`](integration/interval/)
- **Invoice Ninja** - [`invoiceNinja`](integration/invoice-ninja/)
- **Item Lists** - [`itemLists`](integration/item-lists/)
- **Iterable** - [`iterable`](integration/iterable/)
- **Jenkins** - [`jenkins`](integration/jenkins/)
- **Jina AI** - [`jinaAi`](integration/jina-ai/)
- **Jira Software** - [`jira`](integration/jira-software/)
- **JWT** - [`jwt`](integration/jwt/)
- **Kafka** - [`kafka`](integration/kafka/)
- **Keap** - [`keap`](integration/keap/)
- **Kitemaker** - [`kitemaker`](integration/kitemaker/)
- **KoBoToolbox** - [`koBoToolbox`](integration/kobotoolbox/)
- **Ldap** - [`ldap`](integration/ldap/)
- **Lemlist** - [`lemlist`](integration/lemlist/)
- **Line** - [`line`](integration/line/)
- **Linear** - [`linear`](integration/linear/)
- **LingvaNex** - [`lingvaNex`](integration/lingvanex/)
- **LinkedIn** - [`linkedIn`](integration/linkedin/)
- **LoneScale** - [`loneScale`](integration/lonescale/)
- **Magento 2** - [`magento2`](integration/magento-2/)
- **Mailcheck** - [`mailcheck`](integration/mailcheck/)
- **Mailchimp** - [`mailchimp`](integration/mailchimp/)
- **MailerLite** - [`mailerLite`](integration/mailerlite/)
- **Mailgun** - [`mailgun`](integration/mailgun/)
- **Mailjet** - [`mailjet`](integration/mailjet/)
- **Mandrill** - [`mandrill`](integration/mandrill/)
- **Markdown** - [`markdown`](integration/markdown/)
- **Marketstack** - [`marketstack`](integration/marketstack/)
- **Matrix** - [`matrix`](integration/matrix/)
- **Mattermost** - [`mattermost`](integration/mattermost/)
- **Mautic** - [`mautic`](integration/mautic/)
- **Medium** - [`medium`](integration/medium/)
- **MessageBird** - [`messageBird`](integration/messagebird/)
- **Metabase** - [`metabase`](integration/metabase/)
- **Mindee** - [`mindee`](integration/mindee/)
- **MISP** - [`misp`](integration/misp/)
- **Mistral AI** - [`mistralAi`](integration/mistral-ai/)
- **Mocean** - [`mocean`](integration/mocean/)
- **Monday.com** - [`mondayCom`](integration/monday-com/)
- **MongoDB** - [`mongoDb`](integration/mongodb/)
- **Monica CRM** - [`monicaCrm`](integration/monica-crm/)
- **MQTT** - [`mqtt`](integration/mqtt/)
- **MSG91** - [`msg91`](integration/msg91/)
- **MySQL** - [`mySql`](integration/mysql/)
- **n8n** - [`n8n`](integration/n8n/)
- **NASA** - [`nasa`](integration/nasa/)
- **Nextcloud** - [`nextCloud`](integration/nextcloud/)
- **No Operation, do nothing** - [`noOp`](integration/no-operation-do-nothing/)
- **NocoDB** - [`nocoDb`](integration/nocodb/)
- **Notion** - [`notion`](integration/notion/)
- **Npm** - [`npm`](integration/npm/)
- **Odoo** - [`odoo`](integration/odoo/)
- **Okta** - [`okta`](integration/okta/)
- **Onfleet** - [`onfleet`](integration/onfleet/)
- **OpenAI** - [`openAi`](integration/openai/)
- **OpenThesaurus** - [`openThesaurus`](integration/openthesaurus/)
- **OpenWeatherMap** - [`openWeatherMap`](integration/openweathermap/)
- **Orbit** - [`orbit`](integration/orbit/)
- **Oura** - [`oura`](integration/oura/)
- **Paddle** - [`paddle`](integration/paddle/)
- **PagerDuty** - [`pagerDuty`](integration/pagerduty/)
- **PayPal** - [`payPal`](integration/paypal/)
- **Peekalink** - [`peekalink`](integration/peekalink/)
- **Perplexity** - [`perplexity`](integration/perplexity/)
- **Phantombuster** - [`phantombuster`](integration/phantombuster/)
- **Philips Hue** - [`philipsHue`](integration/philips-hue/)
- **Pipedrive** - [`pipedrive`](integration/pipedrive/)
- **Plivo** - [`plivo`](integration/plivo/)
- **PostBin** - [`postBin`](integration/postbin/)
- **Postgres** - [`postgres`](integration/postgres/)
- **PostHog** - [`postHog`](integration/posthog/)
- **ProfitWell** - [`profitWell`](integration/profitwell/)
- **Pushbullet** - [`pushbullet`](integration/pushbullet/)
- **Pushcut** - [`pushcut`](integration/pushcut/)
- **Pushover** - [`pushover`](integration/pushover/)
- **QuestDB** - [`questDb`](integration/questdb/)
- **Quick Base** - [`quickbase`](integration/quick-base/)
- **QuickBooks Online** - [`quickbooks`](integration/quickbooks-online/)
- **QuickChart** - [`quickChart`](integration/quickchart/)
- **RabbitMQ** - [`rabbitmq`](integration/rabbitmq/)
- **Raindrop** - [`raindrop`](integration/raindrop/)
- **Read Binary File** - [`readBinaryFile`](integration/read-binary-file/)
- **Read Binary Files** - [`readBinaryFiles`](integration/read-binary-files/)
- **Read PDF** - [`readPDF`](integration/read-pdf/)
- **Reddit** - [`reddit`](integration/reddit/)
- **Redis** - [`redis`](integration/redis/)
- **Rename Keys** - [`renameKeys`](integration/rename-keys/)
- **Respond With** - [`respondWith`](integration/respond-with/)
- **RocketChat** - [`rocketchat`](integration/rocketchat/)
- **RSS Read** - [`rssFeedRead`](integration/rss-read/)
- **Rundeck** - [`rundeck`](integration/rundeck/)
- **S3** - [`s3`](integration/s3/)
- **Salesforce** - [`salesforce`](integration/salesforce/)
- **Salesmate** - [`salesmate`](integration/salesmate/)
- **SeaTable** - [`seaTable`](integration/seatable/)
- **SecurityScorecard** - [`securityScorecard`](integration/securityscorecard/)
- **Segment** - [`segment`](integration/segment/)
- **Send Email** - [`emailSend`](integration/send-email/)
- **SendGrid** - [`sendGrid`](integration/sendgrid/)
- **Sendy** - [`sendy`](integration/sendy/)
- **Sentry.io** - [`sentryIo`](integration/sentry-io/)
- **ServiceNow** - [`serviceNow`](integration/servicenow/)
- **seven** - [`sms77`](integration/seven/)
- **SIGNL4** - [`signl4`](integration/signl4/)
- **Simulate** - [`simulate`](integration/simulate/)
- **Slack** - [`slack`](integration/slack/)
- **Snowflake** - [`snowflake`](integration/snowflake/)
- **Splunk** - [`splunk`](integration/splunk/)
- **Spontit** - [`spontit`](integration/spontit/)
- **Spreadsheet File** - [`spreadsheetFile`](integration/spreadsheet-file/)
- **SSH** - [`ssh`](integration/ssh/)
- **Stackby** - [`stackby`](integration/stackby/)
- **Start** - [`start`](integration/start/)
- **Sticky Note** - [`stickyNote`](integration/sticky-note/)
- **Storyblok** - [`storyblok`](integration/storyblok/)
- **Strava** - [`strava`](integration/strava/)
- **Stripe** - [`stripe`](integration/stripe/)
- **Supabase** - [`supabase`](integration/supabase/)
- **SurveyMonkey Trigger** - [`surveyMonkeyTrigger`](trigger/surveymonkey-trigger/)
- **SyncroMSP** - [`syncroMsp`](integration/syncromsp/)
- **Taiga** - [`taiga`](integration/taiga/)
- **Tapfiliate** - [`tapfiliate`](integration/tapfiliate/)
- **TheHive** - [`theHive`](integration/thehive/)
- **TheHiveProject** - [`n8n-nodes-base.thehiveproject`](integration/thehiveproject/)
- **TimescaleDB** - [`timescaleDb`](integration/timescaledb/)
- **Todoist** - [`todoist`](integration/todoist/)
- **TOTP** - [`totp`](integration/totp/)
- **TravisCI** - [`travisCi`](integration/travisci/)
- **Trello** - [`trello`](integration/trello/)
- **Twake** - [`twake`](integration/twake/)
- **Twilio** - [`twilio`](integration/twilio/)
- **Twist** - [`twist`](integration/twist/)
- **Unleashed Software** - [`unleashedSoftware`](integration/unleashed-software/)
- **Uplead** - [`uplead`](integration/uplead/)
- **uProc** - [`uproc`](integration/uproc/)
- **UptimeRobot** - [`uptimeRobot`](integration/uptimerobot/)
- **urlscan.io** - [`urlScanIo`](integration/urlscan-io/)
- **Vero** - [`vero`](integration/vero/)
- **Vonage** - [`vonage`](integration/vonage/)
- **Webflow** - [`webflow`](integration/webflow/)
- **Wekan** - [`wekan`](integration/wekan/)
- **WhatsApp Business Cloud** - [`whatsApp`](integration/whatsapp-business-cloud/)
- **Wise** - [`wise`](integration/wise/)
- **WooCommerce** - [`wooCommerce`](integration/woocommerce/)
- **Wordpress** - [`wordpress`](integration/wordpress/)
- **Write Binary File** - [`writeBinaryFile`](integration/write-binary-file/)
- **X (Formerly Twitter)** - [`twitter`](integration/x-formerly-twitter/)
- **Xero** - [`xero`](integration/xero/)
- **XML** - [`xml`](integration/xml/)
- **Yourls** - [`yourls`](integration/yourls/)
- **Zammad** - [`zammad`](integration/zammad/)
- **Zendesk** - [`zendesk`](integration/zendesk/)
- **Zoho CRM** - [`zohoCrm`](integration/zoho-crm/)
- **Zoom** - [`zoom`](integration/zoom/)
- **Zulip** - [`zulip`](integration/zulip/)

#### API Key (4)

- **APITemplate.io** - [`apiTemplateIo`](integration/apitemplate-io/)
- **Facebook Graph API** - [`facebookGraphApi`](integration/facebook-graph-api/)
- **One Simple API** - [`oneSimpleApi`](integration/one-simple-api/)
- **Strapi** - [`strapi`](integration/strapi/)

### Nodes Without Credentials (44)

These nodes work out-of-the-box without external credentials:

- **Acuity Scheduling Trigger** (`acuitySchedulingTrigger`) - [Example](trigger/acuity-scheduling-trigger/)
- **Bitbucket Trigger** (`bitbucketTrigger`) - [Example](trigger/bitbucket-trigger/)
- **Cal.com Trigger** (`calTrigger`) - [Example](trigger/cal-com-trigger/)
- **Calendly Trigger** (`calendlyTrigger`) - [Example](trigger/calendly-trigger/)
- **Clockify** (`clockify`) - [Example](integration/clockify/)
- **Code** (`code`) - [Example](core/code/)
- **Compare Datasets** (`compareDatasets`) - [Example](integration/compare-datasets/)
- **Cron** (`cron`) - [Example](integration/cron/)
- **Drift** (`drift`) - [Example](integration/drift/)
- **Email Trigger (IMAP)** (`emailReadImap`) - [Example](trigger/email-trigger-imap/)
- **Error Trigger** (`errorTrigger`) - [Example](trigger/error-trigger/)
- **Eventbrite Trigger** (`eventbriteTrigger`) - [Example](trigger/eventbrite-trigger/)
- **Figma Trigger (Beta)** (`figmaTrigger`) - [Example](trigger/figma-trigger-beta/)
- **Filter** (`filter`) - [Example](integration/filter/)
- **Form.io Trigger** (`formIoTrigger`) - [Example](trigger/form-io-trigger/)
- **Formstack Trigger** (`formstackTrigger`) - [Example](trigger/formstack-trigger/)
- **Function** (`function`) - [Example](integration/function/)
- **Function Item** (`functionItem`) - [Example](integration/function-item/)
- **Gotify** (`gotify`) - [Example](integration/gotify/)
- **Gumroad Trigger** (`gumroadTrigger`) - [Example](trigger/gumroad-trigger/)
- **If** (`if`) - [Example](core/if/)
- **Jotform Trigger** (`jotFormTrigger`) - [Example](trigger/jotform-trigger/)
- **Limit Wait Time** (`limitWaitTime`) - [Example](integration/limit-wait-time/)
- **Local File Trigger** (`localFileTrigger`) - [Example](trigger/local-file-trigger/)
- **Manual Trigger** (`manualTrigger`) - [Example](trigger/manual-trigger/)
- **Merge** (`merge`) - [Example](core/merge/)
- **n8n Trigger** (`n8nTrigger`) - [Example](trigger/n8n-trigger/)
- **Netlify** (`netlify`) - [Example](integration/netlify/)
- **Postmark Trigger** (`postmarkTrigger`) - [Example](trigger/postmark-trigger/)
- **Schedule Trigger** (`scheduleTrigger`) - [Example](integration/schedule-trigger/)
- **Set** (`set`) - [Example](core/set/)
- **Shopify** (`shopify`) - [Example](integration/shopify/)
- **Split In Batches** (`splitInBatches`) - [Example](integration/split-in-batches/)
- **Spotify** (`spotify`) - [Example](integration/spotify/)
- **SSE Trigger** (`sseTrigger`) - [Example](trigger/sse-trigger/)
- **Stop and Error** (`stopAndError`) - [Example](integration/stop-and-error/)
- **Switch** (`switch`) - [Example](core/switch/)
- **Toggl Trigger** (`togglTrigger`) - [Example](trigger/toggl-trigger/)
- **Typeform Trigger** (`typeformTrigger`) - [Example](trigger/typeform-trigger/)
- **Wait Amount** (`amount`) - [Example](integration/wait-amount/)
- **Webhook** (`webhook`) - [Example](trigger/webhook/)
- **Workable Trigger** (`workableTrigger`) - [Example](trigger/workable-trigger/)
- **Workflow Trigger** (`workflowTrigger`) - [Example](trigger/workflow-trigger/)
- **Wufoo Trigger** (`wufooTrigger`) - [Example](trigger/wufoo-trigger/)

---

## Usage Examples

### Basic Node Usage

Each node can be used in a Terraform workflow:

```hcl
# Example: Using the Code node
resource "n8n_workflow_node" "my_code" {
  name     = "Process Data"
  type     = "n8n-nodes-base.code"
  position = [250, 300]

  parameters = jsonencode({
    mode = "runOnceForAllItems"
    jsCode = "return items;"
  })
}
```

### Complete Workflow Examples

Every node has a complete, tested workflow example in:

```
{category}/{node-slug}/
  ├── main.tf         # Complete workflow with the node
  ├── variables.tf    # Provider configuration
  └── README.md       # Node-specific documentation
```

(All paths are relative to this README location: `examples/nodes/`)

### Testing Your Workflow

```bash
cd core/code
terraform init
terraform validate
terraform plan
```

---

## Contributing

### Reporting Issues

If you find issues with any node:

1. Check the node's example in `examples/nodes/`
2. Review the [n8n node documentation](https://docs.n8n.io/integrations/)
3. Report issues on [GitHub](https://github.com/yourusername/terraform-provider-n8n/issues)

### Adding New Nodes

New nodes are automatically synchronized from the official n8n repository:

1. Run `make nodes` to sync latest nodes
2. Review `NODES_SYNC.md` for changes
3. Test new nodes: `make nodes/test-workflows`
4. Update documentation: `make nodes/docs`

---

## Node Registry

The complete node registry with all properties is available at:

- **Current Registry**: `data/n8n-nodes-registry.json`
- **Previous Registry**: `data/n8n-nodes-registry.previous.json`
- **Sync Report**: `NODES_SYNC.md` (generated, not committed)

### Registry Structure

```json
{
  "version": "n8n@1.x.x",
  "last_sync": "2024-01-01T00:00:00Z",
  "nodes": [
    {
      "name": "Code",
      "type": "n8n-nodes-base.code",
      "category": "core",
      "latest_version": 2,
      "description": "Execute custom JavaScript code",
      "inputs": ["main"],
      "outputs": ["main"],
      "file": "dist/nodes/Code/Code.node.js"
    }
  ]
}
```

---

## Maintenance

This documentation is auto-generated from the node registry. To regenerate:

```bash
make nodes/docs
```

**Last Generated**: 2025-11-25T12:59:33.936Z
