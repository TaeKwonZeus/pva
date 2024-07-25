using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using pva.Grpc;

namespace pva.Application;

public class GrpcService
{
    private readonly GrpcChannel _channel;

    public GrpcService(string addr)
    {
        _channel = GrpcChannel.ForAddress(addr);
    }

    public bool Ping()
    {
        var client = new Main.MainClient(_channel);

        var req = client.Ping(new PingRequest { Name = "Ping" }, deadline: DateTime.UtcNow.AddSeconds(3));
        return req != null;
    }

    public async Task<bool> PingAsync()
    {
        var client = new Main.MainClient(_channel);

        var req = await client.PingAsync(new PingRequest { Name = "Ping" });
        return req != null;
    }
}