import { useState, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export function PnLChart() {
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('30d');

  // Memoized P&L data generation with stable seed
  const data = useMemo(() => {
    const generatePnLData = (days: number) => {
      const data = [];
      const now = new Date();
      
      for (let i = days; i >= 0; i--) {
        const date = new Date(now);
        date.setDate(date.getDate() - i);
        
        // Generate realistic P&L movement with stable seed based on date
        const baseValue = 1000; // $1000 starting value
        const seed = date.getTime() / (1000 * 60 * 60 * 24); // Day-based seed
        const pseudoRandom = Math.sin(seed) * 10000; // Deterministic "random"
        const dailyVariation = (pseudoRandom - Math.floor(pseudoRandom)) * 200 - 100; // ±$100 variation
        const trendFactor = (days - i) * 3; // Gradual upward trend
        const pnl = baseValue + trendFactor + dailyVariation;
        
        data.push({
          date: date.toISOString().split('T')[0],
          timestamp: date.getTime(),
          realized: Math.max(0, pnl * 0.6 + dailyVariation * 0.3),
          unrealized: Math.max(0, pnl * 0.4 + dailyVariation * 0.7),
          total: pnl,
        });
      }
      
      return data;
    };

    return generatePnLData(timeRange === '7d' ? 7 : timeRange === '30d' ? 30 : 90);
  }, [timeRange]); // Only regenerate when timeRange changes

  const formatDate = (timestamp: number) => {
    const date = new Date(timestamp);
    if (timeRange === '7d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const formatValue = (value: number) => {
    return `$${value.toLocaleString()}`;
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Profit & Loss</CardTitle>
          <div className="flex space-x-2">
            {(['7d', '30d', '90d'] as const).map((range) => (
              <Button
                key={range}
                variant={timeRange === range ? 'default' : 'outline'}
                size="sm"
                onClick={() => setTimeRange(range)}
              >
                {range}
              </Button>
            ))}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data}>
              <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
              <XAxis 
                dataKey="timestamp"
                tickFormatter={formatDate}
                stroke="hsl(var(--muted-foreground))"
                fontSize={12}
              />
              <YAxis 
                tickFormatter={formatValue}
                stroke="hsl(var(--muted-foreground))"
                fontSize={12}
              />
              <Tooltip 
                labelFormatter={(timestamp) => formatDate(timestamp as number)}
                formatter={(value, name) => [
                  formatValue(value as number), 
                  name === 'realized' ? 'Realized P&L' : 
                  name === 'unrealized' ? 'Unrealized P&L' : 'Total P&L'
                ]}
                contentStyle={{
                  backgroundColor: 'hsl(var(--card))',
                  border: '1px solid hsl(var(--border))',
                  borderRadius: '6px',
                }}
              />
              <Line 
                type="monotone" 
                dataKey="realized" 
                stroke="hsl(var(--profit))" 
                strokeWidth={2}
                dot={false}
                name="realized"
              />
              <Line 
                type="monotone" 
                dataKey="unrealized" 
                stroke="hsl(var(--primary))" 
                strokeWidth={2}
                dot={false}
                name="unrealized"
              />
              <Line 
                type="monotone" 
                dataKey="total" 
                stroke="hsl(var(--foreground))" 
                strokeWidth={2}
                dot={false}
                name="total"
                strokeDasharray="5 5"
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
        
        <div className="flex justify-center space-x-6 mt-4 text-sm">
          <div className="flex items-center space-x-2">
            <div className="w-3 h-3 bg-profit rounded-full"></div>
            <span>Realized P&L</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-3 h-3 bg-primary rounded-full"></div>
            <span>Unrealized P&L</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-3 h-3 border-2 border-foreground rounded-full"></div>
            <span>Total P&L</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}