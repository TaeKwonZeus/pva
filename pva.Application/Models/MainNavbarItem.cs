using System;
using Avalonia.Media;

namespace pva.Application.Models;

public class MainNavbarItem
{
    public MainNavbarItem(Type modelType, string iconName)
    {
        ModelType = modelType;
        Label = ModelType.Name.Replace("PageViewModel", "");

        Avalonia.Application.Current!.TryGetResource(iconName, null, out var res);
        Icon = (StreamGeometry)res!;
    }

    public string Label { get; }

    public Type ModelType { get; }

    public StreamGeometry Icon { get; }
}