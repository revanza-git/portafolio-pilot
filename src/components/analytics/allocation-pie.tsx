import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';

export function AllocationPie() {
  // Mock allocation data
  const data = [
    { name: 'DeFi Tokens', value: 45, color: 'hsl(var(--primary))' },
    { name: 'Stablecoins', value: 30, color: 'hsl(var(--profit))' },
    { name: 'Blue Chips', value: 20, color: 'hsl(var(--accent))' },
    { name: 'Yield Farming', value: 5, color: 'hsl(var(--warning))' },
  ];

  const formatValue = (value: number) => `${value}%`;

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0];
      return (
        <div className="bg-card border border-border rounded-lg p-3 shadow-lg">
          <p className="font-medium">{data.name}</p>
          <p className="text-sm text-muted-foreground">
            {formatValue(data.value)}
          </p>
        </div>
      );
    }
    return null;
  };

  const CustomLegend = ({ payload }: any) => {
    return (
      <div className="flex flex-wrap justify-center gap-4 mt-4">
        {payload.map((entry: any, index: number) => (
          <div key={index} className="flex items-center space-x-2">
            <div 
              className="w-3 h-3 rounded-full" 
              style={{ backgroundColor: entry.color }}
            />
            <span className="text-sm">{entry.value}</span>
          </div>
        ))}
      </div>
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Portfolio Allocation</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={120}
                paddingAngle={2}
                dataKey="value"
              >
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
              <Legend content={<CustomLegend />} />
            </PieChart>
          </ResponsiveContainer>
        </div>
        
        {/* Allocation Details */}
        <div className="space-y-2 mt-6">
          {data.map((item, index) => (
            <div key={index} className="flex justify-between items-center">
              <div className="flex items-center space-x-2">
                <div 
                  className="w-3 h-3 rounded-full" 
                  style={{ backgroundColor: item.color }}
                />
                <span className="text-sm">{item.name}</span>
              </div>
              <span className="text-sm font-medium">{formatValue(item.value)}</span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}