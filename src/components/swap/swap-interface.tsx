import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ArrowUpDown, Settings } from 'lucide-react';
import { RouteSelector } from './route-selector';
import { useSwapRoutes } from '@/hooks/use-swap-routes';
import { SwapRoute } from '@/lib/api/swap-client';

const TOKENS = [
  { address: '0xA0b86991c431E4dFe7bb8E5f2D5E8b8A8A8b3c8B', symbol: 'USDC', name: 'USD Coin' },
  { address: '0xdAC17F958D2ee523a2206206994597C13D831ec7', symbol: 'USDT', name: 'Tether USD' },
  { address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2', symbol: 'WETH', name: 'Wrapped Ether' },
  { address: '0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599', symbol: 'WBTC', name: 'Wrapped Bitcoin' },
  { address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984', symbol: 'UNI', name: 'Uniswap' },
];

export function SwapInterface() {
  const [fromToken, setFromToken] = useState(TOKENS[0].address);
  const [toToken, setToToken] = useState(TOKENS[2].address);
  const [amount, setAmount] = useState('');
  const [slippage, setSlippage] = useState(0.5);
  const [selectedRoute, setSelectedRoute] = useState<SwapRoute | null>(null);
  const [isExecuting, setIsExecuting] = useState(false);

  const { routes, isLoading, fetchQuotes, executeSwap } = useSwapRoutes();

  const handleSearch = async () => {
    if (!amount || !fromToken || !toToken) return;
    
    await fetchQuotes({
      fromToken: fromToken as `0x${string}`,
      toToken: toToken as `0x${string}`,
      fromAmount: (parseFloat(amount) * Math.pow(10, 6)).toString(), // Assuming 6 decimals
      slippage,
    });
  };

  const handleExecute = async (route: SwapRoute) => {
    setIsExecuting(true);
    try {
      await executeSwap(route);
    } finally {
      setIsExecuting(false);
    }
  };

  const swapTokens = () => {
    const tempToken = fromToken;
    setFromToken(toToken);
    setToToken(tempToken);
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Swap Tokens</CardTitle>
          <CardDescription>
            Get the best rates across multiple DEXs
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* From Section */}
          <div className="space-y-2">
            <Label>From</Label>
            <div className="grid grid-cols-2 gap-4">
              <Select value={fromToken} onValueChange={setFromToken}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TOKENS.map((token) => (
                    <SelectItem key={token.address} value={token.address}>
                      {token.symbol}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Input
                type="number"
                placeholder="0.0"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
              />
            </div>
          </div>

          {/* Swap Button */}
          <div className="flex justify-center">
            <Button variant="outline" size="icon" onClick={swapTokens}>
              <ArrowUpDown className="h-4 w-4" />
            </Button>
          </div>

          {/* To Section */}
          <div className="space-y-2">
            <Label>To</Label>
            <Select value={toToken} onValueChange={setToToken}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {TOKENS.map((token) => (
                  <SelectItem key={token.address} value={token.address}>
                    {token.symbol}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Settings */}
          <div className="flex items-center gap-2">
            <Settings className="h-4 w-4" />
            <Label>Slippage Tolerance: {slippage}%</Label>
            <Input
              type="number"
              className="w-20"
              value={slippage}
              onChange={(e) => setSlippage(parseFloat(e.target.value))}
              min="0.1"
              max="5"
              step="0.1"
            />
          </div>

          <Button 
            onClick={handleSearch} 
            disabled={!amount || isLoading || fromToken === toToken}
            className="w-full"
          >
            {isLoading ? 'Finding Best Rate...' : 'Get Quotes'}
          </Button>
        </CardContent>
      </Card>

      {/* Route Selection */}
      {routes.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Best Routes</CardTitle>
            <CardDescription>
              Choose the optimal swap route for your transaction
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RouteSelector
              routes={routes}
              selectedRoute={selectedRoute}
              onSelectRoute={setSelectedRoute}
              onExecute={handleExecute}
              isExecuting={isExecuting}
            />
          </CardContent>
        </Card>
      )}
    </div>
  );
}