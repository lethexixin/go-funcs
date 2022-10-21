import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Random;
import java.util.zip.Deflater;
import java.util.zip.Inflater;
import java.util.Base64; //JKD1.8

public class AesCdc {

    public static void main(String[] args) {
        String key = "b6c1cd0fe6e55f22fb483096822b5d1c";
        String sdkReq = "{\"name\" : \"xin\"}";
        String encrypt = encryptData(sdkReq, key);
        // 加密结果会不断变化,但最后解密出来的结果都是一样的
        System.out.println("加密结果:" + encrypt);
        System.out.println("解密结果:" + decrypt(encrypt, key.getBytes()));
    }

    // 解密
    public static String decrypt(String srcData, byte[] aesKey) { // base64 --> AES --> 压缩
        if (srcData.length() > 0) {
            byte[] data = Base64.getDecoder().decode(srcData.getBytes());
            byte[] iv = new byte[16];
            System.arraycopy(data, 0, iv, 0, iv.length);
            byte[] aesData = new byte[data.length - iv.length];
            System.arraycopy(data, iv.length, aesData, 0, data.length - iv.length);
            byte[] decryptAesByte = decodeAES(aesData, aesKey, iv); // AES解密
            if (decryptAesByte != null) {
                byte[] decompress = decompress(decryptAesByte);// 解压
                if (decompress != null) {
                    return new String(decompress).trim();
                }
            }
        }

        return "";
    }

    // 压缩base64内容
    public static byte[] compress(byte[] data) {
        byte[] output = null;
        if (data != null && data.length > 0) {
            Deflater def= new Deflater();
            def.setLevel(9);
            def.setInput(data);
            def.finish();

            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            try {
                byte[] buf = new byte[8 * 1024];
                while (!def.finished()) {
                    int byteCount = def.deflate(buf);
                    byteArrayOutputStream.write(buf, 0, byteCount);
                }
                output = byteArrayOutputStream.toByteArray();
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                try {
                    byteArrayOutputStream.close();
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }

            def.end();
        }
        return output;
    }

    // 解压base64内容
    public static byte[] decompress(byte[] data) {
        byte[] output = null;
        if (data != null && data.length > 0) { // data.length > 0一定要加, 否则解压长度为0的字节数组会死循环
            Inflater inf = new Inflater(false);    // no wrap header and tail
            inf.reset();
            inf.setInput(data);

            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream(data.length);
            try {
                byte[] buf = new byte[1024];
                while (!inf.finished()) {
                    int i = inf.inflate(buf);
                    byteArrayOutputStream.write(buf, 0, i);
                }
                output = byteArrayOutputStream.toByteArray();
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                try {
                    byteArrayOutputStream.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }

            inf.end();
        }

        return output;
    }

    public static String encryptData(String data, String key) {
        return encrypt(data, key.getBytes());
    }

    // 加密
    private static String encrypt(String data, byte[] aesKey) { // 压缩--> AES --> base64
        if (data.length() > 0) {
            byte[] compressByte = compress(data.getBytes()); // 压缩
            if (compressByte != null) {
                String randomIv = createRandomIv(16);
                byte[] iv = randomIv.getBytes();
                byte[] aesByte = encodeAES(compressByte, aesKey, iv); // AES加密
                if (aesByte != null) {
                    byte[] newData = byteArrayMerger(iv, aesByte);
                    if (newData != null) {
                        return new String(Base64.getEncoder().encode(newData)).trim(); // base64加密
                    }
                }
            }
        }
        return "";
    }


    private static byte[] byteArrayMerger(byte[] srcByte1, byte[] srcByte2) {
        byte[] targetByte = null;
        if (srcByte1 != null && srcByte2 != null) {
            targetByte = new byte[srcByte1.length + srcByte2.length];
            System.arraycopy(srcByte1, 0, targetByte, 0, srcByte1.length);
            System.arraycopy(srcByte2, 0, targetByte, srcByte1.length, srcByte2.length);
        }

        return targetByte;
    }

    private static byte[] generateAESCBCAlgorithm(int opMode, byte[] encData, byte[] secretKey, byte[] vector) {
        if (encData != null && secretKey != null && vector != null) {
            try {
                SecretKeySpec keySpec = new SecretKeySpec(secretKey, "AES");
                Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
                IvParameterSpec iv = new IvParameterSpec(vector);
                cipher.init(opMode, keySpec, iv);
                return cipher.doFinal(encData);
            } catch (Throwable e) {
                e.printStackTrace();
            }
        }

        return null;
    }

    public static byte[] encodeAES(byte[] encData, byte[] secretKey, byte[] vector) {
        return generateAESCBCAlgorithm(Cipher.ENCRYPT_MODE, encData, secretKey, vector);
    }

    public static byte[] decodeAES(byte[] encData, byte[] secretKey, byte[] vector) {
        return generateAESCBCAlgorithm(Cipher.DECRYPT_MODE, encData, secretKey, vector);
    }

    private static String createRandomIv(int length) {
        StringBuilder sb = new StringBuilder();
        Random random = new Random();
        int data = 0;
        for (int i = 0; i < length; i++) {
            data = random.nextInt(10); // 仅仅会生成0~9
            sb.append(data);
        }

        return sb.toString();
    }
}
