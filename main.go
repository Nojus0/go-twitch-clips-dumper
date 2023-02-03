package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	results := make(chan []Clip)
	jobs := make(chan uint64)

	pages, workers, channel, filePath :=
		flag.Uint64("pages", 1, "Clips max pages"),
		flag.Uint("workers", 1, "Workers amount"),
		flag.String("channel", "twitch", "Twitch channel name"),
		flag.String("file", "output.csv", "File destination")

	flag.Parse()

	if *pages > PAGE_SIZE_LIMIT {
		fmt.Printf("Page size limit is: %d\n", PAGE_SIZE_LIMIT)
		os.Exit(1)
	}

	file, err := os.OpenFile(*filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}

	fileStat, err := file.Stat()

	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(file)
	if fileStat.Size() < 1 {
		err = csvWriter.Write([]string{
			"id",
			"slug",
			"title",
			"views",
			"curatorId",
			"curatorLogin",
			"curatorDisplayName",
			"gameId",
			"gameName",
			"thumbnailUrl",
			"createdAt",
			"duration",
		})

		if err != nil {
			panic(fmt.Sprintf("Could not write csv header: %s", err.Error()))
		}
	}

	// This whole proccess can be encapsulated in a struct,
	// then you won't need to pass the same arguments to multiple functions

	for id := uint(0); id < *workers; id++ {
		go Worker(jobs, results, *channel, id)
	}
	go JobSender(jobs, *pages)
	Writer(*pages, results, csvWriter)

	file.Close()
	csvWriter.Flush()
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
		clipArr, err := fetchClip(page, channel)

		if err != nil || len(clipArr) < 1 {
			panic(err)
		}

		results <- clipArr
	}

}

func Writer(pages uint64, results chan []Clip, csvWriter *csv.Writer) {

	for i := uint64(0); i < pages; i++ {

		clips := <-results

		var clipsCsv [][]string

		for _, clip := range clips {
			clipsCsv = append(clipsCsv, []string{
				clip.ID,
				clip.Slug,
				clip.Title,
				strconv.Itoa(clip.ViewCount),
				clip.Curator.ID,
				clip.Curator.Login,
				clip.Curator.DisplayName,
				clip.Game.ID,
				clip.Game.Name,
				clip.ThumbnailURL,
				clip.CreatedAt.String(),
				strconv.Itoa(clip.DurationSeconds),
			})
		}

		csvWriter.WriteAll(clipsCsv)

	}

}
