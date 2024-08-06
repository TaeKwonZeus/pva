using System.Security.Cryptography;
using System.Text;

namespace pva.Common;

public static class EncryptionUtil
{
    private const int Pbkdf2SaltSize = 16;
    private const int Pbkdf2HashSize = 32;
    private const int Pbkdf2Iterations = 100000;
    private const int GcmTagSize = 16;
    private const int GcmNonceSize = 12;

    public static (string publicKey, string privateKey) GenerateKeypair()
    {
        using var rsa = RSA.Create();
        return (rsa.ExportRSAPublicKeyPem(), rsa.ExportRSAPrivateKeyPem());
    }

    public static (byte[] key, byte[] salt) PasswordToKey(string password)
    {
        var passwordBytes = Encoding.UTF8.GetBytes(password);
        var salt = RandomNumberGenerator.GetBytes(Pbkdf2SaltSize);
        using var pbkdf2 = new Rfc2898DeriveBytes(passwordBytes, salt, Pbkdf2Iterations, HashAlgorithmName.SHA256);

        return (pbkdf2.GetBytes(Pbkdf2HashSize), salt);
    }

    public static byte[] PasswordToKey(string password, string salt)
    {
        var passwordBytes = Encoding.UTF8.GetBytes(password);
        var saltBytes = Convert.FromBase64String(salt);
        using var pbkdf2 = new Rfc2898DeriveBytes(passwordBytes, saltBytes, Pbkdf2Iterations, HashAlgorithmName.SHA256);

        return pbkdf2.GetBytes(Pbkdf2HashSize);
    }

    public static string AesEncrypt(string plaintext, byte[] key, byte[]? aad = null)
    {
        var plaintextBytes = Encoding.UTF8.GetBytes(plaintext);

        using var aes = new AesGcm(key, GcmTagSize);

        var nonce = RandomNumberGenerator.GetBytes(GcmNonceSize);

        var ciphertextBuffer = new byte[plaintextBytes.Length + GcmTagSize + GcmNonceSize];
        var tagBuffer = new byte[GcmTagSize];

        aes.Encrypt(nonce, plaintextBytes, ciphertextBuffer, tagBuffer, aad);

        // Copy tag and nonce to end of ciphertextBuffer
        Buffer.BlockCopy(tagBuffer, 0, ciphertextBuffer, plaintextBytes.Length, GcmTagSize);
        Buffer.BlockCopy(nonce, 0, ciphertextBuffer, plaintextBytes.Length + GcmTagSize, GcmNonceSize);

        return Convert.ToBase64String(ciphertextBuffer);
    }

    public static string AesDecrypt(string ciphertext, byte[] key, byte[]? aad = null)
    {
        var ciphertextBytes = Convert.FromBase64String(ciphertext);
        var ciphertextLength = ciphertextBytes.Length - GcmTagSize - GcmNonceSize;

        var ciphertextBuffer = new byte[ciphertextLength];
        var tagBuffer = new byte[GcmTagSize];
        var nonceBuffer = new byte[GcmNonceSize];
        Buffer.BlockCopy(ciphertextBytes, 0, ciphertextBuffer, 0, ciphertextLength);
        Buffer.BlockCopy(ciphertextBytes, ciphertextLength, tagBuffer, 0, GcmTagSize);
        Buffer.BlockCopy(ciphertextBytes, ciphertextLength + GcmTagSize, nonceBuffer, 0, GcmNonceSize);

        using var aes = new AesGcm(key, GcmTagSize);

        var plaintextBuffer = new byte[ciphertextLength];
        aes.Decrypt(nonceBuffer, ciphertextBuffer, tagBuffer, plaintextBuffer, aad);

        return Encoding.UTF8.GetString(plaintextBuffer);
    }
}