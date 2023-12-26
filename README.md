# **Cache**

## **Introduction**
Made implementation of in-memory cache and persistent storage. Users concurrently add data to the cache or retrieve it from cache. If cache has expired, it retrieves the data from the database and installs it into the cache.
## Installation
Clone this repository
  ```bash
    git clone https://github.com/NRKA/Cache.git
```
## How to Run
```bash
  go run cmd/main.go
```
