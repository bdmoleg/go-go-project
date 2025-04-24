# Asana Extractor Exercise

### Rate limiting

In Asana's docs it's specified that regarding [rate limiting](https://developers.asana.com/docs/rate-limits) we are able to do 150 RPM free tier and 1500 RPM paid tier.  

We should check for HTTP Status Code 429 Too Man Requests in order to catch throttling and apply exponential backoff logic.  

When request is getting throttled there is a header "Retry-After" that specifies number of seconds to wait before making a new request in order to not be throttled by Asana.

There are no any secondary rate limiting applied so we can make 150 of requests concurrently without being throttled.

A solution will be to use all RPMs and when some request will return 429 code then get the Retry-After value and wait number of seconds before continue.