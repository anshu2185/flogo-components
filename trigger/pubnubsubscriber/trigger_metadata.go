package pubnubsubscriber

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = `{
  "name": "pubnub-subscriber",
  "type": "flogo:trigger",
  "ref": "github.com/retgits/flogo-components/trigger/pubnubsubscriber",
  "version": "0.2.0",
  "title": "Receive PubNub Messages",
  "description": "PubNub Subscriber",
  "author": "retgits",
  "homepage": "https://github.com/retgits/flogo-components/trigger/pubnubsubscriber",
  "settings":[
      {
        "name": "publishKey",
        "type": "string",
        "required" : true
      },
      {
        "name": "subscribeKey",
        "type": "string",
        "required" : true
      },
      {
        "name": "uuid",
        "type": "string",
        "required" : false
      }
    ],
    "output": [
      {
        "name": "message",
        "type": "any"
      },
      {
        "name": "channel",
        "type": "string"
      },
      {
        "name": "subscription",
        "type": "string"
      },
      {
        "name": "publisher",
        "type": "string"
      },
      {
        "name": "timeToken",
        "type": "string"
      }
    ],
    "handler": {
      "settings": [
        {
          "name": "channel",
          "type": "string",
          "required" : true
        }
      ]
    }
}`

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
