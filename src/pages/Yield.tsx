import { Navbar } from '@/components/navigation/navbar';
import { PoolTable } from '@/components/yield/pool-table';
import { YieldOverview } from '@/components/yield/yield-overview';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Yield() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Yield Farming</h1>
          <p className="text-muted-foreground mt-2">
            Discover and manage your yield farming positions across DeFi protocols
          </p>
        </div>

        <YieldOverview />
        <div className="mt-8">
          <PoolTable />
        </div>
      </div>
    </div>
  );
}