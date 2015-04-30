CLI to read pusher appIds, Keys and Secrets
===================================================
* Totally work in progress - And also my first golang app :) *

required environment variables
- $PUSHER_EMAIL : your pusher-account email
- $PUSHER_PASSWORD : your pusher-account password


Currently only lists all appIds, Keys and Secrets.
Soon will work on create & delete & find specifc Pair

# Usage:
```
$ export PUSHER_EMAIL="me@somedomain.com"
$ export PUSHER_PASSWORD="supersecret"
$ ./pusher-cli
```

Output:
```
my app|97128|60b08baa60e396c3f|2dd2d3bcf3de234dc311
steve app|68454|728708e8ff7443ff2|0d01d28541419b29ee64
```

# Compile
(sorry I've still to understand the standard golang project directory structure)
## Dependencies
- "github.com/PuerkitoBio/goquery"
- "golang.org/x/net/publicsuffix"

`make install`

## Build
`make`

