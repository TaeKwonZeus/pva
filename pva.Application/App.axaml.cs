using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using pva.Application.ViewModels;
using pva.Application.Views;
using MainWindow = pva.Application.Views.MainWindow;

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

        if (config.ServerAddr != null)
        {
            var grpcService = new GrpcService(config.ServerAddr);

            if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
                desktop.MainWindow = new MainWindow
                {
                    DataContext = new MainWindowViewModel(config, grpcService)
                };
        }
        else
        {
            // TODO dialog window asking for server address
            var window = new ConnectWindow
            {
                DataContext = new ConnectWindowViewModel(config)
            };
            window.Show();
        }

        base.OnFrameworkInitializationCompleted();
    }
}