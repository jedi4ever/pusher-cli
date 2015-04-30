CLI to read pusher appIds, Keys and Secrets
===================================================
* Totally work in progress - And also my first golang app :) *

required environment variables
- $PUSHER_EMAIL : your pusher-account email
- $PUSHER_PASSWORD : your pusher-account password

```
export PUSHER_EMAIL = "me@somedomain.com"
export PUSHER_PASSWORD = "supersecret"
./pusher-cli
```

Output:
```
my app|97128|60b08baa60e396c3f|2dd2d3bcf3de234dc311
steve app|68454|728708e8ff7443ff2|0d01d28541419b29ee64
```
