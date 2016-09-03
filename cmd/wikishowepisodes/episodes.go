package main

import (
	//server "lib-builtin/lib/serveradaptor"

	"lib-builtin/lib/channels"
	"lib-builtin/lib/fatal"
	"lib-builtin/lib/cli"

	//"svc-push-links/lib/tracking"
	"svc-wiki-showepisodes/lib/page"

	"io/ioutil"
	"fmt"
)

const (
	VERSION = "1.3"
	PORT = "3000"
)

// DUP W/ CLI-CACHE-WARMER-CTRL
//var CREDENTIALS = awsProxy.NewStaticCredentials(// IAM: xxx
//	"XXX", // These credentials can xxx
//	"xxx") // AND xxx

var EPISODE_PAGES = []string{
	//"https://en.wikipedia.org/wiki/List_of_Agents_of_S.H.I.E.L.D._episodes",
	//"https://en.wikipedia.org/wiki/List_of_Arrow_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Banshee_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Bones_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Dark_Matter_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Fear_the_Walking_Dead_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Hell_on_Wheels_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Killjoys_episodes",
	"https://en.wikipedia.org/wiki/List_of_Major_Crimes_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Powers_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Rizzoli_%26_Isles_episodes",
	//"https://en.wikipedia.org/wiki/List_of_The_Last_Ship_episodes",
	////
	//"https://en.wikipedia.org/wiki/Blindspot_(TV_series)", // #Episodes
	//"https://en.wikipedia.org/wiki/Legends_of_Tomorrow", // Episodes
	//"https://en.wikipedia.org/wiki/Sense8", // Episodes
	//"https://en.wikipedia.org/wiki/Star_Trek_Continues", // Episodes
	////
	//"https://en.wikipedia.org/wiki/Cleverman",
}

func main() {
	cli.New().ShowVersion(VERSION + "\nListening on localhost:" + PORT)

	zShutDown := channels.CreateShutDown()
	defer zShutDown.Shutdown()

	zPages := []*page.ShowEpisodes{}
	for _, zEpisodeLink := range EPISODE_PAGES {
		zPage, err := page.NewWikiShowEpisodes(zEpisodeLink)
		fatal.IfErrRaw(err)
		zPages = append(zPages, zPage)
		fmt.Println(zPage.GetTitle())
	}

	for _, zPage := range zPages {
		err := zPage.Parse(ioutil.ReadFile(zPage.GetTitle() + ".html"))
		fatal.IfErrRaw(err)
	}

	//zTracker := tracking.NewTracker(server.LoggingDispatchIssues, zShutDown.AsShutDownChannelAccessor())

	//zLinksProcessor := linksProcessor.NewLinksPusher(zShutDown.AsShutDownChannelAccessor(), zTracker, zQueue, 20)

	//fatal.IfErrRaw(zTracker.AddStatusEndpoints(server.NewSimplePathDispatcher(zTracker)).// Tracker 4 Dispatcher == DispatchIssues
	////AddPostHandler("/api/rest/v1/addLinks", linkPusherV1.NewPusher(zTracker, zLinksProcessor).AddLinks).
	//	ListenAndServe(":" + PORT))
}
