# messaging_server
A basic golang implementation of a messaging server using gorilla/websocket and a (very) simple javascript client implemented in a single HTML file.

## Get it running
- You will need a go environnment set up, check out: https://golang.org/doc/install 
- Clone this repo: `git clone https://github.com/Csnigs/messaging_server.git` and cd into it
- Build the app: `go build`
- Copy the sample config file and rename it `cp ./config/config.json.sample ./config/config.json`
- Check out the options `./messaging_server -h`
- Run it: `./messaging_server`
- Open multiple tabs in a web browser http://127.0.0.1:8089/
- I am using Ginkgo and Gomega for my tests. Check it out there: https://github.com/onsi/ginkgo | http://onsi.github.io/gomega/

Keep in mind that the solution is hacky in some places to cover for the lack of real client.
