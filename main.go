package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	results := make(chan []Clip)
	jobs := make(chan uint64)

	pages, workers, channel, filePath :=
		flag.Uint64("pages", 1, "Clips max pages"),
		flag.Uint("workers", 1, "Workers amount"),
		flag.String("channel", "twitch", "Twitch channel name"),
		flag.String("file", "output.json", "File destination")

	flag.Parse()

	outStream, err := os.Create(*filePath)
	outStream.Write([]byte("["))

	if err != nil {
		panic(err)
	}

	// This whole proccess can be encapsulated in a struct,
	// then you won't need to pass the same arguments to multiple functions

	for id := uint(0); id < *workers; id++ {
		go Worker(jobs, results, *channel, id)
	}
	go JobSender(jobs, *pages)
	Writer(*pages, results, outStream)

}

func JobSender(jobs chan uint64, pages uint64) {
	for i := uint64(0); i < pages; i++ {
		fmt.Printf("Task -> Requesting Page %d\n", i)
		jobs <- uint64(i)
	}
	close(jobs)
}

func Worker(jobs chan uint64, results chan []Clip, channel string, id uint) {

	for page := range jobs {

		fmt.Printf("Worker[%d] -> Fetching page %d\n", id, page)
		clipArr, _ := fetchClip(page, channel)

		if len(clipArr) < 1 {
			return
		}

		results <- clipArr
	}

}

func Writer(pages uint64, results chan []Clip, outStream *os.File) {

	for i := uint64(0); i < pages; i++ {

		clips, err := json.Marshal(<-results)

		if err != nil {
			fmt.Println("Error marshaling json:", err.Error())
		}

		data := clips[1 : len(clips)-1]

		if pages > 1 && i != pages-1 {
			data = append(data, byte(','))
		}

		n, err := outStream.Write(data)

		if err != nil {
			fmt.Println("Error writing to outputFile:", err.Error())
			continue
		}

		fmt.Println("Writing -> :", n/1000, "KB")

	}

	outStream.Write([]byte("]"))
	outStream.Close()
}
