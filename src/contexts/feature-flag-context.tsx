import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useToast } from '@/hooks/use-toast';

interface FeatureFlag {
  id: string;
  name: string;
  key: string;
  enabled: boolean;
  description?: string;
  rollout_percentage?: number;
  created_at: string;
  updated_at: string;
}

interface SystemBanner {
  id: string;
  title: string;
  message: string;
  type: 'info' | 'warning' | 'error' | 'success';
  active: boolean;
  dismissible: boolean;
  created_at: string;
  expires_at?: string;
}

interface FeatureFlagContextType {
  flags: Record<string, boolean>;
  banners: SystemBanner[];
  isLoading: boolean;
  refreshFlags: () => Promise<void>;
  refreshBanners: () => Promise<void>;
  dismissBanner: (bannerId: string) => void;
}

const FeatureFlagContext = createContext<FeatureFlagContextType | undefined>(undefined);

interface FeatureFlagProviderProps {
  children: ReactNode;
}

export function FeatureFlagProvider({ children }: FeatureFlagProviderProps) {
  const [flags, setFlags] = useState<Record<string, boolean>>({});
  const [banners, setBanners] = useState<SystemBanner[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dismissedBanners, setDismissedBanners] = useState<string[]>([]);
  const { toast } = useToast();

  // Load dismissed banners from localStorage
  useEffect(() => {
    const dismissed = localStorage.getItem('dismissed_banners');
    if (dismissed) {
      setDismissedBanners(JSON.parse(dismissed));
    }
  }, []);

  const refreshFlags = async () => {
    console.log('FeatureFlagContext: Refreshing flags...');
    try {
      const response = await fetch('/api/admin/feature-flags');
      if (response.ok) {
        const flagsData: FeatureFlag[] = await response.json();
        const flagsMap = flagsData.reduce((acc, flag) => {
          acc[flag.key] = flag.enabled;
          return acc;
        }, {} as Record<string, boolean>);
        console.log('FeatureFlagContext: Flags loaded from API:', flagsMap);
        setFlags(flagsMap);
      } else {
        console.log('FeatureFlagContext: API failed, using local config fallback');
        // Fallback to local config if API fails
        const { config } = await import('@/lib/config');
        setFlags(config.features);
      }
    } catch (error) {
      console.error('FeatureFlagContext: Failed to fetch feature flags:', error);
      console.log('FeatureFlagContext: Using local config fallback');
      // Fallback to local config
      const { config } = await import('@/lib/config');
      setFlags(config.features);
    }
  };

  const refreshBanners = async () => {
    try {
      const response = await fetch('/api/admin/banners');
      if (response.ok) {
        const bannersData: SystemBanner[] = await response.json();
        const activeBanners = bannersData.filter(banner => {
          // Filter active banners that haven't expired
          if (!banner.active) return false;
          if (banner.expires_at && new Date(banner.expires_at) < new Date()) return false;
          // Filter out dismissed banners if they're dismissible
          if (banner.dismissible && dismissedBanners.includes(banner.id)) return false;
          return true;
        });
        setBanners(activeBanners);
      }
    } catch (error) {
      console.error('Failed to fetch banners:', error);
    }
  };

  const dismissBanner = (bannerId: string) => {
    const updatedDismissed = [...dismissedBanners, bannerId];
    setDismissedBanners(updatedDismissed);
    localStorage.setItem('dismissed_banners', JSON.stringify(updatedDismissed));
    setBanners(prev => prev.filter(banner => banner.id !== bannerId));
  };

  // Initial load
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      await Promise.all([refreshFlags(), refreshBanners()]);
      setIsLoading(false);
    };
    loadData();
  }, []);

  // Refresh flags and banners every 5 minutes
  useEffect(() => {
    const interval = setInterval(() => {
      refreshFlags();
      refreshBanners();
    }, 5 * 60 * 1000);
    return () => clearInterval(interval);
  }, []);

  const value: FeatureFlagContextType = {
    flags,
    banners,
    isLoading,
    refreshFlags,
    refreshBanners,
    dismissBanner
  };

  return (
    <FeatureFlagContext.Provider value={value}>
      {children}
    </FeatureFlagContext.Provider>
  );
}

export function useFeatureFlag(flagKey: string): boolean {
  const context = useContext(FeatureFlagContext);
  if (context === undefined) {
    throw new Error('useFeatureFlag must be used within a FeatureFlagProvider');
  }
  return context.flags[flagKey] ?? false;
}

export function useFeatureFlags() {
  const context = useContext(FeatureFlagContext);
  if (context === undefined) {
    throw new Error('useFeatureFlags must be used within a FeatureFlagProvider');
  }
  return context;
}