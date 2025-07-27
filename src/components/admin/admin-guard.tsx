import { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/auth-context';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Shield, AlertCircle } from 'lucide-react';

interface AdminGuardProps {
  children: ReactNode;
}

export function AdminGuard({ children }: AdminGuardProps) {
  const { isAuthenticated, user, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center p-8">
        <Card className="w-full max-w-md">
          <CardContent className="p-6">
            <div className="animate-pulse space-y-4">
              <div className="h-4 bg-muted rounded w-3/4 mx-auto"></div>
              <div className="h-4 bg-muted rounded w-1/2 mx-auto"></div>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/auth/signin" state={{ from: location }} replace />;
  }

  // Check if user is admin (in a real app, this would check a role/permission)
  const isAdmin = user?.address && isAdminAddress(user.address);

  if (!isAdmin) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center p-8">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="mx-auto w-12 h-12 bg-destructive/10 rounded-full flex items-center justify-center mb-4">
              <AlertCircle className="h-6 w-6 text-destructive" />
            </div>
            <CardTitle className="text-xl">Access Denied</CardTitle>
          </CardHeader>
          <CardContent className="text-center space-y-4">
            <p className="text-muted-foreground">
              You don't have permission to access the admin dashboard.
            </p>
            <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
              <Shield className="h-4 w-4" />
              Admin privileges required
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return <>{children}</>;
}

// Helper function to check if an address is an admin
// In a real app, this would check against a database or smart contract
function isAdminAddress(address: string): boolean {
  const adminAddresses = [
    // Add your admin wallet addresses here
    '0x1234567890123456789012345678901234567890', // Example admin address
  ];
  
  return adminAddresses.includes(address.toLowerCase());
}