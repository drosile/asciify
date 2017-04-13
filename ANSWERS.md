# Answers to some relevant questions

## 1. How would you measure the performance of your service?

We have four components of this application:
- accepting web requests / delivering responses
- reading in image data
- resizing/scaling an image
- converting pixels to ASCII

To measure the performance, we'd want to test how many requests per second
we could satisfy. This would be limited by how fast our image processing is
along with our serialization to create the response.

Another check would be to see what the largest image we could process would be,
or how large we could scale an image. The bounds here would be either response
time (resulting in a timeout) or memory (resulting in some sort of overflow
error--I'm not super familiar with Go so I'm not sure how that would look.)

## 2. What are some strategies you would employ to make your service more scalable?

- more web servers to handle more requests
- separate out image scaling so it can be done separately from the converter
- cache processed images so we don't have to recompute for repeated queries
- allow for batched calls to reduce HTTP overhead
- disallow images over a certain size
- for larger images, use asynchronous web requests so the server is not blocked

One last thing you could do is parallelize the image processing. By this I mean
that different sections of the image can be converted to ASCII independently of
each other, so it's certainly possibly to have an ASCII processor for every
vertical pixel, and to do the conversion in one sweep across horizontally, as
one example.

## 3. How would your design change if you needed to store the uploaded images?

The obvious change is that I'd need some kind of database to do the storage.
Once the images are stored, I'd be able to include references to them in the
API payload, such as `image_id`. It's also possible that I might be able to
handle larger images than I would be able to just by using memory by
chunking the processing.

## 4. What are the cost factors of your scaling choices? Which parts of your solution would grow in cost the fastest?

I've chosen to load images fully into memory and then resize, convert, and
print the output. An alternative would be to do some of this processing in a
stream, such that I'm not building up many variables and passing them around.
However, that would possibly lead to a less encapsulated solution.

If we are talking about the scaling strategies in number 2 above, the most
expensive part would likely be the data transfer between the API server and
the image scaling server. This could be mitigated by allowing images to be
stored on the scaling server so that we could at least avoid the upload step
on duplicate requests (or requests for the same image at different scales).

## 5. Where are your critical points of failure and how would you mitigate them?

Basic failures:

- file not an image: check against content-type
- image too large: reject it before processing, or use one of the strategies above
- resized too wide: put an upper bound on resize
- resized too small (negative?): catch & return error

I haven't implemented any of these checks because it's fun to see what happens,
but it would be important to have boundaries to put this into production.

Scaled solution failures:

- cache grows too big: cap the cache to a particular size
- connection between servers breaks down: solve with monitoring

## 6. How given a change to the algorithm what issues do you foresee when upgrading your scaled-out solution?

It's possible that a new algorithm for ASCII conversion wouldn't rely on
scaling at all, making our scaling image service obsolete.

## 7. If you wanted to migrate your scaled-out solution to another cloud provider (with comparable offerings but different API’s) how would you envision this happening? How would deal with data consistency during the transition and rollbacks in the event of failures?

Most providers have transition plans in terms of deployment APIs and such,
because they want to make it easy for people to leave their competition.

Once the service is established on the new provider, I'd probably set up an
internal tool to record incoming requests and their responses, then replay
the requests against the transitioned service to compare the responses. If
that looked good, I'd put up a load balancer to send some small portion of
new requests to the new service, eventually scaling it up to send all
requests there, and finally shutting down the old service. Or, we use the
phrase "sunsetting" I guess. This process would be reversible at every step
in case things go wrong.

If we are storing the images, then I'd set up something to make the databases
eventually consistent with each other--perhaps by having the new service
start its database IDs at some high number (one billion and one, instead
of one), and synchronizing data on some schedule that would allow for the
use cases we need to support.
