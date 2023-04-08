package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	subscriptionKey := "YOURKEY"
	csvFilePath := "input.csv" // Sample csv file

	strings, err := getWordList(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, str := range strings {
		time.Sleep(200 * time.Millisecond) // Wait 0.2 seconds
		query := replaceSpaces(str)
		imageUrl, err := getImage(query, subscriptionKey)
		if err != nil {
			panic(err)
		}
		err = downloadImage(imageUrl, str+".jpg")
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	err = addJpgColumn(csvFilePath)
	if err != nil {
		panic(err)
	}
}

func replaceSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "+")
}

func getWordList(csvFilePath string) ([]string, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var strings []string

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		strings = append(strings, row[0])
	}

	return strings[1:], nil
}

func getImage(query string, subscriptionKey string) (string, error) {
	url := "https://api.bing.microsoft.com/v7.0/images/search?q=" + query + "&count=1&mkt=en-US&safeSearch=Strict"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", subscriptionKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	if images, ok := data["value"].([]interface{}); ok && len(images) > 0 {
		if firstImage, ok := images[0].(map[string]interface{}); ok {
			if contentUrl, ok := firstImage["contentUrl"].(string); ok {
				return contentUrl, nil
			}
		}
	}
	return "", nil
}

func downloadImage(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll("img", 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join("img", filename))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func addJpgColumn(inputFile string) error {
	outputFile := "output.csv"
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		return err
	}

	newHeader := make([]string, len(header)+1)
	copy(newHeader, header[:1])
	newHeader[1] = header[0] + ".jpg"
	copy(newHeader[2:], header[1:])

	f2, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f2.Close()
	writer := csv.NewWriter(f2)
	defer writer.Flush()

	if err := writer.Write(newHeader); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		newRecord := make([]string, len(record)+1)
		copy(newRecord, record[:1])
		newRecord[1] = record[0] + ".jpg"
		copy(newRecord[2:], record[1:])

		if err := writer.Write(newRecord); err != nil {
			return err
		}
	}

	return nil

}
