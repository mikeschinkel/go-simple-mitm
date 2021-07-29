# Simple MitM Proxy Example in Go

This is a simple MitM proxy I implemented with a plan to implement mock/fake/stub/spy testing for the GitHub API. 

However, the client — who hired me to do Go work but was primarily a NodeJS shop themselves — discovered [Mountebank](http://mbtest.org) after I had implemented this and wanted their team to use that instead.

I learned a lot about proxies and the HTTPS protocol than I previously knew in order to write this, and I am publishing it on Github in hopes it helps others more quickly get over the humps that I had to get over.

I also hope to be able to revisit this in the future and use this work in production somewhere. 

## Resources
- [HTTP(S) Proxy in Golang in less than 100 lines of code](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c)
- [HTTPS proxies support in Go 1.10](https://medium.com/@mlowicki/https-proxies-support-in-go-1-10-b956fb501d6b)
- [Golang : Disable security check for HTTPS(SSL) with bad or expired certificate](https://www.socketloop.com/tutorials/golang-disable-security-check-for-http-ssl-with-bad-or-expired-certificate)
