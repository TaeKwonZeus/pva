using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using pva.Application.ViewModels;
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
        if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            desktop.MainWindow = new MainWindow
            {
                DataContext = new MainWindowViewModel()
            };

        base.OnFrameworkInitializationCompleted();
    }
}