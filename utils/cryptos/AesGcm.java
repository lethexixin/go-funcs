import javax.crypto.Cipher;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Base64;
import java.util.Random;
import java.util.zip.Deflater;
import java.util.zip.Inflater;

public class AesGcm {

    private static final int NONCE_LENGTH_BYTE = 12;
    private static final String ENCRYPT_ALGO = "AES/GCM/NoPadding";

    public static void main(String[] args) {
        String data = "{\"name\" : \"xin\"}";
        String key = "b6c1cd0fe6e55f22fb483096822b5d1c";
        String encryptData = gcmEncrypt(data, key.getBytes());
        // 加密结果会不断变化,但最后解密出来的结果都是一样的
        System.out.println("加密后数据: " + encryptData);
        String decryptData = gcmDecrypt(encryptData, key.getBytes());
        System.out.println("解密后数据: " + decryptData);
    }

    private static String gcmEncrypt(String data, byte[] aesKey) {
        if (data.length() > 0 && aesKey.length > 0) {
            byte[] compressByte = compress(data.getBytes()); // 压缩
            if (compressByte != null) {
                String randomNonce = createRandomNonce(NONCE_LENGTH_BYTE);
                byte[] nonceByte = randomNonce.getBytes();
                byte[] dataByte = aesGcmHandle(aesKey, compressByte, nonceByte, Cipher.ENCRYPT_MODE);
                if (dataByte != null) {
                    byte[] mergerByte = byteMerger(nonceByte, dataByte);
                    return new String(Base64.getEncoder().encode(mergerByte)).trim();
                }
            }
        }
        return "";
    }

    private static byte[] byteMerger(byte[] nonceByte, byte[] dataByte) {
        byte[] byteMerger = new byte[nonceByte.length + dataByte.length];
        System.arraycopy(nonceByte, 0, byteMerger, 0, nonceByte.length);
        System.arraycopy(dataByte, 0, byteMerger, nonceByte.length, dataByte.length);
        return byteMerger;
    }

    public static String gcmDecrypt(String srcData, byte[] aesSceneCardSecretKey) {
        if (!srcData.equals("")) {
            byte[] decode = Base64.getMimeDecoder().decode(srcData.getBytes());
            if (decode != null) {
                byte[] nonceByte = getNonceByte(decode);
                byte[] dataByte = getDataByte(decode);
                byte[] bytes = aesGcmHandle(aesSceneCardSecretKey, dataByte, nonceByte, Cipher.DECRYPT_MODE);
                if (bytes != null) {
                    byte[] decompress = decompress(bytes);
                    if (decompress != null) {
                        return new String(decompress).trim();
                    }
                }
            }
        }
        return "";
    }

    // 截取原数据
    private static byte[] getDataByte(byte[] decode) {
        byte[] nonceByte = new byte[decode.length - NONCE_LENGTH_BYTE];
        System.arraycopy(decode, NONCE_LENGTH_BYTE, nonceByte, 0, decode.length - NONCE_LENGTH_BYTE);
        return nonceByte;
    }

    // 截取nonce
    private static byte[] getNonceByte(byte[] decode) {
        byte[] nonceByte = new byte[NONCE_LENGTH_BYTE];
        System.arraycopy(decode, 0, nonceByte, 0, NONCE_LENGTH_BYTE);
        return nonceByte;
    }

    private static byte[] aesGcmHandle(byte[] aesSceneCardSecretKey, byte[] decode, byte[] nonce, int mode) {
        if (aesSceneCardSecretKey.length > 0 && decode.length > 0 && nonce.length > 0) {
            try {
                // AES
                SecretKeySpec secretKeySpec = new SecretKeySpec(aesSceneCardSecretKey, "AES");
                // 获取 AES 密码器(AES/GCM/NoPadding)
                Cipher cipher = Cipher.getInstance(ENCRYPT_ALGO);
                GCMParameterSpec zeroNonce = new GCMParameterSpec(128, nonce);
                cipher.init(mode, secretKeySpec, zeroNonce);
                return cipher.doFinal(decode);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        return null;
    }

    public static byte[] compress(byte[] data) {
        byte[] output = new byte[0];
        if (data.length > 0) {
            Deflater deflater = new Deflater();
            deflater.setLevel(9);
            deflater.setInput(data);
            deflater.finish();
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            try {
                byte[] buf = new byte[1024 * 8];
                while (!deflater.finished()) {
                    int byteCount = deflater.deflate(buf);
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
            deflater.end();
        }
        return output;
    }

    public static byte[] decompress(byte[] data) {
        byte[] output = new byte[0];
        if (data.length > 0) { // data.length > 0一定要加, 否则解压长度为0的字节数组会死循环
            Inflater inflater = new Inflater(false);
            inflater.reset();
            inflater.setInput(data);
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream(data.length);
            try {
                byte[] buf = new byte[1024 * 8];
                while (!inflater.finished()) {
                    int i = inflater.inflate(buf);
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
            inflater.end();
        }
        return output;
    }

    // 根据指定长度生成数字的随机数
    private static String createRandomNonce(int length) {
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
