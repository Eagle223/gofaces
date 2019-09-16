package imageprocess

import (
	"fmt"
	"gofaces/dlib_api"
	"log"
	"path/filepath"
)

const dataDir = "./data"

func ImageProcess() {
	rec, err := dlib_api.NewRecognizer(dataDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	defer rec.Close()
	testImagePristin := filepath.Join(dataDir, "pristin.jpg")
	faces, err := rec.RecognizeFile(testImagePristin)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	var samples []dlib_api.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		cats = append(cats, int32(i))
	}
	labels := []string{
		"Sungyeon", "Yehana", "Roa", "Eunwoo", "Xiyeon",
		"Kyulkyung", "Nayoung", "Rena", "Kyla", "Yuha",
	}
	rec.SetSamples(samples, cats)
	testImageNayong := filepath.Join(dataDir, "nayoung.jpg")
	nayoungFace, err := rec.RecognizeFile(testImageNayong)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if nayoungFace == nil {
		log.Fatalf("Not a single face on the image")
	}
	catID := rec.Classify(nayoungFace[0].Descriptor)
	if catID < 0 {
		log.Fatalf("Can't classify")
	}
	fmt.Println(labels[catID])
}
