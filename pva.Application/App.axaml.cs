using Avalonia.Markup.Xaml;

namespace pva.Application;

public class App : Avalonia.Application
{
    private Config _config = null!;
    private WindowManager _windowManager = null!;

    private static App CurrentApp => (Current as App)!;

    public static Config Config => CurrentApp._config;

    public static WindowManager WindowManager => CurrentApp._windowManager;

    public override void Initialize()
    {
        AvaloniaXamlLoader.Load(this);
    }

    public override void OnFrameworkInitializationCompleted()
    {
        _config = new ConfigService("appsettings.json").Config;

        _windowManager = new WindowManager();

        _windowManager.StartConnect("test");

        base.OnFrameworkInitializationCompleted();
    }
}