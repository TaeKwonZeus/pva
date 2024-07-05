using System.Threading.Tasks;
using Grpc.Core;
using Grpc.Net.Client;
using pva.Grpc;

namespace pva.Application;

public static class GrpcService
{
    private static GrpcChannel? _channel;

    // Returns true if connected successfully
    public static bool Connect(string addr)
    {
        var channel = GrpcChannel.ForAddress(addr);
        if (channel.State != ConnectivityState.Ready) return false;

        _channel = channel;
        return true;
    }

    public static bool Ping()
    {
        if (_channel != null && _channel.State != ConnectivityState.Ready) return false;
        var client = new Main.MainClient(_channel);

        var req = client.Ping(new PingRequest());
        return req != null;
    }

    public static async Task<bool> PingAsync()
    {
        if (_channel != null && _channel.State != ConnectivityState.Ready) return false;

        var client = new Main.MainClient(_channel);

        var req = await client.PingAsync(new PingRequest());
        return req != null;
    }
}