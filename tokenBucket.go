// Copyright Ian Logan 2012, License terms are in the LICENSE file within the source distribution.

//Provides a token bucket useful for rate limiting an application. Upon creation of a tokenBucket, a capacity (# of tokens)
//and a refill rate are established (1 token is added every x ns until the bucket is full). When an application needs to
//consume a resource, tokens are extracted from the bucket. If the bucket is empty, extraction (should) block until a token is available.
package tokenBucket

import (
	"errors"
	"time"
)

type tokenRequest struct {
	count uint64    //How many tokens do you want?
	done  chan bool //Channel to notify on when the request has been granted.
}

type TokenBucket struct {
	req      chan tokenRequest //Channel for requests and responses - it is possible to request multiple tokens at once
	shutdown chan bool         //Send a true to shutdown the timer and kill the goroutine.
	available chan uint64
}

//returns an initialized tokenBucket, creates the goroutine that services the tokenBucket requests.
func NewTokenBucket(r time.Duration, capacity uint64) (TokenBucket, error) {
	var bucket TokenBucket
	var incr uint64
	incr = 1
	if r <= 0 {
		return bucket, errors.New("Duration must be a positive value")
	}
	//We don't want to trigger the Ticker too often, so set the minimum time slice to 50ms and adjust the increment on the bucket if necessary
	if r < 50*time.Millisecond {
		incr = uint64((50 * time.Millisecond) / r)
		r = 50 * time.Millisecond
	}
	
	count := capacity
	clock := time.NewTicker(r)
	bucket = TokenBucket{make(chan tokenRequest, 25), make(chan bool), make(chan uint64)} //our buckets start out full
	go func() {
		for {
			select {
			case <-clock.C:
				if count < capacity {
					count += incr
				}
			case request := <-bucket.req:
				if count >= request.count {
					count -= request.count
					request.done <- true
				} else {
					request.done <- false
				} //end case request := <- bucket.req
			case stop := <-bucket.shutdown:
				if stop == true {
					clock.Stop()
					return
				}
			case bucket.available <- count:
			  
			} //end select {}
		} //end for {}
	}() //end go func {}
	return bucket, nil
}

// Attempts to acquire count many tokens from the bucket, returns true/false on sucess/failure.
// XXX: it would be better if this blocked until count many tokens were available.
func (t *TokenBucket) GetToken(count uint64) bool {
	if count < 1 {
		return false
	}
	request := tokenRequest{count, make(chan bool)}
	t.req <- request
	return <-request.done
}

func (t *TokenBucket) AvailableTokens() uint64 {
  return <- t.available
}

func (t *TokenBucket) Shutdown() {
	t.shutdown <- true
}
