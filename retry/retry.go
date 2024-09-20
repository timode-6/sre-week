package main

import (
    "math/rand"
    "time"
)

func retry(args Args, packets chan Packet, rps int) {
    response := make(chan error)
    packet := Packet{args: args, response: response}
    
    baseTimeout := time.Second / time.Duration(rps)
    maxTimeout := 4 * time.Second 
    
    attempt := 0
    for {
        if send(packet, packets) {
            break
        }
        
        attempt++
        timeout := exponentialBackoff(baseTimeout, attempt, maxTimeout)
        time.Sleep(timeout)
    }
}

func exponentialBackoff(baseTimeout time.Duration, attempt int, maxTimeout time.Duration) time.Duration {
    timeout := baseTimeout * time.Duration(1<<uint(attempt)) 
    
    jitter := time.Duration(rand.Float64() * float64(timeout) * 0.25)
    timeout += jitter
    
    if timeout > maxTimeout {
        timeout = maxTimeout
    }
    
    return timeout
}