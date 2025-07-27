import { useEffect, useCallback } from 'react';
import { useAlertsStore, Alert } from '@/stores/alerts';
import { useTokenPrices } from '@/hooks/use-market-data';
import { AlertEvaluator, NotificationService, AlertEvaluationContext } from '@/lib/alert-evaluator';
import { useToast } from '@/hooks/use-toast';

export function useAlertEvaluator() {
  const { alerts, addHistoryEntry, updateAlert } = useAlertsStore();
  const { toast } = useToast();

  // Get current token prices for all alerts
  const priceTokens = alerts
    .filter(alert => alert.type === 'price' && alert.isActive)
    .map(alert => alert.token.toLowerCase());
  
  const { data: priceData } = useTokenPrices(priceTokens);

  const evaluateAlerts = useCallback(async () => {
    if (!priceData || alerts.length === 0) return;

    // Build evaluation context
    const context: AlertEvaluationContext = {
      currentPrices: Object.fromEntries(
        Object.entries(priceData).map(([symbol, data]) => [symbol, data.price])
      ),
      // TODO: Add real yield and allowance data
      yieldData: {
        'aave pool': { apr: 3.5 },
        'compound usdc': { apr: 2.8 },
        'uniswap v3': { apr: 12.4 }
      },
      allowanceData: {}
    };

    const evaluationResults = AlertEvaluator.evaluateAlerts(alerts, context);

    for (const { alert, result, canTrigger } of evaluationResults) {
      // Update last evaluated time
      updateAlert(alert.id, { lastEvaluated: Date.now() });

      if (result.shouldTrigger && canTrigger) {
        console.log(`Alert ${alert.id} triggered:`, result);

        // Create history entry as pending
        const historyEntry = {
          alertId: alert.id,
          triggeredAt: Date.now(),
          value: result.currentValue,
          threshold: alert.threshold,
          condition: alert.condition,
          token: alert.token,
          channel: alert.channel,
          status: 'pending' as const,
          retryCount: 0
        };

        addHistoryEntry(historyEntry);

        // Send notification
        try {
          const success = await NotificationService.sendNotification(alert, result.currentValue);
          
          // Update history entry with result
          updateAlert(alert.id, { 
            lastTriggered: Date.now(),
            triggerCount: (alert.triggerCount || 0) + 1
          });

          // Note: In a real app, you'd update the history entry status here
          // For demo purposes, we'll show a toast
          if (success) {
            toast({
              title: "Alert Triggered!",
              description: `${alert.token} ${alert.condition} ${alert.threshold} - Notification sent via ${alert.channel}`,
            });
          } else {
            toast({
              title: "Alert Triggered - Notification Failed",
              description: `${alert.token} alert triggered but notification failed to send`,
              variant: "destructive",
            });
          }
        } catch (error) {
          console.error('Error sending notification:', error);
          toast({
            title: "Notification Error",
            description: "Alert triggered but failed to send notification",
            variant: "destructive",
          });
        }
      }
    }
  }, [alerts, priceData, addHistoryEntry, updateAlert, toast]);

  // Manual evaluation trigger (for testing)
  const triggerEvaluation = useCallback(() => {
    evaluateAlerts();
  }, [evaluateAlerts]);

  // Auto-evaluate alerts periodically (in production, this would be a backend worker)
  useEffect(() => {
    // Only run if we have active alerts and price data
    if (alerts.some(a => a.isActive) && priceData) {
      evaluateAlerts();
    }
  }, [evaluateAlerts, alerts, priceData]);

  // Set up periodic evaluation (every 5 minutes for demo - in production this would be backend)
  useEffect(() => {
    const interval = setInterval(() => {
      if (alerts.some(a => a.isActive)) {
        evaluateAlerts();
      }
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(interval);
  }, [evaluateAlerts, alerts]);

  return {
    triggerEvaluation,
    isEvaluating: false // You could add loading state here
  };
}