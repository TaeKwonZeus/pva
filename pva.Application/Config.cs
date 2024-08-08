using System;
using System.Text.Json.Serialization;

namespace pva.Application;

public class Config
{
    public Config(Action updateAction)
    {
        UpdateAction = updateAction;
    }

    public string? Address { get; set; }

    public int? Port { get; set; }

    public string? Username { get; set; }

    public string? Password { get; set; }

    [JsonIgnore] private Action UpdateAction { get; }

    public void Update()
    {
        UpdateAction.Invoke();
    }
}