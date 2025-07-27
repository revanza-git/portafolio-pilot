import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { 
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { 
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Plus, Edit, Flag, RefreshCw } from 'lucide-react';
import { useFeatureFlags } from '@/contexts/feature-flag-context';
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

export default function FeatureFlags() {
  const [flagsList, setFlagsList] = useState<FeatureFlag[]>([]);
  const [loading, setLoading] = useState(true);
  const [newFlag, setNewFlag] = useState({
    name: '',
    key: '',
    description: '',
    enabled: false,
    rollout_percentage: 100,
  });
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  
  const { refreshFlags } = useFeatureFlags();
  const { toast } = useToast();

  useEffect(() => {
    loadFlags();
  }, []);

  const loadFlags = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/admin/feature-flags');
      if (response.ok) {
        const flags = await response.json();
        setFlagsList(flags);
      } else {
        // Mock data for development
        setFlagsList([
          {
            id: '1',
            name: 'Bridge V2',
            key: 'bridgeV2',
            enabled: false,
            description: 'New bridge interface with improved UX',
            rollout_percentage: 25,
            created_at: '2024-01-15T10:00:00Z',
            updated_at: '2024-01-20T14:30:00Z',
          },
          {
            id: '2',
            name: 'Swap Integration',
            key: 'swapEnabled',
            enabled: true,
            description: 'Enable token swapping functionality',
            rollout_percentage: 100,
            created_at: '2024-01-10T09:00:00Z',
            updated_at: '2024-01-18T16:45:00Z',
          },
          {
            id: '3',
            name: 'Push Notifications',
            key: 'notificationsEnabled',
            enabled: false,
            description: 'Real-time push notifications for price alerts',
            rollout_percentage: 0,
            created_at: '2024-01-12T11:15:00Z',
            updated_at: '2024-01-19T13:20:00Z',
          },
        ]);
      }
    } catch (error) {
      console.error('Failed to load feature flags:', error);
      toast({
        title: 'Error loading flags',
        description: 'Failed to fetch feature flags from server',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const toggleFlag = async (flagId: string, enabled: boolean) => {
    try {
      const response = await fetch(`/api/admin/feature-flags/${flagId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ enabled }),
      });

      if (response.ok) {
        setFlagsList(prev => 
          prev.map(flag => 
            flag.id === flagId 
              ? { ...flag, enabled, updated_at: new Date().toISOString() }
              : flag
          )
        );
        await refreshFlags();
        toast({
          title: 'Flag updated',
          description: `Feature flag ${enabled ? 'enabled' : 'disabled'} successfully`,
        });
      } else {
        throw new Error('Failed to update flag');
      }
    } catch (error) {
      console.error('Failed to toggle flag:', error);
      toast({
        title: 'Error updating flag',
        description: 'Failed to update feature flag',
        variant: 'destructive',
      });
    }
  };

  const createFlag = async () => {
    try {
      const response = await fetch('/api/admin/feature-flags', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newFlag),
      });

      if (response.ok) {
        const createdFlag = await response.json();
        setFlagsList(prev => [...prev, createdFlag]);
        setNewFlag({ name: '', key: '', description: '', enabled: false, rollout_percentage: 100 });
        setIsDialogOpen(false);
        await refreshFlags();
        toast({
          title: 'Flag created',
          description: 'New feature flag created successfully',
        });
      } else {
        throw new Error('Failed to create flag');
      }
    } catch (error) {
      console.error('Failed to create flag:', error);
      toast({
        title: 'Error creating flag',
        description: 'Failed to create new feature flag',
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
              <Flag className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold">Feature Flags</h1>
              <p className="text-muted-foreground">Manage application feature toggles</p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={loadFlags} disabled={loading}>
              <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Flag
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Feature Flag</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="name">Name</Label>
                    <Input
                      id="name"
                      value={newFlag.name}
                      onChange={(e) => setNewFlag(prev => ({ ...prev, name: e.target.value }))}
                      placeholder="Bridge V2"
                    />
                  </div>
                  <div>
                    <Label htmlFor="key">Key</Label>
                    <Input
                      id="key"
                      value={newFlag.key}
                      onChange={(e) => setNewFlag(prev => ({ ...prev, key: e.target.value }))}
                      placeholder="bridgeV2"
                    />
                  </div>
                  <div>
                    <Label htmlFor="description">Description</Label>
                    <Textarea
                      id="description"
                      value={newFlag.description}
                      onChange={(e) => setNewFlag(prev => ({ ...prev, description: e.target.value }))}
                      placeholder="Description of the feature..."
                    />
                  </div>
                  <div className="flex items-center space-x-2">
                    <Switch
                      id="enabled"
                      checked={newFlag.enabled}
                      onCheckedChange={(enabled) => setNewFlag(prev => ({ ...prev, enabled }))}
                    />
                    <Label htmlFor="enabled">Enabled by default</Label>
                  </div>
                  <Button onClick={createFlag} className="w-full">
                    Create Flag
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </div>

        {/* Feature Flags Table */}
        <Card>
          <CardHeader>
            <CardTitle>Active Feature Flags</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Key</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Rollout</TableHead>
                  <TableHead>Updated</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {flagsList.map((flag) => (
                  <TableRow key={flag.id}>
                    <TableCell>
                      <div>
                        <div className="font-medium">{flag.name}</div>
                        {flag.description && (
                          <div className="text-sm text-muted-foreground">{flag.description}</div>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <code className="px-2 py-1 rounded bg-muted text-sm">{flag.key}</code>
                    </TableCell>
                    <TableCell>
                      <Badge variant={flag.enabled ? 'default' : 'secondary'}>
                        {flag.enabled ? 'Enabled' : 'Disabled'}
                      </Badge>
                    </TableCell>
                    <TableCell>{flag.rollout_percentage}%</TableCell>
                    <TableCell>
                      {new Date(flag.updated_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Switch
                          checked={flag.enabled}
                          onCheckedChange={(enabled) => toggleFlag(flag.id, enabled)}
                        />
                        <Button variant="ghost" size="sm">
                          <Edit className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}