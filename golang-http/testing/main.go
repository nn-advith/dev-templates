package main

import (
	"encoding/json"
	"fmt"
)

type Dimensions struct {
	Height int
	Width  int
}

type Bird struct {
	Species     string
	Description string
	Dimensions  Dimensions
}

func main() {
	fmt.Printf("\n#----Marshalling and Unmarshalling JSON----#\n\n")

	birdJson := `{"species":"pigeon","description":"likes to perch on rocks", "dimensions":{"height":24,"width":10}}`
	var birds Bird
	json.Unmarshal([]byte(birdJson), &birds)

	fmt.Printf("Species: %v, Desc: %v, Dimensions: %v-%v", birds.Species, birds.Description, birds.Dimensions.Height, birds.Dimensions.Width)

	// for _, i := range birds {
	// 	fmt.Printf("Species:%v \tDescription:%v\n", i.Species, i.Description)
	// }

	// fmt.Printf("Birds : %+v\n", birds)
	// fmt.Println(reflect.TypeOf(birds))

	// fmt.Printf("Species: %s, Description: %s", bird.Species, bird.Description)
	//Species: pigeon, Description: likes to perch on rocks

}
