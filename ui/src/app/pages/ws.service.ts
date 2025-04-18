import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WsService {
  private socket$: WebSocketSubject<any>;

  constructor() {
    this.socket$ = webSocket('ws://localhost:3010/ws');
  }

  getUpdates(): Observable<any> {
    return this.socket$;
  }

  close(): void {
    this.socket$.complete();
  }
}
