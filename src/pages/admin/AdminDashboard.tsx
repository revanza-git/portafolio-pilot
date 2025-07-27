import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  Users, 
  Wallet, 
  AlertTriangle, 
  Activity,
  TrendingUp,
  Database,
  Clock,
  Shield
} from 'lucide-react';

export default function AdminDashboard() {
  // Mock data - replace with real API calls
  const stats = {
    totalUsers: 1247,
    activeWallets: 892,
    totalErrors: 23,
    apiQuota: { used: 15432, limit: 50000 },
    lastUpdate: new Date().toLocaleString(),
  };

  const recentErrors = [
    { id: 1, message: 'Rate limit exceeded for CoinGecko API', time: '2 mins ago', severity: 'warning' },
    { id: 2, message: 'Failed to fetch price data for USDC', time: '15 mins ago', severity: 'error' },
    { id: 3, message: 'Wallet connection timeout', time: '1 hour ago', severity: 'warning' },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
            <Shield className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h1 className="text-3xl font-bold">Admin Dashboard</h1>
            <p className="text-muted-foreground">System overview and management</p>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Users</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.totalUsers.toLocaleString()}</div>
              <p className="text-xs text-muted-foreground">
                <TrendingUp className="inline h-3 w-3 mr-1" />
                +12% from last month
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Wallets</CardTitle>
              <Wallet className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.activeWallets.toLocaleString()}</div>
              <p className="text-xs text-muted-foreground">
                <Activity className="inline h-3 w-3 mr-1" />
                Last 24 hours
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">System Errors</CardTitle>
              <AlertTriangle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-destructive">{stats.totalErrors}</div>
              <p className="text-xs text-muted-foreground">
                <Clock className="inline h-3 w-3 mr-1" />
                Last 24 hours
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">API Quota</CardTitle>
              <Database className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {((stats.apiQuota.used / stats.apiQuota.limit) * 100).toFixed(1)}%
              </div>
              <p className="text-xs text-muted-foreground">
                {stats.apiQuota.used.toLocaleString()} / {stats.apiQuota.limit.toLocaleString()}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Recent Errors */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5" />
              Recent Errors
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentErrors.map((error) => (
                <div key={error.id} className="flex items-center justify-between p-3 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Badge variant={error.severity === 'error' ? 'destructive' : 'secondary'}>
                      {error.severity}
                    </Badge>
                    <span className="font-medium">{error.message}</span>
                  </div>
                  <span className="text-sm text-muted-foreground">{error.time}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* System Status */}
        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="flex items-center justify-between p-3 rounded-lg border">
                <span className="font-medium">Database</span>
                <Badge variant="default" className="bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100">
                  Healthy
                </Badge>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg border">
                <span className="font-medium">API Services</span>
                <Badge variant="default" className="bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100">
                  Operational
                </Badge>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg border">
                <span className="font-medium">External APIs</span>
                <Badge variant="secondary" className="bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-100">
                  Degraded
                </Badge>
              </div>
            </div>
            <p className="text-sm text-muted-foreground mt-4">
              Last updated: {stats.lastUpdate}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}