using System;
using System.Collections.ObjectModel;
using pva.Application.Models;

namespace pva.Application.ViewModels;

public class MainWindowViewModel : ViewModelBase
{
    private readonly Config _config;
    private readonly GrpcService _grpcService;
    private MainNavbarItem _selectedItem;

    public MainWindowViewModel(Config config, GrpcService grpcService)
    {
        _config = config;
        _grpcService = grpcService;
        _selectedItem = NavbarItems[0];
    }

    public MainNavbarItem SelectedItem
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

    public ObservableCollection<MainNavbarItem> NavbarItems { get; } =
    [
        new MainNavbarItem(typeof(PasswordsPageViewModel), "KeyRegular"),
        new MainNavbarItem(typeof(PasswordsPageViewModel), "KeyRegular")
        // TODO add more items
    ];
}