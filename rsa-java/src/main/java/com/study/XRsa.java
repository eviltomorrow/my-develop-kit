package com.study;

import java.io.ByteArrayOutputStream;
import java.io.FileReader;
import java.io.IOException;
import java.io.Reader;
import java.io.Writer;
import java.security.Key;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.Signature;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.HashMap;
import java.util.Map;

import javax.crypto.Cipher;

import org.apache.commons.codec.binary.Base64;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemReader;
import org.bouncycastle.util.io.pem.PemWriter;

public class XRsa {
	public static final String CHARSET = "UTF-8";
	public static final String RSA_ALGORITHM = "RSA";
	public static final String RSA_ALGORITHM_SIGN = "SHA256WithRSA";

	private RSAPublicKey publicKey;
	private RSAPrivateKey privateKey;

	public XRsa(String publicKey, String privateKey) throws Exception {
		try {
			KeyFactory keyFactory = KeyFactory.getInstance(RSA_ALGORITHM);

			// 通过X509编码的Key指令获得公钥对象
			X509EncodedKeySpec x509KeySpec = new X509EncodedKeySpec(Base64.decodeBase64(publicKey));
			this.publicKey = (RSAPublicKey) keyFactory.generatePublic(x509KeySpec);
			// 通过PKCS#8编码的Key指令获得私钥对象
			PKCS8EncodedKeySpec pkcs8KeySpec = new PKCS8EncodedKeySpec(Base64.decodeBase64(privateKey));
			this.privateKey = (RSAPrivateKey) keyFactory.generatePrivate(pkcs8KeySpec);
		} catch (Exception e) {
			throw new Exception("不支持的密钥", e);
		}
	}

	public static String loadPrivateKey(Reader privateReader) throws IOException {
		PemReader pemReader = new PemReader(privateReader);
		PemObject pemObject = pemReader.readPemObject();
		byte[] data = pemObject.getContent();
		return Base64.encodeBase64String(data);
	}

	public static String loadPublicKey(Reader publicReader) throws IOException {
		PemReader pemReader = new PemReader(publicReader);
		PemObject pemObject = pemReader.readPemObject();
		byte[] data = pemObject.getContent();
		return Base64.encodeBase64String(data);
	}

	public static Map<String, String> generateKeys(Writer publicWriter, Writer privateWriter, int keySize)
			throws IOException {
		// 为RSA算法创建一个KeyPairGenerator对象
		KeyPairGenerator kpg;
		try {
			kpg = KeyPairGenerator.getInstance(RSA_ALGORITHM);
		} catch (NoSuchAlgorithmException e) {
			throw new IllegalArgumentException("No such algorithm-->[" + RSA_ALGORITHM + "]");
		}

		// 初始化KeyPairGenerator对象,不要被initialize()源码表面上欺骗,其实这里声明的size是生效的
		kpg.initialize(keySize);
		// 生成密匙对
		KeyPair keyPair = kpg.generateKeyPair();
		// 得到公钥
		Key publicKey = keyPair.getPublic();
		String publicKeyStr = Base64.encodeBase64URLSafeString(publicKey.getEncoded());
		// 得到私钥
		Key privateKey = keyPair.getPrivate();
		String privateKeyStr = Base64.encodeBase64URLSafeString(privateKey.getEncoded());
		Map<String, String> keyPairMap = new HashMap<String, String>();
		keyPairMap.put("publicKey", publicKeyStr);

		if (publicWriter != null) {
			// public
			PemWriter pemWriterp = new PemWriter(publicWriter);
			pemWriterp.writeObject(new PemObject("PUBLIC KEY", publicKey.getEncoded()));
			// 写入硬盘
			pemWriterp.flush();
			pemWriterp.close();
		}

		keyPairMap.put("privateKey", privateKeyStr);

		if (privateWriter != null) {
			// private
			PemWriter pemWriterpri = new PemWriter(privateWriter);
			pemWriterpri.writeObject(new PemObject("PRIVATE KEY", privateKey.getEncoded()));
			// 写入硬盘
			pemWriterpri.flush();
			pemWriterpri.close();
		}

		return keyPairMap;
	}

	public String publicEncrypt(String data) throws Exception {
		try {
			Cipher cipher = Cipher.getInstance(RSA_ALGORITHM);
			cipher.init(Cipher.ENCRYPT_MODE, publicKey);
			return Base64.encodeBase64URLSafeString(rsaSplitCodec(cipher, Cipher.ENCRYPT_MODE, data.getBytes(CHARSET),
					publicKey.getModulus().bitLength()));
		} catch (Exception e) {
			throw new Exception("加密字符串[" + data + "]时遇到异常", e);
		}
	}

	public String privateDecrypt(String data) throws Exception {
		try {
			Cipher cipher = Cipher.getInstance(RSA_ALGORITHM);
			cipher.init(Cipher.DECRYPT_MODE, privateKey);
			return new String(rsaSplitCodec(cipher, Cipher.DECRYPT_MODE, Base64.decodeBase64(data),
					publicKey.getModulus().bitLength()), CHARSET);
		} catch (Exception e) {
			throw new Exception("解密字符串[" + data + "]时遇到异常", e);
		}
	}

