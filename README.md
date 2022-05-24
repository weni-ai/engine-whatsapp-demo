# Engine WhatsApp Demo


## About <a name = "about"></a>

This application is intended to mediate communication between courier and a WhatsApp Business API to enable the use of the same number in multiple communication channels.

## Getting Started <a name = "getting_started"></a>

For this application to work, there must be a github.com/nyaruka/courier and a WhatsApp Business API with which it can communicate and mediate requests.

Also need to have a whatsapp type *channel* created in a *org* on `Rapidpro` where in your configs baseurl is the base url of this application.

### Prerequisites

A system with Go and a mongo database to connect and persist data.


- #### Clone project

```bash
git clone https://github.com/Ilhasoft/engine-whatsapp-demo
```
environment variables

  | Variable              | Required | Default |
  |-----------------------|:--------:|---------|
  | APP_HTTP_PORT         | false    | 9000    |
  | APP_GRPC_PORT         | false    | 7000    |
  | APP_COURIER_BASE_URL  | false    | http://localhost:8000/c/wa |
  | APP_SENTRY_DSN        | false    |    -    |
  | APP_LOG_LEVEL         | false    | debug   |
  | DB_NAME               | false    | whatsapp-router |
  | DB_URI                | false    | mongodb://admin:admin@localhost:27017 |
  | WPP_BASEURL           | true     |    -    |
  | WPP_USERNAME          | true     |    -    |
  | WPP_PASSWORD          | true     |    -    |
  | OIDC_REALM            | false    | gocloak |
  | OIDC_HOST             | false    | http://localhost:8080 |


- #### Run application
```
go run cmd/main.go
```

- #### Setup WhatsApp webhook callback URL
https://developers.facebook.com/docs/whatsapp/api/settings

```
PATCH https://{whatsapp-api-url}/v1/settings/application
```
body:
```json
{
  "webhooks": {
    "url": "https://{engine-whatsap-demo-url}/wr/receive"
  }
}
```

## Usage <a name = "usage"></a>

### Creating a Channel
Call:
```json
gRPC /ChannelService/CreateChannel

Body:
{
	"uuid": "1234-asdf-qwer-qwer",
	"name": "my channel"
}
```
Response:
```json
{
	"token": "weni-demo-BgzokfF65W"
}
```
### Activate token to contact

Start a conversation with the configured contact number from the Whatsapp API and send a message only with the token of a created channel. If the token is valid, the channel will send a confirmation message, and the contact will be able to interact with the number.

### Sending messages
- #### WhatsApp API -> engine-whatsap-demo -> courier

When a contact send a text message to the number configured in the WhatsApp API, it will send a callback http request to the Webhook URL specified in application settings in this case the URL to engine-whatsapp-demo:

```
POST https://foo.bar/wr/receive
```
body:
```json
{
  "contacts": [
    {
      "profile": {
        "name": "user_name"
      },
      "wa_id": "12341341234"
    }
  ],
  "messages": [
    {
      "from": "558299990000",
      "id": "123456",
      "text": {
        "body": "hi dude."
      },
      "timestamp": "623123123123",
      "type": "text"
    }
  ]
}

```

and then engine-whatsapp-demo will send to courier with url to path from registered channel like that:

```
POST https://courier-host.com/c/wa/1234-qwert-asdf-zxcv/receive
```

- #### courier -> engine-whatsapp-demo -> WhatsApp API
When courier send a message to contact in the channel, it will make a request to engine-whatsapp-demo, which will redirect the message to WhatsApp API.
```
POST https://{engine-whatsapp-demo-url}/v1/message
```
Body:
```json
{
  "to": "558299990000",
  "type": "text",
  "text": {
    "body": "Hello World!"
  }
}
```

```
POST https://{whatsapp-business-api}/v1/message
```
