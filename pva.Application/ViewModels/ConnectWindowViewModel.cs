using System;
using System.Linq;
using System.Threading.Tasks;
using Avalonia.Data.Converters;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Grpc.Core;
using pva.Common;

namespace pva.Application.ViewModels;

public partial class ConnectWindowViewModel : ViewModelBase
{
    [ObservableProperty] private string _address = App.Config.Address ?? "";

    [ObservableProperty] private string _message = "";

    [ObservableProperty] private string _password = App.Config.Password ?? "";

    private int? _port = App.Config.Port ?? 5101;

    [ObservableProperty] private string _username = App.Config.Username ?? "";

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
        GrpcService grpcService;

        try
        {
            grpcService = new GrpcService(Address, _port);
            if (!await grpcService.PingAsync())
                throw new RpcException(Status.DefaultSuccess);
        }
        catch (Exception)
        {
            Message = "Failed to connect to server";
            Address = "";
            Port = "5101";
            return;
        }

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

    [RelayCommand]
    private async Task Register()
    {
        GrpcService grpcService;

        try
        {
            grpcService = new GrpcService(Address, _port);
            if (!await grpcService.PingAsync())
                throw new RpcException(Status.DefaultSuccess);
        }
        catch (Exception)
        {
            Message = "Failed to connect to server";
            Address = "";
            Port = "5101";
            return;
        }

        switch (await grpcService.RegisterAsync(Username, Password))
        {
            case RegisterStatus.RegisterOk:
                Message = "Registered successfully; you may log in";
                break;
            case RegisterStatus.RegisterUsernameExists:
                Message = "Account with given username already exists";
                Username = "";
                Password = "";
                break;
            case RegisterStatus.RegisterMissingCredentials:
                Console.WriteLine("Empty credentials; something went wrong");
                break;
            case null:
                Message = "Failed to connect to server";
                Address = "";
                Port = "5101";
                break;
            default:
                throw new ArgumentOutOfRangeException();
        }
    }
}