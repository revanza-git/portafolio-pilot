import { Navbar } from '@/components/navigation/navbar';
import { PnLChart } from '@/components/analytics/pnl-chart';
import { AllocationPie } from '@/components/analytics/allocation-pie';
import { AnalyticsOverview } from '@/components/analytics/analytics-overview';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Analytics() {
  const { isConnected } = useWalletStore();

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
      </div>
    </div>
  );
}