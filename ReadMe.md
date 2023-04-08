# Dental Vocabulary Image Scraper

This is a Go program that scrapes images from a csv file using the Bing Image Search API.

## Installation

To use this program, you must first obtain a subscription key from the Bing Image Search API.

1. Go to the [Bing Image Search API](https://www.microsoft.com/en-us/bing/apis/bing-image-search-api) and create an account.
2. Follow the instructions to create a subscription and obtain a subscription key.

Once you have your subscription key, you can download and run the program:

$ go get github.com/Erkhemee0908/csv-image-scraper
$ go run main.go

## Usage

This program requires a CSV file containing a list of dental vocabulary terms. The file should have one term per line and no header row. An example file is included in the repository (`dental_vocabulary.csv`).

When you run the program, it will loop through the list of terms in the CSV file and scrape an image for each term. The images will be saved to a directory called `img` in the current working directory, and the program will add a new column to the CSV file with the filenames of the corresponding images.