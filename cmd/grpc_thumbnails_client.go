package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/OksidGen/grpc_thumbnail/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"image"
	"image/jpeg"
	"log"
	"net/url"
	"os"
	"sync"
)

func main() {
	serverAddr := flag.String("server_addr", "localhost:50051", "The server address in the format host:port")
	async := flag.Bool("async", false, "Enable asynchronous download")
	outputDir := flag.String("output_dir", "./thumbnails", "Output directory for saving images")
	flag.Parse()

	videoURLs := flag.Args()
	numRequests := len(videoURLs)

	if numRequests == 0 {
		log.Fatal("No video URLs provided")
	}

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := proto.NewThumbnailServiceClient(conn)

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Printf("Failed to create directory: %v", err)
		return
	}

	doRequest := func(videoUrl string) {

		response, err := client.GetThumbnail(context.Background(), &proto.ThumbnailRequest{VideoUrl: videoUrl})
		if err != nil {
			log.Printf("Failed to get thumbnail: %v", err)
			return
		}

		img, _, err := image.Decode(bytes.NewReader(response.ThumbnailData))
		if err != nil {
			log.Printf("Failed to decode thumbnail: %v", err)
			return
		}

		videoID, err := extractVideoID(videoUrl)
		if err != nil {
			log.Printf("Failed to extract video ID: %v", err)
			return
		}

		filePath := fmt.Sprintf("%s/%s.jpg", *outputDir, videoID)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		if err := jpeg.Encode(file, img, nil); err != nil {
			log.Printf("Failed to save image: %v", err)
		} else {
			fmt.Printf("Image saved for video %s\n", videoUrl)
		}
	}

	if *async {
		var wg sync.WaitGroup
		for _, videoURL := range videoURLs {
			wg.Add(1)
			go func(videoURL string) {
				defer wg.Done()
				doRequest(videoURL)
			}(videoURL)
		}

		wg.Wait()
	} else {
		for _, videoURL := range videoURLs {
			doRequest(videoURL)
		}
	}
}

func extractVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	videoID := query.Get("v")

	if videoID == "" {
		segments := parsedURL.Path[1:]
		if len(segments) == 11 {
			videoID = segments
		} else {
			return "", fmt.Errorf("cannot extract video ID from URL: %s", videoURL)
		}
	}

	return videoID, nil
}
