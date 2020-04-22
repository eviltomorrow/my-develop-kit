package com.transfer;

import java.util.concurrent.CountDownLatch;

public class Main {
	public static void main(String[] args) throws InterruptedException {
		String local = "/tmp/local/mysql-8.0.18-linux-glibc2.12-x86_64.tar.xz";
		String remote = "/tmp/remote/mysql-8.0.18-linux-glibc2.12-x86_64.tar.xz";
		String host = "127.0.0.1";
		int port = 8080;
		int n = 1;

//		long start = System.currentTimeMillis();
//		CountDownLatch latch = new CountDownLatch(n);
//		for (int i = 0; i < n; i++) {
//			Executor e = new Executor(local, remote + i, host, port, latch);
//			new Thread(e).start();
//		}
//		latch.await();
//		long end = System.currentTimeMillis();
//		long interval = end - start;
//		System.out.println("运行时长： " + interval / 1000 + " s");
	}
}
