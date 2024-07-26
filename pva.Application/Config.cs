using System;

namespace pva.Application;

public class Config
{
    private int? _port;
    private string? _serverAddr;

    public string? ServerAddr
    {
        get => _serverAddr;
        set
        {
            _serverAddr = value;
            UpdateEvent.Invoke();
        }
    }

    public int? Port
    {
        get => _port;
        set
        {
            _port = value;
            UpdateEvent.Invoke();
        }
    }

    public event Action UpdateEvent = () => { };
}