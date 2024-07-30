using System.Collections.Generic;
using Avalonia.Controls;
using Avalonia.Controls.ApplicationLifetimes;
using pva.Application.ViewModels;
using pva.Application.Views;

namespace pva.Application;

public class WindowManager
{
    private readonly Dictionary<ViewModelBase, Window> _windows = new();

    public void StartMain(ConnectWindowViewModel? connectWindowViewModel, GrpcService grpcService)
    {
        var viewModel = new MainWindowViewModel(grpcService);
        var window = new MainWindow
        {
            DataContext = viewModel
        };
        _windows[viewModel] = window;
        window.Show();
        if (Avalonia.Application.Current!.ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            desktop.MainWindow = window;

        if (connectWindowViewModel != null)
        {
            _windows[connectWindowViewModel].Close();
            _windows.Remove(connectWindowViewModel);
        }
    }

    public void StartConnect(string message = "")
    {
        var viewModel = new ConnectWindowViewModel(message);
        var view = new ConnectWindow
        {
            DataContext = viewModel
        };
        _windows[viewModel] = view;
        view.Show();
    }
}