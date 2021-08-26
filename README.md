# Whatsapp Router


## About <a name = "about"></a>

This application is intended to mediate the communication between courier and whatsapp server

## Getting Started <a name = "getting_started"></a>

For this application to work, there must be a github.com/nyaruka/courier and a whatsapp server with which it can communicate and mediate requests.

Also need to have a whatsapp type *channel* created in a *org* on `Rapidpro` where in your configs baseurl is the base url of this application.

### Prerequisites

A system with Go and a mongo database to connect and persist data. 




Clone project

```bash
git clone https://github.com/Ilhasoft/whatsapp_router
```
fill in the .env file based on the .env.example from the config dir
```env
SERVER_HTTP_PORT=
SERVER_GRPC_PORT=

DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_APP_NAME=

WPP_BASEURL=
WPP_USERNAME=
WPP_PASSWORD=
WPP_AUTHTOKEN=
```
run application
```
go run cmd/main.go
```

## Usage <a name = "usage"></a>

### whatsapp -> router -> courier

when a contact send messate to the number from the whatsapp server, the server will make a request like this to whatsapp router:

```
POST https://foo.bar/wr/receive
```
with a payload like that:
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

and then whatsapp router will send to courier with url to path from registered channel like that:

```
POST https://courier-host.com/c/wa/1234-qwert-asdf-zxcv/receive
```

### courier -> router -> whatsapp

```
POST https://foo.bar/v1/message
```

```
POST https://wa-host.example.mobi
```
