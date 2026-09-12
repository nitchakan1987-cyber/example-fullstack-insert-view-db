# Person Management — Fullstack Example

ระบบจัดการข้อมูลบุคคลตามแบบทดสอบ IT 01 ประกอบด้วยหน้ารายการข้อมูล การเพิ่มข้อมูลผ่าน Modal และการดูรายละเอียดแบบอ่านอย่างเดียว

## เทคโนโลยีและการแบ่งหน้าที่

| ส่วน | เทคโนโลยี | หน้าที่ |
| --- | --- | --- |
| Frontend | Angular 21, TypeScript, CSS | แสดงตาราง ฟอร์มเพิ่มข้อมูล และรายละเอียด |
| Write API | C#, ASP.NET Core 10, EF Core 10 | ตรวจสอบข้อมูลและบันทึกลงฐานข้อมูล |
| Read API | Go 1.26, database/sql, go-mssqldb | อ่านรายการและรายละเอียดบุคคล |
| Database | SQL Server 2019 | เก็บข้อมูลบุคคล |
| Development environment | Docker Compose | รัน Backend และ SQL Server |

แยกงานอ่านและงานบันทึกเพื่อให้ทั้ง C# และ Go มีหน้าที่ในระบบตามข้อกำหนด โดยใช้ฐานข้อมูลเดียวกัน ฝั่ง C# ดูแลโครงสร้างตาราง ส่วน Go อ่านข้อมูลจากตารางดังกล่าว

## ฟังก์ชัน

- แสดง Id, ชื่อ-สกุล, ที่อยู่, วันเกิด, อายุ และปุ่ม View
- กด Add เพื่อเปิด Modal กรอกชื่อ นามสกุล วันเกิด และที่อยู่
- แสดงอายุตามสูตร **ปีปัจจุบัน − ปีเกิด** ตามโจทย์ ไม่หักปีกรณียังไม่ถึงวันเกิด
- บันทึกผ่าน C# API แล้วโหลดรายการใหม่ผ่าน Go API
- เรียงรายการด้วย `Id ASC` เพื่อให้ข้อมูลใหม่อยู่ท้ายตาราง
- กดยกเลิกเพื่อปิด Modal โดยไม่บันทึกข้อมูล
- กด View เพื่ออ่านรายละเอียดตาม ID จาก Go API ทุกช่องเป็น read-only
- แสดงสถานะกำลังโหลด กำลังบันทึก และข้อความเมื่อเกิดข้อผิดพลาด

## โครงสร้างโปรเจกต์

| ตำแหน่ง | รายละเอียด |
| --- | --- |
| `frontend/` | Angular project |
| `frontend/src/app/features/persons/` | Component หน้าข้อมูลบุคคล |
| `backend/Person.WriteApi/` | C# project |
| `backend/Person.WriteApi/Contracts/` | Request DTO และ validation |
| `backend/Person.WriteApi/Controllers/` | HTTP endpoints |
| `backend/Person.WriteApi/Data/` | EF Core DbContext |
| `backend/Person.WriteApi/Models/` | Entity ของฐานข้อมูล |
| `backend/person-read-api/` | Go project, main.go, go.mod และ go.sum |
| `scripts/docker-mssql.yaml` | SQL Server service และ data volume |
| `scripts/docker-dotnet.yaml` | C# API และ .NET SDK service |
| `scripts/docker-go.yaml` | Go API และ Go SDK service |
| `scripts/.env` | ค่าการเชื่อมต่อบนเครื่อง ไม่เก็บใน Git |
| `README.md` | วิธีตั้งค่าและรันระบบ |

## สิ่งที่ต้องมี

- Docker Desktop ที่เปิดใช้งาน Linux containers และ Docker Compose v2
- Node.js ที่รองรับ Angular 21 เช่น Node 22 ตั้งแต่ 22.12.0 ในสาย 22 หรือ Node 24.x พร้อม npm
- พอร์ต `4200`, `5001`, `5002` และ `1433` ว่าง

ไม่จำเป็นต้องติดตั้ง .NET SDK หรือ Go บนเครื่อง เพราะเรียกผ่าน SDK containers ส่วน Angular รันด้วย Node.js บนเครื่อง

คำสั่งด้านล่างใช้ PowerShell และให้รันจากโฟลเดอร์หลัก `example-fullstack-assessment` เว้นแต่ระบุไว้ต่างหาก

## การตั้งค่าฐานข้อมูล

สร้างไฟล์ `scripts/.env` โดยใส่ค่าดังนี้ และแทน `YOUR_PASSWORD` ด้วยรหัสผ่านจริงบนเครื่อง:

