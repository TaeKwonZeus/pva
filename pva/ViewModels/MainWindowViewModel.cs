using System;
using System.Windows.Input;
using ReactiveUI;

namespace pva.GUI.ViewModels;

public class MainWindowViewModel : ViewModelBase
{
    public ICommand ClickCommand { get; }

    public MainWindowViewModel()
    {
        ClickCommand = ReactiveCommand.Create(OnClick);
    }

    private void OnClick()
    {
        Console.WriteLine("PISKA");
    }
}