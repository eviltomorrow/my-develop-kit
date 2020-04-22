package com.transfer.client;

public class UploadFileException extends Exception {
	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public UploadFileException() {
		super();
	}

	public UploadFileException(String message) {
		super(message);
	}

	public UploadFileException(Throwable t) {
		super(t);
	}

	public UploadFileException(Throwable t, String message) {
		super(message, t);
	}
}
