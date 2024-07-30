using System.Data;
using Dapper;
using Grpc.Core;
using pva.Common;

namespace pva.Server.Services;

public class AuthService : Auth.AuthBase
{
    private readonly IDbConnection _db;
    private readonly ILogger<MainService> _logger;

    public AuthService(ILogger<MainService> logger, IDbConnection db)
    {
        _logger = logger;
        _db = db;
    }

    public override async Task<RegisterResponse> Register(RegisterRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.Username) || string.IsNullOrWhiteSpace(request.Password))
            return new RegisterResponse { Status = RegisterStatus.RegisterMissingCredentials };

        if (await _db.QuerySingleOrDefaultAsync("SELECT 1 FROM users WHERE username = ?",
                request.Username) != null)
            return new RegisterResponse { Status = RegisterStatus.RegisterUsernameExists };

        var passwordHash = EncryptionUtil.CreateHash(request.Password);
        var (publicKey, privateKey) = EncryptionUtil.GenerateKeypair();
        var encryptedPrivateKey = EncryptionUtil.EncryptString(privateKey, request.Password);

        await _db.ExecuteAsync(
            """
            INSERT INTO users (username, password_hash, public_key, encrypted_private_key)
            VALUES(@username, @passwordHash, @publicKey, @encryptedPrivateKey);
            """,
            new
            {
                username = request.Username,
                passwordHash,
                publicKey,
                encryptedPrivateKey
            });
        return new RegisterResponse { Status = RegisterStatus.RegisterOk };
        // TODO finish
    }

    public override Task<LoginResponse> Login(LoginRequest request, ServerCallContext context)
    {
        return base.Login(request, context);
    }
}