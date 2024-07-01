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

func mapComics() map[*Comics]int {
	started := time.Now()
	comics := make(map[*Comics]int)
	for i := 1; i < 100; i++ {
		result, err := parseComics(i)
		if err != nil {
			log.Printf("Comics parsing error %d: %v", i, err)
			continue // Skip and go next
		}

		if result == nil {
			log.Printf("Nil result for comic %d", i)
			continue
		}
		comics[result] = result.Num
	}
	ended := time.Now()
	duration := ended.Sub(started)
	fmt.Printf("Парсинг завершён за %.4f сек\n", duration.Seconds())
	return comics
}

func SaveToFile() error {
	file, err := os.Create("comics.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var comics []*Comics
	for i := 1; i < 100; i++ {
		result, err := parseComics(i)
		if err != nil {
			log.Printf("Comics parsing error %d: %v", i, err)
			continue // Skip and go next
		}

		if result == nil {
			log.Printf("Nil result for comic %d", i)
			continue
		}
		comics = append(comics, result)
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(comics)
	if err != nil {
		return err
	}
	return nil
}
