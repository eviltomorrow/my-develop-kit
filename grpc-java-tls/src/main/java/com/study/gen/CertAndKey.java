package com.study.gen;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.Reader;
import java.io.Writer;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.security.Security;
import java.security.Signature;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Date;
import java.util.Locale;
import java.util.UUID;

import org.bouncycastle.asn1.ASN1EncodableVector;
import org.bouncycastle.asn1.ASN1Encoding;
import org.bouncycastle.asn1.ASN1Integer;
import org.bouncycastle.asn1.DERBitString;
import org.bouncycastle.asn1.DERNull;
import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.DERSet;
import org.bouncycastle.asn1.pkcs.CertificationRequest;
import org.bouncycastle.asn1.pkcs.CertificationRequestInfo;
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.Certificate;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.ExtensionsGenerator;
import org.bouncycastle.asn1.x509.GeneralName;
import org.bouncycastle.asn1.x509.GeneralNames;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.asn1.x509.TBSCertificate;
import org.bouncycastle.asn1.x509.Time;
import org.bouncycastle.asn1.x509.V3TBSCertificateGenerator;
import org.bouncycastle.asn1.x9.X9ObjectIdentifiers;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.bc.BcX509ExtensionUtils;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.pkcs.PKCS10CertificationRequest;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemReader;
import org.bouncycastle.util.io.pem.PemWriter;

public class CertAndKey {
	public static final String BC = BouncyCastleProvider.PROVIDER_NAME;
	public static int BITS = 2048;
	static {
		Security.addProvider(new BouncyCastleProvider());
	}

	public static KeyPair generateKey() throws NoSuchAlgorithmException, NoSuchProviderException {
		SecureRandom random = new SecureRandom();
		KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA", BC);
		keyPairGenerator.initialize(BITS, random);
		KeyPair key = keyPairGenerator.generateKeyPair();
		return key;
	}

	public static PrivateKey readPirvateKey(String keyPath) throws Exception {
		Reader reader = null;
		PemReader pemReader = null;
		try {
			reader = new FileReader(keyPath);
			pemReader = new PemReader(reader);
			PemObject pemObject = pemReader.readPemObject();
			byte[] data = pemObject.getContent();
			KeyFactory kf = KeyFactory.getInstance("RSA");
			PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(data);
			PrivateKey keyPair = kf.generatePrivate(keySpec);
			return keyPair;
		} catch (IOException | NoSuchAlgorithmException | InvalidKeySpecException e) {
			throw e;
		} finally {
			if (reader != null) {
				try {
					reader.close();
				} catch (IOException e) {
					e.printStackTrace();
				}
			}
			if (pemReader != null) {
				pemReader.close();
			}
		}

	}

