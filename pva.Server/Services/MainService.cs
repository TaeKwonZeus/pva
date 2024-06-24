using Grpc.Core;
using pva.Grpc;

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
        _logger.LogInformation("Responding to {Name}", request.Name);
        return Task.FromResult(new PingResponse
        {
            Message = "Pong!"
        });
    }
}