import { Navbar } from '@/components/navigation/navbar';
import { BridgeForm } from '@/components/bridge/bridge-form';
import { BridgeHistory } from '@/components/bridge/bridge-history';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Bridge() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Cross-Chain Bridge</h1>
          <p className="text-muted-foreground mt-2">
            Transfer assets between different blockchain networks
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2">
            <BridgeForm />
          </div>
          <div>
            <BridgeHistory />
          </div>
        </div>
      </div>
    </div>
  );
}