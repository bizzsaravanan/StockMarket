import { Component, OnInit, Inject } from '@angular/core';
import { TradeEvaluationService } from '../trade-evaluation.service';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule } from '@angular/material/paginator';
import { HttpClientModule } from '@angular/common/http';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { WsService } from '../ws.service';
import { Subscription } from 'rxjs';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule,
    MatTableModule,
    MatSortModule,
    MatPaginatorModule,
    HttpClientModule,
    MatToolbarModule,
    MatButtonModule,
    CommonModule,
    MatButtonModule,
    FormsModule,
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent implements OnInit {
  displayedColumns: string[] = ['nameSymbol', 'currentStats', 'recommendationStats', 'currentStatus'];

  dataSource: any = [];
  adataSource: any = [];
  private wsSub?: Subscription;
  amount: any = 0
  searchQuery: any = ""

  constructor(private tradeService: TradeEvaluationService, private wsService: WsService, private dialog: MatDialog) { }

  ngOnInit(): void {
    console.log("starting")
  }

  openStartDialog() {
    const dialogRef = this.dialog.open(StartDialogComponent, {
      width: '500px',
      height: '400px',
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Call your API with result.amount and result.cookie
        this.wsSub = this.wsService.getUpdates().subscribe(data => {
          this.dataSource = data.tradeEvaluation;
          this.adataSource = data.atradeEvaluation;
        });
        this.amount = result.amount
        this.tradeService.getEvaluations(result, "Start").subscribe((data: any) => {
          this.dataSource = []
        });
      }
    });
  }

  search() {
    var result = { amount: this.amount, name: this.searchQuery }
    this.tradeService.getEvaluations(result, "Start").subscribe((data: any) => {
      this.searchQuery = ""
    });
  }

  startUpdates() {
    this.tradeService.getEvaluations({}, "Start").subscribe((data: any) => {
      this.dataSource = []
    });
  }

  resetUpdates() {
    this.tradeService.getEvaluations({}, "Reset").subscribe((data: any) => {
      this.searchQuery = ""
      this.dataSource = [];
      this.wsSub?.unsubscribe();
    });
  }
  refresh() {
    this.dataSource = [];
  }

  ngOnDestroy(): void {
    this.wsSub?.unsubscribe();
  }
}

@Component({
  selector: 'app-start-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    HttpClientModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  templateUrl: './start-dialog.component.html',
})
export class StartDialogComponent {
  amount: number = 0;
  cookie: string = '';

  constructor(public dialogRef: MatDialogRef<StartDialogComponent>) { }

  onSubmit(): void {
    const data = { amount: this.amount, cookie: this.cookie };
    this.dialogRef.close(data); // pass data back to parent
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}