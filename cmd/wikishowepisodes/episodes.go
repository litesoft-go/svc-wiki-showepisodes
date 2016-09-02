package main

import (
	server "lib-builtin/lib/serveradaptor"

	"lib-builtin/lib/channels"
	"lib-builtin/lib/fatal"
	"lib-builtin/lib/cli"

	"svc-push-links/lib/tracking"
)

const (
	VERSION = "1.3"
	PORT = "3000"
	QUEUE = "PushedLinks"
)

// DUP W/ CLI-CACHE-WARMER-CTRL
//var CREDENTIALS = awsProxy.NewStaticCredentials(// IAM: xxx
//	"XXX", // These credentials can xxx
//	"xxx") // AND xxx

func main() {
	cli.New().ShowVersion(VERSION + "\nListening on localhost:" + PORT)

	zShutDown := channels.CreateShutDown()
	defer zShutDown.Shutdown()

	zTracker := tracking.NewTracker(server.LoggingDispatchIssues, zShutDown.AsShutDownChannelAccessor())

	//zLinksProcessor := linksProcessor.NewLinksPusher(zShutDown.AsShutDownChannelAccessor(), zTracker, zQueue, 20)

	fatal.IfErrRaw(zTracker.AddStatusEndpoints(server.NewSimplePathDispatcher(zTracker)).// Tracker 4 Dispatcher == DispatchIssues
	//AddPostHandler("/api/rest/v1/addLinks", linkPusherV1.NewPusher(zTracker, zLinksProcessor).AddLinks).
		ListenAndServe(":" + PORT))
}