	public String privateEncrypt(String data) throws Exception {
		try {
			Cipher cipher = Cipher.getInstance(RSA_ALGORITHM);
			cipher.init(Cipher.ENCRYPT_MODE, privateKey);
			return Base64.encodeBase64URLSafeString(rsaSplitCodec(cipher, Cipher.ENCRYPT_MODE, data.getBytes(CHARSET),
					publicKey.getModulus().bitLength()));
		} catch (Exception e) {
			throw new Exception("加密字符串[" + data + "]时遇到异常", e);
		}
	}

	public String publicDecrypt(String data) throws Exception {
		try {
			Cipher cipher = Cipher.getInstance(RSA_ALGORITHM);
			cipher.init(Cipher.DECRYPT_MODE, publicKey);
			return new String(rsaSplitCodec(cipher, Cipher.DECRYPT_MODE, Base64.decodeBase64(data),
					publicKey.getModulus().bitLength()), CHARSET);
		} catch (Exception e) {
			throw new Exception("解密字符串[" + data + "]时遇到异常", e);
		}
	}

	public String sign(String data) throws Exception {
		try {
			// sign
			Signature signature = Signature.getInstance(RSA_ALGORITHM_SIGN);
			signature.initSign(privateKey);
			signature.update(data.getBytes(CHARSET));
			return Base64.encodeBase64URLSafeString(signature.sign());
		} catch (Exception e) {
			throw new Exception("签名字符串[" + data + "]时遇到异常", e);
		}
	}

	public boolean verify(String data, String sign) throws Exception {
		try {
			Signature signature = Signature.getInstance(RSA_ALGORITHM_SIGN);
			signature.initVerify(publicKey);
			signature.update(data.getBytes(CHARSET));
			return signature.verify(Base64.decodeBase64(sign));
		} catch (Exception e) {
			throw new Exception("验签字符串[" + data + "]时遇到异常", e);
		}
	}

	private static byte[] rsaSplitCodec(Cipher cipher, int opmode, byte[] datas, int keySize) throws Exception {
		int maxBlock = 0;
		if (opmode == Cipher.DECRYPT_MODE) {
			maxBlock = keySize / 8;
		} else {
			maxBlock = keySize / 8 - 11;
		}
		ByteArrayOutputStream out = new ByteArrayOutputStream();
		int offSet = 0;
		byte[] buff;
		int i = 0;
		byte[] resultDatas;
		try {
			while (datas.length > offSet) {
				if (datas.length - offSet > maxBlock) {
					buff = cipher.doFinal(datas, offSet, maxBlock);
				} else {
					buff = cipher.doFinal(datas, offSet, datas.length - offSet);
				}
				out.write(buff, 0, buff.length);
				i++;
				offSet = i * maxBlock;
			}
			resultDatas = out.toByteArray();
			out.close();
		} catch (Exception e) {
			throw new Exception("加解密阀值为[" + maxBlock + "]的数据时发生异常", e);
		}

		return resultDatas;

	}

	public static void main(String[] args) throws Exception {
		String publicKeyPath = "/home/shepard/workspace-agent/key/id_rsa_public.key";
		String privateKeyPath = "/home/shepard/workspace-agent/key/pkcs8_id_rsa_private.key";
//		Writer publicWriter = new FileWriter(publicKeyPath);
//		Writer privateWriter = new FileWriter(privateKeyPath);
//		Map<String, String> keys = XRsa.generateKeys(publicWriter, privateWriter, 1024);

//		String publicKey = keys.get("publicKey");
//		String privateKey = keys.get("privateKey");
//		System.out.println("私：" + privateKey);
//		System.out.println("公：" + publicKey);

		Reader publicReader = new FileReader(publicKeyPath);
		String publicKey = loadPrivateKey(publicReader);
		Reader privateReader = new FileReader(privateKeyPath);
		String privateKey = loadPrivateKey(privateReader);
		XRsa rsa = new XRsa(publicKey, privateKey);
		System.out.println();
		/***************************************************************************************************/
		String data = "How are you, 哈哈";
		System.out.println("加密前：" + data);

		String encrypted = rsa.publicEncrypt(data);
		System.out.println("加密后：" + encrypted);

		String decrypted = rsa.privateDecrypt(encrypted);
		System.out.println("解密后：" + decrypted);
		/***************************************************************************************************/
//		System.out.println();
//		Reader privateReader = new FileReader(privateKeyPath);
//		String privateKeyContent = loadPrivateKey(privateReader);
//		System.out.println("私 content: " + privateKeyContent);
//
//		Reader publicReader = new FileReader(publicKeyPath);
//		String publicKeyContent = loadPrivateKey(publicReader);
//		System.out.println("公 content: " + publicKeyContent);
	}
}
