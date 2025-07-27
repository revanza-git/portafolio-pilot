import { useState } from 'react';
import { Bell, BellOff, Trash2, TrendingUp, TrendingDown } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';

interface Alert {
  id: string;
  type: 'price' | 'apr' | 'allowance';
  token: string;
  condition: 'above' | 'below';
  threshold: number;
  isActive: boolean;
  channel: 'email' | 'telegram';
  createdAt: number;
  lastTriggered?: number;
}

export function AlertList() {
  const [alerts, setAlerts] = useState<Alert[]>([
    {
      id: '1',
      type: 'price',
      token: 'ETH',
      condition: 'above',
      threshold: 2600,
      isActive: true,
      channel: 'email',
      createdAt: Date.now() - 86400000,
    },
    {
      id: '2',
      type: 'price',
      token: 'USDC',
      condition: 'below',
      threshold: 0.99,
      isActive: false,
      channel: 'telegram',
      createdAt: Date.now() - 172800000,
    },
    {
      id: '3',
      type: 'apr',
      token: 'AAVE Pool',
      condition: 'below',
      threshold: 3.0,
      isActive: true,
      channel: 'email',
      createdAt: Date.now() - 259200000,
    },
  ]);

  const toggleAlert = (id: string) => {
    setAlerts(alerts.map(alert => 
      alert.id === id ? { ...alert, isActive: !alert.isActive } : alert
    ));
  };

  const deleteAlert = (id: string) => {
    setAlerts(alerts.filter(alert => alert.id !== id));
  };

  const getAlertIcon = (type: string, condition: string) => {
    if (type === 'price') {
      return condition === 'above' ? 
        <TrendingUp className="h-4 w-4 text-profit" /> : 
        <TrendingDown className="h-4 w-4 text-loss" />;
    }
    return <Bell className="h-4 w-4 text-primary" />;
  };

  const formatThreshold = (alert: Alert) => {
    if (alert.type === 'price') {
      return `$${alert.threshold.toLocaleString()}`;
    } else if (alert.type === 'apr') {
      return `${alert.threshold}%`;
    }
    return alert.threshold.toString();
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div className="space-y-4">
      {alerts.map((alert) => (
        <Card key={alert.id}>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  {getAlertIcon(alert.type, alert.condition)}
                  <div>
                    <div className="font-medium">
                      {alert.token} {alert.condition} {formatThreshold(alert)}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {alert.type === 'price' ? 'Price Alert' : 
                       alert.type === 'apr' ? 'APR Alert' : 'Allowance Alert'} • 
                      Created {formatDate(alert.createdAt)}
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  <Badge variant="outline">
                    {alert.channel}
                  </Badge>
                  {alert.lastTriggered && (
                    <Badge variant="secondary">
                      Last triggered {formatDate(alert.lastTriggered)}
                    </Badge>
                  )}
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <div className="flex items-center space-x-2">
                  {alert.isActive ? (
                    <Bell className="h-4 w-4 text-primary" />
                  ) : (
                    <BellOff className="h-4 w-4 text-muted-foreground" />
                  )}
                  <Switch
                    checked={alert.isActive}
                    onCheckedChange={() => toggleAlert(alert.id)}
                  />
                </div>
                
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => deleteAlert(alert.id)}
                  className="text-destructive hover:text-destructive"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
      
      {alerts.length === 0 && (
        <Card>
          <CardContent className="p-12 text-center">
            <Bell className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-medium mb-2">No alerts yet</h3>
            <p className="text-muted-foreground mb-6">
              Create your first alert to get notified about price movements and DeFi events
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}