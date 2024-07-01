package xkcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func parseComics(num int) (*Comics, error) {
	resp, err := http.Get(XkcdLink + strconv.Itoa(num) + "/info.0.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request error: %s", resp.Status)
	}
	var result Comics
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func SaveToFile() error {
	file, err := os.Create("comics.json")
	if err != nil {
		log.Println("Error creating file")
		return err
	}
	defer file.Close()

	var comics []*Comics
	count, errCount := 0, 0
	started := time.Now()
	fmt.Printf("Starting parsing... (%.19v)\n", started)
	for i := 1; ; i++ {
		result, err := parseComics(i)
		if err != nil {
			log.Printf("Comics parsing error %d: %v", i, err)
			errCount++
			if errCount > 2 {
				ended := time.Now()
				duration := ended.Sub(started)
				fmt.Printf("Parsed %d comics in %.2f sec.\n", count, duration.Seconds())
				break
			}
		}
		comics = append(comics, result)
		count++

	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(comics)
	if err != nil {
		return err
	}
	return nil
}
