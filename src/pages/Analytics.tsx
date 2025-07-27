import { useState } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { PnLChart } from '@/components/analytics/pnl-chart';
import { AllocationPie } from '@/components/analytics/allocation-pie';
import { AnalyticsOverview } from '@/components/analytics/analytics-overview';
import { PnLDetailsTable } from '@/components/analytics/pnl-details-table';
import { usePnLCalculator } from '@/hooks/use-pnl-calculator';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';
import { DateRange } from 'react-day-picker';
import { AccountingMethod } from '@/lib/pnl-calculator';

export default function Analytics() {
  const { isConnected } = useWalletStore();
  const [accountingMethod, setAccountingMethod] = useState<AccountingMethod>('fifo');
  const [dateRange, setDateRange] = useState<DateRange | undefined>({
    from: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), // 30 days ago
    to: new Date()
  });

  const { calculation } = usePnLCalculator({
    method: accountingMethod,
    dateRange: {
      start: dateRange?.from || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000),
      end: dateRange?.to || new Date()
    }
  });

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Portfolio Analytics</h1>
          <p className="text-muted-foreground mt-2">
            Advanced insights into your DeFi performance and allocations
          </p>
        </div>

        <AnalyticsOverview />
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mt-8">
          <PnLChart />
          <AllocationPie />
        </div>

        {calculation && (
          <div className="mt-8">
            <PnLDetailsTable calculation={calculation} />
          </div>
        )}
      </div>
    </div>
  );
}