package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	results := make(chan []Clip)
	jobs := make(chan uint64)

	PagesAmount, WorkerAmount, Channel, FilePath := flag.Uint64("pages", 1, "uint64 pages amount"),
		flag.Uint("workers", 1, "int worker amount"),
		flag.String("channel", "twitch", "twitch channel name"),
		flag.String("outputFile", "output.json", "output outputFile for the clips")

	outputFile, err := os.Create(*FilePath)
	outputFile.Write([]byte("["))
	defer outputFile.Close()

	if err != nil {
		fmt.Println("Error creating outputFile:", err.Error())
		return
	}
	flag.Parse()

	for i := uint(0); i < *WorkerAmount; i++ {
		go func() {
			for page := range jobs {
				fmt.Printf("WORKER [%d] PAGE(%d) -> BEGIN JOB\n", i, page)
				clipArr, _ := fetchClip(page, *Channel)

				if len(clipArr) < 1 {
					return
				}

				results <- clipArr
			}
		}()
	}

	go func() {
		for i := uint64(0); i < *PagesAmount; i++ {
			jobs <- uint64(i)
			fmt.Printf("SEND JOB -> PAGE(%d)\n", i)
		}
	}()

	for i := uint64(0); i < *PagesAmount; i++ {
		clips := <-results

		jsonClips, err := json.Marshal(clips)

		if err != nil {
			fmt.Println("Error marshaling json:", err.Error())
		}

		data := jsonClips[1 : len(jsonClips)-1]
		if *PagesAmount > 1 && i != *PagesAmount-1 {
			data = append(data, byte(','))
		}
		n, err := outputFile.Write(data)

		if err != nil {
			fmt.Println("Error writing to outputFile:", err.Error())
		} else {
			fmt.Println("Wrote:", n, "bytes")
		}

	}
	outputFile.Write([]byte("]"))
	close(jobs)

}

func fetchClip(page uint64, channel string) ([]Clip, error) {
	client := http.Client{}
	//fmt.Println("Requesting page:", page)
	var pageSize uint64 = 50
	payload := ClipsPayload{
		OperationName: "ClipsCards__User",
		Variables: Variables{
			Login: channel,
			Limit: int(pageSize),
			Criteria: Criteria{
				Filter: "ALL_TIME",
			},
			Cursor: base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(page*pageSize, 10))),
		},
		Extensions: RequestExtensions{
			PersistedQuery: PersistedQuery{
				Version:    1,
				Sha256Hash: "b73ad2bfaecfd30a9e6c28fada15bd97032c83ec77a0440766a56fe0bd632777",
			},
		},
	}
	payloadJson, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer(payloadJson))

	if err != nil {
		fmt.Println("Failed to create a new http request")
		return nil, err
	}

	request.Header.Add("client-id", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println("Error while sending post request:", err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code is not 200 got:", resp.StatusCode)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading body:", err.Error())
		return nil, err
	}

	var clips TwitchClips

	err = json.Unmarshal(body, &clips)

	if err != nil {
		fmt.Println("Error while unmarshaling json:", err.Error())
		return nil, err
	}

	var normalizedClips []Clip

	for _, node := range clips.Data.User.Clips.Edges {
		normalizedClips = append(normalizedClips, node.Node)
	}

	return normalizedClips, nil
}
