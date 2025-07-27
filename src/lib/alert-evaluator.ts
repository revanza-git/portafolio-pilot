import { Alert, AlertHistory } from '@/stores/alerts';

export interface AlertEvaluationContext {
  currentPrices: Record<string, number>;
  yieldData: Record<string, { apr: number }>;
  allowanceData: Record<string, { amount: number; isUnlimited: boolean }>;
}

export interface AlertEvaluationResult {
  shouldTrigger: boolean;
  currentValue: number;
  reason?: string;
}

export class AlertEvaluator {
  /**
   * Evaluates a single alert against current market conditions
   */
  static evaluate(alert: Alert, context: AlertEvaluationContext): AlertEvaluationResult {
    try {
      switch (alert.type) {
        case 'price':
          return this.evaluatePriceAlert(alert, context);
        case 'apr':
          return this.evaluateAPRAlert(alert, context);
        case 'allowance':
          return this.evaluateAllowanceAlert(alert, context);
        default:
          return {
            shouldTrigger: false,
            currentValue: 0,
            reason: 'Unknown alert type'
          };
      }
    } catch (error) {
      console.error('Error evaluating alert:', error);
      return {
        shouldTrigger: false,
        currentValue: 0,
        reason: `Evaluation error: ${error instanceof Error ? error.message : 'Unknown error'}`
      };
    }
  }

  private static evaluatePriceAlert(alert: Alert, context: AlertEvaluationContext): AlertEvaluationResult {
    const tokenKey = alert.token.toLowerCase();
    const currentPrice = context.currentPrices[tokenKey];
    
    if (currentPrice === undefined) {
      return {
        shouldTrigger: false,
        currentValue: 0,
        reason: `Price data not available for ${alert.token}`
      };
    }

    const shouldTrigger = alert.condition === 'above' 
      ? currentPrice > alert.threshold
      : currentPrice < alert.threshold;

    return {
      shouldTrigger,
      currentValue: currentPrice,
      reason: shouldTrigger ? 
        `Price ${currentPrice} is ${alert.condition} threshold ${alert.threshold}` : 
        `Price ${currentPrice} does not meet condition`
    };
  }

  private static evaluateAPRAlert(alert: Alert, context: AlertEvaluationContext): AlertEvaluationResult {
    const poolKey = alert.token.toLowerCase();
    const yieldInfo = context.yieldData[poolKey];
    
    if (!yieldInfo) {
      return {
        shouldTrigger: false,
        currentValue: 0,
        reason: `APR data not available for ${alert.token}`
      };
    }

    const currentAPR = yieldInfo.apr;
    const shouldTrigger = alert.condition === 'above' 
      ? currentAPR > alert.threshold
      : currentAPR < alert.threshold;

    return {
      shouldTrigger,
      currentValue: currentAPR,
      reason: shouldTrigger ? 
        `APR ${currentAPR}% is ${alert.condition} threshold ${alert.threshold}%` : 
        `APR ${currentAPR}% does not meet condition`
    };
  }

  private static evaluateAllowanceAlert(alert: Alert, context: AlertEvaluationContext): AlertEvaluationResult {
    const tokenKey = alert.token.toLowerCase();
    const allowanceInfo = context.allowanceData[tokenKey];
    
    if (!allowanceInfo) {
      return {
        shouldTrigger: false,
        currentValue: 0,
        reason: `Allowance data not available for ${alert.token}`
      };
    }

    // For allowance alerts, typically we check if allowance is below a certain threshold
    // or if unlimited allowance is detected
    const currentValue = allowanceInfo.isUnlimited ? Infinity : allowanceInfo.amount;
    const shouldTrigger = alert.condition === 'above' 
      ? currentValue > alert.threshold
      : currentValue < alert.threshold;

    return {
      shouldTrigger,
      currentValue: allowanceInfo.isUnlimited ? Infinity : allowanceInfo.amount,
      reason: shouldTrigger ? 
        `Allowance ${allowanceInfo.isUnlimited ? 'unlimited' : currentValue} is ${alert.condition} threshold ${alert.threshold}` : 
        `Allowance does not meet condition`
    };
  }

  /**
   * Checks if enough time has passed since last trigger (cooldown logic)
   */
  static canTrigger(alert: Alert): boolean {
    if (!alert.lastTriggered || !alert.notificationSettings?.cooldown) {
      return true;
    }

    const cooldownMs = alert.notificationSettings.cooldown * 60 * 1000; // Convert minutes to ms
    const timeSinceLastTrigger = Date.now() - alert.lastTriggered;
    
    return timeSinceLastTrigger >= cooldownMs;
  }

  /**
   * Batch evaluate multiple alerts
   */
  static evaluateAlerts(alerts: Alert[], context: AlertEvaluationContext): Array<{
    alert: Alert;
    result: AlertEvaluationResult;
    canTrigger: boolean;
  }> {
    return alerts
      .filter(alert => alert.isActive)
      .map(alert => ({
        alert,
        result: this.evaluate(alert, context),
        canTrigger: this.canTrigger(alert)
      }));
  }
}

export class NotificationService {
  /**
   * Send email notification
   */
  static async sendEmail(alert: Alert, value: number): Promise<boolean> {
    // In production, this would integrate with your email service
    console.log('Sending email notification:', {
      to: alert.notificationSettings?.email,
      subject: `${alert.token} Alert Triggered`,
      body: `Your alert for ${alert.token} has been triggered. Current value: ${value}, threshold: ${alert.threshold}`
    });

    // Simulate email sending with random success/failure
    await new Promise(resolve => setTimeout(resolve, 1000));
    return Math.random() > 0.1; // 90% success rate
  }

  /**
   * Send Telegram notification
   */
  static async sendTelegram(alert: Alert, value: number): Promise<boolean> {
    const webhookUrl = alert.notificationSettings?.webhookUrl;
    
    if (!webhookUrl) {
      console.error('No webhook URL configured for Telegram alert');
      return false;
    }

    try {
      const message = `🚨 Alert Triggered!\n\n${alert.token} is now ${alert.condition} ${alert.threshold}\nCurrent value: ${value}`;
      
      const response = await fetch(webhookUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          text: message,
          timestamp: Date.now()
        }),
      });

      return response.ok;
    } catch (error) {
      console.error('Error sending Telegram notification:', error);
      return false;
    }
  }

  /**
   * Send notification with retry logic
   */
  static async sendNotification(alert: Alert, value: number, retryCount = 0): Promise<boolean> {
    const maxRetries = alert.notificationSettings?.retryAttempts || 3;
    
    let success = false;
    
    if (alert.channel === 'email') {
      success = await this.sendEmail(alert, value);
    } else if (alert.channel === 'telegram') {
      success = await this.sendTelegram(alert, value);
    }

    // Retry on failure
    if (!success && retryCount < maxRetries) {
      console.log(`Notification failed, retrying... (${retryCount + 1}/${maxRetries})`);
      await new Promise(resolve => setTimeout(resolve, 2000 * (retryCount + 1))); // Exponential backoff
      return this.sendNotification(alert, value, retryCount + 1);
    }

    return success;
  }
}