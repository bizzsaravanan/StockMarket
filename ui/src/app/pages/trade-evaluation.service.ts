import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, Subject, map, retry, take, takeLast } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TradeEvaluationService {
  private apiUrl = '/sm/api/StockService/'; // Update to your Go backend

  constructor(private http: HttpClient) {}

  getEvaluations(obj: any, method:any): Observable<any> {
    return  this.http.post<any>(this.apiUrl+method, obj);
  }
}
