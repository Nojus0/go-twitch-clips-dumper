package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var ErrLimitReachedError = errors.New("twitch 100,000 clip limit reached")

const PAGE_SIZE_LIMIT = ((100_000 - PageSize) / PageSize)

const PageSize = 100

func fetchClip(page uint64, channel string) ([]Clip, error) {
	client := http.Client{}

	payload := ClipsPayload{
		OperationName: "ClipsCards__User",
		Variables: Variables{
			Login: channel,
			Limit: int(PageSize),
			Criteria: Criteria{
				Filter: "ALL_TIME",
			},
			Cursor: base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(page*PageSize, 10))),
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

	body, err := io.ReadAll(resp.Body)

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

	var rawClips []Clip

	if clips.Errors != nil {
		return nil, ErrLimitReachedError
	}

	for _, node := range clips.Data.User.Clips.Edges {
		rawClips = append(rawClips, node.Node)
	}

	return rawClips, nil
}
