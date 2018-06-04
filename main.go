package main

import (
	"BooksScrapper/safaryscraper"
	"log"

	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rawCookies := os.Getenv("COOKIES")

	config := &safaryscraper.Config{
		Url:        "https://www.safaribooksonline.com/library/view/cloud-native-programming/9781787125988/f67b9483-1088-4231-b2da-6087d4750b14.xhtml",
		RawCookies: rawCookies,
	}

	bs := safaryscraper.NewBookScrapper(config)
	bs.GetHtmlPages()
}
