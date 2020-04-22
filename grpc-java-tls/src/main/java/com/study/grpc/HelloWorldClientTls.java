package com.study.grpc;

import java.io.File;
import java.util.concurrent.TimeUnit;

import javax.net.ssl.SSLException;

import com.study.proto.GreeterGrpc;
import com.study.proto.HelloReply;
import com.study.proto.HelloRequest;

import io.grpc.ManagedChannel;
import io.grpc.StatusRuntimeException;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NegotiationType;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;

public class HelloWorldClientTls {
	private final ManagedChannel channel;
	private final GreeterGrpc.GreeterBlockingStub blockingStub;

	private static SslContext buildSslContext(String trustCertCollectionFilePath, String clientCertChainFilePath,
			String clientPrivateKeyFilePath) throws SSLException {
		SslContextBuilder builder = GrpcSslContexts.forClient();
		if (trustCertCollectionFilePath != null) {
			builder.trustManager(new File(trustCertCollectionFilePath));
		}
		if (clientCertChainFilePath != null && clientPrivateKeyFilePath != null) {
			builder.keyManager(new File(clientCertChainFilePath), new File(clientPrivateKeyFilePath));
		}
		return builder.build();
	}

	public HelloWorldClientTls(String host, int port, SslContext sslContext) throws SSLException {

		this(NettyChannelBuilder.forAddress(host, port).negotiationType(NegotiationType.TLS).sslContext(sslContext)
				.build());
	}

	HelloWorldClientTls(ManagedChannel channel) {
		this.channel = channel;
		blockingStub = GreeterGrpc.newBlockingStub(channel);
	}

	public void shutdown() throws InterruptedException {
		channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
	}

	public void greet(String name) {
		HelloRequest request = HelloRequest.newBuilder().setName(name).build();
		HelloReply response;
		try {
			response = blockingStub.sayHello(request);
		} catch (StatusRuntimeException e) {
			e.printStackTrace();
			return;
		}
		System.out.println(response.getMessage());
	}

	public static void main(String[] args) throws Exception {

		String path = HelloWorldClientTls.class.getResource("/").toString();
		path = path.replace("file:", "");

		String host = "localhost";
		int port = 8080;
		String caCertFile = path + "certs/ca.crt";
		String certFile = path + "certs/client.crt";
		String keyFile = path + "certs/client.pem";
		HelloWorldClientTls client = null;
		try {
			client = new HelloWorldClientTls(host, port, buildSslContext(caCertFile, certFile, keyFile));
			System.out.println(client == null);
			client.greet("tom");
		} finally {
			client.shutdown();
		}
	}
}
