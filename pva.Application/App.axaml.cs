using System;
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
            var viewModel = new ConnectWindowViewModel(config);
            var window = new ConnectWindow
            {
                DataContext = viewModel
            };
            viewModel.Connected += (config1, grpcService) => OnConnected(config1, grpcService, window);
            viewModel.FailedToConnect += OnFailedToConnect;
            
            window.Show();
        }

        base.OnFrameworkInitializationCompleted();
    }

    private void OnConnected(Config config, GrpcService grpcService, ConnectWindow connectWindow)
    {
        connectWindow.Close();

        var window = new MainWindow()
        {
            DataContext = new MainWindowViewModel(config, grpcService)
        };
        if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
        {
            desktop.MainWindow = window;
        }
    }

    private void OnFailedToConnect(Exception e)
    {
        
    }
}