import { create } from 'zustand';

export interface Alert {
  id: string;
  type: 'price' | 'apr' | 'allowance';
  token: string;
  condition: 'above' | 'below';
  threshold: number;
  isActive: boolean;
  channel: 'email' | 'telegram';
  notificationSettings?: {
    cooldown: number; // minutes between notifications
    retryAttempts: number;
    webhookUrl?: string; // For Telegram
    email?: string; // For email notifications
  };
  createdAt: number;
  lastTriggered?: number;
  lastEvaluated?: number;
  triggerCount: number;
}

export interface AlertHistory {
  id: string;
  alertId: string;
  triggeredAt: number;
  value: number;
  threshold: number;
  condition: string;
  token: string;
  channel: string;
  status: 'sent' | 'failed' | 'pending';
  retryCount: number;
  error?: string;
}

interface AlertsState {
  alerts: Alert[];
  history: AlertHistory[];
  isLoading: boolean;
  
  // Actions
  addAlert: (alert: Omit<Alert, 'id' | 'createdAt' | 'triggerCount'>) => void;
  updateAlert: (id: string, updates: Partial<Alert>) => void;
  deleteAlert: (id: string) => void;
  toggleAlert: (id: string) => void;
  
  // History actions
  addHistoryEntry: (entry: Omit<AlertHistory, 'id'>) => void;
  updateHistoryEntry: (id: string, updates: Partial<AlertHistory>) => void;
  
  // Settings
  setLoading: (loading: boolean) => void;
}

export const useAlertsStore = create<AlertsState>((set, get) => ({
  alerts: [
    {
      id: '1',
      type: 'price',
      token: 'ETH',
      condition: 'above',
      threshold: 2600,
      isActive: true,
      channel: 'email',
      notificationSettings: {
        cooldown: 60,
        retryAttempts: 3,
        email: 'user@example.com'
      },
      createdAt: Date.now() - 86400000,
      triggerCount: 0
    },
    {
      id: '2',
      type: 'price',
      token: 'USDC',
      condition: 'below',
      threshold: 0.99,
      isActive: false,
      channel: 'telegram',
      notificationSettings: {
        cooldown: 30,
        retryAttempts: 2,
        webhookUrl: 'https://api.telegram.org/bot...'
      },
      createdAt: Date.now() - 172800000,
      triggerCount: 0
    }
  ],
  history: [
    {
      id: '1',
      alertId: '1',
      triggeredAt: Date.now() - 3600000,
      value: 2650,
      threshold: 2600,
      condition: 'above',
      token: 'ETH',
      channel: 'email',
      status: 'sent',
      retryCount: 0
    }
  ],
  isLoading: false,
  
  addAlert: (alertData) => {
    const alert: Alert = {
      ...alertData,
      id: crypto.randomUUID(),
      createdAt: Date.now(),
      triggerCount: 0
    };
    set(state => ({ alerts: [...state.alerts, alert] }));
  },
  
  updateAlert: (id, updates) => {
    set(state => ({
      alerts: state.alerts.map(alert =>
        alert.id === id ? { ...alert, ...updates } : alert
      )
    }));
  },
  
  deleteAlert: (id) => {
    set(state => ({
      alerts: state.alerts.filter(alert => alert.id !== id),
      history: state.history.filter(entry => entry.alertId !== id)
    }));
  },
  
  toggleAlert: (id) => {
    const { updateAlert } = get();
    const alert = get().alerts.find(a => a.id === id);
    if (alert) {
      updateAlert(id, { isActive: !alert.isActive });
    }
  },
  
  addHistoryEntry: (entryData) => {
    const entry: AlertHistory = {
      ...entryData,
      id: crypto.randomUUID()
    };
    set(state => ({ history: [entry, ...state.history] }));
    
    // Update alert's last triggered time and count
    const { updateAlert } = get();
    updateAlert(entryData.alertId, {
      lastTriggered: entryData.triggeredAt,
      triggerCount: (get().alerts.find(a => a.id === entryData.alertId)?.triggerCount || 0) + 1
    });
  },
  
  updateHistoryEntry: (id, updates) => {
    set(state => ({
      history: state.history.map(entry =>
        entry.id === id ? { ...entry, ...updates } : entry
      )
    }));
  },
  
  setLoading: (isLoading) => set({ isLoading })
}));