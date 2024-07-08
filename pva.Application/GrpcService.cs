using System.Net;
using System.Threading.Tasks;
using Grpc.Core;
using Grpc.Net.Client;
using pva.Grpc;

namespace pva.Application;

public class GrpcService
{
    private readonly GrpcChannel _channel;

    public GrpcService(string addr)
    {
        _channel = GrpcChannel.ForAddress(addr);
        if (!Ping()) throw new WebException("Failed to connect to the server.");
    }

    public bool Ping()
    {
        var client = new Main.MainClient(_channel);

        var req = client.Ping(new PingRequest());
        return req != null;
    }

    public async Task<bool> PingAsync()
    {

        var client = new Main.MainClient(_channel);

        var req = await client.PingAsync(new PingRequest());
        return req != null;
    }
}