import { TrendingUp, TrendingDown } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import { TokenBalance } from '@/stores/portfolio';

interface TokenTableProps {
  tokens: TokenBalance[];
  isLoading: boolean;
}

export function TokenTable({ tokens, isLoading }: TokenTableProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Token Holdings</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="flex items-center space-x-4">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="space-y-2 flex-1">
                  <Skeleton className="h-4 w-20" />
                  <Skeleton className="h-3 w-32" />
                </div>
                <div className="space-y-2 text-right">
                  <Skeleton className="h-4 w-24" />
                  <Skeleton className="h-3 w-16" />
                </div>
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
        <CardTitle>Token Holdings</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Asset</TableHead>
              <TableHead>Price</TableHead>
              <TableHead>Balance</TableHead>
              <TableHead>Value</TableHead>
              <TableHead className="text-right">24h Change</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {tokens.map((token) => {
              const isPositive = token.change24h >= 0;
              
              return (
                <TableRow key={token.address}>
                  <TableCell className="font-medium">
                    <div className="flex items-center space-x-3">
                      {token.logoUrl ? (
                        <img 
                          src={token.logoUrl} 
                          alt={token.symbol}
                          className="w-8 h-8 rounded-full"
                          onError={(e) => {
                            e.currentTarget.style.display = 'none';
                          }}
                        />
                      ) : (
                        <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                          <span className="text-xs font-medium">
                            {token.symbol.slice(0, 2)}
                          </span>
                        </div>
                      )}
                      <div>
                        <div className="font-medium">{token.symbol}</div>
                        <div className="text-sm text-muted-foreground">{token.name}</div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    ${token.priceUsd.toLocaleString('en-US', { 
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 6 
                    })}
                  </TableCell>
                  <TableCell>{token.balanceFormatted}</TableCell>
                  <TableCell>
                    ${token.usdValue.toLocaleString('en-US', { 
                      minimumFractionDigits: 2 
                    })}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className={`flex items-center justify-end ${
                      isPositive ? 'text-profit' : 'text-loss'
                    }`}>
                      {isPositive ? (
                        <TrendingUp className="h-3 w-3 mr-1" />
                      ) : (
                        <TrendingDown className="h-3 w-3 mr-1" />
                      )}
                      {isPositive ? '+' : ''}{token.change24h.toFixed(2)}%
                    </div>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}