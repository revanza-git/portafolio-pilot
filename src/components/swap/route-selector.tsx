import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { TrendingUp, Zap, DollarSign, ArrowRight, AlertTriangle } from 'lucide-react';
import { SwapRoute } from '@/lib/api/swap-client';

interface RouteSelectorProps {
  routes: SwapRoute[];
  selectedRoute: SwapRoute | null;
  onSelectRoute: (route: SwapRoute) => void;
  onExecute: (route: SwapRoute) => void;
  isExecuting?: boolean;
}

export function RouteSelector({ 
  routes, 
  selectedRoute, 
  onSelectRoute, 
  onExecute,
  isExecuting = false 
}: RouteSelectorProps) {
  const [sortBy, setSortBy] = useState<'amount' | 'gas' | 'impact'>('amount');

  const sortedRoutes = [...routes].sort((a, b) => {
    switch (sortBy) {
      case 'amount':
        return parseFloat(b.toAmount) - parseFloat(a.toAmount);
      case 'gas':
        return parseFloat(a.estimatedGas) - parseFloat(b.estimatedGas);
      case 'impact':
        return a.priceImpact - b.priceImpact;
      default:
        return 0;
    }
  });

  const formatAmount = (amount: string, decimals: number = 6) => {
    const value = parseFloat(amount) / Math.pow(10, decimals);
    return value.toFixed(4);
  };

  const formatPercentage = (value: number) => {
    return `${value.toFixed(2)}%`;
  };

  const getPriceImpactColor = (impact: number) => {
    if (impact < 0.1) return 'text-green-600';
    if (impact < 0.5) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="space-y-4">
      {/* Sort Controls */}
      <div className="flex gap-2">
        <Button
          variant={sortBy === 'amount' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('amount')}
        >
          Best Rate
        </Button>
        <Button
          variant={sortBy === 'gas' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('gas')}
        >
          Lowest Gas
        </Button>
        <Button
          variant={sortBy === 'impact' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('impact')}
        >
          Low Impact
        </Button>
      </div>

      {/* Route Cards */}
      <div className="space-y-3">
        {sortedRoutes.map((route) => (
          <Card
            key={route.id}
            className={`cursor-pointer transition-all hover:shadow-md ${
              selectedRoute?.id === route.id ? 'ring-2 ring-primary' : ''
            }`}
            onClick={() => onSelectRoute(route)}
          >
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="capitalize">
                    {route.provider}
                  </Badge>
                  <Badge variant="secondary">
                    {route.dex}
                  </Badge>
                </div>
                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                  <div className="flex items-center gap-1">
                    <Zap className="h-4 w-4" />
                    {route.estimatedGas} gas
                  </div>
                  <div className="flex items-center gap-1">
                    <TrendingUp className="h-4 w-4" />
                    <span className={getPriceImpactColor(route.priceImpact)}>
                      {formatPercentage(route.priceImpact)}
                    </span>
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="text-lg font-semibold">
                  {formatAmount(route.toAmount)} 
                  <span className="text-sm font-normal text-muted-foreground ml-1">
                    received
                  </span>
                </div>
                <ArrowRight className="h-5 w-5 text-muted-foreground" />
              </div>

              {/* Price Impact Warning */}
              {route.priceImpact > 1 && (
                <div className="mt-2 p-2 rounded bg-yellow-50 border border-yellow-200 flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-yellow-600" />
                  <span className="text-sm text-yellow-800">
                    High price impact ({formatPercentage(route.priceImpact)})
                  </span>
                </div>
              )}

              {/* Fee Breakdown */}
              <div className="mt-3 pt-3 border-t border-border/50">
                <div className="flex justify-between text-sm">
                  <span>Protocol Fee:</span>
                  <span>${route.fees.protocolFee}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Gas Fee:</span>
                  <span>${route.fees.gasFee}</span>
                </div>
                <Separator className="my-2" />
                <div className="flex justify-between text-sm font-medium">
                  <span>Total Fees:</span>
                  <span>${route.fees.total}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Execute Button */}
      {selectedRoute && (
        <Button
          className="w-full"
          size="lg"
          onClick={() => onExecute(selectedRoute)}
          disabled={isExecuting}
        >
          {isExecuting ? 'Executing...' : 'Execute Swap'}
        </Button>
      )}
    </div>
  );
}