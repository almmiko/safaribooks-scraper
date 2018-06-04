# Safaribooks Scraper
> Safaribooks Scraper Example.

## Usage

* Rename .env.example -> .env
* Login to safaribooks
* Add you safaribooks cookies to .env file
> Note: Remove all links from your cookies

> You can use developer console for getting cookies -> ```document.cookie```

```
COOKIES=BrowserCookie=...
```

* Add a book path

```
	config := &safaryscraper.Config{
		Url:        "https://www.safaribooksonline.com/library/view/cloud-native-programming/9781787125988/f67b9483-1088-4231-b2da-6087d4750b14.xhtml",
		RawCookies: rawCookies,
	}

	bs := safaryscraper.NewBookScrapper(config)
	bs.GetHtmlPages()
```

## Disclaimer

This project shows how you can scrap web sites using golang. Safarybooks prevents full books scraping via javascript, using safaribooks-scraper you can only get limited page content
