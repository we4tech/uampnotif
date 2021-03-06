# Welcome to UampNotif Project

[WIP] More details coming
[WIP] A simplified notification dispatching system.

![UampNotif](https://user-images.githubusercontent.com/4054/102446843-b2dbf480-3ffc-11eb-840e-2a0ccbf08d28.png)

## Build

```bash
make build-mac
```

## Usages

```bash
Usage of ./uampnotif:
  -d string
    	(Required) Locate receivers directory
  -n string
    	(Required) Locate notification.yml file 
```

## Terms

WIP: Confusing. Needs more simplification

- Receiver: A unit of notification recipient
- Notification: A configuration to glue together all notification receivers.

## Example configuration

```yaml
##
# Default settings for all receivers.
#
default_settings:
  retries: 3
  # Do you prefer uampnotif to call all notifcfg concurrently?
  async: true
  # What do you want us to do in case of a failure to notify one of the
  # downstream systems?
  #
  #   Accepted values:
  #    - ignore: Do nothing
  #    - fatal: Exit with 1 if on_error_receivers is set, notify before exiting.
  #    - no_error_receivers: Mute notifying the on_error_receivers
  #
  on_error: ignore
  # TODO: Future feature
  # Do you prefer to notify a list of notifcfg if error arises?
  on_error_receivers:
    - id: slack
      params:
        callback_url: https://httpbin.org/post

##
# A list of notifcfg to be invoked whenever uampnotif is launched.
#
receivers:
  ##
  # id: Must match with corresponding integration id:
  #
  - id: newrelic
    ##
    # You can use go-template with .Env for all param value.
    #
    params:
      callback_url: https://httpbin.org/post
      api_key: '{{.FindEnv "NEWRELIC_API_KEY"}}'

  - id: rollbar
    params:
      access_token: hello-access-token
      environment: staging
      rollbar_username: hello-rollbar-user
      comment: Exception from Rollbar

  - id: slack
    params:
      callback_url: https://httpbin.org/post

  - id: sox-auditor
    params:
      repo_name: test/repo
      secret: hello-world-secret
    ##
    # Override global settings. You can only add the one you may prefer.
    # Uampnotif merges the settings below with the globally exposed settings.
    #
    settings:
      on_error: fatal

```

## Individual receiver spec

```yaml
name: Test
id: test
request:
  valid_http_codes:
    - 201
    - 200
  params:
    - name: app_id
      label: App ID
    - name: api_key
      label: API Key
    - name: hostname
      label: Hostname
  method: POST
  url_tmpl: '{{.FindEnv "NOTIF_SERVER_URL"}}/v2/applications/{{.FindParam "app_id"}}/deployments.json'
  headers:
    - name: Content-Type
      value_tmpl: application/json
    - name: X-Api-Key
      value_tmpl: '{{.FindParam "api_key"}}'
  body_tmpl: |
    { "deployment": { "revision": "{{.FindEnv "COMMIT_HASH"}}" } }
```
