//unit tests for the tokenBucket class.

package tokenBucket

import (
  "testing"
  "time"
  "fmt"
  )
  
  //TestTokenBucketSub50ms(): the maximum update frequency for a TokenBucket is fixed at 20Hz.
  //For update frequencies higher than 20Hz, the number of tokens per update is increased so that
  //the effective update frequency (assuming one token per update) is equal to the requested frequency.
  //This test requests an update frequency of 40Hz, so we should see 2 tokens/50ms going into the bucket.
  func TestTokenBucketSub50ms(t *testing.T){
    duration, _ := time.ParseDuration("25ms")
    pause, _    := time.ParseDuration("120ms")
    capacity    := uint64(1000)
    bucket, err := NewTokenBucket(duration,capacity)
    if err != nil {
      t.Fatalf("NewTocketBucket(%v,%v) failed: %v", duration, capacity, err)
    }
    if bucket.GetToken(1000) != true {
      t.Fatalf("Failed to get 1000 tokens at once.")
    }
    time.Sleep(pause)
    
    count := bucket.AvailableTokens()
    
    if count != 4 {
      t.Fatalf("bucket.AvailableTokens() == %v, should be 4",count)
    }
    
    
  }
  
  //TestTokenBucketFillIn5(): this test drains the bucket and ensures that it fills in the correct time.
  func TestTokenBucketFillIn5(t *testing.T) {
    pause6, _    := time.ParseDuration("6s")
    duration, _ := time.ParseDuration("500ms")
    capacity    := uint64(10)
    bucket, err := NewTokenBucket(duration,capacity)
    if err != nil {
      t.Fatalf("NewTocketBucket(%v,%v) failed: %v", duration, capacity, err)
    }
    
    count := bucket.AvailableTokens()
    if count != 10 {
      t.Fatalf("bucket.AvailableTokens() == %v, should be 10",count)
    }
    time.Sleep(duration)
    count = bucket.AvailableTokens()
    if count != 10 {
      t.Fatalf("bucket.AvailableTokens() == %v, should be 10",count)
    }
    
    if bucket.GetToken(10) != true {
      t.Fatalf("Failed to get 10 tokens at once.")
    }
    
    count = bucket.AvailableTokens()
    
    if count != 0 {
      t.Fatalf("Bucket has %v tokens, should be 0", count)
    }
    time.Sleep(pause6)
    
    count = bucket.AvailableTokens()
    
    if count != 10 {
      t.Fatalf("Bucket has %v tokens, should be 10", count)
    }
    
    
  }
  
  //This test drains the bucket faster than it fills to ensure proper behaviour.
  func TestTokenBucketDrainBucketIn5(t *testing.T){
    pause, _    := time.ParseDuration("300ms")
    duration, _ := time.ParseDuration("500ms")
    capacity    := uint64(10)
    bucket, err := NewTokenBucket(duration,capacity)
    now         := time.Now()
    if err != nil {
      t.Fatalf("NewTocketBucket(%v,%v) failed: %v", duration, capacity, err)
    }
    for i := 0; i < 20; i++ {
      ok := bucket.GetToken(1)
      if ok == true {
        fmt.Printf("Got token (%v)\n", i)
      }else{
        fmt.Printf("Failed to get a token, waiting 300ms\n")
        i--
        time.Sleep(pause)
      }
    }
    
    elapsed := time.Since(now).Seconds()
    if elapsed < 5.0 || elapsed > 5.5 {
      t.Fatalf("Run time(t) for test should be 5.0 < t < 5.5, was t=%v", elapsed)
    }
    
    bucket.Shutdown()
    
  }