export interface TradeEvaluation {
    symbol: string;
    currentPrice: number;
    high5Min: number;
    averageVolume: number;
    currentVolume: number;
    movingAvg50: number;
    rsi: number;
    resistanceLevel: number;
    supportLevel: number;
    profitTarget: number;
    stopLoss: number;
    buyScore: number;
    sellScore: number;
    recommendation: string;
    evaluatedAt: string; // ISO string
  }
  