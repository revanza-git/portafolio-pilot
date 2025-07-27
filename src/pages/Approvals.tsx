import { useState, useEffect } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { AllowanceTable } from '@/components/approvals/allowance-table';
import { RevokeConfirmDialog } from '@/components/approvals/revoke-confirm-dialog';
import { useWalletStore } from '@/stores/wallet';
import { usePortfolioStore } from '@/stores/portfolio';
import { generateMockAllowances } from '@/lib/mock-data';
import { Navigate } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertTriangle } from 'lucide-react';

export default function Approvals() {
  const { isConnected } = useWalletStore();
  const { 
    allowances, 
    allowancesLoading, 
    setAllowances, 
    setAllowancesLoading 
  } = usePortfolioStore();
  
  const [selectedAllowance, setSelectedAllowance] = useState<string | null>(null);

  useEffect(() => {
    if (isConnected) {
      // TODO: Replace with real API calls
      setAllowancesLoading(true);
      
      setTimeout(() => {
        const mockAllowances = generateMockAllowances();
        setAllowances(mockAllowances);
        setAllowancesLoading(false);
      }, 800);
    }
  }, [isConnected, setAllowances, setAllowancesLoading]);

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  const handleRevoke = (allowanceId: string) => {
    setSelectedAllowance(allowanceId);
  };

  const handleRevokeConfirm = async () => {
    // TODO: Implement actual revoke transaction
    console.log('Revoking allowance:', selectedAllowance);
    setSelectedAllowance(null);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Token Approvals</h1>
          <p className="text-muted-foreground mt-2">
            Manage your token approvals and revoke unnecessary permissions
          </p>
        </div>

        {/* Security Warning */}
        <Card className="mb-6 border-warning/20 bg-warning/5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-warning">
              <AlertTriangle className="h-5 w-5" />
              Security Notice
            </CardTitle>
            <CardDescription>
              Token approvals give smart contracts permission to spend your tokens. 
              Revoke approvals you no longer need to protect your assets.
            </CardDescription>
          </CardHeader>
        </Card>

        <AllowanceTable 
          allowances={allowances}
          isLoading={allowancesLoading}
          onRevoke={handleRevoke}
        />

        <RevokeConfirmDialog 
          isOpen={!!selectedAllowance}
          onClose={() => setSelectedAllowance(null)}
          onConfirm={handleRevokeConfirm}
          allowance={allowances.find(a => a.id === selectedAllowance)}
        />
      </div>
    </div>
  );
}