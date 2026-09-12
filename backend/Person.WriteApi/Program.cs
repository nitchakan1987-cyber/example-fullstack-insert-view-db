using Microsoft.EntityFrameworkCore;
using Person.WriteApi.Data;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddOpenApi();

var connectionString = builder.Configuration["DB_CONNECTION"]
    ?? throw new InvalidOperationException(
        "Missing DB_CONNECTION configuration.");

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlServer(
        connectionString,
        sql => sql.UseCompatibilityLevel(150)));

builder.Services.AddCors(options =>
{
    options.AddPolicy("AngularLocal", policy =>
        policy.WithOrigins("http://localhost:4200")
              .AllowAnyHeader()
              .AllowAnyMethod());
});

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();

    using var scope = app.Services.CreateScope();
    var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();

    await db.Database.EnsureCreatedAsync();

    app.Logger.LogInformation("Database initialization completed.");
}
app.UseCors("AngularLocal");
app.MapControllers();

app.Run();