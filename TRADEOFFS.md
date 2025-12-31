# Explaining the Pros and Cons of the Design Choices of this Project

## Architecture Choices

Hexagonal architecture approach for modularity benefits (client, processor, database, converter)

### Why?

- While this project implements Gemini client in `internal/client/gemini`, it's easy to swap AI providers later.
  I reckon I could implement an interface to make this even more flexible when the need arises. This approach absolves the
  core logic from being tampered with.
- Testing is easier as a result of the clear boundaries.
- Simple, explicit dependencies

### Alternatives

- Single package approach with everything in package main.
