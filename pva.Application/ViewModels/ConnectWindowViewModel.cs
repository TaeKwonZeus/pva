using System;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Grpc.Core;

namespace pva.Application.ViewModels;

public partial class ConnectWindowViewModel : ViewModelBase
{
    [ObservableProperty] private string _address;

    [ObservableProperty] private string _message;

    private int? _port;

    public ConnectWindowViewModel(string message = "")
    {
        Message = message;
        _address = App.Config.ServerAddr ?? "";
        _port = App.Config.Port ?? 5101;
    }

    public string Port
    {
        get => _port?.ToString() ?? "";
        set
        {
            if (string.IsNullOrWhiteSpace(value))
                _port = null;
            else if (int.TryParse(value, out var val))
                SetProperty(ref _port, val);
        }
    }


    public bool Remember { get; set; }

    [RelayCommand(CanExecute = nameof(CanConnect))]
    private async Task Connect(object _)
    {
        try
        {
            var grpcService = new GrpcService(Address, _port);
            if (!await grpcService.PingAsync())
                throw new RpcException(Status.DefaultSuccess);

            if (Remember)
            {
                App.Config.ServerAddr = Address;
                App.Config.Port = _port;
            }

            App.WindowManager.StartMain(this, grpcService);
        }
        catch (UriFormatException)
        {
            Message = "Invalid URL";
            Address = "";
            Port = "5101";
        }
        catch (RpcException)
        {
            Message = "Failed to connect to server";
            Address = "";
            Port = "5101";
        }
    }

    private bool CanConnect(object _)
    {
        return !string.IsNullOrWhiteSpace(Address);
    }
}