	public static void writePrivateKey(String path, KeyPair key) {
		Writer writer = null;
		PemWriter pemWriterKey = null;
		try {
			writer = new FileWriter(path);
			pemWriterKey = new PemWriter(writer);
			pemWriterKey.writeObject(new PemObject("RSA PRIVATE KEY", key.getPrivate().getEncoded()));
			pemWriterKey.flush();
		} catch (IOException e) {
			e.printStackTrace();
		} finally {
			try {
				pemWriterKey.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
			if (writer != null) {
				try {
					writer.close();
				} catch (IOException e) {
					e.printStackTrace();
				}
			}
		}

	}

	public static Certificate generateCARootCertificate(X500Name name, int expirationDay, KeyPair selfKeyPair)
			throws Exception {
		SubjectPublicKeyInfo subjectPublicKeyInfo = SubjectPublicKeyInfo
				.getInstance(selfKeyPair.getPublic().getEncoded());
		BcX509ExtensionUtils extUtils = new BcX509ExtensionUtils();
		ExtensionsGenerator extensionsGenerator = new ExtensionsGenerator();
		extensionsGenerator.addExtension(Extension.basicConstraints, true, new BasicConstraints(true));
		extensionsGenerator.addExtension(Extension.keyUsage, true,
				new KeyUsage(KeyUsage.digitalSignature | KeyUsage.keyCertSign | KeyUsage.cRLSign));
		extensionsGenerator.addExtension(Extension.subjectKeyIdentifier, false,
				extUtils.createSubjectKeyIdentifier(subjectPublicKeyInfo));
		extensionsGenerator.addExtension(Extension.authorityKeyIdentifier, false,
				extUtils.createAuthorityKeyIdentifier(subjectPublicKeyInfo));
		V3TBSCertificateGenerator tbsGen = new V3TBSCertificateGenerator();
		tbsGen.setSerialNumber(new ASN1Integer(UUID.randomUUID().getMostSignificantBits() & Long.MAX_VALUE));
		tbsGen.setIssuer(name);

		Date notBefore = new Date();
		Date notAfter = new Date(System.currentTimeMillis() + 1000L * 60L * 60L * 24L * expirationDay); // xx 天

		tbsGen.setStartDate(new Time(notBefore, Locale.CHINA));
		tbsGen.setEndDate(new Time(notAfter, Locale.CHINA));
		tbsGen.setSubject(name);
		tbsGen.setSubjectPublicKeyInfo(subjectPublicKeyInfo);
		tbsGen.setExtensions(extensionsGenerator.generate());
		tbsGen.setSignature(getSignAlgo(subjectPublicKeyInfo.getAlgorithm())); // 签名算法标识等于密钥算法标识
		TBSCertificate tbs = tbsGen.generateTBSCertificate();
		Certificate certificate = assembleCert(tbs, subjectPublicKeyInfo, selfKeyPair.getPrivate());
		return certificate;
	}

	public static Certificate generateCertificate(X500Name name, int expirationDay, X509Certificate caCertificate,
			PrivateKey caPrivateKey, KeyPair selfKeyPair, String[] dns, String[] ip) throws Exception {
		PKCS10CertificationRequest csr = generateCSR(name, selfKeyPair.getPublic(), selfKeyPair.getPrivate());
		Date notBefore = new Date();
		Date notAfter = new Date(System.currentTimeMillis() + 1000L * 60L * 60L * 24L * expirationDay); // xx 天

		X509CertificateHolder issuer = new X509CertificateHolder(caCertificate.getEncoded());
		if (!verifyCSR(csr))
			throw new IllegalArgumentException("证书请求验证失败");

		BcX509ExtensionUtils extUtils = new BcX509ExtensionUtils();
		ExtensionsGenerator extensionsGenerator = new ExtensionsGenerator();
		extensionsGenerator.addExtension(Extension.basicConstraints, true, new BasicConstraints(false));
		extensionsGenerator.addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.digitalSignature));
		extensionsGenerator.addExtension(Extension.authorityKeyIdentifier, false,
				extUtils.createAuthorityKeyIdentifier(issuer)); // 授权密钥标识
		extensionsGenerator.addExtension(Extension.subjectKeyIdentifier, false,
				extUtils.createSubjectKeyIdentifier(csr.getSubjectPublicKeyInfo())); // 使用者密钥标识

		if (ip != null || dns != null) {
			int ipLength = 0;
			int dnsLength = 0;
			if (ip != null) {
				ipLength = ip.length;
			}
			if (dns != null) {
				dnsLength = dns.length;
			}
			GeneralName[] names = new GeneralName[ipLength + dnsLength];
			int i = 0;
			for (String val : ip) {
				GeneralName gn = new GeneralName(GeneralName.iPAddress, val);
				names[i] = gn;
				i++;
			}
			for (String val : dns) {
				GeneralName gn = new GeneralName(GeneralName.dNSName, val);
				names[i] = gn;
				i++;
			}

			GeneralNames subjectAltNames = new GeneralNames(names);
			extensionsGenerator.addExtension(Extension.subjectAlternativeName, false, subjectAltNames);
		}

