# Explaining the Pros and Cons of the Design Choices of this Project

## Architecture Choices

Hexagonal architecture approach for modularity benefits (client, processor, database, converter)

### Why?

- While this project implements Gemini client in `internal/client/gemini`, it's easy to swap AI providers later.
  An interface-based approach would make this even more flexible when the need arises. This approach keeps the core logic from being tampered with.
- Testing is easier as a result of the clear boundaries.
- Simple, explicit dependencies

### Alternatives

- Single package approach with everything in package main.
  Rejected as it's unfit for this kind of project. Makes it impossible to reuse components and testing is going to be a nightmare.

## Concurrency Strategy

Configurable concurrency API calls with worker pools, waitgroups and channels, defaulting to conservative settings for Gemini free tier

### Why?

Gemini's free tier has five requests per minute limit (12 seconds between requests). Under this constraint, multiple workers provide no speed benefit since the rate limiter forces sequential execution.

However, I decided to go wuth the concurrency infrastructure with configurable workers because of the following:

- Users on paid tiers can set `NUM_WORKERS=10` experience immediate 10x speed.
- Scalable and does not need to be refactored later when the needs arise.

### Alternatives

- Could have gone the simpler route of keeping everything synchronous because of the strict free tier limit, but that is neither scalable nor future proof. As such, I thought the added complexity of concurrency was worth it. No need to refactor for higher rate limits. Just adjust the configuration.

## Batching

Dynamic batch sizing based on daily API limits (default of about 156 tweets per batch)

### Why?

Again, this decision was informed by the limits of the Gemini free tier. Free tier allows only 20 requests per day. With 3128 tweets, batch size must be 3128 ÷ 20 (156 tweets) per request to process everything in one day. Again, batch size is configurable. The batch size could potentially be larger given the 1M token input limit on the free tier.

Users can override batch size for different use cases and depending on the tier.

A smaller batch (20-50) requres paid tier (100+ requests/day) and offers better granularity, better error isolation and arguably more detailed analysis. It would however require a lot more API calls that would exceed the free tier's 20 requests/day. Only makes sense with paid tier.

### Alternatives

- Smaller batch sizes require more API calls exceeding daily limits and reducing API efficiency.

## Error Handling Strategy

Retry with exponential backoff, log and continue

### Why?

Resilience. Without a error handling strategy, the code fails on the first network downtime.
Implemented a 3 retry attempts with 1s, 4s, 9s backoff (exponential). Added a smart retry logic taht knows not to retry invalid API keys but attempts to retry rate limits. Approach respects cancellation during retries.

Smart retry logic ensures the system doesn't waste retries on unrecoverable errors. The system is resilient to transient API blips.

## Technology Choices

### Why Go?

- Concurrency primitives (goroutines, channels)
- Strong typing
- Robust standard library

### Why NOT Cobra?

- Standard library flag package is sufficient
- External dependencies kept minimal
- Simple two-command tool doesn't need a complex CLI framework
