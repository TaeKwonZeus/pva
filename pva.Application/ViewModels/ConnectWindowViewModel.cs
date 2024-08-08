using System;
using System.Linq;
using System.Threading.Tasks;
using Avalonia.Data.Converters;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Grpc.Core;

namespace pva.Application.ViewModels;

public partial class ConnectWindowViewModel : ViewModelBase
{
    [ObservableProperty] private string _address = "";

    [ObservableProperty] private string _message = "";

    [ObservableProperty] private string _password = "";

    private int? _port = 5101;

    [ObservableProperty] private string _username = "";

    public ConnectWindowViewModel(string message = "")
    {
        Message = message;
    }

    public static IMultiValueConverter FormConverter { get; } =
        new FuncMultiValueConverter<string, bool>(values => values.All(v => !string.IsNullOrWhiteSpace(v)));

    public string Port
    {
        get => _port?.ToString() ?? "";
        set
        {
            if (string.IsNullOrWhiteSpace(value))
                SetProperty(ref _port, null);
            else if (int.TryParse(value, out var val))
                SetProperty(ref _port, val);
        }
    }

    public bool Remember { get; set; }

    // [RelayCommand(CanExecute = nameof(CanConnect))]
    // private async Task Connect(object _)
    // {
    //     try
    //     {
    //         var grpcService = new GrpcService(Address, _port);
    //         if (!await grpcService.PingAsync())
    //             throw new RpcException(Status.DefaultSuccess);
    //
    //         if (Remember)
    //         {
    //             App.Config.ServerAddr = Address;
    //             App.Config.Port = _port;
    //         }
    //
    //         App.WindowManager.StartMain(this, grpcService);
    //     }
    //     catch (UriFormatException)
    //     {
    //         Message = "Invalid URL";
    //         Address = "";
    //         Port = "5101";
    //     }
    //     catch (RpcException)
    //     {
    //         Message = "Failed to connect to server";
    //         Address = "";
    //         Port = "5101";
    //     }
    // }
    //
    // private bool CanConnect(object _)
    // {
    //     return !string.IsNullOrWhiteSpace(Address);
    // }

    [RelayCommand]
    private async Task Login()
    {
        try
        {
            var grpcService = new GrpcService(Address, _port);
            if (!await grpcService.PingAsync())
                throw new RpcException(Status.DefaultSuccess);

            if (await grpcService.LoginAsync(Username, Password))
            {
                if (Remember)
                {
                    App.Config.Address = Address;
                    App.Config.Port = _port;
                    App.Config.Username = Username;
                    App.Config.Password = Password;
                    App.Config.Update();
                }

                App.WindowManager.StartMain(this, grpcService);
            }
            else
            {
                Message = "Failed to log in";
                Username = "";
                Password = "";
            }
        }
        catch (Exception)
        {
            Message = "Failed to connect to server";
            Address = "";
            Port = "5101";
        }
    }

    [RelayCommand]
    private async Task Register()
    {
    }
}