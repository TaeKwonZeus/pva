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

        if (config.ServerAddr != null)
            try
            {
                var grpcService = new GrpcService(config.ServerAddr);
                windowManager.StartMain(null, grpcService);
            }
            catch (RpcException e)
            {
                config.ServerAddr = null;
                windowManager.StartConnect();
            }
        else
            windowManager.StartConnect();

        base.OnFrameworkInitializationCompleted();
    }
}