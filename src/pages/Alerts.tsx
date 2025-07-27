import { Navbar } from '@/components/navigation/navbar';
import { AlertList } from '@/components/alerts/alert-list';
import { CreateAlertButton } from '@/components/alerts/create-alert-button';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Alerts() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold">Price Alerts</h1>
            <p className="text-muted-foreground mt-2">
              Set up notifications for price movements and DeFi events
            </p>
          </div>
          <CreateAlertButton />
        </div>

        <AlertList />
      </div>
    </div>
  );
}