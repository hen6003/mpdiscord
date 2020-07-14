package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/hugolgst/rich-go/client"
)

func main() {
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = client.Login("[REMOVED]")
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	random := ""
	// Loop printing the current status of MPD.
	for {
		status, err := conn.Status()
		if err != nil {
			log.Fatalln(err)
		}
		song, err := conn.CurrentSong()
		if err != nil {
			log.Fatalln(err)
		}

		if status["random"] == "1" {
			random = "On"
		} else {
			random = "Off"
		}

		songNumber, err := strconv.ParseInt(status["song"], 10, 64)
		songNumber++
		songNumberStr := strconv.FormatInt(songNumber, 10)

		songName := song["file"]
		songNameSplit := strings.Split(songName, "/")
		songName = songNameSplit[len(songNameSplit)-1]
		songNameSplit = strings.Split(songName, ".")
		songName = songNameSplit[0]

		songTime := strings.Split(status["time"], ":")

		songStartNum, err := strconv.ParseInt(songTime[0], 10, 64)
		songStart := time.Now().Unix() - songStartNum

		songEnd, err := strconv.ParseInt(songTime[1], 10, 64)
		songEnd = time.Now().Unix() + songEnd - songStartNum

		songStartUnix := time.Unix(songStart, 0)
		songEndUnix := time.Unix(songEnd, 0)

		state := ""
		if status["state"] != "play" {
			state = "Paused - "
		}

		state += songNumberStr + "/" + status["playlistlength"] + " - Random: " + random

		err = client.SetActivity(client.Activity{
			State:   state,
			Details: songName,

			Timestamps: &client.Timestamps{
				Start: &songStartUnix,
				End:   &songEndUnix,
			},
		})
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(1e9)
	}
}
