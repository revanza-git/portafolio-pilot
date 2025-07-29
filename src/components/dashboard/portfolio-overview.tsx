import { TrendingUp, TrendingDown, DollarSign, Target } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { usePortfolioStore } from '@/stores/portfolio';
import { usePortfolioWithRealPrices } from '@/hooks/use-market-data';

interface PortfolioOverviewProps {
  totalValue?: number;
  change24h?: number;
  isLoading?: boolean;
}

// Safety function to prevent astronomical portfolio values
function safeTotalValue(value: number): number {
  // If value is unrealistically high (over $1 trillion), return 0
  if (value > 1000000000000) {
    console.warn(`Unrealistic portfolio value detected: $${value}. Displaying $0 instead.`);
    return 0;
  }
  return value;
}

export function PortfolioOverview({ 
  totalValue: propTotalValue, 
  change24h: propChange24h, 
  isLoading: propIsLoading 
}: PortfolioOverviewProps) {
  const { tokens } = usePortfolioStore();
  const { data: realData, isLoading: realDataLoading } = usePortfolioWithRealPrices(tokens);
  
  // Use real data if available, fallback to props
  const rawTotalValue = realData?.totalValue ?? propTotalValue ?? 0;
  const totalValue = safeTotalValue(rawTotalValue);
  const change24h = realData?.change24h ?? propChange24h ?? 0;
  const isLoading = realDataLoading || propIsLoading || false;
  const isPositive = change24h >= 0;
  const changePercent = totalValue > 0 ? (change24h / totalValue) * 100 : 0;

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[...Array(3)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-4 w-4" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-32 mb-2" />
              <Skeleton className="h-4 w-24" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {/* Total Portfolio Value */}
      <Card className="bg-gradient-card shadow-card border-0">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Total Portfolio Value
          </CardTitle>
          <DollarSign className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            ${totalValue.toLocaleString('en-US', { minimumFractionDigits: 2 })}
          </div>
          <div className={`text-xs flex items-center mt-1 ${
            isPositive ? 'text-profit' : 'text-loss'
          }`}>
            {isPositive ? (
              <TrendingUp className="h-3 w-3 mr-1" />
            ) : (
              <TrendingDown className="h-3 w-3 mr-1" />
            )}
            {isPositive ? '+' : ''}${change24h.toFixed(2)} ({changePercent.toFixed(2)}%) 24h
          </div>
        </CardContent>
      </Card>

      {/* Token Count */}
      <Card className="bg-gradient-card shadow-card border-0">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Assets
          </CardTitle>
          <Target className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">4</div>
          <p className="text-xs text-muted-foreground">
            Across 1 network
          </p>
        </CardContent>
      </Card>

      {/* Best Performer */}
      <Card className="bg-gradient-card shadow-card border-0">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Best Performer 24h
          </CardTitle>
          <TrendingUp className="h-4 w-4 text-profit" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">UNI</div>
          <p className="text-xs text-profit">
            +5.8% ($7.00)
          </p>
        </CardContent>
      </Card>
    </div>
  );
}