import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FiveMinuteBreakOutComponent } from './five-minute-break-out.component';

describe('FiveMinuteBreakOutComponent', () => {
  let component: FiveMinuteBreakOutComponent;
  let fixture: ComponentFixture<FiveMinuteBreakOutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FiveMinuteBreakOutComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FiveMinuteBreakOutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
