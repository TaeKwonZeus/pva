using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using pva.Common;

namespace pva.Application;

public class GrpcService
{
    private readonly Auth.AuthClient _authClient;
    private readonly GrpcChannel _channel;
    private readonly Main.MainClient _mainClient;

    private string _token;

    public GrpcService(string addr, int? port)
    {
        var uri = new UriBuilder
        {
            Scheme = "http",
            Host = addr
        };
        if (port != null) uri.Port = port.Value;
        _channel = GrpcChannel.ForAddress(uri.Uri);
        _mainClient = new Main.MainClient(_channel);
        _authClient = new Auth.AuthClient(_channel);
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
        catch (Exception)
        {
            return false;
        }
    }
}