import { X, AlertTriangle, Info, CheckCircle, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useFeatureFlags } from '@/contexts/feature-flag-context';

const bannerIcons = {
  info: Info,
  warning: AlertTriangle,
  error: AlertCircle,
  success: CheckCircle,
};

const bannerStyles = {
  info: 'border-blue-200 bg-blue-50 text-blue-900 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-100',
  warning: 'border-yellow-200 bg-yellow-50 text-yellow-900 dark:border-yellow-800 dark:bg-yellow-950 dark:text-yellow-100',
  error: 'border-red-200 bg-red-50 text-red-900 dark:border-red-800 dark:bg-red-950 dark:text-red-100',
  success: 'border-green-200 bg-green-50 text-green-900 dark:border-green-800 dark:bg-green-950 dark:text-green-100',
};

export function SystemBanners() {
  const { banners, dismissBanner } = useFeatureFlags();

  if (banners.length === 0) return null;

  return (
    <div className="space-y-2">
      {banners.map((banner) => {
        const Icon = bannerIcons[banner.type];
        
        return (
          <Alert 
            key={banner.id} 
            className={`${bannerStyles[banner.type]} relative`}
          >
            <Icon className="h-4 w-4" />
            <AlertTitle className="flex items-center justify-between">
              {banner.title}
              {banner.dismissible && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0 hover:bg-background/20"
                  onClick={() => dismissBanner(banner.id)}
                  aria-label="Dismiss banner"
                >
                  <X className="h-3 w-3" />
                </Button>
              )}
            </AlertTitle>
            <AlertDescription>{banner.message}</AlertDescription>
          </Alert>
        );
      })}
    </div>
  );
}