		V3TBSCertificateGenerator tbsGen = new V3TBSCertificateGenerator();
		tbsGen.setSerialNumber(new ASN1Integer(UUID.randomUUID().getMostSignificantBits() & Long.MAX_VALUE));
		tbsGen.setIssuer(issuer.getSubject());
		tbsGen.setStartDate(new Time(notBefore, Locale.CHINA));
		tbsGen.setEndDate(new Time(notAfter, Locale.CHINA));
		tbsGen.setSubject(csr.getSubject());
		tbsGen.setSubjectPublicKeyInfo(csr.getSubjectPublicKeyInfo());
		tbsGen.setExtensions(extensionsGenerator.generate());

		tbsGen.setSignature(getSignAlgo(issuer.getSubjectPublicKeyInfo().getAlgorithm()));
		TBSCertificate tbs = tbsGen.generateTBSCertificate();
		Certificate certificate = assembleCert(tbs, issuer.getSubjectPublicKeyInfo(), caPrivateKey);
		return certificate;
	}

	public static X509Certificate readCertificate(String path) throws Exception {
		CertificateFactory cf;
		FileInputStream in = null;
		try {
			cf = CertificateFactory.getInstance("X.509");
			in = new FileInputStream(path);
			return (X509Certificate) cf.generateCertificate(in);
		} catch (CertificateException | FileNotFoundException e) {
			throw e;
		} finally {
			if (in != null) {
				in.close();
			}
		}
	}

	public static void writeCertificate(String path, Certificate certificate) {
		Writer writer = null;
		PemWriter pemWriterCert = null;
		try {
			writer = new FileWriter(path);
			pemWriterCert = new PemWriter(writer);
			pemWriterCert.writeObject(new PemObject("CERTIFICATE", certificate.getEncoded()));
			pemWriterCert.flush();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} finally {
			try {
				pemWriterCert.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
			if (writer != null) {
				try {
					writer.close();
				} catch (IOException e) {
					e.printStackTrace();
				}
			}
		}
	}

	private static PKCS10CertificationRequest generateCSR(X500Name subject, PublicKey publicKey,
			PrivateKey privateKey) {
		try {
			SubjectPublicKeyInfo subjectPublicKeyInfo = SubjectPublicKeyInfo.getInstance(publicKey.getEncoded());
			CertificationRequestInfo info = new CertificationRequestInfo(subject, subjectPublicKeyInfo, new DERSet());
			byte[] signature;
			AlgorithmIdentifier signAlgo = getSignAlgo(subjectPublicKeyInfo.getAlgorithm());
			if (signAlgo.getAlgorithm().equals(PKCSObjectIdentifiers.sha256WithRSAEncryption)) {
				signature = RSASign(info.getEncoded(ASN1Encoding.DER), privateKey);
			} else {
				throw new IllegalArgumentException("密钥算法不支持");
			}
			return new PKCS10CertificationRequest(
					new CertificationRequest(info, signAlgo, new DERBitString(signature)));
		} catch (Exception e) {
			e.printStackTrace();
			throw new IllegalArgumentException("密钥结构错误");
		}
	}

	private static AlgorithmIdentifier getSignAlgo(AlgorithmIdentifier asymAlgo) { // 根据公钥算法标识返回对应签名算法标识
		if (asymAlgo.getAlgorithm().equals(PKCSObjectIdentifiers.rsaEncryption)) {
			return new AlgorithmIdentifier(PKCSObjectIdentifiers.sha256WithRSAEncryption, DERNull.INSTANCE);
		} else
			throw new IllegalArgumentException("密钥算法不支持");
	}

	private static Certificate assembleCert(TBSCertificate tbsCertificate,
			SubjectPublicKeyInfo issuerSubjectPublicKeyInfo, PrivateKey issuerPrivateKey) throws Exception {
		byte[] signature = null;
		if (issuerPrivateKey.getAlgorithm().equalsIgnoreCase("RSA")) {
			signature = RSASign(tbsCertificate.getEncoded(), issuerPrivateKey);
		} else {
			throw new IllegalArgumentException("不支持的密钥算法");
		}

		ASN1EncodableVector v = new ASN1EncodableVector();
		v.add(tbsCertificate);
		v.add(getSignAlgo(issuerSubjectPublicKeyInfo.getAlgorithm()));
		v.add(new DERBitString(signature));
		return Certificate.getInstance(new DERSequence(v));
	}

	private static boolean verifyCSR(PKCS10CertificationRequest csr) throws Exception {
		byte[] signature = csr.getSignature();
		if (csr.getSignatureAlgorithm().getAlgorithm().equals(PKCSObjectIdentifiers.sha256WithRSAEncryption)) {
			return RSAVerifySign(csr.toASN1Structure().getCertificationRequestInfo().getEncoded(ASN1Encoding.DER),
					signature, getPublicKey(csr.getSubjectPublicKeyInfo()));
		} else {
			throw new IllegalArgumentException("不支持的签名算法");
		}
	}

	private static byte[] RSASign(byte[] inData, PrivateKey privateKey) throws Exception {
		Signature signer = Signature.getInstance("SHA256WITHRSA", BC);
		signer.initSign(privateKey);
		signer.update(inData);
		return signer.sign();
	}

	private static PublicKey getPublicKey(SubjectPublicKeyInfo subjectPublicKeyInfo) throws Exception {
		BouncyCastleProvider bouncyCastleProvider = ((BouncyCastleProvider) Security.getProvider(BC));
		bouncyCastleProvider.addKeyInfoConverter(PKCSObjectIdentifiers.rsaEncryption,
				new org.bouncycastle.jcajce.provider.asymmetric.rsa.KeyFactorySpi());
		bouncyCastleProvider.addKeyInfoConverter(X9ObjectIdentifiers.id_ecPublicKey,
				new org.bouncycastle.jcajce.provider.asymmetric.ec.KeyFactorySpi.EC());
		return BouncyCastleProvider.getPublicKey(subjectPublicKeyInfo);
	}

	private static boolean RSAVerifySign(byte[] inData, byte[] signature, PublicKey publicKey) throws Exception {
		Signature signer = Signature.getInstance("SHA256WITHRSA", BC);
		signer.initVerify(publicKey);
		signer.update(inData);
		return signer.verify(signature);
	}

	public static void main(String[] args) throws Exception {
		String subject = "C = CN, ST = BeiJing, L = BeiJing, O = Apple Inc, OU = Dev, CN = localhost";
		String caKeyPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/ca.key";
		String caCertPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/ca.crt";
		KeyPair caKeyPair = generateKey();
		writePrivateKey(caKeyPath, caKeyPair);

		Certificate certificate = generateCARootCertificate(new X500Name(subject), 365, caKeyPair);
		writeCertificate(caCertPath, certificate);

		String serverKeyPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/server.pem";
		String serverCertPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/server.crt";
		KeyPair serverKeyPair = generateKey();
		writePrivateKey(serverKeyPath, serverKeyPair);

		X509Certificate caCert = readCertificate(caCertPath);
		String[] dns = { "localhost" };
		String[] ip = { "127.0.0.1" };
		certificate = generateCertificate(new X500Name(subject), 365, caCert, caKeyPair.getPrivate(), serverKeyPair,
				dns, ip);
		writeCertificate(serverCertPath, certificate);

		String clientKeyPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/client.pem";
		String clientCertPath = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/client.crt";
		KeyPair clientKeyPair = generateKey();
		writePrivateKey(clientKeyPath, clientKeyPair);

		caCert = readCertificate(caCertPath);
		certificate = generateCertificate(new X500Name(subject), 365, caCert, caKeyPair.getPrivate(), clientKeyPair,
				null, null);
		writeCertificate(clientCertPath, certificate);
	}
}
