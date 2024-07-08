using System;
using System.Net;
using Avalonia.Controls.ApplicationLifetimes;
using pva.Application.Views;

namespace pva.Application.ViewModels;

public class ConnectWindowViewModel : ViewModelBase
{
    private readonly Config _config;

    public event Action<Config, GrpcService> Connected = (_, _) => { };
    public event Action<Exception> FailedToConnect = (_) => { };

    public ConnectWindowViewModel(Config config)
    {
        _config = config;
    }

    public string Url { get; set; } = "";
    public bool Remember { get; set; } = true;

    public void Connect()
    {
        if (Remember) _config.ServerAddr = Url;

        try
        {
            var grpcService = new GrpcService(Url);

            Connected.Invoke(_config, grpcService);
        }
        catch (Exception e)
        {
            Console.WriteLine(e);
            FailedToConnect.Invoke(e);
        }
    }
}