```dotenv
DB_CONNECTION='Server=sqlserver,1433;Database=PersonAssessment;User Id=sa;Password=YOUR_PASSWORD;Encrypt=True;TrustServerCertificate=True'
GO_DB_HOST=sqlserver:1433
GO_DB_NAME=PersonAssessment
GO_DB_USER=sa
GO_DB_PASSWORD='YOUR_PASSWORD'
```

สำหรับ SQL Server ใหม่ ให้ตั้ง `MSSQL_SA_PASSWORD` ใน `scripts/docker-mssql.yaml` เป็นรหัสผ่านเดียวกันและเป็นไปตาม password policy ของ SQL Server สำหรับ SQL Server ที่มีข้อมูลใน Volume อยู่แล้ว ให้ใช้รหัสผ่านปัจจุบันของบัญชี `sa` การแก้ environment variable อย่างเดียวไม่ได้เปลี่ยนรหัสผ่านในฐานข้อมูลเดิม

ไฟล์ Compose ของ `write-api` และ `read-api` ต้องมี:

```yaml
env_file:
  - .env
```

เส้นทางนี้อ้างอิงจากไฟล์ Compose ใน `scripts/` ส่วน SQL Server ใช้ค่าจากไฟล์ Compose ของตนเอง ไม่ได้อ่าน `GO_DB_PASSWORD` อัตโนมัติ

เพิ่ม `.env` ใน `.gitignore` ที่โฟลเดอร์หลัก และอย่า commit รหัสผ่านจริงในไฟล์ Compose หากต้องการใช้ตัวแปรแทนรหัสผ่านใน Compose ให้ปรับไฟล์และวิธีโหลด environment ให้ตรงกันก่อนส่งงาน

`TrustServerCertificate=True` และบัญชี `sa` ใช้สำหรับสภาพแวดล้อมพัฒนานี้

## วิธีรันระบบ

### 1. เปิด SQL Server

```powershell
docker compose -f scripts/docker-mssql.yaml up -d sqlserver
docker compose -f scripts/docker-mssql.yaml logs --tail 50 sqlserver
```

รอข้อความว่า SQL Server พร้อมรับการเชื่อมต่อก่อนเปิด Backend โดยเฉพาะครั้งแรกที่สร้าง Container

### 2. เปิด C# API

```powershell
docker compose -f scripts/docker-dotnet.yaml up -d write-api
docker compose -f scripts/docker-dotnet.yaml logs -f write-api
```

รอ `Database initialization completed.` และ `Now listening on: http://0.0.0.0:8080` แล้วกด `Ctrl + C` เพื่อออกจากการดู Log โดย Container ยังทำงานต่อ

ใน Development ฝั่ง C# ใช้ `EnsureCreatedAsync()` สร้างฐานข้อมูล `PersonAssessment` และตาราง `Persons` หากยังไม่มี จึงต้องเปิด C# สำเร็จก่อนเปิด Go ครั้งแรก

### 3. เปิด Go API

```powershell
docker compose -f scripts/docker-go.yaml up -d read-api
docker compose -f scripts/docker-go.yaml logs -f read-api
```

รอ `SQL Server connected` และ `Go API listening on :8080` แล้วกด `Ctrl + C` เพื่อออกจากการดู Log

### 4. เปิด Angular

```powershell
cd frontend
npm ci
npm start
```

