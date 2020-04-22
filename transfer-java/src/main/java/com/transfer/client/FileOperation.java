package com.transfer.client;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.security.MessageDigest;
import java.util.List;

import org.apache.commons.codec.binary.Hex;

import com.google.common.util.concurrent.RateLimiter;

import pb.Transfer.FileInfo;
import pb.Transfer.FileInfo.Builder;
import pb.Transfer.FilePart;

public class FileOperation {
	public final static RateLimiter bucket = RateLimiter.create(1024 * 1024 * 10);

	public static FileResult upload(String host, int port, String localPath, String remotePath, int taskNum)
			throws UploadFileException {
		FileResult result = new FileResult();
		UploadClient client = null;
		try {
			File localFile = new File(localPath);
			if (!localFile.exists()) {
				throw new UploadFileException("The file path [" + localPath + "] is not exist");
			}
			if (localFile.isDirectory()) {
				throw new UploadFileException("The file path [" + localPath + "] is a dierectory");
			}

			long size = localFile.length();
			long lastMod = localFile.lastModified();
			String localMD5 = getMD5(localFile);

			client = new UploadClient(host, port);

			FileInfo remoteInfo = client.getFileInfo(remotePath);
			if (remoteInfo.getPath() != "" && !remoteInfo.getIsDir()) {
				if (localMD5.equals(remoteInfo.getMd5())) {
					result.setLocalPath(localPath);
					result.setRemotePath(remotePath);
					result.setSize(size);
					result.setMd5(localMD5);
					result.setLastMod(remoteInfo.getLastMod());
					return result;
				}
			}

			Builder remoteFileInfoBuilder = FileInfo.newBuilder();
			remoteFileInfoBuilder.setPath(remotePath);
			remoteFileInfoBuilder.setIsDir(false);
			remoteFileInfoBuilder.setSize(size);
			remoteFileInfoBuilder.setMd5(localMD5);
			remoteFileInfoBuilder.setLastMod(lastMod);

			List<FilePart> parts = client.setCheckPoint(remotePath + ".cpf", remoteFileInfoBuilder.build());
			if (parts.size() == 0) {
				throw new UploadFileException("System internal error: the file parts's size is 0");
			}

			client.uploadFile(localPath, remotePath + ".cpf", taskNum, remoteFileInfoBuilder.build(), parts);

			remoteInfo = client.getFileInfo(remotePath);
			if (remoteInfo.getPath() != "" && !remoteInfo.getIsDir()) {
				if (localMD5.equals(remoteInfo.getMd5())) {
					result.setLocalPath(localPath);
					result.setRemotePath(remotePath);
					result.setSize(size);
					result.setMd5(localMD5);
					result.setLastMod(remoteInfo.getLastMod());
					return result;
				}
			}

			throw new UploadFileException("System internal error: check remote file failure");
		} catch (Exception e) {
			e.printStackTrace();
			throw new UploadFileException(e);
		} finally {
			if (client != null) {
				try {
					client.shutdown();
				} catch (InterruptedException e) {
					e.printStackTrace();
				}
			}
		}
	}

	private static String getMD5(File file) throws UploadFileException {
		FileInputStream fileInputStream = null;
		try {
			MessageDigest MD5 = MessageDigest.getInstance("MD5");
			fileInputStream = new FileInputStream(file);
			byte[] buffer = new byte[8192];
			int length;
			while ((length = fileInputStream.read(buffer)) != -1) {
				MD5.update(buffer, 0, length);
			}
			return new String(Hex.encodeHex(MD5.digest()));
		} catch (Exception e) {
			throw new UploadFileException(e, "calculate md5 for the file path [" + file.getPath() + "] failure");
		} finally {
			try {
				if (fileInputStream != null) {
					fileInputStream.close();
				}
			} catch (IOException e) {
				e.printStackTrace();
			}
		}
	}

	public static void main(String[] args) {

		try {
			upload("127.0.0.1", 8080, "/home/shepard/workspace-tmp/local/go1.13.8.linux-amd64.tar.gz",
					"/home/shepard/workspace-tmp/remote/go1.13.8.linux-amd64.tar.gz", 1);
		} catch (UploadFileException e) {
			e.printStackTrace();
		}

	}

}
