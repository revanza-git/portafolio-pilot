import { Navbar } from '@/components/navigation/navbar';
import { WatchlistGrid } from '@/components/watchlist/watchlist-grid';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Watchlist() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Watchlist</h1>
          <p className="text-muted-foreground mt-2">
            Keep track of your favorite tokens, pools, and protocols
          </p>
        </div>

        <WatchlistGrid />
      </div>
    </div>
  );
}