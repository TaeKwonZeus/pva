namespace pva.Application.ViewModels;

public class ConnectWindowViewModel : ViewModelBase
{
    private readonly Config _config;

    public ConnectWindowViewModel(Config config)
    {
        _config = config;
    }

    public string Url { get; set; } = "";
    public bool Remember { get; set; } = true;

    public void Connect()
    {
        // TODO Connect to gRPC server, launch MainWindow and set desktop.MainWindow
    }
}