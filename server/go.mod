module server

go 1.18

require (
	github.com/gocolly/colly v1.2.0
	github.com/joho/godotenv v1.4.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/shawnyu5/check_coin_product v1.0.0
	github.com/shawnyu5/penguin-search v1.0.0
	github.com/shawnyu5/penguin-utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/antchfx/htmlquery v1.2.4 // indirect
	github.com/antchfx/xmlquery v1.3.10 // indirect
	github.com/antchfx/xpath v1.2.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

replace github.com/shawnyu5/penguin-utils v0.0.0-00010101000000-000000000000 => ./utils/

replace github.com/shawnyu5/check_coin_product v1.0.0 => ./coin_products

replace github.com/shawnyu5/penguin-search => ./search
