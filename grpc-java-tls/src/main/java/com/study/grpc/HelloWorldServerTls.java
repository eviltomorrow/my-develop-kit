package com.study.grpc;

import java.io.File;
import java.io.IOException;
import java.net.InetSocketAddress;

import com.study.proto.GreeterGrpc.GreeterImplBase;
import com.study.proto.HelloReply;
import com.study.proto.HelloRequest;

import io.grpc.Server;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyServerBuilder;
import io.grpc.stub.StreamObserver;
import io.netty.handler.ssl.ClientAuth;
import io.netty.handler.ssl.SslContextBuilder;
import io.netty.handler.ssl.SslProvider;

public class HelloWorldServerTls {
	private Server server;

	private final String host;
	private final int port;
	private final String certChainFilePath;
	private final String privateKeyFilePath;
	private final String trustCertCollectionFilePath;

	public HelloWorldServerTls(String host, int port, String certChainFilePath, String privateKeyFilePath,
			String trustCertCollectionFilePath) {
		this.host = host;
		this.port = port;
		this.certChainFilePath = certChainFilePath;
		this.privateKeyFilePath = privateKeyFilePath;
		this.trustCertCollectionFilePath = trustCertCollectionFilePath;
	}

	private SslContextBuilder getSslContextBuilder() {
		SslContextBuilder sslClientContextBuilder = SslContextBuilder.forServer(new File(certChainFilePath),
				new File(privateKeyFilePath));
		if (trustCertCollectionFilePath != null) {
			sslClientContextBuilder.trustManager(new File(trustCertCollectionFilePath));
			sslClientContextBuilder.clientAuth(ClientAuth.REQUIRE);
		}
		return GrpcSslContexts.configure(sslClientContextBuilder, SslProvider.OPENSSL);
	}

	private void start() throws IOException {
		server = NettyServerBuilder.forAddress(new InetSocketAddress(host, port)).addService(new GreeterImpl())
				.sslContext(getSslContextBuilder().build()).build().start();
		Runtime.getRuntime().addShutdownHook(new Thread() {
			@Override
			public void run() {
				System.err.println("*** shutting down gRPC server since JVM is shutting down");
				HelloWorldServerTls.this.stop();
				System.err.println("*** server shut down");
			}
		});
	}

	private void stop() {
		if (server != null) {
			server.shutdown();
		}
	}

	private void blockUntilShutdown() throws InterruptedException {
		if (server != null) {
			server.awaitTermination();
		}
	}

	public static void main(String[] args) throws IOException, InterruptedException {
		String path = HelloWorldServerTls.class.getResource("/").toString();
		path = path.replace("file:", "");

		String host = "localhost";
		int port = 8080;
		String certFile = path + "certs/server.crt";
		String keyFile = path + "certs/server.pem";
		String caCertFile = path + "certs/ca.crt";
		final HelloWorldServerTls server = new HelloWorldServerTls(host, port, certFile, keyFile, caCertFile);
		server.start();
		server.blockUntilShutdown();
	}

	static class GreeterImpl extends GreeterImplBase {
		public void sayHello(HelloRequest req, StreamObserver<HelloReply> responseObserver) {
			HelloReply reply = HelloReply.newBuilder().setMessage("Hello " + req.getName()).build();
			System.out.println(req.getName());
			responseObserver.onNext(reply);
			responseObserver.onCompleted();
		}
	}
}
