import { useState } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { SwapInterface } from '@/components/swap/swap-interface';
import { SwapHistory } from '@/components/swap/swap-history';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Construction } from 'lucide-react';

export default function Swap() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Swap Tokens</h1>
          <p className="text-muted-foreground mt-2">
            Exchange tokens using the best available rates
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Swap Interface */}
          <div className="lg:col-span-2">
            <SwapInterface />
          </div>

          {/* Swap History */}
          <div>
            <SwapHistory />
          </div>
        </div>

        {/* Development Notice */}
        <Card className="mt-8 border-primary/20 bg-primary/5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-primary">
              <Construction className="h-5 w-5" />
              Under Development
            </CardTitle>
            <CardDescription>
              The swap functionality is currently in development. 
              Integration with 0x and 1inch protocols coming soon.
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    </div>
  );
}