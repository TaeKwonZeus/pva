using System.Data;
using System.Diagnostics;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Security.Principal;
using System.Text;
using Dapper;
using Microsoft.Data.Sqlite;
using pva.Server.Services;

string? dataSource = null;

// Verify app is running as sudo/admin
if (OperatingSystem.IsLinux() || OperatingSystem.IsMacOS())
{
    [DllImport("libc")]
    static extern int geteuid();

    // Root check
    if (geteuid() != 0)
    {
        // Console.WriteLine("The server needs to be run as root");
        // Environment.Exit(-1);
        dataSource = "Data Source=db.sqlite";
    }
    else
    {
        dataSource = "Data Source=/var/lib/pva-server/db.sqlite";
        Directory.CreateDirectory("/var/lib/pva-server");

        Process.Start("chown", "root /var/lib/pva-server");
        Process.Start("chmod", "700 /var/lib/pva-server");
    }
}
else if (OperatingSystem.IsWindows())
{
    // Admin check
    using var identity = WindowsIdentity.GetCurrent();
    var principal = new WindowsPrincipal(identity);
    if (!principal.IsInRole(WindowsBuiltInRole.Administrator))
    {
        Console.WriteLine("The server needs to be run as root");
        Environment.Exit(-1);
    }

    // TODO create admin-only DB file for Windows
    // dataSource = "/";
}
else
{
    Console.WriteLine("The server can only be run on Windows, MacOS and Linux");
    Environment.Exit(-1);
}

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddGrpc();
builder.Services.AddTransient<IDbConnection>(_ =>
{
    var conn = new SqliteConnection(dataSource!);
    conn.Open();
    return conn;
});

var app = builder.Build();

using (var conn = app.Services.GetService<IDbConnection>()!)
{
    // Startup SQL query if database file does not exist
    var assembly = Assembly.GetExecutingAssembly();
    using var stream = assembly.GetManifestResourceStream($"{assembly.GetName().Name}.Resources.startup.sql")!;
    using var streamReader = new StreamReader(stream, Encoding.UTF8);
    streamReader.ReadToEnd();
    conn.Execute(streamReader.ReadToEnd(), commandType: CommandType.Text);
}

// Configure the HTTP request pipeline.
app.MapGrpcService<MainService>();
app.MapGet("/",
    () =>
        "Communication with gRPC endpoints must be made through a gRPC client. To learn how to create a client, visit: https://go.microsoft.com/fwlink/?linkid=2086909");

app.Run();