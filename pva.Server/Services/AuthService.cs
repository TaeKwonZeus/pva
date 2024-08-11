using System.Data;
using System.Security.Cryptography;
using System.Text;
using Dapper;
using Grpc.Core;
using Microsoft.IdentityModel.JsonWebTokens;
using Microsoft.IdentityModel.Tokens;
using pva.Common;

namespace pva.Server.Services;

public class AuthService : Auth.AuthBase
{
    private readonly SigningCredentials _credentials;
    private readonly IDbConnection _db;
    private readonly ILogger<MainService> _logger;

    public AuthService(ILogger<MainService> logger, IDbConnection db,
        SigningCredentials credentials)
    {
        _logger = logger;
        _db = db;
        _credentials = credentials;
    }

    public override async Task<RegisterResponse> Register(RegisterRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Register: {}", context.Peer);

        if (string.IsNullOrWhiteSpace(request.Username) || string.IsNullOrWhiteSpace(request.Password))
            return new RegisterResponse { Status = RegisterStatus.RegisterMissingCredentials };

        if (await _db.QuerySingleOrDefaultAsync<int?>("SELECT 1 FROM users WHERE username = @username",
                new { username = request.Username }) != null)
            return new RegisterResponse { Status = RegisterStatus.RegisterUsernameExists };

        var (encryptionKey, salt) = EncryptionUtil.PasswordToKey(request.Password);
        var (publicKey, privateKey) = EncryptionUtil.GenerateKeypair();
        var encryptedPrivateKey =
            EncryptionUtil.AesEncrypt(privateKey, encryptionKey, Encoding.UTF8.GetBytes(request.Username));

        await _db.ExecuteAsync(
            """
            INSERT INTO users (username, public_key, salt, encrypted_private_key)
            VALUES(@username, @publicKey, @salt, @encryptedPrivateKey);
            """,
            new
            {
                username = request.Username,
                publicKey,
                salt = Convert.ToBase64String(salt),
                encryptedPrivateKey
            });
        return new RegisterResponse { Status = RegisterStatus.RegisterOk };
        // TODO finish
    }

    public override async Task<LoginResponse> Login(LoginRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Login: {}", context.Peer);

        if (string.IsNullOrWhiteSpace(request.Username) || string.IsNullOrWhiteSpace(request.Password))
            return new LoginResponse { Status = LoginStatus.LoginFailed };

        var userData = await _db.QuerySingleOrDefaultAsync(
            "SELECT id, salt, encrypted_private_key FROM users WHERE username = @username",
            new { username = request.Username });

        if (userData == null) return new LoginResponse { Status = LoginStatus.LoginFailed };

        long id = userData.id;
        string salt = userData.salt;
        string encryptedPrivateKey = userData.encrypted_private_key;

        try
        {
            EncryptionUtil.AesDecrypt(encryptedPrivateKey, EncryptionUtil.PasswordToKey(request.Password, salt),
                Encoding.UTF8.GetBytes(request.Username));
        }
        catch (CryptographicException)
        {
            return new LoginResponse { Status = LoginStatus.LoginFailed };
        }

        var token = new JsonWebTokenHandler().CreateToken(new SecurityTokenDescriptor
        {
            Claims = new Dictionary<string, object>
            {
                { JwtRegisteredClaimNames.Iss, context.Host },
                { JwtRegisteredClaimNames.Aud, context.Peer },
                { JwtRegisteredClaimNames.Sub, id }
            },
            Expires = null,
            SigningCredentials = _credentials
        })!;

        return new LoginResponse { Status = LoginStatus.LoginOk, Token = token };
    }
}