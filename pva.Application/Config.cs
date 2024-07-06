using System;

namespace pva.Application;

public class Config
{
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

    public event Action UpdateEvent = () => { };
}