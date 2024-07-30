using System.Security.Cryptography;
using System.Text;

namespace pva.Common;

public static class EncryptionUtil
{
    private const int SaltSize = 16;
    private const int HashSize = 32;
    private const int Iterations = 100000;

    private const int IvSize = 16;

    public static string CreateHash(string password)
    {
        var salt = RandomNumberGenerator.GetBytes(SaltSize);

        using var pbkdf2 = new Rfc2898DeriveBytes(password, salt, Iterations, HashAlgorithmName.SHA256);
        return Convert.ToBase64String(pbkdf2.GetBytes(HashSize));
    }

    public static (string publicKey, string privateKey) GenerateKeypair()
    {
        using var rsa = RSA.Create();
        return (rsa.ExportRSAPublicKeyPem(), rsa.ExportRSAPrivateKeyPem());
    }

    public static string EncryptString(string input, string key)
    {
        var inputBytes = Encoding.UTF8.GetBytes(input);
        var keyBytes = SHA256.HashData(Encoding.UTF8.GetBytes(key));

        var iv = RandomNumberGenerator.GetBytes(IvSize);

        using var aes = Aes.Create();
        aes.Key = keyBytes;
        aes.IV = iv;
        aes.Mode = CipherMode.CBC;
        aes.Padding = PaddingMode.PKCS7;

        using var encryptor = aes.CreateEncryptor(aes.Key, aes.IV);
        using var ms = new MemoryStream();
        using var cs = new CryptoStream(ms, encryptor, CryptoStreamMode.Write);
        using var sw = new StreamWriter(cs);
        sw.Write(inputBytes);

        var encryptedBytes = ms.ToArray();

        var result = new byte[IvSize + encryptedBytes.Length];
        Buffer.BlockCopy(iv, 0, result, 0, IvSize);
        Buffer.BlockCopy(encryptedBytes, 0, result, IvSize, encryptedBytes.Length);

        return Convert.ToBase64String(result);
    }

    public static string DecryptString(string input, string key)
    {
        var inputBytes = Convert.FromBase64String(input);
        var keyBytes = SHA256.HashData(Encoding.UTF8.GetBytes(key));

        var iv = new byte[IvSize];
        var encryptedBytes = new byte[inputBytes.Length - IvSize];
        Buffer.BlockCopy(inputBytes, 0, iv, 0, IvSize);
        Buffer.BlockCopy(inputBytes, IvSize, encryptedBytes, 0, encryptedBytes.Length);

        using var aes = Aes.Create();
        aes.Key = keyBytes;
        aes.IV = iv;
        aes.Mode = CipherMode.CBC;
        aes.Padding = PaddingMode.PKCS7;

        using var decryptor = aes.CreateDecryptor(aes.Key, aes.IV);
        using var ms = new MemoryStream();
        using var cs = new CryptoStream(ms, decryptor, CryptoStreamMode.Write);
        using var sw = new StreamWriter(cs);
        sw.Write(encryptedBytes);

        return Encoding.UTF8.GetString(ms.ToArray());
    }
}