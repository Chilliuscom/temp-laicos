package websoct

var NotificationHub *Hub

func Notify(eventname string, from int, where string, usernames []string) {
	//eventpayload.name = event name - followrequest, invitation, joinrequest, newevent, like, comment, post "... has sent you a NAME"
	//eventpayload.from = user_id (USERNAME has sent you a ...)
	//eventpayload.where = group/event id (... has requested to join your group WHERE)
	//eventpayload.link = contains a link to redirect to the event/comment/post or API endpoint that lets user accept or decline straight from notification
	//eventpayload.to = the user_id that should receive this notification
	var payload SocketEventStruct
	payload.EventName = "notification"
	payload.EventPayload = map[string]interface{}{
		"name":  eventname,
		"from":  from,
		"where": where,
	}
	for i := range usernames {
		EmitToSpecificClient(NotificationHub, payload, getUserIDByUsername(NotificationHub, usernames[i]))

	}
}
