using System.ComponentModel.DataAnnotations;

namespace Person.WriteApi.Contracts;

public class CreatePersonRequest
{
    [Required]
    [StringLength(100)]
    public string FirstName { get; set; } = string.Empty;

    [Required]
    [StringLength(100)]
    public string LastName { get; set; } = string.Empty;

    [Required]
    public DateOnly? BirthDate { get; set; }

    [Required]
    [StringLength(1000)]
    public string Address { get; set; } = string.Empty;
}