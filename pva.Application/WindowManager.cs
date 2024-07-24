using System.Collections.Generic;
using Avalonia.Controls;
using Avalonia.Controls.ApplicationLifetimes;
using pva.Application.ViewModels;
using pva.Application.Views;

namespace pva.Application;

public class WindowManager
{
    private readonly IDictionary<ViewModelBase, Window>
        _windows = new Dictionary<ViewModelBase, Window>();

    public WindowManager(Config config)
    {
        Config = config;
    }

    private Config Config { get; }

    public void StartMain(ConnectWindowViewModel? connectWindowViewModel, GrpcService grpcService)
    {
        var viewModel = new MainWindowViewModel(Config, grpcService);
        var window = new MainWindow
        {
            DataContext = viewModel
        };
        _windows[viewModel] = window;
        if (Avalonia.Application.Current!.ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            desktop.MainWindow = window;

        if (connectWindowViewModel != null)
            _windows[connectWindowViewModel].Close();
    }

    public void StartConnect()
    {
        var viewModel = new ConnectWindowViewModel(Config, this);
        var view = new ConnectWindow
        {
            DataContext = viewModel
        };
        _windows[viewModel] = view;
        view.Show();
    }
}