เปิด [http://localhost:4200](http://localhost:4200) และปล่อย Terminal นี้ทำงานอยู่ `npm ci` ใช้ติดตั้ง dependencies จาก lockfile หลัง Clone หรือเมื่อ dependencies เปลี่ยน ไม่ต้องรันใหม่ทุกครั้ง

เมื่อ Clone โปรเจกต์ที่มีโค้ดแล้ว ไม่ต้องรัน `dotnet new`, `go mod init` หรือ `ng new` ซ้ำ

### พอร์ตและ Network

| ส่วน | URL จากเครื่องผู้ใช้ | ที่อยู่ภายใน Docker |
| --- | --- | --- |
| Angular | http://localhost:4200 | รันบนเครื่อง |
| C# API | http://localhost:5001 | write-api:8080 |
| Go API | http://localhost:5002 | read-api:8080 |
| SQL Server | localhost,1433 | sqlserver:1433 |

Compose ทั้งสามไฟล์อยู่ใน `scripts/` และต้องใช้ Compose project เดียวกัน ในการตั้งค่าปัจจุบัน Network คือ `scripts_default` จึงเรียกฐานข้อมูลด้วยชื่อ `sqlserver` ได้ ไม่ใช้ `localhost` จากภายใน Backend container

หากกำหนด `-p` หรือ `COMPOSE_PROJECT_NAME` เอง ต้องใช้ค่าเดียวกันทุกคำสั่ง สำหรับระบบเดิมไม่ควรเปลี่ยนชื่อ project โดยไม่ได้วางแผนย้าย Volume เพราะชื่อ project มีผลต่อชื่อ Network และ Volume

CORS ของ Backend อนุญาต `http://localhost:4200` จึงควรเปิด Angular ด้วย URL นี้

## API

| Method | Endpoint | Backend | ผลลัพธ์ |
| --- | --- | --- | --- |
| POST | http://localhost:5001/api/persons | C# | เพิ่มข้อมูลและตอบ 201 |
| GET | http://localhost:5002/api/persons | Go | รายการทั้งหมด เรียง Id ASC |
| GET | http://localhost:5002/api/persons/1 | Go | รายละเอียด ID 1 หรือ 404 หากไม่พบ |
| GET | http://localhost:5002/health | Go | สถานะ HTTP service |

`/health` ยืนยันว่า Go HTTP service ตอบสนอง ไม่ได้ตรวจการเชื่อมต่อฐานข้อมูลใหม่ทุกครั้ง หากต้องการตรวจเส้นทางอ่านฐานข้อมูลให้เรียก `/api/persons`

### ตัวอย่างเพิ่มข้อมูล

```json
{
  "firstName": "สมหญิง",
  "lastName": "ใจดี",
  "birthDate": "1995-10-30",
  "address": "กรุงเทพมหานคร"
}
```

ทดลองผ่าน PowerShell (จะเพิ่มข้อมูลจริงหนึ่งรายการต่อการเรียก):

```powershell
$body = @{
    firstName = "สมหญิง"
    lastName = "ใจดี"
    birthDate = "1995-10-30"
    address = "กรุงเทพมหานคร"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "http://localhost:5001/api/persons" `
    -Method Post `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

ผลตอบกลับประกอบด้วย `id`, `firstName`, `lastName`, `birthDate`, `age` และ `address` โดยอายุขึ้นกับปีที่เรียก API และ ID ขึ้นกับข้อมูลในฐานข้อมูล

### Validation

- ชื่อและนามสกุลต้องไม่ว่าง ความยาวไม่เกินช่องละ 100 ตัวอักษร
- ที่อยู่ต้องไม่ว่าง ความยาวไม่เกิน 1,000 ตัวอักษร
- ต้องระบุวันเกิดที่ถูกต้องและไม่เป็นอนาคต
- ข้อมูลไม่ผ่าน validation ตอบ HTTP 400
- Go ตอบ HTTP 400 เมื่อ ID ไม่ถูกต้อง และ HTTP 404 เมื่อไม่พบรายการ
- ไม่มีเงื่อนไขอายุขั้นต่ำตามโจทย์

## โครงสร้างฐานข้อมูล

Database: `PersonAssessment` — Table: `dbo.Persons`

| Column | SQL Server type | รายละเอียด |
| --- | --- | --- |
| Id | int IDENTITY PRIMARY KEY | รหัสที่ฐานข้อมูลสร้างให้ |
| FirstName | nvarchar(100) NOT NULL | ชื่อ |
| LastName | nvarchar(100) NOT NULL | นามสกุล |
| BirthDate | date NOT NULL | วันเกิด |
| Address | nvarchar(1000) NOT NULL | ที่อยู่ |

ไม่จัดเก็บอายุในตาราง แต่คำนวณเมื่อแสดงผล Backend ใช้ปีปัจจุบันตามเวลา UTC+7 ส่วน Preview ใน Angular ใช้เวลาของ Browser

ข้อมูล SQL Server เก็บใน named volume `mssql_data` ตาม Compose โดยชื่อจริงอาจมี prefix ของ project เช่น `scripts_mssql_data` การหยุด Container ไม่ทำให้ข้อมูลใน Volume หาย

## การพัฒนาและหยุดระบบ

หลังแก้โค้ด Backend ให้ Restart เพราะตั้งค่าเป็น `dotnet run` และ `go run` ไม่ได้เปิด watch mode:

```powershell
docker compose -f scripts/docker-dotnet.yaml restart write-api
docker compose -f scripts/docker-go.yaml restart read-api
```

หลังแก้ `.env` หรือ Compose ให้สร้าง Container ของ Service นั้นใหม่เพื่อโหลดค่า:

```powershell
docker compose -f scripts/docker-dotnet.yaml up -d --force-recreate write-api
docker compose -f scripts/docker-go.yaml up -d --force-recreate read-api
```

หยุด Backend ก่อน แล้วหยุด SQL Server:

```powershell
docker compose -f scripts/docker-go.yaml stop read-api
docker compose -f scripts/docker-dotnet.yaml stop write-api
docker compose -f scripts/docker-mssql.yaml stop sqlserver
```

หยุด Angular ด้วย `Ctrl + C` ใน Terminal ที่รัน `npm start`

## รายการตรวจสอบการทำงาน

รายการนี้เป็นขั้นตอนตรวจรับ ไม่ใช่ผลการรัน automated tests

| กรณี | ผลที่คาดหวัง |
| --- | --- |
| เปิดหน้าเว็บ | โหลดรายการจาก Go API |
| Add และบันทึกข้อมูลครบ | บันทึกผ่าน C# แล้วรายการใหม่อยู่ท้ายตาราง |
| Refresh หน้าเว็บ | ข้อมูลที่บันทึกยังอยู่ |
| กดยกเลิก Add | ไม่เพิ่มข้อมูล เปิดใหม่แล้วฟอร์มว่าง |
| เลือกวันเกิด | อายุเป็นปีปัจจุบันลบปีเกิด |
| กรอกชื่อเป็นช่องว่าง | ไม่สามารถบันทึกได้ |
| ระบุวันเกิดในอนาคต | ไม่สามารถบันทึกได้ |
| เปิด View สลับรายการ | ข้อมูลตรงกับ ID ที่เลือก |
| ลองแก้ไขช่องใน View | แก้ไขไม่ได้ มีปุ่มปิด |
| ปิด Go แล้วโหลดตาราง | แสดงข้อผิดพลาดและปุ่มลองใหม่ |
| ปิด C# แล้วกดบันทึก | แสดงข้อผิดพลาดและเก็บฟอร์มไว้ |

## ข้อจำกัดของเวอร์ชันนี้

- เป็นระบบสำหรับแบบทดสอบและการรันในเครื่อง ยังไม่มี Authentication, Pagination, แก้ไข หรือลบข้อมูล เพราะอยู่นอกขอบเขตที่เลือก
- ใช้ `EnsureCreatedAsync()` แทน EF Migrations จึงไม่อัปเดตโครงสร้างตารางเดิมเมื่อ Model เปลี่ยน หากพัฒนาต่อควรวางแผนใช้ Migrations
- URL ของ API ถูกกำหนดใน Angular สำหรับ localhost หากเปลี่ยนเครื่องหรือพอร์ตต้องปรับ URL และ CORS
- SDK images และ bind mounts ใช้เพื่อพัฒนา หากนำไป deploy ควรจัดทำ production build และ runtime images แยก
- README อธิบายระบบตามโครงสร้างนี้ ควรตรวจไฟล์จริงและรันตามขั้นตอนจาก checkout ใหม่ก่อนส่งมอบ

## แก้ปัญหาที่พบบ่อย

| อาการ | วิธีตรวจสอบ |
| --- | --- |
| หาไฟล์ Compose ไม่พบ | รันคำสั่งจากโฟลเดอร์หลักของโปรเจกต์ |
| `additional properties ... not allowed` | ตรวจ indent YAML ให้ทุก Service อยู่ใต้ `services:` |
| `dotnet` หรือ `go` ไม่พบบน Windows | ใช้ SDK service ผ่าน Docker Compose ตามที่โปรเจกต์กำหนด |
| `Found orphan containers` | เกิดได้เมื่อแยก Compose หลายไฟล์ใน project เดียวกัน อย่าใช้ `--remove-orphans` เพราะอาจลบ Service อื่นที่ยังใช้งานอยู่ |
| Go ต่อฐานข้อมูลไม่ได้ | ตรวจ SQL Server พร้อมใช้งาน รหัสผ่าน Network และ C# สร้างฐานข้อมูลแล้ว |
| API เปิดผ่าน Browser ได้ แต่ Angular เรียกไม่ได้ | ตรวจ CORS และเปิด Angular ที่ `http://localhost:4200` |
| แก้ `.env` แล้วค่าไม่เปลี่ยน | ใช้ `up -d --force-recreate` กับ Service ที่เกี่ยวข้อง |
| Template แจ้งจะเขียนทับไฟล์ | มีโปรเจกต์แล้ว ไม่ต้องสร้างใหม่หรือใช้ `--force` |

อย่าใช้ `docker compose down -v` หากต้องการเก็บข้อมูลเดิม เพราะตัวเลือก `-v` ลบ named volumes ที่ Compose จัดการ
