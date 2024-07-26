using Avalonia.Markup.Xaml;
using Grpc.Core;

namespace pva.Application;

public class App : Avalonia.Application
{
    public override void Initialize()
    {
        AvaloniaXamlLoader.Load(this);
    }

    public override void OnFrameworkInitializationCompleted()
    {
        var config = new ConfigService("appsettings.json").Config;

        var windowManager = new WindowManager(config);

        if (config.ServerAddr != null && config.Port != null)
            try
            {
                var grpcService = new GrpcService(config.ServerAddr, config.Port.Value);
                if (!grpcService.Ping())
                    throw new RpcException(Status.DefaultCancelled);
                windowManager.StartMain(null, grpcService);
            }
            catch (RpcException)
            {
                config.ServerAddr = null;
                windowManager.StartConnect("Failed to connect with saved configuration");
            }
        else
            windowManager.StartConnect("test");

        base.OnFrameworkInitializationCompleted();
    }
}