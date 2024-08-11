using Grpc.Core;
using pva.Common;

namespace pva.Server.Services;

public class MainService : Main.MainBase
{
    private readonly ILogger<MainService> _logger;

    public MainService(ILogger<MainService> logger)
    {
        _logger = logger;
    }

    public override Task<PingResponse> Ping(PingRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Ping: {}", context.Peer);
        return Task.FromResult(new PingResponse
        {
            Message = "Pong!"
        });
    }
}