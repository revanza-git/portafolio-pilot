import { AlertTriangle, ExternalLink } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Allowance } from '@/stores/portfolio';

interface AllowanceTableProps {
  allowances: Allowance[];
  isLoading: boolean;
  onRevoke: (allowanceId: string) => void;
}

export function AllowanceTable({ allowances, isLoading, onRevoke }: AllowanceTableProps) {
  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const openEtherscan = (address: string) => {
    window.open(`https://etherscan.io/address/${address}`, '_blank');
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Token Allowances</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="flex items-center space-x-4 p-4 border rounded-lg">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-3 w-48" />
                </div>
                <Skeleton className="h-8 w-20" />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Active Allowances</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Token</TableHead>
              <TableHead>Spender</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Last Updated</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {allowances.map((allowance) => (
              <TableRow key={allowance.id}>
                <TableCell>
                  <div className="flex items-center space-x-3">
                    {allowance.token.logoUrl ? (
                      <img 
                        src={allowance.token.logoUrl} 
                        alt={allowance.token.symbol}
                        className="w-8 h-8 rounded-full"
                        onError={(e) => {
                          e.currentTarget.style.display = 'none';
                        }}
                      />
                    ) : (
                      <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                        <span className="text-xs font-medium">
                          {allowance.token.symbol.slice(0, 2)}
                        </span>
                      </div>
                    )}
                    <div>
                      <div className="font-medium">{allowance.token.symbol}</div>
                      <div className="text-sm text-muted-foreground">{allowance.token.name}</div>
                    </div>
                  </div>
                </TableCell>
                
                <TableCell>
                  <div className="flex items-center space-x-2">
                    <div>
                      <div className="font-medium">{allowance.spender.name}</div>
                      <div className="text-xs text-muted-foreground">
                        {allowance.spender.address.slice(0, 6)}...{allowance.spender.address.slice(-4)}
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => openEtherscan(allowance.spender.address)}
                      className="h-auto p-1"
                    >
                      <ExternalLink className="h-3 w-3" />
                    </Button>
                  </div>
                </TableCell>
                
                <TableCell>
                  <div className="flex items-center space-x-2">
                    <span className="font-medium">{allowance.amountFormatted}</span>
                    {allowance.isUnlimited && (
                      <Badge variant="destructive" className="h-5">
                        <AlertTriangle className="h-3 w-3 mr-1" />
                        Unlimited
                      </Badge>
                    )}
                  </div>
                </TableCell>
                
                <TableCell>
                  {formatDate(allowance.lastUpdated)}
                </TableCell>
                
                <TableCell className="text-right">
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => onRevoke(allowance.id)}
                  >
                    Revoke
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        
        {allowances.length === 0 && (
          <div className="text-center py-8">
            <p className="text-muted-foreground">
              No active allowances found.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}