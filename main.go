package main

import (
	"fmt"
	"github.com/ynori7/MusicNewReleases/filter"
	"io/ioutil"
	"log"
	"os"

	"github.com/ynori7/MusicNewReleases/config"
	"github.com/ynori7/MusicNewReleases/music"
)

func main() {
	if len(os.Args) < 2 {
		//log.Fatal(errors.New("you must specify the path to the config file"))
	}

	//Get the config
	data, err := ioutil.ReadFile("config.yml")//os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		log.Fatal(err)
	}

	newReleases, err := music.GetPotentiallyInterestingNewReleases(conf)
	if err != nil {
		log.Fatal(err)
	}

	filterer := filter.NewFilterer(conf, newReleases)
	interestingDiscographies := filterer.FilterAndEnrich()

	for _, r := range interestingDiscographies {
		fmt.Println(r.Artist.Name)
		fmt.Println(r.Artist.Genres)
		fmt.Printf("Best Rating: %d, Avg Rating: %d\n", r.BestRating, r.AverageRating)
		fmt.Println(r.Albums[len(r.Albums)-1]) //the newest one
	}
}
