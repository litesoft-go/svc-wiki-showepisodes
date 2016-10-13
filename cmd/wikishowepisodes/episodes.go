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
	"os"
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

	// SOT Working:
	"https://en.wikipedia.org/wiki/List_of_Agents_of_S.H.I.E.L.D._episodes",
	"https://en.wikipedia.org/wiki/List_of_Arrow_episodes",
	"https://en.wikipedia.org/wiki/List_of_Banshee_episodes",
	"https://en.wikipedia.org/wiki/List_of_Between_episodes", // S01E06 S02E01 - ?? (2016)
	"https://en.wikipedia.org/wiki/List_of_Bones_episodes",
	"https://en.wikipedia.org/wiki/List_of_Code_Black_episodes", // S01E18
	"https://en.wikipedia.org/wiki/List_of_Dark_Matter_episodes",
	"https://en.wikipedia.org/wiki/List_of_Elementary_episodes", // S04E24
	"https://en.wikipedia.org/wiki/List_of_Fargo_episodes", // S02E10
	"https://en.wikipedia.org/wiki/List_of_Fear_the_Walking_Dead_episodes",
	"https://en.wikipedia.org/wiki/List_of_Greys_Anatomy_episodes", // S12E24
	"https://en.wikipedia.org/wiki/List_of_Grimm_episodes", // S05E22
	"https://en.wikipedia.org/wiki/List_of_Halt_and_Catch_Fire_episodes", // S02E10 S03E01 - 10 (2016)
	"https://en.wikipedia.org/wiki/List_of_Hawaii_Five-0_episodes", // S06E25
	"https://en.wikipedia.org/wiki/List_of_Homeland_episodes", // S05E12
	"https://en.wikipedia.org/wiki/List_of_How_to_Get_Away_with_Murder_episodes", // S02E15
	"https://en.wikipedia.org/wiki/List_of_iZombie_episodes", // S02E19
	"https://en.wikipedia.org/wiki/List_of_Jane_the_Virgin_episodes", // S02E22
	"https://en.wikipedia.org/wiki/List_of_Killjoys_episodes",
	"https://en.wikipedia.org/wiki/List_of_Life_in_Pieces_episodes", // S01E22
	"https://en.wikipedia.org/wiki/List_of_Madam_Secretary_episodes", // S02E23
	"https://en.wikipedia.org/wiki/List_of_Major_Crimes_episodes",
	"https://en.wikipedia.org/wiki/List_of_Modern_Family_episodes", // S07E22
	"https://en.wikipedia.org/wiki/List_of_Mr_Robot_episodes", // S01E10 S02E01 - ?? (2016 June?)
	"https://en.wikipedia.org/wiki/List_of_Powers_episodes",
	"https://en.wikipedia.org/wiki/List_of_Orphan_Black_episodes", // S04E10
	"https://en.wikipedia.org/wiki/List_of_Person_of_Interest_episodes", // S05E13
	"https://en.wikipedia.org/wiki/List_of_Quantico_episodes", // S01E22
	"https://en.wikipedia.org/wiki/List_of_Reign_episodes", // S03E10
	"https://en.wikipedia.org/wiki/List_of_Rizzoli_%26_Isles_episodes",
	"https://en.wikipedia.org/wiki/List_of_Scorpion_episodes", // S02E24
	"https://en.wikipedia.org/wiki/List_of_Sleepy_Hollow_episodes", // S03E18
	"https://en.wikipedia.org/wiki/List_of_Stitchers_episodes", // S02E10
	"https://en.wikipedia.org/wiki/List_of_Switched_at_Birth_episodes", // S04E20
	"https://en.wikipedia.org/wiki/List_of_The_Flash_episodes", // S02E23
	"https://en.wikipedia.org/wiki/List_of_The_Last_Ship_episodes",
	"https://en.wikipedia.org/wiki/List_of_The_Mysteries_of_Laura_episodes", // S02E16
	"https://en.wikipedia.org/wiki/List_of_The_Strain_episodes", // S02E13 S03E01 - 10 (2016)
	"https://en.wikipedia.org/wiki/List_of_True_Detective_episodes", // S02E08 S03???

	"https://en.wikipedia.org/wiki/Crossing_Lines", // S03E12
	"https://en.wikipedia.org/wiki/Rosewood_(TV_series)", // S01E22
	"https://en.wikipedia.org/wiki/Through_the_Wormhole", // S06E06
	//
	//
	// SOT, NOT Working - Actuals - No Match:
	"https://en.wikipedia.org/wiki/List_of_Broadchurch_episodes", // S02E08 (?? ??? ????) S03E01
	//"https://en.wikipedia.org/wiki/List_of_Call_the_Midwife_episodes", // S05E08
	//"https://en.wikipedia.org/wiki/List_of_Doc_Martin_episodes", // S07E08
	//"https://en.wikipedia.org/wiki/List_of_Game_of_Thrones_episodes", // S06E10
	//"https://en.wikipedia.org/wiki/List_of_Grace_and_Frankie_episodes", // S02E13 (2017 May?)
	//"https://en.wikipedia.org/wiki/List_of_Grand_Designs_episodes", // S16E09
	//"https://en.wikipedia.org/wiki/List_of_Lexx_episodes", // S01-04
	//"https://en.wikipedia.org/wiki/List_of_Longmire_episodes", // S04E10 S05E01 - ?? (2016)
	//"https://en.wikipedia.org/wiki/List_of_Orange_Is_the_New_Black_episodes", // S03E13 S04E01 - ?? (2016)
	//"https://en.wikipedia.org/wiki/List_of_Supergirl_episodes", // S01E20
	//"https://en.wikipedia.org/wiki/List_of_The_Big_Bang_Theory_episodes", // S09E24
	//"https://en.wikipedia.org/wiki/List_of_The_Good_Wife_episodes", // S07E22
	//"https://en.wikipedia.org/wiki/List_of_The_X-Files_episodes", // S10E06
	////
	////
	//// SOT, NOT Working - Headers Match - Body Errors:
	//"https://en.wikipedia.org/wiki/List_of_CSI: Cyber_episodes", // S02E18
	//"https://en.wikipedia.org/wiki/List_of_Grand_Designs_Australia_episodes", // S0?E??
	//"https://en.wikipedia.org/wiki/List_of_Guardians_of_the_Galaxy_episodes", // S01E10
	//"https://en.wikipedia.org/wiki/List_of_Hell_on_Wheels_episodes",
	//"https://en.wikipedia.org/wiki/List_of_Last_Man_Standing_episodes", // S05E21
	//"https://en.wikipedia.org/wiki/List_of_NCIS_episodes", // S12E24
	//"https://en.wikipedia.org/wiki/List_of_NCIS: Los Angeles_episodes", // S07E24
	//"https://en.wikipedia.org/wiki/List_of_NCIS: New Orleans_episodes", // S02E24
	//"https://en.wikipedia.org/wiki/List_of_Outlander_episodes", // S02E13
	//"https://en.wikipedia.org/wiki/List_of_The_Librarians_episodes", // S02E10
	//"https://en.wikipedia.org/wiki/List_of_The_Walking_Dead_episodes", // S06E16
	//
	//
	// Direct Season Table:
	//"https://en.wikipedia.org/wiki/List_of_Colony_episodes", // 	  	 S01E10
	//"https://en.wikipedia.org/wiki/List_of_Daredevil_episodes", // S02 S03???
	//"https://en.wikipedia.org/wiki/List_of_Shadowhunters_episodes", // S01E13
	//"https://en.wikipedia.org/wiki/List_of_Minority_Report_episodes", // S01E10
	//"https://en.wikipedia.org/wiki/List_of_The_Shannara_Chronicles_episodes", // S01E10
	//
	//"https://en.wikipedia.org/wiki/Agent_Carter_(TV_series)", // S02E10
	//"https://en.wikipedia.org/wiki/Blindspot_(TV_series)", // #Episodes
	//"https://en.wikipedia.org/wiki/Cleverman",
	//"https://en.wikipedia.org/wiki/Designated_Survivor_(TV_series)", // S01E01
	//"https://en.wikipedia.org/wiki/Doctor_Thorne_(TV_series)", // S01E03
	//"https://en.wikipedia.org/wiki/Humans_(TV_series)", // S01E08 S02E01 - 08 (2016)
	//"https://en.wikipedia.org/wiki/Legends_of_Tomorrow", // Episodes
	//"https://en.wikipedia.org/wiki/Limitless_(TV_series)", // S01E22
	//"https://en.wikipedia.org/wiki/Luke_Cage_(TV_series)", // S01E13
	//"https://en.wikipedia.org/wiki/MacGyver_(2016_TV_series)", // S01E01
	//"https://en.wikipedia.org/wiki/Pitch_(TV_series)", // S01E01
	//"https://en.wikipedia.org/wiki/Second_Chance_(2016_TV_series)", // S01E11
	//"https://en.wikipedia.org/wiki/Sense8", // Episodes
	//"https://en.wikipedia.org/wiki/Shades_of_Blue_(TV_series)", // S01E13
	//"https://en.wikipedia.org/wiki/Speechless_(TV_series)", // S01E01
	//"https://en.wikipedia.org/wiki/Stan_Lee%27s_Lucky_Man", // S01E10
	//"https://en.wikipedia.org/wiki/Star_Trek_Continues", // Episodes
	//"https://en.wikipedia.org/wiki/The_Expanse_(TV_series)", // S01E10
	//"https://en.wikipedia.org/wiki/The_Grinder_(TV_series)", // S01E22
	//"https://en.wikipedia.org/wiki/The_Nightmare_Worlds_of_H._G._Wells", // S01E04
	//"https://en.wikipedia.org/wiki/Timeless_(TV_series)", // S01E01
	//"https://en.wikipedia.org/wiki/Westworld_(TV_series)", // S01E01
	//
	// "Machines: How They Work", // ??????
	// "Treehouse Masters", // S05E09
	// "Tripped", // S01E04 (channel 4 BBC)
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
		zFileName := zPage.GetTitle() + ".html"
		_, err := os.Stat(zFileName);
		if os.IsNotExist(err) {
			err = zPage.PullAndParse(zFileName)
		} else {
			err = zPage.Parse(ioutil.ReadFile(zFileName))
		}
		if err != nil {
			fmt.Println("******", zPage.GetTitle(), "| Error:", err)
		}
		//fatal.IfErrRaw(err)
		fmt.Println(zPage)
	}

	//zTracker := tracking.NewTracker(server.LoggingDispatchIssues, zShutDown.AsShutDownChannelAccessor())

	//zLinksProcessor := linksProcessor.NewLinksPusher(zShutDown.AsShutDownChannelAccessor(), zTracker, zQueue, 20)

	//fatal.IfErrRaw(zTracker.AddStatusEndpoints(server.NewSimplePathDispatcher(zTracker)).// Tracker 4 Dispatcher == DispatchIssues
	////AddPostHandler("/api/rest/v1/addLinks", linkPusherV1.NewPusher(zTracker, zLinksProcessor).AddLinks).
	//	ListenAndServe(":" + PORT))
}
