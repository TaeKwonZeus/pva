using System;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Grpc.Core;

namespace pva.Application.ViewModels;

public partial class ConnectWindowViewModel : ViewModelBase
{
    private readonly Config _config;
    private readonly WindowManager _windowManager;

    [ObservableProperty] private string _message;

    public ConnectWindowViewModel(Config config, WindowManager windowManager, string message = "")
    {
        _config = config;
        _windowManager = windowManager;
        _message = message;
    }

    public string Url { get; set; } = "";

    public bool Remember { get; set; } = true;


    [RelayCommand]
    public void Connect()
    {
        if (Remember) _config.ServerAddr = Url;

        try
        {
            var grpcService = new GrpcService(Url);

            _windowManager.StartMain(this, grpcService);
        }
        catch (RpcException e)
        {
            Console.WriteLine(e);
        }
    }
}