import { Component, inject, signal, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';


interface PersonItem {
  id: number;
  firstName: string;
  lastName: string;
  birthDate: string;
  age: number;
  address: string;
}

@Component({
  selector: 'app-persons',
  imports: [FormsModule],
  templateUrl: './persons.html',
  styleUrl: './persons.css'
})
export class Persons implements OnInit {
  private readonly http = inject(HttpClient);

  readonly selectedPerson = signal<PersonItem | null>(null);
  readonly viewLoading = signal(false);
  readonly viewError = signal('');

  readonly saving = signal(false);
  readonly saveError = signal('');
  readonly successMessage = signal('');

  readonly persons = signal<PersonItem[]>([]);
  readonly loading = signal(false);
  readonly loadError = signal('');

  firstName = '';
  lastName = '';
  birthDate = '';
  address = '';

  ngOnInit(): void {

    this.loadPersons();
  }
  loadPersons(): void {
    this.loading.set(true);
    this.loadError.set('');

    this.http.get<PersonItem[]>(
      'http://localhost:5002/api/persons'
    ).subscribe({
      next: (items) => {
        this.persons.set(items);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.loadError.set(
          'โหลดรายการไม่สำเร็จ กรุณาตรวจว่า Go API เปิดอยู่'
        );
      }
    });
  }

  formatBirthDate(value: string): string {
    const [year, month, day] = value.split('-');
    return `${day}/${month}/${year}`;
  }



  get age(): number | null {
    if (!this.birthDate) return null;
    return new Date().getFullYear() - Number(this.birthDate.slice(0, 4));
  }

  get today(): string {
    const date = new Date();
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
  }

  get canSave(): boolean {
    return (
      this.firstName.trim().length > 0 &&
      this.lastName.trim().length > 0 &&
      this.address.trim().length > 0 &&
      this.birthDate !== '' &&
      this.birthDate <= this.today
    );
  }

  openAdd(dialog: HTMLDialogElement): void {
    this.firstName = '';
    this.lastName = '';
    this.birthDate = '';
    this.address = '';
    this.saveError.set('');
    this.successMessage.set('');
    dialog.showModal();
  }

  save(dialog: HTMLDialogElement): void {
    if (!this.canSave || this.saving()) return;

    this.saving.set(true);
    this.saveError.set('');
    this.successMessage.set('');

    const body = {
      firstName: this.firstName.trim(),
      lastName: this.lastName.trim(),
      birthDate: this.birthDate,
      address: this.address.trim()
    };

    this.http.post<{ id: number }>(
      'http://localhost:5001/api/persons',
      body
    ).subscribe({
      next: (result) => {
        this.saving.set(false);
        this.successMessage.set(
          `บันทึกข้อมูลสำเร็จ รหัสรายการ ${result.id}`
        );
        dialog.close();
        this.loadPersons();
      },
      error: (error: HttpErrorResponse) => {
        this.saving.set(false);

        if (error.status === 400) {
          this.saveError.set(
            'ข้อมูลไม่ถูกต้อง กรุณาตรวจช่องกรอกและวันเกิด'
          );
        } else if (error.status === 0) {
          this.saveError.set(
            'เชื่อมต่อ API ไม่ได้ กรุณาตรวจว่า API เปิดอยู่และตั้งค่า CORS แล้ว'
          );
        } else {
          this.saveError.set(
            'บันทึกไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'
          );
        }
      }
    });
  }

  openView(id: number, dialog: HTMLDialogElement): void {
    if (this.viewLoading()) return;

    this.selectedPerson.set(null);
    this.viewError.set('');
    this.viewLoading.set(true);
    dialog.showModal();

    this.http.get<PersonItem>(
      `http://localhost:5002/api/persons/${id}`
    ).subscribe({
      next: (person) => {
        this.selectedPerson.set(person);
        this.viewLoading.set(false);
      },
      error: (error: HttpErrorResponse) => {
        this.viewLoading.set(false);
        this.viewError.set(
          error.status === 404
            ? 'ไม่พบข้อมูลรายการนี้'
            : 'โหลดรายละเอียดไม่สำเร็จ กรุณาปิดแล้วลองใหม่'
        );
      }
    });
  }
}