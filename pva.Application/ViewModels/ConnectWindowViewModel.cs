using System;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Grpc.Core;

namespace pva.Application.ViewModels;

public partial class ConnectWindowViewModel : ViewModelBase
{
    private readonly Config _config;
    private readonly WindowManager _windowManager;

    [ObservableProperty] private string _message;

    [ObservableProperty] private string _url = "";

    public ConnectWindowViewModel(Config config, WindowManager windowManager, string message = "")
    {
        _config = config;
        _windowManager = windowManager;
        Message = message;
    }

    public bool Remember { get; set; } = false;

    [RelayCommand(CanExecute = nameof(CanConnect))]
    private async Task Connect(object a)
    {
        try
        {
            var grpcService = new GrpcService(Url);
            if (!await grpcService.PingAsync())
                throw new RpcException(Status.DefaultCancelled);

            if (Remember) _config.ServerAddr = Url;

            _windowManager.StartMain(this, grpcService);
        }
        catch (UriFormatException e)
        {
            Message = "Invalid URL";
            Url = "";
        }
        catch (RpcException e)
        {
            Message = "Failed to connect to server";
            Url = "";
        }
    }

    private bool CanConnect(object _)
    {
        return !string.IsNullOrWhiteSpace(Url);
    }
}