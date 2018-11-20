// Sample video_quickstart uses the Google Cloud Video Intelligence API to label a video.
// Sriram Kaushik - Techops

package main

import (
	video "cloud.google.com/go/videointelligence/apiv1"
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
	"io/ioutil"
	"log"
)

func main() {

	var input *string = flag.String("input", "IMG_4387.MOV", "input video to analyze")
	flag.Parse()

	ctx := context.Background()

	//create a client

	client, err := video.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	//Write code for input content file here:

	video_data, err := ioutil.ReadFile(*input) //dereference the pointer

	//_, err = ioutil.ReadFile(*input) //dereference the pointer. Comment if using a hardcoded URL below

	if err != nil {
		log.Fatalf("Unable to decode the video input file: %v", err)
	}

	//uncomment the below InputUri definition if you want to check the hardcoded URL.
	//request := &videopb.AnnotateVideoRequest{InputUri: "gs://demomaker/cat.mp4", Features: []videopb.Feature{videopb.Feature_LABEL_DETECTION}}

	request := &videopb.AnnotateVideoRequest{InputContent: video_data, Features: []videopb.Feature{videopb.Feature_LABEL_DETECTION}}

	op, err := client.AnnotateVideo(ctx, request)

	if err != nil {
		log.Fatalf("Failed to start annotator job:%v", err)
	}

	resp, err := op.Wait(ctx)

	if err != nil {
		log.Fatalf("Failed to annotate: %v", err)
	}

	//Video processed. Let's analyze the results

	result := resp.GetAnnotationResults()[0] //first result

	for _, annotation := range result.SegmentLabelAnnotations { //all entities found are listed as in array
		fmt.Printf("Description: %s\n", annotation.Entity.Description) //description of the entity found in the video with the ID

		for _, category := range annotation.CategoryEntities {
			fmt.Printf("\tCategory: %s\n", category.Description) //description of the category of the entity found in the video
		}

		for _, segment := range annotation.Segments { //start time and end time of the entities.
			start, _ := ptypes.Duration(segment.Segment.StartTimeOffset)
			end, _ := ptypes.Duration(segment.Segment.EndTimeOffset)
			fmt.Printf("\tSegment: %s to %s\n", start, end)
			fmt.Printf("\tConfidence: %v\n", segment.Confidence)
		}
	}

	return
}
