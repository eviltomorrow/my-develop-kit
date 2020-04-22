package com.transfer.client;

import java.io.FileInputStream;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.Callable;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

import com.google.protobuf.BoolValue;
import com.google.protobuf.ByteString;
import com.google.protobuf.StringValue;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import pb.Transfer.CheckPoint;
import pb.Transfer.FileChannel;
import pb.Transfer.FileChannel.UploadStrategy;
import pb.Transfer.FileInfo;
import pb.Transfer.FilePart;
import pb.UploadFileGrpc;
import pb.UploadFileGrpc.UploadFileBlockingStub;
import pb.UploadFileGrpc.UploadFileStub;

public class UploadClient {
	private final ManagedChannel channel;

	public UploadClient(String host, int port) {
		channel = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build();
	}

	public FileInfo getFileInfo(String path) throws UploadFileException {
		StringValue request = StringValue.newBuilder().setValue(path).build();
		UploadFileBlockingStub bs = UploadFileGrpc.newBlockingStub(channel);
		FileInfo info = null;
		try {
			info = bs.getFileInfo(request);
		} catch (StatusRuntimeException e) {
			e.printStackTrace();
			throw new UploadFileException(e, "Get remote file info failure");
		}
		return info;
	}

	public List<FilePart> setCheckPoint(String path, FileInfo info) throws InterruptedException, UploadFileException {
		CheckPoint cp = CheckPoint.newBuilder().setPath(path).setFileInfo(info).build();

		final Map<String, Throwable> errors = new HashMap<String, Throwable>();
		final CountDownLatch latch = new CountDownLatch(1);
		final List<FilePart> partList = new ArrayList<FilePart>();

		UploadFileStub stub = UploadFileGrpc.newStub(channel);
		StreamObserver<FilePart> rp = new StreamObserver<FilePart>() {

			@Override
			public void onNext(FilePart value) {
				partList.add(value);
			}

			@Override
			public void onError(Throwable t) {
				errors.put("error", t);
				latch.countDown();
			}

			@Override
			public void onCompleted() {
				latch.countDown();
			}
		};

		stub.setCheckPoint(cp, rp);
		latch.await();

		Throwable t = errors.get("error");
		if (t != null) {
			throw new UploadFileException(t, "set checkpoint failure");
		}
		return partList;
	}

	public void uploadFile(String localPath, String checkpointPath, int taskNum, FileInfo info, List<FilePart> parts)
			throws InterruptedException, UploadFileException {
		ArrayList<Future<PartResult>> futures = new ArrayList<Future<PartResult>>();
		List<PartResult> taskResults = new ArrayList<PartResult>();

		final Map<String, Throwable> errors = new HashMap<String, Throwable>();
		final CountDownLatch finishLatch = new CountDownLatch(1);
		ExecutorService service = Executors.newFixedThreadPool(taskNum);
		UploadFileStub stub = UploadFileGrpc.newStub(channel);
		StreamObserver<BoolValue> rp = new StreamObserver<BoolValue>() {

			@Override
			public void onNext(BoolValue value) {
				System.out.println("Upload file: " + value.getValue());
			}

			@Override
			public void onError(Throwable t) {
				errors.put("error", t);
				t.printStackTrace();
				finishLatch.countDown();
			}

			@Override
			public void onCompleted() {
				System.out.println("Upload file complete");
				finishLatch.countDown();
			}
		};

		StreamObserver<FileChannel> rq = stub.uploadFile(rp);
		DataStream stream = new DataStream(rq);

		for (int i = 0; i < parts.size(); i++) {
			FilePart part = parts.get(i);
			if (part.getIsCompleted()) {
				PartResult tr = new PartResult(i, part.getOffset(), part.getSize());
				taskResults.add(tr);
			} else {
				futures.add(service.submit(new Task(i, localPath, checkpointPath, info, parts.get(i), stream)));
			}
		}

		finishLatch.await();
		service.shutdownNow();

		rq.onCompleted();

		Throwable t = errors.get("error");
		if (t != null) {
			throw new UploadFileException(t, "Upload file part failure");
		}

		for (Future<PartResult> future : futures) {
			try {
				PartResult tr = future.get();
				taskResults.add(tr);
			} catch (ExecutionException e) {
				throw new UploadFileException(t, "Execute upload file task failure");
			}
		}

		Collections.sort(taskResults, new Comparator<PartResult>() {
			@Override
			public int compare(PartResult p1, PartResult p2) {
				return p1.getNumber() - p2.getNumber();
			}
		});

		for (PartResult partResult : taskResults) {
			if (partResult.isFailed()) {
				throw new UploadFileException(partResult.getException());
			}
		}
	}

	public void shutdown() throws InterruptedException {
		channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
	}

	static class DataStream {
		private StreamObserver<FileChannel> rq;

		public DataStream(StreamObserver<FileChannel> rq) {
			this.rq = rq;
		}

		public synchronized void write(FileChannel channel) {
			this.rq.onNext(channel);
		}
	}

	static class PartResult {
		private int number;
		private long offset;
		private long size;
		private boolean isFailed;
		private Exception exception;

		public PartResult(int number, long offset, long size) {
			this.number = number;
			this.offset = offset;
			this.size = size;
		}

		public int getNumber() {
			return number;
		}

		public void setNumber(int number) {
			this.number = number;
		}

		public long getOffset() {
			return offset;
		}

		public void setOffset(long offset) {
			this.offset = offset;
		}

		public long getSize() {
			return size;
		}

		public void setSize(long size) {
			this.size = size;
		}

		public boolean isFailed() {
			return isFailed;
		}

		public void setFailed(boolean isFailed) {
			this.isFailed = isFailed;
		}

		public Exception getException() {
			return exception;
		}

		public void setException(Exception exception) {
			this.exception = exception;
		}

	}

	static class Task implements Callable<PartResult> {
		private FileInfo fileInfo;
		private FilePart filePart;
		private String localPath;
		private String checkpointPath;
		private int number;
		private DataStream stream;

		public Task(int number, String localPath, String checkpointPath, FileInfo fileInfo, FilePart filePart,
				DataStream stream) {
			this.number = number;
			this.localPath = localPath;
			this.checkpointPath = checkpointPath;
			this.fileInfo = fileInfo;
			this.filePart = filePart;
			this.stream = stream;
		}

		@Override
		public PartResult call() throws Exception {
			PartResult tr = new PartResult(this.number, this.filePart.getOffset(), this.filePart.getSize());
			InputStream instream = null;
			try {
				instream = new FileInputStream(this.localPath);
				instream.skip(this.filePart.getOffset());

				byte[] buffer = new byte[(int) filePart.getSize()];
				int n = instream.read(buffer);

				FileChannel.Builder fileChannelBuilder = FileChannel.newBuilder();
				FilePart.Builder filePartBuilder = FilePart.newBuilder();
				filePartBuilder.setOffset(this.filePart.getOffset());
				filePartBuilder.setSize(this.filePart.getSize());
				filePartBuilder.setNum(this.filePart.getNum());
				filePartBuilder.setData(ByteString.copyFrom(buffer, 0, n));
				filePartBuilder.build();

				fileChannelBuilder.setStrategy(UploadStrategy.EXIST_FAILURE).setCheckpoint(this.checkpointPath)
						.setFileInfo(fileInfo).setFilePart(filePartBuilder.build());

				FileOperation.bucket.acquire((int) this.filePart.getSize());

				this.stream.write(fileChannelBuilder.build());

			} catch (Exception e) {
				tr.setFailed(true);
				tr.setException(e);
			} finally {
				if (instream != null) {
					instream.close();
				}
			}
			return tr;
		}

	}

}
