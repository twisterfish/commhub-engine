package switchboard

//////////////////////////////////////////////////////////////////////////////////////////
// The master signal switching map for processing
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"authorize"
	"invitations"
	"perishables"
	"pricing"
	"products"
	"profiles"
	"properties"
	"uploads"
	"workspaces"
	"worktickets"

	"github.com/aws/aws-lambda-go/events"
)

//////////////////////////////////////////////////////////////////////////////////////////
//
//////////////////////////////////////////////////////////////////////////////////////////
func RouteSignal(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch signal {
	case "authorize":
		return authorize.DoAction(action, payload)
	case "workspaces":
		return workspaces.DoAction(signal, action, payload, request)
	case "products":
		return products.DoAction(signal, action, payload, request)
	case "pricing":
		return pricing.DoAction(signal, action, payload, request)
	case "properties":
		return properties.DoAction(signal, action, payload, request)
	case "profiles":
		return profiles.DoAction(signal, action, payload, request)
	case "perishables":
		return perishables.DoAction(signal, action, payload)
	case "uploads":
		return uploads.DoAction(signal, action, payload, request)
	case "worktickets":
		return worktickets.DoAction(signal, action, payload, request)
	case "invitations":
		return invitations.DoAction(signal, action, payload, request)
	default:
		return "{\"signal\":\"error\",\"action\":\"invalid_signal...\"}"
	}

}
