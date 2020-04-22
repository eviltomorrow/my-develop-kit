package com.transfer.limit;

import com.google.common.util.concurrent.RateLimiter;

public class RateLimiterExample {
	public static void testAcquire() {
		RateLimiter limiter = RateLimiter.create(1);

		while(true) {
			limiter.acquire(3);
			System.out.println("acquire");
		}
	}
	
	public static void main(String[] args) {
		testAcquire();
	}
}
