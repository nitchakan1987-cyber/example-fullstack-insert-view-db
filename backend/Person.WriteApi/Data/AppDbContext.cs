using Microsoft.EntityFrameworkCore;
using PersonEntity = global::Person.WriteApi.Models.Person;

namespace Person.WriteApi.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options)
        : base(options)
    {
    }

    public DbSet<PersonEntity> Persons => Set<PersonEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        var person = modelBuilder.Entity<PersonEntity>();

        person.ToTable("Persons");
        person.HasKey(p => p.Id);

        person.Property(p => p.FirstName)
            .HasMaxLength(100)
            .IsRequired();

        person.Property(p => p.LastName)
            .HasMaxLength(100)
            .IsRequired();

        person.Property(p => p.BirthDate)
            .HasColumnType("date");

        person.Property(p => p.Address)
            .HasMaxLength(1000)
            .IsRequired();
    }
}