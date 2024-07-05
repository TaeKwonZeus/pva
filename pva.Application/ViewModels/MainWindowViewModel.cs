using System;
using System.Collections.ObjectModel;
using System.Reactive.Concurrency;
using Avalonia.Media;
using ReactiveUI;

namespace pva.Application.ViewModels;

public class MainWindowViewModel : ViewModelBase
{
    private SidebarItem _selectedItem;

    public MainWindowViewModel()
    {
        _selectedItem = SidebarItems[0];
        RxApp.TaskpoolScheduler.Schedule(EstablishGrpcConnection);
    }

    public SidebarItem SelectedItem
    {
        get => _selectedItem;
        set
        {
            if (value == _selectedItem) return;
            _selectedItem = value;
            var instance = Activator.CreateInstance(SelectedItem.ModelType);
            if (instance is null) return;
            CurrentPage = (ViewModelBase)instance;
            Console.WriteLine("Tab changed");
        }
    }

    public ViewModelBase CurrentPage { get; private set; } = new PasswordsPageViewModel();

    public ObservableCollection<SidebarItem> SidebarItems { get; } =
    [
        new SidebarItem(typeof(PasswordsPageViewModel), "KeyRegular"),
        new SidebarItem(typeof(PasswordsPageViewModel), "KeyRegular")
        // TODO add more items
    ];

    private void EstablishGrpcConnection()
    {
    }
}

public class SidebarItem
{
    public SidebarItem(Type modelType, string iconName)
    {
        ModelType = modelType;
        Label = ModelType.Name.Replace("PageViewModel", "");

        Avalonia.Application.Current!.TryGetResource(iconName, null, out var res);
        Icon = (StreamGeometry)res!;
    }

    public string Label { get; }

    public Type ModelType { get; }

    public StreamGeometry Icon { get; }
}