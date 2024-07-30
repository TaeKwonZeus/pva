using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using pva.Common;

namespace pva.Application;

public class GrpcService
{
    private readonly GrpcChannel _channel;

    public GrpcService(string addr, int? port)
    {
        var uri = new UriBuilder
        {
            Scheme = "http",
            Host = addr
        };
        if (port != null) uri.Port = port.Value;
        _channel = GrpcChannel.ForAddress(uri.Uri);
    }

    public bool Ping()
    {
        var client = new Main.MainClient(_channel);

        var req = client.Ping(new PingRequest { Name = "Ping" }, deadline: DateTime.UtcNow.AddMilliseconds(120));
        return req != null;
    }

    public async Task<bool> PingAsync()
    {
        var client = new Main.MainClient(_channel);

        var req = await client.PingAsync(new PingRequest { Name = "Ping" });
        return req != null;
    }
}