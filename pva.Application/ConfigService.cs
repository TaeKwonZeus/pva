using System;
using System.IO;
using System.Text.Json;

namespace pva.Application;

public class ConfigService
{
    private readonly string _addr;

    private readonly JsonSerializerOptions _prettyPrint = new()
    {
        WriteIndented = true
    };

    public ConfigService(string addr)
    {
        _addr = addr;
        if (File.Exists(_addr))
        {
            try
            {
                var config = JsonSerializer.Deserialize<Config>(File.ReadAllText(_addr))!;
                Config = config;
            }
            catch (Exception)
            {
                Config = new Config();
                File.WriteAllText(_addr, JsonSerializer.Serialize(Config, _prettyPrint));
            }
        }
        else
        {
            Config = new Config();
            File.WriteAllText(_addr, JsonSerializer.Serialize(Config, _prettyPrint));
        }

        Config.UpdateEvent += UpdateConfig;
    }

    public Config Config { get; }

    private void UpdateConfig()
    {
        File.WriteAllText(_addr, JsonSerializer.Serialize(Config, _prettyPrint));
    }
}