import { useState } from 'react';
import { TrendingUp, TrendingDown, DollarSign, BarChart3, Download, Calendar } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { DateRangePicker } from './date-range-picker';
import { usePnLCalculator } from '@/hooks/use-pnl-calculator';
import { formatCurrency, formatPercent, AccountingMethod } from '@/lib/pnl-calculator';
import { DateRange } from 'react-day-picker';
import { toast } from '@/hooks/use-toast';

export function AnalyticsOverview() {
  const [accountingMethod, setAccountingMethod] = useState<AccountingMethod>('fifo');
  const [dateRange, setDateRange] = useState<DateRange | undefined>({
    from: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), // 30 days ago
    to: new Date()
  });

  const { calculation, isCalculating } = usePnLCalculator({
    method: accountingMethod,
    dateRange: {
      start: dateRange?.from || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000),
      end: dateRange?.to || new Date()
    }
  });

  const handleExportCSV = () => {
    if (!calculation) {
      toast({
        title: "No data to export",
        description: "Please wait for the calculation to complete or check your date range.",
        variant: "destructive"
      });
      return;
    }

    // Generate CSV content
    const headers = ['Date', 'Asset', 'Type', 'Quantity', 'Price', 'Realized P&L'];
    const csvContent = [
      headers.join(','),
      ...calculation.trades.map(trade => [
        new Date(trade.timestamp).toISOString().split('T')[0],
        trade.symbol,
        trade.type.toUpperCase(),
        trade.quantity.toString(),
        trade.price.toString(),
        (trade.realizedPnL || 0).toString()
      ].join(','))
    ].join('\n');

    // Create and download file
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `pnl-analysis-${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);

    toast({
      title: "CSV exported successfully",
      description: "Your P&L analysis has been downloaded."
    });
  };

  return (
    <div className="space-y-6">
      {/* Controls */}
      <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div className="flex flex-wrap items-center gap-4">
          <DateRangePicker
            dateRange={dateRange}
            onDateRangeChange={setDateRange}
          />
          <Select 
            value={accountingMethod} 
            onValueChange={(value: AccountingMethod) => setAccountingMethod(value)}
          >
            <SelectTrigger className="w-32">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="fifo">FIFO</SelectItem>
              <SelectItem value="lifo">LIFO</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button 
          variant="outline" 
          onClick={handleExportCSV}
          disabled={isCalculating || !calculation}
        >
          <Download className="mr-2 h-4 w-4" />
          Export CSV
        </Button>
      </div>

      {/* Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card className="bg-gradient-card shadow-elegant border-0">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Realized P&L</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isCalculating ? (
              <div className="space-y-2">
                <Skeleton className="h-8 w-24" />
                <Skeleton className="h-4 w-20" />
              </div>
            ) : (
              <>
                <div className={`text-2xl font-bold ${(calculation?.realizedPnL || 0) >= 0 ? 'text-profit' : 'text-loss'}`}>
                  {formatCurrency(calculation?.realizedPnL || 0)}
                </div>
                <p className="text-xs text-muted-foreground">
                  From {calculation?.trades.filter(t => t.type === 'sell').length || 0} sell trades
                </p>
              </>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gradient-card shadow-elegant border-0">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Unrealized P&L</CardTitle>
            {(calculation?.unrealizedPnL || 0) >= 0 ? (
              <TrendingUp className="h-4 w-4 text-profit" />
            ) : (
              <TrendingDown className="h-4 w-4 text-loss" />
            )}
          </CardHeader>
          <CardContent>
            {isCalculating ? (
              <div className="space-y-2">
                <Skeleton className="h-8 w-24" />
                <Skeleton className="h-4 w-16" />
              </div>
            ) : (
              <>
                <div className={`text-2xl font-bold ${(calculation?.unrealizedPnL || 0) >= 0 ? 'text-profit' : 'text-loss'}`}>
                  {formatCurrency(calculation?.unrealizedPnL || 0)}
                </div>
                <p className="text-xs text-muted-foreground">
                  Current positions
                </p>
              </>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gradient-card shadow-elegant border-0">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Return</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isCalculating ? (
              <div className="space-y-2">
                <Skeleton className="h-8 w-20" />
                <Skeleton className="h-4 w-24" />
              </div>
            ) : (
              <>
                <div className={`text-2xl font-bold ${(calculation?.totalReturnPercent || 0) >= 0 ? 'text-profit' : 'text-loss'}`}>
                  {formatPercent(calculation?.totalReturnPercent || 0)}
                </div>
                <p className="text-xs text-muted-foreground">
                  {formatCurrency(calculation?.totalReturn || 0)}
                </p>
              </>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gradient-card shadow-elegant border-0">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Lots</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isCalculating ? (
              <div className="space-y-2">
                <Skeleton className="h-8 w-8" />
                <Skeleton className="h-4 w-28" />
              </div>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {calculation?.lots.length || 0}
                </div>
                <p className="text-xs text-muted-foreground">
                  Open positions ({accountingMethod.toUpperCase()})
                </p>
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}