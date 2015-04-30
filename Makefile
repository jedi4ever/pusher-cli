build:
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go build pusher-cli.go

install:
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go get github.com/PuerkitoBio/goquery
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go get golang.org/x/net/publicsuffix
