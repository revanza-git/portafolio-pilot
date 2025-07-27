import { useState } from 'react';
import { ArrowUpDown, Settings, Info } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';

export function SwapInterface() {
  const [fromToken, setFromToken] = useState('USDC');
  const [toToken, setToToken] = useState('WETH');
  const [fromAmount, setFromAmount] = useState('');
  const [toAmount, setToAmount] = useState('');
  const [slippage, setSlippage] = useState('0.5');

  const tokens = [
    { symbol: 'USDC', name: 'USD Coin', balance: '5,000.00' },
    { symbol: 'WETH', name: 'Wrapped Ether', balance: '1.50' },
    { symbol: 'WBTC', name: 'Wrapped Bitcoin', balance: '0.05' },
    { symbol: 'UNI', name: 'Uniswap', balance: '500.00' },
  ];

  const handleSwapTokens = () => {
    setFromToken(toToken);
    setToToken(fromToken);
    setFromAmount(toAmount);
    setToAmount(fromAmount);
  };

  const handleFromAmountChange = (value: string) => {
    setFromAmount(value);
    // TODO: Calculate toAmount based on current exchange rate
    if (value && !isNaN(Number(value))) {
      // Mock calculation - replace with actual price feed
      const mockRate = fromToken === 'USDC' ? 0.0004 : 2500;
      setToAmount((Number(value) * mockRate).toFixed(6));
    } else {
      setToAmount('');
    }
  };

  const getTokenBalance = (symbol: string) => {
    return tokens.find(t => t.symbol === symbol)?.balance || '0.00';
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Swap Tokens</CardTitle>
          <Button variant="ghost" size="icon">
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* From Token */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <label className="text-sm font-medium">From</label>
            <span className="text-xs text-muted-foreground">
              Balance: {getTokenBalance(fromToken)}
            </span>
          </div>
          <div className="flex space-x-2">
            <div className="flex-1">
              <Input
                placeholder="0.0"
                value={fromAmount}
                onChange={(e) => handleFromAmountChange(e.target.value)}
                className="text-lg h-12"
              />
            </div>
            <Select value={fromToken} onValueChange={setFromToken}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {tokens.map((token) => (
                  <SelectItem key={token.symbol} value={token.symbol}>
                    {token.symbol}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Swap Button */}
        <div className="flex justify-center">
          <Button
            variant="ghost"
            size="icon"
            onClick={handleSwapTokens}
            className="rounded-full border"
          >
            <ArrowUpDown className="h-4 w-4" />
          </Button>
        </div>

        {/* To Token */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <label className="text-sm font-medium">To</label>
            <span className="text-xs text-muted-foreground">
              Balance: {getTokenBalance(toToken)}
            </span>
          </div>
          <div className="flex space-x-2">
            <div className="flex-1">
              <Input
                placeholder="0.0"
                value={toAmount}
                readOnly
                className="text-lg h-12 bg-muted/50"
              />
            </div>
            <Select value={toToken} onValueChange={setToToken}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {tokens.map((token) => (
                  <SelectItem key={token.symbol} value={token.symbol}>
                    {token.symbol}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Swap Details */}
        {fromAmount && toAmount && (
          <>
            <Separator />
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Exchange Rate</span>
                <span>1 {fromToken} = {(Number(toAmount) / Number(fromAmount)).toFixed(6)} {toToken}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Price Impact</span>
                <span className="text-profit">{'<0.01%'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Max Slippage</span>
                <span>{slippage}%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Estimated Gas</span>
                <span>~$15.20</span>
              </div>
            </div>
          </>
        )}

        {/* Swap Button */}
        <Button 
          size="lg" 
          className="w-full bg-gradient-primary hover:opacity-90"
          disabled={!fromAmount || !toAmount}
        >
          {!fromAmount || !toAmount ? 'Enter Amount' : `Swap ${fromToken} for ${toToken}`}
        </Button>

        {/* Development Notice */}
        <div className="flex items-start space-x-2 p-3 bg-primary/5 border border-primary/20 rounded-lg">
          <Info className="h-4 w-4 text-primary mt-0.5 flex-shrink-0" />
          <div className="text-xs">
            <p className="font-medium text-primary mb-1">Development Mode</p>
            <p className="text-muted-foreground">
              This is a UI demonstration. Actual swap execution will be available when 
              0x/1inch integration is completed.
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}