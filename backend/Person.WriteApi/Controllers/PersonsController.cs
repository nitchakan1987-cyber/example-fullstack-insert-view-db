using Microsoft.AspNetCore.Mvc;
using Person.WriteApi.Contracts;
using Person.WriteApi.Data;
using PersonEntity = global::Person.WriteApi.Models.Person;

namespace Person.WriteApi.Controllers;

[ApiController]
[Route("api/persons")]
public class PersonsController : ControllerBase
{
    private readonly AppDbContext _db;

    public PersonsController(AppDbContext db)
    {
        _db = db;
    }

    [HttpPost]
    public async Task<IActionResult> Create(
        [FromBody] CreatePersonRequest request,
        CancellationToken cancellationToken)
    {

        var today = DateOnly.FromDateTime(
            DateTimeOffset.UtcNow.ToOffset(
                TimeSpan.FromHours(7)).DateTime);

        if (request.BirthDate is null)
        {
            ModelState.AddModelError(
                nameof(request.BirthDate),
                "กรุณาระบุวันเกิด");

            return ValidationProblem(ModelState);
        }

        if (request.BirthDate.Value > today)
        {
            ModelState.AddModelError(
                nameof(request.BirthDate),
                "วันเกิดต้องไม่เป็นวันที่ในอนาคต");

            return ValidationProblem(ModelState);
        }

        var person = new PersonEntity
        {
            FirstName = request.FirstName.Trim(),
            LastName = request.LastName.Trim(),
            BirthDate = request.BirthDate.Value,
            Address = request.Address.Trim()
        };

        _db.Persons.Add(person);
        await _db.SaveChangesAsync(cancellationToken);

        return StatusCode(StatusCodes.Status201Created, new
        {
            person.Id,
            person.FirstName,
            person.LastName,
            person.BirthDate,
            age = today.Year - person.BirthDate.Year,
            person.Address
        });
    }
}