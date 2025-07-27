import { Bell, BellOff, Trash2, TrendingUp, TrendingDown, Settings } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { useAlertsStore } from '@/stores/alerts';
import { useAlertEvaluator } from '@/hooks/use-alert-evaluator';

export function AlertList() {
  const { alerts, toggleAlert, deleteAlert } = useAlertsStore();
  const { triggerEvaluation } = useAlertEvaluator();

  const getAlertIcon = (type: string, condition: string) => {
    if (type === 'price') {
      return condition === 'above' ? 
        <TrendingUp className="h-4 w-4 text-profit" /> : 
        <TrendingDown className="h-4 w-4 text-loss" />;
    }
    return <Bell className="h-4 w-4 text-primary" />;
  };

  const formatThreshold = (type: string, threshold: number) => {
    if (type === 'price') {
      return `$${threshold.toLocaleString()}`;
    } else if (type === 'apr') {
      return `${threshold}%`;
    }
    return threshold.toString();
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getCooldownStatus = (alert: any) => {
    if (!alert.lastTriggered || !alert.notificationSettings?.cooldown) {
      return null;
    }

    const cooldownMs = alert.notificationSettings.cooldown * 60 * 1000;
    const timeSinceLastTrigger = Date.now() - alert.lastTriggered;
    
    if (timeSinceLastTrigger < cooldownMs) {
      const remainingMs = cooldownMs - timeSinceLastTrigger;
      const remainingMinutes = Math.ceil(remainingMs / (60 * 1000));
      return `Cooldown: ${remainingMinutes}m`;
    }
    
    return null;
  };

  return (
    <div className="space-y-4">
      {/* Debug Controls */}
      <Card className="bg-muted/20 border-primary/20">
        <CardContent className="p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">
              Alert Evaluator (Demo Mode)
            </div>
            <Button 
              variant="outline" 
              size="sm"
              onClick={triggerEvaluation}
            >
              Test Evaluation
            </Button>
          </div>
        </CardContent>
      </Card>

      {alerts.map((alert) => (
        <Card key={alert.id} className={`${alert.isActive ? 'bg-gradient-card' : 'bg-muted/50'} shadow-card border-0`}>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  {getAlertIcon(alert.type, alert.condition)}
                  <div>
                    <div className="font-medium">
                      {alert.token} {alert.condition} {formatThreshold(alert.type, alert.threshold)}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {alert.type === 'price' ? 'Price Alert' : 
                       alert.type === 'apr' ? 'APR Alert' : 'Allowance Alert'} • 
                      Created {formatDate(alert.createdAt)}
                    </div>
                    {alert.notificationSettings && (
                      <div className="text-xs text-muted-foreground mt-1">
                        Cooldown: {alert.notificationSettings.cooldown}m • 
                        Retries: {alert.notificationSettings.retryAttempts}
                      </div>
                    )}
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  <Badge variant="outline" className="capitalize">
                    {alert.channel}
                  </Badge>
                  {alert.lastTriggered && (
                    <Badge variant="secondary">
                      Last triggered {formatDate(alert.lastTriggered)}
                    </Badge>
                  )}
                  {alert.triggerCount > 0 && (
                    <Badge variant="outline">
                      {alert.triggerCount} triggers
                    </Badge>
                  )}
                  {getCooldownStatus(alert) && (
                    <Badge variant="outline" className="text-warning">
                      {getCooldownStatus(alert)}
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