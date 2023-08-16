package main

type conf struct {
	Player player_conf
	Screen screen_conf
}

type player_conf struct {
	Enabled bool
}

type screen_conf struct {
	Width  int
	Height int
}

var BackgroundDimensions img_dims = get_image_dimensions(StartImageSrc)

var Config = conf{
	Player: player_conf{
		Enabled: false,
	},
	Screen: screen_conf{
		Width:  BackgroundDimensions.Width,
		Height: BackgroundDimensions.Height,
	},
}

// import (
// 	"log"
// 	"os"

// 	"gopkg.in/yaml.v3"
// )

// func (c *conf) getConf() *conf {

// 	data, err := os.ReadFile("config.yml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = yaml.Unmarshal(data, c)
// 	if err != nil {
// 		log.Fatalf("Unmarshal: %v", err)
// 	}

// 	return c
// }
