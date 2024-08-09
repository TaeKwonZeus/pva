using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using pva.Common;

namespace pva.Application;

public class GrpcService
{
    private readonly Auth.AuthClient _authClient;
    private readonly Main.MainClient _mainClient;

    private string? _token;

    public GrpcService(string addr, int? port)
    {
        var uri = new UriBuilder
        {
            Scheme = "http",
            Host = addr
        };
        if (port != null) uri.Port = port.Value;
        var channel = GrpcChannel.ForAddress(uri.Uri);
        _mainClient = new Main.MainClient(channel);
        _authClient = new Auth.AuthClient(channel);
    }

    public async Task<bool> PingAsync()
    {
        var req = await _mainClient.PingAsync(new PingRequest { Name = "Ping" });
        return req != null;
    }

    public async Task<bool> LoginAsync(string username, string password)
    {
        try
        {
            var res = await _authClient.LoginAsync(new LoginRequest
            {
                Username = username, Password = password
            });

            if (res.Status == LoginStatus.LoginFailed) return false;

            _token = res.Token;
            return true;
        }
        catch
        {
            return false;
        }
    }

    // Null is registration failed due to server error
    public async Task<RegisterStatus?> RegisterAsync(string username, string password)
    {
        try
        {
            var res = await _authClient.RegisterAsync(new RegisterRequest
            {
                Username = username,
                Password = password
            });

            return res.Status;
        }
        catch
        {
            return null;
        }
    }
}